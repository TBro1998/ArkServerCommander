package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"ark-server-commander/database"
	"ark-server-commander/models"
	"ark-server-commander/service/docker_manager"
	"ark-server-commander/utils"

	"github.com/gin-gonic/gin"
)

// GetServers 获取服务器列表
// @Summary 获取服务器列表
// @Description 获取当前用户的所有服务器列表
// @Tags 服务器管理
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} map[string][]models.ServerResponse "服务器列表"
// @Failure 401 {object} map[string]string "未授权"
// @Failure 500 {object} map[string]string "服务器错误"
// @Router /servers [get]
func GetServers(c *gin.Context) {
	userID := c.GetUint("user_id")

	var servers []models.Server
	if err := database.DB.Where("user_id = ?", userID).Find(&servers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取服务器列表失败"})
		return
	}

	// 转换为响应格式并获取实时状态
	dockerManager, err := docker_manager.GetDockerManager()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取Docker管理器失败"})
		return
	}

	var serverResponses []models.ServerResponse
	for _, server := range servers {
		// 获取Docker容器实时状态
		containerName := utils.GetServerContainerName(server.ID)
		realTimeStatus := server.Status // 默认使用数据库状态

		// 检查容器是否存在
		containerExists, err := dockerManager.ContainerExists(containerName)
		if err == nil && containerExists {
			if dockerStatus, err := dockerManager.GetContainerStatus(containerName); err == nil {
				// 如果能获取到Docker状态，使用实时状态
				realTimeStatus = dockerStatus

				// 如果实时状态与数据库状态不同，更新数据库（异步）
				if realTimeStatus != server.Status {
					go func(s models.Server, status string) {
						database.DB.Model(&s).Update("status", status)
					}(server, realTimeStatus)
				}
			}
		} else if err == nil && !containerExists && server.Status == "running" {
			// 如果容器不存在但数据库状态是运行中，更新为停止状态
			realTimeStatus = "stopped"
			go func(s models.Server) {
				database.DB.Model(&s).Update("status", "stopped")
			}(server)
		}

		serverResponses = append(serverResponses, models.ServerResponse{
			ID:            server.ID,
			Identifier:    server.Identifier,
			SessionName:   server.SessionName,
			ClusterID:     server.ClusterID,
			Port:          server.Port,
			QueryPort:     server.QueryPort,
			RCONPort:      server.RCONPort,
			AdminPassword: server.AdminPassword,
			Map:           server.Map,
			MaxPlayers:    server.MaxPlayers,
			GameModIds:    server.GameModIds,
			Status:        realTimeStatus,
			AutoRestart:   server.AutoRestart,
			UserID:        server.UserID,
			CreatedAt:     server.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:     server.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "获取成功",
		"data":    serverResponses,
	})
}

