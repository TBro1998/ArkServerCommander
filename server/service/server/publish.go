package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
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

// 发布到服务器收录站（arkfgf.com）。
//
// 开启 publish 的服务器每次变成 ready（查询端口有响应、玩家可以连接）时发布一次，之后不再发请求：
// 收录站每半小时自己向查询端口发 A2S_INFO 更新在线状态和人数，停服后查不到就显示为离线。
// 关闭发布或删除服务器时调用 /publish/unpublish 删除收录站上的记录，成功前会一直重试（令牌存在数据库里，管理端重启后继续）。
// 收录站发布时也会从它那边查询端口验证，所以查询端口（UDP）必须能从公网访问。
//
// 防止刷屏：
//   - 每次可以连接只发布一次；运行中短暂查询不到（notReadyGrace 内）不算断开，不会重新发布
//   - 失败后指数退避（1 分钟起，最长 30 分钟），收录站返回 429 / 503 时按 Retry-After 等待；退避时间在停服后保留
//   - 请求串行发出，同一时刻只有一个请求
//   - 收录站那边还有按 IP、端口、全站的限流，见 web/backend/controllers/publish.go

const (
	publishTickInterval   = 15 * time.Second
	publishBackoffBase    = time.Minute
	publishBackoffMax     = 30 * time.Minute
	publishRequestTimeout = 15 * time.Second // 收录站会同步做 A2S 查询（最多 3 秒），留足余量
	// 运行中但暂时查询不到（监控的一次 A2S 超时会让阶段短暂变回 launching）时，超过这个时间才算断开，
	// 避免恢复后马上重新发布，触发收录站的发布频率限制
	notReadyGrace = 2 * time.Minute
)

// publishState 一台服务器的发布状态，只在内存里；管理端重启后带着数据库里的令牌重新发布
type publishState struct {
	published bool      // 这一次可以连接之后已经发布过
	nextAt    time.Time // 失败后下一次重试的时间
	failures  int
	errCode   string
	errMsg    string
	lastOK    time.Time
	url       string
	badge     string // 状态徽章图片地址（旧版收录站不返回）
	// 运行中但不是 ready 的开始时间，见 notReadyGrace
	notReadySince time.Time
}

// Publisher 定期检查开启了发布的服务器：可以连接时发布，关闭发布或删除后取消发布
type Publisher struct {
	mu     sync.Mutex
	states map[uint]*publishState
	client *http.Client
	api    func() string
	now    func() time.Time
	// ready 判断服务器能否连接；profile 生成发布内容。测试时替换
	ready   func(models.Server) bool
	profile func(models.Server) publishProfile
}

var publisher = &Publisher{
	states:  map[uint]*publishState{},
	client:  &http.Client{Timeout: publishRequestTimeout},
	api:     func() string { return config.PublishAPI },
	now:     time.Now,
	ready:   serverReady,
	profile: buildPublishProfile,
}

// StartPublisher 启动后台发布协程。ARKFGF_API 为空时不启动
func StartPublisher() {
	if publisher.api() == "" {
		utils.Info("服务器收录站发布已关闭（ARKFGF_API 为空）")
		return
	}
	go func() {
		ticker := time.NewTicker(publishTickInterval)
		defer ticker.Stop()
		for range ticker.C {
			publisher.tick()
		}
	}()
}

func serverReady(server models.Server) bool {
	if server.Status != "running" {
		return false
	}
	info := runtimeMonitor.Get(server.ID)
	return info != nil && info.Phase == models.PhaseReady
}

