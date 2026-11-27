package server

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	"ark-server-commander/database"
	"ark-server-commander/models"
	"ark-server-commander/service/docker_manager"
	"ark-server-commander/utils"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// 玩家站（ArkPlayerPortal，子模块）：管理端为玩家开一个只读网页，显示公告、服务器状态和 ArkShop 商店。
// PortalManager 每 portalInterval 把要给玩家看的内容写成快照（数据目录 portal/，格式见 ArkPlayerPortal/docs/snapshot.md），
// 并确保玩家站容器在运行。玩家站容器只读这个目录，和管理端隔离（见 docker_manager/docker_portal.go）。
// 内容只在变化时重写：玩家站发现文件变了，就给打开着的页面推送更新

const (
	portalInterval        = 15 * time.Second
	portalHashLabel       = "com.arkservercommander.portal.hash"
	defaultPortalPort     = 8090
	maxPortalNameRunes    = 64
	maxPortalDescRunes    = 2000
	maxPortalLinks        = 10
	maxPortalURLLength    = 512
	maxAnnouncementTitle  = 100
	maxAnnouncementBody   = 10000
	maxPortalAnnouncement = 200
)

var (
	ErrInvalidPortalSettings = errors.New("玩家站设置无效")
	ErrAnnouncementNotFound  = errors.New("公告不存在")
	ErrInvalidAnnouncement   = errors.New("公告无效：标题不能为空、最多 100 字，内容最多 10000 字")
	ErrTooManyAnnouncements  = errors.New("公告最多 200 条，请先删除旧的公告")
)

// portalDocker 玩家站用到的 Docker 操作，测试时替换
type portalDocker interface {
	// ImageReady 镜像已存在时返回 true；不存在时在后台开始下载并返回 false
	ImageReady(image string) (bool, error)
	ImageID(image string) (string, error)
	Inspect(name string) (docker_manager.LinkContainerInfo, error)
	Create(spec docker_manager.PortalContainerSpec) error
	Remove() error
}

// PortalManager 管理玩家站：写快照、玩家站容器
type PortalManager struct {
	mu     sync.Mutex
	status models.PortalStatus
	dir    func() string
	docker portalDocker
	now    func() time.Time
	wakeCh chan struct{}

	// 写入的快照内容摘要（不含 updated_at），没变化时不重写；key 是相对路径
	written map[string]string
	// 最近一次查到的游戏版本，停服后保留
	versions map[uint]string
}

var portalManager = &PortalManager{
	status:   models.PortalStatus{State: models.PortalDisabled},
	dir:      docker_manager.PortalDataDir,
	docker:   dockerPortalOps{},
	now:      time.Now,
	wakeCh:   make(chan struct{}, 1),
	written:  map[string]string{},
	versions: map[uint]string{},
}

// StartPortalManager 启动玩家站后台协程
func StartPortalManager() {
	go func() {
		t := time.NewTicker(portalInterval)
		defer t.Stop()
		for {
			portalManager.tick()
			select {
			case <-t.C:
			case <-portalManager.wakeCh:
			}
		}
	}()
}

// wake 设置或公告变化后马上处理一轮
func (m *PortalManager) wake() {
	select {
	case m.wakeCh <- struct{}{}:
	default:
	}
}

func loadPortalSettings() (models.PortalSettings, error) {
	var s models.PortalSettings
	err := database.DB.Where(models.PortalSettings{ID: 1}).Attrs(models.PortalSettings{Port: defaultPortalPort, LinksJSON: "[]"}).FirstOrCreate(&s).Error
	return s, err
}

func (m *PortalManager) tick() {
	settings, err := loadPortalSettings()
	if err != nil {
		utils.Warn("读取玩家站设置失败", zap.Error(err))
		return
	}
	if !settings.Enabled {
		m.disable()
		return
	}
	if err := m.writeSnapshot(settings); err != nil {
		utils.Warn("写入玩家站快照失败", zap.Error(err))
		m.setStatus(models.PortalStatus{State: models.PortalError, ErrorCode: "local_error", Error: "写入玩家站数据失败：" + err.Error()})
		return
	}
	m.ensureContainer(settings)
}

