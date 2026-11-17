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

func (f *fakeListing) token(i int) string {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.tokens[i]
}

func newTestPublisher(t *testing.T, f *fakeListing, ready *bool, now *time.Time) *Publisher {
	t.Helper()
	utils.InitLogger()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		f.mu.Lock()
		f.requests = append(f.requests, r.URL.Path)
		f.tokens = append(f.tokens, strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer "))
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

func okHandler(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/publish":
		json.NewEncoder(w).Encode(map[string]any{"id": 9, "token": "afp_tok", "url": "https://arkfgf.com/servers"})
	default:
		json.NewEncoder(w).Encode(map[string]any{"ok": true})
	}
}

func errorHandler(status int, code string, retryAfter string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if retryAfter != "" {
			w.Header().Set("Retry-After", retryAfter)
		}
		w.WriteHeader(status)
		json.NewEncoder(w).Encode(map[string]string{"code": code, "error": code})
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

	// ready 后发布一次，保存令牌
	ready = true
	p.tick()
	expectCalls(t, f, "/publish")
	if got := loadPublishServer(t, s.ID).PublishToken; got != "afp_tok" {
		t.Fatalf("token not saved: %q", got)
	}
	if st := p.statusFor(loadPublishServer(t, s.ID)); st.State != models.PublishPublished || st.URL == "" {
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

	// 再次可以连接：重新发布，带上已保存的令牌
	ready = true
	database.DB.Model(&s).Update("status", "running")
	p.tick()
	expectCalls(t, f, "/publish", "/publish")
	if f.token(1) != "afp_tok" {
		t.Fatalf("republish should send saved token, got %q", f.token(1))
	}

	// 关闭发布：取消发布一次，清空令牌，之后不再请求
	database.DB.Model(&s).Update("publish", false)
	p.tick()
	p.tick()
	expectCalls(t, f, "/publish", "/publish", "/publish/unpublish")
	if f.token(2) != "afp_tok" {
		t.Fatalf("unpublish should send token, got %q", f.token(2))
	}
	if got := loadPublishServer(t, s.ID).PublishToken; got != "" {
		t.Fatalf("token should be cleared: %q", got)
	}
	if st := p.statusFor(loadPublishServer(t, s.ID)); st != nil {
		t.Fatalf("status should be nil when publish is off: %+v", st)
	}
}

func TestPublisherUnpublishDeleted(t *testing.T) {
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

	// 收录站已经没有这条记录：也算成功，清空令牌
	f.mu.Lock()
	f.handler = errorHandler(http.StatusUnauthorized, "invalid_token", "")
	f.mu.Unlock()
	now = now.Add(time.Minute)
	p.tick()
	p.tick()
	expectCalls(t, f, "/publish/unpublish", "/publish/unpublish")
	if got := loadPublishServer(t, s.ID).PublishToken; got != "" {
		t.Fatalf("token should be cleared: %q", got)
	}
}

func TestPublisherBackoff(t *testing.T) {
	setupTestDB(t)
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
	f.mu.Lock()
	f.handler = errorHandler(http.StatusTooManyRequests, "rate_limited", "1800")
	f.mu.Unlock()
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
}
