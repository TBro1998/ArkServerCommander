package server

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	"ark-server-commander/config"
	"ark-server-commander/database"
	"ark-server-commander/models"
	"ark-server-commander/service/docker_manager"
	"ark-server-commander/utils"

	"go.uber.org/zap"
)

// P2P 联机房（arkfgf.com 联机助手）。
// 需要先绑定 arkfgf.com 账号（account.go），房间的请求都用账号令牌；
// 未绑定时开启了 link 的服务器状态为 login_required，不登记、不心跳、不运行节点。
//
// 每个管理端一个房间（models.LinkRoom，只有一行）：所有开启了 link 的服务器是同一个房间里的多台服务器，
// 共用一个 EasyTier 节点容器 ase-link-room（宿主机网络，--p2p-only，只放行这些服务器的游戏端口、+1 和查询端口的 UDP）
// 和一个凭证池。玩家用联机助手从网站拿一次性凭证加入这个 EasyTier 网络，直连（打洞）到节点，再经宿主机
// 127.0.0.1 上 Docker 的端口映射到达各台 ARK 容器，按端口区分服务器。
//
// LinkManager 每 15 秒检查一次：
//   - 有服务器开启 link：没有登记（room_id 为 0，或还留着旧版房间令牌）或资料变化（房间资料、各服务器的名称、地图、Mod、
//     倍率、端口等）→ POST /link/rooms（账号令牌；还留着旧版房间令牌时带 legacy_token 让网站删除旧房间，成功后清空）
//   - 任意一台开启 link 的服务器运行中且可以连接过 → 确保节点容器在运行（参数或镜像变化时重建）；全部不在运行时删除节点。
//     节点放行所有开启 link 的服务器的端口（包括已停止的），单台服务器启停不会重建节点、不影响其他服务器上的玩家
//   - 心跳：在线时每 30 秒、离线时每 60 秒，带每台服务器的在线和人数；在线状态变化、有凭证要上传或吊销要确认时提前发
//   - 按心跳返回的 want_credentials 在节点里签发凭证（每轮最多 20 个），处理 revoke 后在下次心跳回报 revoked_ack
//   - 心跳等返回 404 room_not_found → 清掉房间 ID 和网络名、删除节点和凭证，下一轮重新登记
//   - 没有服务器开启 link：已登记 → 用账号令牌 DELETE 下架；还留着旧版房间令牌 → 用旧令牌 DELETE（不需要绑定）。
//     之后删除容器、凭证目录，清空房间 ID、网络名、网络密钥和旧令牌
//   - 任何请求返回 401 invalid_account → 视为已解绑，清空本地账号和登记（account.go 的 clearLocked）
//   - 旧版每台服务器一个房间（servers.link_token 不为空）→ 用旧令牌下架后清空；旧的 ase-link-<id> 容器一律删除
//   - 请求失败按 1、2、4…30 分钟退避，429 / 503 按 Retry-After

const (
	linkTickInterval        = 15 * time.Second
	linkHeartbeatOnline     = 30 * time.Second
	linkHeartbeatOffline    = 60 * time.Second
	linkBackoffBase         = time.Minute
	linkBackoffMax          = 30 * time.Minute
	linkRequestTimeout      = 15 * time.Second
	linkExecTimeout         = 15 * time.Second
	linkCredentialTTL       = 30 * 24 * time.Hour
	linkCredentialsPerTick  = 20 // 每轮最多签发的凭证数
	linkCredentialsPerBeat  = 50 // 每次心跳最多上传的凭证数（网站的上限）
	linkVirtualIP           = "10.126.126.1"
	linkCredentialFile      = "/data/cred.json"
	linkCommandHashLabel    = "com.arkservercommander.link.hash"
	linkRPCPortBase         = 20000 // 节点 RPC 端口，和已知端口冲突时顺延
	linkNetworkSecretLength = 32
	linkRoomID              = 1 // models.LinkRoom 只有这一行
	linkRoomNameMax         = 64
)

// linkCredential 已签发、等待随心跳上传的凭证
type linkCredential struct {
	ID        string `json:"id"`
	Secret    string `json:"secret"`
	ExpiresAt int64  `json:"expires_at"`
}

// linkServerProfile 房间里一台服务器的资料（POST /link/rooms 的 servers[]）
type linkServerProfile struct {
	Key          string   `json:"key"`  // 服务器 ID
	Game         string   `json:"game"` // ase 或 asa；ASA 没有查询端口（query_port 为 0），玩家用控制台命令 open 进服
	Name         string   `json:"name"`
	Map          string   `json:"map"`
	MaxPlayers   int      `json:"max_players"`
	Password     bool     `json:"password"`
	PvP          bool     `json:"pvp"`
	ArkApi       bool     `json:"arkapi"`
	Cluster      string   `json:"cluster"`
	Mods         []string `json:"mods"`
	Version      string   `json:"version"`
	RateXP       float64  `json:"rate_xp"`
	RateTame     float64  `json:"rate_tame"`
	RateHarvest  float64  `json:"rate_harvest"`
	RateBreeding float64  `json:"rate_breeding"`
	GamePort     int      `json:"game_port"`
	QueryPort    int      `json:"query_port"`
}

// linkRoomProfile 登记房间时提交的资料（POST /link/rooms）
type linkRoomProfile struct {
	Name        string              `json:"name"`
	Description string              `json:"description"`
	Tags        []string            `json:"tags"`
	QQGroup     string              `json:"qq_group"`
	Discord     string              `json:"discord"`
	ServerIP    string              `json:"server_ip"`
	// 玩家站的 TCP 端口（开启了玩家站时），加入的玩家经联机助手把本机端口转发过来；0 表示没有
	PortalPort int                 `json:"portal_port"`
	Servers    []linkServerProfile `json:"servers"`
}

// linkServerBeat 心跳里一台服务器的状态
type linkServerBeat struct {
	id      uint
	Key     string `json:"key"`
	Online  bool   `json:"online"`
	Players int    `json:"players"`
}

// linkRuntime 服务器当前能否连接、在线人数和游戏版本
type linkRuntime struct {
	Ready   bool
	Players int
	Version string
}

// linkDocker 联机房用到的 Docker 操作，测试时替换
type linkDocker interface {
	// ImageReady 镜像已存在时返回 true；不存在时在后台开始下载并返回 false
	ImageReady(image string) (bool, error)
	ImageID(image string) (string, error)
	Inspect(name string) (docker_manager.LinkContainerInfo, error)
	Create(spec docker_manager.LinkContainerSpec) error
	Remove(name string) error
	List() ([]string, error)
	Exec(name string, argv []string, timeout time.Duration) (docker_manager.ExecResult, error)
	// RemoveDir 删除房间的凭证目录
	RemoveDir() error
	// RemoveLegacyDir 删除旧版每台服务器一个房间时的凭证目录
	RemoveLegacyDir(serverID uint) error
}

// linkBackoff 网站请求失败的退避
type linkBackoff struct {
	nextAt   time.Time
	failures int
	errCode  string
	errMsg   string
}

