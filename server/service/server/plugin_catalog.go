package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"ark-server-commander/config"
	"ark-server-commander/utils"

	"go.uber.org/zap"
)

// 插件目录：ASE 的 ArkApi 插件和 ASA 的 AsaApi 插件都发布在 ArkApi 官方站 ark-server-api.com 上，arkfgf.com 整理成目录
// （GET {ARKFGF_API}/plugins/catalog，ASA 带 ?game=asa，公开接口，不需要绑定账号），Commander 在插件页浏览和一键安装。
// 两个游戏的插件不通用，目录各缓存一份，服务器只用自己游戏的目录。
// 插件文件从 ark-server-api.com 的下载地址直接下载，安装流程和"从地址安装"相同。
// 从目录安装的插件在目录里写一个 .commander-source.json，记下资源 ID 和版本，用来判断有没有更新

const (
	pluginCatalogTTL   = time.Hour
	pluginCatalogRetry = 5 * time.Minute
	pluginSourceFile   = ".commander-source.json"
)

// catalogDownloadHosts 目录里的下载地址只接受这些站点（https）
var catalogDownloadHosts = []string{"ark-server-api.com"}

// allowCatalogDownload 测试时替换成接受 httptest 地址
var allowCatalogDownload = func(u *url.URL) bool {
	if u.Scheme != "https" {
		return false
	}
	host := strings.ToLower(u.Hostname())
	for _, h := range catalogDownloadHosts {
		if host == h {
			return true
		}
	}
	return false
}

// CatalogPluginInfo 资源压缩包里的一个插件（目录名和 PluginInfo.json）
type CatalogPluginInfo struct {
	Name          string   `json:"name"`
	FullName      string   `json:"full_name,omitempty"`
	Version       string   `json:"version,omitempty"`
	MinAPIVersion string   `json:"min_api_version,omitempty"`
	Dependencies  []string `json:"dependencies,omitempty"`
}

// CatalogPlugin 目录里的一个资源
type CatalogPlugin struct {
	ID          uint                `json:"id"`
	Title       string              `json:"title"`
	Tagline     string              `json:"tagline,omitempty"`
	Author      string              `json:"author,omitempty"`
	Category    string              `json:"category"` // official / free / paid
	Version     string              `json:"version,omitempty"`
	Downloads   int                 `json:"downloads"`
	Rating      float64             `json:"rating,omitempty"`
	UpdatedAt   *time.Time          `json:"updated_at,omitempty"`
	PageURL     string              `json:"page_url"`
	DownloadURL string              `json:"download_url,omitempty"` // 为空表示不能一键安装（付费、需要到网站下载）
	FileName    string              `json:"file_name,omitempty"`
	FileSize    int64               `json:"file_size,omitempty"`
	Plugins     []CatalogPluginInfo `json:"plugins,omitempty"`
}