// CreateServer 创建服务器
// @Summary 创建新服务器
// @Description 创建一个新的ARK服务器配置
// @Tags 服务器管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param server body models.ServerRequest true "服务器配置"
// @Success 201 {object} map[string]models.ServerResponse "创建成功"
// @Failure 400 {object} map[string]string "请求错误"
// @Failure 401 {object} map[string]string "未授权"
// @Failure 500 {object} map[string]string "服务器错误"
// @Router /servers [post]
func CreateServer(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req models.ServerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	// 检查服务器标识是否已存在
	var existingServer models.Server
	if err := database.DB.Where("identifier = ? AND user_id = ?", req.Identifier, userID).First(&existingServer).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "服务器标识已存在"})
		return
	}

	// 设置默认值
	if req.Map == "" {
		req.Map = "TheIsland"
	}
	// 设置默认最大玩家数
	if req.MaxPlayers == 0 {
		req.MaxPlayers = 70
	}
	// 设置默认自动重启值
	if req.AutoRestart == nil {
		defaultVal := true
		req.AutoRestart = &defaultVal
	}

	// 开始数据库事务
	tx := database.DB.Begin()
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "数据库事务启动失败"})
		return
	}

	// 确保在函数结束时处理事务
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器创建过程中发生异常"})
		}
	}()

	// 创建服务器
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
		argsJson, err := json.Marshal(req.ServerArgs)
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": "启动参数格式错误"})
			return
		}
		server.ServerArgsJSON = string(argsJson)
	} else {
		server.ServerArgsJSON = "{}"
	}

	if err := tx.Create(&server).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器创建失败"})
		return
	}

	// 检查Docker环境
	if err := docker_manager.CheckDockerStatus(); err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Docker环境检查失败: %v", err)})
		return
	}

	// 创建Docker卷（不创建容器）
	dockerManager, err := docker_manager.GetDockerManager()
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取Docker管理器失败"})
		return
	}

	var volumeName string

	// 创建Docker卷
	volumeName, err = dockerManager.CreateVolume(server.ID)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("创建Docker卷失败: %v", err)})
		return
	}
	fmt.Printf("Created Docker volume: %s\n", volumeName)

	// 写入配置文件到Docker卷中
	// 如果请求中包含配置文件，使用请求中的配置；否则使用默认配置
	var gameUserSettings string
	var gameIni string

	if req.GameUserSettings != "" {
		// 验证配置文件格式
		if err := utils.ValidateINIContent(req.GameUserSettings); err != nil {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("GameUserSettings.ini格式错误: %v", err)})
			return
		}
		gameUserSettings = req.GameUserSettings
	} else {
		gameUserSettings = utils.GetDefaultGameUserSettings(server.Identifier, server.Map, 70)
	}

	if req.GameIni != "" {
		// 验证配置文件格式
		if err := utils.ValidateINIContent(req.GameIni); err != nil {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Game.ini格式错误: %v", err)})
			return
		}
		gameIni = req.GameIni
	} else {
		gameIni = utils.GetDefaultGameIni()
	}

	// 写入配置文件
	if err := dockerManager.WriteConfigFile(server.ID, utils.GameUserSettingsFileName, gameUserSettings); err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("写入GameUserSettings.ini失败: %v", err)})
		return
	}

	if err := dockerManager.WriteConfigFile(server.ID, utils.GameIniFileName, gameIni); err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("写入Game.ini失败: %v", err)})
		return
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		// 提交失败，清理Docker资源
		dockerManager.RemoveVolume(volumeName)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "数据库提交失败"})
		return
	}

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

	// 从Docker卷读取配置文件内容并添加到响应中
	dockerManager2, err := docker_manager.GetDockerManager()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取Docker管理器失败"})
		return
	}
	if gameUserSettings, err := dockerManager2.ReadConfigFile(uint(server.ID), utils.GameUserSettingsFileName); err == nil {
		response.GameUserSettings = gameUserSettings
	}
	if gameIni, err := dockerManager2.ReadConfigFile(uint(server.ID), utils.GameIniFileName); err == nil {
		response.GameIni = gameIni
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "服务器创建成功",
		"data":    response,
	})
}