// linkServerState 房间里一台服务器的状态，只在内存里
type linkServerState struct {
	wasReady      bool      // 这次运行期间已经可以连接过
	notReadySince time.Time // 运行中但查询不到的开始时间
	version       string    // 最近一次查到的游戏版本，停服后保留，避免资料来回变化
	reported      bool      // 最近一次成功的心跳里报告为可以连接
}

// linkState 房间的状态，只在内存里；管理端重启后第一轮重新登记（拿回 peers 和房间地址）
type linkState struct {
	linkBackoff

	profileHash string          // 最近一次登记成功时的资料摘要，变化时重新登记
	keys        map[string]bool // 最近一次登记成功时房间里的服务器（key），发布到收录站前要等服务器移出房间（inRoom）
	url         string          // 网站上的房间页面
	peers       []string        // 公共节点地址，由网站返回

	// 心跳
	lastBeat    time.Time // 最近一次心跳成功的时间
	lastOnline  bool      // 最近一次成功心跳报告的 online
	lastServers string    // 最近一次成功心跳里各服务器的在线情况，变化时提前心跳
	beatNow     bool      // 踢人后尽快发心跳，拿到要吊销的凭证
	want        int       // 网站要补充的凭证数
	revoke      []string  // 网站要求吊销、还没确认的凭证
	acks        []string  // 已在节点吊销、等待回报的凭证
	pending     []linkCredential

	// 节点
	nodeErrCode string
	nodeErr     string
}

// LinkManager 管理联机房：登记房间、心跳、节点容器、凭证
type LinkManager struct {
	mu      sync.Mutex
	room    *linkState
	servers map[uint]*linkServerState
	legacy  map[uint]*linkBackoff // 旧版每台服务器一个房间，下架失败时的退避
	client  *http.Client
	api     func() string
	now     func() time.Time
	docker  linkDocker
	runtime func(models.Server) linkRuntime
	profile func(models.Server) linkServerProfile
	newID   func() string
	wakeCh  chan struct{} // 绑定账号后马上处理一轮

	// 已经检查过管理端启动前留下的节点容器（只在第一次处理房间时看一次，见 process）
	adoptChecked bool
}

var linkManager = &LinkManager{
	room:    &linkState{},
	servers: map[uint]*linkServerState{},
	legacy:  map[uint]*linkBackoff{},
	client:  &http.Client{Timeout: linkRequestTimeout},
	api:     func() string { return config.LinkAPI },
	now:     time.Now,
	docker:  dockerLinkOps{},
	runtime: linkRuntimeOf,
	profile: buildLinkServerProfile,
	newID:   newUUID,
	wakeCh:  make(chan struct{}, 1),
}

// StartLinkManager 启动联机房后台协程。ARKFGF_API 为空时只清理残留的节点容器
func StartLinkManager() {
	if linkManager.api() == "" {
		utils.Info("P2P 联机房已关闭（ARKFGF_API 为空）")
	}
	go func() {
		linkManager.tick()
		ticker := time.NewTicker(linkTickInterval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
			case <-linkManager.wakeCh:
			}
			linkManager.tick()
		}
	}()
}

// wake 让后台协程马上处理一轮
func (m *LinkManager) wake() {
	select {
	case m.wakeCh <- struct{}{}:
	default:
	}
}

// reset 清空房间的内存状态和退避（绑定或解除绑定账号后），服务器的在线报告也一并清除
func (m *LinkManager) reset() {
	m.mu.Lock()
	defer m.mu.Unlock()
	*m.room = linkState{}
	for _, ss := range m.servers {
		ss.reported = false
	}
}

// linkServerColumns 联机房用到的服务器字段
var linkServerColumns = []string{"id", "game", "session_name", "map", "max_players", "port", "query_port", "rcon_port", "cluster_id",
	"game_mod_ids", "arkapi", "server_args_json", "status", "link", "listing_description", "listing_tags", "listing_qq", "listing_discord"}

// loadLinkRoom 读取联机房记录，不存在时创建
func loadLinkRoom() (models.LinkRoom, error) {
	var room models.LinkRoom
	err := database.DB.Where(models.LinkRoom{ID: linkRoomID}).FirstOrCreate(&room).Error
	return room, err
}

func (m *LinkManager) tick() {
	arkfgfMu.Lock()
	defer arkfgfMu.Unlock()
	wantNode := false
	linked := map[uint]bool{}

	if m.api() != "" {
		m.cleanupLegacy()
		token := accountManager.token()

		// 发布到收录站和联机房只能二选一：从发布切换过来的服务器要等收录站取消发布（published 变成 false）后才放进房间，
		// 否则网站拒绝登记（409 published）
		var servers []models.Server
		if err := database.DB.Select(linkServerColumns).Where("link = ? AND published = ?", true, false).Order("id").Find(&servers).Error; err != nil {
			utils.Warn("联机房：读取服务器列表失败", zap.Error(err))
			return
		}
		room, err := loadLinkRoom()
		if err != nil {
			utils.Warn("联机房：读取房间失败", zap.Error(err))
			return
		}
		for _, server := range servers {
			linked[server.ID] = true
		}
		switch {
		case len(servers) > 0:
			// 未绑定账号：状态为 login_required（见 statusFor），不登记、不运行节点
			if token != "" {
				wantNode = m.process(&room, servers, token)
			}
		case room.Token != "":
			// 没有服务器开启联机房了，还留着旧版房间令牌（升级后还没登记过）：用旧令牌下架，不需要绑定
			m.unregister(room, room.Token)
		case room.RoomID != 0 && token != "":
			// 没有服务器开启联机房了：用账号令牌下架
			m.unregister(room, token)
		}
	}

	// 不需要的节点容器（没有服务器在运行、关闭了联机房、旧版的 ase-link-<id>）一律删除
	if names, err := m.docker.List(); err != nil {
		utils.Warn("联机房：列出节点容器失败", zap.Error(err))
	} else {
		for _, name := range names {
			if wantNode && name == utils.LinkContainerName {
				continue
			}
			if err := m.docker.Remove(name); err != nil {
				utils.Warn("联机房：删除节点容器失败", zap.String("container", name), zap.Error(err))
			}
		}
	}

	m.mu.Lock()
	for id := range m.servers {
		if !linked[id] {
			delete(m.servers, id)
		}
	}
	m.mu.Unlock()
}

func (m *LinkManager) serverState(id uint) *linkServerState {
	m.mu.Lock()
	defer m.mu.Unlock()
	st, ok := m.servers[id]
	if !ok {
		st = &linkServerState{}
		m.servers[id] = st
	}
	return st
}

func linkServerKey(id uint) string { return strconv.FormatUint(uint64(id), 10) }

// linkRegistered 房间已用账号登记：有房间 ID，并且旧版房间令牌已经清空（还留着旧令牌时要先带 legacy_token 登记一次，网站删除旧房间）
func linkRegistered(room models.LinkRoom) bool {
	return room.RoomID != 0 && room.Token == ""
}

