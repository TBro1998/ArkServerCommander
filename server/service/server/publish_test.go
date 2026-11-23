package server

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"ark-server-commander/database"
	"ark-server-commander/models"
	"ark-server-commander/utils"
)

// fakeListing 模拟收录站：记录收到的请求，按 handler 返回
type fakeListing struct {
	mu       sync.Mutex
	requests []string
	tokens   []string
	bodies   []map[string]any
	handler  func(w http.ResponseWriter, r *http.Request)
}

func (f *fakeListing) calls() []string {
	f.mu.Lock()
	defer f.mu.Unlock()
	return append([]string{}, f.requests...)
}

func (f *fakeListing) token(i int) string {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.tokens[i]
}

func (f *fakeListing) body(i int) map[string]any {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.bodies[i]
}

func (f *fakeListing) setHandler(h func(w http.ResponseWriter, r *http.Request)) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.handler = h
}

func newTestPublisher(t *testing.T, f *fakeListing, ready *bool, now *time.Time) *Publisher {
	t.Helper()
	utils.InitLogger()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, _ := io.ReadAll(r.Body)
		var body map[string]any
		_ = json.Unmarshal(data, &body)
		r.Body = io.NopCloser(bytes.NewReader(data))
		f.mu.Lock()
		f.requests = append(f.requests, r.URL.Path)
		f.tokens = append(f.tokens, strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer "))
		f.bodies = append(f.bodies, body)
		handler := f.handler
		f.mu.Unlock()
		handler(w, r)
	}))
	t.Cleanup(srv.Close)
	return &Publisher{
		states:  map[uint]*publishState{},
		client:  srv.Client(),
		api:     func() string { return srv.URL },
		now:     func() time.Time { return *now },
		ready:   func(models.Server) bool { return *ready },
		profile: func(s models.Server) publishProfile { return publishProfile{Port: s.Port, QueryPort: s.QueryPort} },
	}
}

// okHandler 模拟新版收录站：发布、取消发布和统计都要账号令牌，取消发布还接受旧版发布令牌
func okHandler(w http.ResponseWriter, r *http.Request) {
	token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
	body := bodyOf(r)
	switch {
	case r.URL.Path == "/publish/unpublish" && strings.HasPrefix(token, "afp_"):
		json.NewEncoder(w).Encode(map[string]any{"ok": true})
	case token != testAccountToken:
		writeAPIError(w, http.StatusUnauthorized, "login_required")
	case r.URL.Path == "/publish":
		if body["key"] == nil || body["token"] != nil {
			writeAPIError(w, http.StatusBadRequest, "invalid_input")
			return
		}
		json.NewEncoder(w).Encode(map[string]any{"id": 9, "url": "https://arkfgf.com/servers/9", "badge": "https://arkfgf.com/badge/9.svg"})
	case r.URL.Path == "/publish/unpublish":
		if body["key"] == nil {
			writeAPIError(w, http.StatusBadRequest, "invalid_input")
			return
		}
		json.NewEncoder(w).Encode(map[string]any{"ok": true})
	case r.URL.Path == "/publish/players":
		if body["servers"] == nil {
			writeAPIError(w, http.StatusBadRequest, "invalid_input")
			return
		}
		json.NewEncoder(w).Encode(map[string]any{"updated": len(body["servers"].([]any))})
	case r.URL.Path == "/publish/stats":
		if r.URL.Query().Get("key") == "" {
			writeAPIError(w, http.StatusNotFound, "not_published")
			return
		}
		json.NewEncoder(w).Encode(map[string]any{"id": 9, "url": "https://arkfgf.com/servers/9", "week": map[string]int{"views": 3}})
	default:
		writeAPIError(w, http.StatusNotFound, "not_found")
	}
}

// bodyOf 解析请求体（JSON 对象），没有时为空
func bodyOf(r *http.Request) map[string]any {
	var body map[string]any
	_ = json.NewDecoder(r.Body).Decode(&body)
	return body
}

