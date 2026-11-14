// Package apitoken 管理 MCP 客户端使用的长期访问令牌
package apitoken

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"strings"
	"time"
	"unicode/utf8"

	"ark-server-commander/database"
	"ark-server-commander/models"

	"gorm.io/gorm"
)

// TokenPrefix 令牌固定前缀，便于识别和密钥扫描
const TokenPrefix = "asc_"

// 每个用户最多保存的令牌数
const maxTokensPerUser = 50

// 距上次记录超过这个时间才更新 last_used_at，避免每次调用都写数据库
const lastUsedInterval = time.Minute

var (
	ErrTokenNotFound    = errors.New("令牌不存在")
	ErrInvalidTokenName = errors.New("令牌名称不能为空，且不超过 64 个字符")
	ErrTooManyTokens    = errors.New("令牌数量已达上限，请先删除不用的令牌")
	ErrInvalidToken     = errors.New("无效的访问令牌")
)

func hashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

// Create 生成新令牌，返回的 Token 字段是明文，只在这里出现一次
func Create(userID uint, name string, readOnly bool) (*models.APITokenCreated, error) {
	name = strings.TrimSpace(name)
	if name == "" || utf8.RuneCountInString(name) > 64 {
		return nil, ErrInvalidTokenName
	}

	var count int64
	if err := database.DB.Model(&models.APIToken{}).Where("user_id = ?", userID).Count(&count).Error; err != nil {
		return nil, err
	}
	if count >= maxTokensPerUser {
		return nil, ErrTooManyTokens
	}

	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return nil, err
	}
	token := TokenPrefix + base64.RawURLEncoding.EncodeToString(buf)

	record := models.APIToken{
		UserID:    userID,
		Name:      name,
		Prefix:    token[:len(TokenPrefix)+6],
		TokenHash: hashToken(token),
		ReadOnly:  readOnly,
	}
	if err := database.DB.Create(&record).Error; err != nil {
		return nil, err
	}
	// ReadOnly 没有带 default:true，false 不会被 GORM 替换，这里不需要写回
	return &models.APITokenCreated{APIToken: record, Token: token}, nil
}

// List 列出用户的令牌（不含明文和哈希），最新的在前
func List(userID uint) ([]models.APIToken, error) {
	tokens := []models.APIToken{}
	err := database.DB.Where("user_id = ?", userID).Order("id DESC").Find(&tokens).Error
	return tokens, err
}

// Delete 撤销令牌，立即失效
func Delete(userID, id uint) error {
	res := database.DB.Where("id = ? AND user_id = ?", id, userID).Delete(&models.APIToken{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrTokenNotFound
	}
	return nil
}

// Verify 校验明文令牌，返回令牌记录。令牌所属用户已删除时视为无效
func Verify(token string) (*models.APIToken, error) {
	if !strings.HasPrefix(token, TokenPrefix) {
		return nil, ErrInvalidToken
	}

	var record models.APIToken
	err := database.DB.Where("token_hash = ?", hashToken(token)).First(&record).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrInvalidToken
	}
	if err != nil {
		return nil, err
	}

	var users int64
	if err := database.DB.Model(&models.User{}).Where("id = ?", record.UserID).Count(&users).Error; err != nil {
		return nil, err
	}
	if users == 0 {
		return nil, ErrInvalidToken
	}

	now := time.Now()
	if record.LastUsedAt == nil || now.Sub(*record.LastUsedAt) > lastUsedInterval {
		// 记录失败不影响本次调用
		database.DB.Model(&models.APIToken{}).Where("id = ?", record.ID).Update("last_used_at", now)
		record.LastUsedAt = &now
	}
	return &record, nil
}
