package models

import "time"

// LinkRoom P2P 联机房。每个管理端只有一个房间（只有一行，ID 固定为 1）：所有开启了 link 的服务器
// 在 arkfgf.com 上是同一个房间里的多台服务器，共用一个 EasyTier 节点容器和凭证池，玩家加入一次就能连所有服务器
type LinkRoom struct {
	ID uint `gorm:"primaryKey"`
	// 旧版网站返回的房间令牌（alr_...）。现在用绑定的 arkfgf.com 账号令牌调用网站，这个令牌只用于登记时让网站
	// 认领旧房间（legacy_token），或没有服务器开启联机房时下架旧房间，之后清空，新代码不再写入
	Token string `gorm:"default:''"`
	// 网站上的房间 ID，不为 0 表示已用账号登记（Token 为空时）
	RoomID uint `gorm:"default:0"`
	// 网站分配的 EasyTier 网络名（ark-xxxxxxxxxxxxxxxx）
	NetworkName string `gorm:"default:''"`
	// EasyTier 网络密钥，本地生成（32 字节随机数 hex），只存在本机
	NetworkSecret string `gorm:"default:''"`
	// 节点 RPC 端口（宿主机 127.0.0.1），见 service/server/link.go 的 assignLinkRPCPort
	RPCPort int `gorm:"column:rpc_port;default:0"`
	// 房间资料（设置页填写）。名称为空时用第一台服务器的名称；简介、标签、QQ 群、Discord 全部为空时用第一台服务器的收录资料
	Name        string `gorm:"default:''"`
	Description string `gorm:"type:text;default:''"`
	Tags        string `gorm:"default:''"` // 逗号分隔
	QQGroup     string `gorm:"column:qq_group;default:''"`
	Discord     string `gorm:"default:''"`
	UpdatedAt   time.Time
}

// LinkRoomSettings 联机房资料（设置页）
type LinkRoomSettings struct {
	Name    string         `json:"name"`    // 房间名称，为空时用第一台服务器的名称
	Listing ListingProfile `json:"listing"` // 简介、标签、QQ 群、Discord；全部为空时用第一台服务器的收录资料
}

// LinkRoomInfo 联机房资料和状态
type LinkRoomInfo struct {
	LinkRoomSettings
	Status  *LinkStatus      `json:"status"`  // 房间状态；没有服务器开启联机房时为 nil
	Servers []LinkRoomMember `json:"servers"` // 房间里的服务器（开启了 link 的服务器）
}

// LinkRoomMember 房间里的一台服务器
type LinkRoomMember struct {
	ID          uint   `json:"id"`
	Game        string `json:"game"` // ase 或 asa
	SessionName string `json:"session_name"`
	Port        int    `json:"port"`
	QueryPort   int    `json:"query_port"`
	Status      string `json:"status"`
	Online      bool   `json:"online"` // 最近一次心跳报告为可以连接
}
