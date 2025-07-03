package models

import (
	"time"

	"gorm.io/gorm"
)

// ServerArgs ARK服务器命令行参数结构体
type ServerArgs struct {
	// 基础服务器参数
	SessionName             string `json:"session_name,omitempty"`                // 服务器会话名称
	ServerPassword          string `json:"server_password,omitempty"`             // 服务器密码
	ServerAdminPassword     string `json:"server_admin_password,omitempty"`       // 管理员密码
	MaxPlayers              int    `json:"max_players,omitempty"`                 // 最大玩家数
	ReservedPlayerSlots     int    `json:"reserved_player_slots,omitempty"`       // 保留玩家槽位
	QueryPort               int    `json:"query_port,omitempty"`                  // 查询端口
	Port                    int    `json:"port,omitempty"`                        // 游戏端口
	RCONPort                int    `json:"rcon_port,omitempty"`                   // RCON端口
	RCONEnabled             bool   `json:"rcon_enabled,omitempty"`                // 启用RCON
	RCONServerGameLogBuffer int    `json:"rcon_server_game_log_buffer,omitempty"` // RCON游戏日志缓冲区

	// 游戏设置
	DifficultyOffset                   float64 `json:"difficulty_offset,omitempty"`                     // 难度偏移
	OverrideOfficialDifficulty         float64 `json:"override_official_difficulty,omitempty"`          // 覆盖官方难度
	ForceAllowCaveFlyers               bool    `json:"force_allow_cave_flyers,omitempty"`               // 强制允许洞穴飞行生物
	AllowCaveBuildingPvE               bool    `json:"allow_cave_building_pve,omitempty"`               // PvE允许洞穴建筑
	AllowCaveBuildingPvP               bool    `json:"allow_cave_building_pvp,omitempty"`               // PvP允许洞穴建筑
	DisableStructurePlacementCollision bool    `json:"disable_structure_placement_collision,omitempty"` // 禁用建筑放置碰撞

	// 玩家设置
	ShowMapPlayerLocation    bool `json:"show_map_player_location,omitempty"`    // 显示地图玩家位置
	AllowThirdPersonPlayer   bool `json:"allow_third_person_player,omitempty"`   // 允许第三人称
	AlwaysNotifyPlayerLeft   bool `json:"always_notify_player_left,omitempty"`   // 总是通知玩家离开
	AlwaysNotifyPlayerJoined bool `json:"always_notify_player_joined,omitempty"` // 总是通知玩家加入
	ServerCrosshair          bool `json:"server_crosshair,omitempty"`            // 服务器准星
	ServerForceNoHUD         bool `json:"server_force_no_hud,omitempty"`         // 强制无HUD

	// 建筑和结构
	TheMaxStructuresInRange         int `json:"the_max_structures_in_range,omitempty"`         // 范围内最大建筑数
	MaxPlatformSaddleStructureLimit int `json:"max_platform_saddle_structure_limit,omitempty"` // 平台鞍座最大建筑限制

	// 游戏规则
	EnablePvPGamma     bool `json:"enable_pvp_gamma,omitempty"`      // 启用PvP伽马
	DisablePvEGamma    bool `json:"disable_pve_gamma,omitempty"`     // 禁用PvE伽马
	AllowFlyerCarryPvE bool `json:"allow_flyer_carry_pve,omitempty"` // PvE允许飞行生物携带
	AllowFlyerCarryPvP bool `json:"allow_flyer_carry_pvp,omitempty"` // PvP允许飞行生物携带

	// 模组支持
	ActiveMods string `json:"active_mods,omitempty"` // 激活的模组ID列表（逗号分隔）
	ModMapID   int    `json:"mod_map_id,omitempty"`  // 模组地图ID

	// 性能设置
	UseBattlEye             bool   `json:"use_battleye,omitempty"`               // 使用BattlEye
	NoBattlEye              bool   `json:"no_battleye,omitempty"`                // 不使用BattlEye
	NoTransferFromFiltering bool   `json:"no_transfer_from_filtering,omitempty"` // 禁用传输过滤
	ClusterDirOverride      string `json:"cluster_dir_override,omitempty"`       // 集群目录覆盖

	// 自定义参数（用于存储其他未列出的参数）
	CustomArgs map[string]string `json:"custom_args,omitempty"` // 自定义参数
}

