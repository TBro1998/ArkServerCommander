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

	"ark-server-commander/config"
	"ark-server-commander/database"
	"ark-server-commander/models"
	"ark-server-commander/service/docker_manager"
	"ark-server-commander/utils"

	"go.uber.org/zap"
)

// P2P 联机房（arkfgf.com 联机助手）。接口约定见外层仓库 docs/p2p-link.md 的"Commander 侧"。
//
// 开启 link 的服务器在网站上登记一个房间，并在服务器可以连接时运行一个 EasyTier 节点容器 ase-link-<id>
// （宿主机网络，--p2p-only，只放行游戏端口、+1 和查询端口的 UDP）。玩家用联机助手从网站拿一次性凭证加入
// 这个 EasyTier 网络，直连（打洞）到节点，再经宿主机 127.0.0.1 上 Docker 的端口映射到达 ARK 容器。
//
// LinkManager 每 15 秒检查一次：
//   - 没有登记或资料变化（名称、地图、Mod、倍率、收录资料、端口等）→ POST /link/rooms（带令牌即更新）
//   - 服务器运行中且可以连接过 → 确保节点容器在运行（参数或镜像变化时重建）；运行中短暂查询不到不删容器
//   - 心跳：在线时每 30 秒，离线时每 60 秒；在线状态变化、有凭证要上传或吊销要确认时提前发
//   - 按心跳返回的 want_credentials 在节点里签发凭证（每轮最多 20 个），处理 revoke 后在下次心跳回报 revoked_ack
//   - 关闭 link 或删除服务器（有令牌）→ DELETE 下架，删除容器、凭证目录，清空令牌
//   - 其他 ase-link-* 容器（服务器没在运行、已删除）一律删除
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
	linkServerIDLabel       = "com.arkservercommander.link.server"
	linkRPCPortBase         = 20000 // 节点 RPC 端口 = 20000 + 服务器 ID，和已知端口冲突时顺延
	linkNetworkSecretLength = 32
)

// linkCredential 已签发、等待随心跳上传的凭证
type linkCredential struct {
	ID        string `json:"id"`
	Secret    string `json:"secret"`
	ExpiresAt int64  `json:"expires_at"`
}

// linkRoomProfile 登记房间时提交的资料（POST /link/rooms）
type linkRoomProfile struct {
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
	Description  string   `json:"description"`
	Tags         []string `json:"tags"`
	QQGroup      string   `json:"qq_group"`
	Discord      string   `json:"discord"`
	GamePort     int      `json:"game_port"`
	QueryPort    int      `json:"query_port"`
	ServerIP     string   `json:"server_ip"`
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
	RemoveDir(serverID uint) error
}

// linkState 一台服务器的联机房状态，只在内存里；管理端重启后带着数据库里的令牌重新登记（拿回 peers 和房间地址）
type linkState struct {
	profileHash string   // 最近一次登记成功时的资料摘要，变化时重新登记
	url         string   // 网站上的房间页面
	peers       []string // 公共节点地址，由网站返回
	version     string   // 最近一次查到的游戏版本，停服后保留，避免资料来回变化

	// 网站请求失败的退避
	nextAt   time.Time
	failures int
	errCode  string
	errMsg   string

	// 心跳
	lastBeat   time.Time // 最近一次心跳成功的时间
	lastOnline bool      // 最近一次成功心跳报告的 online
	beatNow    bool      // 踢人后尽快发心跳，拿到要吊销的凭证
	want       int       // 网站要补充的凭证数
	revoke     []string  // 网站要求吊销、还没确认的凭证
	acks       []string  // 已在节点吊销、等待回报的凭证
	pending    []linkCredential

	// 节点
	wasReady      bool      // 这次运行期间已经可以连接过
	notReadySince time.Time // 运行中但查询不到的开始时间
	nodeErrCode   string
	nodeErr       string
}

// LinkManager 管理开启了联机房的服务器：登记房间、心跳、节点容器、凭证
type LinkManager struct {
	mu      sync.Mutex
	states  map[uint]*linkState
	client  *http.Client
	api     func() string
	now     func() time.Time
	docker  linkDocker
	runtime func(models.Server) linkRuntime
	profile func(models.Server) linkRoomProfile
	newID   func() string
}