// GetServer 获取单个服务器信息
// @Summary 获取服务器详情
// @Description 根据ID获取指定服务器的详细信息（包括配置文件内容）
// @Tags 服务器管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "服务器ID"
// @Success 200 {object} map[string]models.ServerResponse "服务器信息（包含配置文件）"
// @Failure 400 {object} map[string]string "请求错误"
// @Failure 404 {object} map[string]string "服务器不存在"
// @Failure 401 {object} map[string]string "未授权"
// @Router /servers/{id} [get]
func GetServer(c *gin.Context) {
	userID := c.GetUint("user_id")
	serverID := c.Param("id")

	id, err := strconv.ParseUint(serverID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的服务器ID"})
		return
	}

	var server models.Server
	if err := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&server).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "服务器不存在"})
		return
	}

	// 解析启动参数
	var serverArgs *models.ServerArgs
	if server.ServerArgsJSON != "" && server.ServerArgsJSON != "{}" {
		serverArgs = models.NewServerArgs()
		if err := json.Unmarshal([]byte(server.ServerArgsJSON), serverArgs); err != nil {
			// 如果解析失败，使用默认参数
			serverArgs = models.FromServer(server)
		}
	} else {
		serverArgs = models.FromServer(server)
	}

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
		MaxPlayers:    server.MaxPlayers,
		GameModIds:    server.GameModIds,
		Status:        server.Status,
		AutoRestart:   server.AutoRestart,
		UserID:        server.UserID,
		CreatedAt:     server.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:     server.UpdatedAt.Format("2006-01-02 15:04:05"),
		ServerArgs:    serverArgs,
		GeneratedArgs: serverArgs.GenerateArgsString(server),
	}

	// 从Docker卷读取配置文件内容（如果存在）
	dockerManager, err := docker_manager.GetDockerManager()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取Docker管理器失败"})
		return
	}
	if gameUserSettings, err := dockerManager.ReadConfigFile(uint(id), utils.GameUserSettingsFileName); err == nil {
		response.GameUserSettings = gameUserSettings
	}
	if gameIni, err := dockerManager.ReadConfigFile(uint(id), utils.GameIniFileName); err == nil {
		response.GameIni = gameIni
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "获取成功",
		"data":    response,
	})
}

// GetServerRCON 获取服务器RCON连接信息
// @Summary 获取服务器RCON信息
// @Description 获取指定服务器的RCON连接信息（包括密码）
// @Tags 服务器管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "服务器ID"
// @Success 200 {object} map[string]interface{} "RCON信息"
// @Failure 400 {object} map[string]string "请求错误"
// @Failure 404 {object} map[string]string "服务器不存在"
// @Failure 401 {object} map[string]string "未授权"
// @Router /servers/{id}/rcon [get]
func GetServerRCON(c *gin.Context) {
	userID := c.GetUint("user_id")
	serverID := c.Param("id")

	id, err := strconv.ParseUint(serverID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的服务器ID"})
		return
	}

	var server models.Server
	if err := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&server).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "服务器不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "获取成功",
		"data": gin.H{
			"server_id":         server.ID,
			"server_identifier": server.Identifier,
			"rcon_port":         server.RCONPort,
			"admin_password":    server.AdminPassword,
		},
	})
}

