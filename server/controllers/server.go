package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"ark-server-manager/database"
	"ark-server-manager/models"
	"ark-server-manager/utils"

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
	dockerManager := utils.NewDockerManager()
	var serverResponses []models.ServerResponse
	for _, server := range servers {
		// 获取Docker容器实时状态
		containerName := utils.GetServerContainerName(server.ID)
		realTimeStatus := server.Status // 默认使用数据库状态
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

		serverResponses = append(serverResponses, models.ServerResponse{
			ID:            server.ID,
			Identifier:    server.Identifier,
			Port:          server.Port,
			QueryPort:     server.QueryPort,
			RCONPort:      server.RCONPort,
			AdminPassword: server.AdminPassword,
			Map:           server.Map,
			Status:        realTimeStatus,
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

	// 创建服务器
	server := models.Server{
		Identifier:    req.Identifier,
		Port:          req.Port,
		QueryPort:     req.QueryPort,
		RCONPort:      req.RCONPort,
		AdminPassword: req.AdminPassword,
		Map:           req.Map,
		Status:        "stopped",
		UserID:        userID,
	}

	if err := database.DB.Create(&server).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器创建失败"})
		return
	}

	// 检查Docker环境
	if err := utils.CheckDockerStatus(); err != nil {
		// 删除数据库记录并返回错误
		database.DB.Delete(&server)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Docker环境检查失败: %v", err)})
		return
	}

	// 创建Docker卷和容器
	dockerManager := utils.NewDockerManager()

	// 创建Docker卷
	volumeName, err := dockerManager.CreateVolume(server.ID)
	if err != nil {
		// 删除数据库记录并返回错误
		database.DB.Delete(&server)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("创建Docker卷失败: %v", err)})
		return
	}
	fmt.Printf("Created Docker volume: %s\n", volumeName)

	// 创建Docker容器（不立即启动）
	containerID, err := dockerManager.CreateContainer(server.ID, server.Identifier, server.Port, server.QueryPort, server.RCONPort, server.AdminPassword, server.Map)
	if err != nil {
		// 清理已创建的卷
		dockerManager.RemoveVolume(volumeName)
		database.DB.Delete(&server)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("创建Docker容器失败: %v", err)})
		return
	}
	fmt.Printf("Created Docker container: %s\n", containerID)

	// 停止容器（创建时会自动启动）
	containerName := utils.GetServerContainerName(server.ID)
	if err := dockerManager.StopContainer(containerName); err != nil {
		fmt.Printf("Warning: Failed to stop container after creation: %v\n", err)
	}

	// 创建默认配置文件到Docker卷中
	gameUserSettings := utils.GetDefaultGameUserSettings(server.Identifier, server.Map, server.Port, server.QueryPort, server.RCONPort, 70, server.AdminPassword)
	if err := dockerManager.WriteConfigFile(server.ID, utils.GameUserSettingsFileName, gameUserSettings); err != nil {
		fmt.Printf("Warning: Failed to create default GameUserSettings.ini: %v\n", err)
	}

	gameIni := utils.GetDefaultGameIni()
	if err := dockerManager.WriteConfigFile(server.ID, utils.GameIniFileName, gameIni); err != nil {
		fmt.Printf("Warning: Failed to create default Game.ini: %v\n", err)
	}

	response := models.ServerResponse{
		ID:            server.ID,
		Identifier:    server.Identifier,
		Port:          server.Port,
		QueryPort:     server.QueryPort,
		RCONPort:      server.RCONPort,
		AdminPassword: server.AdminPassword,
		Map:           server.Map,
		Status:        server.Status,
		UserID:        server.UserID,
		CreatedAt:     server.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:     server.UpdatedAt.Format("2006-01-02 15:04:05"),
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

	response := models.ServerResponse{
		ID:            server.ID,
		Identifier:    server.Identifier,
		Port:          server.Port,
		QueryPort:     server.QueryPort,
		RCONPort:      server.RCONPort,
		AdminPassword: server.AdminPassword,
		Map:           server.Map,
		Status:        server.Status,
		UserID:        server.UserID,
		CreatedAt:     server.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:     server.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	// 从Docker卷读取配置文件内容（如果存在）
	dockerManager := utils.NewDockerManager()
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
	dockerManager := utils.NewDockerManager()
	containerName := utils.GetServerContainerName(server.ID)

	// 构建RCON命令
	rconCommand := fmt.Sprintf("echo '%s' | /ark/rcon -H localhost -P 32330 -p '%s'", req.Command, server.AdminPassword)

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

	if err := database.DB.Save(&server).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器更新失败"})
		return
	}

	// 更新数据时必须有配置文件数据
	if req.GameUserSettings != "" && req.GameIni != "" {
		// 验证配置文件格式
		if req.GameUserSettings != "" {
			if err := utils.ValidateINIContent(req.GameUserSettings); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("GameUserSettings.ini格式错误: %v", err)})
				return
			}
		}
		if req.GameIni != "" {
			if err := utils.ValidateINIContent(req.GameIni); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Game.ini格式错误: %v", err)})
				return
			}
		}

		// 写入配置文件到Docker卷
		dockerManager := utils.NewDockerManager()
		if req.GameUserSettings != "" {
			if err := dockerManager.WriteConfigFile(uint(id), utils.GameUserSettingsFileName, req.GameUserSettings); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("写入GameUserSettings.ini失败: %v", err)})
				return
			}
		}
		if req.GameIni != "" {
			if err := dockerManager.WriteConfigFile(uint(id), utils.GameIniFileName, req.GameIni); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("写入Game.ini失败: %v", err)})
				return
			}
		}
	}

	response := models.ServerResponse{
		ID:            server.ID,
		Identifier:    server.Identifier,
		Port:          server.Port,
		QueryPort:     server.QueryPort,
		RCONPort:      server.RCONPort,
		AdminPassword: server.AdminPassword,
		Map:           server.Map,
		Status:        server.Status,
		UserID:        server.UserID,
		CreatedAt:     server.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:     server.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	// 从Docker卷读取配置文件内容并添加到响应中
	dockerManager2 := utils.NewDockerManager()
	if gameUserSettings, err := dockerManager2.ReadConfigFile(uint(id), utils.GameUserSettingsFileName); err == nil {
		response.GameUserSettings = gameUserSettings
	}
	if gameIni, err := dockerManager2.ReadConfigFile(uint(id), utils.GameIniFileName); err == nil {
		response.GameIni = gameIni
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "服务器更新成功",
		"data":    response,
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
	dockerManager := utils.NewDockerManager()
	containerName := utils.GetServerContainerName(server.ID)
	volumeName := utils.GetServerVolumeName(server.ID)

	// 删除容器
	if err := dockerManager.RemoveContainer(containerName); err != nil {
		fmt.Printf("Warning: Failed to remove Docker container: %v\n", err)
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
	dockerManager := utils.NewDockerManager()
	containerName := utils.GetServerContainerName(server.ID)

	go func() {
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
	dockerManager := utils.NewDockerManager()
	containerName := utils.GetServerContainerName(server.ID)

	go func() {
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
	dockerManager := utils.NewDockerManager()
	volumeName := utils.GetServerVolumeName(server.ID)
	containerName := utils.GetServerContainerName(server.ID)

	// 检查卷是否存在
	volumeExists, _ := dockerManager.VolumeExists(volumeName)

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
			"server_id":        server.ID,
			"server_name":      server.Identifier,
			"volume_name":      volumeName,
			"volume_exists":    volumeExists,
			"container_name":   containerName,
			"container_exists": containerExists,
			"container_status": containerStatus,
		},
	})
}