func (p *Publisher) tick() {
	if p.api() == "" {
		return
	}
	keep := map[uint]bool{}

	var servers []models.Server
	if err := database.DB.Select("id", "session_name", "map", "max_players", "port", "query_port", "cluster_id",
		"game_mod_ids", "arkapi", "server_args_json", "status", "publish", "publish_token",
		"listing_description", "listing_tags", "listing_qq", "listing_discord").
		Where("publish = ?", true).Find(&servers).Error; err != nil {
		utils.Warn("发布到收录站：读取服务器列表失败", zap.Error(err))
		return
	}
	for _, server := range servers {
		keep[server.ID] = true
		p.process(server)
	}

	// 关闭了发布或已删除、但收录站上还有记录（有令牌）的服务器：取消发布
	var withdrawn []models.Server
	if err := database.DB.Unscoped().Select("id", "query_port", "publish_token").
		Where("publish_token <> '' AND (publish = ? OR deleted_at IS NOT NULL)", false).
		Find(&withdrawn).Error; err != nil {
		utils.Warn("发布到收录站：读取待取消发布的服务器失败", zap.Error(err))
	}
	for _, server := range withdrawn {
		keep[server.ID] = true
		p.unpublish(server)
	}

	p.mu.Lock()
	for id := range p.states {
		if !keep[id] {
			delete(p.states, id)
		}
	}
	p.mu.Unlock()
}

// process 处理一台开启了发布的服务器：这一次可以连接后还没发布过就发布
func (p *Publisher) process(server models.Server) {
	st := p.state(server.ID)
	ready := p.ready(server)

	p.mu.Lock()
	now := p.now()
	if !ready {
		if server.Status == "running" && st.published {
			if st.notReadySince.IsZero() {
				st.notReadySince = now
			}
			if now.Sub(st.notReadySince) < notReadyGrace {
				p.mu.Unlock()
				return
			}
		}
		// 停止或长时间不能连接：下次可以连接时重新发布。失败退避的时间保留，避免反复启停时连续请求
		st.published = false
		st.notReadySince = time.Time{}
		st.errCode, st.errMsg = "", ""
		p.mu.Unlock()
		return
	}
	st.notReadySince = time.Time{}
	due := !st.published && !now.Before(st.nextAt)
	p.mu.Unlock()
	if due {
		p.publish(server)
	}
}

func (p *Publisher) state(id uint) *publishState {
	p.mu.Lock()
	defer p.mu.Unlock()
	st, ok := p.states[id]
	if !ok {
		st = &publishState{}
		p.states[id] = st
	}
	return st
}

// publishProfile 发布内容。名称、地图、人数和是否有密码收录站会以 A2S 查询结果为准，这里的值只作补充
type publishProfile struct {
	Port             int      `json:"port"`
	QueryPort        int      `json:"query_port"`
	Name             string   `json:"name"`
	Map              string   `json:"map"`
	MaxPlayers       int      `json:"max_players"`
	PvP              bool     `json:"pvp"`
	ArkApi           bool     `json:"arkapi"`
	Cluster          string   `json:"cluster,omitempty"`
	Mods             []string `json:"mods"`
	RateXP           float64  `json:"rate_xp"`
	RateTame         float64  `json:"rate_tame"`
	RateHarvest      float64  `json:"rate_harvest"`
	RateBreeding     float64  `json:"rate_breeding"`
	CommanderVersion string   `json:"commander_version"`
	// 收录资料，总是发送（空值表示清空）；收录站对没有这些字段的旧版请求保持原值
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
	QQGroup     string   `json:"qq_group"`
	Discord     string   `json:"discord"`
}

func (p *Publisher) publish(server models.Server) {
	profile := p.profile(server)
	var resp struct {
		ID    uint   `json:"id"`
		Token string `json:"token"`
		URL   string `json:"url"`
		Badge string `json:"badge"`
	}
	status, retryAfter, code, err := p.post("/publish", server.PublishToken, profile, &resp)
	if err != nil {
		p.fail(server, status, retryAfter, code, err)
		return
	}
	if resp.Token != "" {
		if err := database.DB.Model(&models.Server{}).Where("id = ?", server.ID).Update("publish_token", resp.Token).Error; err != nil {
			utils.Warn("保存收录站发布令牌失败", zap.Uint("server_id", server.ID), zap.Error(err))
		}
	}
	p.mu.Lock()
	st := p.states[server.ID]
	st.published = true
	st.url = resp.URL
	st.badge = resp.Badge
	p.succeed(st)
	p.mu.Unlock()
	utils.Info("已发布到服务器收录站", zap.Uint("server_id", server.ID), zap.Uint("listing_id", resp.ID))
}

