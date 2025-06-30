package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"ark-server-manager/database"
	"ark-server-manager/models"

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

	// 转换为响应格式
	var serverResponses []models.ServerResponse
	for _, server := range servers {
		serverResponses = append(serverResponses, models.ServerResponse{
			ID:            server.ID,
			Name:          server.Name,
			Description:   server.Description,
			Port:          server.Port,
			QueryPort:     server.QueryPort,
			RCONPort:      server.RCONPort,
			AdminPassword: server.AdminPassword,
			Map:           server.Map,
			MaxPlayers:    server.MaxPlayers,
			Status:        server.Status,
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

// GetServer 获取单个服务器信息
// @Summary 获取服务器详情
// @Description 根据ID获取指定服务器的详细信息
// @Tags 服务器管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "服务器ID"
// @Success 200 {object} map[string]models.ServerResponse "服务器信息"
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
		Name:          server.Name,
		Description:   server.Description,
		Port:          server.Port,
		QueryPort:     server.QueryPort,
		RCONPort:      server.RCONPort,
		AdminPassword: server.AdminPassword,
		Map:           server.Map,
		MaxPlayers:    server.MaxPlayers,
		Status:        server.Status,
		UserID:        server.UserID,
		CreatedAt:     server.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:     server.UpdatedAt.Format("2006-01-02 15:04:05"),
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
			"server_id":      server.ID,
			"server_name":    server.Name,
			"rcon_port":      server.RCONPort,
			"admin_password": server.AdminPassword,
		},
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

	// 检查服务器名称是否已存在
	var existingServer models.Server
	if err := database.DB.Where("name = ? AND user_id = ?", req.Name, userID).First(&existingServer).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "服务器名称已存在"})
		return
	}

	// 设置默认值
	if req.Map == "" {
		req.Map = "TheIsland"
	}

	// 创建服务器
	server := models.Server{
		Name:          req.Name,
		Description:   req.Description,
		Port:          req.Port,
		QueryPort:     req.QueryPort,
		RCONPort:      req.RCONPort,
		AdminPassword: req.AdminPassword,
		Map:           req.Map,
		MaxPlayers:    req.MaxPlayers,
		Status:        "stopped",
		UserID:        userID,
	}

	if err := database.DB.Create(&server).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器创建失败"})
		return
	}

	response := models.ServerResponse{
		ID:            server.ID,
		Name:          server.Name,
		Description:   server.Description,
		Port:          server.Port,
		QueryPort:     server.QueryPort,
		RCONPort:      server.RCONPort,
		AdminPassword: server.AdminPassword,
		Map:           server.Map,
		MaxPlayers:    server.MaxPlayers,
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

// UpdateServer 更新服务器
// @Summary 更新服务器配置
// @Description 更新指定服务器的配置信息
// @Tags 服务器管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "服务器ID"
// @Param server body models.ServerUpdateRequest true "更新的服务器配置"
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

	// 检查名称是否与其他服务器冲突
	if req.Name != "" && req.Name != server.Name {
		var existingServer models.Server
		if err := database.DB.Where("name = ? AND user_id = ? AND id != ?", req.Name, userID, id).First(&existingServer).Error; err == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "服务器名称已存在"})
			return
		}
		server.Name = req.Name
	}

	// 更新字段
	if req.Description != "" {
		server.Description = req.Description
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

	if err := database.DB.Save(&server).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器更新失败"})
		return
	}

	response := models.ServerResponse{
		ID:            server.ID,
		Name:          server.Name,
		Description:   server.Description,
		Port:          server.Port,
		QueryPort:     server.QueryPort,
		RCONPort:      server.RCONPort,
		AdminPassword: server.AdminPassword,
		Map:           server.Map,
		MaxPlayers:    server.MaxPlayers,
		Status:        server.Status,
		UserID:        server.UserID,
		CreatedAt:     server.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:     server.UpdatedAt.Format("2006-01-02 15:04:05"),
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

	// TODO: 这里应该实现实际的服务器启动逻辑
	// 目前只是模拟启动过程
	go func() {
		// 模拟启动过程
		time.Sleep(3 * time.Second)

		// 更新状态为运行中
		if err := database.DB.Model(&server).Update("status", "running").Error; err != nil {
			// 记录错误日志
			fmt.Printf("Failed to update server status to running: %v\n", err)
		}
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

	// TODO: 这里应该实现实际的服务器停止逻辑
	// 目前只是模拟停止过程
	go func() {
		// 模拟停止过程
		time.Sleep(2 * time.Second)

		// 更新状态为已停止
		if err := database.DB.Model(&server).Update("status", "stopped").Error; err != nil {
			// 记录错误日志
			fmt.Printf("Failed to update server status to stopped: %v\n", err)
		}
	}()

	c.JSON(http.StatusOK, gin.H{
		"message": "服务器停止命令已发送",
	})
}
