package server

import (
	"bytes"
	"errors"
	"net"
	"net/netip"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode"

	"ark-server-commander/database"
	"ark-server-commander/models"
	"ark-server-commander/utils"

	"go.uber.org/zap"
)

// 公网可达性检测。
//
// 运行监控（runtime.go）的 A2S 是在容器里向 127.0.0.1 发的，只说明游戏进程起来了，不说明玩家能从公网连进来：
// 服主可能没有公网 IP（运营商 NAT），路由器没做端口转发，或者防火墙没放行。
// 所以服务器可以连接后，管理端再用玩家连接的地址（PUBLIC_HOST，或自动探测的公网 IP，见 PublicHost）
// 向查询端口发一次 A2S_INFO：没有公网 IP 时，探测到的是运营商 NAT 的出口地址，发过去没有响应。
//
// 这是从管理端自己查的，结果有两种不准：
//   - 路由器不支持 NAT 回环时，端口转发正确也查不到（报成访问不到）
//   - 机器自己就有公网 IP 时，查自己不经过外部的防火墙 / 云服务器安全组（安全组没放行也报成能访问）
//
// 服务器每次变成可以连接（ready）后检测一次；运行中短暂查询不到（notReadyGrace 内）不算断开，不重新检测。
// 有结论后不再自动检测，服主修好端口转发后在详情页手动重新检测；
// 拿不到公网地址时按 1、2、4…30 分钟退避重试。只查 IPv4（ARK 客户端只能用 IPv4 连接）。

const (
	reachTickInterval   = 15 * time.Second
	reachBackoffBase    = time.Minute
	reachBackoffMax     = 30 * time.Minute
	reachQueryTimeout   = 3 * time.Second
	reachManualInterval = 10 * time.Second // 手动重新检测的最短间隔
)

// reachState 一台服务器的检测状态，只在内存里；停服后清除，下次可以连接时重新检测
type reachState struct {
	done     bool // 这一次可以连接之后已经有结论
	checking bool
	result   models.Reachability
	nextAt   time.Time // 拿不到公网地址时下一次重试的时间
	failures int
	lastAt   time.Time // 最近一次检测的时间
	// 运行中但不是 ready 的开始时间，见 notReadyGrace
	notReadySince time.Time
}

// a2sReply 公网查询到的服务器信息
type a2sReply struct {
	Name string
	Ark  bool // folder / game 里有 ark
}

// ReachabilityChecker 定期检测运行中服务器的公网可达性
type ReachabilityChecker struct {
	mu     sync.Mutex
	states map[uint]*reachState
	now    func() time.Time
	// 以下测试时替换：target 返回要检测的 IPv4 和它是不是公网地址；query 向 ip:port 发 A2S_INFO；
	// ready 判断服务器能否连接；name 返回运行监控在容器里查到的服务器名称
	target func() (string, bool, error)
	query  func(ip string, port int) (*a2sReply, error)
	ready  func(models.Server) bool
	name   func(models.Server) string
}

var reachability = &ReachabilityChecker{
	states: map[uint]*reachState{},
	now:    time.Now,
	target: publicTarget,
	query: func(ip string, port int) (*a2sReply, error) {
		return queryA2S(ip, port, reachQueryTimeout)
	},
	ready: serverReady,
	name: func(server models.Server) string {
		if info := runtimeMonitor.Get(server.ID); info != nil {
			return info.ServerName
		}
		return ""
	},
}

var errNoPublicHost = errors.New("没有拿到公网地址")

// publicTarget 玩家连接的地址（PUBLIC_HOST 或自动探测的公网 IP）对应的 IPv4，以及它是不是公网地址。
// PUBLIC_HOST 是域名时解析成 IPv4
func publicTarget() (string, bool, error) {
	host := PublicHost()
	if host == "" {
		return "", false, errNoPublicHost
	}
	addr, err := netip.ParseAddr(host)
	if err != nil {
		ips, err := net.LookupIP(host)
		if err != nil {
			return "", false, err
		}
		for _, ip := range ips {
			if v4 := ip.To4(); v4 != nil {
				addr = netip.AddrFrom4([4]byte(v4))
				break
			}
		}
		if !addr.IsValid() {
			return "", false, errors.New(host + " 没有 IPv4 地址")
		}
	}
	addr = addr.Unmap()
	if !addr.Is4() {
		return "", false, errors.New(host + " 不是 IPv4 地址")
	}
	return addr.String(), isPublicIPv4(addr), nil
}

var nonPublicPrefixes = []netip.Prefix{
	netip.MustParsePrefix("0.0.0.0/8"),
	netip.MustParsePrefix("100.64.0.0/10"), // 运营商级 NAT
	netip.MustParsePrefix("192.0.0.0/24"),
	netip.MustParsePrefix("198.18.0.0/15"),
	netip.MustParsePrefix("240.0.0.0/4"),
}