// unpublish 删除收录站上的记录，成功后清空令牌。收录站已经没有这条记录（invalid_token）也算成功
func (p *Publisher) unpublish(server models.Server) {
	st := p.state(server.ID)
	p.mu.Lock()
	wait := p.now().Before(st.nextAt)
	p.mu.Unlock()
	if wait {
		return
	}
	status, retryAfter, code, err := p.post("/publish/unpublish", server.PublishToken, nil, nil)
	if err != nil && code != "invalid_token" {
		p.fail(server, status, retryAfter, code, err)
		return
	}
	if err := database.DB.Unscoped().Model(&models.Server{}).Where("id = ?", server.ID).Update("publish_token", "").Error; err != nil {
		utils.Warn("清除收录站发布令牌失败", zap.Uint("server_id", server.ID), zap.Error(err))
		return
	}
	p.mu.Lock()
	delete(p.states, server.ID)
	p.mu.Unlock()
	utils.Info("已从服务器收录站取消发布", zap.Uint("server_id", server.ID))
}

func (p *Publisher) succeed(st *publishState) {
	st.failures = 0
	st.errCode, st.errMsg = "", ""
	st.lastOK = p.now()
	st.nextAt = time.Time{}
}

// fail 记录失败并安排重试：指数退避，收录站给了 Retry-After 时取两者中较长的
func (p *Publisher) fail(server models.Server, status int, retryAfter time.Duration, code string, err error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	st := p.states[server.ID]
	st.failures++
	wait := publishBackoffBase << min(st.failures-1, 5)
	wait = min(wait, publishBackoffMax)
	if retryAfter > wait {
		wait = min(retryAfter, 2*time.Hour)
	}
	st.nextAt = p.now().Add(wait)
	if code == "" {
		code = "network"
	}
	st.errCode = code
	st.errMsg = publishErrorMessage(code, server)
	if code == "network" || st.errMsg == "" {
		st.errMsg = err.Error()
	}
	if st.failures == 1 || st.failures%10 == 0 {
		utils.Warn("发布到服务器收录站失败", zap.Uint("server_id", server.ID), zap.Int("status", status),
			zap.String("code", code), zap.Int("failures", st.failures), zap.Error(err))
	}
}

