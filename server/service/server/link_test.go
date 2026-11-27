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
	mu                sync.Mutex
	imageReady        bool
	imageID           string
	containers        map[string]*fakeLinkContainer
	creates           int
	removed           []string
	removedDirs       int
	removedLegacyDirs []uint
	generated         []string
	revoked           []string
	execErr           error
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
	d.containers[utils.LinkContainerName] = &fakeLinkContainer{running: true, imageID: d.imageID, labels: spec.Labels, args: spec.Args}
	return nil
}

// addContainer 放一个残留的容器（例如旧版的 ase-link-<id>）
func (d *fakeLinkDocker) addContainer(name string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.containers[name] = &fakeLinkContainer{running: true}
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

func (d *fakeLinkDocker) RemoveDir() error {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.removedDirs++
	return nil
}

func (d *fakeLinkDocker) RemoveLegacyDir(id uint) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.removedLegacyDirs = append(d.removedLegacyDirs, id)
	return nil
}

func (d *fakeLinkDocker) container(name string) *fakeLinkContainer {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.containers[name]
}

func (d *fakeLinkDocker) createCount() int {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.creates
}

// fakeLinkSite 模拟网站的联机房接口：登记、心跳、成员管理都要账号令牌，下架还接受旧版房间令牌
type fakeLinkSite struct {
	mu       sync.Mutex
	calls    []string // "METHOD path"
	tokens   []string
	beats    []linkBeat
	profiles []linkRoomProfile
	legacy   []string // 每次登记带的 legacy_token
	want     int
	revoke   []string
	roomID   int // 登记返回的房间 ID，0 为 12
	// fail 不为空时所有请求按它返回
	fail func(w http.ResponseWriter)
	// beatFail 不为空时心跳按它返回
	beatFail func(w http.ResponseWriter)
}

type linkBeat struct {
	Online      bool             `json:"online"`
	Servers     []linkServerBeat `json:"servers"`
	Credentials []linkCredential `json:"credentials"`
	RevokedAck  []string         `json:"revoked_ack"`
}

// server 心跳里某台服务器的状态
func (b linkBeat) server(key string) (linkServerBeat, bool) {
	for _, s := range b.Servers {
		if s.Key == key {
			return s, true
		}
	}
	return linkServerBeat{}, false
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
	roomID := f.roomID
	if roomID == 0 {
		roomID = 12
	}
	room := fmt.Sprintf("/link/rooms/%d", roomID)
	// 旧版房间令牌只能下架
	if strings.HasPrefix(token, "alr_") && r.Method == http.MethodDelete {
		json.NewEncoder(w).Encode(map[string]any{"ok": true})
		return
	}
	if token != testAccountToken {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"code": "login_required", "error": "login required"})
		return
	}
	switch {
	case r.Method == http.MethodPost && r.URL.Path == "/link/rooms":
		var p linkRegisterRequest
		_ = json.NewDecoder(r.Body).Decode(&p)
		f.profiles = append(f.profiles, p.linkRoomProfile)
		f.legacy = append(f.legacy, p.LegacyToken)
		json.NewEncoder(w).Encode(map[string]any{"id": roomID, "network_name": "ark-0123456789abcdef", "peers": peers,
			"url": fmt.Sprintf("https://arkfgf.com/rooms/%d", roomID)})
	case r.Method == http.MethodPost && r.URL.Path == room+"/heartbeat" && f.beatFail != nil:
		f.beatFail(w)
	case r.Method == http.MethodPost && r.URL.Path == room+"/heartbeat":
		var b linkBeat
		_ = json.NewDecoder(r.Body).Decode(&b)
		f.beats = append(f.beats, b)
		// 网站收到凭证后减少缺口，收到确认后从 revoke 里去掉
		f.want = max(f.want-len(b.Credentials), 0)
		f.revoke = slices.DeleteFunc(f.revoke, func(id string) bool { return slices.Contains(b.RevokedAck, id) })
		json.NewEncoder(w).Encode(map[string]any{"want_credentials": f.want, "revoke": f.revoke, "network_name": "ark-0123456789abcdef", "peers": peers})
	case r.Method == http.MethodDelete && r.URL.Path == room:
		json.NewEncoder(w).Encode(map[string]any{"ok": true})
	case r.Method == http.MethodGet && r.URL.Path == room+"/members":
		json.NewEncoder(w).Encode(map[string]any{"members": []map[string]any{{"steam_id": "76561198000000001", "name": "p1", "credential_id": "c1", "joined_at": 1797600000}}, "bans": []any{}})
	case r.Method == http.MethodPost && r.URL.Path == room+"/kick":
		json.NewEncoder(w).Encode(map[string]any{"ok": true})
	case strings.HasPrefix(r.URL.Path, "/link/rooms/") && !strings.HasPrefix(r.URL.Path, room+"/bans/"):
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"code": "room_not_found", "error": "room not found"})
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

func (f *fakeLinkSite) lastProfile() linkRoomProfile {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.profiles[len(f.profiles)-1]
}

func (f *fakeLinkSite) beatCount() int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return len(f.beats)
}

func (f *fakeLinkSite) count(call string) int {
	f.mu.Lock()
	defer f.mu.Unlock()
	n := 0
	for _, c := range f.calls {
		if c == call {
			n++
		}
	}
	return n
}

func (f *fakeLinkSite) set(fn func(f *fakeLinkSite)) {
	f.mu.Lock()
	defer f.mu.Unlock()
	fn(f)
}

// fakeRuntimes 按服务器 ID 返回运行信息
type fakeRuntimes struct {
	mu sync.Mutex
	rt map[uint]linkRuntime
}

func (f *fakeRuntimes) set(id uint, rt linkRuntime) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.rt[id] = rt
}

func (f *fakeRuntimes) get(s models.Server) linkRuntime {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.rt[s.ID]
}