// disable 关闭玩家站：删除容器和快照
func (m *PortalManager) disable() {
	m.mu.Lock()
	wasOn := m.status.State != models.PortalDisabled
	m.mu.Unlock()
	if err := m.docker.Remove(); err != nil {
		utils.Warn("删除玩家站容器失败", zap.Error(err))
		return
	}
	if wasOn {
		if err := os.RemoveAll(m.dir()); err != nil {
			utils.Warn("删除玩家站数据失败", zap.Error(err))
		}
	}
	m.mu.Lock()
	m.status = models.PortalStatus{State: models.PortalDisabled}
	m.written = map[string]string{}
	m.mu.Unlock()
}

func (m *PortalManager) setStatus(st models.PortalStatus) {
	m.mu.Lock()
	m.status = st
	m.mu.Unlock()
}

// Status 玩家站状态
func (m *PortalManager) Status() models.PortalStatus {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.status
}

// ensureContainer 确保玩家站容器按当前的镜像和端口在运行，镜像或端口变化时重建
func (m *PortalManager) ensureContainer(settings models.PortalSettings) {
	url := portalURL(settings)
	ready, err := m.docker.ImageReady(utils.PortalImage)
	if err != nil {
		m.setStatus(models.PortalStatus{State: models.PortalError, URL: url, ErrorCode: "image_error", Error: "检查玩家站镜像失败：" + err.Error()})
		return
	}
	if !ready {
		m.setStatus(models.PortalStatus{State: models.PortalStarting, URL: url})
		return
	}
	imageID, err := m.docker.ImageID(utils.PortalImage)
	if err != nil {
		m.setStatus(models.PortalStatus{State: models.PortalError, URL: url, ErrorCode: "image_error", Error: "检查玩家站镜像失败：" + err.Error()})
		return
	}
	sum := sha256.Sum256([]byte(utils.PortalImage + "\x00" + imageID + "\x00" + strconv.Itoa(settings.Port)))
	hash := hex.EncodeToString(sum[:])

	info, err := m.docker.Inspect(utils.PortalContainerName)
	if err != nil {
		m.setStatus(models.PortalStatus{State: models.PortalError, URL: url, ErrorCode: "container_error", Error: err.Error()})
		return
	}
	if info.Exists && info.Running && info.Labels[portalHashLabel] == hash {
		m.setStatus(models.PortalStatus{State: models.PortalRunning, URL: url})
		return
	}
	spec := docker_manager.PortalContainerSpec{
		Image:    utils.PortalImage,
		HostPort: settings.Port,
		Labels:   map[string]string{portalHashLabel: hash},
	}
	if err := m.docker.Create(spec); err != nil {
		m.setStatus(models.PortalStatus{State: models.PortalError, URL: url, ErrorCode: "container_error", Error: err.Error()})
		return
	}
	m.setStatus(models.PortalStatus{State: models.PortalRunning, URL: url})
}

// portalURL 玩家访问玩家站的地址
func portalURL(s models.PortalSettings) string {
	if s.PublicURL != "" {
		return s.PublicURL
	}
	host := PublicHost()
	if host == "" {
		return ""
	}
	return "http://" + hostPort(host, s.Port)
}

func hostPort(host string, port int) string {
	if strings.Contains(host, ":") && !strings.HasPrefix(host, "[") {
		host = "[" + host + "]"
	}
	return host + ":" + strconv.Itoa(port)
}

// ---- 快照 ----

type portalSnapshot struct {
	UpdatedAt     string                     `json:"updated_at"`
	Site          portalSite                 `json:"site"`
	Announcements []portalSnapshotAnnounce   `json:"announcements"`
	Servers       []portalSnapshotServerInfo `json:"servers"`
}

type portalSite struct {
	Name        string              `json:"name"`
	Description string              `json:"description"`
	QQGroup     string              `json:"qq_group"`
	Discord     string              `json:"discord"`
	Links       []models.PortalLink `json:"links"`
}

