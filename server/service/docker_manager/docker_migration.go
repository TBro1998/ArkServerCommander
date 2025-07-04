package docker_manager

import (
	"context"
	"fmt"
	"log"

	"ark-server-manager/database"
	"ark-server-manager/models"
	"ark-server-manager/utils"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/client"
)

// InitializeDockerForExistingServers 为现有服务器初始化Docker容器和卷
// 在系统启动时调用，确保所有数据库中的服务器都有对应的Docker资源
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

	log.Printf("开始为 %d 个服务器初始化Docker资源...", len(servers))

	dockerManager, err := GetDockerManager()
	if err != nil {
		return fmt.Errorf("获取Docker管理器失败: %v", err)
	}

	// 批量检查现有资源，减少API调用
	existingVolumes, existingContainers, err := batchCheckDockerResources(dockerManager, servers)
	if err != nil {
		log.Printf("Warning: Failed to batch check Docker resources: %v", err)
	}

	createdVolumes := 0
	createdContainers := 0
	createdConfigs := 0
	errorCount := 0

	for _, server := range servers {
		log.Printf("Checking Docker resources for server %d (%s)...", server.ID, server.Identifier)

		// 检查并创建Docker卷
		volumeName := utils.GetServerVolumeName(server.ID)
		if !existingVolumes[volumeName] {
			// 创建卷
			if _, err := dockerManager.CreateVolume(server.ID); err != nil {
				log.Printf("Warning: Failed to create volume for server %d: %v", server.ID, err)
				errorCount++
			} else {
				log.Printf("Created Docker volume for server %d: %s", server.ID, volumeName)
				createdVolumes++
			}
		}

		// 检查并创建Docker容器
		containerName := utils.GetServerContainerName(server.ID)
		if !existingContainers[containerName] {
			// 创建容器
			containerID, err := dockerManager.CreateContainer(
				server.ID,
				server.Identifier,
				server.Port,
				server.QueryPort,
				server.RCONPort,
				server.AdminPassword,
				server.Map,
				server.GameModIds,
				server.AutoRestart,
			)
			if err != nil {
				log.Printf("Warning: Failed to create container for server %d: %v", server.ID, err)
				errorCount++
			} else {
				log.Printf("Created Docker container for server %d: %s", server.ID, containerID)
				createdContainers++
			}
		}

		// 为没有配置文件的服务器创建默认配置文件
		if err := ensureDefaultConfigFiles(dockerManager, server); err != nil {
			log.Printf("Warning: Failed to ensure config files for server %d: %v", server.ID, err)
		} else {
			createdConfigs++
		}
	}

	log.Printf("Docker initialization completed: %d volumes created, %d containers created, %d config files ensured, %d errors",
		createdVolumes, createdContainers, createdConfigs, errorCount)

	return nil
}

// batchCheckDockerResources 批量检查Docker资源，减少API调用次数
func batchCheckDockerResources(dockerManager *DockerManager, servers []models.Server) (map[string]bool, map[string]bool, error) {
	existingVolumes := make(map[string]bool)
	existingContainers := make(map[string]bool)

	// 获取所有卷
	volumes, err := dockerManager.client.VolumeList(dockerManager.ctx, volume.ListOptions{})
	if err != nil {
		return nil, nil, fmt.Errorf("获取卷列表失败: %v", err)
	}

	// 获取所有容器
	containers, err := dockerManager.client.ContainerList(dockerManager.ctx, container.ListOptions{All: true})
	if err != nil {
		return nil, nil, fmt.Errorf("获取容器列表失败: %v", err)
	}

	// 构建卷名称集合
	for _, volume := range volumes.Volumes {
		existingVolumes[volume.Name] = true
	}

	// 构建容器名称集合
	for _, container := range containers {
		existingContainers[container.Names[0][1:]] = true // 移除开头的 "/"
	}

	return existingVolumes, existingContainers, nil
}

// ensureDefaultConfigFiles 确保服务器有默认配置文件
func ensureDefaultConfigFiles(dockerManager *DockerManager, server models.Server) error {
	// 检查是否已有配置文件
	if _, err := dockerManager.ReadConfigFile(server.ID, utils.GameUserSettingsFileName); err != nil {
		// 文件不存在，创建默认配置
		gameUserSettings := utils.GetDefaultGameUserSettings(
			server.Identifier,
			server.Map,
			server.Port,
			server.QueryPort,
			server.RCONPort,
			70, // 默认最大玩家数
			server.AdminPassword,
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
