package server

import (
	"ark-server-commander/database"
	"ark-server-commander/models"
	"ark-server-commander/service/docker_manager"
	"ark-server-commander/utils"
	"fmt"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// createDockerResources 创建 Docker 资源（卷和配置文件）
func (s *ServerService) createDockerResources(userID uint, server *models.Server, req models.ServerRequest, dockerRollback *docker_manager.RollbackManager) (*models.ServerResponse, error) {
	var err error
	
	// 获取 Docker 管理器
	dockerManager, getErr := docker_manager.GetDockerManager()
	if getErr != nil {
		// Docker 管理器获取失败，需要删除数据库记录
		database.DB.Unscoped().Delete(server)
		return nil, fmt.Errorf("获取Docker管理器失败: %w", getErr)
	}
	
	// 步骤1: 创建 Docker 卷
	utils.Info("创建Docker卷", zap.Uint("server_id", server.ID))
	volumeName, volErr := dockerManager.CreateVolume(server.ID)
	if volErr != nil {
		// 卷创建失败，删除数据库记录
		database.DB.Unscoped().Delete(server)
		return nil, fmt.Errorf("创建Docker卷失败: %w", volErr)
	}
	
	utils.Info("Docker卷创建成功", zap.String("volume", volumeName))
	
	// 注册回滚操作：删除 Docker 卷
	dockerRollback.AddAction("volume", volumeName, "删除Docker卷", func() error {
		return dockerManager.RemoveVolume(volumeName)
	})
	
	// 步骤2: 处理配置文件
	var gameUserSettings string
	var gameIni string
	
	if req.GameUserSettings != "" {
		if validateErr := utils.ValidateINIContent(req.GameUserSettings); validateErr != nil {
			err = fmt.Errorf("GameUserSettings.ini格式错误: %w", validateErr)
			// 触发 Docker 回滚 + 删除数据库记录
			database.DB.Unscoped().Delete(server)
			return nil, err
		}
		gameUserSettings = req.GameUserSettings
	} else {
		gameUserSettings = utils.GetDefaultGameUserSettings(server.Identifier, server.Map, 70)
	}
	
	if req.GameIni != "" {
		if validateErr := utils.ValidateINIContent(req.GameIni); validateErr != nil {
			err = fmt.Errorf("Game.ini格式错误: %w", validateErr)
			database.DB.Unscoped().Delete(server)
			return nil, err
		}
		gameIni = req.GameIni
	} else {
		gameIni = utils.GetDefaultGameIni()
	}
	
	// 步骤3: 写入配置文件
	utils.Info("写入配置文件", zap.Uint("server_id", server.ID))
	if writeErr := dockerManager.WriteConfigFile(server.ID, utils.GameUserSettingsFileName, gameUserSettings); writeErr != nil {
		err = fmt.Errorf("写入GameUserSettings.ini失败: %w", writeErr)
		database.DB.Unscoped().Delete(server)
		return nil, err
	}
	
	if writeErr := dockerManager.WriteConfigFile(server.ID, utils.GameIniFileName, gameIni); writeErr != nil {
		err = fmt.Errorf("写入Game.ini失败: %w", writeErr)
		database.DB.Unscoped().Delete(server)
		return nil, err
	}
	
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
	if content, readErr := dockerManager.ReadConfigFile(server.ID, utils.GameUserSettingsFileName); readErr == nil {
		response.GameUserSettings = content
	}
	if content, readErr := dockerManager.ReadConfigFile(server.ID, utils.GameIniFileName); readErr == nil {
		response.GameIni = content
	}
	
	// 成功完成，清空 Docker 回滚操作
	dockerRollback.Clear()
	
	return &response, nil
}
