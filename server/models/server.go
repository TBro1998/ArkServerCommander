package models

import (
	"time"

	"gorm.io/gorm"
)

type Server struct {
	ID            uint           `json:"id" gorm:"primarykey"`
	Identifier    string         `json:"identifier" gorm:"not null"`
	Port          int            `json:"port" gorm:"not null;default:7777"`
	QueryPort     int            `json:"query_port" gorm:"not null;default:27015"`
	RCONPort      int            `json:"rcon_port" gorm:"not null;default:32330"`
	AdminPassword string         `json:"admin_password" gorm:"not null;default:password"`
	Map           string         `json:"map" gorm:"default:'TheIsland'"`
	GameModIds    string         `json:"game_mod_ids" gorm:"default:''"` // 游戏模组ID列表，用逗号分隔
	Status        string         `json:"status" gorm:"default:'stopped'"`
	AutoRestart   bool           `json:"auto_restart" gorm:"default:true"`
	UserID        uint           `json:"user_id" gorm:"not null"`
	User          User           `json:"user" gorm:"foreignKey:UserID"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index"`
}

type ServerRequest struct {
	Identifier    string `json:"identifier" binding:"required"`
	Port          int    `json:"port" binding:"required,min=1,max=65535"`
	QueryPort     int    `json:"query_port" binding:"required,min=1,max=65535"`
	RCONPort      int    `json:"rcon_port" binding:"required,min=1,max=65535"`
	AdminPassword string `json:"admin_password" binding:"required"`
	Map           string `json:"map"`
	GameModIds    string `json:"game_mod_ids"` // 游戏模组ID列表，用逗号分隔
	AutoRestart   *bool  `json:"auto_restart"` // 是否自动重启（可选）
}

type ServerResponse struct {
	ID            uint   `json:"id"`
	Identifier    string `json:"identifier"`
	Port          int    `json:"port"`
	QueryPort     int    `json:"query_port"`
	RCONPort      int    `json:"rcon_port"`
	AdminPassword string `json:"admin_password"`
	Map           string `json:"map"`
	GameModIds    string `json:"game_mod_ids"`
	Status        string `json:"status"`
	AutoRestart   bool   `json:"auto_restart"`
	UserID        uint   `json:"user_id"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
	// 配置文件内容
	GameUserSettings string `json:"game_user_settings,omitempty"` // GameUserSettings.ini 文件内容
	GameIni          string `json:"game_ini,omitempty"`           // Game.ini 文件内容
}

type ServerUpdateRequest struct {
	Identifier    string `json:"identifier"`
	Port          int    `json:"port" binding:"min=1,max=65535"`
	QueryPort     int    `json:"query_port" binding:"min=1,max=65535"`
	RCONPort      int    `json:"rcon_port" binding:"min=1,max=65535"`
	AdminPassword string `json:"admin_password"`
	Map           string `json:"map"`
	GameModIds    string `json:"game_mod_ids"` // 游戏模组ID列表，用逗号分隔
	AutoRestart   *bool  `json:"auto_restart"`
	// 配置文件内容（可选）
	GameUserSettings string `json:"game_user_settings,omitempty"` // GameUserSettings.ini 文件内容
	GameIni          string `json:"game_ini,omitempty"`           // Game.ini 文件内容
}

// ServerConfigFileInfo 服务器配置文件信息
type ServerConfigFileInfo struct {
	FileName     string `json:"file_name"`     // 文件名
	FilePath     string `json:"file_path"`     // 文件完整路径
	FileSize     int64  `json:"file_size"`     // 文件大小（字节）
	LastModified string `json:"last_modified"` // 最后修改时间
	Exists       bool   `json:"exists"`        // 文件是否存在
}