func writeAPIError(w http.ResponseWriter, status int, code string) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"code": code, "error": code})
}

func errorHandler(status int, code string, retryAfter string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if retryAfter != "" {
			w.Header().Set("Retry-After", retryAfter)
		}
		writeAPIError(w, status, code)
	}
}

func loadPublishServer(t *testing.T, id uint) models.Server {
	t.Helper()
	var s models.Server
	if err := database.DB.Unscoped().First(&s, id).Error; err != nil {
		t.Fatal(err)
	}
	return s
}

func expectCalls(t *testing.T, f *fakeListing, want ...string) {
	t.Helper()
	got := f.calls()
	if strings.Join(got, ",") != strings.Join(want, ",") {
		t.Fatalf("calls = %v, want %v", got, want)
	}
}

func TestPublisherLifecycle(t *testing.T) {
	setupTestDB(t)
	bindTestAccount(t)
	s := createTestServer(t, 7777, 27015, 32330)
	database.DB.Model(&s).Updates(map[string]any{"publish": true, "status": "running"})

	ready := false
	now := time.Now()
	f := &fakeListing{handler: okHandler}
	p := newTestPublisher(t, f, &ready, &now)

	// 还没 ready：不发请求，状态为等待
	p.tick()
	expectCalls(t, f)
	if st := p.statusFor(loadPublishServer(t, s.ID)); st.State != models.PublishWaiting {
		t.Fatalf("state = %s", st.State)
	}

	// ready 后用账号令牌发布一次，带 key，记下 published
	ready = true
	p.tick()
	expectCalls(t, f, "/publish")
	if f.token(0) != testAccountToken || f.body(0)["key"] != publishKey(s.ID) || f.body(0)["legacy_token"] != nil {
		t.Fatalf("publish request: token %q body %v", f.token(0), f.body(0))
	}
	if got := loadPublishServer(t, s.ID); !got.Published || got.PublishToken != "" {
		t.Fatalf("published not saved: %+v", got)
	}
	if st := p.statusFor(loadPublishServer(t, s.ID)); st.State != models.PublishPublished || st.URL == "" || st.Badge == "" {
		t.Fatalf("status = %+v", st)
	}

	// 之后不再发请求（没有心跳）
	for range 50 {
		now = now.Add(time.Minute)
		p.tick()
	}
	expectCalls(t, f, "/publish")

	// 运行中短暂查询不到：宽限期内恢复不算断开，不重新发布
	ready = false
	now = now.Add(time.Minute)
	p.tick()
	ready = true
	now = now.Add(15 * time.Second)
	p.tick()
	expectCalls(t, f, "/publish")

	// 停止：不发请求（收录站自己会查到离线）
	ready = false
	database.DB.Model(&s).Update("status", "stopped")
	p.tick()
	expectCalls(t, f, "/publish")
	if st := p.statusFor(loadPublishServer(t, s.ID)); st.State != models.PublishWaiting {
		t.Fatalf("stopped server should be waiting: %+v", st)
	}

	// 再次可以连接：重新发布
	ready = true
	database.DB.Model(&s).Update("status", "running")
	p.tick()
	expectCalls(t, f, "/publish", "/publish")
	if f.token(1) != testAccountToken {
		t.Fatalf("republish should use the account token, got %q", f.token(1))
	}

	// 关闭发布：用 key 取消发布一次，published 置 false，之后不再请求
	database.DB.Model(&s).Update("publish", false)
	p.tick()
	p.tick()
	expectCalls(t, f, "/publish", "/publish", "/publish/unpublish")
	if f.token(2) != testAccountToken || f.body(2)["key"] != publishKey(s.ID) {
		t.Fatalf("unpublish request: token %q body %v", f.token(2), f.body(2))
	}
	if got := loadPublishServer(t, s.ID); got.Published {
		t.Fatal("published should be cleared")
	}
	if st := p.statusFor(loadPublishServer(t, s.ID)); st != nil {
		t.Fatalf("status should be nil when publish is off: %+v", st)
	}
}