// post 向收录站发 JSON 请求。返回 HTTP 状态码、Retry-After、错误码和错误
func (p *Publisher) post(path, token string, body any, out any) (int, time.Duration, string, error) {
	var buf bytes.Buffer
	if body != nil {
		if err := json.NewEncoder(&buf).Encode(body); err != nil {
			return 0, 0, "", err
		}
	}
	req, err := http.NewRequest(http.MethodPost, p.api()+path, &buf)
	if err != nil {
		return 0, 0, "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "ArkServerCommander/"+utils.Version)
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	resp, err := p.client.Do(req)
	if err != nil {
		return 0, 0, "", fmt.Errorf("连接收录站失败: %w", err)
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(io.LimitReader(resp.Body, 64<<10))
	if err != nil {
		return resp.StatusCode, 0, "", fmt.Errorf("读取收录站响应失败: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		var e struct {
			Error string `json:"error"`
			Code  string `json:"code"`
		}
		_ = json.Unmarshal(data, &e)
		retry := time.Duration(0)
		if secs, err := strconv.Atoi(resp.Header.Get("Retry-After")); err == nil && secs > 0 {
			retry = time.Duration(secs) * time.Second
		}
		msg := e.Error
		if msg == "" {
			msg = http.StatusText(resp.StatusCode)
		}
		return resp.StatusCode, retry, e.Code, fmt.Errorf("收录站返回 %d: %s", resp.StatusCode, msg)
	}
	if out != nil {
		if err := json.Unmarshal(data, out); err != nil {
			return resp.StatusCode, 0, "", fmt.Errorf("解析收录站响应失败: %w", err)
		}
	}
	return resp.StatusCode, 0, "", nil
}

func publishErrorMessage(code string, server models.Server) string {
	switch code {
	case "unreachable":
		return fmt.Sprintf("收录站无法通过查询端口 %d（UDP）访问服务器，请在防火墙或云服务器安全组中放行该端口", server.QueryPort)
	case "not_ark":
		return "查询端口返回的不是 ARK 服务器"
	case "ip_not_public":
		return "发布请求不是从公网 IPv4 地址发出的，收录站无法确认服务器地址"
	case "rate_limited", "busy":
		return "收录站暂时繁忙或请求过于频繁，稍后自动重试"
	case "too_many_servers":
		return "这个公网 IP 发布的服务器已达上限"
	case "owned":
		return "这台服务器已由 arkfgf.com 上的服主账号管理"
	case "invalid_input":
		return "发布内容不符合收录站要求"
	}
	return ""
}

// statusFor 返回服务器的发布状态，未开启发布时为 nil
func (p *Publisher) statusFor(server models.Server) *models.PublishStatus {
	if !server.Publish {
		return nil
	}
	if p.api() == "" {
		return &models.PublishStatus{State: models.PublishDisabled}
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	st, ok := p.states[server.ID]
	status := &models.PublishStatus{State: models.PublishWaiting}
	if !ok {
		return status
	}
	status.URL = st.url
	status.Badge = st.badge
	if !st.lastOK.IsZero() {
		status.PublishedAt = st.lastOK.Format(time.RFC3339)
	}
	switch {
	case st.errCode != "":
		status.State = models.PublishError
		status.ErrorCode, status.Error = st.errCode, st.errMsg
		if st.failures > 0 && !st.nextAt.IsZero() {
			status.RetryAt = st.nextAt.Format(time.RFC3339)
		}
	case st.published:
		status.State = models.PublishPublished
	}
	return status
}

// buildPublishProfile 从服务器记录、启动参数和配置文件得到发布内容。
// 倍率和 PvE 先看启动参数里的 ?Key=Value（优先级高于配置文件），再看 GameUserSettings.ini / Game.ini
func buildPublishProfile(server models.Server) publishProfile {
	profile := publishProfile{
		Port:             server.Port,
		QueryPort:        server.QueryPort,
		Name:             server.SessionName,
		Map:              server.Map,
		MaxPlayers:       server.MaxPlayers,
		ArkApi:           server.ArkApi,
		Cluster:          server.ClusterID,
		Mods:             []string{},
		CommanderVersion: utils.Version,
	}
	listing := listingOf(server)
	profile.Description, profile.Tags, profile.QQGroup, profile.Discord = listing.Description, listing.Tags, listing.QQGroup, listing.Discord
	for _, id := range strings.Split(server.GameModIds, ",") {
		if id = strings.TrimSpace(id); id != "" {
			profile.Mods = append(profile.Mods, id)
		}
	}

	var args models.ServerArgs
	_ = json.Unmarshal([]byte(server.ServerArgsJSON), &args)
	var gus, game string
	if dm, err := docker_manager.GetDockerManager(); err == nil {
		gus, _ = dm.ReadConfigFile(server.ID, "GameUserSettings.ini")
		game, _ = dm.ReadConfigFile(server.ID, "Game.ini")
	}
	lookup := func(key, file, section string) string {
		for k, v := range args.QueryParams {
			if strings.EqualFold(k, key) {
				return v
			}
		}
		value := ""
		for _, e := range utils.FindIniEntries(file, key) {
			if strings.EqualFold(e.Section, section) {
				value = e.Value // 同名键以最后一个为准
			}
		}
		return value
	}
	const gusSection, gameSection = "ServerSettings", "/script/shootergame.shootergamemode"
	profile.PvP = !parseIniBool(lookup("serverPVE", gus, gusSection))
	profile.RateXP = parseRate(lookup("XPMultiplier", gus, gusSection))
	profile.RateTame = parseRate(lookup("TamingSpeedMultiplier", gus, gusSection))
	profile.RateHarvest = parseRate(lookup("HarvestAmountMultiplier", gus, gusSection))
	profile.RateBreeding = parseRate(lookup("BabyMatureSpeedMultiplier", game, gameSection))
	return profile
}

func parseIniBool(v string) bool {
	return strings.EqualFold(strings.TrimSpace(v), "true") || strings.TrimSpace(v) == "1"
}

// parseRate 没有设置时游戏按 1 倍处理
func parseRate(v string) float64 {
	f, err := strconv.ParseFloat(strings.TrimSpace(v), 64)
	if err != nil || math.IsNaN(f) || math.IsInf(f, 0) || f <= 0 {
		return 1
	}
	return f
}