// CatalogArkApi 目录里的框架本体：ASE 目录里是 ArkApi，ASA 目录里是 AsaApi
type CatalogArkApi struct {
	Version     string     `json:"version"`
	DownloadURL string     `json:"download_url"`
	FileName    string     `json:"file_name,omitempty"`
	PageURL     string     `json:"page_url"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
}

// PluginCatalog 插件目录
type PluginCatalog struct {
	Game        string          `json:"game"` // ase / asa
	Source      string          `json:"source"`
	SourceURL   string          `json:"source_url"`
	GeneratedAt time.Time       `json:"generated_at"`
	Plugins     []CatalogPlugin `json:"plugins"`
	ArkApi      *CatalogArkApi  `json:"arkapi,omitempty"`
}

// PluginSource 从目录安装的插件的来源
type PluginSource struct {
	Source      string `json:"source"`
	ResourceID  uint   `json:"resource_id"`
	Version     string `json:"version"`
	InstalledAt string `json:"installed_at"`
}

type pluginCatalogCache struct {
	game   string
	api    func() string
	client *http.Client

	mu      sync.Mutex
	data    *PluginCatalog
	fetched time.Time
	tried   time.Time
	lastErr error
}

var (
	pluginCatalogClient = &http.Client{Timeout: 15 * time.Second}
	asePluginCatalog    = &pluginCatalogCache{game: utils.GameASE, api: func() string { return config.AccountAPI }, client: pluginCatalogClient}
	asaPluginCatalog    = &pluginCatalogCache{game: utils.GameASA, api: func() string { return config.AccountAPI }, client: pluginCatalogClient}
)

// pluginCatalogFor 游戏的插件目录（空值视为 ASE）
func pluginCatalogFor(game string) *pluginCatalogCache {
	if utils.NormalizeGame(game) == utils.GameASA {
		return asaPluginCatalog
	}
	return asePluginCatalog
}

// get 返回目录：缓存 1 小时；获取失败时有旧数据就返回旧数据，5 分钟内不重试
func (c *pluginCatalogCache) get(refresh bool) (*PluginCatalog, error) {
	api := c.api()
	if api == "" {
		return nil, ErrArkfgfDisabled
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	now := time.Now()
	if !refresh && c.data != nil && now.Sub(c.fetched) < pluginCatalogTTL {
		return c.data, nil
	}
	if !refresh && c.lastErr != nil && now.Sub(c.tried) < pluginCatalogRetry {
		if c.data != nil {
			return c.data, nil
		}
		return nil, c.lastErr
	}
	var out PluginCatalog
	endpoint := api + "/plugins/catalog"
	if c.game != utils.GameASE {
		endpoint += "?game=" + c.game
	}
	_, _, _, err := arkfgfRequest(c.client, http.MethodGet, endpoint, "", nil, &out, "arkfgf.com")
	// 不认识 game 参数的网站会返回 ASE 的目录（没有 game 字段），不能拿来给 ASA 服务器用
	if err == nil && utils.NormalizeGame(out.Game) != c.game {
		err = fmt.Errorf("网站返回的是 %s 的插件目录", utils.NormalizeGame(out.Game))
	}
	c.tried = now
	if err != nil {
		utils.Warn("获取插件目录失败", zap.String("game", c.game), zap.Error(err))
		c.lastErr = fmt.Errorf("%w: %v", ErrPluginCatalogUnavailable, err)
		if c.data != nil {
			return c.data, nil
		}
		return nil, c.lastErr
	}
	if out.Plugins == nil {
		out.Plugins = []CatalogPlugin{}
	}
	out.Game = c.game
	c.data, c.fetched, c.lastErr = &out, now, nil
	return c.data, nil
}

// peek 只返回已经缓存的目录，不发请求
func (c *pluginCatalogCache) peek() *PluginCatalog {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.data
}

// forget 测试换数据时清空缓存
func (c *pluginCatalogCache) forget() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data, c.fetched, c.tried, c.lastErr = nil, time.Time{}, time.Time{}, nil
}

// GetPluginCatalog 返回游戏（ase / asa，空为 ase）的插件目录，refresh 时重新向 arkfgf.com 获取
func (s *ServerService) GetPluginCatalog(game string, refresh bool) (*PluginCatalog, error) {
	if !utils.ValidGame(game) {
		return nil, ErrInvalidGame
	}
	return pluginCatalogFor(game).get(refresh)
}

// InstallCatalogPlugin 从目录安装（或更新）一个资源，返回安装的插件名称。只能在服务器停止时进行，
// 覆盖安装保留原来的 config.json 和启用状态
func (s *ServerService) InstallCatalogPlugin(userID uint, serverID string, resourceID uint) ([]string, error) {
	server, err := stoppedServer(userID, serverID)
	if err != nil {
		return nil, err
	}
	// 两个游戏的插件不通用，只在服务器自己游戏的目录里找
	catalog, err := pluginCatalogFor(server.Game).get(false)
	if err != nil {
		return nil, err
	}
	var entry *CatalogPlugin
	for i := range catalog.Plugins {
		if catalog.Plugins[i].ID == resourceID {
			entry = &catalog.Plugins[i]
			break
		}
	}
	if entry == nil {
		return nil, ErrCatalogPluginNotFound
	}
	if entry.DownloadURL == "" {
		return nil, ErrCatalogPluginNotInstallable
	}
	u, err := url.Parse(entry.DownloadURL)
	if err != nil || !allowCatalogDownload(u) {
		return nil, fmt.Errorf("%w: 下载地址不是 ark-server-api.com", ErrCatalogPluginNotInstallable)
	}
	data, err := downloadPluginZip(entry.DownloadURL)
	if err != nil {
		return nil, err
	}
	names, err := s.installPluginsData(userID, serverID, data, "")
	if err != nil {
		return names, err
	}

	source := PluginSource{
		Source:      catalog.Source,
		ResourceID:  entry.ID,
		Version:     entry.Version,
		InstalledAt: time.Now().UTC().Format(time.RFC3339),
	}
	content, _ := json.Marshal(source)
	unlock := lockPlugins(server.ID)
	defer unlock()
	for _, name := range names {
		dir, _, err := findPlugin(server.ID, name)
		if err != nil {
			continue
		}
		if err := os.WriteFile(filepath.Join(dir, pluginSourceFile), content, 0o644); err != nil {
			utils.Warn("记录插件来源失败", zap.String("plugin", name), zap.Error(err))
		}
	}
	utils.Info("从插件目录安装了插件", zap.Uint("server_id", server.ID), zap.Uint("resource_id", entry.ID), zap.Strings("plugins", names))
	return names, nil
}

// readPluginSource 读取插件目录里的来源记录，没有或读不了时返回 nil
func readPluginSource(dir string) *PluginSource {
	data, err := os.ReadFile(filepath.Join(dir, pluginSourceFile))
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			utils.Warn("读取插件来源失败", zap.String("dir", dir), zap.Error(err))
		}
		return nil
	}
	var src PluginSource
	if json.Unmarshal(data, &src) != nil || src.ResourceID == 0 {
		return nil
	}
	return &src
}

// catalogArkApi 游戏的目录里的框架本体（版本号合法、下载地址在允许的站点上时），只看已缓存的目录
func catalogArkApi(game string) *CatalogArkApi {
	c := pluginCatalogFor(game).peek()
	if c == nil || c.ArkApi == nil || !arkApiVersionPattern.MatchString(c.ArkApi.Version) {
		return nil
	}
	u, err := url.Parse(c.ArkApi.DownloadURL)
	if err != nil || !allowCatalogDownload(u) {
		return nil
	}
	return c.ArkApi
}