func newTestLinkManager(t *testing.T, site *fakeLinkSite, dock *fakeLinkDocker, rt *fakeRuntimes, now *time.Time) *LinkManager {
	t.Helper()
	utils.InitLogger()
	srv := httptest.NewServer(http.HandlerFunc(site.handler))
	t.Cleanup(srv.Close)
	n := 0
	return &LinkManager{
		room:    &linkState{},
		servers: map[uint]*linkServerState{},
		legacy:  map[uint]*linkBackoff{},
		client:  srv.Client(),
		api:     func() string { return srv.URL },
		now:     func() time.Time { return *now },
		docker:  dock,
		runtime: rt.get,
		profile: func(s models.Server) linkServerProfile {
			return linkServerProfile{Name: s.SessionName, Map: s.Map, Mods: []string{}, GamePort: s.Port, QueryPort: s.QueryPort}
		},
		newID: func() string { n++; return fmt.Sprintf("00000000-0000-4000-8000-%012d", n) },
	}
}

type linkTestEnv struct {
	s    models.Server
	site *fakeLinkSite
	dock *fakeLinkDocker
	rt   *fakeRuntimes
	now  *time.Time
	m    *LinkManager
}

func (e *linkTestEnv) advance(d time.Duration) { *e.now = e.now.Add(d) }

func setupLinkTest(t *testing.T) *linkTestEnv {
	setupTestDB(t)
	bindTestAccount(t)
	s := createTestServer(t, 7777, 27015, 32330)
	database.DB.Model(&s).Update("link", true)
	s.Link = true
	site := &fakeLinkSite{}
	dock := newFakeLinkDocker()
	rt := &fakeRuntimes{rt: map[uint]linkRuntime{}}
	now := time.Now()
	m := newTestLinkManager(t, site, dock, rt, &now)
	return &linkTestEnv{s: s, site: site, dock: dock, rt: rt, now: &now, m: m}
}

func loadTestLinkRoom(t *testing.T) models.LinkRoom {
	t.Helper()
	room, err := loadLinkRoom()
	if err != nil {
		t.Fatal(err)
	}
	return room
}

func expectLinkCalls(t *testing.T, site *fakeLinkSite, want ...string) {
	t.Helper()
	if got := site.callList(); strings.Join(got, ",") != strings.Join(want, ",") {
		t.Fatalf("calls = %v\nwant    %v", got, want)
	}
}

func TestLinkRegisterHeartbeatAndNode(t *testing.T) {
	e := setupLinkTest(t)
	s, site, dock, m := e.s, e.site, e.dock, e.m
	name := utils.LinkContainerName

	// 已停止：用账号令牌登记房间并保存房间 ID 等，发一次离线心跳，不创建节点
	m.tick()
	expectLinkCalls(t, site, "POST /link/rooms", "POST /link/rooms/12/heartbeat")
	if site.tokens[0] != testAccountToken || site.tokens[1] != testAccountToken || site.legacy[0] != "" {
		t.Fatalf("tokens = %v legacy = %v", site.tokens, site.legacy)
	}
	room := loadTestLinkRoom(t)
	if room.Token != "" || room.RoomID != 12 || room.NetworkName != "ark-0123456789abcdef" {
		t.Fatalf("room not saved: %+v", room)
	}
	if len(room.NetworkSecret) != 64 || room.RPCPort != linkRPCPortBase {
		t.Fatalf("local params: secret %q port %d", room.NetworkSecret, room.RPCPort)
	}
	p := site.lastProfile()
	if p.Name != "test" || p.ServerIP != linkVirtualIP || len(p.Servers) != 1 || p.Servers[0].Key != fmt.Sprint(s.ID) ||
		p.Servers[0].GamePort != 7777 || p.Servers[0].QueryPort != 27015 {
		t.Fatalf("profile = %+v", p)
	}
	b := site.lastBeat()
	if sb, ok := b.server(fmt.Sprint(s.ID)); b.Online || !ok || sb.Online || sb.Players != 0 || len(b.Credentials) != 0 {
		t.Fatalf("offline beat = %+v", b)
	}
	if st := m.statusFor(s); st.State != models.LinkOffline || st.RoomURL != "https://arkfgf.com/rooms/12" {
		t.Fatalf("status = %+v", st)
	}
	if dock.createCount() != 0 {
		t.Fatal("node should not run while stopped")
	}

	// 离线时每 60 秒心跳一次
	e.advance(30 * time.Second)
	m.tick()
	if site.beatCount() != 1 {
		t.Fatal("offline heartbeat before 60s")
	}
	e.advance(30 * time.Second)
	m.tick()
	if site.beatCount() != 2 {
		t.Fatal("offline heartbeat missing at 60s")
	}

	// 运行中但还没 ready：不创建节点，状态为等待
	database.DB.Model(&s).Update("status", "running")
	s.Status = "running"
	e.advance(15 * time.Second)
	m.tick()
	if dock.createCount() != 0 {
		t.Fatal("node should wait until ready")
	}
	if st := m.statusFor(s); st.State != models.LinkWaiting {
		t.Fatalf("status = %+v", st)
	}

	// ready：创建节点，立即发在线心跳
	e.rt.set(s.ID, linkRuntime{Ready: true, Players: 3})
	e.advance(15 * time.Second)
	m.tick()
	c := dock.container(name)
	if c == nil || dock.createCount() != 1 {
		t.Fatalf("node not created: %d", dock.createCount())
	}
	args := strings.Join(c.args, " ")
	for _, want := range []string{
		"--network-name ark-0123456789abcdef", "--network-secret " + room.NetworkSecret, "--secure-mode",
		"--credential-file /data/cred.json", "--peers tcp://et.arkfgf.com:11010 udp://et.arkfgf.com:11010 --ipv4",
		"--ipv4 10.126.126.1 --no-tun --p2p-only", "--tcp-whitelist 0 --udp-whitelist 7777 7778 27015 --listeners",
		"--listeners tcp://0.0.0.0:0 udp://0.0.0.0:0", fmt.Sprintf("--rpc-portal 127.0.0.1:%d", room.RPCPort), "--hostname test",
	} {
		if !strings.Contains(args, want) {
			t.Fatalf("args missing %q: %s", want, args)
		}
	}
	b = site.lastBeat()
	if sb, _ := b.server(fmt.Sprint(s.ID)); !b.Online || !sb.Online || sb.Players != 3 {
		t.Fatalf("online beat = %+v", b)
	}
	if st := m.statusFor(s); st.State != models.LinkOnline {
		t.Fatalf("status = %+v", st)
	}

	// 在线时每 30 秒心跳一次，节点不重复创建
	beats := site.beatCount()
	e.advance(15 * time.Second)
	m.tick()
	if site.beatCount() != beats {
		t.Fatal("online heartbeat before 30s")
	}
	e.advance(15 * time.Second)
	m.tick()
	if site.beatCount() != beats+1 || dock.createCount() != 1 {
		t.Fatalf("beats %d creates %d", site.beatCount(), dock.createCount())
	}

	// 短暂查询不到：节点保留，宽限期内仍报告在线
	e.rt.set(s.ID, linkRuntime{})
	e.advance(15 * time.Second)
	m.tick()
	if dock.container(name) == nil {
		t.Fatal("node removed on a short A2S miss")
	}
	e.rt.set(s.ID, linkRuntime{Ready: true, Players: 2})

	// 资料变化（改名）：带令牌重新登记；房间名称跟着第一台服务器变，节点参数（hostname）变化，重建容器
	database.DB.Model(&s).Update("session_name", "renamed")
	e.advance(15 * time.Second)
	m.tick()
	if site.count("POST /link/rooms") != 2 || site.tokens[slices.Index(site.callList()[1:], "POST /link/rooms")+1] != testAccountToken {
		t.Fatalf("expected a re-registration with the account token: %v %v", site.callList(), site.tokens)
	}
	if p := site.lastProfile(); p.Name != "renamed" || p.Servers[0].Name != "renamed" {
		t.Fatalf("profile = %+v", p)
	}
	if dock.createCount() != 2 || !slices.Contains(dock.container(name).args, "renamed") {
		t.Fatalf("node not recreated after args change: %d", dock.createCount())
	}

	// 镜像更新（镜像 ID 变化）也重建
	dock.mu.Lock()
	dock.imageID = "sha256:et2"
	dock.mu.Unlock()
	e.advance(15 * time.Second)
	m.tick()
	if dock.createCount() != 3 {
		t.Fatalf("node not recreated after image update: %d", dock.createCount())
	}

	// 停服：删除节点，立即发离线心跳
	database.DB.Model(&s).Update("status", "stopped")
	s.Status = "stopped"
	e.rt.set(s.ID, linkRuntime{})
	beats = site.beatCount()
	e.advance(15 * time.Second)
	m.tick()
	if dock.container(name) != nil {
		t.Fatal("node should be removed when the server stops")
	}
	if site.beatCount() != beats+1 || site.lastBeat().Online {
		t.Fatalf("expected offline beat: %+v", site.lastBeat())
	}
	if st := m.statusFor(s); st.State != models.LinkOffline {
		t.Fatalf("status = %+v", st)
	}
}

