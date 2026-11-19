package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync"
	"testing"
	"time"

	"ark-server-commander/database"
	"ark-server-commander/models"
	"ark-server-commander/utils"
)

const testAccountToken = "acm_test"

// bindTestAccount 直接写入一个已绑定的账号
func bindTestAccount(t *testing.T) {
	t.Helper()
	if err := database.DB.Save(&models.ArkfgfAccount{ID: accountRowID, Token: testAccountToken, CommanderID: 7,
		CommanderName: "test", SteamID: "76561198000000009", UserName: "owner", BoundAt: time.Now()}).Error; err != nil {
		t.Fatal(err)
	}
	accountManager.forget()
}

// useTestLinkManager 把全局的 linkManager 换成使用假 Docker 的实例（清除本地登记时会删除节点容器和凭证目录）
func useTestLinkManager(t *testing.T) *fakeLinkDocker {
	t.Helper()
	dock := newFakeLinkDocker()
	now := time.Now()
	m := newTestLinkManager(t, &fakeLinkSite{}, dock, &fakeRuntimes{rt: map[uint]linkRuntime{}}, &now)
	old := linkManager
	linkManager = m
	t.Cleanup(func() { linkManager = old })
	return dock
}

// fakeAccountSite 模拟网站的登录和账号接口
type fakeAccountSite struct {
	mu        sync.Mutex
	calls     []string
	tokens    []string
	login     map[string]string // POST /commander/login 的请求体
	finish    map[string]string // POST /commander/login/finish 的 params
	finishErr string            // finish 返回的错误码，空为成功
	deleteErr string            // DELETE /commander/me 返回的错误码，空为成功
	starts    int               // 签发过的 state 数，第 n 个为 acs_<n>
	noState   bool              // 返回的 Steam 地址里不带 state
	// 不为 nil 时，login/finish 等它关闭后才回复（模拟网站向 Steam 核验要几秒）。只在启动前设置
	finishGate chan struct{}
}

func (f *fakeAccountSite) handler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/commander/login/finish" && f.finishGate != nil {
		<-f.finishGate
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	f.calls = append(f.calls, r.Method+" "+r.URL.Path)
	token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
	f.tokens = append(f.tokens, token)
	switch {
	case r.Method == http.MethodPost && r.URL.Path == "/commander/login":
		_ = json.NewDecoder(r.Body).Decode(&f.login)
		f.starts++
		returnTo := f.login["return_to"] + fmt.Sprintf("?state=acs_%d", f.starts)
		if f.noState {
			returnTo = f.login["return_to"]
		}
		json.NewEncoder(w).Encode(map[string]any{
			"steam_url": "https://steamcommunity.com/openid/login?openid.mode=checkid_setup&openid.return_to=" + url.QueryEscape(returnTo) +
				"&openid.realm=" + url.QueryEscape("http://192.168.1.5:8080"),
			"expires_in": 600,
		})
	case r.Method == http.MethodPost && r.URL.Path == "/commander/login/finish":
		var body struct {
			Params map[string]string `json:"params"`
		}
		_ = json.NewDecoder(r.Body).Decode(&body)
		f.finish = body.Params
		switch f.finishErr {
		case "":
			json.NewEncoder(w).Encode(map[string]any{
				"token":         "acm_new",
				"commander":     map[string]any{"id": 7, "name": "我的开服机"},
				"user":          map[string]any{"steam_id": "76561198000000009", "name": "Owner", "avatar_url": "https://a/x.jpg"},
				"dashboard_url": "https://arkfgf.com/dashboard",
			})
		case "too_many_commanders":
			writeAPIError(w, http.StatusConflict, f.finishErr)
		case "busy":
			writeAPIError(w, http.StatusServiceUnavailable, f.finishErr)
		default:
			writeAPIError(w, http.StatusBadRequest, f.finishErr)
		}
	case r.Method == http.MethodDelete && r.URL.Path == "/commander/me":
		if f.deleteErr != "" {
			writeAPIError(w, http.StatusUnauthorized, f.deleteErr)
			return
		}
		json.NewEncoder(w).Encode(map[string]any{"ok": true})
	case r.Method == http.MethodGet && r.URL.Path == "/commander/me":
		if token != testAccountToken {
			writeAPIError(w, http.StatusUnauthorized, "invalid_account")
			return
		}
		json.NewEncoder(w).Encode(map[string]any{
			"commander": map[string]any{"id": 7, "name": "renamed"},
			"user":      map[string]any{"steam_id": "76561198000000009", "name": "New Nick", "avatar_url": "https://a/y.jpg"},
		})
	default:
		writeAPIError(w, http.StatusNotFound, "not_found")
	}
}

