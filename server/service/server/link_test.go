package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"slices"
	"strings"
	"sync"
	"testing"
	"time"

	"ark-server-commander/database"
	"ark-server-commander/models"
	"ark-server-commander/service/docker_manager"
	"ark-server-commander/utils"
)

// fakeLinkContainer 模拟的节点容器
type fakeLinkContainer struct {
	running bool
	imageID string
	labels  map[string]string
	args    []string
}

// fakeLinkDocker 模拟 Docker：记录创建、删除和 easytier-cli 调用
type fakeLinkDocker struct {
	mu          sync.Mutex
	imageReady  bool
	imageID     string
	containers  map[string]*fakeLinkContainer
	creates     int
	removed     []string
	removedDirs []uint
	generated   []string
	revoked     []string
	execErr     error
}

func newFakeLinkDocker() *fakeLinkDocker {
	return &fakeLinkDocker{imageReady: true, imageID: "sha256:et1", containers: map[string]*fakeLinkContainer{}}
}

func (d *fakeLinkDocker) ImageReady(string) (bool, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.imageReady, nil
}

func (d *fakeLinkDocker) ImageID(string) (string, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.imageID, nil
}

func (d *fakeLinkDocker) Inspect(name string) (docker_manager.LinkContainerInfo, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	c, ok := d.containers[name]
	if !ok {
		return docker_manager.LinkContainerInfo{}, nil
	}
	return docker_manager.LinkContainerInfo{Exists: true, Running: c.running, ImageID: c.imageID, Labels: c.labels}, nil
}

func (d *fakeLinkDocker) Create(spec docker_manager.LinkContainerSpec) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.creates++
	d.containers[utils.GetLinkContainerName(spec.ServerID)] = &fakeLinkContainer{running: true, imageID: d.imageID, labels: spec.Labels, args: spec.Args}
	return nil
}

func (d *fakeLinkDocker) Remove(name string) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	if _, ok := d.containers[name]; ok {
		delete(d.containers, name)
		d.removed = append(d.removed, name)
	}
	return nil
}

func (d *fakeLinkDocker) List() ([]string, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	var names []string
	for name := range d.containers {
		names = append(names, name)
	}
	return names, nil
}

func (d *fakeLinkDocker) Exec(name string, argv []string, _ time.Duration) (docker_manager.ExecResult, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.execErr != nil {
		return docker_manager.ExecResult{}, d.execErr
	}
	if _, ok := d.containers[name]; !ok {
		return docker_manager.ExecResult{}, fmt.Errorf("no such container %s", name)
	}
	if !slices.Equal(argv[:5], []string{"easytier-cli", "-p", argv[2], "-o", "json"}) || !strings.HasPrefix(argv[2], "127.0.0.1:") {
		return docker_manager.ExecResult{ExitCode: 2, Stderr: "bad args"}, nil
	}
	switch {
	case slices.Equal(argv[5:7], []string{"credential", "generate"}):
		i := slices.Index(argv, "--credential-id")
		if i < 0 || !slices.Contains(argv, "--reusable") || !slices.Contains(argv, "2592000") {
			return docker_manager.ExecResult{ExitCode: 2, Stderr: "bad args"}, nil
		}
		id := argv[i+1]
		d.generated = append(d.generated, id)
		out, _ := json.Marshal(map[string]string{"credential_id": id, "credential_secret": "secret-" + id})
		return docker_manager.ExecResult{Stdout: string(out)}, nil
	case slices.Equal(argv[5:7], []string{"credential", "revoke"}):
		d.revoked = append(d.revoked, argv[7])
		return docker_manager.ExecResult{Stdout: `{"success": true}`}, nil
	}
	return docker_manager.ExecResult{ExitCode: 2, Stderr: "unknown command"}, nil
}

func (d *fakeLinkDocker) RemoveDir(id uint) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.removedDirs = append(d.removedDirs, id)
	return nil
}

func (d *fakeLinkDocker) container(name string) *fakeLinkContainer {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.containers[name]
}

// fakeLinkSite 模拟网站的联机房接口
type fakeLinkSite struct {
	mu       sync.Mutex
	calls    []string // "METHOD path"
	tokens   []string
	beats    []linkBeat
	profiles []linkRoomProfile
	want     int
	revoke   []string
	// fail 不为空时所有请求按它返回
	fail func(w http.ResponseWriter)
}