// process 处理房间（servers 是开启了联机房的服务器，按 ID 排序），返回节点容器是否应该保留。
// 账号失效（本地已清空）时返回 false
func (m *LinkManager) process(room *models.LinkRoom, servers []models.Server, token string) bool {
	if err := ensureLinkLocal(room); err != nil {
		utils.Warn("联机房：生成本地参数失败", zap.Error(err))
		return false
	}
	st := m.room
	wantNode := false
	// 管理端重启后第一次处理：节点容器还在运行，说明重启前房间里有服务器可以连接。"可以连接过"只记在内存里，
	// 运行监控还没查到结果时会被当成不需要节点而删掉，房间里的玩家全部断开。所以这一次让运行中的服务器接上
	// 重启前的状态：节点保留（参数不变就不重建）；之后和平常一样，查询不到超过 notReadyGrace 时报告这台离线
	adopt := false
	if !m.adoptChecked {
		m.adoptChecked = true
		if info, err := m.docker.Inspect(utils.LinkContainerName); err == nil && info.Exists && info.Running {
			adopt = true
		}
	}
	beats := make([]linkServerBeat, len(servers))
	profiles := make([]linkServerProfile, len(servers))
	for i, server := range servers {
		ss := m.serverState(server.ID)
		running := server.Status == "running"
		rt := linkRuntime{}
		if running {
			rt = m.runtime(server)
		}

		m.mu.Lock()
		now := m.now()
		if adopt && running {
			ss.wasReady = true
		}
		switch {
		case !running:
			ss.wasReady, ss.notReadySince = false, time.Time{}
		case rt.Ready:
			ss.wasReady, ss.notReadySince = true, time.Time{}
		case ss.wasReady && ss.notReadySince.IsZero():
			ss.notReadySince = now
		}
		if rt.Version != "" {
			ss.version = rt.Version
		}
		// 可以连接过一次后，整个运行期间保留节点（A2S 偶尔超时不删容器，玩家不掉线）；
		// 但查询不到超过 notReadyGrace（崩溃重启、卡死）时向网站报告这台服务器离线
		wantNode = wantNode || (running && ss.wasReady)
		online := running && (rt.Ready || (ss.wasReady && now.Sub(ss.notReadySince) < notReadyGrace))
		version := ss.version
		m.mu.Unlock()

		beats[i] = linkServerBeat{id: server.ID, Key: linkServerKey(server.ID), Online: online}
		if online {
			beats[i].Players = rt.Players
		}
		p := m.profile(server)
		p.Key = linkServerKey(server.ID)
		p.Version = version
		profiles[i] = p
	}
	if !wantNode {
		m.mu.Lock()
		st.nodeErrCode, st.nodeErr = "", ""
		m.mu.Unlock()
	}

	// 登记或更新房间
	profile := buildLinkRoomProfile(*room, servers, profiles)
	hash := hashJSON(profile)
	m.mu.Lock()
	needRegister := (!linkRegistered(*room) || st.profileHash != hash) && !m.now().Before(st.nextAt)
	m.mu.Unlock()
	if needRegister && !m.register(room, profile, hash, token) {
		return false
	}

	// 节点容器
	nodeUp := false
	if wantNode && linkRegistered(*room) && room.NetworkName != "" {
		m.mu.Lock()
		peers := slices.Clone(st.peers)
		m.mu.Unlock()
		if len(peers) > 0 {
			nodeUp = m.ensureNode(*room, servers, peers)
		}
	}

	if linkRegistered(*room) {
		if !nodeUp {
			for i := range beats {
				beats[i].Online, beats[i].Players = false, 0
			}
		}
		if !m.heartbeat(room, beats, nodeUp, token) {
			return false
		}
	}
	return wantNode
}

// linkRegisterRequest 登记房间的请求：房间资料，加上还没清理的旧版房间令牌
type linkRegisterRequest struct {
	linkRoomProfile
	LegacyToken string `json:"legacy_token,omitempty"`
}

// register 用账号令牌登记或更新房间（POST /link/rooms），保存房间 ID、网络名和 peers；
// 带了旧版房间令牌（网站删除旧房间）的，成功后清空。账号失效时返回 false
func (m *LinkManager) register(room *models.LinkRoom, profile linkRoomProfile, hash, token string) bool {
	st := m.room
	var resp struct {
		ID          uint     `json:"id"`
		NetworkName string   `json:"network_name"`
		Peers       []string `json:"peers"`
		URL         string   `json:"url"`
	}
	body := linkRegisterRequest{linkRoomProfile: profile, LegacyToken: room.Token}
	status, retryAfter, code, err := m.request(http.MethodPost, "/link/rooms", token, body, &resp)
	if err != nil {
		if isInvalidAccount(status, code) {
			accountManager.invalidateLocked()
			return false
		}
		m.fail(&st.linkBackoff, "room", status, retryAfter, code, err)
		return true
	}
	if resp.ID == 0 {
		m.fail(&st.linkBackoff, "room", status, 0, "invalid_response", errors.New("网站没有返回房间 ID"))
		return true
	}

	oldID := room.RoomID
	fresh := !linkRegistered(*room) // 第一次用账号登记
	updates := map[string]any{}
	if room.Token != "" {
		updates["token"] = ""
	}
	if resp.ID != room.RoomID {
		updates["room_id"] = resp.ID
	}
	if resp.NetworkName != "" && resp.NetworkName != room.NetworkName {
		updates["network_name"] = resp.NetworkName
	}
	if len(updates) > 0 {
		if err := database.DB.Model(&models.LinkRoom{}).Where("id = ?", room.ID).Updates(updates).Error; err != nil {
			// 没保存下来，下一轮再登记一次（网站按 Commander 找到同一个房间）
			m.fail(&st.linkBackoff, "room", 0, 0, "local_error", fmt.Errorf("保存房间登记失败: %w", err))
			return true
		}
		room.Token, room.RoomID = "", resp.ID
		if resp.NetworkName != "" {
			room.NetworkName = resp.NetworkName
		}
	}
	newRoom := oldID != 0 && resp.ID != oldID
	if newRoom {
		// 网站给的是另一个房间：旧房间的凭证（含被封禁玩家的）不能留在节点里
		m.resetNode()
	}

	m.mu.Lock()
	if newRoom || fresh {
		st.lastBeat, st.lastOnline, st.lastServers, st.want = time.Time{}, false, "", 0
		st.revoke, st.acks, st.pending = nil, nil, nil
	}
	st.profileHash = hash
	st.keys = map[string]bool{}
	for _, server := range profile.Servers {
		st.keys[server.Key] = true
	}
	st.url = resp.URL
	if len(resp.Peers) > 0 {
		st.peers = resp.Peers
	}
	m.succeed(&st.linkBackoff)
	m.mu.Unlock()
	if fresh || newRoom {
		utils.Info("已登记联机房", zap.Uint("room_id", room.RoomID), zap.Int("servers", len(profile.Servers)),
			zap.Bool("legacy_token", body.LegacyToken != ""))
	}
	return true
}