// 一个房间里多台服务器：一个节点放行所有服务器的端口；单台启停不重建节点；按服务器心跳
func TestLinkMultipleServers(t *testing.T) {
	e := setupLinkTest(t)
	a, site, dock, m := e.s, e.site, e.dock, e.m
	b := createTestServer(t, 7779, 27017, 32332)
	database.DB.Model(&b).Updates(map[string]any{"link": true, "session_name": "second"})
	b.Link, b.SessionName = true, "second"
	// 没开联机房的服务器不在房间里
	createTestServer(t, 7781, 27019, 32334)
	keyA, keyB := fmt.Sprint(a.ID), fmt.Sprint(b.ID)

	// a 运行中可以连接，b 停止
	database.DB.Model(&a).Update("status", "running")
	a.Status = "running"
	e.rt.set(a.ID, linkRuntime{Ready: true, Players: 4, Version: "358.24"})
	m.tick()

	if site.count("POST /link/rooms") != 1 {
		t.Fatalf("calls = %v", site.callList())
	}
	p := site.lastProfile()
	if p.Name != "test" || len(p.Servers) != 2 || p.Servers[0].Key != keyA || p.Servers[1].Key != keyB ||
		p.Servers[1].Name != "second" || p.Servers[1].GamePort != 7779 || p.Servers[0].Version != "358.24" {
		t.Fatalf("profile = %+v", p)
	}
	c := dock.container(utils.LinkContainerName)
	if c == nil || dock.createCount() != 1 {
		t.Fatal("node not created")
	}
	// 已停止的 b 的端口也放行，b 启动时不用重建节点
	if args := strings.Join(c.args, " "); !strings.Contains(args, "--udp-whitelist 7777 7778 27015 7779 7780 27017 --listeners") {
		t.Fatalf("args = %s", args)
	}
	beat := site.lastBeat()
	sa, _ := beat.server(keyA)
	sb, okB := beat.server(keyB)
	if !beat.Online || !sa.Online || sa.Players != 4 || !okB || sb.Online || len(beat.Servers) != 2 {
		t.Fatalf("beat = %+v", beat)
	}
	if st := m.statusFor(a); st.State != models.LinkOnline {
		t.Fatalf("a status = %+v", st)
	}
	if st := m.statusFor(b); st.State != models.LinkOffline || st.RoomURL == "" {
		t.Fatalf("b status = %+v", st)
	}

	// b 启动并可以连接：节点不重建，立即心跳报告 b 在线
	beats := site.beatCount()
	database.DB.Model(&b).Update("status", "running")
	b.Status = "running"
	e.rt.set(b.ID, linkRuntime{Ready: true, Players: 1})
	e.advance(15 * time.Second)
	m.tick()
	if dock.createCount() != 1 {
		t.Fatalf("starting another server should not recreate the node: %d", dock.createCount())
	}
	if site.beatCount() != beats+1 {
		t.Fatal("server state change should trigger a heartbeat")
	}
	if sb, _ := site.lastBeat().server(keyB); !sb.Online || sb.Players != 1 {
		t.Fatalf("beat = %+v", site.lastBeat())
	}
	if st := m.statusFor(b); st.State != models.LinkOnline {
		t.Fatalf("b status = %+v", st)
	}

	// a 停止：节点保留（b 还在运行），a 报告离线
	database.DB.Model(&a).Update("status", "stopped")
	a.Status = "stopped"
	e.rt.set(a.ID, linkRuntime{})
	e.advance(15 * time.Second)
	m.tick()
	if dock.container(utils.LinkContainerName) == nil || dock.createCount() != 1 {
		t.Fatal("node should stay while another server is running")
	}
	beat = site.lastBeat()
	if sa, _ := beat.server(keyA); !beat.Online || sa.Online {
		t.Fatalf("beat = %+v", beat)
	}

	// 关闭 b 的联机房：重新登记只剩 a，节点按新端口重建；房间不下架
	database.DB.Model(&b).Update("link", false)
	e.advance(15 * time.Second)
	m.tick()
	if p := site.lastProfile(); len(p.Servers) != 1 || p.Servers[0].Key != keyA {
		t.Fatalf("profile = %+v", p)
	}
	if site.count("DELETE /link/rooms/12") != 0 {
		t.Fatal("room should not be removed while other servers are linked")
	}
	// 剩下的 a 没在运行，节点删除
	if dock.container(utils.LinkContainerName) != nil {
		t.Fatal("node should be removed when no linked server is running")
	}
}