type linkBeat struct {
	Online      bool             `json:"online"`
	Players     int              `json:"players"`
	Credentials []linkCredential `json:"credentials"`
	RevokedAck  []string         `json:"revoked_ack"`
}

func (f *fakeLinkSite) handler(w http.ResponseWriter, r *http.Request) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.calls = append(f.calls, r.Method+" "+r.URL.Path)
	token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
	f.tokens = append(f.tokens, token)
	if f.fail != nil {
		f.fail(w)
		return
	}
	peers := []string{"tcp://et.arkfgf.com:11010", "udp://et.arkfgf.com:11010"}
	switch {
	case r.Method == http.MethodPost && r.URL.Path == "/link/rooms":
		var p linkRoomProfile
		_ = json.NewDecoder(r.Body).Decode(&p)
		f.profiles = append(f.profiles, p)
		resp := map[string]any{"id": 12, "network_name": "ark-0123456789abcdef", "peers": peers, "url": "https://arkfgf.com/rooms/12"}
		if token != "alr_tok" {
			resp["token"] = "alr_tok"
		}
		json.NewEncoder(w).Encode(resp)
	case r.Method == http.MethodPost && r.URL.Path == "/link/rooms/12/heartbeat":
		var b linkBeat
		_ = json.NewDecoder(r.Body).Decode(&b)
		f.beats = append(f.beats, b)
		// 网站收到凭证后减少缺口，收到确认后从 revoke 里去掉
		f.want = max(f.want-len(b.Credentials), 0)
		f.revoke = slices.DeleteFunc(f.revoke, func(id string) bool { return slices.Contains(b.RevokedAck, id) })
		json.NewEncoder(w).Encode(map[string]any{"want_credentials": f.want, "revoke": f.revoke, "network_name": "ark-0123456789abcdef", "peers": peers})
	case r.Method == http.MethodDelete && r.URL.Path == "/link/rooms/12":
		json.NewEncoder(w).Encode(map[string]any{"ok": true})
	case r.Method == http.MethodGet && r.URL.Path == "/link/rooms/12/members":
		json.NewEncoder(w).Encode(map[string]any{"members": []map[string]any{{"steam_id": "76561198000000001", "name": "p1", "credential_id": "c1", "joined_at": 1797600000}}, "bans": []any{}})
	case r.Method == http.MethodPost && r.URL.Path == "/link/rooms/12/kick":
		json.NewEncoder(w).Encode(map[string]any{"ok": true})
	default:
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"code": "not_found", "error": "not found"})
	}
}

func (f *fakeLinkSite) callList() []string {
	f.mu.Lock()
	defer f.mu.Unlock()
	return slices.Clone(f.calls)
}

func (f *fakeLinkSite) lastBeat() linkBeat {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.beats[len(f.beats)-1]
}

func (f *fakeLinkSite) beatCount() int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return len(f.beats)
}

func (f *fakeLinkSite) set(fn func(f *fakeLinkSite)) {
	f.mu.Lock()
	defer f.mu.Unlock()
	fn(f)
}

func newTestLinkManager(t *testing.T, site *fakeLinkSite, dock *fakeLinkDocker, rt *linkRuntime, now *time.Time) *LinkManager {
	t.Helper()
	utils.InitLogger()
	srv := httptest.NewServer(http.HandlerFunc(site.handler))
	t.Cleanup(srv.Close)
	n := 0
	return &LinkManager{
		states:  map[uint]*linkState{},
		client:  srv.Client(),
		api:     func() string { return srv.URL },
		now:     func() time.Time { return *now },
		docker:  dock,
		runtime: func(models.Server) linkRuntime { return *rt },
		profile: func(s models.Server) linkRoomProfile {
			return linkRoomProfile{Name: s.SessionName, Map: s.Map, GamePort: s.Port, QueryPort: s.QueryPort, ServerIP: linkVirtualIP}
		},
		newID: func() string { n++; return fmt.Sprintf("00000000-0000-4000-8000-%012d", n) },
	}
}

