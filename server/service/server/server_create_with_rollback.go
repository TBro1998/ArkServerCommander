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

// CreateServerWithRollback 创建服务器（带完整回滚机制）
// 如果任何步骤失败，会自动清理已创建的所有资源
func (s *ServerService) CreateServerWithRollback(userID uint, req models.ServerRequest) (*models.ServerResponse, error) {
	// 创建回滚管理器
	rollback := docker_manager.NewRollbackManager()
	var err error

	// 使用 defer 确保发生错误时执行回滚
	defer func() {
		if err != nil && rollback.Count() > 0 {
			utils.Warn("服务器创建失败，开始执行回滚", zap.Error(err))
			if rollbackErr := rollback.Rollback(); rollbackErr != nil {
				utils.Error("回滚过程中发生错误", zap.Error(rollbackErr))
			}
		}
	}()

	utils.Info("开始创建服务器（带回滚保护）",
		zap.String("identifier", req.Identifier),
		zap.Uint("user_id", userID))

	// 步骤1: 检查服务器标识是否已存在
	var existingServer models.Server
	if checkErr := database.DB.Where("identifier = ? AND user_id = ?", req.Identifier, userID).First(&existingServer).Error; checkErr == nil {
		err = fmt.Errorf("服务器标识已存在")
		return nil, err
	}

	// 步骤2: 设置默认值
	if req.Map == "" {
		req.Map = "TheIsland"
	}
	if req.MaxPlayers == 0 {
		req.MaxPlayers = 70
	}
	if req.AutoRestart == nil {
		defaultVal := true
		req.AutoRestart = &defaultVal
	}

	// 步骤3: 开始数据库事务
	tx := database.DB.Begin()
	if tx.Error != nil {
		err = fmt.Errorf("数据库事务启动失败: %w", tx.Error)
		return nil, err
	}

	// 添加回滚操作：回滚数据库事务
	rollback.AddAction("database", "transaction", "回滚数据库事务", func() error {
		tx.Rollback()
		return nil
	})

	// 步骤4: 创建服务器记录
	server := models.Server{
		Identifier:    req.Identifier,
		SessionName:   req.SessionName,
		ClusterID:     req.ClusterID,
		Port:          req.Port,
		QueryPort:     req.QueryPort,
		RCONPort:      req.RCONPort,
		AdminPassword: req.AdminPassword,
		Map:           req.Map,
		MaxPlayers:    req.MaxPlayers,
		GameModIds:    req.GameModIds,
		Status:        "stopped",
		AutoRestart:   *req.AutoRestart,
		UserID:        userID,
	}

	if req.ServerArgs != nil {
		argsJson, marshalErr := json.Marshal(req.ServerArgs)
		if marshalErr != nil {
			err = fmt.Errorf("启动参数格式错误: %w", marshalErr)
			return nil, err
		}
		server.ServerArgsJSON = string(argsJson)
	} else {
		server.ServerArgsJSON = "{}"
	}

	if createErr := tx.Create(&server).Error; createErr != nil {
		err = fmt.Errorf("服务器创建失败: %w", createErr)
		return nil, err
	}

	utils.Info("服务器记录创建成功", zap.Uint("server_id", server.ID))

	// 添加回滚操作：删除服务器记录
	serverID := server.ID
	rollback.AddAction("database", fmt.Sprintf("server_%d", serverID), "删除服务器记录", func() error {
		return database.DB.Unscoped().Delete(&models.Server{}, serverID).Error
	})

	// 继续下一部分...
	return s.createServerContinue(userID, &server, req, tx, rollback)
}