func TestLinkCredentialsAndRevoke(t *testing.T) {
	e := setupLinkTest(t)
	s, site, dock, m := e.s, e.site, e.dock, e.m
	database.DB.Model(&s).Update("status", "running")
	s.Status = "running"
	e.rt.set(s.ID, linkRuntime{Ready: true})
	site.set(func(f *fakeLinkSite) { f.want = 25 })

	// 第一轮：登记、建节点、在线心跳，网站要 25 个凭证
	m.tick()
	if site.beatCount() != 1 || len(site.lastBeat().Credentials) != 0 {
		t.Fatalf("first beat = %+v", site.lastBeat())
	}

	// 下一轮签发 20 个（每轮上限），随心跳上传
	e.advance(15 * time.Second)
	m.tick()
	b := site.lastBeat()
	if len(b.Credentials) != 20 || len(dock.generated) != 20 {
		t.Fatalf("credentials = %d generated = %d", len(b.Credentials), len(dock.generated))
	}
	c := b.Credentials[0]
	if c.ID != dock.generated[0] || c.Secret != "secret-"+c.ID {
		t.Fatalf("credential = %+v", c)
	}
	if want := e.now.Add(30 * 24 * time.Hour).Unix(); c.ExpiresAt != want {
		t.Fatalf("expires_at = %d, want %d", c.ExpiresAt, want)
	}

	// 再下一轮补剩下的 5 个
	e.advance(15 * time.Second)
	m.tick()
	if b := site.lastBeat(); len(b.Credentials) != 5 || len(dock.generated) != 25 {
		t.Fatalf("credentials = %d generated = %d", len(b.Credentials), len(dock.generated))
	}

	// 池满后不再提前心跳
	beats := site.beatCount()
	e.advance(15 * time.Second)
	m.tick()
	if site.beatCount() != beats {
		t.Fatal("no extra heartbeat expected when the pool is full")
	}

	// 踢人：网站下次心跳返回 revoke；之后在节点里吊销，再下一次心跳回报 revoked_ack
	site.set(func(f *fakeLinkSite) { f.revoke = []string{"cred-x"} })
	if err := m.linkCall(http.MethodPost, "/kick", models.LinkKickRequest{SteamID: "76561198000000001"}, nil); err != nil {
		t.Fatal(err)
	}
	m.poke()
	e.advance(15 * time.Second)
	m.tick() // 提前心跳，拿到 revoke
	e.advance(15 * time.Second)
	m.tick() // 吊销并回报
	if !slices.Equal(dock.revoked, []string{"cred-x"}) {
		t.Fatalf("revoked = %v", dock.revoked)
	}
	if b := site.lastBeat(); !slices.Equal(b.RevokedAck, []string{"cred-x"}) {
		t.Fatalf("revoked_ack = %v", b.RevokedAck)
	}
	// 确认后不再重复吊销和回报
	e.advance(30 * time.Second)
	m.tick()
	if len(dock.revoked) != 1 || len(site.lastBeat().RevokedAck) != 0 {
		t.Fatalf("revoked = %v ack = %v", dock.revoked, site.lastBeat().RevokedAck)
	}

	// 节点出错：状态显示错误，心跳照发（不上传凭证）
	site.set(func(f *fakeLinkSite) { f.want = 3 })
	e.advance(30 * time.Second)
	m.tick()
	dock.mu.Lock()
	dock.execErr = errors.New("rpc down")
	dock.mu.Unlock()
	beats = site.beatCount()
	e.advance(15 * time.Second)
	m.tick()
	if site.beatCount() != beats+1 || len(site.lastBeat().Credentials) != 0 {
		t.Fatalf("beat = %+v", site.lastBeat())
	}
	if st := m.statusFor(s); st.State != models.LinkError || st.ErrorCode != "node_error" {
		t.Fatalf("status = %+v", st)
	}
}

func TestLinkUnregisterOnDisable(t *testing.T) {
	e := setupLinkTest(t)
	s, site, dock, m := e.s, e.site, e.dock, e.m
	database.DB.Model(&s).Update("status", "running")
	e.rt.set(s.ID, linkRuntime{Ready: true})
	m.tick()
	name := utils.LinkContainerName
	if dock.container(name) == nil {
		t.Fatal("node not running")
	}

	// 最后一台服务器关闭联机房：用账号令牌下架，删除节点和凭证目录，清空房间 ID 等，之后不再请求
	database.DB.Model(&models.LinkRoom{}).Where("id = ?", linkRoomID).Update("name", "My Room")
	database.DB.Model(&s).Update("link", false)
	e.advance(15 * time.Second)
	m.tick()
	m.tick()
	expectLinkCalls(t, site, "POST /link/rooms", "POST /link/rooms/12/heartbeat", "DELETE /link/rooms/12")
	if site.tokens[2] != testAccountToken {
		t.Fatalf("delete token = %q", site.tokens[2])
	}
	if dock.container(name) != nil || dock.removedDirs == 0 {
		t.Fatal("node and data dir should be removed")
	}
	room := loadTestLinkRoom(t)
	if room.Token != "" || room.RoomID != 0 || room.NetworkName != "" || room.NetworkSecret != "" {
		t.Fatalf("link fields not cleared: %+v", room)
	}
	// 房间资料保留
	if room.Name != "My Room" {
		t.Fatalf("room settings should be kept: %+v", room)
	}
	s.Link = false
	if m.statusFor(s) != nil {
		t.Fatal("status should be nil when link is off")
	}
}

