package server

import (
	"encoding/json"
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
	handler  func(w http.ResponseWriter, r *http.Request)
}

func (f *fakeListing) calls() []string {
	f.mu.Lock()
	defer f.mu.Unlock()
	return append([]string{}, f.requests...)
}

func newTestPublisher(t *testing.T, f *fakeListing, ready *bool, now *time.Time) *Publisher {
	t.Helper()
	utils.InitLogger()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		f.mu.Lock()
		f.requests = append(f.requests, r.URL.Path)
		f.tokens = append(f.tokens, strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer "))
		f.mu.Unlock()
		f.handler(w, r)
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

func okHandler(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/publish":
		json.NewEncoder(w).Encode(map[string]any{"id": 9, "token": "afp_tok", "heartbeat_interval": 300, "url": "https://arkfgf.com/servers"})
	case "/publish/heartbeat":
		json.NewEncoder(w).Encode(map[string]any{"ok": true, "online": true, "heartbeat_interval": 300})
	default:
		json.NewEncoder(w).Encode(map[string]any{"ok": true})
	}
}

func loadPublishServer(t *testing.T, id uint) models.Server {
	t.Helper()
	var s models.Server
	if err := database.DB.First(&s, id).Error; err != nil {
		t.Fatal(err)
	}
	return s
}

func TestPublisherLifecycle(t *testing.T) {
	setupTestDB(t)
	s := createTestServer(t, 7777, 27015, 32330)
	database.DB.Model(&s).Updates(map[string]any{"publish": true, "status": "running"})

	ready := false
	now := time.Now()
	f := &fakeListing{handler: okHandler}
	p := newTestPublisher(t, f, &ready, &now)

	// 还没 ready：不发请求，状态为等待
	p.process(loadPublishServer(t, s.ID))
	if len(f.calls()) != 0 {
		t.Fatalf("should not publish before ready: %v", f.calls())
	}
	if st := p.statusFor(loadPublishServer(t, s.ID)); st.State != models.PublishWaiting {
		t.Fatalf("state = %s", st.State)
	}

	// ready 后发布一次，保存令牌
	ready = true
	p.process(loadPublishServer(t, s.ID))
	if got := loadPublishServer(t, s.ID).PublishToken; got != "afp_tok" {
		t.Fatalf("token not saved: %q", got)
	}
	if st := p.statusFor(loadPublishServer(t, s.ID)); st.State != models.PublishPublished || st.URL == "" {
		t.Fatalf("status = %+v", st)
	}

	// 心跳间隔内不再发请求
	for i := 0; i < 5; i++ {
		now = now.Add(30 * time.Second)
		p.process(loadPublishServer(t, s.ID))
	}
	if calls := f.calls(); len(calls) != 1 {
		t.Fatalf("expected only the publish call, got %v", calls)
	}
	// 到时间后发心跳，带令牌
	now = now.Add(3 * time.Minute)
	p.process(loadPublishServer(t, s.ID))
	if calls := f.calls(); len(calls) != 2 || calls[1] != "/publish/heartbeat" || f.tokens[1] != "afp_tok" {
		t.Fatalf("expected heartbeat with token, got %v %v", calls, f.tokens)
	}

	// 运行中短暂查询不到：宽限期内不下线
	ready = false
	now = now.Add(time.Minute)
	p.process(loadPublishServer(t, s.ID))
	if len(f.calls()) != 2 {
		t.Fatalf("should not go offline within grace: %v", f.calls())
	}
	ready = true
	now = now.Add(15 * time.Second)
	p.process(loadPublishServer(t, s.ID))
	if len(f.calls()) != 2 {
		t.Fatalf("recovered server should keep heartbeat schedule: %v", f.calls())
	}

	// 停止：通知下线一次，之后不再请求
	ready = false
	database.DB.Model(&s).Update("status", "stopped")
	p.process(loadPublishServer(t, s.ID))
	p.process(loadPublishServer(t, s.ID))
	if calls := f.calls(); len(calls) != 3 || calls[2] != "/publish/offline" {
		t.Fatalf("expected one offline call, got %v", calls)
	}

	// 再次启动：重新发布，带上已保存的令牌
	ready = true
	database.DB.Model(&s).Update("status", "running")
	p.process(loadPublishServer(t, s.ID))
	if calls := f.calls(); len(calls) != 4 || calls[3] != "/publish" || f.tokens[3] != "afp_tok" {
		t.Fatalf("expected republish with token, got %v %v", calls, f.tokens)
	}

	// 关闭发布：通知下线，状态为空
	database.DB.Model(&s).Update("publish", false)
	p.process(loadPublishServer(t, s.ID))
	if calls := f.calls(); len(calls) != 5 || calls[4] != "/publish/offline" {
		t.Fatalf("expected offline after disabling, got %v", calls)
	}
	if st := p.statusFor(loadPublishServer(t, s.ID)); st != nil {
		t.Fatalf("status should be nil when publish is off: %+v", st)
	}
}

func TestPublisherBackoff(t *testing.T) {
	setupTestDB(t)
	s := createTestServer(t, 7777, 27015, 32330)
	database.DB.Model(&s).Updates(map[string]any{"publish": true, "status": "running"})

	ready := true
	now := time.Now()
	f := &fakeListing{handler: func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(map[string]string{"code": "unreachable", "error": "no A2S response"})
	}}
	p := newTestPublisher(t, f, &ready, &now)

	// 失败后按 1、2、4 分钟退避，期间不重试
	expect := 0
	for _, wait := range []time.Duration{time.Minute, 2 * time.Minute, 4 * time.Minute} {
		p.process(loadPublishServer(t, s.ID))
		expect++
		now = now.Add(wait - time.Second)
		p.process(loadPublishServer(t, s.ID))
		if got := len(f.calls()); got != expect {
			t.Fatalf("retried before backoff %v: %d calls", wait, got)
		}
		now = now.Add(time.Second)
	}
	st := p.statusFor(loadPublishServer(t, s.ID))
	if st.State != models.PublishError || st.ErrorCode != "unreachable" || !strings.Contains(st.Error, "27015") || st.RetryAt == "" {
		t.Fatalf("status = %+v", st)
	}

	// 429 时按 Retry-After 等待（比退避时间长）
	f.handler = func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Retry-After", "1800")
		w.WriteHeader(http.StatusTooManyRequests)
		json.NewEncoder(w).Encode(map[string]string{"code": "rate_limited"})
	}
	p.process(loadPublishServer(t, s.ID))
	calls := len(f.calls())
	now = now.Add(29 * time.Minute)
	p.process(loadPublishServer(t, s.ID))
	if len(f.calls()) != calls {
		t.Fatal("should wait for Retry-After")
	}
}

func TestPublisherHeartbeatInvalidToken(t *testing.T) {
	setupTestDB(t)
	s := createTestServer(t, 7777, 27015, 32330)
	database.DB.Model(&s).Updates(map[string]any{"publish": true, "status": "running"})

	ready := true
	now := time.Now()
	f := &fakeListing{handler: okHandler}
	p := newTestPublisher(t, f, &ready, &now)
	p.process(loadPublishServer(t, s.ID))

	// 收录站清理了记录：心跳返回 invalid_token 后下一轮重新发布
	f.handler = func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/publish/heartbeat" {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"code": "invalid_token"})
			return
		}
		okHandler(w, r)
	}
	now = now.Add(6 * time.Minute)
	p.process(loadPublishServer(t, s.ID))
	p.process(loadPublishServer(t, s.ID))
	if calls := f.calls(); len(calls) != 3 || calls[2] != "/publish" {
		t.Fatalf("expected republish after invalid token, got %v", calls)
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
}