// ExecuteRCONCommand 执行RCON命令
// @Summary 执行RCON命令
// @Description 在指定服务器上执行RCON命令
// @Tags 服务器管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "服务器ID"
// @Param command body map[string]string true "RCON命令 {\"command\": \"ListPlayers\"}"
// @Success 200 {object} map[string]string "命令执行结果"
// @Failure 400 {object} map[string]string "请求错误"
// @Failure 404 {object} map[string]string "服务器不存在"
// @Failure 401 {object} map[string]string "未授权"
// @Failure 500 {object} map[string]string "服务器错误"
// @Router /servers/{id}/rcon/execute [post]
func ExecuteRCONCommand(c *gin.Context) {
	userID := c.GetUint("user_id")
	serverID := c.Param("id")

	id, err := strconv.ParseUint(serverID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的服务器ID"})
		return
	}

	// 解析请求体
	var req struct {
		Command string `json:"command" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	var server models.Server
	if err := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&server).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "服务器不存在"})
		return
	}

	// 检查服务器状态
	if server.Status != "running" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "服务器未运行，无法执行RCON命令"})
		return
	}

	// 执行RCON命令
	dockerManager, err := docker_manager.GetDockerManager()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取Docker管理器失败"})
		return
	}
	containerName := utils.GetServerContainerName(server.ID)

	// 检查容器是否存在
	containerExists, err := dockerManager.ContainerExists(containerName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "检查容器状态失败"})
		return
	}
	if !containerExists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "服务器容器不存在，无法执行RCON命令"})
		return
	}

	// 构建RCON命令 - 使用服务器的RCON端口
	rconCommand := fmt.Sprintf("echo '%s' | /ark/rcon -H localhost -P %d -p '%s'", req.Command, server.RCONPort, server.AdminPassword)

	output, err := dockerManager.ExecuteCommand(containerName, rconCommand)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("执行RCON命令失败: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "RCON命令执行成功",
		"data": gin.H{
			"command": req.Command,
			"output":  output,
		},
	})
}

// UpdateServer 更新服务器
// @Summary 更新服务器配置
// @Description 更新指定服务器的配置信息（包括配置文件）
// @Tags 服务器管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "服务器ID"
// @Param server body models.ServerUpdateRequest true "更新的服务器配置（可包含配置文件内容）"
// @Success 200 {object} map[string]models.ServerResponse "更新成功"
// @Failure 400 {object} map[string]string "请求错误"
// @Failure 404 {object} map[string]string "服务器不存在"
// @Failure 401 {object} map[string]string "未授权"
// @Failure 500 {object} map[string]string "服务器错误"
// @Router /servers/{id} [put]
func UpdateServer(c *gin.Context) {
	userID := c.GetUint("user_id")
	serverID := c.Param("id")

	id, err := strconv.ParseUint(serverID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的服务器ID"})
		return
	}

	var req models.ServerUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	// 查找服务器
	var server models.Server
	if err := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&server).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "服务器不存在"})
		return
	}

	// 检查标识是否与其他服务器冲突
	if req.Identifier != "" && req.Identifier != server.Identifier {
		var existingServer models.Server
		if err := database.DB.Where("identifier = ? AND user_id = ? AND id != ?", req.Identifier, userID, id).First(&existingServer).Error; err == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "服务器标识已存在"})
			return
		}
		server.Identifier = req.Identifier
	}

	// 更新字段
	if req.SessionName != "" {
		server.SessionName = req.SessionName
	}
	if req.ClusterID != "" {
		server.ClusterID = req.ClusterID
	}
	if req.Port > 0 {
		server.Port = req.Port
	}
	if req.QueryPort > 0 {
		server.QueryPort = req.QueryPort
	}
	if req.RCONPort > 0 {
		server.RCONPort = req.RCONPort
	}
	if req.AdminPassword != "" {
		server.AdminPassword = req.AdminPassword
	}
	if req.Map != "" {
		server.Map = req.Map
	}
	if req.MaxPlayers > 0 {
		server.MaxPlayers = req.MaxPlayers
	}
	if req.GameModIds != "" {
		server.GameModIds = req.GameModIds
	}
	// 检查启动参数是否发生变化
	argsChanged := false
	if req.ServerArgs != nil {
		argsJson, err := json.Marshal(req.ServerArgs)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "启动参数格式错误"})
			return
		}
		newArgsJSON := string(argsJson)
		if server.ServerArgsJSON != newArgsJSON {
			argsChanged = true
			server.ServerArgsJSON = newArgsJSON
		}
	}

	if err := database.DB.Save(&server).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器更新失败"})
		return
	}

	// 处理配置文件更新（允许单独更新任一配置文件）
	if req.GameUserSettings != "" || req.GameIni != "" {
		// 获取Docker管理器
		dockerManager, err := docker_manager.GetDockerManager()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "获取Docker管理器失败"})
			return
		}

		// 处理GameUserSettings.ini
		if req.GameUserSettings != "" {
			// 验证配置文件格式
			if err := utils.ValidateINIContent(req.GameUserSettings); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("GameUserSettings.ini格式错误: %v", err)})
				return
			}
			// 写入配置文件
			if err := dockerManager.WriteConfigFile(uint(id), utils.GameUserSettingsFileName, req.GameUserSettings); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("写入GameUserSettings.ini失败: %v", err)})
				return
			}
		}

		// 处理Game.ini
		if req.GameIni != "" {
			// 验证配置文件格式
			if err := utils.ValidateINIContent(req.GameIni); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Game.ini格式错误: %v", err)})
				return
			}
			// 写入配置文件
			if err := dockerManager.WriteConfigFile(uint(id), utils.GameIniFileName, req.GameIni); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("写入Game.ini失败: %v", err)})
				return
			}
		}
	}

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
		MaxPlayers:    server.MaxPlayers,
		GameModIds:    server.GameModIds,
		Status:        server.Status,
		AutoRestart:   server.AutoRestart,
		UserID:        server.UserID,
		CreatedAt:     server.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:     server.UpdatedAt.Format("2006-01-02 15:04:05"),
		ServerArgs:    models.FromServer(server),
	}

	// 从Docker卷读取配置文件内容并添加到响应中
	dockerManager2, err := docker_manager.GetDockerManager()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取Docker管理器失败"})
		return
	}
	if gameUserSettings, err := dockerManager2.ReadConfigFile(uint(id), utils.GameUserSettingsFileName); err == nil {
		response.GameUserSettings = gameUserSettings
	}
	if gameIni, err := dockerManager2.ReadConfigFile(uint(id), utils.GameIniFileName); err == nil {
		response.GameIni = gameIni
	}

	// 构建响应消息
	message := "服务器更新成功"
	if argsChanged && server.Status == "running" {
		message = "服务器更新成功，启动参数已修改。由于服务器正在运行，需要重启服务器以应用新的启动参数。"
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      message,
		"data":         response,
		"args_changed": argsChanged,
	})
}

// DeleteServer 删除服务器
// @Summary 删除服务器
// @Description 删除指定的服务器配置（仅允许删除已停止的服务器）
// @Tags 服务器管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "服务器ID"
// @Success 200 {object} map[string]string "删除成功"
// @Failure 400 {object} map[string]string "请求错误"
// @Failure 404 {object} map[string]string "服务器不存在"
// @Failure 401 {object} map[string]string "未授权"
// @Failure 500 {object} map[string]string "服务器错误"
// @Router /servers/{id} [delete]
func DeleteServer(c *gin.Context) {
	userID := c.GetUint("user_id")
	serverID := c.Param("id")

	id, err := strconv.ParseUint(serverID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的服务器ID"})
		return
	}

	// 查找服务器
	var server models.Server
	if err := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&server).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "服务器不存在"})
		return
	}

	// 检查服务器状态，如果正在运行则不允许删除
	if server.Status == "running" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无法删除正在运行的服务器，请先停止服务器"})
		return
	}

	// 软删除服务器
	if err := database.DB.Delete(&server).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器删除失败"})
		return
	}

	// 删除Docker容器和卷
	dockerManager, err := docker_manager.GetDockerManager()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取Docker管理器失败"})
		return
	}
	containerName := utils.GetServerContainerName(server.ID)
	volumeName := utils.GetServerVolumeName(server.ID)

	// 检查并删除容器（如果存在）
	containerExists, err := dockerManager.ContainerExists(containerName)
	if err != nil {
		fmt.Printf("Warning: Failed to check container existence: %v\n", err)
	} else if containerExists {
		if err := dockerManager.RemoveContainer(containerName); err != nil {
			fmt.Printf("Warning: Failed to remove Docker container: %v\n", err)
		}
	}

	// 删除卷
	if err := dockerManager.RemoveVolume(volumeName); err != nil {
		fmt.Printf("Warning: Failed to remove Docker volume: %v\n", err)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "服务器删除成功",
	})
}

// StartServer 启动服务器
// @Summary 启动服务器
// @Description 启动指定的ARK服务器
// @Tags 服务器管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "服务器ID"
// @Success 200 {object} map[string]string "启动成功"
// @Failure 400 {object} map[string]string "请求错误"
// @Failure 404 {object} map[string]string "服务器不存在"
// @Failure 401 {object} map[string]string "未授权"
// @Failure 500 {object} map[string]string "服务器错误"
// @Router /servers/{id}/start [post]
func StartServer(c *gin.Context) {
	userID := c.GetUint("user_id")
	serverID := c.Param("id")

	id, err := strconv.ParseUint(serverID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的服务器ID"})
		return
	}

	// 查找服务器
	var server models.Server
	if err := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&server).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "服务器不存在"})
		return
	}

	// 检查服务器状态
	if server.Status == "running" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "服务器已在运行中"})
		return
	}

	if server.Status == "starting" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "服务器正在启动中"})
		return
	}

	// 更新服务器状态为启动中
	server.Status = "starting"
	if err := database.DB.Save(&server).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新服务器状态失败"})
		return
	}

	// 启动Docker容器
	dockerManager, err := docker_manager.GetDockerManager()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取Docker管理器失败"})
		return
	}
	containerName := utils.GetServerContainerName(server.ID)

	go func() {
		// 检查镜像是否存在
		imageName := "tbro98/ase-server:latest"
		exists, err := dockerManager.ImageExists(imageName)
		if err != nil {
			fmt.Printf("Failed to check image existence: %v\n", err)
			database.DB.Model(&server).Update("status", "stopped")
			return
		}
		if !exists {
			fmt.Printf("Image %s does not exist, cannot start server\n", imageName)
			database.DB.Model(&server).Update("status", "stopped")
			return
		}

		// 检查容器是否存在
		containerExists, err := dockerManager.ContainerExists(containerName)
		if err != nil {
			fmt.Printf("Failed to check container existence: %v\n", err)
			database.DB.Model(&server).Update("status", "stopped")
			return
		}

		needRecreateContainer := false

		if containerExists {
			// 容器存在，检查是否需要重建（比较启动参数）
			envVars, err := dockerManager.GetContainerEnvVars(containerName)
			if err != nil {
				fmt.Printf("Failed to get container environment variables: %v\n", err)
				// 如果无法获取环境变量，为了安全起见重建容器
				needRecreateContainer = true
			} else {
				// 获取当前服务器的启动参数
				var serverArgs *models.ServerArgs
				if server.ServerArgsJSON != "" && server.ServerArgsJSON != "{}" {
					serverArgs = models.NewServerArgs()
					if err := json.Unmarshal([]byte(server.ServerArgsJSON), serverArgs); err != nil {
						serverArgs = models.FromServer(server)
					}
				} else {
					serverArgs = models.FromServer(server)
				}
				currentArgsString := serverArgs.GenerateArgsString(server)

				// 比较环境变量中的SERVER_ARGS
				if containerArgsString, exists := envVars["SERVER_ARGS"]; exists {
					if containerArgsString != currentArgsString {
						fmt.Printf("Server arguments changed, need to recreate container. Old: %s, New: %s\n", containerArgsString, currentArgsString)
						needRecreateContainer = true
					}
				} else {
					// 如果容器中没有SERVER_ARGS环境变量，重建容器
					fmt.Printf("Container missing SERVER_ARGS environment variable, need to recreate\n")
					needRecreateContainer = true
				}

				// 检查其他可能影响容器配置的参数
				if !needRecreateContainer {
					// 检查GameModIds
					if server.GameModIds != envVars["GameModIds"] {
						fmt.Printf("GameModIds changed, need to recreate container. Old: %s, New: %s\n", envVars["GameModIds"], server.GameModIds)
						needRecreateContainer = true
					}
				}
			}

			if needRecreateContainer {
				fmt.Printf("Recreating container due to configuration changes\n")
				// 删除现有容器
				if err := dockerManager.RemoveContainer(containerName); err != nil {
					fmt.Printf("Failed to remove existing container: %v\n", err)
					database.DB.Model(&server).Update("status", "stopped")
					return
				}
			}
		}

		// 如果容器不存在或需要重建，创建新容器
		if !containerExists || needRecreateContainer {
			// 创建容器
			_, err = dockerManager.CreateContainer(server.ID, server.Identifier, server.Port, server.QueryPort, server.RCONPort, server.AdminPassword, server.Map, server.GameModIds, server.AutoRestart)
			if err != nil {
				fmt.Printf("Failed to create Docker container: %v\n", err)
				database.DB.Model(&server).Update("status", "stopped")
				return
			}
		}

		// 启动容器
		if err := dockerManager.StartContainer(containerName); err != nil {
			fmt.Printf("Failed to start Docker container: %v\n", err)
			// 更新状态为停止
			database.DB.Model(&server).Update("status", "stopped")
			return
		}

		// 等待容器启动并验证状态
		for i := 0; i < 30; i++ { // 最多等待30秒
			time.Sleep(1 * time.Second)
			status, err := dockerManager.GetContainerStatus(containerName)
			if err != nil {
				fmt.Printf("Failed to get container status: %v\n", err)
				continue
			}

			if status == "running" {
				// 更新状态为运行中
				if err := database.DB.Model(&server).Update("status", "running").Error; err != nil {
					fmt.Printf("Failed to update server status to running: %v\n", err)
				}
				return
			}
		}

		// 超时或启动失败
		fmt.Printf("Container failed to start or timeout\n")
		database.DB.Model(&server).Update("status", "stopped")
	}()

	c.JSON(http.StatusOK, gin.H{
		"message": "服务器启动命令已发送",
	})
}

// StopServer 停止服务器
// @Summary 停止服务器
// @Description 停止指定的ARK服务器
// @Tags 服务器管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "服务器ID"
// @Success 200 {object} map[string]string "停止成功"
// @Failure 400 {object} map[string]string "请求错误"
// @Failure 404 {object} map[string]string "服务器不存在"
// @Failure 401 {object} map[string]string "未授权"
// @Failure 500 {object} map[string]string "服务器错误"
// @Router /servers/{id}/stop [post]
func StopServer(c *gin.Context) {
	userID := c.GetUint("user_id")
	serverID := c.Param("id")

	id, err := strconv.ParseUint(serverID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的服务器ID"})
		return
	}

	// 查找服务器
	var server models.Server
	if err := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&server).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "服务器不存在"})
		return
	}

	// 检查服务器状态
	if server.Status == "stopped" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "服务器已经停止"})
		return
	}

	if server.Status == "stopping" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "服务器正在停止中"})
		return
	}

	// 更新服务器状态为停止中
	server.Status = "stopping"
	if err := database.DB.Save(&server).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新服务器状态失败"})
		return
	}

	// 停止Docker容器
	dockerManager, err := docker_manager.GetDockerManager()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取Docker管理器失败"})
		return
	}
	containerName := utils.GetServerContainerName(server.ID)

	go func() {
		// 检查容器是否存在
		containerExists, err := dockerManager.ContainerExists(containerName)
		if err != nil {
			fmt.Printf("Failed to check container existence: %v\n", err)
			// 即使检查失败，也更新状态为停止
			database.DB.Model(&server).Update("status", "stopped")
			return
		}

		if !containerExists {
			// 容器不存在，直接更新状态为停止
			database.DB.Model(&server).Update("status", "stopped")
			return
		}

		// 停止容器
		if err := dockerManager.StopContainer(containerName); err != nil {
			fmt.Printf("Failed to stop Docker container: %v\n", err)
			// 即使停止失败，也更新状态（可能容器已经停止）
		}

		// 验证容器状态
		for i := 0; i < 30; i++ { // 最多等待30秒
			time.Sleep(1 * time.Second)
			status, err := dockerManager.GetContainerStatus(containerName)
			if err != nil {
				fmt.Printf("Failed to get container status: %v\n", err)
				continue
			}

			if status == "stopped" || status == "not_found" {
				break
			}
		}

		// 更新状态为已停止
		if err := database.DB.Model(&server).Update("status", "stopped").Error; err != nil {
			fmt.Printf("Failed to update server status to stopped: %v\n", err)
		}
	}()

	c.JSON(http.StatusOK, gin.H{
		"message": "服务器停止命令已发送",
	})
}

// GetServerFolderInfo 获取服务器文件夹信息（内部管理使用）
// @Summary 获取服务器文件夹信息
// @Description 获取指定服务器的文件夹路径和大小（仅供内部管理使用）
// @Tags 服务器管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "服务器ID"
// @Success 200 {object} map[string]interface{} "文件夹信息"
// @Failure 400 {object} map[string]string "请求错误"
// @Failure 404 {object} map[string]string "服务器不存在"
// @Failure 401 {object} map[string]string "未授权"
// @Router /servers/{id}/folder [get]
func GetServerFolderInfo(c *gin.Context) {
	userID := c.GetUint("user_id")
	serverID := c.Param("id")

	id, err := strconv.ParseUint(serverID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的服务器ID"})
		return
	}

	var server models.Server
	if err := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&server).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "服务器不存在"})
		return
	}

	// 获取Docker卷信息
	dockerManager, err := docker_manager.GetDockerManager()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取Docker管理器失败"})
		return
	}
	volumeName := utils.GetServerVolumeName(server.ID)
	pluginsVolumeName := utils.GetServerPluginsVolumeName(server.ID)
	containerName := utils.GetServerContainerName(server.ID)

	// 检查卷是否存在
	volumeExists, _ := dockerManager.VolumeExists(volumeName)
	pluginsVolumeExists, _ := dockerManager.VolumeExists(pluginsVolumeName)

	// 检查容器是否存在
	containerExists, _ := dockerManager.ContainerExists(containerName)

	// 获取容器状态
	containerStatus := "not_found"
	if containerExists {
		if status, err := dockerManager.GetContainerStatus(containerName); err == nil {
			containerStatus = status
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "获取成功",
		"data": gin.H{
			"server_id":             server.ID,
			"server_name":           server.Identifier,
			"volume_name":           volumeName,
			"volume_exists":         volumeExists,
			"plugins_volume_name":   pluginsVolumeName,
			"plugins_volume_exists": pluginsVolumeExists,
			"container_name":        containerName,
			"container_exists":      containerExists,
			"container_status":      containerStatus,
		},
	})
}

// GetImageStatus 获取镜像状态和拉取进度
// @Summary 获取镜像状态
// @Description 获取ARK服务器镜像的状态信息（镜像在后台异步拉取）
// @Tags 服务器管理
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} map[string]interface{} "镜像状态信息"
// @Failure 401 {object} map[string]string "未授权"
// @Failure 500 {object} map[string]string "服务器错误"
// @Router /servers/images/status [get]
func GetImageStatus(c *gin.Context) {
	// 创建Docker管理器
	dockerManager, err := docker_manager.GetDockerManager()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取Docker管理器失败"})
		return
	}

	// 需要检查的镜像列表
	requiredImages := []string{
		"tbro98/ase-server:latest", // ARK服务器镜像
		"alpine:latest",            // Alpine镜像（用于配置文件操作）
	}

	// 获取每个镜像的状态
	imageStatuses := make(map[string]*docker_manager.ImageStatus)
	allReady := true
	anyPulling := false
	pullingCount := 0

	for _, imageName := range requiredImages {
		status := dockerManager.GetImageStatus(imageName)
		imageStatuses[imageName] = status

		if !status.Ready {
			allReady = false
		}

		if status.Pulling {
			anyPulling = true
			pullingCount++
		}
	}

	// 生成总体状态描述
	var overallStatus string
	if allReady {
		overallStatus = "所有镜像已就绪"
	} else if anyPulling {
		overallStatus = "正在下载镜像"
	} else {
		overallStatus = "镜像未就绪，等待下载"
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "获取成功",
		"data": gin.H{
			"images":            imageStatuses,
			"any_pulling":       anyPulling, // 检查是否有镜像正在拉取中
			"any_not_ready":     !allReady,
			"can_create_server": true,     // 现在创建服务器不需要镜像就绪，只需要卷
			"can_start_server":  allReady, // 启动服务器需要镜像就绪
			"overall_status":    overallStatus,
			"pulling_count":     pullingCount,
			"total_images":      len(requiredImages),
		},
	})
}