// 未绑定账号：开启发布的服务器状态为 login_required，不发请求；绑定后正常发布
func TestPublisherReportsPlayers(t *testing.T) {
	setupTestDB(t)
	bindTestAccount(t)
	a := createTestServer(t, 7777, 27015, 32330)
	b := createTestServer(t, 7779, 27017, 32332)
	database.DB.Model(&models.Server{}).Where("id IN ?", []uint{a.ID, b.ID}).Updates(map[string]any{"publish": true, "status": "running"})

	ready := true
	now := time.Now()
	f := &fakeListing{handler: okHandler}
	p := newTestPublisher(t, f, &ready, &now)
	// a 有 2 个玩家；b 的 RCON 连不上，跳过
	p.players = func(s models.Server) ([]Player, error) {
		if s.ID == b.ID {
			return nil, ErrRconUnavailable
		}
		return []Player{{Name: "Alice", SteamID: "76561198000000001"}, {Name: "Bob", SteamID: "76561198000000002"}}, nil
	}

	// 发布后同一轮就上报一次在线玩家（ID 和名称），只带列得出来的
	p.tick()
	expectCalls(t, f, "/publish", "/publish", "/publish/players")
	entries := f.body(2)["servers"].([]any)
	if f.token(2) != testAccountToken || len(entries) != 1 || entries[0].(map[string]any)["key"] != publishKey(a.ID) {
		t.Fatalf("players request: %v", f.body(2))
	}
	list := entries[0].(map[string]any)["players"].([]any)
	if len(list) != 2 || list[0].(map[string]any)["id"] != "76561198000000001" || list[1].(map[string]any)["name"] != "Bob" {
		t.Fatalf("players list: %v", list)
	}

	// 间隔内不再上报，到时间再报
	now = now.Add(playersReportInterval - time.Second)
	p.tick()
	expectCalls(t, f, "/publish", "/publish", "/publish/players")
	now = now.Add(time.Second)
	p.tick()
	expectCalls(t, f, "/publish", "/publish", "/publish/players", "/publish/players")

	// 不能连接的（停服、断开）不上报；一台都没有时不发请求
	ready = false
	now = now.Add(playersReportInterval + notReadyGrace)
	p.tick()
	expectCalls(t, f, "/publish", "/publish", "/publish/players", "/publish/players")
}

func TestPublisherUnboundLoginRequired(t *testing.T) {
	setupTestDB(t)
	s := createTestServer(t, 7777, 27015, 32330)
	database.DB.Model(&s).Updates(map[string]any{"publish": true, "status": "running"})

	ready := true
	now := time.Now()
	f := &fakeListing{handler: okHandler}
	p := newTestPublisher(t, f, &ready, &now)

	for range 3 {
		p.tick()
		now = now.Add(time.Minute)
	}
	expectCalls(t, f)
	if st := p.statusFor(loadPublishServer(t, s.ID)); st.State != models.PublishLoginRequired || st.Error == "" {
		t.Fatalf("status = %+v", st)
	}

	bindTestAccount(t)
	p.tick()
	expectCalls(t, f, "/publish")
	if st := p.statusFor(loadPublishServer(t, s.ID)); st.State != models.PublishPublished {
		t.Fatalf("status = %+v", st)
	}
}