// resetNode 删除节点容器和凭证文件，下次按新参数重建
func (m *LinkManager) resetNode() {
	if err := m.docker.Remove(utils.LinkContainerName); err != nil {
		utils.Warn("联机房：删除节点容器失败", zap.Error(err))
	}
	if err := m.docker.RemoveDir(); err != nil {
		utils.Warn("联机房：删除凭证目录失败", zap.Error(err))
	}
}

// ensureNode 确保节点容器按当前参数运行，返回节点是否在运行
func (m *LinkManager) ensureNode(room models.LinkRoom, servers []models.Server, peers []string) bool {
	st := m.room
	spec := linkContainerSpec(room, servers, peers)
	name := utils.LinkContainerName

	ready, err := m.docker.ImageReady(spec.Image)
	if err != nil {
		m.nodeFail(st, "image_error", fmt.Errorf("检查镜像 %s 失败: %w", spec.Image, err))
		return false
	}
	if !ready {
		return false // 正在下载镜像，下一轮再创建
	}
	info, err := m.docker.Inspect(name)
	if err != nil {
		m.nodeFail(st, "node_error", err)
		return false
	}
	imageID, _ := m.docker.ImageID(spec.Image)
	upToDate := info.Exists && info.Running && info.Labels[linkCommandHashLabel] == spec.Labels[linkCommandHashLabel] &&
		(imageID == "" || info.ImageID == imageID)
	if !upToDate {
		if err := m.docker.Create(spec); err != nil {
			m.nodeFail(st, "node_error", err)
			return false
		}
		if info.Exists {
			utils.Info("联机房节点参数或镜像有变化，已重建")
		}
	}
	return true
}

func (m *LinkManager) nodeFail(st *linkState, code string, err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if st.nodeErrCode == "" {
		utils.Warn("联机房节点操作失败", zap.String("code", code), zap.Error(err))
	}
	st.nodeErrCode, st.nodeErr = code, err.Error()
}

// beatSignature 各服务器的在线情况，变化时提前心跳
func beatSignature(beats []linkServerBeat) string {
	var b strings.Builder
	for _, s := range beats {
		b.WriteString(s.Key)
		if s.Online {
			b.WriteString(":1,")
		} else {
			b.WriteString(":0,")
		}
	}
	return b.String()
}

// heartbeat 到期时发心跳；节点在运行时先吊销网站要求吊销的凭证、签发网站要补充的凭证，随心跳上传。
// 账号失效时返回 false
func (m *LinkManager) heartbeat(room *models.LinkRoom, beats []linkServerBeat, nodeUp bool, token string) bool {
	st := m.room
	online := false
	for _, b := range beats {
		online = online || b.Online
	}
	sig := beatSignature(beats)

	m.mu.Lock()
	now := m.now()
	interval := linkHeartbeatOffline
	if online {
		interval = linkHeartbeatOnline
	}
	work := nodeUp && (st.want > len(st.pending) || len(st.revoke) > 0)
	due := st.lastBeat.IsZero() || !now.Before(st.lastBeat.Add(interval)) || st.lastOnline != online || st.lastServers != sig ||
		st.beatNow || work || len(st.pending) > 0 || len(st.acks) > 0
	if !due || now.Before(st.nextAt) {
		m.mu.Unlock()
		return true
	}
	revoke := slices.Clone(st.revoke)
	acked := slices.Clone(st.acks)
	toIssue := 0
	if nodeUp {
		toIssue = min(st.want-len(st.pending), linkCredentialsPerTick, linkCredentialsPerBeat-len(st.pending))
	}
	m.mu.Unlock()

	if nodeUp {
		var nodeErr error
		var newAcks []string
		for _, id := range revoke {
			if slices.Contains(acked, id) {
				continue
			}
			if err := m.revokeCredential(room.RPCPort, id); err != nil {
				nodeErr = err
				break
			}
			newAcks = append(newAcks, id)
		}
		var issued []linkCredential
		for i := 0; i < toIssue && nodeErr == nil; i++ {
			cred, err := m.generateCredential(room.RPCPort)
			if err != nil {
				nodeErr = err
				break
			}
			issued = append(issued, cred)
		}
		m.mu.Lock()
		st.acks = append(st.acks, newAcks...)
		st.pending = append(st.pending, issued...)
		if nodeErr != nil {
			if st.nodeErrCode == "" {
				utils.Warn("联机房：签发或吊销凭证失败", zap.Error(nodeErr))
			}
			st.nodeErrCode, st.nodeErr = "node_error", nodeErr.Error()
		} else {
			st.nodeErrCode, st.nodeErr = "", ""
		}
		m.mu.Unlock()
	}

	m.mu.Lock()
	body := struct {
		Online      bool             `json:"online"`
		Servers     []linkServerBeat `json:"servers"`
		Credentials []linkCredential `json:"credentials"`
		RevokedAck  []string         `json:"revoked_ack"`
	}{Online: online, Servers: beats, Credentials: []linkCredential{}, RevokedAck: []string{}}
	body.Credentials = append(body.Credentials, st.pending[:min(len(st.pending), linkCredentialsPerBeat)]...)
	body.RevokedAck = append(body.RevokedAck, st.acks...)
	st.beatNow = false
	m.mu.Unlock()

	var resp struct {
		WantCredentials int      `json:"want_credentials"`
		Revoke          []string `json:"revoke"`
		NetworkName     string   `json:"network_name"`
		Peers           []string `json:"peers"`
	}
	path := fmt.Sprintf("/link/rooms/%d/heartbeat", room.RoomID)
	status, retryAfter, code, err := m.request(http.MethodPost, path, token, body, &resp)
	if err != nil {
		switch {
		case isInvalidAccount(status, code):
			accountManager.invalidateLocked()
			return false
		case code == "room_not_found" || status == http.StatusNotFound:
			// 网站上已经没有这个房间（或不属于这个账号了）：下一轮重新登记
			m.roomGone(room)
		default:
			m.fail(&st.linkBackoff, "room", status, retryAfter, code, err)
		}
		return true
	}
	if resp.NetworkName != "" && resp.NetworkName != room.NetworkName {
		if err := database.DB.Model(&models.LinkRoom{}).Where("id = ?", room.ID).
			Update("network_name", resp.NetworkName).Error; err != nil {
			utils.Warn("联机房：保存网络名失败", zap.Error(err))
		} else {
			room.NetworkName = resp.NetworkName
		}
	}

	m.mu.Lock()
	defer m.mu.Unlock()
	st.pending = st.pending[len(body.Credentials):]
	st.acks = slices.DeleteFunc(st.acks, func(id string) bool { return slices.Contains(body.RevokedAck, id) })
	st.want = max(resp.WantCredentials, 0)
	st.revoke = resp.Revoke
	if len(resp.Peers) > 0 {
		st.peers = resp.Peers
	}
	st.lastBeat = m.now()
	st.lastOnline = online
	st.lastServers = sig
	for _, b := range beats {
		if ss, ok := m.servers[b.id]; ok {
			ss.reported = b.Online
		}
	}
	m.succeed(&st.linkBackoff)
	return true
}

