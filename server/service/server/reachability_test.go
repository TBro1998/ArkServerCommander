package server

import (
	"bytes"
	"errors"
	"net"
	"net/netip"
	"sync/atomic"
	"testing"
	"time"

	"ark-server-commander/database"
	"ark-server-commander/models"
	"ark-server-commander/utils"
)

// fakeArkQueryPort 在本机 UDP 端口上模拟 ARK 查询端口：先回 challenge，带上 challenge 后回 A2S_INFO。
// 返回端口和收到的查询包个数
func fakeArkQueryPort(t *testing.T, name, folder string) (int, *atomic.Int32) {
	t.Helper()
	conn, err := net.ListenPacket("udp4", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { conn.Close() })
	var count atomic.Int32
	go func() {
		buf := make([]byte, 1500)
		for {
			n, addr, err := conn.ReadFrom(buf)
			if err != nil {
				return
			}
			count.Add(1)
			if !bytes.HasSuffix(buf[:n], []byte{1, 2, 3, 4}) {
				conn.WriteTo([]byte{0xFF, 0xFF, 0xFF, 0xFF, 0x41, 1, 2, 3, 4}, addr)
				continue
			}
			resp := []byte{0xFF, 0xFF, 0xFF, 0xFF, 0x49, 17}
			game := "Counter-Strike"
			if folder == "ark_survival_evolved" {
				game = "ARK: Survival Evolved"
			}
			for _, s := range []string{name, "TheIsland", folder, game} {
				resp = append(resp, append([]byte(s), 0)...)
			}
			resp = append(resp, 0, 0, 0, 70, 0, 'd', 'l', 0, 0)
			conn.WriteTo(resp, addr)
		}
	}()
	return conn.LocalAddr().(*net.UDPAddr).Port, &count
}

