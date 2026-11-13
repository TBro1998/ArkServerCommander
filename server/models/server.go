package models

import (
	"time"

	"gorm.io/gorm"
)

type Server struct {
	ID            uint   `json:"id" gorm:"primarykey"`
	Identifier    string `json:"identifier" gorm:"not null"`
	SessionName   string `json:"session_name" gorm:"default:'ARK Server'"` // 服务器名称
	ClusterID     string `json:"cluster_id" gorm:"default:''"`             // 集群ID
	Port          int    `json:"port" gorm:"not null;default:7777"`
	QueryPort     int    `json:"query_port" gorm:"not null;default:27015"`
	RCONPort      int    `json:"rcon_port" gorm:"not null;default:32330"`
	AdminPassword string `json:"admin_password" gorm:"not null;default:password"`
	Map           string `json:"map" gorm:"default:'TheIsland'"`
	MaxPlayers    int    `json:"max_players" gorm:"not null;default:70"` // 最大玩家数
	GameModIds    string `json:"game_mod_ids" gorm:"default:''"`         // 游戏模组ID列表，用逗号分隔
	Status        string `json:"status" gorm:"default:'stopped'"`
	LastError     string `json:"last_error" gorm:"default:''"` // 最近一次启停失败的原因，下次启停时清空
	AutoRestart   bool   `json:"auto_restart" gorm:"default:true"`
	AutoUpdate    bool   `json:"auto_update" gorm:"default:false"` // 每次启动前检查游戏本体和 Mod 更新
	// 启用 ArkApi：用 Proton 镜像运行 Windows 版服务端并加载 ArkApi 插件；关闭时用原生 Linux 版服务端。
	// 加这个字段之前创建的服务器在迁移时设为 true（原来都是 Proton 版并带 ArkApi）
	ArkApi    bool           `json:"arkapi" gorm:"column:arkapi;default:false"`
	UserID    uint           `json:"user_id" gorm:"not null"`
	User      User           `json:"user" gorm:"foreignKey:UserID"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
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
	AutoUpdate    *bool  `json:"auto_update"`                         // 每次启动前检查游戏本体和 Mod 更新（可选）
	ArkApi        *bool  `json:"arkapi"`                              // 启用 ArkApi（Proton 镜像），默认关闭（Linux 镜像）
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
	AutoUpdate  bool   `json:"auto_update"`
	ArkApi      bool   `json:"arkapi"` // 启用 ArkApi（Proton 镜像运行 Windows 版服务端）
	Image       string `json:"image"`  // 服务器使用的镜像，由 arkapi 决定
	// 运行中服务器的实时信息（阶段、在线人数、资源占用），由后台监控定期采集
	Runtime *RuntimeInfo `json:"runtime,omitempty"`
	// 有已保存但尚未写入数据卷的配置文件修改，下次启动时生效
	ConfigPending bool `json:"config_pending"`
	// 容器仍在使用旧镜像，下次启动时会重建
	ImageOutdated bool `json:"image_outdated"`
	// 创意工坊上有更新的 Mod ID（更新时间晚于安装时间），由后台定期检查
	ModsOutdated []string `json:"mods_outdated,omitempty"`
	UserID       uint     `json:"user_id"`
	CreatedAt    string   `json:"created_at"`
	UpdatedAt    string   `json:"updated_at"`
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
	AutoUpdate    *bool   `json:"auto_update"`
	ArkApi        *bool   `json:"arkapi"` // 切换后下次启动时重建容器
	// 配置文件内容（可选）
	GameUserSettings string `json:"game_user_settings,omitempty"` // GameUserSettings.ini 文件内容
	GameIni          string `json:"game_ini,omitempty"`           // Game.ini 文件内容
	// 启动参数（可选）
	ServerArgs *ServerArgsRequest `json:"server_args,omitempty"` // 启动参数结构
}

// 运行阶段：容器运行后，入口脚本依次下载服务端、安装 Mod、启动游戏，游戏加载完地图后才能连接
const (
	PhasePullingImage   = "pulling_image"   // 启动前下载服务器需要的镜像（此时状态为 starting）
	PhaseCopying        = "copying"         // 启动前从另一台已安装的服务器复制服务端文件（此时状态为 starting）
	PhaseDownloading    = "downloading"     // 下载或校验服务端文件
	PhaseInstallingMods = "installing_mods" // 下载、安装 Mod
	PhaseLaunching      = "launching"       // 游戏进程已启动，正在加载地图
	PhaseReady          = "ready"           // 查询端口有响应，可以连接
)

// RuntimeInfo 运行中服务器的实时信息
type RuntimeInfo struct {
	Phase       string  `json:"phase"`                  // 见 Phase* 常量，无法判断时为空
	Progress    float64 `json:"progress,omitempty"`     // 下载服务端时的进度（0-100）
	ModCurrent  int     `json:"mod_current,omitempty"`  // 正在处理第几个 Mod
	ModTotal    int     `json:"mod_total,omitempty"`    // Mod 总数
	Players     int     `json:"players"`                // 在线人数（查询端口返回）
	MaxPlayers  int     `json:"max_players,omitempty"`  // 查询端口返回的人数上限
	ServerName  string  `json:"server_name,omitempty"`  // 查询端口返回的名称，带游戏版本号
	CopySource  string  `json:"copy_source,omitempty"`  // copying 阶段的来源服务器名称
	Image       string  `json:"image,omitempty"`        // pulling_image 阶段正在下载的镜像
	CPUPercent  float64 `json:"cpu_percent"`            // CPU 占用率，100 表示占满一个核心
	MemoryBytes uint64  `json:"memory_bytes"`           // 内存占用
	MemoryLimit uint64  `json:"memory_limit,omitempty"` // 内存上限（宿主机内存）
	UpdatedAt   string  `json:"updated_at"`
}