func (f *fakeAccountSite) set(fn func(f *fakeAccountSite)) {
	f.mu.Lock()
	defer f.mu.Unlock()
	fn(f)
}

func (f *fakeAccountSite) callList() []string {
	f.mu.Lock()
	defer f.mu.Unlock()
	return append([]string{}, f.calls...)
}

// setupAccountTest 换上指向模拟网站的 AccountManager 和假的联机房 Docker
func setupAccountTest(t *testing.T) (*fakeAccountSite, *AccountManager, *fakeLinkDocker) {
	t.Helper()
	setupTestDB(t)
	utils.InitLogger()
	site := &fakeAccountSite{}
	srv := httptest.NewServer(http.HandlerFunc(site.handler))
	t.Cleanup(srv.Close)
	now := time.Now()
	a := &AccountManager{client: srv.Client(), api: func() string { return srv.URL }, now: func() time.Time { return now }}
	old := accountManager
	accountManager = a
	t.Cleanup(func() { accountManager = old })
	return site, a, useTestLinkManager(t)
}

// callbackParams Steam 回调参数，state 为第 n 次开始登录得到的
func callbackParams(n int) map[string]string {
	return map[string]string{
		"state": fmt.Sprintf("acs_%d", n), "openid.ns": "http://specs.openid.net/auth/2.0", "openid.mode": "id_res",
		"openid.claimed_id": "https://steamcommunity.com/openid/id/76561198000000009",
	}
}

func startTestLogin(t *testing.T, a *AccountManager) {
	t.Helper()
	if _, err := a.StartLogin("", "http://192.168.1.5:8080/arkfgf-callback"); err != nil {
		t.Fatal(err)
	}
}

func finishCalls(site *fakeAccountSite) int {
	n := 0
	for _, c := range site.callList() {
		if c == "POST /commander/login/finish" {
			n++
		}
	}
	return n
}

func TestAccountLoginStart(t *testing.T) {
	site, a, _ := setupAccountTest(t)
	steamURL, err := a.StartLogin("  ", "http://192.168.1.5:8080/settings/arkfgf-callback")
	if err != nil || !strings.HasPrefix(steamURL, "https://steamcommunity.com/openid/login?") {
		t.Fatalf("steam_url = %q err = %v", steamURL, err)
	}
	site.mu.Lock()
	if site.login["name"] != accountDefaultName || site.login["version"] != utils.Version ||
		site.login["return_to"] != "http://192.168.1.5:8080/settings/arkfgf-callback" || site.tokens[0] != "" {
		t.Fatalf("login request = %v token %q", site.login, site.tokens[0])
	}
	site.mu.Unlock()

	// 回调地址和名称校验，不发请求
	for _, bad := range []string{"", "ftp://x/cb", "/settings/arkfgf-callback", "http:///cb", "http://h/cb?x=1", "http://h/cb#a",
		"http://user@h/cb", "http://h/" + strings.Repeat("a", 260)} {
		if _, err := a.StartLogin("", bad); !errors.Is(err, ErrInvalidReturnTo) {
			t.Errorf("return_to %q: err = %v", bad, err)
		}
	}
	if _, err := a.StartLogin(strings.Repeat("长", 65), "https://h/cb"); !errors.Is(err, ErrInvalidCommanderName) {
		t.Fatalf("err = %v", err)
	}
	if name, _ := normalizeCommanderName(" a\x00b\n "); name != "ab" {
		t.Fatalf("name = %q", name)
	}
	if got := site.callList(); len(got) != 1 {
		t.Fatalf("calls = %v", got)
	}

	// 已绑定
	bindTestAccount(t)
	if _, err := a.StartLogin("", "https://h/cb"); !errors.Is(err, ErrAccountBound) {
		t.Fatalf("err = %v", err)
	}
	// ARKFGF_API 为空
	a.api = func() string { return "" }
	if _, err := a.StartLogin("", "https://h/cb"); !errors.Is(err, ErrArkfgfDisabled) {
		t.Fatalf("err = %v", err)
	}
	if st := a.Status(false); st.Enabled {
		t.Fatalf("status = %+v", st)
	}

	// 连不上网站
	database.DB.Where("id = ?", accountRowID).Delete(&models.ArkfgfAccount{})
	a.forget()
	a.api = func() string { return "http://127.0.0.1:1" }
	if _, err := a.StartLogin("", "https://h/cb"); !errors.Is(err, ErrArkfgfUnavailable) {
		t.Fatalf("err = %v", err)
	}
}

