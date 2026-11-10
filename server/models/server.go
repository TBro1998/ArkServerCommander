package models

import (
	"time"

	"gorm.io/gorm"
)

type Server struct {
	ID            uint           `json:"id" gorm:"primarykey"`
	Identifier    string         `json:"identifier" gorm:"not null"`
	SessionName   string         `json:"session_name" gorm:"default:'ARK Server'"` // 服务器名称
	ClusterID     string         `json:"cluster_id" gorm:"default:''"`             // 集群ID
	Port          int            `json:"port" gorm:"not null;default:7777"`
	QueryPort     int            `json:"query_port" gorm:"not null;default:27015"`
	RCONPort      int            `json:"rcon_port" gorm:"not null;default:32330"`
	AdminPassword string         `json:"admin_password" gorm:"not null;default:password"`
	Map           string         `json:"map" gorm:"default:'TheIsland'"`
	MaxPlayers    int            `json:"max_players" gorm:"not null;default:70"` // 最大玩家数
	GameModIds    string         `json:"game_mod_ids" gorm:"default:''"`         // 游戏模组ID列表，用逗号分隔
	Status        string         `json:"status" gorm:"default:'stopped'"`
	LastError     string         `json:"last_error" gorm:"default:''"` // 最近一次启停失败的原因，下次启停时清空
	AutoRestart   bool           `json:"auto_restart" gorm:"default:true"`
	UserID        uint           `json:"user_id" gorm:"not null"`
	User          User           `json:"user" gorm:"foreignKey:UserID"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index"`
	// 服务器未停止时保存的配置文件，下次启动前写入数据卷后清空（停服时 ARK 会覆盖 GameUserSettings.ini）
	PendingGameUserSettings string `json:"-" gorm:"type:text;default:''"`
	PendingGameIni          string `json:"-" gorm:"type:text;default:''"`
	// 启动参数（JSON格式存储）
	ServerArgsJSON string `json:"server_args_json" gorm:"default:'{}'"` // 启动参数的JSON字符串
}

type ServerRequest struct {
	Identifier    string `json:"identifier" binding:"required"`
	SessionName   string `json:"session_name"` // 服务器名称
	ClusterID     string `json:"cluster_id"`   // 集群ID
	Port          int    `json:"port" binding:"required,min=1,max=65535"`
	QueryPort     int    `json:"query_port" binding:"required,min=1,max=65535"`
	RCONPort      int    `json:"rcon_port" binding:"required,min=1,max=65535"`
	AdminPassword string `json:"admin_password" binding:"required"`
	Map           string `json:"map"`
	MaxPlayers    int    `json:"max_players" binding:"min=1,max=200"` // 最大玩家数
	GameModIds    string `json:"game_mod_ids"`                        // 游戏模组ID列表，用逗号分隔
	AutoRestart   *bool  `json:"auto_restart"`                        // 是否自动重启（可选）
	// 配置文件内容（可选）
	GameUserSettings string `json:"game_user_settings,omitempty"` // GameUserSettings.ini 文件内容
	GameIni          string `json:"game_ini,omitempty"`           // Game.ini 文件内容
	// 启动参数（可选）
	ServerArgs *ServerArgsRequest `json:"server_args,omitempty"`
}

type ServerResponse struct {
	ID          uint   `json:"id"`
	Identifier  string `json:"identifier"`
	SessionName string `json:"session_name"` // 服务器名称
	ClusterID   string `json:"cluster_id"`   // 集群ID
	Port        int    `json:"port"`
	QueryPort   int    `json:"query_port"`
	RCONPort    int    `json:"rcon_port"`
	Map         string `json:"map"`
	MaxPlayers  int    `json:"max_players"` // 最大玩家数
	GameModIds  string `json:"game_mod_ids"`
	Status      string `json:"status"`
	LastError   string `json:"last_error,omitempty"` // 最近一次启停失败的原因
	AutoRestart bool   `json:"auto_restart"`
	// 有已保存但尚未写入数据卷的配置文件修改，下次启动时生效
	ConfigPending bool `json:"config_pending"`
	// 容器仍在使用旧镜像，下次启动时会重建
	ImageOutdated bool   `json:"image_outdated"`
	UserID        uint   `json:"user_id"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
	// 配置文件内容
	GameUserSettings string `json:"game_user_settings,omitempty"` // GameUserSettings.ini 文件内容
	GameIni          string `json:"game_ini,omitempty"`           // Game.ini 文件内容
	// 启动参数
	ServerArgs    *ServerArgs `json:"server_args,omitempty"`    // 启动参数结构
	GeneratedArgs string      `json:"generated_args,omitempty"` // 生成的完整启动参数字符串
}

// ServerUpdateRequest 更新服务器请求：所有字段都是可选的，没传（或为零值）的字段保持不变。
// ClusterID 和 GameModIds 使用指针，传空字符串表示清空
type ServerUpdateRequest struct {
	Identifier    string  `json:"identifier"`
	SessionName   string  `json:"session_name"` // 服务器名称
	ClusterID     *string `json:"cluster_id"`   // 集群ID，传空字符串表示清空
	Port          int     `json:"port" binding:"omitempty,min=1,max=65535"`
	QueryPort     int     `json:"query_port" binding:"omitempty,min=1,max=65535"`
	RCONPort      int     `json:"rcon_port" binding:"omitempty,min=1,max=65535"`
	AdminPassword string  `json:"admin_password"`
	Map           string  `json:"map"`
	MaxPlayers    int     `json:"max_players" binding:"omitempty,min=1,max=200"` // 最大玩家数
	GameModIds    *string `json:"game_mod_ids"`                                  // 游戏模组ID列表，用逗号分隔，传空字符串表示清空
	AutoRestart   *bool   `json:"auto_restart"`
	// 配置文件内容（可选）
	GameUserSettings string `json:"game_user_settings,omitempty"` // GameUserSettings.ini 文件内容
	GameIni          string `json:"game_ini,omitempty"`           // Game.ini 文件内容
	// 启动参数（可选）
	ServerArgs *ServerArgsRequest `json:"server_args,omitempty"` // 启动参数结构
}