type portalSnapshotAnnounce struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Body        string `json:"body"`
	Pinned      bool   `json:"pinned"`
	PublishedAt string `json:"published_at"`
	UpdatedAt   string `json:"updated_at"`
}

type portalRates struct {
	XP       float64 `json:"xp"`
	Taming   float64 `json:"taming"`
	Harvest  float64 `json:"harvest"`
	Breeding float64 `json:"breeding"`
}

type portalSnapshotServerInfo struct {
	Key           string      `json:"key"`
	Game          string      `json:"game"`
	Name          string      `json:"name"`
	Map           string      `json:"map"`
	Status        string      `json:"status"`
	Players       int         `json:"players"`
	MaxPlayers    int         `json:"max_players"`
	Password      bool        `json:"password"`
	PvP           bool        `json:"pvp"`
	ArkApi        bool        `json:"arkapi"`
	ClusterGroup  int         `json:"cluster_group"`
	Mods          []string    `json:"mods"`
	Rates         portalRates `json:"rates"`
	Version       string      `json:"version"`
	Address       string      `json:"address"`
	QueryPort     int         `json:"query_port"`
	RoomID        uint        `json:"room_id"`
	NextRestartAt *string     `json:"next_restart_at"`
	Shop          bool        `json:"shop"`
}

// portalServers 玩家站上显示的服务器（按 ID，去掉隐藏的）
func portalServers(settings models.PortalSettings) ([]models.Server, error) {
	var servers []models.Server
	if err := database.DB.Order("id").Find(&servers).Error; err != nil {
		return nil, err
	}
	hidden := portalHidden(settings)
	return slices.DeleteFunc(servers, func(s models.Server) bool { return hidden[s.ID] }), nil
}

func portalHidden(settings models.PortalSettings) map[uint]bool {
	hidden := map[uint]bool{}
	for p := range strings.SplitSeq(settings.HiddenServers, ",") {
		if id, err := strconv.ParseUint(strings.TrimSpace(p), 10, 64); err == nil {
			hidden[uint(id)] = true
		}
	}
	return hidden
}

func portalLinks(settings models.PortalSettings) []models.PortalLink {
	var links []models.PortalLink
	_ = json.Unmarshal([]byte(settings.LinksJSON), &links)
	if links == nil {
		links = []models.PortalLink{}
	}
	return links
}

// writeSnapshot 生成快照，有变化的文件才重写（先写临时文件再改名），删掉不再需要的商店文件
func (m *PortalManager) writeSnapshot(settings models.PortalSettings) error {
	servers, err := portalServers(settings)
	if err != nil {
		return fmt.Errorf("读取服务器列表失败: %w", err)
	}
	var announcements []models.PortalAnnouncement
	if err := database.DB.Order("published_at desc, id desc").Find(&announcements).Error; err != nil {
		return fmt.Errorf("读取公告失败: %w", err)
	}
	var room models.LinkRoom
	_ = database.DB.Limit(1).Find(&room, 1).Error

	snap := portalSnapshot{
		Site: portalSite{
			Name:        settings.Name,
			Description: settings.Description,
			QQGroup:     settings.QQGroup,
			Discord:     settings.Discord,
			Links:       portalLinks(settings),
		},
		Announcements: []portalSnapshotAnnounce{},
		Servers:       []portalSnapshotServerInfo{},
	}
	for _, a := range announcements {
		snap.Announcements = append(snap.Announcements, portalSnapshotAnnounce{
			ID:          strconv.FormatUint(uint64(a.ID), 10),
			Title:       a.Title,
			Body:        a.Body,
			Pinned:      a.Pinned,
			PublishedAt: a.PublishedAt.UTC().Format(time.RFC3339),
			UpdatedAt:   a.UpdatedAt.UTC().Format(time.RFC3339),
		})
	}

	shops := map[string]*portalShop{}
	groups := map[string]int{}
	nextRestart := portalNextRestarts()
	for _, s := range servers {
		info := m.serverInfo(s, room, groups, nextRestart)
		if s.ArkApi {
			if data, err := readArkShopConfig(s.ID); err == nil {
				if shop, err := buildPortalShop(data); err == nil {
					shops[info.Key] = shop
					info.Shop = true
				}
			}
		}
		snap.Servers = append(snap.Servers, info)
	}

	dir := m.dir()
	if err := os.MkdirAll(filepath.Join(dir, "shops"), 0o755); err != nil {
		return err
	}
	now := m.now().UTC().Format(time.RFC3339)
	if err := m.writeIfChanged("portal.json", &snap, func() { snap.UpdatedAt = now }); err != nil {
		return err
	}
	for key, shop := range shops {
		if err := m.writeIfChanged("shops/"+key+".json", shop, func() { shop.UpdatedAt = now }); err != nil {
			return err
		}
	}
	// 删除不再需要的商店文件（服务器隐藏、删除，ArkShop 停用、卸载）
	entries, _ := os.ReadDir(filepath.Join(dir, "shops"))
	for _, e := range entries {
		key := strings.TrimSuffix(e.Name(), ".json")
		if _, ok := shops[key]; ok && strings.HasSuffix(e.Name(), ".json") {
			continue
		}
		_ = os.Remove(filepath.Join(dir, "shops", e.Name()))
		m.mu.Lock()
		delete(m.written, "shops/"+e.Name())
		m.mu.Unlock()
	}
	return nil
}

