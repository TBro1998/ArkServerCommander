package models

import (
	"time"

	"gorm.io/gorm"
)

type Server struct {
	ID            uint   `json:"id" gorm:"primarykey"`
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
	ArkApi bool `json:"arkapi" gorm:"column:arkapi;default:false"`
	// 游戏版本（Steam 分支）：空为蜗牛海洋正式版，preaquatica 为通配符原版（海洋 DLC 更新前），见 utils.GameBranches。
	// 切换后下次启动时重建容器并重新校验服务端文件
	GameBranch string `json:"game_branch" gorm:"default:''"`
	// 发布到服务器收录站（arkfgf.com）：服务器每次可以连接后用绑定的账号发布一次，关闭发布或删除后取消发布
	Publish bool `json:"publish" gorm:"default:false"`
	// 收录站上有这台服务器的记录（用账号令牌发布成功过、还没取消发布）
	Published bool `json:"-" gorm:"default:false"`
	// 旧版匿名发布时收录站返回的发布令牌（afp_...）：只用于下次发布时让收录站清理旧记录（legacy_token），
	// 或关闭发布 / 删除后取消旧记录，之后清空，新代码不再写入
	PublishToken string `json:"-" gorm:"default:''"`
	// P2P 联机房（arkfgf.com 联机助手）：加入管理端的联机房（见 LinkRoom），没有公网 IP 的玩家也能通过 P2P 连进来
	Link bool `json:"link" gorm:"default:false"`
	// 以下 Link* 字段是旧版"每台服务器一个房间"留下的，只用于升级后下架旧房间（service/server/link.go 的 cleanupLegacy），
	// 下架后清空，新代码不再写入
	LinkToken         string `json:"-" gorm:"column:link_token;default:''"`
	LinkRoomID        uint   `json:"-" gorm:"column:link_room_id;default:0"`
	LinkNetworkName   string `json:"-" gorm:"column:link_network_name;default:''"`
	LinkNetworkSecret string `json:"-" gorm:"column:link_network_secret;default:''"`
	LinkRPCPort       int    `json:"-" gorm:"column:link_rpc_port;default:0"`
	// 收录资料：随发布一起提交，显示在收录站的详情页上。标签用逗号分隔
	ListingDescription string         `json:"-" gorm:"type:text;default:''"`
	ListingTags        string         `json:"-" gorm:"default:''"`
	ListingQQ          string         `json:"-" gorm:"default:''"`
	ListingDiscord     string         `json:"-" gorm:"default:''"`
	UserID             uint           `json:"user_id" gorm:"not null"`
	User               User           `json:"user" gorm:"foreignKey:UserID"`
	CreatedAt          time.Time      `json:"created_at"`
	UpdatedAt          time.Time      `json:"updated_at"`
	DeletedAt          gorm.DeletedAt `gorm:"index"`
	// 服务器未停止时保存的配置文件，下次启动前写入数据卷后清空（停服时 ARK 会覆盖 GameUserSettings.ini）
	PendingGameUserSettings string `json:"-" gorm:"type:text;default:''"`
	PendingGameIni          string `json:"-" gorm:"type:text;default:''"`
	// 启动参数（JSON格式存储）
	ServerArgsJSON string `json:"server_args_json" gorm:"default:'{}'"` // 启动参数的JSON字符串
}

type ServerRequest struct {
	SessionName   string          `json:"session_name"` // 服务器名称
	ClusterID     string          `json:"cluster_id"`   // 集群ID
	Port          int             `json:"port" binding:"required,min=1,max=65535"`
	QueryPort     int             `json:"query_port" binding:"required,min=1,max=65535"`
	RCONPort      int             `json:"rcon_port" binding:"required,min=1,max=65535"`
	AdminPassword string          `json:"admin_password" binding:"required"`
	Map           string          `json:"map"`
	MaxPlayers    int             `json:"max_players" binding:"min=1,max=200"` // 最大玩家数
	GameModIds    string          `json:"game_mod_ids"`                        // 游戏模组ID列表，用逗号分隔
	AutoRestart   *bool           `json:"auto_restart"`                        // 是否自动重启（可选）
	AutoUpdate    *bool           `json:"auto_update"`                         // 每次启动前检查游戏本体和 Mod 更新（可选）
	ArkApi        *bool           `json:"arkapi"`                              // 启用 ArkApi（Proton 镜像），默认关闭（Linux 镜像）
	GameBranch    string          `json:"game_branch"`                         // 游戏版本：空为蜗牛海洋正式版，preaquatica 为通配符原版
	Publish       *bool           `json:"publish"`                             // 发布到服务器收录站，默认关闭
	Link          *bool           `json:"link"`                                // P2P 联机房，默认关闭
	Listing       *ListingProfile `json:"listing,omitempty"`                   // 收录资料（可选）
	// 配置文件内容（可选）
	GameUserSettings string `json:"game_user_settings,omitempty"` // GameUserSettings.ini 文件内容
	GameIni          string `json:"game_ini,omitempty"`           // Game.ini 文件内容
	// 启动参数（可选）
	ServerArgs *ServerArgsRequest `json:"server_args,omitempty"`
}

