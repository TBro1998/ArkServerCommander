package utils

import (
	"fmt"
	"log"

	"ark-server-manager/database"
	"ark-server-manager/models"
)

// InitializeDockerForExistingServers 为现有服务器初始化Docker容器和卷
// 在系统启动时调用，确保所有现有服务器都有对应的Docker容器和卷
func InitializeDockerForExistingServers() error {
	// 首先检查Docker环境
	if err := CheckDockerStatus(); err != nil {
		return fmt.Errorf("Docker环境不可用: %v", err)
	}

	var servers []models.Server

	// 获取所有活跃服务器（不包括软删除的）
	if err := database.DB.Find(&servers).Error; err != nil {
		return fmt.Errorf("获取服务器列表失败: %v", err)
	}

	if len(servers) == 0 {
		log.Println("No existing servers found, skipping Docker initialization")
		return nil
	}

	dockerManager, err := NewDockerManager()
	if err != nil {
		return fmt.Errorf("创建Docker管理器失败: %v", err)
	}
	defer dockerManager.Close()

	createdVolumes := 0
	createdContainers := 0
	errorCount := 0

	for _, server := range servers {
		log.Printf("Checking Docker resources for server %d (%s)...", server.ID, server.Identifier)

		// 检查并创建Docker卷
		volumeName := GetServerVolumeName(server.ID)
		if exists, err := dockerManager.VolumeExists(volumeName); err != nil {
			log.Printf("Warning: Failed to check volume existence for server %d: %v", server.ID, err)
			errorCount++
		} else if !exists {
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
		containerName := GetServerContainerName(server.ID)
		if exists, err := dockerManager.ContainerExists(containerName); err != nil {
			log.Printf("Warning: Failed to check container existence for server %d: %v", server.ID, err)
			errorCount++
		} else if !exists {
			// 创建容器
			containerID, err := dockerManager.CreateContainer(
				server.ID,
				server.Identifier,
				server.Port,
				server.QueryPort,
				server.RCONPort,
				server.AdminPassword,
				server.Map,
			)
			if err != nil {
				log.Printf("Warning: Failed to create container for server %d: %v", server.ID, err)
				errorCount++
			} else {
				log.Printf("Created Docker container for server %d: %s", server.ID, containerID)
				createdContainers++

				// 停止容器（创建时会自动启动）
				if err := dockerManager.StopContainer(containerName); err != nil {
					log.Printf("Warning: Failed to stop container after creation for server %d: %v", server.ID, err)
				}
			}
		}

		// 为没有配置文件的服务器创建默认配置文件
		if err := ensureDefaultConfigFiles(dockerManager, server); err != nil {
			log.Printf("Warning: Failed to ensure config files for server %d: %v", server.ID, err)
		}
	}

	log.Printf("Docker initialization completed: %d volumes created, %d containers created, %d errors",
		createdVolumes, createdContainers, errorCount)

	return nil
}

// ensureDefaultConfigFiles 确保服务器有默认配置文件
func ensureDefaultConfigFiles(dockerManager *DockerManager, server models.Server) error {
	// 检查是否已有配置文件
	if _, err := dockerManager.ReadConfigFile(server.ID, GameUserSettingsFileName); err != nil {
		// 文件不存在，创建默认配置
		gameUserSettings := GetDefaultGameUserSettings(
			server.Identifier,
			server.Map,
			server.Port,
			server.QueryPort,
			server.RCONPort,
			70, // 默认最大玩家数
			server.AdminPassword,
		)
		if err := dockerManager.WriteConfigFile(server.ID, GameUserSettingsFileName, gameUserSettings); err != nil {
			return fmt.Errorf("创建默认GameUserSettings.ini失败: %v", err)
		}
		log.Printf("Created default GameUserSettings.ini for server %d", server.ID)
	}

	if _, err := dockerManager.ReadConfigFile(server.ID, GameIniFileName); err != nil {
		// 文件不存在，创建默认配置
		gameIni := GetDefaultGameIni()
		if err := dockerManager.WriteConfigFile(server.ID, GameIniFileName, gameIni); err != nil {
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
		return fmt.Errorf("Docker环境不可用: %v", err)
	}

	var servers []models.Server

	// 获取所有活跃服务器
	if err := database.DB.Find(&servers).Error; err != nil {
		return fmt.Errorf("获取服务器列表失败: %v", err)
	}

	dockerManager, err := NewDockerManager()
	if err != nil {
		return fmt.Errorf("创建Docker管理器失败: %v", err)
	}
	defer dockerManager.Close()

	updatedCount := 0

	for _, server := range servers {
		containerName := GetServerContainerName(server.ID)

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