var linkManager = &LinkManager{
	states:  map[uint]*linkState{},
	client:  &http.Client{Timeout: linkRequestTimeout},
	api:     func() string { return config.PublishAPI },
	now:     time.Now,
	docker:  dockerLinkOps{},
	runtime: linkRuntimeOf,
	profile: buildLinkRoomProfile,
	newID:   newUUID,
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
		for range ticker.C {
			linkManager.tick()
		}
	}()
}

func (m *LinkManager) tick() {
	desired := map[string]bool{}
	keep := map[uint]bool{}

	if m.api() != "" {
		var servers []models.Server
		if err := database.DB.Select("id", "session_name", "map", "max_players", "port", "query_port", "rcon_port", "cluster_id",
			"game_mod_ids", "arkapi", "server_args_json", "status", "link", "link_token", "link_room_id", "link_network_name",
			"link_network_secret", "link_rpc_port", "listing_description", "listing_tags", "listing_qq", "listing_discord").
			Where("link = ?", true).Find(&servers).Error; err != nil {
			utils.Warn("联机房：读取服务器列表失败", zap.Error(err))
			return
		}
		for _, server := range servers {
			keep[server.ID] = true
			if m.process(server) {
				desired[utils.GetLinkContainerName(server.ID)] = true
			}
		}

		// 关闭了联机房或已删除、但网站上还有房间（有令牌）的服务器：下架
		var withdrawn []models.Server
		if err := database.DB.Unscoped().Select("id", "link_token", "link_room_id").
			Where("link_token <> '' AND (link = ? OR deleted_at IS NOT NULL)", false).
			Find(&withdrawn).Error; err != nil {
			utils.Warn("联机房：读取待下架的服务器失败", zap.Error(err))
			return
		}
		for _, server := range withdrawn {
			keep[server.ID] = true
			m.unregister(server)
		}
	}

	// 不需要的节点容器（服务器没在运行、关闭了联机房、已删除）一律删除
	if names, err := m.docker.List(); err != nil {
		utils.Warn("联机房：列出节点容器失败", zap.Error(err))
	} else {
		for _, name := range names {
			if !desired[name] {
				if err := m.docker.Remove(name); err != nil {
					utils.Warn("联机房：删除节点容器失败", zap.String("container", name), zap.Error(err))
				}
			}
		}
	}

	m.mu.Lock()
	for id := range m.states {
		if !keep[id] {
			delete(m.states, id)
		}
	}
	m.mu.Unlock()
}

func (m *LinkManager) state(id uint) *linkState {
	m.mu.Lock()
	defer m.mu.Unlock()
	st, ok := m.states[id]
	if !ok {
		st = &linkState{}
		m.states[id] = st
	}
	return st
}

// process 处理一台开启了联机房的服务器，返回节点容器是否应该保留
func (m *LinkManager) process(server models.Server) bool {
	if err := ensureLinkLocal(&server); err != nil {
		utils.Warn("联机房：生成本地参数失败", zap.Uint("server_id", server.ID), zap.Error(err))
		return false
	}
	st := m.state(server.ID)
	rt := linkRuntime{}
	running := server.Status == "running"
	if running {
		rt = m.runtime(server)
	}

	m.mu.Lock()
	now := m.now()
	switch {
	case !running:
		st.wasReady, st.notReadySince = false, time.Time{}
		st.nodeErrCode, st.nodeErr = "", ""
	case rt.Ready:
		st.wasReady, st.notReadySince = true, time.Time{}
	case st.wasReady && st.notReadySince.IsZero():
		st.notReadySince = now
	}
	if rt.Version != "" {
		st.version = rt.Version
	}
	// 可以连接过一次后，整个运行期间保留节点（A2S 偶尔超时不删容器，玩家不掉线）；
	// 但查询不到超过 notReadyGrace（崩溃重启、卡死）时向网站报告离线
	wantNode := running && st.wasReady
	online := running && (rt.Ready || (st.wasReady && now.Sub(st.notReadySince) < notReadyGrace))
	version := st.version
	m.mu.Unlock()

	// 登记或更新房间
	profile := m.profile(server)
	profile.Version = version
	hash := hashJSON(profile)
	m.mu.Lock()
	needRegister := (server.LinkToken == "" || server.LinkRoomID == 0 || st.profileHash != hash) && !m.now().Before(st.nextAt)
	m.mu.Unlock()
	if needRegister {
		m.register(&server, st, profile, hash)
	}

	// 节点容器
	nodeUp := false
	if wantNode && server.LinkNetworkName != "" {
		m.mu.Lock()
		peers := slices.Clone(st.peers)
		m.mu.Unlock()
		if len(peers) > 0 {
			nodeUp = m.ensureNode(server, st, peers)
		}
	}

	if server.LinkToken != "" && server.LinkRoomID != 0 {
		m.heartbeat(&server, st, online && nodeUp, rt.Players, nodeUp)
	}
	return wantNode
}

