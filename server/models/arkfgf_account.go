package models

import "time"

// ArkfgfAccount 绑定的 arkfgf.com 账号。每个管理端最多绑定一个账号（只有一行，ID 固定为 1，未绑定时没有这一行）。
// 发布到收录站和 P2P 联机房都用账号令牌调用网站，见外层仓库 docs/commander-account.md
type ArkfgfAccount struct {
	ID uint `gorm:"primaryKey"`
	// 账号令牌（acm_...），只存在本机，不返回给浏览器
	Token string `gorm:"default:''"`
	// 网站上这个 Commander 的 ID 和名称
	CommanderID   uint   `gorm:"default:0"`
	CommanderName string `gorm:"default:''"`
	// 绑定的 Steam 账号
	SteamID   string `gorm:"default:''"`
	UserName  string `gorm:"default:''"`
	AvatarURL string `gorm:"default:''"`
	// 网站的服主后台地址
	DashboardURL string `gorm:"default:''"`
	BoundAt      time.Time
	UpdatedAt    time.Time
}

// ArkfgfAccountStatus arkfgf.com 账号的绑定状态（设置页）
type ArkfgfAccountStatus struct {
	Enabled      bool               `json:"enabled"`                 // ARKFGF_API 不为空
	Bound        bool               `json:"bound"`                   // 已绑定账号
	Account      *ArkfgfAccountInfo `json:"account,omitempty"`       // 已绑定时的账号信息
	DashboardURL string             `json:"dashboard_url,omitempty"` // 网站的服主后台
	Login        ArkfgfLoginStatus  `json:"login"`                   // 进行中的 Steam 登录
}

// ArkfgfLoginStatus 进行中的 Steam 登录：在 Steam 客户端里登录时，回调页在客户端的内置浏览器里打开，
// 设置页轮询这里知道结果
type ArkfgfLoginStatus struct {
	Pending   bool   `json:"pending"`              // 有开始了、还没完成也没过期的登录
	ExpiresAt string `json:"expires_at,omitempty"` // 最近一次开始的登录的过期时间
	ErrorCode string `json:"error_code,omitempty"` // 最近一次完成登录失败的错误码（开始新的登录或绑定成功时清除）
	Error     string `json:"error,omitempty"`      // 失败原因（中文）
}

// ArkfgfAccountInfo 绑定的账号信息（不含令牌）
type ArkfgfAccountInfo struct {
	CommanderID   uint   `json:"commander_id"`
	CommanderName string `json:"commander_name"` // 这个 Commander 在网站上的名称
	SteamID       string `json:"steam_id"`
	UserName      string `json:"user_name"` // Steam 昵称
	AvatarURL     string `json:"avatar_url"`
	BoundAt       string `json:"bound_at"`
}

// ArkfgfLoginRequest 开始用 Steam 登录绑定
type ArkfgfLoginRequest struct {
	Name     string `json:"name"`      // 这个 Commander 在网站上的名称，为空时用 ArkServerCommander，最多 64 字
	ReturnTo string `json:"return_to"` // Steam 登录后回到的 Commander 页面（origin + /settings/arkfgf-callback），不带 query
}

// ArkfgfLoginResponse 浏览器要跳转的 Steam 登录地址
type ArkfgfLoginResponse struct {
	SteamURL string `json:"steam_url"`
}

// ArkfgfLoginFinishRequest Steam 登录后回调地址上的全部 query 参数（state、openid.*）
type ArkfgfLoginFinishRequest struct {
	Params map[string]string `json:"params" binding:"required"`
}