type ServerResponse struct {
	ID          uint   `json:"id"`
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
	ArkApi      bool   `json:"arkapi"`      // 启用 ArkApi（Proton 镜像运行 Windows 版服务端）
	Image       string `json:"image"`       // 服务器使用的镜像，由 arkapi 决定
	GameBranch  string `json:"game_branch"` // 游戏版本：空为蜗牛海洋正式版，preaquatica 为通配符原版
	Publish     bool   `json:"publish"`     // 发布到服务器收录站
	// 收录资料：简介、标签、QQ 群、Discord
	Listing ListingProfile `json:"listing"`
	// 发布状态，未开启发布时为空
	PublishStatus *PublishStatus `json:"publish_status,omitempty"`
	Link          bool           `json:"link"` // P2P 联机房
	// 联机房状态，未开启联机房时为空
	LinkStatus *LinkStatus `json:"link_status,omitempty"`
	// 公网可达性：管理端用公网地址查询这台服务器的查询端口的结果。未运行时为空
	Reachability *Reachability `json:"reachability,omitempty"`
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
	ArkApi        *bool   `json:"arkapi"`      // 切换后下次启动时重建容器
	GameBranch    *string `json:"game_branch"` // 游戏版本，传空字符串表示蜗牛海洋正式版；切换后下次启动时重新校验服务端文件，只下载差异
	Publish       *bool   `json:"publish"`     // 发布到服务器收录站
	Link          *bool   `json:"link"`        // P2P 联机房，立即生效
	// 收录资料，传了就整体替换；运行中的服务器会重新发布
	Listing *ListingProfile `json:"listing,omitempty"`
	// 配置文件内容（可选）
	GameUserSettings string `json:"game_user_settings,omitempty"` // GameUserSettings.ini 文件内容
	GameIni          string `json:"game_ini,omitempty"`           // Game.ini 文件内容
	// 启动参数（可选）
	ServerArgs *ServerArgsRequest `json:"server_args,omitempty"` // 启动参数结构
}

// ListingProfile 收录资料，显示在 arkfgf.com 的服务器详情页上
type ListingProfile struct {
	Description string   `json:"description"` // 简介，最多 500 字，可以换行
	Tags        []string `json:"tags"`        // 标签，最多 8 个，每个最多 16 字
	QQGroup     string   `json:"qq_group"`    // QQ 群号
	Discord     string   `json:"discord"`     // Discord 邀请链接（https://discord.gg/...）
}

// PublishStats 收录站上的访问统计，由收录站按账号令牌和服务器 key 返回
type PublishStats struct {
	ID    uint              `json:"id"`
	URL   string            `json:"url"`
	Badge string            `json:"badge"`
	Week  PublishStatTotals `json:"week"`
	Month PublishStatTotals `json:"month"`
	Days  []PublishStatDay  `json:"days"`
}

type PublishStatTotals struct {
	Views    int `json:"views"`
	Connects int `json:"connects"`
	Copies   int `json:"copies"`
}

type PublishStatDay struct {
	Day      string `json:"day"`
	Views    int    `json:"views"`
	Connects int    `json:"connects"`
	Copies   int    `json:"copies"`
}

// 发布状态
const (
	PublishDisabled  = "disabled"  // 管理端关闭了发布功能（ARKFGF_API 为空）
	PublishWaiting   = "waiting"   // 等待服务器启动完成（可以连接后才发布）
	PublishPublished = "published" // 已发布，收录站上显示为在线
	PublishError     = "error"     // 最近一次发布或心跳失败，会自动重试
	// 还没有绑定 arkfgf.com 账号，不发布（在设置页绑定后自动开始）
	PublishLoginRequired = "login_required"
)