// register 登记或更新房间（POST /link/rooms），保存令牌、房间 ID、网络名和 peers
func (m *LinkManager) register(server *models.Server, st *linkState, profile linkRoomProfile, hash string) {
	var resp struct {
		ID          uint     `json:"id"`
		Token       string   `json:"token"`
		NetworkName string   `json:"network_name"`
		Peers       []string `json:"peers"`
		URL         string   `json:"url"`
	}
	status, retryAfter, code, err := m.request(http.MethodPost, "/link/rooms", server.LinkToken, profile, &resp)
	if err != nil {
		m.fail(server.ID, status, retryAfter, code, err)
		return
	}
	if resp.ID == 0 && server.LinkRoomID == 0 {
		m.fail(server.ID, status, 0, "invalid_response", errors.New("网站没有返回房间 ID"))
		return
	}

	oldToken := server.LinkToken
	updates := map[string]any{}
	if resp.Token != "" && resp.Token != server.LinkToken {
		updates["link_token"] = resp.Token
		server.LinkToken = resp.Token
	}
	if resp.ID != 0 && resp.ID != server.LinkRoomID {
		updates["link_room_id"] = resp.ID
		server.LinkRoomID = resp.ID
	}
	if resp.NetworkName != "" && resp.NetworkName != server.LinkNetworkName {
		updates["link_network_name"] = resp.NetworkName
		server.LinkNetworkName = resp.NetworkName
	}
	if len(updates) > 0 {
		if err := database.DB.Model(&models.Server{}).Where("id = ?", server.ID).Updates(updates).Error; err != nil {
			// 令牌没保存下来，下次会再新建一个房间；网站上的旧房间 30 天没心跳后由网站清理
			m.fail(server.ID, 0, 0, "local_error", fmt.Errorf("保存房间令牌失败: %w", err))
			return
		}
	}
	newRoom := oldToken != "" && server.LinkToken != oldToken
	if newRoom {
		// 旧令牌已失效，网站新建了房间：旧房间的凭证（含被封禁玩家的）不能留在节点里
		m.resetNode(server.ID)
	}

	m.mu.Lock()
	if newRoom || oldToken == "" {
		st.lastBeat, st.lastOnline, st.want = time.Time{}, false, 0
		st.revoke, st.acks, st.pending = nil, nil, nil
	}
	st.profileHash = hash
	st.url = resp.URL
	if len(resp.Peers) > 0 {
		st.peers = resp.Peers
	}
	m.succeed(st)
	m.mu.Unlock()
	if oldToken == "" || newRoom {
		utils.Info("已登记联机房", zap.Uint("server_id", server.ID), zap.Uint("room_id", server.LinkRoomID))
	}
}

// resetNode 删除节点容器和凭证文件，下次按新参数重建
func (m *LinkManager) resetNode(serverID uint) {
	if err := m.docker.Remove(utils.GetLinkContainerName(serverID)); err != nil {
		utils.Warn("联机房：删除节点容器失败", zap.Uint("server_id", serverID), zap.Error(err))
	}
	if err := m.docker.RemoveDir(serverID); err != nil {
		utils.Warn("联机房：删除凭证目录失败", zap.Uint("server_id", serverID), zap.Error(err))
	}
}