// 旧版匿名发布的令牌：发布时作为 legacy_token 只带一次，成功后清空
func TestPublisherLegacyTokenSentOnce(t *testing.T) {
	setupTestDB(t)
	bindTestAccount(t)
	s := createTestServer(t, 7777, 27015, 32330)
	database.DB.Model(&s).Updates(map[string]any{"publish": true, "status": "running", "publish_token": "afp_old"})

	ready := true
	now := time.Now()
	f := &fakeListing{handler: okHandler}
	p := newTestPublisher(t, f, &ready, &now)

	p.tick()
	expectCalls(t, f, "/publish")
	if f.token(0) != testAccountToken || f.body(0)["legacy_token"] != "afp_old" {
		t.Fatalf("first publish: token %q body %v", f.token(0), f.body(0))
	}
	if got := loadPublishServer(t, s.ID); got.PublishToken != "" || !got.Published {
		t.Fatalf("legacy token should be cleared: %+v", got)
	}

	// 停服再启动：重新发布时不再带 legacy_token
	ready = false
	database.DB.Model(&s).Update("status", "stopped")
	p.tick()
	ready = true
	database.DB.Model(&s).Update("status", "running")
	p.tick()
	expectCalls(t, f, "/publish", "/publish")
	if _, ok := f.body(1)["legacy_token"]; ok {
		t.Fatalf("legacy_token sent again: %v", f.body(1))
	}
}

// 关闭发布或删除后，还留着旧版发布令牌的服务器不需要绑定账号也用旧令牌取消发布
func TestPublisherLegacyUnpublishWithoutAccount(t *testing.T) {
	setupTestDB(t)
	s := createTestServer(t, 7777, 27015, 32330)
	database.DB.Model(&s).Updates(map[string]any{"publish": true, "publish_token": "afp_old"})
	database.DB.Delete(&models.Server{}, s.ID)

	ready := false
	now := time.Now()
	f := &fakeListing{handler: errorHandler(http.StatusServiceUnavailable, "busy", "")}
	p := newTestPublisher(t, f, &ready, &now)

	// 失败后退避，期间不重试
	p.tick()
	now = now.Add(30 * time.Second)
	p.tick()
	expectCalls(t, f, "/publish/unpublish")
	if f.token(0) != "afp_old" {
		t.Fatalf("legacy unpublish token = %q", f.token(0))
	}

	// 收录站已经没有这条记录：也算成功，清空令牌
	f.setHandler(errorHandler(http.StatusUnauthorized, "invalid_token", ""))
	now = now.Add(time.Minute)
	p.tick()
	p.tick()
	expectCalls(t, f, "/publish/unpublish", "/publish/unpublish")
	if got := loadPublishServer(t, s.ID).PublishToken; got != "" {
		t.Fatalf("token should be cleared: %q", got)
	}
}

// 删除的服务器：已用账号发布过的用 key 取消发布
func TestPublisherUnpublishDeleted(t *testing.T) {
	setupTestDB(t)
	bindTestAccount(t)
	s := createTestServer(t, 7777, 27015, 32330)
	database.DB.Model(&s).Updates(map[string]any{"publish": true, "published": true})
	database.DB.Delete(&models.Server{}, s.ID)

	ready := false
	now := time.Now()
	f := &fakeListing{handler: okHandler}
	p := newTestPublisher(t, f, &ready, &now)
	p.tick()
	p.tick()
	expectCalls(t, f, "/publish/unpublish")
	if f.body(0)["key"] != publishKey(s.ID) {
		t.Fatalf("body = %v", f.body(0))
	}
	if loadPublishServer(t, s.ID).Published {
		t.Fatal("published should be cleared")
	}
}

// 收录站返回 invalid_account：视为已解绑，清空本地账号和所有服务器的 published
func TestPublisherInvalidAccountClearsLocal(t *testing.T) {
	setupTestDB(t)
	bindTestAccount(t)
	useTestLinkManager(t)
	s := createTestServer(t, 7777, 27015, 32330)
	other := createTestServer(t, 7779, 27017, 32332)
	database.DB.Model(&s).Updates(map[string]any{"publish": true, "status": "running"})
	database.DB.Model(&other).Updates(map[string]any{"published": true})

	ready := true
	now := time.Now()
	f := &fakeListing{handler: errorHandler(http.StatusUnauthorized, "invalid_account", "")}
	p := newTestPublisher(t, f, &ready, &now)
	old := publisher
	publisher = p
	t.Cleanup(func() { publisher = old })

	p.tick()
	expectCalls(t, f, "/publish")
	if accountManager.bound() {
		t.Fatal("account should be cleared")
	}
	var n int64
	database.DB.Model(&models.ArkfgfAccount{}).Count(&n)
	if n != 0 {
		t.Fatal("account row should be deleted")
	}
	if loadPublishServer(t, other.ID).Published {
		t.Fatal("published should be cleared for all servers")
	}
	if st := p.statusFor(loadPublishServer(t, s.ID)); st.State != models.PublishLoginRequired {
		t.Fatalf("status = %+v", st)
	}
	// 之后不再发请求
	now = now.Add(time.Hour)
	p.tick()
	expectCalls(t, f, "/publish")
}

