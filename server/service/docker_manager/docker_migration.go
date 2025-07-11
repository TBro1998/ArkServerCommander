package docker_manager

import (
	"context"
	"fmt"
	"log"

	"ark-server-commander/database"
	"ark-server-commander/models"
	"ark-server-commander/utils"

	"github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/client"
)

// InitializeDockerForExistingServers 为现有服务器初始化Docker卷和配置文件
// 在系统启动时调用，确保所有数据库中的服务器都有对应的Docker卷和配置文件
// 注意：不创建容器，容器只在服务器启动时创建
func InitializeDockerForExistingServers() error {
	// 首先检查Docker环境
	if err := CheckDockerStatus(); err != nil {
		return fmt.Errorf("docker环境不可用: %v", err)
	}

	var servers []models.Server

	// 获取所有活跃服务器
	if err := database.DB.Find(&servers).Error; err != nil {
		return fmt.Errorf("获取服务器列表失败: %v", err)
	}

	if len(servers) == 0 {
		log.Println("没有找到需要初始化的服务器")
		return nil
	}

	log.Printf("开始为 %d 个服务器初始化Docker卷和配置文件...", len(servers))

	dockerManager, err := GetDockerManager()
	if err != nil {
		return fmt.Errorf("获取Docker管理器失败: %v", err)
	}

	// 批量检查现有卷，减少API调用
	existingVolumes, err := batchCheckDockerVolumes(dockerManager)
	if err != nil {
		log.Printf("Warning: Failed to batch check Docker volumes: %v", err)
	}

	createdVolumes := 0
	createdConfigs := 0
	errorCount := 0

	for _, server := range servers {
		log.Printf("Checking Docker resources for server %d (%s)...", server.ID, server.Identifier)

		// 检查并创建Docker卷
		volumeName := utils.GetServerVolumeName(server.ID)
		pluginsVolumeName := utils.GetServerPluginsVolumeName(server.ID)

		// 检查游戏数据卷
		if !existingVolumes[volumeName] {
			// 创建卷（包括游戏数据卷和插件卷）
			if _, err := dockerManager.CreateVolume(server.ID); err != nil {
				log.Printf("Warning: Failed to create volumes for server %d: %v", server.ID, err)
				errorCount++
			} else {
				log.Printf("Created Docker volumes for server %d: game=%s, plugins=%s", server.ID, volumeName, pluginsVolumeName)
				createdVolumes += 2 // 创建了两个卷
			}
		} else {
			// 游戏数据卷存在，检查插件卷
			if !existingVolumes[pluginsVolumeName] {
				// 只创建插件卷
				if err := dockerManager.createSingleVolume(pluginsVolumeName); err != nil {
					log.Printf("Warning: Failed to create plugins volume for server %d: %v", server.ID, err)
					errorCount++
				} else {
					log.Printf("Created Docker plugins volume for server %d: %s", server.ID, pluginsVolumeName)
					createdVolumes++
				}
			}
		}

		// 为没有配置文件的服务器创建默认配置文件
		if err := ensureDefaultConfigFiles(dockerManager, server); err != nil {
			log.Printf("Warning: Failed to ensure config files for server %d: %v", server.ID, err)
		} else {
			createdConfigs++
		}
	}

	log.Printf("Docker initialization completed: %d volumes created, %d config files ensured, %d errors",
		createdVolumes, createdConfigs, errorCount)

	return nil
}

// batchCheckDockerVolumes 批量检查Docker卷，减少API调用次数
func batchCheckDockerVolumes(dockerManager *DockerManager) (map[string]bool, error) {
	existingVolumes := make(map[string]bool)

	// 获取所有卷
	volumes, err := dockerManager.client.VolumeList(dockerManager.ctx, volume.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("获取卷列表失败: %v", err)
	}

	// 构建卷名称集合
	for _, volume := range volumes.Volumes {
		existingVolumes[volume.Name] = true
	}

	return existingVolumes, nil
}