func setupLinkTest(t *testing.T) (models.Server, *fakeLinkSite, *fakeLinkDocker, *linkRuntime, *time.Time, *LinkManager) {
	setupTestDB(t)
	s := createTestServer(t, 7777, 27015, 32330)
	database.DB.Model(&s).Update("link", true)
	site := &fakeLinkSite{}
	dock := newFakeLinkDocker()
	rt := &linkRuntime{}
	now := time.Now()
	m := newTestLinkManager(t, site, dock, rt, &now)
	return s, site, dock, rt, &now, m
}

func expectLinkCalls(t *testing.T, site *fakeLinkSite, want ...string) {
	t.Helper()
	if got := site.callList(); strings.Join(got, ",") != strings.Join(want, ",") {
		t.Fatalf("calls = %v\nwant    %v", got, want)
	}
}

func TestLinkRegisterHeartbeatAndNode(t *testing.T) {
	s, site, dock, rt, now, m := setupLinkTest(t)
	name := utils.GetLinkContainerName(s.ID)

	// 已停止：登记房间并保存令牌等，发一次离线心跳，不创建节点
	m.tick()
	expectLinkCalls(t, site, "POST /link/rooms", "POST /link/rooms/12/heartbeat")
	if site.tokens[0] != "" || site.tokens[1] != "alr_tok" {
		t.Fatalf("tokens = %v", site.tokens)
	}
	got := loadPublishServer(t, s.ID)
	if got.LinkToken != "alr_tok" || got.LinkRoomID != 12 || got.LinkNetworkName != "ark-0123456789abcdef" {
		t.Fatalf("room not saved: %+v", got)
	}
	if len(got.LinkNetworkSecret) != 64 || got.LinkRPCPort != linkRPCPortBase+int(s.ID) {
		t.Fatalf("local params: secret %q port %d", got.LinkNetworkSecret, got.LinkRPCPort)
	}
	if b := site.lastBeat(); b.Online || b.Players != 0 || len(b.Credentials) != 0 {
		t.Fatalf("offline beat = %+v", b)
	}
	if st := m.statusFor(got); st.State != models.LinkOffline || st.RoomURL != "https://arkfgf.com/rooms/12" {
		t.Fatalf("status = %+v", st)
	}
	if dock.creates != 0 {
		t.Fatal("node should not run while stopped")
	}

	// 离线时每 60 秒心跳一次
	*now = now.Add(30 * time.Second)
	m.tick()
	if site.beatCount() != 1 {
		t.Fatal("offline heartbeat before 60s")
	}
	*now = now.Add(30 * time.Second)
	m.tick()
	if site.beatCount() != 2 {
		t.Fatal("offline heartbeat missing at 60s")
	}

	// 运行中但还没 ready：不创建节点，状态为等待
	database.DB.Model(&s).Update("status", "running")
	*now = now.Add(15 * time.Second)
	m.tick()
	if dock.creates != 0 {
		t.Fatal("node should wait until ready")
	}
	if st := m.statusFor(loadPublishServer(t, s.ID)); st.State != models.LinkWaiting {
		t.Fatalf("status = %+v", st)
	}

	// ready：创建节点，立即发在线心跳
	*rt = linkRuntime{Ready: true, Players: 3}
	*now = now.Add(15 * time.Second)
	m.tick()
	c := dock.container(name)
	if c == nil || dock.creates != 1 {
		t.Fatalf("node not created: %d", dock.creates)
	}
	args := strings.Join(c.args, " ")
	for _, want := range []string{
		"--network-name ark-0123456789abcdef", "--network-secret " + got.LinkNetworkSecret, "--secure-mode",
		"--credential-file /data/cred.json", "--peers tcp://et.arkfgf.com:11010 udp://et.arkfgf.com:11010 --ipv4",
		"--ipv4 10.126.126.1 --no-tun --p2p-only", "--tcp-whitelist 0 --udp-whitelist 7777 7778 27015",
		"--listeners tcp://0.0.0.0:0 udp://0.0.0.0:0", fmt.Sprintf("--rpc-portal 127.0.0.1:%d", got.LinkRPCPort), "--hostname test",
	} {
		if !strings.Contains(args, want) {
			t.Fatalf("args missing %q: %s", want, args)
		}
	}
	if b := site.lastBeat(); !b.Online || b.Players != 3 {
		t.Fatalf("online beat = %+v", b)
	}
	if st := m.statusFor(loadPublishServer(t, s.ID)); st.State != models.LinkOnline {
		t.Fatalf("status = %+v", st)
	}

	// 在线时每 30 秒心跳一次，节点不重复创建
	beats := site.beatCount()
	*now = now.Add(15 * time.Second)
	m.tick()
	if site.beatCount() != beats {
		t.Fatal("online heartbeat before 30s")
	}
	*now = now.Add(15 * time.Second)
	m.tick()
	if site.beatCount() != beats+1 || dock.creates != 1 {
		t.Fatalf("beats %d creates %d", site.beatCount(), dock.creates)
	}

	// 短暂查询不到：节点保留，宽限期内仍报告在线
	*rt = linkRuntime{}
	*now = now.Add(15 * time.Second)
	m.tick()
	if dock.container(name) == nil {
		t.Fatal("node removed on a short A2S miss")
	}
	*rt = linkRuntime{Ready: true, Players: 2}

	// 资料变化（改名）：带令牌重新登记，节点参数（hostname）变化，重建容器
	database.DB.Model(&s).Update("session_name", "renamed")
	*now = now.Add(15 * time.Second)
	m.tick()
	calls := site.callList()
	regs := 0
	for _, c := range calls {
		if c == "POST /link/rooms" {
			regs++
		}
	}
	if regs != 2 || site.tokens[slices.Index(calls[1:], "POST /link/rooms")+1] != "alr_tok" {
		t.Fatalf("expected a re-registration with token: %v %v", calls, site.tokens)
	}
	if site.profiles[1].Name != "renamed" {
		t.Fatalf("profile = %+v", site.profiles[1])
	}
	if dock.creates != 2 || !slices.Contains(dock.container(name).args, "renamed") {
		t.Fatalf("node not recreated after args change: %d", dock.creates)
	}

	// 镜像更新（镜像 ID 变化）也重建
	dock.mu.Lock()
	dock.imageID = "sha256:et2"
	dock.mu.Unlock()
	*now = now.Add(15 * time.Second)
	m.tick()
	if dock.creates != 3 {
		t.Fatalf("node not recreated after image update: %d", dock.creates)
	}

	// 停服：删除节点，立即发离线心跳
	database.DB.Model(&s).Update("status", "stopped")
	*rt = linkRuntime{}
	beats = site.beatCount()
	*now = now.Add(15 * time.Second)
	m.tick()
	if dock.container(name) != nil {
		t.Fatal("node should be removed when the server stops")
	}
	if site.beatCount() != beats+1 || site.lastBeat().Online {
		t.Fatalf("expected offline beat: %+v", site.lastBeat())
	}
	if st := m.statusFor(loadPublishServer(t, s.ID)); st.State != models.LinkOffline {
		t.Fatalf("status = %+v", st)
	}
}