// writeIfChanged 内容（不含 updated_at）和上次写入的不同、或者文件不在了时才写；stamp 在写入前填上 updated_at
func (m *PortalManager) writeIfChanged(name string, v any, stamp func()) error {
	body, err := json.Marshal(v)
	if err != nil {
		return err
	}
	sum := sha256.Sum256(body)
	hash := hex.EncodeToString(sum[:])
	path := filepath.Join(m.dir(), filepath.FromSlash(name))
	m.mu.Lock()
	same := m.written[name] == hash
	m.mu.Unlock()
	if same {
		if _, err := os.Stat(path); err == nil {
			return nil
		}
	}
	stamp()
	if body, err = json.MarshalIndent(v, "", "  "); err != nil {
		return err
	}
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, body, 0o644); err != nil {
		return err
	}
	if err := os.Rename(tmp, path); err != nil {
		_ = os.Remove(tmp)
		return err
	}
	m.mu.Lock()
	m.written[name] = hash
	m.mu.Unlock()
	return nil
}

// serverInfo 一台服务器在快照里的信息
func (m *PortalManager) serverInfo(s models.Server, room models.LinkRoom, groups map[string]int, nextRestart map[uint]time.Time) portalSnapshotServerInfo {
	sp := buildServerProfile(s)
	game := utils.NormalizeGame(s.Game)
	info := portalSnapshotServerInfo{
		Key:        strconv.FormatUint(uint64(s.ID), 10),
		Game:       game,
		Name:       s.SessionName,
		Map:        s.Map,
		Status:     "offline",
		MaxPlayers: s.MaxPlayers,
		Password:   sp.Password,
		PvP:        sp.PvP,
		ArkApi:     s.ArkApi,
		Mods:       sp.Mods,
		Rates:      portalRates{XP: sp.RateXP, Taming: sp.RateTame, Harvest: sp.RateHarvest, Breeding: sp.RateBreeding},
	}
	if s.ClusterID != "" {
		if _, ok := groups[s.ClusterID]; !ok {
			groups[s.ClusterID] = len(groups) + 1
		}
		info.ClusterGroup = groups[s.ClusterID]
	}

	switch s.Status {
	case "running":
		rt := linkRuntimeOf(s)
		if rt.Ready {
			info.Status = "online"
			info.Players = rt.Players
		} else {
			info.Status = "starting"
		}
		if rt.Version != "" {
			m.mu.Lock()
			m.versions[s.ID] = rt.Version
			m.mu.Unlock()
		}
	case "starting":
		info.Status = "starting"
	}
	m.mu.Lock()
	info.Version = m.versions[s.ID]
	m.mu.Unlock()

	if s.Published {
		if host := PublicHost(); host != "" {
			info.Address = hostPort(host, s.Port)
			if game == utils.GameASE {
				info.QueryPort = s.QueryPort
			}
		}
	}
	if s.Link && linkRegistered(room) && linkManager.inRoom(s.ID) {
		info.RoomID = room.RoomID
	}
	if t, ok := nextRestart[s.ID]; ok {
		v := t.UTC().Format(time.RFC3339)
		info.NextRestartAt = &v
	}
	return info
}

