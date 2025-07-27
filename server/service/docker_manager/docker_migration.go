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
