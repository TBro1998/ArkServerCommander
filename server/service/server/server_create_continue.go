package server

import (
	"ark-server-commander/models"
	"ark-server-commander/service/docker_manager"
	"ark-server-commander/utils"
	"fmt"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// createServerContinue 继续创建服务器（第二部分：Docker 资源）
func (s *ServerService) createServerContinue(userID uint, server *models.Server, req models.ServerRequest, tx *gorm.DB, rollback *docker_manager.RollbackManager) (*models.ServerResponse, error) {
	var err error

	// 步骤5: 获取 Docker 管理器
	dockerManager, getErr := docker_manager.GetDockerManager()
	if getErr != nil {
		err = fmt.Errorf("获取Docker管理器失败: %w", getErr)
		return nil, err
	}

	// 步骤6: 创建 Docker 卷
	utils.Info("创建Docker卷", zap.Uint("server_id", server.ID))
	volumeName, volErr := dockerManager.CreateVolume(server.ID)
	if volErr != nil {
		err = fmt.Errorf("创建Docker卷失败: %w", volErr)
		return nil, err
	}

	utils.Info("Docker卷创建成功", zap.String("volume", volumeName))

	// 添加回滚操作：删除 Docker 卷
	rollback.AddAction("volume", volumeName, "删除Docker卷", func() error {
		return dockerManager.RemoveVolume(volumeName)
	})

	// 步骤7: 处理配置文件
	var gameUserSettings string
	var gameIni string

	if req.GameUserSettings != "" {
		if validateErr := utils.ValidateINIContent(req.GameUserSettings); validateErr != nil {
			err = fmt.Errorf("GameUserSettings.ini格式错误: %w", validateErr)
			return nil, err
		}
		gameUserSettings = req.GameUserSettings
	} else {
		gameUserSettings = utils.GetDefaultGameUserSettings(server.Identifier, server.Map, 70)
	}

	if req.GameIni != "" {
		if validateErr := utils.ValidateINIContent(req.GameIni); validateErr != nil {
			err = fmt.Errorf("Game.ini格式错误: %w", validateErr)
			return nil, err
		}
		gameIni = req.GameIni
	} else {
		gameIni = utils.GetDefaultGameIni()
	}

	// 步骤8: 写入 GameUserSettings.ini
	utils.Info("写入GameUserSettings.ini", zap.Uint("server_id", server.ID))
	if writeErr := dockerManager.WriteConfigFile(server.ID, utils.GameUserSettingsFileName, gameUserSettings); writeErr != nil {
		err = fmt.Errorf("写入GameUserSettings.ini失败: %w", writeErr)
		return nil, err
	}

	// 添加回滚操作：删除配置文件（通过删除卷来实现）
	// 注意：配置文件的回滚已经包含在卷的回滚中

	// 步骤9: 写入 Game.ini
	utils.Info("写入Game.ini", zap.Uint("server_id", server.ID))
	if writeErr := dockerManager.WriteConfigFile(server.ID, utils.GameIniFileName, gameIni); writeErr != nil {
		err = fmt.Errorf("写入Game.ini失败: %w", writeErr)
		return nil, err
	}

	// 步骤10: 提交数据库事务
	utils.Info("提交数据库事务", zap.Uint("server_id", server.ID))
	if commitErr := tx.Commit().Error; commitErr != nil {
		err = fmt.Errorf("数据库提交失败: %w", commitErr)
		return nil, err
	}

	// 事务提交成功，移除数据库回滚操作
	// 但保留 Docker 资源的回滚操作，以防后续步骤失败

	utils.Info("服务器创建成功",
		zap.Uint("server_id", server.ID),
		zap.String("identifier", server.Identifier))

	// 构建响应
	response := models.ServerResponse{
		ID:            server.ID,
		Identifier:    server.Identifier,
		SessionName:   server.SessionName,
		ClusterID:     server.ClusterID,
		Port:          server.Port,
		QueryPort:     server.QueryPort,
		RCONPort:      server.RCONPort,
		AdminPassword: server.AdminPassword,
		Map:           server.Map,
		GameModIds:    server.GameModIds,
		Status:        server.Status,
		AutoRestart:   server.AutoRestart,
		UserID:        server.UserID,
		CreatedAt:     server.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:     server.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	// 读取配置文件内容
	if gameUserSettingsContent, readErr := dockerManager.ReadConfigFile(server.ID, utils.GameUserSettingsFileName); readErr == nil {
		response.GameUserSettings = gameUserSettingsContent
	}
	if gameIniContent, readErr := dockerManager.ReadConfigFile(server.ID, utils.GameIniFileName); readErr == nil {
		response.GameIni = gameIniContent
	}

	// 成功完成，清空回滚操作
	rollback.Clear()

	return &response, nil
}