// PublishStatus 发布到服务器收录站的状态
type PublishStatus struct {
	State       string `json:"state"`                  // 见 Publish* 常量
	URL         string `json:"url,omitempty"`          // 收录站页面
	Badge       string `json:"badge,omitempty"`        // 状态徽章图片
	ErrorCode   string `json:"error_code,omitempty"`   // 收录站返回的错误码，或 network
	Error       string `json:"error,omitempty"`        // 错误详情（中文）
	PublishedAt string `json:"published_at,omitempty"` // 最近一次发布或心跳成功的时间
	RetryAt     string `json:"retry_at,omitempty"`     // 失败后下一次重试的时间
}

// 联机房状态
const (
	LinkDisabled = "disabled" // 管理端关闭了与 arkfgf.com 的连接（ARKFGF_API 为空）
	LinkWaiting  = "waiting"  // 还没有登记房间，或服务器正在启动、联机节点正在准备
	LinkOnline   = "online"   // 服务器可以连接，联机节点在运行，心跳正常
	LinkOffline  = "offline"  // 房间已登记，服务器未运行（网站上显示离线）
	LinkError    = "error"    // 最近一次请求或联机节点操作失败，会自动重试
	// 还没有绑定 arkfgf.com 账号，不登记房间、不运行联机节点（在设置页绑定后自动开始）
	LinkLoginRequired = "login_required"
)

// LinkStatus P2P 联机房的状态
type LinkStatus struct {
	State     string `json:"state"`                // 见 Link* 常量
	RoomURL   string `json:"room_url,omitempty"`   // 网站上的房间页面
	ErrorCode string `json:"error_code,omitempty"` // 网站返回的错误码，或 network / node_error / image_error
	Error     string `json:"error,omitempty"`      // 错误详情（中文）
	RetryAt   string `json:"retry_at,omitempty"`   // 请求失败后下一次重试的时间
}

// 公网可达性检测结果
const (
	ReachChecking    = "checking"    // 等待检测或正在检测
	ReachReachable   = "reachable"   // 从公网能查询到这台服务器
	ReachUnreachable = "unreachable" // 从公网查询不到：没有公网 IP、没做端口转发或防火墙没放行
	ReachMismatch    = "mismatch"    // 公网上这个端口回应的是另一台服务器（端口转发到了别的机器）
	ReachNotPublic   = "not_public"  // 连接地址（PUBLIC_HOST）是内网地址，公网上的玩家连不上
	ReachError       = "error"       // 没拿到公网地址，会自动重试
)

// Reachability 公网可达性：管理端向玩家连接的地址（PUBLIC_HOST 或自动探测的公网 IP）的查询端口发 A2S_INFO
type Reachability struct {
	State     string `json:"state"`                // 见 Reach* 常量
	IP        string `json:"ip,omitempty"`         // 检测的 IPv4（PUBLIC_HOST 是域名时为解析结果）
	Name      string `json:"name,omitempty"`       // 公网上查询到的服务器名称（mismatch 时用来提示是哪台）
	CheckedAt string `json:"checked_at,omitempty"` // 最近一次得出结论的时间
	ErrorCode string `json:"error_code,omitempty"` // public_ip_unknown
	Error     string `json:"error,omitempty"`      // 错误详情（中文）
	RetryAt   string `json:"retry_at,omitempty"`   // 下一次重试的时间
}

// LinkMember 联机房成员：当前持有未吊销凭证的玩家
type LinkMember struct {
	SteamID      string `json:"steam_id"`
	Name         string `json:"name"`
	AvatarURL    string `json:"avatar_url"`
	CredentialID string `json:"credential_id"`
	JoinedAt     any    `json:"joined_at"` // 网站返回的时间，原样转发
}

// LinkBan 联机房封禁记录
type LinkBan struct {
	SteamID   string `json:"steam_id"`
	Name      string `json:"name"`
	AvatarURL string `json:"avatar_url"`
	BannedAt  any    `json:"banned_at"` // 网站返回的时间，原样转发
}

// LinkMembers 联机房成员和封禁列表
type LinkMembers struct {
	Members []LinkMember `json:"members"`
	Bans    []LinkBan    `json:"bans"`
}

// LinkKickRequest 踢出联机房成员
type LinkKickRequest struct {
	SteamID string `json:"steam_id" binding:"required"` // 玩家的 SteamID64
	Ban     bool   `json:"ban"`                         // 同时封禁
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
