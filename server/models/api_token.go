package models

import "time"

// APIToken 供 MCP 客户端（Agent）使用的长期访问令牌。数据库只保存令牌的 SHA-256，明文只在创建时返回一次
type APIToken struct {
	ID         uint       `json:"id" gorm:"primarykey"`
	UserID     uint       `json:"-" gorm:"index;not null"`
	Name       string     `json:"name" gorm:"not null"`
	Prefix     string     `json:"prefix" gorm:"not null"`                          // 令牌开头几位，用于在列表中辨认
	TokenHash  string     `json:"-" gorm:"column:token_hash;uniqueIndex;not null"` // 令牌的 SHA-256（十六进制）
	ReadOnly   bool       `json:"read_only" gorm:"not null;default:false"`         // 只读令牌只能调用查询类工具
	LastUsedAt *time.Time `json:"last_used_at"`
	CreatedAt  time.Time  `json:"created_at"`
}

// APITokenRequest 创建令牌请求
type APITokenRequest struct {
	Name     string `json:"name" binding:"required"`
	ReadOnly bool   `json:"read_only"`
}

// APITokenCreated 创建令牌的响应，Token 是明文，之后无法再次获取
type APITokenCreated struct {
	APIToken
	Token string `json:"token"`
}
