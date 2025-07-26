package auth

import (
	"net/http"

	"ark-server-commander/database"
	"ark-server-commander/models"
	"ark-server-commander/utils"

	"github.com/gin-gonic/gin"
)

// CheckInit 检查是否已初始化用户
// @Summary 检查系统初始化状态
// @Description 检查系统是否已经初始化过用户
// @Tags 认证
// @Accept json
// @Produce json
// @Success 200 {object} map[string]bool "初始化状态"
// @Router /auth/check-init [get]
func CheckInit(c *gin.Context) {
	var count int64
	database.DB.Model(&models.User{}).Count(&count)

	c.JSON(http.StatusOK, gin.H{
		"initialized": count > 0,
	})
}

// InitUser 初始化用户
// @Summary 初始化系统用户
// @Description 创建第一个管理员用户，只能在系统未初始化时调用
// @Tags 认证
// @Accept json
// @Produce json
// @Param user body models.UserRequest true "用户信息"
// @Success 200 {object} map[string]interface{} "初始化成功"
// @Failure 400 {object} map[string]string "请求错误"
// @Failure 500 {object} map[string]string "服务器错误"
// @Router /auth/init [post]
func InitUser(c *gin.Context) {
	var req models.UserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	// 检查是否已有用户
	var count int64
	database.DB.Model(&models.User{}).Count(&count)
	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "系统已初始化"})
		return
	}

	// 加密密码
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "密码加密失败"})
		return
	}

	// 创建用户
	user := models.User{
		Username: req.Username,
		Password: hashedPassword,
	}

	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "用户创建失败"})
		return
	}

	// 生成token
	token, err := utils.GenerateToken(user.ID, user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "令牌生成失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "初始化成功",
		"token":   token,
		"user": models.UserResponse{
			ID:       user.ID,
			Username: user.Username,
		},
	})
}

// Login 用户登录
// @Summary 用户登录
// @Description 使用用户名和密码登录系统
// @Tags 认证
// @Accept json
// @Produce json
// @Param credentials body models.UserRequest true "登录凭据"
// @Success 200 {object} map[string]interface{} "登录成功"
// @Failure 400 {object} map[string]string "请求错误"
// @Failure 401 {object} map[string]string "认证失败"
// @Failure 500 {object} map[string]string "服务器错误"
// @Router /auth/login [post]
func Login(c *gin.Context) {
	var req models.UserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	// 查找用户
	var user models.User
	if err := database.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	// 验证密码
	if !utils.CheckPassword(req.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	// 生成token
	token, err := utils.GenerateToken(user.ID, user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "令牌生成失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "登录成功",
		"token":   token,
		"user": models.UserResponse{
			ID:       user.ID,
			Username: user.Username,
		},
	})
}

// GetProfile 获取用户信息
// @Summary 获取当前用户信息
// @Description 获取当前登录用户的基本信息
// @Tags 用户
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} map[string]models.UserResponse "用户信息"
// @Failure 401 {object} map[string]string "未授权"
// @Router /profile [get]
func GetProfile(c *gin.Context) {
	userID := c.GetUint("user_id")
	username := c.GetString("username")

	c.JSON(http.StatusOK, gin.H{
		"user": models.UserResponse{
			ID:       userID,
			Username: username,
		},
	})
}