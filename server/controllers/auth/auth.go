package auth

import (
	"fmt"
	"math"
	"net/http"
	"sync"

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
// initMu 让"检查是否已有用户"和"创建用户"成为一个整体，避免并发请求创建出多个管理员
var initMu sync.Mutex

func InitUser(c *gin.Context) {
	var req models.UserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	// 限制用户名和密码长度：bcrypt 只处理前 72 字节，超过会直接报错
	if len(req.Username) == 0 || len(req.Username) > 64 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户名长度需要在 1 到 64 个字节之间"})
		return
	}
	if len(req.Password) < 8 || len(req.Password) > 72 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "密码长度需要在 8 到 72 个字节之间"})
		return
	}

	initMu.Lock()
	defer initMu.Unlock()

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

	ip := c.ClientIP()
	if wait := limiter.retryAfter(ip, req.Username); wait > 0 {
		minutes := int(math.Ceil(wait.Minutes()))
		c.Header("Retry-After", fmt.Sprint(int(wait.Seconds())))
		c.JSON(http.StatusTooManyRequests, gin.H{
			"error":   fmt.Sprintf("登录失败次数过多，请 %d 分钟后再试", minutes),
			"code":    "too_many_attempts",
			"minutes": minutes,
		})
		return
	}

	// 查找用户
	var user models.User
	if err := database.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		limiter.fail(ip, req.Username)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误", "code": "invalid_credentials"})
		return
	}

	// 验证密码
	if !utils.CheckPassword(req.Password, user.Password) {
		limiter.fail(ip, req.Username)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误", "code": "invalid_credentials"})
		return
	}
	limiter.success(ip, req.Username)

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

// ChangePasswordRequest 修改密码请求
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

// ChangePassword 修改登录密码
// @Summary 修改登录密码
// @Description 校验当前密码后设置新密码（8 到 72 字节）
// @Tags 用户
// @Accept json
// @Produce json
// @Security Bearer
// @Param body body ChangePasswordRequest true "当前密码和新密码"
// @Success 200 {object} map[string]string "修改成功"
// @Failure 400 {object} map[string]string "请求错误或当前密码错误"
// @Router /auth/password [put]
func ChangePassword(c *gin.Context) {
	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误", "code": "bad_request"})
		return
	}
	if len(req.NewPassword) < 8 || len(req.NewPassword) > 72 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "密码长度需要在 8 到 72 个字节之间", "code": "invalid_password_length"})
		return
	}

	var user models.User
	if err := database.DB.First(&user, c.GetUint("user_id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}
	if !utils.CheckPassword(req.OldPassword, user.Password) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "当前密码错误", "code": "wrong_password"})
		return
	}

	hashed, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "密码加密失败"})
		return
	}
	if err := database.DB.Model(&user).Update("password", hashed).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "密码修改失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "密码已修改"})
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
