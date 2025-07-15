package controllers

import (
	"net/http"

	"ark-server-commander/models"
	"ark-server-commander/service/server"

	"github.com/gin-gonic/gin"
)

var serverService = server.NewServerService()

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

	serverResponses, err := serverService.GetServers(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
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

	response, err := serverService.CreateServer(userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
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

	response, err := serverService.GetServer(userID, serverID)
	if err != nil {
		if err.Error() == "无效的服务器ID" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err.Error() == "服务器不存在" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
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

	data, err := serverService.GetServerRCON(userID, serverID)
	if err != nil {
		if err.Error() == "无效的服务器ID" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err.Error() == "服务器不存在" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "获取成功",
		"data":    data,
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

	var req struct {
		Command string `json:"command" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	data, err := serverService.ExecuteRCONCommand(userID, serverID, req.Command)
	if err != nil {
		if err.Error() == "无效的服务器ID" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err.Error() == "服务器不存在" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if err.Error() == "服务器未运行，无法执行RCON命令" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err.Error() == "服务器容器不存在，无法执行RCON命令" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "RCON命令执行成功",
		"data":    data,
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

	var req models.ServerUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	response, argsChanged, err := serverService.UpdateServer(userID, serverID, req)
	if err != nil {
		if err.Error() == "无效的服务器ID" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err.Error() == "服务器不存在" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if err.Error() == "服务器标识已存在" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 构建响应消息
	message := "服务器更新成功"
	if argsChanged && response.Status == "running" {
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

	err := serverService.DeleteServer(userID, serverID)
	if err != nil {
		if err.Error() == "无效的服务器ID" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err.Error() == "服务器不存在" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if err.Error() == "无法删除正在运行的服务器，请先停止服务器" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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

	err := serverService.StartServer(userID, serverID)
	if err != nil {
		if err.Error() == "无效的服务器ID" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err.Error() == "服务器不存在" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if err.Error() == "服务器已在运行中" || err.Error() == "服务器正在启动中" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

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

	err := serverService.StopServer(userID, serverID)
	if err != nil {
		if err.Error() == "无效的服务器ID" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err.Error() == "服务器不存在" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if err.Error() == "服务器已经停止" || err.Error() == "服务器正在停止中" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

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

	data, err := serverService.GetServerFolderInfo(userID, serverID)
	if err != nil {
		if err.Error() == "无效的服务器ID" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err.Error() == "服务器不存在" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "获取成功",
		"data":    data,
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
	data, err := serverService.GetImageStatus()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "获取成功",
		"data":    data,
	})
}