func TestLinkUnregisterDeleted(t *testing.T) {
	e := setupLinkTest(t)
	s, site, m := e.s, e.site, e.m
	m.tick()
	database.DB.Delete(&models.Server{}, s.ID)

	// 网站暂时不可用：退避，期间不重试
	site.set(func(f *fakeLinkSite) {
		f.fail = func(w http.ResponseWriter) {
			w.WriteHeader(http.StatusServiceUnavailable)
			json.NewEncoder(w).Encode(map[string]string{"code": "busy", "error": "busy"})
		}
	})
	e.advance(15 * time.Second)
	m.tick()
	e.advance(30 * time.Second)
	m.tick()
	expectLinkCalls(t, site, "POST /link/rooms", "POST /link/rooms/12/heartbeat", "DELETE /link/rooms/12")

	// 网站已经没有这个房间（room_not_found）：视同下架成功
	site.set(func(f *fakeLinkSite) {
		f.fail = func(w http.ResponseWriter) {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"code": "room_not_found", "error": "room not found"})
		}
	})
	e.advance(time.Minute)
	m.tick()
	m.tick()
	expectLinkCalls(t, site, "POST /link/rooms", "POST /link/rooms/12/heartbeat", "DELETE /link/rooms/12", "DELETE /link/rooms/12")
	if room := loadTestLinkRoom(t); room.RoomID != 0 || room.NetworkSecret != "" {
		t.Fatalf("room should be cleared: %+v", room)
	}
}

// 升级后下架旧版每台服务器一个房间时登记的房间，删除旧的节点容器和凭证目录
func TestLinkLegacyCleanup(t *testing.T) {
	e := setupLinkTest(t)
	s, site, dock, m := e.s, e.site, e.dock, e.m
	database.DB.Model(&s).Updates(map[string]any{"link_token": "alr_old", "link_room_id": 7, "link_network_name": "ark-old",
		"link_network_secret": "sec", "link_rpc_port": 20001})
	// 已删除的服务器的旧房间也要下架
	gone := createTestServer(t, 7781, 27019, 32334)
	database.DB.Model(&gone).Updates(map[string]any{"link_token": "alr_gone", "link_room_id": 9})
	database.DB.Delete(&models.Server{}, gone.ID)
	dock.addContainer(fmt.Sprintf("ase-link-%d", s.ID))

	m.tick()
	calls := site.callList()
	if !slices.Contains(calls, "DELETE /link/rooms/7") || !slices.Contains(calls, "DELETE /link/rooms/9") {
		t.Fatalf("calls = %v", calls)
	}
	if i := slices.Index(calls, "DELETE /link/rooms/7"); site.tokens[i] != "alr_old" {
		t.Fatalf("legacy token = %q", site.tokens[i])
	}
	// 网站没有 9 号房间（404）也算下架完成
	var got []models.Server
	database.DB.Unscoped().Where("link_token <> ''").Find(&got)
	if len(got) != 0 {
		t.Fatalf("legacy tokens not cleared: %+v", got)
	}
	var reloaded models.Server
	database.DB.First(&reloaded, s.ID)
	if reloaded.LinkRoomID != 0 || reloaded.LinkNetworkSecret != "" || reloaded.LinkRPCPort != 0 || !reloaded.Link {
		t.Fatalf("legacy fields = %+v", reloaded)
	}
	if !slices.Contains(dock.removedLegacyDirs, s.ID) || !slices.Contains(dock.removedLegacyDirs, gone.ID) {
		t.Fatalf("legacy dirs = %v", dock.removedLegacyDirs)
	}
	if dock.container(fmt.Sprintf("ase-link-%d", s.ID)) != nil {
		t.Fatal("legacy node should be removed")
	}
	// 新房间照常登记
	if room := loadTestLinkRoom(t); room.RoomID != 12 || room.Token != "" {
		t.Fatalf("room = %+v", room)
	}
	// 之后不再请求旧房间
	e.advance(time.Minute)
	m.tick()
	if site.count("DELETE /link/rooms/7") != 1 {
		t.Fatalf("calls = %v", site.callList())
	}
}

func TestLinkBackoff(t *testing.T) {
	e := setupLinkTest(t)
	s, site, m := e.s, e.site, e.m
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
		e.advance(wait - time.Second)
		m.tick()
		if got := len(site.callList()); got != expect {
			t.Fatalf("retried before backoff %v: %d calls", wait, got)
		}
		e.advance(time.Second)
	}
	st := m.statusFor(s)
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
	e.advance(time.Hour)
	m.tick()
	calls := len(site.callList())
	e.advance(29 * time.Minute)
	m.tick()
	if len(site.callList()) != calls {
		t.Fatal("should wait for Retry-After")
	}

	// 恢复后登记成功，错误清除
	site.set(func(f *fakeLinkSite) { f.fail = nil })
	e.advance(2 * time.Minute)
	m.tick()
	if room := loadTestLinkRoom(t); room.RoomID != 12 {
		t.Fatalf("room = %+v", room)
	}
	if st := m.statusFor(s); st.State != models.LinkOffline || st.ErrorCode != "" {
		t.Fatalf("status = %+v", st)
	}
}

