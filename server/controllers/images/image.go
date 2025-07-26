package images

import (
	"net/http"

	"ark-server-commander/service/server"

	"github.com/gin-gonic/gin"
)

var serverService = server.NewServerService()

// PullImage 手动拉取镜像
// @Summary 手动拉取Docker镜像
// @Description 用户主动触发镜像下载操作
// @Tags 镜像管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body map[string]string true "镜像信息 {\"image_name\": \"tbro98/ase-server:latest\"}"
// @Success 200 {object} map[string]interface{} "拉取状态"
// @Failure 400 {object} map[string]string "请求错误"
// @Failure 401 {object} map[string]string "未授权"
// @Failure 500 {object} map[string]string "服务器错误"
// @Router /servers/images/pull [post]
func PullImage(c *gin.Context) {
	var req struct {
		ImageName string `json:"image_name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	err := serverService.PullImage(req.ImageName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "镜像拉取已开始",
		"data": map[string]interface{}{
			"image_name": req.ImageName,
			"status":     "pulling",
		},
	})
}

// CheckImageUpdates 检查镜像更新
// @Summary 检查镜像是否有更新
// @Description 检查所有管理的镜像是否有新版本
// @Tags 镜像管理
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} map[string]bool "镜像更新状态映射"
// @Failure 401 {object} map[string]string "未授权"
// @Failure 500 {object} map[string]string "服务器错误"
// @Router /servers/images/check-updates [get]
func CheckImageUpdates(c *gin.Context) {
	data, err := serverService.CheckImageUpdates()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "检查完成",
		"data":    data,
	})
}

// UpdateImage 更新镜像
// @Summary 更新Docker镜像
// @Description 更新指定镜像并处理相关容器
// @Tags 镜像管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body map[string]string true "镜像信息 {\"image_name\": \"tbro98/ase-server:latest\"}"
// @Success 200 {object} map[string]interface{} "更新状态"
// @Failure 400 {object} map[string]string "请求错误"
// @Failure 401 {object} map[string]string "未授权"
// @Failure 500 {object} map[string]string "服务器错误"
// @Router /servers/images/update [post]
func UpdateImage(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req struct {
		ImageName string `json:"image_name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	affectedServers, err := serverService.UpdateImage(req.ImageName, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "镜像更新已开始",
		"data": map[string]interface{}{
			"image_name":       req.ImageName,
			"affected_servers": affectedServers,
			"status":           "updating",
		},
	})
}

// GetAffectedServers 获取使用指定镜像的服务器列表
// @Summary 获取影响的服务器列表
// @Description 获取使用指定镜像的服务器列表
// @Tags 镜像管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param image_name query string true "镜像名称"
// @Success 200 {object} map[string]interface{} "服务器列表"
// @Failure 400 {object} map[string]string "请求错误"
// @Failure 401 {object} map[string]string "未授权"
// @Failure 500 {object} map[string]string "服务器错误"
// @Router /servers/images/affected [get]
func GetAffectedServers(c *gin.Context) {
	userID := c.GetUint("user_id")
	imageName := c.Query("image_name")

	if imageName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "镜像名称不能为空"})
		return
	}

	servers, err := serverService.GetAffectedServers(imageName, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "获取成功",
		"data": map[string]interface{}{
			"image_name": imageName,
			"servers":    servers,
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