func TestLinkCredentialsAndRevoke(t *testing.T) {
	s, site, dock, rt, now, m := setupLinkTest(t)
	database.DB.Model(&s).Update("status", "running")
	*rt = linkRuntime{Ready: true}
	site.set(func(f *fakeLinkSite) { f.want = 25 })

	// 第一轮：登记、建节点、在线心跳，网站要 25 个凭证
	m.tick()
	if site.beatCount() != 1 || len(site.lastBeat().Credentials) != 0 {
		t.Fatalf("first beat = %+v", site.lastBeat())
	}

	// 下一轮签发 20 个（每轮上限），随心跳上传
	*now = now.Add(15 * time.Second)
	m.tick()
	b := site.lastBeat()
	if len(b.Credentials) != 20 || len(dock.generated) != 20 {
		t.Fatalf("credentials = %d generated = %d", len(b.Credentials), len(dock.generated))
	}
	c := b.Credentials[0]
	if c.ID != dock.generated[0] || c.Secret != "secret-"+c.ID {
		t.Fatalf("credential = %+v", c)
	}
	if want := now.Add(30 * 24 * time.Hour).Unix(); c.ExpiresAt != want {
		t.Fatalf("expires_at = %d, want %d", c.ExpiresAt, want)
	}

	// 再下一轮补剩下的 5 个
	*now = now.Add(15 * time.Second)
	m.tick()
	if b := site.lastBeat(); len(b.Credentials) != 5 || len(dock.generated) != 25 {
		t.Fatalf("credentials = %d generated = %d", len(b.Credentials), len(dock.generated))
	}

	// 池满后不再提前心跳
	beats := site.beatCount()
	*now = now.Add(15 * time.Second)
	m.tick()
	if site.beatCount() != beats {
		t.Fatal("no extra heartbeat expected when the pool is full")
	}

	// 踢人：网站下次心跳返回 revoke；之后在节点里吊销，再下一次心跳回报 revoked_ack
	site.set(func(f *fakeLinkSite) { f.revoke = []string{"cred-x"} })
	srv := loadPublishServer(t, s.ID)
	srv.UserID = 1
	if err := m.linkCall(srv, http.MethodPost, "/kick", models.LinkKickRequest{SteamID: "76561198000000001"}, nil); err != nil {
		t.Fatal(err)
	}
	m.poke(s.ID)
	*now = now.Add(15 * time.Second)
	m.tick() // 提前心跳，拿到 revoke
	*now = now.Add(15 * time.Second)
	m.tick() // 吊销并回报
	if !slices.Equal(dock.revoked, []string{"cred-x"}) {
		t.Fatalf("revoked = %v", dock.revoked)
	}
	if b := site.lastBeat(); !slices.Equal(b.RevokedAck, []string{"cred-x"}) {
		t.Fatalf("revoked_ack = %v", b.RevokedAck)
	}
	// 确认后不再重复吊销和回报
	*now = now.Add(30 * time.Second)
	m.tick()
	if len(dock.revoked) != 1 || len(site.lastBeat().RevokedAck) != 0 {
		t.Fatalf("revoked = %v ack = %v", dock.revoked, site.lastBeat().RevokedAck)
	}

	// 节点出错：状态显示错误，心跳照发（不上传凭证）
	site.set(func(f *fakeLinkSite) { f.want = 3 })
	*now = now.Add(30 * time.Second)
	m.tick()
	dock.mu.Lock()
	dock.execErr = errors.New("rpc down")
	dock.mu.Unlock()
	beats = site.beatCount()
	*now = now.Add(15 * time.Second)
	m.tick()
	if site.beatCount() != beats+1 || len(site.lastBeat().Credentials) != 0 {
		t.Fatalf("beat = %+v", site.lastBeat())
	}
	if st := m.statusFor(loadPublishServer(t, s.ID)); st.State != models.LinkError || st.ErrorCode != "node_error" {
		t.Fatalf("status = %+v", st)
	}
}