// ensureNode 确保节点容器按当前参数运行，返回节点是否在运行
func (m *LinkManager) ensureNode(server models.Server, st *linkState, peers []string) bool {
	spec := linkContainerSpec(server, peers)
	name := utils.GetLinkContainerName(server.ID)

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
			utils.Info("联机房节点参数或镜像有变化，已重建", zap.Uint("server_id", server.ID))
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

// heartbeat 到期时发心跳；节点在运行时先吊销网站要求吊销的凭证、签发网站要补充的凭证，随心跳上传
func (m *LinkManager) heartbeat(server *models.Server, st *linkState, online bool, players int, nodeUp bool) {
	m.mu.Lock()
	now := m.now()
	interval := linkHeartbeatOffline
	if online {
		interval = linkHeartbeatOnline
	}
	work := nodeUp && (st.want > len(st.pending) || len(st.revoke) > 0)
	due := st.lastBeat.IsZero() || !now.Before(st.lastBeat.Add(interval)) || st.lastOnline != online || st.beatNow ||
		work || len(st.pending) > 0 || len(st.acks) > 0
	if !due || now.Before(st.nextAt) {
		m.mu.Unlock()
		return
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
		container := utils.GetLinkContainerName(server.ID)
		var newAcks []string
		for _, id := range revoke {
			if slices.Contains(acked, id) {
				continue
			}
			if err := m.revokeCredential(container, server.LinkRPCPort, id); err != nil {
				nodeErr = err
				break
			}
			newAcks = append(newAcks, id)
		}
		var issued []linkCredential
		for i := 0; i < toIssue && nodeErr == nil; i++ {
			cred, err := m.generateCredential(container, server.LinkRPCPort)
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
				utils.Warn("联机房：签发或吊销凭证失败", zap.Uint("server_id", server.ID), zap.Error(nodeErr))
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
		Players     int              `json:"players"`
		Credentials []linkCredential `json:"credentials"`
		RevokedAck  []string         `json:"revoked_ack"`
	}{Online: online, Credentials: []linkCredential{}, RevokedAck: []string{}}
	if online {
		body.Players = players
	}
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
	path := fmt.Sprintf("/link/rooms/%d/heartbeat", server.LinkRoomID)
	status, retryAfter, code, err := m.request(http.MethodPost, path, server.LinkToken, body, &resp)
	if err != nil {
		if code == "invalid_token" || status == http.StatusNotFound {
			m.roomGone(server, st)
		}
		m.fail(server.ID, status, retryAfter, code, err)
		return
	}
	if resp.NetworkName != "" && resp.NetworkName != server.LinkNetworkName {
		if err := database.DB.Model(&models.Server{}).Where("id = ?", server.ID).
			Update("link_network_name", resp.NetworkName).Error; err != nil {
			utils.Warn("联机房：保存网络名失败", zap.Uint("server_id", server.ID), zap.Error(err))
		} else {
			server.LinkNetworkName = resp.NetworkName
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
	m.succeed(st)
}

// roomGone 网站上已经没有这个房间（令牌无效）：清空令牌和节点，下一轮重新登记一个新房间
func (m *LinkManager) roomGone(server *models.Server, st *linkState) {
	utils.Warn("联机房：网站上的房间已不存在，将重新登记", zap.Uint("server_id", server.ID))
	if err := database.DB.Model(&models.Server{}).Where("id = ?", server.ID).
		Updates(map[string]any{"link_token": "", "link_room_id": 0, "link_network_name": ""}).Error; err != nil {
		utils.Warn("联机房：清除房间令牌失败", zap.Uint("server_id", server.ID), zap.Error(err))
		return
	}
	m.resetNode(server.ID)
	m.mu.Lock()
	st.profileHash, st.url = "", ""
	st.lastBeat, st.lastOnline, st.want = time.Time{}, false, 0
	st.revoke, st.acks, st.pending = nil, nil, nil
	m.mu.Unlock()
}

// unregister 下架房间（DELETE /link/rooms/:id），成功后删除节点容器和凭证目录，清空令牌。
// 网站已经没有这个房间（invalid_token / 404）也算成功
func (m *LinkManager) unregister(server models.Server) {
	st := m.state(server.ID)
	m.mu.Lock()
	wait := m.now().Before(st.nextAt)
	m.mu.Unlock()
	if wait {
		return
	}
	status, retryAfter, code, err := m.request(http.MethodDelete, fmt.Sprintf("/link/rooms/%d", server.LinkRoomID), server.LinkToken, nil, nil)
	if err != nil && code != "invalid_token" && status != http.StatusNotFound {
		m.fail(server.ID, status, retryAfter, code, err)
		return
	}
	m.resetNode(server.ID)
	if err := database.DB.Unscoped().Model(&models.Server{}).Where("id = ?", server.ID).
		Updates(map[string]any{"link_token": "", "link_room_id": 0, "link_network_name": "", "link_network_secret": ""}).Error; err != nil {
		utils.Warn("联机房：清除房间令牌失败", zap.Uint("server_id", server.ID), zap.Error(err))
		return
	}
	m.mu.Lock()
	delete(m.states, server.ID)
	m.mu.Unlock()
	utils.Info("已下架联机房", zap.Uint("server_id", server.ID))
}

// generateCredential 在节点里签发一个不可复用的凭证（TTL 30 天）
func (m *LinkManager) generateCredential(container string, rpcPort int) (linkCredential, error) {
	id := m.newID()
	out, err := m.cli(container, rpcPort, "credential", "generate", "--ttl", strconv.Itoa(int(linkCredentialTTL/time.Second)),
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
func (m *LinkManager) revokeCredential(container string, rpcPort int, id string) error {
	_, err := m.cli(container, rpcPort, "credential", "revoke", id)
	if err != nil && !strings.Contains(strings.ToLower(err.Error()), "not found") {
		return fmt.Errorf("吊销凭证失败: %w", err)
	}
	return nil
}

// cli 在节点容器里执行 easytier-cli（JSON 输出），返回标准输出
func (m *LinkManager) cli(container string, rpcPort int, args ...string) (string, error) {
	argv := append([]string{"easytier-cli", "-p", fmt.Sprintf("127.0.0.1:%d", rpcPort), "-o", "json"}, args...)
	res, err := m.docker.Exec(container, argv, linkExecTimeout)
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
func (m *LinkManager) succeed(st *linkState) {
	st.failures = 0
	st.errCode, st.errMsg = "", ""
	st.nextAt = time.Time{}
}

// fail 记录网站请求失败并安排重试：指数退避，网站给了 Retry-After 时取两者中较长的
func (m *LinkManager) fail(serverID uint, status int, retryAfter time.Duration, code string, err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	st, ok := m.states[serverID]
	if !ok {
		return
	}
	st.failures++
	wait := min(linkBackoffBase<<min(st.failures-1, 5), linkBackoffMax)
	if retryAfter > wait {
		wait = min(retryAfter, 2*time.Hour)
	}
	st.nextAt = m.now().Add(wait)
	if code == "" {
		code = "network"
	}
	st.errCode = code
	st.errMsg = linkErrorMessage(code)
	if st.errMsg == "" {
		st.errMsg = err.Error()
	}
	if st.failures == 1 || st.failures%10 == 0 {
		utils.Warn("联机房请求失败", zap.Uint("server_id", serverID), zap.Int("status", status),
			zap.String("code", code), zap.Int("failures", st.failures), zap.Error(err))
	}
}

func linkErrorMessage(code string) string {
	switch code {
	case "rate_limited", "busy":
		return "arkfgf.com 暂时繁忙或请求过于频繁，稍后自动重试"
	case "too_many_rooms":
		return "这个公网 IP 登记的联机房已达上限"
	case "invalid_token":
		return "网站上的房间已不存在，正在重新登记"
	case "invalid_input":
		return "房间资料不符合网站要求"
	}
	return ""
}

// statusFor 返回服务器的联机房状态，未开启时为 nil
func (m *LinkManager) statusFor(server models.Server) *models.LinkStatus {
	if !server.Link {
		return nil
	}
	if m.api() == "" {
		return &models.LinkStatus{State: models.LinkDisabled}
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	status := &models.LinkStatus{State: models.LinkWaiting}
	st, ok := m.states[server.ID]
	if !ok {
		return status
	}
	status.RoomURL = st.url
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
	case st.lastBeat.IsZero():
	case st.lastOnline:
		status.State = models.LinkOnline
	case server.Status != "running":
		status.State = models.LinkOffline
	}
	return status
}

// poke 让下一轮马上发心跳（踢人后尽快拿到要吊销的凭证）
func (m *LinkManager) poke(serverID uint) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if st, ok := m.states[serverID]; ok {
		st.beatNow = true
	}
}

// ---- 节点参数 ----

// linkNodeArgs easytier-core 的参数，见 docs/p2p-link.md 的"服主节点"
func linkNodeArgs(server models.Server, peers []string) []string {
	hostname := strings.TrimSpace(server.SessionName)
	if hostname == "" {
		hostname = fmt.Sprintf("ark-%d", server.ID)
	}
	args := []string{
		"--network-name", server.LinkNetworkName,
		"--network-secret", server.LinkNetworkSecret,
		"--secure-mode",
		"--credential-file", linkCredentialFile,
		"--peers",
	}
	args = append(args, peers...)
	args = append(args,
		"--ipv4", linkVirtualIP, "--no-tun", "--p2p-only",
		"--tcp-whitelist", "0",
		"--udp-whitelist", strconv.Itoa(server.Port), strconv.Itoa(server.Port+1), strconv.Itoa(server.QueryPort),
		"--listeners", "tcp://0.0.0.0:0", "udp://0.0.0.0:0",
		"--rpc-portal", fmt.Sprintf("127.0.0.1:%d", server.LinkRPCPort),
		"--hostname", hostname,
	)
	return args
}

// linkContainerSpec 节点容器配置；标签里的摘要包含镜像和全部参数，参数变化时重建
func linkContainerSpec(server models.Server, peers []string) docker_manager.LinkContainerSpec {
	args := linkNodeArgs(server, peers)
	sum := sha256.Sum256([]byte(utils.EasyTierImage + "\x00" + strings.Join(args, "\x00")))
	return docker_manager.LinkContainerSpec{
		ServerID: server.ID,
		Image:    utils.EasyTierImage,
		Args:     args,
		Labels: map[string]string{
			linkCommandHashLabel: hex.EncodeToString(sum[:]),
			linkServerIDLabel:    strconv.FormatUint(uint64(server.ID), 10),
		},
	}
}

// ensureLinkLocal 第一次处理时生成网络密钥、分配 RPC 端口并保存
func ensureLinkLocal(server *models.Server) error {
	updates := map[string]any{}
	if server.LinkNetworkSecret == "" {
		buf := make([]byte, linkNetworkSecretLength)
		if _, err := rand.Read(buf); err != nil {
			return err
		}
		server.LinkNetworkSecret = hex.EncodeToString(buf)
		updates["link_network_secret"] = server.LinkNetworkSecret
	}
	if server.LinkRPCPort == 0 {
		port, err := assignLinkRPCPort(server.ID)
		if err != nil {
			return err
		}
		server.LinkRPCPort = port
		updates["link_rpc_port"] = port
	}
	if len(updates) == 0 {
		return nil
	}
	return database.DB.Model(&models.Server{}).Where("id = ?", server.ID).Updates(updates).Error
}

// assignLinkRPCPort 节点 RPC 端口（只监听宿主机 127.0.0.1）：20000 + 服务器 ID。
// 管理端运行在容器里，探测不了宿主机端口，只避开数据库里已知的端口（各服务器的游戏端口、+1、查询、RCON 和其他节点的 RPC 端口），冲突时顺延
func assignLinkRPCPort(serverID uint) (int, error) {
	var servers []models.Server
	if err := database.DB.Select("id", "port", "query_port", "rcon_port", "link_rpc_port").Find(&servers).Error; err != nil {
		return 0, err
	}
	used := map[int]bool{}
	for _, s := range servers {
		used[s.Port], used[s.Port+1], used[s.QueryPort], used[s.RCONPort] = true, true, true, true
		if s.ID != serverID && s.LinkRPCPort > 0 {
			used[s.LinkRPCPort] = true
		}
	}
	port := linkRPCPortBase + int(serverID%40000)
	for used[port] {
		port++
		if port > 65535 {
			port = linkRPCPortBase
		}
	}
	return port, nil
}

// ---- 资料 ----

// buildLinkRoomProfile 登记房间的资料：和发布到收录站相同的服务器资料，加上端口和节点虚拟 IP
func buildLinkRoomProfile(server models.Server) linkRoomProfile {
	sp := buildServerProfile(server)
	return linkRoomProfile{
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
		Description:  sp.Listing.Description,
		Tags:         sp.Listing.Tags,
		QQGroup:      sp.Listing.QQGroup,
		Discord:      sp.Listing.Discord,
		GamePort:     server.Port,
		QueryPort:    server.QueryPort,
		ServerIP:     linkVirtualIP,
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

// ---- 成员管理（代理到网站） ----

var steamIDRe = regexp.MustCompile(`^[0-9]{17}$`)

// linkCall 用房间令牌调用网站的成员管理接口
func (m *LinkManager) linkCall(server models.Server, method, path string, body, out any) error {
	if m.api() == "" {
		return fmt.Errorf("%w：管理端已关闭与 arkfgf.com 的连接（ARKFGF_API 为空）", ErrLinkUnavailable)
	}
	if server.LinkToken == "" || server.LinkRoomID == 0 {
		return ErrLinkNotRegistered
	}
	status, _, code, err := m.request(method, fmt.Sprintf("/link/rooms/%d%s", server.LinkRoomID, path), server.LinkToken, body, out)
	if err == nil {
		return nil
	}
	switch {
	case code == "invalid_token" || status == http.StatusUnauthorized:
		return ErrLinkNotRegistered
	case status == http.StatusBadRequest || status == http.StatusNotFound || status == http.StatusConflict || status == http.StatusUnprocessableEntity:
		return fmt.Errorf("%w：%v", ErrLinkRequestFailed, err)
	}
	return fmt.Errorf("%w：%v", ErrLinkUnavailable, err)
}

// GetLinkMembers 联机房的成员（持有未吊销凭证的玩家）和封禁列表
func (s *ServerService) GetLinkMembers(userID uint, serverID string) (*models.LinkMembers, error) {
	server, err := getUserServer(userID, serverID)
	if err != nil {
		return nil, err
	}
	var out models.LinkMembers
	if err := linkManager.linkCall(server, http.MethodGet, "/members", nil, &out); err != nil {
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
func (s *ServerService) KickLinkMember(userID uint, serverID string, req models.LinkKickRequest) error {
	if !steamIDRe.MatchString(req.SteamID) {
		return ErrInvalidSteamID
	}
	server, err := getUserServer(userID, serverID)
	if err != nil {
		return err
	}
	if err := linkManager.linkCall(server, http.MethodPost, "/kick", req, nil); err != nil {
		return err
	}
	linkManager.poke(server.ID)
	return nil
}

// UnbanLinkMember 解除联机房封禁
func (s *ServerService) UnbanLinkMember(userID uint, serverID, steamID string) error {
	if !steamIDRe.MatchString(steamID) {
		return ErrInvalidSteamID
	}
	server, err := getUserServer(userID, serverID)
	if err != nil {
		return err
	}
	return linkManager.linkCall(server, http.MethodDelete, "/bans/"+url.PathEscape(steamID), nil, nil)
}

// removeLinkContainer 删除服务器时顺手删掉节点容器（下架由后台的 tick 完成）
func removeLinkContainer(serverID uint) {
	if err := linkManager.docker.Remove(utils.GetLinkContainerName(serverID)); err != nil {
		utils.Warn("删除联机房容器失败", zap.Uint("server_id", serverID), zap.Error(err))
	}
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

func (dockerLinkOps) RemoveDir(serverID uint) error { return docker_manager.RemoveLinkDir(serverID) }
