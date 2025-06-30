package controllers

import (
	"net/http"
	"strconv"

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

	// 转换为响应格式
	var serverResponses []models.ServerResponse
	for _, server := range servers {
		serverResponses = append(serverResponses, models.ServerResponse{
			ID:          server.ID,
			Name:        server.Name,
			Description: server.Description,
			IP:          server.IP,
			Port:        server.Port,
			QueryPort:   server.QueryPort,
			RCONPort:    server.RCONPort,
			Map:         server.Map,
			MaxPlayers:  server.MaxPlayers,
			Status:      server.Status,
			UserID:      server.UserID,
			CreatedAt:   server.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:   server.UpdatedAt.Format("2006-01-02 15:04:05"),
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
		ID:          server.ID,
		Name:        server.Name,
		Description: server.Description,
		IP:          server.IP,
		Port:        server.Port,
		QueryPort:   server.QueryPort,
		RCONPort:    server.RCONPort,
		Map:         server.Map,
		MaxPlayers:  server.MaxPlayers,
		Status:      server.Status,
		UserID:      server.UserID,
		CreatedAt:   server.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   server.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "获取成功",
		"data":    response,
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

	// 加密RCON密码
	hashedPassword, err := utils.HashPassword(req.RCONPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "密码加密失败"})
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
		Name:         req.Name,
		Description:  req.Description,
		IP:           req.IP,
		Port:         req.Port,
		QueryPort:    req.QueryPort,
		RCONPort:     req.RCONPort,
		RCONPassword: hashedPassword,
		Map:          req.Map,
		MaxPlayers:   req.MaxPlayers,
		Status:       "stopped",
		UserID:       userID,
	}

	if err := database.DB.Create(&server).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器创建失败"})
		return
	}

	response := models.ServerResponse{
		ID:          server.ID,
		Name:        server.Name,
		Description: server.Description,
		IP:          server.IP,
		Port:        server.Port,
		QueryPort:   server.QueryPort,
		RCONPort:    server.RCONPort,
		Map:         server.Map,
		MaxPlayers:  server.MaxPlayers,
		Status:      server.Status,
		UserID:      server.UserID,
		CreatedAt:   server.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   server.UpdatedAt.Format("2006-01-02 15:04:05"),
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
	if req.IP != "" {
		server.IP = req.IP
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
		ID:          server.ID,
		Name:        server.Name,
		Description: server.Description,
		IP:          server.IP,
		Port:        server.Port,
		QueryPort:   server.QueryPort,
		RCONPort:    server.RCONPort,
		Map:         server.Map,
		MaxPlayers:  server.MaxPlayers,
		Status:      server.Status,
		UserID:      server.UserID,
		CreatedAt:   server.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   server.UpdatedAt.Format("2006-01-02 15:04:05"),
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