func TestAccountLoginFinish(t *testing.T) {
	site, a, _ := setupAccountTest(t)
	// 绑定后清掉发布的旧状态并马上处理一轮
	old := publisher
	publisher = &Publisher{states: map[uint]*publishState{7: {errCode: "x"}}, api: func() string { return "" }, wakeCh: make(chan struct{}, 1)}
	t.Cleanup(func() { publisher = old })

	startTestLogin(t, a)
	if st := a.Status(false); !st.Login.Pending || st.Login.ExpiresAt == "" || st.Login.ErrorCode != "" {
		t.Fatalf("login status = %+v", st.Login)
	}
	if err := a.FinishLogin(callbackParams(1)); err != nil {
		t.Fatal(err)
	}
	site.mu.Lock()
	if site.finish["state"] != "acs_1" || site.finish["openid.claimed_id"] != callbackParams(1)["openid.claimed_id"] || site.tokens[1] != "" {
		t.Fatalf("finish params = %v tokens %v", site.finish, site.tokens)
	}
	site.mu.Unlock()
	st := a.Status(false)
	if !st.Bound || st.Account == nil || st.Account.UserName != "Owner" || st.Account.CommanderName != "我的开服机" ||
		st.Account.CommanderID != 7 || st.DashboardURL != "https://arkfgf.com/dashboard" || st.Account.BoundAt == "" || st.Login.Pending {
		t.Fatalf("status = %+v account = %+v", st, st.Account)
	}
	// 账号令牌保存在数据库里，状态里不返回
	var saved models.ArkfgfAccount
	database.DB.First(&saved, accountRowID)
	if saved.Token != "acm_new" || saved.SteamID != "76561198000000009" || saved.AvatarURL != "https://a/x.jpg" {
		t.Fatalf("saved = %+v", saved)
	}
	if data, _ := json.Marshal(st); strings.Contains(string(data), "acm_new") {
		t.Fatalf("token leaked: %s", data)
	}
	if len(publisher.states) != 0 || len(publisher.wakeCh) != 1 {
		t.Fatalf("publisher not reset / woken: %d states, %d wakes", len(publisher.states), len(publisher.wakeCh))
	}
}