func TestLinkUnregisterOnDisable(t *testing.T) {
	s, site, dock, rt, now, m := setupLinkTest(t)
	database.DB.Model(&s).Update("status", "running")
	*rt = linkRuntime{Ready: true}
	m.tick()
	name := utils.GetLinkContainerName(s.ID)
	if dock.container(name) == nil {
		t.Fatal("node not running")
	}

	// 关闭联机房：带令牌下架，删除节点和凭证目录，清空令牌，之后不再请求
	database.DB.Model(&s).Update("link", false)
	*now = now.Add(15 * time.Second)
	m.tick()
	m.tick()
	expectLinkCalls(t, site, "POST /link/rooms", "POST /link/rooms/12/heartbeat", "DELETE /link/rooms/12")
	if site.tokens[2] != "alr_tok" {
		t.Fatalf("delete token = %q", site.tokens[2])
	}
	if dock.container(name) != nil || !slices.Contains(dock.removedDirs, s.ID) {
		t.Fatal("node and data dir should be removed")
	}
	got := loadPublishServer(t, s.ID)
	if got.LinkToken != "" || got.LinkRoomID != 0 || got.LinkNetworkName != "" || got.LinkNetworkSecret != "" {
		t.Fatalf("link fields not cleared: %+v", got)
	}
	if m.statusFor(got) != nil {
		t.Fatal("status should be nil when link is off")
	}
}

