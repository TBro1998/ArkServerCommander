package utils

import (
	"fmt"
	"log"
	"os"

	"ark-server-manager/database"
	"ark-server-manager/models"
)

// InitializeExistingServerFolders 为现有服务器初始化文件夹
// 在系统启动时调用，确保所有现有服务器都有对应的文件夹
func InitializeExistingServerFolders() error {
	var servers []models.Server

	// 获取所有服务器（包括软删除的）
	if err := database.DB.Unscoped().Find(&servers).Error; err != nil {
		return fmt.Errorf("获取服务器列表失败: %v", err)
	}

	createdCount := 0
	errorCount := 0

	for _, server := range servers {
		// 获取服务器文件夹路径
		folderPath := GetServerFolderPath(server.ID)

		// 检查文件夹是否存在
		if !FolderExists(folderPath) {
			// 如果服务器已被软删除，不创建文件夹
			if server.DeletedAt.Valid {
				continue
			}

			// 创建文件夹
			_, err := CreateServerFolder(server.ID)
			if err != nil {
				log.Printf("Warning: Failed to create folder for server %d (%s): %v", server.ID, server.Name, err)
				errorCount++
			} else {
				log.Printf("Created folder for existing server %d (%s): %s", server.ID, server.Name, folderPath)
				createdCount++
			}
		}
	}

	if createdCount > 0 {
		log.Printf("Initialized %d server folders for existing servers", createdCount)
	}
	if errorCount > 0 {
		log.Printf("Warning: Failed to create %d server folders", errorCount)
	}

	return nil
}

// CleanupOrphanedFolders 清理孤立的服务器文件夹
// 删除没有对应数据库记录的服务器文件夹
func CleanupOrphanedFolders() error {
	// 确保Works目录存在
	_, err := EnsureWorksDirectory()
	if err != nil {
		return fmt.Errorf("获取Works目录失败: %v", err)
	}

	// 获取所有服务器ID（包括软删除的）
	var serverIDs []uint
	if err := database.DB.Unscoped().Model(&models.Server{}).Pluck("id", &serverIDs).Error; err != nil {
		return fmt.Errorf("获取服务器ID列表失败: %v", err)
	}

	// 创建ID集合用于快速查找
	serverIDMap := make(map[string]bool)
	for _, id := range serverIDs {
		serverIDMap[fmt.Sprintf("%d", id)] = true
	}

	// TODO: 扫描Works目录下的文件夹，删除不在数据库中的文件夹
	// 这里需要实现目录扫描和清理逻辑
	// 暂时不实现，避免意外删除重要数据

	return nil
}

// MigrateServerConfigs 为现有服务器生成默认的INI配置文件
// 在系统启动时调用，为没有配置文件的服务器创建默认配置文件
func MigrateServerConfigs() error {
	var servers []models.Server

	// 获取所有服务器
	if err := database.DB.Find(&servers).Error; err != nil {
		return fmt.Errorf("获取服务器列表失败: %v", err)
	}

	createdCount := 0

	for _, server := range servers {
		// 检查配置文件是否存在
		gameUserSettingsPath := GetConfigFilePath(server.ID, GameUserSettingsFileName)
		gameIniPath := GetConfigFilePath(server.ID, GameIniFileName)

		gameUserSettingsExists := false
		gameIniExists := false

		if _, err := os.Stat(gameUserSettingsPath); err == nil {
			gameUserSettingsExists = true
		}

		if _, err := os.Stat(gameIniPath); err == nil {
			gameIniExists = true
		}

		// 如果配置文件不存在，创建默认配置
		if !gameUserSettingsExists || !gameIniExists {
			if err := CreateDefaultConfigFiles(
				server.ID,
				server.Name,
				server.Map,
				server.Port,
				server.QueryPort,
				server.RCONPort,
				server.MaxPlayers,
				server.AdminPassword,
			); err != nil {
				log.Printf("Warning: Failed to create config files for server %d (%s): %v", server.ID, server.Name, err)
			} else {
				log.Printf("Created default config files for server %d (%s)", server.ID, server.Name)
				createdCount++
			}
		}
	}

	if createdCount > 0 {
		log.Printf("Created default configuration files for %d servers", createdCount)
	}

	return nil
}