// resetRoomState 清空房间的内存状态（保留退避），服务器的在线报告也一并清除
func (m *LinkManager) resetRoomState() {
	m.mu.Lock()
	defer m.mu.Unlock()
	backoff := m.room.linkBackoff
	*m.room = linkState{linkBackoff: backoff}
	for _, ss := range m.servers {
		ss.reported = false
	}
}

// roomGone 网站上已经没有这个房间（room_not_found）：清掉房间 ID、网络名和节点，下一轮重新登记
func (m *LinkManager) roomGone(room *models.LinkRoom) {
	utils.Warn("联机房：网站上的房间已不存在，将重新登记", zap.Uint("room_id", room.RoomID))
	if err := database.DB.Model(&models.LinkRoom{}).Where("id = ?", room.ID).
		Updates(map[string]any{"room_id": 0, "network_name": ""}).Error; err != nil {
		utils.Warn("联机房：清除房间登记失败", zap.Error(err))
		return
	}
	room.RoomID, room.NetworkName = 0, ""
	m.resetNode()
	m.resetRoomState()
}

// unregister 下架房间（DELETE /link/rooms/:id），token 是账号令牌或旧版房间令牌。成功后删除节点容器和凭证目录，
// 清空房间 ID、网络名、网络密钥和旧令牌（保留房间资料）。网站已经没有这个房间（room_not_found / invalid_token / 404）也算成功
func (m *LinkManager) unregister(room models.LinkRoom, token string) {
	st := m.room
	m.mu.Lock()
	wait := m.now().Before(st.nextAt)
	m.mu.Unlock()
	if wait {
		return
	}
	status, retryAfter, code, err := m.request(http.MethodDelete, fmt.Sprintf("/link/rooms/%d", room.RoomID), token, nil, nil)
	if err != nil && isInvalidAccount(status, code) {
		accountManager.invalidateLocked()
		return
	}
	if err != nil && code != "invalid_token" && code != "room_not_found" && status != http.StatusNotFound {
		m.fail(&st.linkBackoff, "room", status, retryAfter, code, err)
		return
	}
	m.resetNode()
	if err := database.DB.Model(&models.LinkRoom{}).Where("id = ?", room.ID).
		Updates(map[string]any{"token": "", "room_id": 0, "network_name": "", "network_secret": ""}).Error; err != nil {
		utils.Warn("联机房：清除房间令牌失败", zap.Error(err))
		return
	}
	m.mu.Lock()
	*m.room = linkState{}
	m.mu.Unlock()
	utils.Info("已下架联机房", zap.Uint("room_id", room.RoomID))
}

// cleanupLegacy 下架旧版每台服务器一个房间时登记的房间（servers.link_token 不为空，含已删除的服务器），
// 成功（或网站已经没有）后删除旧的凭证目录、清空这些字段。旧的 ase-link-<id> 容器由 tick 删除
func (m *LinkManager) cleanupLegacy() {
	var legacy []models.Server
	if err := database.DB.Unscoped().Select("id", "link_token", "link_room_id").
		Where("link_token <> ''").Find(&legacy).Error; err != nil {
		utils.Warn("联机房：读取旧版房间失败", zap.Error(err))
		return
	}
	for _, server := range legacy {
		m.mu.Lock()
		b, ok := m.legacy[server.ID]
		if !ok {
			b = &linkBackoff{}
			m.legacy[server.ID] = b
		}
		wait := m.now().Before(b.nextAt)
		m.mu.Unlock()
		if wait {
			continue
		}
		status, retryAfter, code, err := m.request(http.MethodDelete, fmt.Sprintf("/link/rooms/%d", server.LinkRoomID), server.LinkToken, nil, nil)
		if err != nil && code != "invalid_token" && status != http.StatusNotFound {
			m.fail(b, "legacy-"+linkServerKey(server.ID), status, retryAfter, code, err)
			continue
		}
		if err := m.docker.RemoveLegacyDir(server.ID); err != nil {
			utils.Warn("联机房：删除旧版凭证目录失败", zap.Uint("server_id", server.ID), zap.Error(err))
		}
		if err := database.DB.Unscoped().Model(&models.Server{}).Where("id = ?", server.ID).
			Updates(map[string]any{"link_token": "", "link_room_id": 0, "link_network_name": "", "link_network_secret": "", "link_rpc_port": 0}).Error; err != nil {
			utils.Warn("联机房：清除旧版房间令牌失败", zap.Uint("server_id", server.ID), zap.Error(err))
			continue
		}
		m.mu.Lock()
		delete(m.legacy, server.ID)
		m.mu.Unlock()
		utils.Info("已下架旧版的单服务器联机房", zap.Uint("server_id", server.ID), zap.Uint("room_id", server.LinkRoomID))
	}
}

// generateCredential 在节点里签发一个不可复用的凭证（TTL 30 天）
func (m *LinkManager) generateCredential(rpcPort int) (linkCredential, error) {
	id := m.newID()
	out, err := m.cli(rpcPort, "credential", "generate", "--ttl", strconv.Itoa(int(linkCredentialTTL/time.Second)),
		"--credential-id", id, "--reusable", "false")
	if err != nil {
		return linkCredential{}, fmt.Errorf("签发凭证失败: %w", err)
	}
	var res struct {
		CredentialID     string `json:"credential_id"`
		CredentialSecret string `json:"credential_secret"`
	}
	if err := json.Unmarshal([]byte(out), &res); err != nil || res.CredentialSecret == "" {
		return linkCredential{}, fmt.Errorf("签发凭证失败：无法解析 easytier-cli 的输出 %q", truncate(out, 200))
	}
	if res.CredentialID != "" {
		id = res.CredentialID
	}
	return linkCredential{ID: id, Secret: res.CredentialSecret, ExpiresAt: m.now().Add(linkCredentialTTL).Unix()}, nil
}

// revokeCredential 在节点里吊销凭证。凭证不存在（{"success": false}、not found）也算完成
func (m *LinkManager) revokeCredential(rpcPort int, id string) error {
	_, err := m.cli(rpcPort, "credential", "revoke", id)
	if err != nil && !strings.Contains(strings.ToLower(err.Error()), "not found") {
		return fmt.Errorf("吊销凭证失败: %w", err)
	}
	return nil
}

// cli 在节点容器里执行 easytier-cli（JSON 输出），返回标准输出
func (m *LinkManager) cli(rpcPort int, args ...string) (string, error) {
	argv := append([]string{"easytier-cli", "-p", fmt.Sprintf("127.0.0.1:%d", rpcPort), "-o", "json"}, args...)
	res, err := m.docker.Exec(utils.LinkContainerName, argv, linkExecTimeout)
	if err != nil {
		return "", err
	}
	if res.ExitCode != 0 {
		msg := strings.TrimSpace(res.Stderr)
		if msg == "" {
			msg = strings.TrimSpace(res.Stdout)
		}
		if msg == "" {
			msg = fmt.Sprintf("退出码 %d", res.ExitCode)
		}
		return "", fmt.Errorf("easytier-cli: %s", truncate(msg, 300))
	}
	return strings.TrimSpace(res.Stdout), nil
}

