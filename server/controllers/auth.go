package controllers

import (
	"net/http"

	"ark-server-manager/database"
	"ark-server-manager/models"
	"ark-server-manager/utils"

	"github.com/gin-gonic/gin"
)

// CheckInit 检查是否已初始化用户
func CheckInit(c *gin.Context) {
	var count int64
	database.DB.Model(&models.User{}).Count(&count)

	c.JSON(http.StatusOK, gin.H{
		"initialized": count > 0,
	})
}

// InitUser 初始化用户
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