// 网站核验期间（state 已经取走）设置页轮询看到的仍是登录中，不会误报过期；
// 结束后要么已绑定，要么带着失败原因
func TestAccountLoginPendingWhileFinishing(t *testing.T) {
	for _, fail := range []bool{false, true} {
		site, a, _ := setupAccountTest(t)
		gate := make(chan struct{})
		site.finishGate = gate
		if fail {
			site.finishErr = "steam_verify_failed"
		}
		startTestLogin(t, a)
		done := make(chan error, 1)
		go func() { done <- a.FinishLogin(callbackParams(1)) }()
		// 等 state 被取走
		deadline := time.Now().Add(5 * time.Second)
		for {
			a.mu.Lock()
			taken := len(a.logins) == 0 && a.finishing == 1
			a.mu.Unlock()
			if taken {
				break
			}
			if time.Now().After(deadline) {
				t.Fatal("state not consumed")
			}
			time.Sleep(5 * time.Millisecond)
		}
		if st := a.Status(false); st.Bound || !st.Login.Pending || st.Login.ErrorCode != "" {
			t.Fatalf("fail=%v: status while finishing = %+v", fail, st)
		}
		close(gate)
		err := <-done
		st := a.Status(false)
		if st.Login.Pending {
			t.Fatalf("fail=%v: still pending after finish: %+v", fail, st.Login)
		}
		if fail {
			if err == nil || st.Bound || st.Login.ErrorCode != "steam_verify_failed" {
				t.Fatalf("err = %v status = %+v", err, st)
			}
		} else if err != nil || !st.Bound {
			t.Fatalf("err = %v status = %+v", err, st)
		}
	}
}

// 网站返回的 Steam 地址里没有 state：无法完成登录，报错
func TestAccountLoginStartWithoutState(t *testing.T) {
	site, a, _ := setupAccountTest(t)
	site.set(func(f *fakeAccountSite) { f.noState = true })
	if _, err := a.StartLogin("", "https://h/cb"); !errors.Is(err, ErrArkfgfUnavailable) {
		t.Fatalf("err = %v", err)
	}
	if st := a.Status(false); st.Login.Pending {
		t.Fatalf("login status = %+v", st.Login)
	}
}

// 完成登录的接口是公开的：只接受本机开始的、没过期、没用过的 state，否则不联系网站
func TestAccountLoginFinishRequiresLocalState(t *testing.T) {
	site, a, _ := setupAccountTest(t)

	// 没有开始过登录
	if err := a.FinishLogin(callbackParams(1)); !errors.Is(err, ErrLoginExpired) {
		t.Fatalf("err = %v", err)
	}
	// 没有 state
	if err := a.FinishLogin(map[string]string{"openid.mode": "id_res"}); !errors.Is(err, ErrLoginExpired) {
		t.Fatalf("err = %v", err)
	}
	startTestLogin(t, a)
	// 参数太多
	many := callbackParams(1)
	for i := range accountParamsMax {
		many[fmt.Sprintf("x%d", i)] = "y"
	}
	if err := a.FinishLogin(many); !errors.Is(err, ErrLoginExpired) {
		t.Fatalf("err = %v", err)
	}
	// 别的 state
	if err := a.FinishLogin(callbackParams(9)); !errors.Is(err, ErrLoginExpired) {
		t.Fatalf("err = %v", err)
	}
	if finishCalls(site) != 0 {
		t.Fatalf("calls = %v", site.callList())
	}
	// 没有匹配的 state 时不记失败原因（谁都能调用这个接口）
	if st := a.Status(false); st.Login.ErrorCode != "" || !st.Login.Pending {
		t.Fatalf("login status = %+v", st.Login)
	}

	// state 只能用一次：第一次网站核验失败，记下原因；再用同一个 state 不再联系网站
	site.set(func(f *fakeAccountSite) { f.finishErr = "steam_verify_failed" })
	if err := a.FinishLogin(callbackParams(1)); !errors.Is(err, ErrSteamVerifyFailed) {
		t.Fatalf("err = %v", err)
	}
	if st := a.Status(false); st.Login.Pending || st.Login.ErrorCode != "steam_verify_failed" || st.Login.Error == "" {
		t.Fatalf("login status = %+v", st.Login)
	}
	site.set(func(f *fakeAccountSite) { f.finishErr = "" })
	if err := a.FinishLogin(callbackParams(1)); !errors.Is(err, ErrLoginExpired) || finishCalls(site) != 1 {
		t.Fatalf("err = %v calls = %v", err, site.callList())
	}

	// 开始新的登录清除失败原因；旧的 state 在过期前仍然有效
	startTestLogin(t, a) // acs_2
	startTestLogin(t, a) // acs_3
	if st := a.Status(false); st.Login.ErrorCode != "" || !st.Login.Pending {
		t.Fatalf("login status = %+v", st.Login)
	}
	// 过期
	now := a.now()
	a.now = func() time.Time { return now.Add(accountLoginTTL) }
	if err := a.FinishLogin(callbackParams(2)); !errors.Is(err, ErrLoginExpired) || finishCalls(site) != 1 {
		t.Fatalf("err = %v", err)
	}
	if st := a.Status(false); st.Login.Pending {
		t.Fatalf("login status = %+v", st.Login)
	}
	a.now = func() time.Time { return now.Add(accountLoginTTL + time.Second) }
	startTestLogin(t, a) // acs_4；acs_2 / acs_3 已过期被清掉
	if err := a.FinishLogin(callbackParams(3)); !errors.Is(err, ErrLoginExpired) {
		t.Fatalf("err = %v", err)
	}
	if err := a.FinishLogin(callbackParams(4)); err != nil {
		t.Fatal(err)
	}
	if !a.bound() || finishCalls(site) != 2 {
		t.Fatalf("bound = %v calls = %v", a.bound(), site.callList())
	}
}