// 心跳返回 room_not_found：清掉房间 ID 和网络名、删除节点和凭证，下一轮重新登记
func TestLinkHeartbeatRoomNotFoundReregisters(t *testing.T) {
	e := setupLinkTest(t)
	s, site, dock, m := e.s, e.site, e.dock, e.m
	database.DB.Model(&s).Update("status", "running")
	e.rt.set(s.ID, linkRuntime{Ready: true})
	m.tick()
	if dock.container(utils.LinkContainerName) == nil {
		t.Fatal("node not running")
	}
	site.set(func(f *fakeLinkSite) {
		f.beatFail = func(w http.ResponseWriter) {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"code": "room_not_found", "error": "room not found"})
		}
	})
	e.advance(time.Minute)
	m.tick()
	if room := loadTestLinkRoom(t); room.RoomID != 0 || room.NetworkName != "" || room.NetworkSecret == "" {
		t.Fatalf("room id should be cleared: %+v", room)
	}
	if dock.removedDirs == 0 || dock.container(utils.LinkContainerName) != nil {
		t.Fatal("node and credential dir should be reset")
	}
	// 下一轮马上重新登记（不等退避），节点重建
	site.set(func(f *fakeLinkSite) { f.beatFail = nil; f.roomID = 13 })
	e.advance(15 * time.Second)
	m.tick()
	calls := site.callList()
	if calls[len(calls)-2] != "POST /link/rooms" || calls[len(calls)-1] != "POST /link/rooms/13/heartbeat" {
		t.Fatalf("calls = %v", calls)
	}
	if room := loadTestLinkRoom(t); room.RoomID != 13 {
		t.Fatalf("room = %+v", room)
	}
	if dock.container(utils.LinkContainerName) == nil {
		t.Fatal("node should be recreated")
	}
}

// 未绑定账号：状态为 login_required，不发请求，删除残留的节点容器；绑定后正常登记
func TestLinkUnboundLoginRequired(t *testing.T) {
	e := setupLinkTest(t)
	s, site, dock, m := e.s, e.site, e.dock, e.m
	database.DB.Where("id = ?", accountRowID).Delete(&models.ArkfgfAccount{})
	accountManager.forget()
	database.DB.Model(&s).Update("status", "running")
	s.Status = "running"
	e.rt.set(s.ID, linkRuntime{Ready: true})
	dock.addContainer(utils.LinkContainerName)

	m.tick()
	e.advance(time.Minute)
	m.tick()
	expectLinkCalls(t, site)
	if dock.container(utils.LinkContainerName) != nil || dock.createCount() != 0 {
		t.Fatal("node should not run without an account")
	}
	if st := m.statusFor(s); st.State != models.LinkLoginRequired || st.Error == "" {
		t.Fatalf("status = %+v", st)
	}
	if info, err := (&ServerService{}).GetLinkRoom(); err != nil || info.Status.State != models.LinkLoginRequired {
		t.Fatalf("room status = %+v err = %v", info, err)
	}

	bindTestAccount(t)
	e.advance(15 * time.Second)
	m.tick()
	expectLinkCalls(t, site, "POST /link/rooms", "POST /link/rooms/12/heartbeat")
	if dock.container(utils.LinkContainerName) == nil {
		t.Fatal("node should run after binding")
	}
}

// 升级后第一次登记：带旧版房间令牌认领旧房间（同一个房间 ID，节点和凭证保留），成功后清空，之后不再带
func TestLinkLegacyRoomClaim(t *testing.T) {
	e := setupLinkTest(t)
	site, dock, m := e.site, e.dock, e.m
	database.DB.Save(&models.LinkRoom{ID: linkRoomID, Token: "alr_old", RoomID: 12, NetworkName: "ark-0123456789abcdef",
		NetworkSecret: strings.Repeat("ab", 32), RPCPort: 20000})

	m.tick()
	expectLinkCalls(t, site, "POST /link/rooms", "POST /link/rooms/12/heartbeat")
	if site.tokens[0] != testAccountToken || site.legacy[0] != "alr_old" {
		t.Fatalf("tokens = %v legacy = %v", site.tokens, site.legacy)
	}
	if room := loadTestLinkRoom(t); room.Token != "" || room.RoomID != 12 {
		t.Fatalf("room = %+v", room)
	}
	if dock.removedDirs != 0 {
		t.Fatal("claimed room should keep its credentials")
	}

	// 资料变化重新登记时不再带 legacy_token
	database.DB.Model(&e.s).Update("session_name", "renamed")
	e.advance(15 * time.Second)
	m.tick()
	if len(site.legacy) != 2 || site.legacy[1] != "" {
		t.Fatalf("legacy = %v", site.legacy)
	}
}

// 旧房间认领不了（网站新建了另一个房间）：删除节点和凭证
func TestLinkLegacyRoomReplaced(t *testing.T) {
	e := setupLinkTest(t)
	site, dock, m := e.site, e.dock, e.m
	database.DB.Save(&models.LinkRoom{ID: linkRoomID, Token: "alr_old", RoomID: 7, NetworkName: "ark-old",
		NetworkSecret: strings.Repeat("ab", 32), RPCPort: 20000})
	dock.addContainer(utils.LinkContainerName)

	m.tick()
	if site.legacy[0] != "alr_old" {
		t.Fatalf("legacy = %v", site.legacy)
	}
	if room := loadTestLinkRoom(t); room.Token != "" || room.RoomID != 12 || room.NetworkName != "ark-0123456789abcdef" {
		t.Fatalf("room = %+v", room)
	}
	if dock.removedDirs == 0 {
		t.Fatal("credentials of the old room should be removed")
	}
}

// 没有服务器开启联机房、还留着旧版房间令牌：不需要绑定也用旧令牌下架
func TestLinkLegacyRoomDeleteWithoutAccount(t *testing.T) {
	e := setupLinkTest(t)
	s, site, dock, m := e.s, e.site, e.dock, e.m
	database.DB.Where("id = ?", accountRowID).Delete(&models.ArkfgfAccount{})
	accountManager.forget()
	database.DB.Model(&s).Update("link", false)
	database.DB.Save(&models.LinkRoom{ID: linkRoomID, Token: "alr_old", RoomID: 7, NetworkName: "ark-old",
		NetworkSecret: "sec", RPCPort: 20000, Name: "My Room"})

	m.tick()
	m.tick()
	expectLinkCalls(t, site, "DELETE /link/rooms/7")
	if site.tokens[0] != "alr_old" {
		t.Fatalf("token = %q", site.tokens[0])
	}
	room := loadTestLinkRoom(t)
	if room.Token != "" || room.RoomID != 0 || room.NetworkName != "" || room.NetworkSecret != "" || room.Name != "My Room" {
		t.Fatalf("room = %+v", room)
	}
	if dock.removedDirs == 0 {
		t.Fatal("credential dir should be removed")
	}
}