func (m *LinkManager) request(method, path, token string, body, out any) (int, time.Duration, string, error) {
	return arkfgfRequest(m.client, method, m.api()+path, token, body, out, " arkfgf.com ")
}

// succeed 清除失败记录，调用时持有 m.mu
func (m *LinkManager) succeed(b *linkBackoff) {
	b.failures = 0
	b.errCode, b.errMsg = "", ""
	b.nextAt = time.Time{}
}

// fail 记录网站请求失败并安排重试：指数退避，网站给了 Retry-After 时取两者中较长的
func (m *LinkManager) fail(b *linkBackoff, target string, status int, retryAfter time.Duration, code string, err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	b.failures++
	wait := min(linkBackoffBase<<min(b.failures-1, 5), linkBackoffMax)
	if retryAfter > wait {
		wait = min(retryAfter, 2*time.Hour)
	}
	b.nextAt = m.now().Add(wait)
	if code == "" {
		code = "network"
	}
	b.errCode = code
	b.errMsg = linkErrorMessage(code)
	if b.errMsg == "" {
		b.errMsg = err.Error()
	}
	if b.failures == 1 || b.failures%10 == 0 {
		utils.Warn("联机房请求失败", zap.String("target", target), zap.Int("status", status),
			zap.String("code", code), zap.Int("failures", b.failures), zap.Error(err))
	}
}

func linkErrorMessage(code string) string {
	switch code {
	case "rate_limited", "busy":
		return "arkfgf.com 暂时繁忙或请求过于频繁，稍后自动重试"
	case "too_many_rooms":
		return "这个公网 IP 登记的联机房已达上限"
	case "room_not_found":
		return "网站上的房间已不存在，正在重新登记"
	case "published":
		return "房间里有服务器还发布在收录站上（同一台服务器只能发布到收录站或放进联机房），取消发布后会自动重试"
	case "login_required":
		return "arkfgf.com 要求先绑定账号"
	case "invalid_input":
		return "房间资料不符合网站要求"
	}
	return ""
}

// roomStatusLocked 房间的状态，调用时持有 m.mu
func (m *LinkManager) roomStatusLocked() *models.LinkStatus {
	st := m.room
	status := &models.LinkStatus{State: models.LinkWaiting, RoomURL: st.url}
	switch {
	case st.errCode != "":
		status.State = models.LinkError
		status.ErrorCode, status.Error = st.errCode, st.errMsg
		if !st.nextAt.IsZero() {
			status.RetryAt = st.nextAt.Format(time.RFC3339)
		}
	case st.nodeErrCode != "":
		status.State = models.LinkError
		status.ErrorCode, status.Error = st.nodeErrCode, st.nodeErr
	}
	return status
}

// inRoom 服务器是否还在网站上的房间里（按最近一次登记成功的资料）。从联机房切换到发布的服务器，
// 要等重新登记（或下架）把它移出房间后才发布，否则收录站拒绝发布（409 in_room）
func (m *LinkManager) inRoom(id uint) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.room.keys[linkServerKey(id)]
}

// statusFor 返回服务器的联机房状态，未开启时为 nil。房间出错时所有服务器都显示错误；
// 否则这台服务器在最近一次心跳里报告为可以连接时为 online，没在运行时为 offline，其他情况为 waiting
func (m *LinkManager) statusFor(server models.Server) *models.LinkStatus {
	if !server.Link {
		return nil
	}
	if m.api() == "" {
		return &models.LinkStatus{State: models.LinkDisabled}
	}
	if !accountManager.bound() {
		return linkLoginRequired()
	}
	// 从发布切换过来、收录站还没取消发布：还没放进房间
	if server.Published {
		return &models.LinkStatus{State: models.LinkWaiting}
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	status := m.roomStatusLocked()
	if status.State == models.LinkError || m.room.lastBeat.IsZero() {
		return status
	}
	if ss, ok := m.servers[server.ID]; ok && ss.reported {
		status.State = models.LinkOnline
	} else if server.Status != "running" {
		status.State = models.LinkOffline
	}
	return status
}

// roomStatus 房间的状态（设置页），没有服务器开启联机房时为 nil
func (m *LinkManager) roomStatus(hasServers bool) *models.LinkStatus {
	if !hasServers {
		return nil
	}
	if m.api() == "" {
		return &models.LinkStatus{State: models.LinkDisabled}
	}
	if !accountManager.bound() {
		return linkLoginRequired()
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	status := m.roomStatusLocked()
	if status.State == models.LinkError || m.room.lastBeat.IsZero() {
		return status
	}
	if m.room.lastOnline {
		status.State = models.LinkOnline
	} else {
		status.State = models.LinkOffline
	}
	return status
}

func linkLoginRequired() *models.LinkStatus {
	return &models.LinkStatus{State: models.LinkLoginRequired, ErrorCode: "login_required", Error: loginRequiredMessage}
}

// serverReported 这台服务器在最近一次成功的心跳里是否报告为可以连接
func (m *LinkManager) serverReported(id uint) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	ss, ok := m.servers[id]
	return ok && ss.reported
}

// poke 让下一轮马上发心跳（踢人后尽快拿到要吊销的凭证）
func (m *LinkManager) poke() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.room.beatNow = true
}

// ---- 节点参数 ----

// linkRoomName 房间名称：设置里的名称，为空时用第一台服务器的名称
func linkRoomName(room models.LinkRoom, servers []models.Server) string {
	if name := strings.TrimSpace(room.Name); name != "" {
		return name
	}
	for _, s := range servers {
		if name := strings.TrimSpace(s.SessionName); name != "" {
			return name
		}
	}
	return "ARK Server"
}

// linkNodePorts 节点放行的 UDP 端口：每台服务器的游戏端口、+1 和查询端口（ASA 只有游戏端口），去重，按服务器顺序
func linkNodePorts(servers []models.Server) []string {
	seen := map[int]bool{}
	var ports []string
	for _, s := range servers {
		forward := []int{s.Port, s.Port + 1, s.QueryPort}
		if utils.NormalizeGame(s.Game) == utils.GameASA {
			forward = []int{s.Port}
		}
		for _, p := range forward {
			if p > 0 && !seen[p] {
				seen[p] = true
				ports = append(ports, strconv.Itoa(p))
			}
		}
	}
	return ports
}

// linkNodeArgs easytier-core 的参数
func linkNodeArgs(room models.LinkRoom, servers []models.Server, peers []string) []string {
	args := []string{
		"--network-name", room.NetworkName,
		"--network-secret", room.NetworkSecret,
		"--secure-mode",
		"--credential-file", linkCredentialFile,
		"--peers",
	}
	args = append(args, peers...)
	args = append(args,
		"--ipv4", linkVirtualIP, "--no-tun", "--p2p-only",
		"--tcp-whitelist", portalLinkTCPPort(),
		"--udp-whitelist",
	)
	args = append(args, linkNodePorts(servers)...)
	args = append(args,
		"--listeners", "tcp://0.0.0.0:0", "udp://0.0.0.0:0",
		"--rpc-portal", fmt.Sprintf("127.0.0.1:%d", room.RPCPort),
		"--hostname", linkRoomName(room, servers),
	)
	return args
}