// portalNextRestarts 每台服务器下一次定时重启、更新或停服的时间
func portalNextRestarts() map[uint]time.Time {
	var schedules []models.Schedule
	if err := database.DB.Where("enabled = ? AND action IN ?", true,
		[]string{models.ScheduleRestart, models.ScheduleUpdate, models.ScheduleStop}).Find(&schedules).Error; err != nil {
		return nil
	}
	out := map[uint]time.Time{}
	for _, sch := range schedules {
		if next := withNextRun(sch).NextRunAt; next != nil {
			if cur, ok := out[sch.ServerID]; !ok || next.Before(cur) {
				out[sch.ServerID] = *next
			}
		}
	}
	return out
}

// ---- 设置和公告（接口） ----

// GetPortal 玩家站设置、状态和服务器
func (s *ServerService) GetPortal() (*models.PortalInfo, error) {
	settings, err := loadPortalSettings()
	if err != nil {
		return nil, fmt.Errorf("读取玩家站设置失败: %w", err)
	}
	var servers []models.Server
	if err := database.DB.Order("id").Find(&servers).Error; err != nil {
		return nil, fmt.Errorf("读取服务器列表失败: %w", err)
	}
	hidden := portalHidden(settings)
	info := &models.PortalInfo{
		Settings: settings,
		Links:    portalLinks(settings),
		Status:   portalManager.Status(),
		Servers:  []models.PortalServer{},
		Image:    utils.PortalImage,
	}
	if info.Status.State == models.PortalDisabled {
		info.Status.URL = portalURL(settings)
	}
	for _, srv := range servers {
		shop, shopErr := portalShopStatus(srv)
		info.Servers = append(info.Servers, models.PortalServer{
			ID: srv.ID, Name: srv.SessionName, Game: utils.NormalizeGame(srv.Game),
			Hidden: hidden[srv.ID], Shop: shop, ShopError: shopErr,
		})
	}
	return info, nil
}

// UpdatePortal 修改玩家站设置（整体替换）
func (s *ServerService) UpdatePortal(req models.PortalSettingsRequest) (*models.PortalInfo, error) {
	settings, err := loadPortalSettings()
	if err != nil {
		return nil, fmt.Errorf("读取玩家站设置失败: %w", err)
	}
	req, err = normalizePortalSettings(req)
	if err != nil {
		return nil, err
	}
	if req.Enabled {
		if err := checkPortalPort(req.Port); err != nil {
			return nil, err
		}
	}
	links, _ := json.Marshal(req.Links)
	hidden := make([]string, 0, len(req.HiddenServers))
	for _, id := range req.HiddenServers {
		hidden = append(hidden, strconv.FormatUint(uint64(id), 10))
	}
	updates := map[string]any{
		"enabled":        req.Enabled,
		"port":           req.Port,
		"public_url":     req.PublicURL,
		"name":           req.Name,
		"description":    req.Description,
		"qq_group":       req.QQGroup,
		"discord":        req.Discord,
		"links_json":     string(links),
		"hidden_servers": strings.Join(hidden, ","),
	}
	if err := database.DB.Model(&settings).Updates(updates).Error; err != nil {
		return nil, fmt.Errorf("保存玩家站设置失败: %w", err)
	}
	portalManager.wake()
	return s.GetPortal()
}

func invalidPortal(msg string) error {
	return fmt.Errorf("%w：%s", ErrInvalidPortalSettings, msg)
}