// 网站返回 invalid_account：视为已解绑，清空本地账号、房间登记和节点
func TestLinkInvalidAccountClearsLocal(t *testing.T) {
	e := setupLinkTest(t)
	s, site, dock, m := e.s, e.site, e.dock, e.m
	old := linkManager
	linkManager = m
	t.Cleanup(func() { linkManager = old })
	database.DB.Model(&s).Update("status", "running")
	s.Status = "running"
	// 另一台发布到收录站的服务器（发布和联机房只能二选一）：账号失效时 published 也要清除
	pub := createTestServer(t, 7790, 27030, 32340)
	database.DB.Model(&pub).Updates(map[string]any{"publish": true, "published": true})
	e.rt.set(s.ID, linkRuntime{Ready: true})
	m.tick()
	if dock.container(utils.LinkContainerName) == nil {
		t.Fatal("node not running")
	}

	site.set(func(f *fakeLinkSite) {
		f.fail = func(w http.ResponseWriter) {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"code": "invalid_account", "error": "invalid account"})
		}
	})
	e.advance(time.Minute)
	m.tick()
	if accountManager.bound() {
		t.Fatal("account should be cleared")
	}
	room := loadTestLinkRoom(t)
	if room.RoomID != 0 || room.NetworkName != "" || room.NetworkSecret != "" {
		t.Fatalf("room = %+v", room)
	}
	if dock.container(utils.LinkContainerName) != nil || dock.removedDirs == 0 {
		t.Fatal("node and credentials should be removed")
	}
	if loadPublishServer(t, pub.ID).Published {
		t.Fatal("published should be cleared")
	}
	if st := m.statusFor(s); st.State != models.LinkLoginRequired {
		t.Fatalf("status = %+v", st)
	}
	// 之后不再请求
	calls := len(site.callList())
	e.advance(time.Hour)
	m.tick()
	if len(site.callList()) != calls {
		t.Fatalf("calls = %v", site.callList())
	}
}

func TestLinkMembersProxy(t *testing.T) {
	e := setupLinkTest(t)
	site, m := e.site, e.m
	old := linkManager
	linkManager = m
	t.Cleanup(func() { linkManager = old })
	svc := NewServerService()

	// 还没登记：409
	if _, err := svc.GetLinkMembers(); !errors.Is(err, ErrLinkNotRegistered) {
		t.Fatalf("err = %v", err)
	}
	m.tick()
	// 没有绑定账号：409 account_required
	database.DB.Model(&models.ArkfgfAccount{}).Where("id = ?", accountRowID).Update("token", "")
	accountManager.forget()
	if _, err := svc.GetLinkMembers(); !errors.Is(err, ErrAccountRequired) {
		t.Fatalf("err = %v", err)
	}
	bindTestAccount(t)
	members, err := svc.GetLinkMembers()
	if err != nil || len(members.Members) != 1 || members.Members[0].SteamID != "76561198000000001" || members.Bans == nil {
		t.Fatalf("members = %+v err = %v", members, err)
	}
	if err := svc.KickLinkMember(models.LinkKickRequest{SteamID: "bad"}); !errors.Is(err, ErrInvalidSteamID) {
		t.Fatalf("err = %v", err)
	}
	if err := svc.KickLinkMember(models.LinkKickRequest{SteamID: "76561198000000001", Ban: true}); err != nil {
		t.Fatal(err)
	}
	// 网站没有这个接口时按网站拒绝处理
	if err := svc.UnbanLinkMember("76561198000000001"); !errors.Is(err, ErrLinkRequestFailed) {
		t.Fatalf("err = %v", err)
	}
	if !slices.Contains(site.callList(), "POST /link/rooms/12/kick") {
		t.Fatalf("calls = %v", site.callList())
	}
}

// 房间资料：设置页修改名称和收录资料后重新登记；都为空时用第一台服务器的
func TestLinkRoomSettings(t *testing.T) {
	e := setupLinkTest(t)
	s, site, m := e.s, e.site, e.m
	old := linkManager
	linkManager = m
	t.Cleanup(func() { linkManager = old })
	svc := NewServerService()
	database.DB.Model(&s).Updates(map[string]any{"listing_description": "server desc", "listing_tags": "pve,new"})

	m.tick()
	p := site.lastProfile()
	if p.Name != "test" || p.Description != "server desc" || !slices.Equal(p.Tags, []string{"pve", "new"}) {
		t.Fatalf("fallback profile = %+v", p)
	}

	// 校验
	if _, err := svc.UpdateLinkRoom(models.LinkRoomSettings{Name: strings.Repeat("长", 65)}); !errors.Is(err, ErrInvalidListing) {
		t.Fatalf("err = %v", err)
	}
	if _, err := svc.UpdateLinkRoom(models.LinkRoomSettings{Listing: models.ListingProfile{QQGroup: "abc"}}); !errors.Is(err, ErrInvalidListing) {
		t.Fatalf("err = %v", err)
	}

	info, err := svc.UpdateLinkRoom(models.LinkRoomSettings{Name: "  Our Cluster ",
		Listing: models.ListingProfile{Description: "room desc", Tags: []string{"pvp", "PVP", " x "}, QQGroup: "123456"}})
	if err != nil {
		t.Fatal(err)
	}
	if info.Name != "Our Cluster" || !slices.Equal(info.Listing.Tags, []string{"pvp", "x"}) || len(info.Servers) != 1 ||
		info.Servers[0].ID != s.ID || info.Status == nil {
		t.Fatalf("info = %+v", info)
	}
	e.advance(15 * time.Second)
	m.tick()
	p = site.lastProfile()
	if site.count("POST /link/rooms") != 2 || p.Name != "Our Cluster" || p.Description != "room desc" || p.QQGroup != "123456" ||
		!slices.Equal(p.Tags, []string{"pvp", "x"}) {
		t.Fatalf("profile = %+v", p)
	}
	// 服务器名称不受房间名称影响
	if p.Servers[0].Name != "test" {
		t.Fatalf("server profile = %+v", p.Servers[0])
	}

	// 没有服务器开启联机房时状态为 nil
	database.DB.Model(&s).Update("link", false)
	if info, err := svc.GetLinkRoom(); err != nil || info.Status != nil || len(info.Servers) != 0 {
		t.Fatalf("info = %+v err = %v", info, err)
	}
}