func TestPublisherBackoff(t *testing.T) {
	setupTestDB(t)
	bindTestAccount(t)
	s := createTestServer(t, 7777, 27015, 32330)
	database.DB.Model(&s).Updates(map[string]any{"publish": true, "status": "running"})

	ready := true
	now := time.Now()
	f := &fakeListing{handler: errorHandler(http.StatusUnprocessableEntity, "unreachable", "")}
	p := newTestPublisher(t, f, &ready, &now)

	// 失败后按 1、2、4 分钟退避，期间不重试
	expect := 0
	for _, wait := range []time.Duration{time.Minute, 2 * time.Minute, 4 * time.Minute} {
		p.tick()
		expect++
		now = now.Add(wait - time.Second)
		p.tick()
		if got := len(f.calls()); got != expect {
			t.Fatalf("retried before backoff %v: %d calls", wait, got)
		}
		now = now.Add(time.Second)
	}
	st := p.statusFor(loadPublishServer(t, s.ID))
	if st.State != models.PublishError || st.ErrorCode != "unreachable" || !strings.Contains(st.Error, "27015") || st.RetryAt == "" {
		t.Fatalf("status = %+v", st)
	}

	// 停服再启动不会跳过退避
	ready = false
	database.DB.Model(&s).Update("status", "stopped")
	p.tick()
	ready = true
	database.DB.Model(&s).Update("status", "running")
	p.tick()
	if got := len(f.calls()); got != expect+1 {
		// 上一轮循环结束时正好到期，这里允许发一次
		t.Fatalf("expected one retry at the due time, got %d calls", got)
	}
	p.tick()
	if got := len(f.calls()); got != expect+1 {
		t.Fatalf("restart should not skip backoff: %d calls", got)
	}

	// 429 时按 Retry-After 等待（比退避时间长）
	f.setHandler(errorHandler(http.StatusTooManyRequests, "rate_limited", "1800"))
	now = now.Add(time.Hour)
	p.tick()
	calls := len(f.calls())
	now = now.Add(29 * time.Minute)
	p.tick()
	if len(f.calls()) != calls {
		t.Fatal("should wait for Retry-After")
	}
	now = now.Add(2 * time.Minute)
	p.tick()
	if len(f.calls()) != calls+1 {
		t.Fatal("should retry after Retry-After")
	}
}

func TestBuildPublishProfileArgsOverride(t *testing.T) {
	s := models.Server{
		ID: 123456, SessionName: "S", Map: "TheIsland", Port: 7777, QueryPort: 27015, GameModIds: "1, 2,,3",
		ServerArgsJSON: `{"query_params":{"XPMultiplier":"2.5","serverPVE":"True"}}`,
	}
	p := buildPublishProfile(s)
	if p.RateXP != 2.5 || p.PvP || p.RateTame != 1 || len(p.Mods) != 3 || p.Mods[2] != "3" {
		t.Fatalf("profile = %+v", p)
	}
	// 发布请求把 key 和 legacy_token 放在同一层
	data, _ := json.Marshal(publishRequest{publishProfile: p, Key: "123456"})
	var m map[string]any
	_ = json.Unmarshal(data, &m)
	if m["key"] != "123456" || m["port"] != float64(7777) {
		t.Fatalf("request = %s", data)
	}
	if _, ok := m["legacy_token"]; ok {
		t.Fatalf("empty legacy_token should be omitted: %s", data)
	}
}