// isPublicIPv4 内网、回环、链路本地、运营商级 NAT 等地址公网上的玩家连不上
func isPublicIPv4(addr netip.Addr) bool {
	if !addr.Is4() || !addr.IsGlobalUnicast() || addr.IsPrivate() || addr.IsLoopback() || addr.IsLinkLocalUnicast() {
		return false
	}
	for _, p := range nonPublicPrefixes {
		if p.Contains(addr) {
			return false
		}
	}
	return true
}

// StartReachabilityChecker 启动后台检测协程
func StartReachabilityChecker() {
	go func() {
		ticker := time.NewTicker(reachTickInterval)
		defer ticker.Stop()
		for range ticker.C {
			reachability.tick()
		}
	}()
}

func (r *ReachabilityChecker) tick() {
	var servers []models.Server
	if err := database.DB.Select("id", "status", "query_port").Where("status = ?", "running").Find(&servers).Error; err != nil {
		utils.Warn("公网可达性检测：读取服务器列表失败", zap.Error(err))
		return
	}
	keep := map[uint]bool{}
	for _, server := range servers {
		keep[server.ID] = true
		if r.due(server) {
			r.check(server)
		}
	}
	r.mu.Lock()
	for id := range r.states {
		if !keep[id] {
			delete(r.states, id)
		}
	}
	r.mu.Unlock()
}

// due 这台服务器现在该不该检测；顺便处理断开（超过宽限期不能连接就清掉结论，下次可以连接时重新检测）
func (r *ReachabilityChecker) due(server models.Server) bool {
	ready := r.ready(server)
	r.mu.Lock()
	defer r.mu.Unlock()
	st := r.stateLocked(server.ID)
	now := r.now()
	if !ready {
		if st.done || st.failures > 0 {
			if st.notReadySince.IsZero() {
				st.notReadySince = now
			}
			if now.Sub(st.notReadySince) < notReadyGrace {
				return false
			}
		}
		*st = reachState{}
		return false
	}
	st.notReadySince = time.Time{}
	return !st.done && !st.checking && !now.Before(st.nextAt)
}

func (r *ReachabilityChecker) stateLocked(id uint) *reachState {
	st, ok := r.states[id]
	if !ok {
		st = &reachState{}
		r.states[id] = st
	}
	return st
}

// check 用公网地址查询一次，结果记在内存里。同一台服务器同一时刻只有一个检测
func (r *ReachabilityChecker) check(server models.Server) {
	r.mu.Lock()
	st := r.stateLocked(server.ID)
	if st.checking {
		r.mu.Unlock()
		return
	}
	st.checking = true
	st.lastAt = r.now()
	r.mu.Unlock()

	result, err := r.probe(server)

	r.mu.Lock()
	defer r.mu.Unlock()
	st = r.stateLocked(server.ID)
	st.checking = false
	now := r.now()
	if err != nil {
		st.failures++
		st.nextAt = now.Add(min(reachBackoffBase<<min(st.failures-1, 5), reachBackoffMax))
		st.result = models.Reachability{State: models.ReachError, ErrorCode: "public_ip_unknown",
			Error: "没有拿到公网地址（" + err.Error() + "），可通过环境变量 PUBLIC_HOST 设置", RetryAt: st.nextAt.Format(time.RFC3339)}
		if st.failures == 1 || st.failures%10 == 0 {
			utils.Warn("公网可达性检测：没有拿到公网地址", zap.Uint("server_id", server.ID), zap.Int("failures", st.failures), zap.Error(err))
		}
		return
	}
	st.failures, st.nextAt, st.done = 0, time.Time{}, true
	result.CheckedAt = now.Format(time.RFC3339)
	st.result = result
	utils.Info("公网可达性检测完成", zap.Uint("server_id", server.ID), zap.String("state", result.State), zap.String("ip", result.IP))
}

// probe 向公网地址的查询端口发 A2S_INFO（UDP 可能丢包，没响应时再试一次）。拿不到公网地址时返回错误
func (r *ReachabilityChecker) probe(server models.Server) (models.Reachability, error) {
	ip, public, err := r.target()
	if err != nil {
		return models.Reachability{}, err
	}
	if !public {
		// PUBLIC_HOST 设成了内网地址：查得到也只说明局域网里能连
		return models.Reachability{State: models.ReachNotPublic, IP: ip}, nil
	}
	var reply *a2sReply
	for range 2 {
		if reply, err = r.query(ip, server.QueryPort); err == nil {
			break
		}
	}
	if err != nil {
		return models.Reachability{State: models.ReachUnreachable, IP: ip}, nil
	}
	result := models.Reachability{State: models.ReachReachable, IP: ip, Name: reply.Name}
	if !reply.Ark || !sameServerName(reply.Name, r.name(server)) {
		result.State = models.ReachMismatch
	}
	return result, nil
}

