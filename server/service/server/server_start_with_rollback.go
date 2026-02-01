package server

import (
	"ark-server-commander/database"
	"ark-server-commander/models"
	"ark-server-commander/service/docker_manager"
	"ark-server-commander/utils"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
)

// startServerAsyncWithRollback 异步启动服务器（带回滚机制）
func (s *ServerService) startServerAsyncWithRollback(server models.Server, dockerManager *docker_manager.DockerManager, containerName string) error {
	// 创建回滚管理器
	rollback := docker_manager.NewRollbackManager()
	var err error

	// 使用 defer 确保发生错误时执行回滚
	defer func() {
		if err != nil {
			utils.Warn("服务器启动失败，开始执行回滚", zap.Error(err))
			// 更新服务器状态为停止
			database.DB.Model(&server).Update("status", "stopped")

			if rollback.Count() > 0 {
				if rollbackErr := rollback.Rollback(); rollbackErr != nil {
					utils.Error("回滚过程中发生错误", zap.Error(rollbackErr))
				}
			}
		}
	}()

	utils.Info("开始启动服务器（带回滚保护）",
		zap.String("container", containerName),
		zap.Uint("server_id", server.ID))

	// 步骤1: 验证必要镜像是否存在
	missingImages, validateErr := dockerManager.ValidateRequiredImages()
	if validateErr != nil {
		err = fmt.Errorf("验证镜像失败: %w", validateErr)
		return err
	}
	if len(missingImages) > 0 {
		err = fmt.Errorf("无法启动服务器，缺失必要镜像: %v。请手动下载镜像后再启动服务器", missingImages)
		return err
	}

	// 步骤2: 检查容器是否存在
	containerExists, checkErr := dockerManager.ContainerExists(containerName)
	if checkErr != nil {
		err = fmt.Errorf("检查容器是否存在失败: %w", checkErr)
		return err
	}

	needRecreateContainer := false

	// 步骤3: 如果容器存在，检查是否需要重建
	if containerExists {
		envVars, envErr := dockerManager.GetContainerEnvVars(containerName)
		if envErr != nil {
			needRecreateContainer = true
			utils.Info("无法获取容器环境变量，将重建容器", zap.Error(envErr))
		} else {
			// 获取当前服务器的启动参数
			var serverArgs *models.ServerArgs
			if server.ServerArgsJSON != "" && server.ServerArgsJSON != "{}" {
				serverArgs = models.NewServerArgs()
				if unmarshalErr := json.Unmarshal([]byte(server.ServerArgsJSON), serverArgs); unmarshalErr != nil {
					serverArgs = models.FromServer(server)
				}
			} else {
				serverArgs = models.FromServer(server)
			}
			currentArgsString := serverArgs.GenerateArgsString(server)

			// 比较环境变量
			if containerArgsString, exists := envVars["SERVER_ARGS"]; exists {
				if containerArgsString != currentArgsString {
					needRecreateContainer = true
					utils.Info("服务器参数已变更，需要重建容器")
				}
			} else {
				needRecreateContainer = true
			}

			// 检查其他参数
			if !needRecreateContainer {
				if server.GameModIds != envVars["GameModIds"] {
					needRecreateContainer = true
					utils.Info("Mod列表已变更，需要重建容器")
				}
			}
		}

		// 如果需要重建，删除现有容器
		if needRecreateContainer {
			utils.Info("删除现有容器以便重建", zap.String("container", containerName))
			if removeErr := dockerManager.RemoveContainer(containerName); removeErr != nil {
				err = fmt.Errorf("删除现有容器失败: %w", removeErr)
				return err
			}
		}
	}

	// 继续下一部分...
	return s.startServerAsyncContinue(server, dockerManager, containerName, containerExists, needRecreateContainer, rollback)
}