// linkContainerSpec 节点容器配置；标签里的摘要包含镜像和全部参数，参数变化时重建
func linkContainerSpec(room models.LinkRoom, servers []models.Server, peers []string) docker_manager.LinkContainerSpec {
	args := linkNodeArgs(room, servers, peers)
	sum := sha256.Sum256([]byte(utils.EasyTierImage + "\x00" + strings.Join(args, "\x00")))
	return docker_manager.LinkContainerSpec{
		Image:  utils.EasyTierImage,
		Args:   args,
		Labels: map[string]string{linkCommandHashLabel: hex.EncodeToString(sum[:])},
	}
}

// ensureLinkLocal 第一次处理时生成网络密钥、分配 RPC 端口并保存
func ensureLinkLocal(room *models.LinkRoom) error {
	updates := map[string]any{}
	if room.NetworkSecret == "" {
		buf := make([]byte, linkNetworkSecretLength)
		if _, err := rand.Read(buf); err != nil {
			return err
		}
		room.NetworkSecret = hex.EncodeToString(buf)
		updates["network_secret"] = room.NetworkSecret
	}
	if room.RPCPort == 0 {
		port, err := assignLinkRPCPort()
		if err != nil {
			return err
		}
		room.RPCPort = port
		updates["rpc_port"] = port
	}
	if len(updates) == 0 {
		return nil
	}
	return database.DB.Model(&models.LinkRoom{}).Where("id = ?", room.ID).Updates(updates).Error
}

// assignLinkRPCPort 节点 RPC 端口（只监听宿主机 127.0.0.1）：从 20000 开始。
// 管理端运行在容器里，探测不了宿主机端口，只避开数据库里已知的端口（各服务器的游戏端口、+1、查询、RCON，
// 以及旧版节点还没删掉时的 RPC 端口），冲突时顺延
func assignLinkRPCPort() (int, error) {
	var servers []models.Server
	if err := database.DB.Select("id", "port", "query_port", "rcon_port", "link_rpc_port").Find(&servers).Error; err != nil {
		return 0, err
	}
	used := map[int]bool{}
	for _, s := range servers {
		used[s.Port], used[s.Port+1], used[s.QueryPort], used[s.RCONPort] = true, true, true, true
		if s.LinkRPCPort > 0 {
			used[s.LinkRPCPort] = true
		}
	}
	port := linkRPCPortBase
	for used[port] {
		port++
		if port > 65535 {
			return 0, errors.New("找不到空闲的节点 RPC 端口")
		}
	}
	return port, nil
}

// ---- 资料 ----

// linkRoomListing 设置里的房间收录资料；全部为空时用第一台服务器的收录资料
func linkRoomListing(room models.LinkRoom, servers []models.Server) models.ListingProfile {
	p := models.ListingProfile{Description: room.Description, QQGroup: room.QQGroup, Discord: room.Discord, Tags: []string{}}
	for _, tag := range strings.Split(room.Tags, ",") {
		if tag = strings.TrimSpace(tag); tag != "" {
			p.Tags = append(p.Tags, tag)
		}
	}
	if p.Description == "" && p.QQGroup == "" && p.Discord == "" && len(p.Tags) == 0 && len(servers) > 0 {
		return listingOf(servers[0])
	}
	return p
}

// buildLinkRoomProfile 登记房间的资料：房间名称和收录资料，加上每台服务器的资料
func buildLinkRoomProfile(room models.LinkRoom, servers []models.Server, profiles []linkServerProfile) linkRoomProfile {
	listing := linkRoomListing(room, servers)
	return linkRoomProfile{
		Name:        linkRoomName(room, servers),
		Description: listing.Description,
		Tags:        listing.Tags,
		QQGroup:     listing.QQGroup,
		Discord:     listing.Discord,
		ServerIP:    linkVirtualIP,
		PortalPort:  portalPortForLink(),
		Servers:     profiles,
	}
}

// buildLinkServerProfile 房间里一台服务器的资料：和发布到收录站相同的服务器资料，加上端口
func buildLinkServerProfile(server models.Server) linkServerProfile {
	sp := buildServerProfile(server)
	return linkServerProfile{
		Key:          linkServerKey(server.ID),
		Game:         utils.NormalizeGame(server.Game),
		Name:         sp.Name,
		Map:          sp.Map,
		MaxPlayers:   sp.MaxPlayers,
		Password:     sp.Password,
		PvP:          sp.PvP,
		ArkApi:       sp.ArkApi,
		Cluster:      sp.Cluster,
		Mods:         sp.Mods,
		RateXP:       sp.RateXP,
		RateTame:     sp.RateTame,
		RateHarvest:  sp.RateHarvest,
		RateBreeding: sp.RateBreeding,
		GamePort:     server.Port,
		QueryPort:    server.QueryPort,
	}
}

var serverVersionRe = regexp.MustCompile(`\(v([0-9][0-9.]*)\)\s*$`)

// linkRuntimeOf 从运行信息监控取能否连接、人数和游戏版本（A2S 返回的名称末尾是 "(v358.24)"）
func linkRuntimeOf(server models.Server) linkRuntime {
	info := runtimeMonitor.Get(server.ID)
	if info == nil {
		return linkRuntime{}
	}
	rt := linkRuntime{Ready: info.Phase == models.PhaseReady, Players: info.Players}
	if m := serverVersionRe.FindStringSubmatch(info.ServerName); m != nil {
		rt.Version = m[1]
	}
	return rt
}

func hashJSON(v any) string {
	data, _ := json.Marshal(v)
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}

// newUUID 随机 UUID（v4），作为凭证 ID
func newUUID() string {
	var b [16]byte
	_, _ = rand.Read(b[:])
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	h := hex.EncodeToString(b[:])
	return h[0:8] + "-" + h[8:12] + "-" + h[12:16] + "-" + h[16:20] + "-" + h[20:32]
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "…"
}

// ---- 房间资料和成员管理（成员管理代理到网站） ----

var steamIDRe = regexp.MustCompile(`^[0-9]{17}$`)