// Recheck 手动重新检测（服主修改了端口转发或防火墙之后），同步等待结果
func (r *ReachabilityChecker) Recheck(server models.Server) (*models.Reachability, error) {
	if !r.ready(server) {
		return nil, ErrServerNotReady
	}
	r.mu.Lock()
	st := r.stateLocked(server.ID)
	recent := st.checking || (!st.lastAt.IsZero() && r.now().Sub(st.lastAt) < reachManualInterval)
	if !recent {
		st.done, st.nextAt, st.failures = false, time.Time{}, 0
	}
	r.mu.Unlock()
	if !recent {
		r.check(server)
	}
	return r.statusFor(server), nil
}

// statusFor 返回服务器的公网可达性，没有可显示的结果时为 nil
func (r *ReachabilityChecker) statusFor(server models.Server) *models.Reachability {
	if server.Status != "running" {
		return nil
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	st, ok := r.states[server.ID]
	switch {
	case ok && st.checking:
		result := st.result
		result.State = models.ReachChecking
		return &result
	case ok && st.result.State != "":
		result := st.result
		return &result
	case r.ready(server):
		// 可以连接了但下一轮才检测
		return &models.Reachability{State: models.ReachChecking}
	}
	return nil
}

// queryA2S 从管理端向 ip:port 发 A2S_INFO，返回名称（去掉版本号）和是不是 ARK。
// 新版协议可能先返回 challenge（0x41），需要带上 challenge 重发
func queryA2S(ip string, port int, timeout time.Duration) (*a2sReply, error) {
	conn, err := net.Dial("udp4", net.JoinHostPort(ip, strconv.Itoa(port)))
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	_ = conn.SetDeadline(time.Now().Add(timeout))
	req := []byte("\xFF\xFF\xFF\xFFTSource Engine Query\x00")
	if _, err := conn.Write(req); err != nil {
		return nil, err
	}
	buf := make([]byte, 1400)
	n, err := conn.Read(buf)
	if err != nil {
		return nil, err
	}
	for range 3 {
		if n < 9 || buf[4] != 0x41 {
			break
		}
		if _, err := conn.Write(append(append([]byte{}, req...), buf[5:9]...)); err != nil {
			return nil, err
		}
		if n, err = conn.Read(buf); err != nil {
			return nil, err
		}
	}
	return parseA2SReply(buf[:n])
}

// parseA2SReply 解析 A2S_INFO 响应（0x49）：协议版本之后依次是名称、地图、目录、游戏名
func parseA2SReply(data []byte) (*a2sReply, error) {
	if len(data) < 6 || !bytes.Equal(data[:4], []byte{0xFF, 0xFF, 0xFF, 0xFF}) || data[4] != 0x49 {
		return nil, errors.New("unexpected A2S response")
	}
	fields := make([]string, 0, 4)
	rest := data[6:]
	for range 4 {
		i := bytes.IndexByte(rest, 0)
		if i < 0 {
			return nil, errors.New("truncated A2S response")
		}
		fields = append(fields, strings.ToValidUTF8(string(rest[:i]), ""))
		rest = rest[i+1:]
	}
	name := strings.TrimSpace(serverNameVersionRe.ReplaceAllString(fields[0], ""))
	name = strings.Map(func(r rune) rune {
		if unicode.IsControl(r) {
			return -1
		}
		return r
	}, name)
	if r := []rune(name); len(r) > 64 {
		name = string(r[:64])
	}
	return &a2sReply{Name: name, Ark: strings.Contains(strings.ToLower(fields[2]+" "+fields[3]), "ark")}, nil
}

var serverNameVersionRe = regexp.MustCompile(`\s*-?\s*\(v[0-9][0-9.]*\)\s*$`)

// sameServerName 公网上查到的名称是不是这台服务器：只比较字母和数字，一边是另一边的前缀也算相同；
// 任一边为空时无法判断，算相同
func sameServerName(remote, local string) bool {
	norm := func(s string) string {
		s = serverNameVersionRe.ReplaceAllString(s, "")
		return strings.Map(func(r rune) rune {
			if unicode.IsLetter(r) || unicode.IsDigit(r) {
				return unicode.ToLower(r)
			}
			return -1
		}, s)
	}
	a, b := norm(remote), norm(local)
	if a == "" || b == "" {
		return true
	}
	return strings.HasPrefix(a, b) || strings.HasPrefix(b, a)
}