func TestLinkUnregisterDeleted(t *testing.T) {
	s, site, dock, _, now, m := setupLinkTest(t)
	m.tick()
	// 残留的节点容器（例如服务器被删除时没删掉）也会被清理
	dock.Create(docker_manager.LinkContainerSpec{ServerID: s.ID})
	database.DB.Delete(&models.Server{}, s.ID)

	// 网站暂时不可用：退避，期间不重试
	site.set(func(f *fakeLinkSite) {
		f.fail = func(w http.ResponseWriter) {
			w.WriteHeader(http.StatusServiceUnavailable)
			json.NewEncoder(w).Encode(map[string]string{"code": "busy", "error": "busy"})
		}
	})
	*now = now.Add(15 * time.Second)
	m.tick()
	if dock.container(utils.GetLinkContainerName(s.ID)) != nil {
		t.Fatal("orphan node should be removed")
	}
	*now = now.Add(30 * time.Second)
	m.tick()
	expectLinkCalls(t, site, "POST /link/rooms", "POST /link/rooms/12/heartbeat", "DELETE /link/rooms/12")

	// 网站已经没有这个房间（invalid_token）：视同下架成功
	site.set(func(f *fakeLinkSite) {
		f.fail = func(w http.ResponseWriter) {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"code": "invalid_token", "error": "invalid token"})
		}
	})
	*now = now.Add(time.Minute)
	m.tick()
	m.tick()
	expectLinkCalls(t, site, "POST /link/rooms", "POST /link/rooms/12/heartbeat", "DELETE /link/rooms/12", "DELETE /link/rooms/12")
	if got := loadPublishServer(t, s.ID); got.LinkToken != "" {
		t.Fatalf("token should be cleared: %q", got.LinkToken)
	}
}

func TestLinkBackoff(t *testing.T) {
	s, site, _, _, now, m := setupLinkTest(t)
	site.set(func(f *fakeLinkSite) {
		f.fail = func(w http.ResponseWriter) {
			w.WriteHeader(http.StatusServiceUnavailable)
			json.NewEncoder(w).Encode(map[string]string{"code": "busy", "error": "busy"})
		}
	})

	// 失败后按 1、2、4 分钟退避，期间不重试
	expect := 0
	for _, wait := range []time.Duration{time.Minute, 2 * time.Minute, 4 * time.Minute} {
		m.tick()
		expect++
		*now = now.Add(wait - time.Second)
		m.tick()
		if got := len(site.callList()); got != expect {
			t.Fatalf("retried before backoff %v: %d calls", wait, got)
		}
		*now = now.Add(time.Second)
	}
	st := m.statusFor(loadPublishServer(t, s.ID))
	if st.State != models.LinkError || st.ErrorCode != "busy" || st.RetryAt == "" {
		t.Fatalf("status = %+v", st)
	}

	// 429 时按 Retry-After 等待（比退避时间长）
	site.set(func(f *fakeLinkSite) {
		f.fail = func(w http.ResponseWriter) {
			w.Header().Set("Retry-After", "1800")
			w.WriteHeader(http.StatusTooManyRequests)
			json.NewEncoder(w).Encode(map[string]string{"code": "rate_limited", "error": "slow down"})
		}
	})
	*now = now.Add(time.Hour)
	m.tick()
	calls := len(site.callList())
	*now = now.Add(29 * time.Minute)
	m.tick()
	if len(site.callList()) != calls {
		t.Fatal("should wait for Retry-After")
	}

	// 恢复后登记成功，错误清除
	site.set(func(f *fakeLinkSite) { f.fail = nil })
	*now = now.Add(2 * time.Minute)
	m.tick()
	if got := loadPublishServer(t, s.ID); got.LinkToken != "alr_tok" {
		t.Fatalf("token = %q", got.LinkToken)
	}
	if st := m.statusFor(loadPublishServer(t, s.ID)); st.State != models.LinkOffline || st.ErrorCode != "" {
		t.Fatalf("status = %+v", st)
	}
}