// linkCall 用账号令牌调用网站的成员管理接口
func (m *LinkManager) linkCall(method, path string, body, out any) error {
	if m.api() == "" {
		return fmt.Errorf("%w：管理端已关闭与 arkfgf.com 的连接（ARKFGF_API 为空）", ErrLinkUnavailable)
	}
	token := accountManager.token()
	if token == "" {
		return ErrAccountRequired
	}
	room, err := loadLinkRoom()
	if err != nil {
		return fmt.Errorf("读取联机房失败: %w", err)
	}
	if !linkRegistered(room) {
		return ErrLinkNotRegistered
	}
	status, _, code, err := m.request(method, fmt.Sprintf("/link/rooms/%d%s", room.RoomID, path), token, body, out)
	if err == nil {
		return nil
	}
	switch {
	case isInvalidAccount(status, code):
		accountManager.invalidate()
		return ErrAccountRequired
	case status == http.StatusUnauthorized:
		return ErrAccountRequired
	case code == "room_not_found":
		// 下一次心跳也会收到 room_not_found，届时重新登记
		return ErrLinkNotRegistered
	case status == http.StatusBadRequest || status == http.StatusNotFound || status == http.StatusConflict || status == http.StatusUnprocessableEntity:
		return fmt.Errorf("%w：%v", ErrLinkRequestFailed, err)
	}
	return fmt.Errorf("%w：%v", ErrLinkUnavailable, err)
}

// GetLinkRoom 联机房的资料、状态和房间里的服务器
func (s *ServerService) GetLinkRoom() (*models.LinkRoomInfo, error) {
	room, err := loadLinkRoom()
	if err != nil {
		return nil, err
	}
	var servers []models.Server
	if err := database.DB.Select("id", "game", "session_name", "port", "query_port", "status").
		Where("link = ?", true).Order("id").Find(&servers).Error; err != nil {
		return nil, err
	}
	info := &models.LinkRoomInfo{
		LinkRoomSettings: models.LinkRoomSettings{
			Name:    room.Name,
			Listing: models.ListingProfile{Description: room.Description, QQGroup: room.QQGroup, Discord: room.Discord, Tags: []string{}},
		},
		Status:  linkManager.roomStatus(len(servers) > 0),
		Servers: []models.LinkRoomMember{},
	}
	for _, tag := range strings.Split(room.Tags, ",") {
		if tag = strings.TrimSpace(tag); tag != "" {
			info.Listing.Tags = append(info.Listing.Tags, tag)
		}
	}
	for _, sv := range servers {
		info.Servers = append(info.Servers, models.LinkRoomMember{
			ID: sv.ID, Game: utils.NormalizeGame(sv.Game), SessionName: sv.SessionName, Port: sv.Port, QueryPort: sv.QueryPort, Status: sv.Status,
			Online: linkManager.serverReported(sv.ID),
		})
	}
	return info, nil
}

// UpdateLinkRoom 修改联机房的名称和收录资料，下一轮重新登记时生效
func (s *ServerService) UpdateLinkRoom(req models.LinkRoomSettings) (*models.LinkRoomInfo, error) {
	name := strings.TrimSpace(req.Name)
	if utf8.RuneCountInString(name) > linkRoomNameMax {
		return nil, fmt.Errorf("%w：房间名称最多 %d 字", ErrInvalidListing, linkRoomNameMax)
	}
	listing, err := normalizeListing(req.Listing)
	if err != nil {
		return nil, err
	}
	if _, err := loadLinkRoom(); err != nil {
		return nil, err
	}
	if err := database.DB.Model(&models.LinkRoom{}).Where("id = ?", linkRoomID).Updates(map[string]any{
		"name": name, "description": listing.Description, "tags": strings.Join(listing.Tags, ","),
		"qq_group": listing.QQGroup, "discord": listing.Discord,
	}).Error; err != nil {
		return nil, err
	}
	return s.GetLinkRoom()
}

// GetLinkMembers 联机房的成员（持有未吊销凭证的玩家）和封禁列表
func (s *ServerService) GetLinkMembers() (*models.LinkMembers, error) {
	var out models.LinkMembers
	if err := linkManager.linkCall(http.MethodGet, "/members", nil, &out); err != nil {
		return nil, err
	}
	if out.Members == nil {
		out.Members = []models.LinkMember{}
	}
	if out.Bans == nil {
		out.Bans = []models.LinkBan{}
	}
	return &out, nil
}

// KickLinkMember 踢出联机房成员（吊销凭证，下次心跳时在节点里吊销），ban 时同时封禁
func (s *ServerService) KickLinkMember(req models.LinkKickRequest) error {
	if !steamIDRe.MatchString(req.SteamID) {
		return ErrInvalidSteamID
	}
	if err := linkManager.linkCall(http.MethodPost, "/kick", req, nil); err != nil {
		return err
	}
	linkManager.poke()
	return nil
}

// UnbanLinkMember 解除联机房封禁
func (s *ServerService) UnbanLinkMember(steamID string) error {
	if !steamIDRe.MatchString(steamID) {
		return ErrInvalidSteamID
	}
	return linkManager.linkCall(http.MethodDelete, "/bans/"+url.PathEscape(steamID), nil, nil)
}

// ---- Docker 操作 ----

// dockerLinkOps 用 docker_manager 实现 linkDocker
type dockerLinkOps struct{}

func (dockerLinkOps) dm() (*docker_manager.DockerManager, error) {
	return docker_manager.GetDockerManager()
}

func (o dockerLinkOps) ImageReady(image string) (bool, error) {
	dm, err := o.dm()
	if err != nil {
		return false, err
	}
	exists, err := dm.ImageExists(image)
	if err != nil || exists {
		return exists, err
	}
	if docker_manager.ClaimPull(image) {
		utils.Info("联机房使用的镜像不存在，开始下载", zap.String("image", image))
		go func() {
			if err := dm.PullClaimedImage(image); err != nil {
				utils.Warn("下载联机房镜像失败", zap.String("image", image), zap.Error(err))
			}
		}()
	}
	return false, nil
}

func (o dockerLinkOps) ImageID(image string) (string, error) {
	dm, err := o.dm()
	if err != nil {
		return "", err
	}
	return dm.GetImageID(image)
}

func (o dockerLinkOps) Inspect(name string) (docker_manager.LinkContainerInfo, error) {
	dm, err := o.dm()
	if err != nil {
		return docker_manager.LinkContainerInfo{}, err
	}
	return dm.InspectLinkContainer(name)
}

func (o dockerLinkOps) Create(spec docker_manager.LinkContainerSpec) error {
	dm, err := o.dm()
	if err != nil {
		return err
	}
	return dm.CreateLinkContainer(spec)
}

func (o dockerLinkOps) Remove(name string) error {
	dm, err := o.dm()
	if err != nil {
		return err
	}
	return dm.RemoveLinkContainer(name)
}

func (o dockerLinkOps) List() ([]string, error) {
	dm, err := o.dm()
	if err != nil {
		return nil, err
	}
	return dm.ListLinkContainers()
}

func (o dockerLinkOps) Exec(name string, argv []string, timeout time.Duration) (docker_manager.ExecResult, error) {
	dm, err := o.dm()
	if err != nil {
		return docker_manager.ExecResult{}, err
	}
	return dm.ExecCommand(name, argv, nil, timeout)
}

func (dockerLinkOps) RemoveDir() error { return docker_manager.RemoveLinkDir() }

func (dockerLinkOps) RemoveLegacyDir(serverID uint) error {
	return docker_manager.RemoveLegacyLinkDir(serverID)
}