func normalizePortalSettings(req models.PortalSettingsRequest) (models.PortalSettingsRequest, error) {
	req.Name = strings.TrimSpace(req.Name)
	req.Description = strings.TrimSpace(req.Description)
	req.QQGroup = strings.TrimSpace(req.QQGroup)
	req.Discord = strings.TrimSpace(req.Discord)
	req.PublicURL = strings.TrimRight(strings.TrimSpace(req.PublicURL), "/")
	if req.Port == 0 {
		req.Port = defaultPortalPort
	}
	switch {
	case req.Port < 1024 || req.Port > 65535:
		return req, invalidPortal("端口要在 1024–65535 之间")
	case utf8.RuneCountInString(req.Name) > maxPortalNameRunes:
		return req, invalidPortal("站点名称最多 64 字")
	case utf8.RuneCountInString(req.Description) > maxPortalDescRunes:
		return req, invalidPortal("简介最多 2000 字")
	case req.QQGroup != "" && !qqGroupRe.MatchString(req.QQGroup):
		return req, invalidPortal("QQ 群号应为 5–12 位数字")
	case req.Discord != "" && !validHTTPURL(req.Discord):
		return req, invalidPortal("Discord 链接要以 http:// 或 https:// 开头")
	case req.PublicURL != "" && !validHTTPURL(req.PublicURL):
		return req, invalidPortal("访问地址要以 http:// 或 https:// 开头")
	case len(req.Links) > maxPortalLinks:
		return req, invalidPortal("链接最多 10 个")
	}
	links := make([]models.PortalLink, 0, len(req.Links))
	for _, l := range req.Links {
		l.Title, l.URL = strings.TrimSpace(l.Title), strings.TrimSpace(l.URL)
		if l.Title == "" && l.URL == "" {
			continue
		}
		if !validHTTPURL(l.URL) {
			return req, invalidPortal("链接要以 http:// 或 https:// 开头")
		}
		if utf8.RuneCountInString(l.Title) > maxPortalNameRunes {
			return req, invalidPortal("链接标题最多 64 字")
		}
		links = append(links, l)
	}
	req.Links = links
	return req, nil
}

func validHTTPURL(v string) bool {
	if len(v) > maxPortalURLLength || strings.ContainsAny(v, " \t\r\n") {
		return false
	}
	u, err := url.Parse(v)
	return err == nil && (u.Scheme == "http" || u.Scheme == "https") && u.Host != ""
}

// checkPortalPort 玩家站端口不能和服务器的端口、联机房节点的 RPC 端口冲突
func checkPortalPort(port int) error {
	var servers []models.Server
	if err := database.DB.Select("id", "session_name", "port", "query_port", "rcon_port").Find(&servers).Error; err != nil {
		return fmt.Errorf("检查端口冲突失败: %w", err)
	}
	for _, s := range servers {
		for _, u := range serverPorts(s.Port, s.QueryPort, s.RCONPort) {
			if u.port == port {
				return fmt.Errorf("%w: 玩家站端口 %d 已被服务器「%s」的%s占用", ErrPortConflict, port, s.SessionName, u.name)
			}
		}
	}
	var room models.LinkRoom
	if err := database.DB.Limit(1).Find(&room, 1).Error; err == nil && room.RPCPort > 0 && room.RPCPort == port {
		return fmt.Errorf("%w: 玩家站端口 %d 已被联机房节点占用", ErrPortConflict, port)
	}
	return nil
}

// portalPortInUse 新建、修改服务器时检查端口是否和开启着的玩家站冲突
func portalPortInUse(ports []portUse) error {
	var settings models.PortalSettings
	if err := database.DB.Limit(1).Find(&settings, 1).Error; err != nil {
		return fmt.Errorf("检查端口冲突失败: %w", err)
	}
	if !settings.Enabled {
		return nil
	}
	for _, p := range ports {
		if p.port == settings.Port {
			return fmt.Errorf("%w: %s %d 已被玩家站占用", ErrPortConflict, p.name, p.port)
		}
	}
	return nil
}

// ListAnnouncements 公告，从新到旧
func (s *ServerService) ListAnnouncements() ([]models.PortalAnnouncement, error) {
	var list []models.PortalAnnouncement
	if err := database.DB.Order("published_at desc, id desc").Find(&list).Error; err != nil {
		return nil, fmt.Errorf("读取公告失败: %w", err)
	}
	return list, nil
}