type Server struct {
	ID            uint   `json:"id" gorm:"primarykey"`
	Identifier    string `json:"identifier" gorm:"not null"`
	Port          int    `json:"port" gorm:"not null;default:7777"`
	QueryPort     int    `json:"query_port" gorm:"not null;default:27015"`
	RCONPort      int    `json:"rcon_port" gorm:"not null;default:32330"`
	AdminPassword string `json:"admin_password" gorm:"not null;default:password"`
	Map           string `json:"map" gorm:"default:'TheIsland'"`
	Status        string `json:"status" gorm:"default:'stopped'"`
	AutoRestart   bool   `json:"auto_restart" gorm:"default:true"`
	UserID        uint   `json:"user_id" gorm:"not null"`
	User          User   `json:"user" gorm:"foreignKey:UserID"`
	// 新增字段
	ServerArgs   string         `json:"server_args" gorm:"type:text"`      // SERVER_ARGS 参数（JSON格式）
	UpdateServer bool           `json:"update_server" gorm:"default:true"` // UPDATE_SERVER 是否自动更新服务器
	UpdateMods   bool           `json:"update_mods" gorm:"default:true"`   // UPDATE_MODS 是否自动更新模组
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}

type ServerRequest struct {
	Identifier    string      `json:"identifier" binding:"required"`
	Port          int         `json:"port" binding:"required,min=1,max=65535"`
	QueryPort     int         `json:"query_port" binding:"required,min=1,max=65535"`
	RCONPort      int         `json:"rcon_port" binding:"required,min=1,max=65535"`
	AdminPassword string      `json:"admin_password" binding:"required"`
	Map           string      `json:"map"`
	AutoRestart   *bool       `json:"auto_restart"`            // 是否自动重启（可选）
	ServerArgs    *ServerArgs `json:"server_args,omitempty"`   // 服务器参数
	UpdateServer  *bool       `json:"update_server,omitempty"` // 是否自动更新服务器
	UpdateMods    *bool       `json:"update_mods,omitempty"`   // 是否自动更新模组
}

type ServerResponse struct {
	ID            uint        `json:"id"`
	Identifier    string      `json:"identifier"`
	Port          int         `json:"port"`
	QueryPort     int         `json:"query_port"`
	RCONPort      int         `json:"rcon_port"`
	AdminPassword string      `json:"admin_password"`
	Map           string      `json:"map"`
	Status        string      `json:"status"`
	AutoRestart   bool        `json:"auto_restart"`
	UserID        uint        `json:"user_id"`
	ServerArgs    *ServerArgs `json:"server_args,omitempty"` // 服务器参数
	UpdateServer  bool        `json:"update_server"`         // 是否自动更新服务器
	UpdateMods    bool        `json:"update_mods"`           // 是否自动更新模组
	CreatedAt     string      `json:"created_at"`
	UpdatedAt     string      `json:"updated_at"`
	// 配置文件内容
	GameUserSettings string `json:"game_user_settings,omitempty"` // GameUserSettings.ini 文件内容
	GameIni          string `json:"game_ini,omitempty"`           // Game.ini 文件内容
}

type ServerUpdateRequest struct {
	Identifier    string      `json:"identifier"`
	Port          int         `json:"port" binding:"min=1,max=65535"`
	QueryPort     int         `json:"query_port" binding:"min=1,max=65535"`
	RCONPort      int         `json:"rcon_port" binding:"min=1,max=65535"`
	AdminPassword string      `json:"admin_password"`
	Map           string      `json:"map"`
	AutoRestart   *bool       `json:"auto_restart"`
	ServerArgs    *ServerArgs `json:"server_args,omitempty"`   // 服务器参数
	UpdateServer  *bool       `json:"update_server,omitempty"` // 是否自动更新服务器
	UpdateMods    *bool       `json:"update_mods,omitempty"`   // 是否自动更新模组
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
