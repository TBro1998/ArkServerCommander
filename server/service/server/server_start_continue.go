package server

import (
	"ark-server-commander/database"
	"ark-server-commander/models"
	"ark-server-commander/service/docker_manager"
	"ark-server-commander/utils"
	"fmt"
	"go.uber.org/zap"
	"time"
)

// startServerAsyncContinue 继续启动服务器（第二部分）
func (s *ServerService) startServerAsyncContinue(server models.Server, dockerManager *docker_manager.DockerManager, containerName string, containerExists, needRecreateContainer bool, rollback *docker_manager.RollbackManager) error {
	var err error
	
	// 步骤4: 创建容器（如果需要）
	if !containerExists || needRecreateContainer {
		utils.Info("创建新容器", zap.String("container", containerName))
		
		// 使用带回滚的容器创建方法
		containerID, createErr := dockerManager.CreateContainerWithRollback(
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
		
		if createErr != nil {
			err = fmt.Errorf("创建容器失败: %w", createErr)
			return err
		}
		
		utils.Info("容器创建成功", 
			zap.String("container_id", containerID),
			zap.String("container_name", containerName))
		
		// 添加回滚操作：删除刚创建的容器
		rollback.AddAction("container", containerName, "删除容器", func() error {
			return dockerManager.RemoveContainer(containerName)
		})
	}
	
	// 步骤5: 启动容器
	utils.Info("启动容器", zap.String("container", containerName))
	if startErr := dockerManager.StartContainer(containerName); startErr != nil {
		err = fmt.Errorf("启动容器失败: %w", startErr)
		return err
	}
	
	// 添加回滚操作：停止容器
	rollback.AddAction("container", containerName, "停止容器", func() error {
		return dockerManager.StopContainer(containerName)
	})
	
	// 步骤6: 等待容器启动
	utils.Info("等待容器启动", zap.String("container", containerName))
	for i := 0; i < 30; i++ {
		time.Sleep(1 * time.Second)
		status, statusErr := dockerManager.GetContainerStatus(containerName)
		if statusErr != nil {
			utils.Debug("获取容器状态失败，继续等待", zap.Error(statusErr))
			continue
		}
		
		if status == "running" {
			utils.Info("容器启动成功", 
				zap.String("container", containerName),
				zap.Int("wait_seconds", i+1))
			
			// 更新数据库状态
			if updateErr := database.DB.Model(&server).Update("status", "running").Error; updateErr != nil {
				utils.Error("更新服务器状态为running失败", zap.Error(updateErr))
			}
			
			// 成功完成，清空回滚操作
			rollback.Clear()
			return nil
		}
	}
	
	err = fmt.Errorf("容器启动超时")
	return err
}