// closedUDPPort 一个刚关闭的本机 UDP 端口，向它发包收不到响应
func closedUDPPort(t *testing.T) int {
	t.Helper()
	conn, err := net.ListenPacket("udp4", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	port := conn.LocalAddr().(*net.UDPAddr).Port
	conn.Close()
	return port
}

func TestQueryA2S(t *testing.T) {
	port, _ := fakeArkQueryPort(t, "My \x01Server - (v358.24)", "ark_survival_evolved")
	reply, err := queryA2S("127.0.0.1", port, time.Second)
	if err != nil || reply.Name != "My Server" || !reply.Ark {
		t.Fatalf("reply %+v err %v", reply, err)
	}
	if _, err := queryA2S("127.0.0.1", closedUDPPort(t), 200*time.Millisecond); err == nil {
		t.Fatal("closed port should fail")
	}
	if _, err := parseA2SReply([]byte{0xFF, 0xFF, 0xFF, 0xFF, 0x49, 17, 'a'}); err == nil {
		t.Fatal("truncated reply should fail")
	}
}

type reachEnv struct {
	r      *ReachabilityChecker
	server models.Server
	ready  *bool
	now    *time.Time
	ip     *string
	public *bool
	tgtErr *error
	port   *int // 实际向这个端口发查询（测试里把服务器的查询端口映射到本机的假端口）
	calls  *int
}

func setupReachTest(t *testing.T) *reachEnv {
	t.Helper()
	utils.InitLogger()
	setupTestDB(t)
	s := createTestServer(t, 7777, 27015, 32330)
	database.DB.Model(&s).Update("status", "running")
	s.Status = "running"

	ready, now := false, time.Now()
	ip, public := "127.0.0.1", true
	var tgtErr error
	port, calls := 0, 0
	r := &ReachabilityChecker{
		states: map[uint]*reachState{},
		now:    func() time.Time { return now },
		target: func() (string, bool, error) { return ip, public, tgtErr },
		query: func(qip string, qport int) (*a2sReply, error) {
			calls++
			if qport != 27015 {
				t.Errorf("query port = %d", qport)
			}
			return queryA2S(qip, port, 200*time.Millisecond)
		},
		ready: func(models.Server) bool { return ready },
		name:  func(models.Server) string { return "My Server - (v358.24)" },
	}
	return &reachEnv{r: r, server: s, ready: &ready, now: &now, ip: &ip, public: &public, tgtErr: &tgtErr, port: &port, calls: &calls}
}

func (e *reachEnv) state() *models.Reachability { return e.r.statusFor(e.server) }

func TestReachabilityChecksOncePerReady(t *testing.T) {
	e := setupReachTest(t)
	*e.port, _ = fakeArkQueryPort(t, "My Server - (v358.24)", "ark_survival_evolved")

	// 还不能连接：不检测，没有状态
	e.r.tick()
	if *e.calls != 0 || e.state() != nil {
		t.Fatalf("not ready: calls %d state %+v", *e.calls, e.state())
	}

	*e.ready = true
	if st := e.state(); st == nil || st.State != models.ReachChecking {
		t.Fatalf("before first check: %+v", st)
	}
	e.r.tick()
	st := e.state()
	if *e.calls != 1 || st.State != models.ReachReachable || st.IP != "127.0.0.1" || st.Name != "My Server" || st.CheckedAt == "" {
		t.Fatalf("reachable: calls %d %+v", *e.calls, st)
	}

	// 之后不再自动检测
	for range 20 {
		*e.now = e.now.Add(time.Minute)
		e.r.tick()
	}
	if *e.calls != 1 {
		t.Fatalf("should check once, got %d", *e.calls)
	}

	// 短暂查询不到（宽限期内）：不重新检测，结果保留
	*e.ready = false
	*e.now = e.now.Add(30 * time.Second)
	e.r.tick()
	*e.ready = true
	e.r.tick()
	if *e.calls != 1 || e.state().State != models.ReachReachable {
		t.Fatalf("grace: calls %d state %+v", *e.calls, e.state())
	}

	// 停服：清除结果，再次可以连接时重新检测
	*e.ready = false
	e.server.Status = "stopped"
	database.DB.Model(&e.server).Update("status", "stopped")
	e.r.tick()
	if e.state() != nil {
		t.Fatalf("stopped server should have no reachability: %+v", e.state())
	}
	e.server.Status = "running"
	database.DB.Model(&e.server).Update("status", "running")
	*e.ready = true
	e.r.tick()
	if *e.calls != 2 {
		t.Fatalf("should check again after restart, got %d", *e.calls)
	}
}

func TestReachabilityUnreachableAndMismatch(t *testing.T) {
	e := setupReachTest(t)
	*e.ready = true

	// 没有响应：重试一次后判为访问不到，之后不自动重试
	*e.port = closedUDPPort(t)
	e.r.tick()
	if st := e.state(); *e.calls != 2 || st.State != models.ReachUnreachable || st.IP != "127.0.0.1" || st.RetryAt != "" {
		t.Fatalf("unreachable: calls %d %+v", *e.calls, st)
	}
	for range 10 {
		*e.now = e.now.Add(10 * time.Minute)
		e.r.tick()
	}
	if *e.calls != 2 {
		t.Fatalf("unreachable should not be retried automatically, got %d", *e.calls)
	}

	// 手动重新检测：公网上这个端口是另一台服务器
	*e.port, _ = fakeArkQueryPort(t, "Neighbour - (v358.24)", "ark_survival_evolved")
	st, err := e.r.Recheck(e.server)
	if err != nil || st.State != models.ReachMismatch || st.Name != "Neighbour" {
		t.Fatalf("mismatch: %+v %v", st, err)
	}
	// 10 秒内再点不重新查询
	if _, err := e.r.Recheck(e.server); err != nil || *e.calls != 3 {
		t.Fatalf("recheck too soon: calls %d err %v", *e.calls, err)
	}
	// 回应的不是 ARK 也算 mismatch
	*e.now = e.now.Add(time.Minute)
	*e.port, _ = fakeArkQueryPort(t, "My Server", "csgo")
	if st, _ := e.r.Recheck(e.server); st.State != models.ReachMismatch {
		t.Fatalf("not ark: %+v", st)
	}
}

func TestReachabilityNotPublicAndNoAddress(t *testing.T) {
	e := setupReachTest(t)
	*e.ready = true

	// PUBLIC_HOST 是内网地址：不查询，提示公网玩家连不上
	*e.ip, *e.public = "192.168.1.10", false
	e.r.tick()
	if st := e.state(); *e.calls != 0 || st.State != models.ReachNotPublic || st.IP != "192.168.1.10" {
		t.Fatalf("not public: calls %d %+v", *e.calls, st)
	}

	// 拿不到公网地址：error，退避后重试
	*e.tgtErr = errNoPublicHost
	*e.now = e.now.Add(time.Minute)
	if _, err := e.r.Recheck(e.server); err != nil {
		t.Fatal(err)
	}
	st := e.state()
	if st.State != models.ReachError || st.ErrorCode != "public_ip_unknown" || st.RetryAt == "" {
		t.Fatalf("no address: %+v", st)
	}
	*e.tgtErr, *e.ip, *e.public = nil, "127.0.0.1", true
	*e.port, _ = fakeArkQueryPort(t, "My Server", "ark_survival_evolved")
	*e.now = e.now.Add(30 * time.Second)
	e.r.tick()
	if e.state().State != models.ReachError {
		t.Fatalf("should wait for backoff: %+v", e.state())
	}
	*e.now = e.now.Add(time.Minute)
	e.r.tick()
	if st := e.state(); st.State != models.ReachReachable || st.ErrorCode != "" {
		t.Fatalf("after retry: %+v", st)
	}
}

func TestReachabilityRecheckRequiresReady(t *testing.T) {
	e := setupReachTest(t)
	if _, err := e.r.Recheck(e.server); !errors.Is(err, ErrServerNotReady) {
		t.Fatalf("not ready: %v", err)
	}
}

func TestIsPublicIPv4(t *testing.T) {
	for ip, want := range map[string]bool{
		"8.8.8.8": true, "203.0.113.5": true, "192.168.1.1": false, "10.0.0.1": false,
		"172.16.0.1": false, "127.0.0.1": false, "100.64.1.1": false, "169.254.1.1": false,
	} {
		t.Run(ip, func(t *testing.T) {
			addr := netip.MustParseAddr(ip)
			if got := isPublicIPv4(addr); got != want {
				t.Errorf("isPublicIPv4(%s) = %v", ip, got)
			}
		})
	}
}

func TestSameServerName(t *testing.T) {
	for _, c := range []struct {
		remote, local string
		same          bool
	}{
		{"My Server", "My Server - (v358.24)", true},
		{"中文服务器", "中文服务器 - (v358.24)", true},
		{"Long name truncated", "Long name truncated at sixty four - (v358.24)", true},
		{"", "My Server", true},
		{"Neighbour", "My Server - (v358.24)", false},
	} {
		if got := sameServerName(c.remote, c.local); got != c.same {
			t.Errorf("sameServerName(%q, %q) = %v", c.remote, c.local, got)
		}
	}
}
