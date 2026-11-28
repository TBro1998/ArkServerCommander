package models

import "time"

// 玩家站（ArkPlayerPortal）：Commander 把要给玩家看的内容写成快照，玩家站容器只读这些文件提供网页和接口，
// 见 service/server/portal.go 和子模块 ArkPlayerPortal/docs/snapshot.md

// PortalSettings 玩家站设置，只有 ID 为 1 的一行（loadPortalSettings 不存在时创建）
type PortalSettings struct {
	ID      uint `json:"-" gorm:"primarykey"`
	Enabled bool `json:"enabled" gorm:"default:false"`
	// 玩家站在宿主机上的 TCP 端口，映射到容器的 8080
	Port int `json:"port" gorm:"not null;default:8090"`
	// 玩家访问的地址（服主自己的域名或反向代理），为空时用 http://<公网地址>:<端口>
	PublicURL   string `json:"public_url" gorm:"default:''"`
	Name        string `json:"name" gorm:"default:''"` // 站点名称，为空时玩家站显示第一台服务器的名称
	Description string `json:"description" gorm:"type:text;default:''"`
	QQGroup     string `json:"qq_group" gorm:"default:''"`
	Discord     string `json:"discord" gorm:"default:''"`
	// 链接（[{title, url}] 的 JSON）
	LinksJSON string `json:"-" gorm:"column:links_json;type:text;default:'[]'"`
	// 不在玩家站上显示的服务器 ID，逗号分隔
	HiddenServers string    `json:"-" gorm:"default:''"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// PortalLink 站点资料里的一个链接
type PortalLink struct {
	Title string `json:"title"`
	URL   string `json:"url"`
}

// PortalAnnouncement 玩家站公告
type PortalAnnouncement struct {
	ID          uint      `json:"id" gorm:"primarykey"`
	Title       string    `json:"title" gorm:"not null"`
	Body        string    `json:"body" gorm:"type:text;default:''"`
	Pinned      bool      `json:"pinned" gorm:"default:false"`
	PublishedAt time.Time `json:"published_at"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// PortalSettingsRequest 修改玩家站设置（整体替换）
type PortalSettingsRequest struct {
	Enabled       bool         `json:"enabled"`
	Port          int          `json:"port"`
	PublicURL     string       `json:"public_url"`
	Name          string       `json:"name"`
	Description   string       `json:"description"`
	QQGroup       string       `json:"qq_group"`
	Discord       string       `json:"discord"`
	Links         []PortalLink `json:"links"`
	HiddenServers []uint       `json:"hidden_servers"`
}

// PortalAnnouncementRequest 新建或修改公告
type PortalAnnouncementRequest struct {
	Title  string `json:"title"`
	Body   string `json:"body"`
	Pinned bool   `json:"pinned"`
}

// 玩家站状态
const (
	PortalDisabled = "disabled" // 没有开启
	PortalStarting = "starting" // 正在下载镜像或启动容器
	PortalRunning  = "running"
	PortalError    = "error"
)

// PortalStatus 玩家站容器的状态，只在内存里
type PortalStatus struct {
	State     string `json:"state"`
	URL       string `json:"url,omitempty"` // 玩家访问的地址
	ErrorCode string `json:"error_code,omitempty"`
	Error     string `json:"error,omitempty"`
}

// PortalServer 玩家站设置页里的一台服务器
type PortalServer struct {
	ID     uint   `json:"id"`
	Name   string `json:"name"`
	Game   string `json:"game"`
	Hidden bool   `json:"hidden"`
	// 这台服务器有 ArkShop（启用了 ArkApi，装了 ArkShop 并且没有停用），玩家站显示它的商店
	Shop bool `json:"shop"`
	// 读取 ArkShop 的 config.json 失败的原因
	ShopError string `json:"shop_error,omitempty"`
}

// PortalInfo GET /api/portal 的响应
type PortalInfo struct {
	Settings PortalSettings `json:"settings"`
	Links    []PortalLink   `json:"links"`
	Status   PortalStatus   `json:"status"`
	Servers  []PortalServer `json:"servers"`
	Image    string         `json:"image"`
}

// PortalPage 玩家站上服主写的页面（规则、新手指南等），内容是纯文本：「# 」开头的行是标题，「- 」开头的行是列表
type PortalPage struct {
	ID        uint      `json:"id" gorm:"primarykey"`
	Title     string    `json:"title" gorm:"not null"`
	Body      string    `json:"body" gorm:"type:text;default:''"`
	Sort      int       `json:"sort" gorm:"default:0"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// PortalPageRequest 新建或修改页面
type PortalPageRequest struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}