func TestLinkNodeArgsAndRPCPort(t *testing.T) {
	setupTestDB(t)
	// 另一台服务器的游戏端口正好占了 20000
	createTestServer(t, linkRPCPortBase, 27017, 32332)
	port, err := assignLinkRPCPort()
	if err != nil {
		t.Fatal(err)
	}
	if port != linkRPCPortBase+2 { // 顺延跳过游戏端口和 +1
		t.Fatalf("port = %d", port)
	}

	room := models.LinkRoom{NetworkName: "ark-x", NetworkSecret: "sec", RPCPort: 20003}
	// 查询端口和游戏端口重叠时去重
	servers := []models.Server{{ID: 3, SessionName: "", Port: 7777, QueryPort: 7777}}
	spec1 := linkContainerSpec(room, servers, []string{"tcp://a:1"})
	args := strings.Join(spec1.Args, " ")
	if !strings.Contains(args, "--udp-whitelist 7777 7778 --listeners") || !strings.Contains(args, "--hostname ARK Server") ||
		spec1.Image != utils.EasyTierImage {
		t.Fatalf("spec = %+v", spec1)
	}
	servers[0].Port = 7779
	spec2 := linkContainerSpec(room, servers, []string{"tcp://a:1"})
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

// 管理端重启后，节点容器还在运行、服务器还在运行：不删除节点，也不重建（参数不变），宽限期内仍报告在线；
// 之后一直查询不到时报告离线，但节点保留到服务器停止
func TestLinkKeepsNodeAcrossRestart(t *testing.T) {
	e := setupLinkTest(t)
	s, site, dock := e.s, e.site, e.dock
	name := utils.LinkContainerName
	database.DB.Model(&s).Update("status", "running")
	s.Status = "running"
	e.rt.set(s.ID, linkRuntime{Ready: true, Players: 3})
	e.m.tick()
	if dock.container(name) == nil || dock.createCount() != 1 {
		t.Fatalf("node not created: %d", dock.createCount())
	}

	// 重启：内存里的状态都没了，运行监控还没查到结果
	e.rt.set(s.ID, linkRuntime{})
	m := newTestLinkManager(t, site, dock, e.rt, e.now)
	e.advance(15 * time.Second)
	m.tick()
	if dock.container(name) == nil || dock.createCount() != 1 || len(dock.removed) != 0 {
		t.Fatalf("node should be kept across a restart: creates %d removed %v", dock.createCount(), dock.removed)
	}
	if sb, _ := site.lastBeat().server(fmt.Sprint(s.ID)); !sb.Online {
		t.Fatalf("server should still be reported online within the grace period: %+v", site.lastBeat())
	}

	// 一直查询不到：超过宽限期报告离线，节点保留
	e.advance(notReadyGrace + time.Second)
	m.tick()
	if dock.container(name) == nil {
		t.Fatal("node removed while the server is still running")
	}
	if sb, _ := site.lastBeat().server(fmt.Sprint(s.ID)); sb.Online {
		t.Fatalf("server should be reported offline after the grace period: %+v", site.lastBeat())
	}

	// 停服：照常删除节点
	database.DB.Model(&s).Update("status", "stopped")
	e.advance(15 * time.Second)
	m.tick()
	if dock.container(name) != nil {
		t.Fatal("node should be removed when the server stops")
	}
}

// 重启时服务器已经停了：残留的节点删除；重启前没有节点：照常等服务器可以连接后再创建
func TestLinkRestartWithoutRunningServer(t *testing.T) {
	e := setupLinkTest(t)
	s, dock := e.s, e.dock
	name := utils.LinkContainerName
	dock.addContainer(name)
	e.m.tick()
	if dock.container(name) != nil {
		t.Fatal("leftover node should be removed when no server is running")
	}

	m := newTestLinkManager(t, e.site, dock, e.rt, e.now)
	database.DB.Model(&s).Update("status", "running")
	e.advance(15 * time.Second)
	m.tick()
	if dock.createCount() != 0 {
		t.Fatal("node should wait until the server is ready")
	}
}

// 发布和联机房只能二选一：从发布切换过来的服务器等收录站取消发布（published 变成 false）后才放进房间；
// inRoom 按最近一次登记成功的房间回答，供发布等待服务器移出房间
func TestLinkWaitsForUnpublish(t *testing.T) {
	e := setupLinkTest(t)
	s, site, m := e.s, e.site, e.m
	database.DB.Model(&s).Updates(map[string]any{"status": "running", "published": true})
	e.rt.set(s.ID, linkRuntime{Ready: true})

	m.tick()
	expectLinkCalls(t, site)
	if st := m.statusFor(loadPublishServer(t, s.ID)); st.State != models.LinkWaiting {
		t.Fatalf("status = %+v", st)
	}
	if m.inRoom(s.ID) {
		t.Fatal("server should not be in the room yet")
	}

	database.DB.Model(&s).Update("published", false)
	e.advance(15 * time.Second)
	m.tick()
	expectLinkCalls(t, site, "POST /link/rooms", "POST /link/rooms/12/heartbeat")
	if !m.inRoom(s.ID) {
		t.Fatal("server should be in the room")
	}

	// 关闭联机房（切换到发布）：下架后不在房间里
	database.DB.Model(&s).Update("link", false)
	e.advance(15 * time.Second)
	m.tick()
	expectLinkCalls(t, site, "POST /link/rooms", "POST /link/rooms/12/heartbeat", "DELETE /link/rooms/12")
	if m.inRoom(s.ID) {
		t.Fatal("server should have left the room")
	}
}
