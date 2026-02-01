package server

import (
	"ark-server-commander/database"
	"ark-server-commander/models"
	"ark-server-commander/service/docker_manager"
	"ark-server-commander/utils"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// CreateServerWithTransaction 创建服务器（使用数据库事务 + Docker 回滚）
// 改进版：数据库操作使用事务，Docker 资源使用回滚管理器
func (s *ServerService) CreateServerWithTransaction(userID uint, req models.ServerRequest) (*models.ServerResponse, error) {
	var server models.Server
	var err error
	
	// Docker 资源回滚管理器
	dockerRollback := docker_manager.NewRollbackManager()
	
	// 使用 defer 确保 Docker 资源在失败时清理
	defer func() {
		if err != nil && dockerRollback.Count() > 0 {
			utils.Warn("创建失败，清理 Docker 资源", zap.Error(err))
			if rollbackErr := dockerRollback.Rollback(); rollbackErr != nil {
				utils.Error("Docker 资源清理失败", zap.Error(rollbackErr))
			}
		}
	}()
	
	utils.Info("开始创建服务器（事务模式）", 
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
	
	// 步骤3: 数据库事务 - 创建服务器记录
	err = database.DB.Transaction(func(tx *gorm.DB) error {
		// 创建服务器记录
		server = models.Server{
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
				return fmt.Errorf("启动参数格式错误: %w", marshalErr)
			}
			server.ServerArgsJSON = string(argsJson)
		} else {
			server.ServerArgsJSON = "{}"
		}
		
		// 在事务中创建记录
		if createErr := tx.Create(&server).Error; createErr != nil {
			return fmt.Errorf("服务器创建失败: %w", createErr)
		}
		
		utils.Info("服务器记录创建成功（事务中）", zap.Uint("server_id", server.ID))
		return nil
	})
	
	// 数据库事务失败，直接返回（GORM 已自动回滚）
	if err != nil {
		return nil, err
	}
	
	// 继续创建 Docker 资源...
	return s.createDockerResources(userID, &server, req, dockerRollback)
}