// ensureDefaultConfigFiles 确保服务器有默认配置文件
func ensureDefaultConfigFiles(dockerManager *DockerManager, server models.Server) error {
	// 检查是否已有配置文件
	if _, err := dockerManager.ReadConfigFile(server.ID, utils.GameUserSettingsFileName); err != nil {
		// 文件不存在，创建默认配置
		gameUserSettings := utils.GetDefaultGameUserSettings(
			server.Identifier,
			server.Map,
			70, // 默认最大玩家数
		)
		if err := dockerManager.WriteConfigFile(server.ID, utils.GameUserSettingsFileName, gameUserSettings); err != nil {
			return fmt.Errorf("创建默认GameUserSettings.ini失败: %v", err)
		}
		log.Printf("Created default GameUserSettings.ini for server %d", server.ID)
	}

	if _, err := dockerManager.ReadConfigFile(server.ID, utils.GameIniFileName); err != nil {
		// 文件不存在，创建默认配置
		gameIni := utils.GetDefaultGameIni()
		if err := dockerManager.WriteConfigFile(server.ID, utils.GameIniFileName, gameIni); err != nil {
			return fmt.Errorf("创建默认Game.ini失败: %v", err)
		}
		log.Printf("Created default Game.ini for server %d", server.ID)
	}

	return nil
}

// CleanupOrphanedDockerResources 清理孤立的Docker资源
// 删除没有对应数据库记录的Docker容器和卷
func CleanupOrphanedDockerResources() error {
	// 获取所有服务器ID（包括软删除的）
	var serverIDs []uint
	if err := database.DB.Unscoped().Model(&models.Server{}).Pluck("id", &serverIDs).Error; err != nil {
		return fmt.Errorf("获取服务器ID列表失败: %v", err)
	}

	// 创建ID集合用于快速查找
	serverIDMap := make(map[uint]bool)
	for _, id := range serverIDs {
		serverIDMap[id] = true
	}

	// dockerManager := NewDockerManager()

	// TODO: 实现清理逻辑
	// 这里需要列出所有ark-server-*的容器和卷，然后检查对应的服务器ID是否存在
	// 暂时不实现，避免意外删除重要数据

	log.Println("Docker cleanup completed (not implemented for safety)")
	return nil
}

// SyncServerStatusWithDocker 同步服务器状态与Docker容器状态
func SyncServerStatusWithDocker() error {
	// 首先检查Docker环境
	if err := CheckDockerStatus(); err != nil {
		return fmt.Errorf("docker环境不可用: %v", err)
	}

	var servers []models.Server

	// 获取所有活跃服务器
	if err := database.DB.Find(&servers).Error; err != nil {
		return fmt.Errorf("获取服务器列表失败: %v", err)
	}

	dockerManager, err := GetDockerManager()
	if err != nil {
		return fmt.Errorf("获取Docker管理器失败: %v", err)
	}

	updatedCount := 0

	for _, server := range servers {
		containerName := utils.GetServerContainerName(server.ID)

		// 获取Docker容器状态
		dockerStatus, err := dockerManager.GetContainerStatus(containerName)
		if err != nil {
			log.Printf("Warning: Failed to get container status for server %d: %v", server.ID, err)
			continue
		}

		// 处理容器不存在的情况
		if dockerStatus == "not_found" {
			// 如果容器不存在但数据库状态是运行中，更新为停止状态
			if server.Status == "running" || server.Status == "starting" {
				if err := database.DB.Model(&server).Update("status", "stopped").Error; err != nil {
					log.Printf("Warning: Failed to update status for server %d: %v", server.ID, err)
				} else {
					log.Printf("Updated server %d status from %s to stopped (container not found)", server.ID, server.Status)
					updatedCount++
				}
			}
			continue
		}

		// 如果状态不同，更新数据库
		if dockerStatus != server.Status {
			if err := database.DB.Model(&server).Update("status", dockerStatus).Error; err != nil {
				log.Printf("Warning: Failed to update status for server %d: %v", server.ID, err)
			} else {
				log.Printf("Updated server %d status from %s to %s", server.ID, server.Status, dockerStatus)
				updatedCount++
			}
		}
	}

	if updatedCount > 0 {
		log.Printf("Synchronized status for %d servers", updatedCount)
	}

	return nil
}

// CheckDockerStatus 检查Docker环境状态
// 返回: Docker是否可用和错误信息
func CheckDockerStatus() error {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return fmt.Errorf("docker客户端创建失败: %v", err)
	}
	defer cli.Close()

	// 检查Docker服务是否运行
	_, err = cli.Ping(context.Background())
	if err != nil {
		return fmt.Errorf("docker服务未运行或无法连接: %v", err)
	}

	return nil
}