// 最多记 accountPendingMax 个进行中的登录，超出时丢掉最早的
func TestAccountLoginPendingLimit(t *testing.T) {
	_, a, _ := setupAccountTest(t)
	for range accountPendingMax + 1 {
		startTestLogin(t, a)
	}
	if err := a.FinishLogin(callbackParams(1)); !errors.Is(err, ErrLoginExpired) {
		t.Fatalf("oldest state should be dropped: %v", err)
	}
	if err := a.FinishLogin(callbackParams(2)); err != nil {
		t.Fatal(err)
	}
}

func TestAccountLoginFinishErrors(t *testing.T) {
	cases := map[string]error{
		"expired_login":       ErrLoginExpired,
		"login_cancelled":     ErrLoginCancelled,
		"steam_verify_failed": ErrSteamVerifyFailed,
		"too_many_commanders": ErrTooManyCommanders,
		"busy":                ErrArkfgfUnavailable,
		"invalid_input":       ErrArkfgfRejected,
	}
	for code, want := range cases {
		t.Run(code, func(t *testing.T) {
			site, a, _ := setupAccountTest(t)
			site.set(func(f *fakeAccountSite) { f.finishErr = code })
			startTestLogin(t, a)
			if err := a.FinishLogin(callbackParams(1)); !errors.Is(err, want) {
				t.Fatalf("err = %v, want %v", err, want)
			}
			if a.bound() {
				t.Fatal("should not be bound")
			}
			if st := a.Status(false); st.Login.ErrorCode == "" || st.Login.Error == "" {
				t.Fatalf("login status = %+v", st.Login)
			}
		})
	}
	// 已绑定：消耗 state 后不再转发
	site, a, _ := setupAccountTest(t)
	startTestLogin(t, a)
	database.DB.Save(&models.ArkfgfAccount{ID: accountRowID, Token: testAccountToken})
	a.mu.Lock()
	a.loaded = false // 重新读取数据库，保留进行中的登录
	a.mu.Unlock()
	if err := a.FinishLogin(callbackParams(1)); !errors.Is(err, ErrAccountBound) || finishCalls(site) != 0 {
		t.Fatalf("err = %v calls = %v", err, site.callList())
	}
}

// prepareLocalRegistrations 已绑定并有本地的发布、联机房登记和节点容器
func prepareLocalRegistrations(t *testing.T, dock *fakeLinkDocker) models.Server {
	t.Helper()
	bindTestAccount(t)
	s := createTestServer(t, 7777, 27015, 32330)
	database.DB.Model(&s).Updates(map[string]any{"publish": true, "published": true, "link": true})
	database.DB.Save(&models.LinkRoom{ID: linkRoomID, RoomID: 12, NetworkName: "ark-x", NetworkSecret: "sec", RPCPort: 20000, Name: "My Room"})
	dock.addContainer(utils.LinkContainerName)
	return s
}