func normalizeAnnouncement(req models.PortalAnnouncementRequest) (models.PortalAnnouncementRequest, error) {
	req.Title = strings.TrimSpace(req.Title)
	req.Body = strings.TrimSpace(strings.ReplaceAll(req.Body, "\r\n", "\n"))
	if req.Title == "" || utf8.RuneCountInString(req.Title) > maxAnnouncementTitle || utf8.RuneCountInString(req.Body) > maxAnnouncementBody ||
		strings.ContainsAny(req.Title, "\r\n") {
		return req, ErrInvalidAnnouncement
	}
	return req, nil
}

// CreateAnnouncement 发布公告
func (s *ServerService) CreateAnnouncement(req models.PortalAnnouncementRequest) (*models.PortalAnnouncement, error) {
	req, err := normalizeAnnouncement(req)
	if err != nil {
		return nil, err
	}
	var n int64
	if err := database.DB.Model(&models.PortalAnnouncement{}).Count(&n).Error; err != nil {
		return nil, fmt.Errorf("保存公告失败: %w", err)
	}
	if n >= maxPortalAnnouncement {
		return nil, ErrTooManyAnnouncements
	}
	a := models.PortalAnnouncement{Title: req.Title, Body: req.Body, Pinned: req.Pinned, PublishedAt: time.Now()}
	if err := database.DB.Create(&a).Error; err != nil {
		return nil, fmt.Errorf("保存公告失败: %w", err)
	}
	portalManager.wake()
	return &a, nil
}

// UpdateAnnouncement 修改公告
func (s *ServerService) UpdateAnnouncement(id string, req models.PortalAnnouncementRequest) (*models.PortalAnnouncement, error) {
	req, err := normalizeAnnouncement(req)
	if err != nil {
		return nil, err
	}
	var a models.PortalAnnouncement
	if err := database.DB.First(&a, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrAnnouncementNotFound
		}
		return nil, fmt.Errorf("读取公告失败: %w", err)
	}
	if err := database.DB.Model(&a).Updates(map[string]any{"title": req.Title, "body": req.Body, "pinned": req.Pinned}).Error; err != nil {
		return nil, fmt.Errorf("保存公告失败: %w", err)
	}
	portalManager.wake()
	return &a, nil
}

// DeleteAnnouncement 删除公告
func (s *ServerService) DeleteAnnouncement(id string) error {
	res := database.DB.Delete(&models.PortalAnnouncement{}, "id = ?", id)
	if res.Error != nil {
		return fmt.Errorf("删除公告失败: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return ErrAnnouncementNotFound
	}
	portalManager.wake()
	return nil
}

// ---- Docker ----

type dockerPortalOps struct{}

func (dockerPortalOps) ImageReady(image string) (bool, error) {
	return dockerLinkOps{}.ImageReady(image)
}

func (dockerPortalOps) ImageID(image string) (string, error) {
	return dockerLinkOps{}.ImageID(image)
}

func (dockerPortalOps) Inspect(name string) (docker_manager.LinkContainerInfo, error) {
	return dockerLinkOps{}.Inspect(name)
}

func (dockerPortalOps) Create(spec docker_manager.PortalContainerSpec) error {
	dm, err := docker_manager.GetDockerManager()
	if err != nil {
		return err
	}
	return dm.CreatePortalContainer(spec)
}

func (dockerPortalOps) Remove() error {
	dm, err := docker_manager.GetDockerManager()
	if err != nil {
		return err
	}
	return dm.RemovePortalContainer()
}

// portalLinkTCPPort 联机房节点放行的 TCP 端口：开启了玩家站时是玩家站的端口（加入房间的玩家经联机助手访问玩家站），
// 否则为 0（拒绝所有 TCP 入站）
func portalLinkTCPPort() string {
	var settings models.PortalSettings
	if err := database.DB.Limit(1).Find(&settings, 1).Error; err != nil || !settings.Enabled {
		return "0"
	}
	return strconv.Itoa(settings.Port)
}
