package models

import (
	"time"

	"gorm.io/gorm"
)

type Server struct {
	ID            uint           `json:"id" gorm:"primarykey"`
	Name          string         `json:"name" gorm:"not null"`
	Description   string         `json:"description"`
	Port          int            `json:"port" gorm:"not null;default:7777"`
	QueryPort     int            `json:"query_port" gorm:"not null;default:27015"`
	RCONPort      int            `json:"rcon_port" gorm:"not null;default:27020"`
	AdminPassword string         `json:"admin_password" gorm:"not null;default:password"`
	Map           string         `json:"map" gorm:"default:'TheIsland'"`
	MaxPlayers    int            `json:"max_players" gorm:"default:70"`
	Status        string         `json:"status" gorm:"default:'stopped'"`
	UserID        uint           `json:"user_id" gorm:"not null"`
	User          User           `json:"user" gorm:"foreignKey:UserID"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index"`
}

type ServerRequest struct {
	Name          string `json:"name" binding:"required"`
	Description   string `json:"description"`
	Port          int    `json:"port" binding:"required,min=1,max=65535"`
	QueryPort     int    `json:"query_port" binding:"required,min=1,max=65535"`
	RCONPort      int    `json:"rcon_port" binding:"required,min=1,max=65535"`
	AdminPassword string `json:"admin_password" binding:"required"`
	Map           string `json:"map"`
	MaxPlayers    int    `json:"max_players" binding:"required,min=1,max=200"`
}

type ServerResponse struct {
	ID            uint   `json:"id"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	Port          int    `json:"port"`
	QueryPort     int    `json:"query_port"`
	RCONPort      int    `json:"rcon_port"`
	AdminPassword string `json:"admin_password"`
	Map           string `json:"map"`
	MaxPlayers    int    `json:"max_players"`
	Status        string `json:"status"`
	UserID        uint   `json:"user_id"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
}

type ServerUpdateRequest struct {
	Name          string `json:"name"`
	Description   string `json:"description"`
	Port          int    `json:"port" binding:"min=1,max=65535"`
	QueryPort     int    `json:"query_port" binding:"min=1,max=65535"`
	RCONPort      int    `json:"rcon_port" binding:"min=1,max=65535"`
	AdminPassword string `json:"admin_password"`
	Map           string `json:"map"`
	MaxPlayers    int    `json:"max_players" binding:"min=1,max=200"`
}

// ServerConfigRequest 服务器配置文件更新请求
type ServerConfigRequest struct {
	GameUserSettings string `json:"game_user_settings"` // GameUserSettings.ini 文件内容
	GameIni          string `json:"game_ini"`           // Game.ini 文件内容
}

// ServerConfigResponse 服务器配置文件响应
type ServerConfigResponse struct {
	ServerID             uint   `json:"server_id"`
	ServerName           string `json:"server_name"`
	GameUserSettings     string `json:"game_user_settings"`      // GameUserSettings.ini 文件内容
	GameIni              string `json:"game_ini"`                // Game.ini 文件内容
	GameUserSettingsPath string `json:"game_user_settings_path"` // GameUserSettings.ini 文件路径
	GameIniPath          string `json:"game_ini_path"`           // Game.ini 文件路径
	UpdatedAt            string `json:"updated_at"`
}

// ServerConfigFileInfo 服务器配置文件信息
type ServerConfigFileInfo struct {
	FileName     string `json:"file_name"`     // 文件名
	FilePath     string `json:"file_path"`     // 文件完整路径
	FileSize     int64  `json:"file_size"`     // 文件大小（字节）
	LastModified string `json:"last_modified"` // 最后修改时间
	Exists       bool   `json:"exists"`        // 文件是否存在
}