func expectLocalCleared(t *testing.T, a *AccountManager, dock *fakeLinkDocker, s models.Server) {
	t.Helper()
	if a.bound() {
		t.Fatal("account should be cleared")
	}
	var n int64
	database.DB.Model(&models.ArkfgfAccount{}).Count(&n)
	if n != 0 {
		t.Fatal("account row should be deleted")
	}
	if got := loadPublishServer(t, s.ID); got.Published || !got.Publish || !got.Link {
		t.Fatalf("server = %+v", got)
	}
	room := loadTestLinkRoom(t)
	if room.RoomID != 0 || room.NetworkName != "" || room.NetworkSecret != "" || room.Name != "My Room" || room.RPCPort != 20000 {
		t.Fatalf("room = %+v", room)
	}
	if dock.container(utils.LinkContainerName) != nil || dock.removedDirs == 0 {
		t.Fatal("node and credential dir should be removed")
	}
}

func TestAccountUnbind(t *testing.T) {
	site, a, dock := setupAccountTest(t)
	s := prepareLocalRegistrations(t, dock)
	if err := a.Unbind(false); err != nil {
		t.Fatal(err)
	}
	site.mu.Lock()
	if strings.Join(site.calls, ",") != "DELETE /commander/me" || site.tokens[0] != testAccountToken {
		t.Fatalf("calls = %v tokens = %v", site.calls, site.tokens)
	}
	site.mu.Unlock()
	expectLocalCleared(t, a, dock, s)
	// 未绑定时再解除绑定什么也不做
	if err := a.Unbind(false); err != nil {
		t.Fatal(err)
	}
}

// 网站已经不认识这个账号（invalid_account）：视同解除成功
func TestAccountUnbindInvalidAccount(t *testing.T) {
	site, a, dock := setupAccountTest(t)
	s := prepareLocalRegistrations(t, dock)
	site.set(func(f *fakeAccountSite) { f.deleteErr = "invalid_account" })
	if err := a.Unbind(false); err != nil {
		t.Fatal(err)
	}
	expectLocalCleared(t, a, dock, s)
}

// 连不上网站：不强制时报错、保留绑定；强制时只清本地
func TestAccountUnbindNetworkError(t *testing.T) {
	_, a, dock := setupAccountTest(t)
	s := prepareLocalRegistrations(t, dock)
	a.api = func() string { return "http://127.0.0.1:1" }
	if err := a.Unbind(false); !errors.Is(err, ErrArkfgfUnavailable) {
		t.Fatalf("err = %v", err)
	}
	if !a.bound() || loadTestLinkRoom(t).RoomID != 12 || !loadPublishServer(t, s.ID).Published {
		t.Fatal("nothing should be cleared without force")
	}
	if err := a.Unbind(true); err != nil {
		t.Fatal(err)
	}
	expectLocalCleared(t, a, dock, s)
}

// 刷新账号信息：更新昵称；网站返回 invalid_account 时清空本地
func TestAccountRefresh(t *testing.T) {
	site, a, dock := setupAccountTest(t)
	s := prepareLocalRegistrations(t, dock)
	st := a.Status(true)
	if st.Account == nil || st.Account.UserName != "New Nick" || st.Account.CommanderName != "renamed" {
		t.Fatalf("status = %+v", st)
	}
	site.mu.Lock()
	if site.calls[0] != "GET /commander/me" {
		t.Fatalf("calls = %v", site.calls)
	}
	site.mu.Unlock()

	database.DB.Model(&models.ArkfgfAccount{}).Where("id = ?", accountRowID).Update("token", "acm_removed")
	a.forget()
	if st := a.Status(true); st.Bound {
		t.Fatalf("status = %+v", st)
	}
	expectLocalCleared(t, a, dock, s)
}