func TestLinkHeartbeatInvalidTokenReregisters(t *testing.T) {
	s, site, dock, _, now, m := setupLinkTest(t)
	m.tick()
	site.set(func(f *fakeLinkSite) {
		f.fail = func(w http.ResponseWriter) {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"code": "invalid_token", "error": "invalid token"})
		}
	})
	*now = now.Add(time.Minute)
	m.tick()
	if got := loadPublishServer(t, s.ID); got.LinkToken != "" || got.LinkRoomID != 0 {
		t.Fatalf("token should be cleared: %+v", got)
	}
	if !slices.Contains(dock.removedDirs, s.ID) {
		t.Fatal("credential dir should be reset")
	}
	// 退避结束后重新登记（不带令牌）
	site.set(func(f *fakeLinkSite) { f.fail = nil })
	*now = now.Add(time.Minute)
	m.tick()
	calls := site.callList()
	if calls[len(calls)-2] != "POST /link/rooms" || site.tokens[len(calls)-2] != "" {
		t.Fatalf("calls = %v tokens = %v", calls, site.tokens)
	}
}

func TestLinkMembersProxy(t *testing.T) {
	s, site, _, _, _, m := setupLinkTest(t)
	old := linkManager
	linkManager = m
	t.Cleanup(func() { linkManager = old })
	svc := NewServerService()
	id := fmt.Sprint(s.ID)

	// 还没登记：409
	if _, err := svc.GetLinkMembers(1, id); !errors.Is(err, ErrLinkNotRegistered) {
		t.Fatalf("err = %v", err)
	}
	m.tick()
	members, err := svc.GetLinkMembers(1, id)
	if err != nil || len(members.Members) != 1 || members.Members[0].SteamID != "76561198000000001" || members.Bans == nil {
		t.Fatalf("members = %+v err = %v", members, err)
	}
	if err := svc.KickLinkMember(1, id, models.LinkKickRequest{SteamID: "bad"}); !errors.Is(err, ErrInvalidSteamID) {
		t.Fatalf("err = %v", err)
	}
	if err := svc.KickLinkMember(1, id, models.LinkKickRequest{SteamID: "76561198000000001", Ban: true}); err != nil {
		t.Fatal(err)
	}
	// 网站没有这个接口时按网站拒绝处理
	if err := svc.UnbanLinkMember(1, id, "76561198000000001"); !errors.Is(err, ErrLinkRequestFailed) {
		t.Fatalf("err = %v", err)
	}
	// 别人的服务器
	if _, err := svc.GetLinkMembers(2, id); !errors.Is(err, ErrServerNotFound) {
		t.Fatalf("err = %v", err)
	}
	if !slices.Contains(site.callList(), "POST /link/rooms/12/kick") {
		t.Fatalf("calls = %v", site.callList())
	}
}

func TestLinkNodeArgsAndRPCPort(t *testing.T) {
	setupTestDB(t)
	a := createTestServer(t, 7777, 27015, 32330)
	// 另一台服务器的游戏端口正好占了 20000 + a.ID
	createTestServer(t, linkRPCPortBase+int(a.ID), 27017, 32332)
	port, err := assignLinkRPCPort(a.ID)
	if err != nil {
		t.Fatal(err)
	}
	if port != linkRPCPortBase+int(a.ID)+2 { // 顺延跳过游戏端口和 +1
		t.Fatalf("port = %d", port)
	}

	s := models.Server{ID: 3, SessionName: "", Port: 7777, QueryPort: 27015, LinkNetworkName: "ark-x", LinkNetworkSecret: "sec", LinkRPCPort: 20003}
	spec1 := linkContainerSpec(s, []string{"tcp://a:1"})
	if !slices.Contains(spec1.Args, "ark-3") || spec1.Image != utils.EasyTierImage {
		t.Fatalf("spec = %+v", spec1)
	}
	s.Port = 7779
	spec2 := linkContainerSpec(s, []string{"tcp://a:1"})
	if spec1.Labels[linkCommandHashLabel] == spec2.Labels[linkCommandHashLabel] {
		t.Fatal("hash should change with args")
	}
	if id := newUUID(); len(id) != 36 || id[14] != '4' {
		t.Fatalf("uuid = %s", id)
	}
	rt := serverVersionRe.FindStringSubmatch("My Server - (v358.24)")
	if rt == nil || rt[1] != "358.24" {
		t.Fatalf("version = %v", rt)
	}
}
