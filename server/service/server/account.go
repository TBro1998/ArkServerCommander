package server

import (
	"crypto/subtle"
	"fmt"
	"net/http"
	"net/url"
	"slices"
	"strings"
	"sync"
	"time"
	"unicode"
	"unicode/utf8"

	"ark-server-commander/config"
	"ark-server-commander/database"
	"ark-server-commander/models"
	"ark-server-commander/utils"

	"go.uber.org/zap"
)

// 绑定 arkfgf.com 账号（在 Commander 里直接用 Steam 登录）。接口约定见外层仓库 docs/commander-account.md。
//
// 发布到收录站和 P2P 联机房都要求先绑定。设置页点"用 Steam 登录"：前端把回调地址（本页面所在的
// origin + /settings/arkfgf-callback）和名称交给后端，后端转给 POST /commander/login 拿到 Steam 登录地址，
// 浏览器跳过去；Steam 登录后回到回调页，前端把 URL 上的全部参数交给后端，后端转给 POST /commander/login/finish
// （网站核验 state 和 Steam 签名），成功后保存账号令牌（acm_...）。Commander 不核验 Steam 签名，只转发。
// 之后 Publisher 和 LinkManager 用账号令牌调用网站。账号令牌只存在本机数据库，不返回给浏览器。
//
// 解除绑定（或任何请求返回 401 invalid_account，即服主在网站后台移除了这个 Commander）时清空本地账号和本地的
// 发布 / 联机房登记：servers.published、房间 ID、网络名、网络密钥、节点容器和凭证目录（房间资料和旧版令牌保留）。

const (
	accountRowID          = 1 // models.ArkfgfAccount 只有这一行
	accountDefaultName    = "ArkServerCommander"
	accountNameMax        = 64
	accountReturnToMax    = 256
	accountParamsMax      = 64 // 回调参数最多几个（Steam OpenID 的参数十来个）
	accountParamMax       = 2048
	accountRequestTimeout = 15 * time.Second // 网站在 finish 里同步向 Steam 核验，留足余量
	accountLoginTTL       = 10 * time.Minute // 和网站上 state 的有效期一致
	accountPendingMax     = 5
)

// arkfgfMu 串行化发布和联机房的每一轮处理，以及账号的保存和清除：
// 清除本地登记时不会有正在进行的一轮把刚清掉的状态（published、房间 ID、节点容器）又写回去
var arkfgfMu sync.Mutex

// pendingLogin 已开始、还没完成的 Steam 登录：完成登录的接口是公开的（Steam 客户端内置浏览器里没有管理端的登录 cookie），
// 只接受这里记下的 state，用一次就删
type pendingLogin struct {
	state     string
	startedAt time.Time
	expiresAt time.Time
}

// AccountManager 绑定的 arkfgf.com 账号，缓存在内存里
type AccountManager struct {
	mu      sync.Mutex
	loaded  bool
	account models.ArkfgfAccount // Token 为空表示未绑定

	// 进行中的登录（最多 accountPendingMax 个）和最近一次完成登录失败的原因
	logins   []pendingLogin
	loginErr string
	loginMsg string
	// 已取走 state、正在向网站核验的登录数。核验期间（网站要向 Steam 核验，可能要几秒）也算登录中，
	// 否则设置页轮询时会在这几秒里看到"没有进行中的登录、也没绑定"，误报登录过期并停止等待
	finishing int

	client *http.Client
	api    func() string
	now    func() time.Time
}

var accountManager = &AccountManager{
	client: &http.Client{Timeout: accountRequestTimeout},
	api:    func() string { return config.PublishAPI },
	now:    time.Now,
}

// loadLocked 第一次使用时从数据库读取账号，调用时持有 a.mu
func (a *AccountManager) loadLocked() {
	if a.loaded {
		return
	}
	var account models.ArkfgfAccount
	if err := database.DB.Where("id = ?", accountRowID).Limit(1).Find(&account).Error; err != nil {
		utils.Warn("读取 arkfgf.com 账号失败", zap.Error(err))
		return
	}
	a.account = account
	a.loaded = true
}

// forget 丢弃内存里的账号，下次使用时重新读取数据库（测试换数据库时用）
func (a *AccountManager) forget() {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.loaded = false
	a.account = models.ArkfgfAccount{}
	a.logins = nil
	a.loginErr, a.loginMsg = "", ""
	a.finishing = 0
}

// pruneLoginsLocked 去掉过期的登录，调用时持有 a.mu
func (a *AccountManager) pruneLoginsLocked() {
	now := a.now()
	a.logins = slices.DeleteFunc(a.logins, func(l pendingLogin) bool { return !now.Before(l.expiresAt) })
}

// addLogin 记下一次开始的登录，超过上限时丢掉最早的；清除上次失败的原因
func (a *AccountManager) addLogin(state string) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.pruneLoginsLocked()
	now := a.now()
	a.logins = append(a.logins, pendingLogin{state: state, startedAt: now, expiresAt: now.Add(accountLoginTTL)})
	if len(a.logins) > accountPendingMax {
		a.logins = a.logins[len(a.logins)-accountPendingMax:]
	}
	a.loginErr, a.loginMsg = "", ""
}

// consumeLogin 找到并删除和 state 一致、还没过期的登录（常数时间比较），找不到返回 false
func (a *AccountManager) consumeLogin(state string) bool {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.pruneLoginsLocked()
	found := -1
	for i, l := range a.logins {
		if subtle.ConstantTimeCompare([]byte(l.state), []byte(state)) == 1 {
			found = i
		}
	}
	if found < 0 {
		return false
	}
	a.logins = slices.Delete(a.logins, found, found+1)
	a.finishing++
	return true
}

// finishDone 一次核验结束（成功或失败），和 consumeLogin 成对调用
func (a *AccountManager) finishDone() {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.finishing > 0 {
		a.finishing--
	}
}

// loginFailed 记下完成登录失败的原因，设置页（在普通浏览器里）轮询状态时显示
func (a *AccountManager) loginFailed(err error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.loginErr, a.loginMsg = ErrorCode(err), err.Error()
	if a.loginErr == "" {
		a.loginErr = "login_failed"
	}
}

// loginStatusLocked 登录状态，调用时持有 a.mu
func (a *AccountManager) loginStatusLocked() models.ArkfgfLoginStatus {
	a.pruneLoginsLocked()
	st := models.ArkfgfLoginStatus{Pending: len(a.logins) > 0 || a.finishing > 0, ErrorCode: a.loginErr, Error: a.loginMsg}
	if len(a.logins) > 0 {
		st.ExpiresAt = a.logins[len(a.logins)-1].expiresAt.Format(time.RFC3339)
	}
	return st
}

// stateFromSteamURL 从网站返回的 Steam 登录地址里取出 state（openid.return_to 的 query 参数）
func stateFromSteamURL(steamURL string) string {
	u, err := url.Parse(steamURL)
	if err != nil {
		return ""
	}
	ret, err := url.Parse(u.Query().Get("openid.return_to"))
	if err != nil {
		return ""
	}
	return ret.Query().Get("state")
}

// token 账号令牌，未绑定时为空
func (a *AccountManager) token() string {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.loadLocked()
	return a.account.Token
}

// bound 是否已绑定账号
func (a *AccountManager) bound() bool {
	return a.token() != ""
}

// Status 账号的绑定状态。refresh 为 true 且已绑定时先向网站查询账号信息（GET /commander/me），
// 网站返回 invalid_account 时视为已解绑；连不上网站时返回本地保存的信息
func (a *AccountManager) Status(refresh bool) models.ArkfgfAccountStatus {
	if refresh && a.api() != "" {
		if token := a.token(); token != "" {
			a.refresh(token)
		}
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	a.loadLocked()
	status := models.ArkfgfAccountStatus{
		Enabled: a.api() != "",
		Bound:   a.account.Token != "",
		Login:   a.loginStatusLocked(),
	}
	if status.Bound {
		acc := a.account
		status.DashboardURL = acc.DashboardURL
		status.Account = &models.ArkfgfAccountInfo{
			CommanderID:   acc.CommanderID,
			CommanderName: acc.CommanderName,
			SteamID:       acc.SteamID,
			UserName:      acc.UserName,
			AvatarURL:     acc.AvatarURL,
		}
		if !acc.BoundAt.IsZero() {
			status.Account.BoundAt = acc.BoundAt.Format(time.RFC3339)
		}
	}
	return status
}

// accountMeResponse GET /commander/me 和登录成功时的账号信息
type accountMeResponse struct {
	Commander struct {
		ID   uint   `json:"id"`
		Name string `json:"name"`
	} `json:"commander"`
	User struct {
		SteamID   string `json:"steam_id"`
		Name      string `json:"name"`
		AvatarURL string `json:"avatar_url"`
	} `json:"user"`
	DashboardURL string `json:"dashboard_url"`
}

// refresh 向网站查询账号信息并更新本地保存的名称和头像
func (a *AccountManager) refresh(token string) {
	var resp accountMeResponse
	status, _, code, err := a.request(http.MethodGet, "/commander/me", token, nil, &resp)
	if err != nil {
		if isInvalidAccount(status, code) {
			a.invalidate()
		}
		return
	}
	arkfgfMu.Lock()
	defer arkfgfMu.Unlock()
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.account.Token != token {
		return
	}
	updates := map[string]any{}
	set := func(column string, field *string, value string) {
		if value != "" && *field != value {
			*field = value
			updates[column] = value
		}
	}
	set("commander_name", &a.account.CommanderName, resp.Commander.Name)
	set("steam_id", &a.account.SteamID, resp.User.SteamID)
	set("user_name", &a.account.UserName, resp.User.Name)
	set("avatar_url", &a.account.AvatarURL, resp.User.AvatarURL)
	set("dashboard_url", &a.account.DashboardURL, resp.DashboardURL)
	if resp.Commander.ID != 0 && resp.Commander.ID != a.account.CommanderID {
		a.account.CommanderID = resp.Commander.ID
		updates["commander_id"] = resp.Commander.ID
	}
	if len(updates) > 0 {
		if err := database.DB.Model(&models.ArkfgfAccount{}).Where("id = ?", accountRowID).Updates(updates).Error; err != nil {
			utils.Warn("保存 arkfgf.com 账号信息失败", zap.Error(err))
		}
	}
}

// normalizeCommanderName 去掉控制字符和首尾空白，为空时用默认名称
func normalizeCommanderName(name string) (string, error) {
	name = strings.TrimSpace(strings.Map(func(r rune) rune {
		if unicode.IsControl(r) {
			return -1
		}
		return r
	}, name))
	if name == "" {
		return accountDefaultName, nil
	}
	if utf8.RuneCountInString(name) > accountNameMax {
		return "", ErrInvalidCommanderName
	}
	return name, nil
}

// validateReturnTo 回调地址：http(s)、有 host、不带 query 和 fragment、最多 256 字（和网站的要求一致）
func validateReturnTo(returnTo string) error {
	if returnTo == "" || len(returnTo) > accountReturnToMax || strings.ContainsAny(returnTo, "?# \t\r\n") {
		return ErrInvalidReturnTo
	}
	u, err := url.Parse(returnTo)
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" || u.User != nil {
		return ErrInvalidReturnTo
	}
	return nil
}

// StartLogin 开始用 Steam 登录绑定：向网站申请登录（POST /commander/login），返回浏览器要跳转的 Steam 登录地址
func (a *AccountManager) StartLogin(name, returnTo string) (string, error) {
	if a.api() == "" {
		return "", ErrArkfgfDisabled
	}
	name, err := normalizeCommanderName(name)
	if err != nil {
		return "", err
	}
	returnTo = strings.TrimSpace(returnTo)
	if err := validateReturnTo(returnTo); err != nil {
		return "", err
	}
	if a.bound() {
		return "", ErrAccountBound
	}
	var resp struct {
		SteamURL string `json:"steam_url"`
	}
	body := map[string]string{"return_to": returnTo, "name": name, "version": utils.Version}
	status, _, code, err := a.request(http.MethodPost, "/commander/login", "", body, &resp)
	if err != nil {
		return "", loginError(status, code, err)
	}
	if !strings.HasPrefix(resp.SteamURL, "https://") {
		return "", fmt.Errorf("%w：网站没有返回 Steam 登录地址", ErrArkfgfUnavailable)
	}
	state := stateFromSteamURL(resp.SteamURL)
	if state == "" {
		return "", fmt.Errorf("%w：网站返回的 Steam 登录地址里没有 state", ErrArkfgfUnavailable)
	}
	a.addLogin(state)
	return resp.SteamURL, nil
}

// FinishLogin Steam 登录后回到回调页：把回调地址上的全部参数转给网站（POST /commander/login/finish），
// 网站核验通过后保存账号令牌，让发布和联机房马上开始。
// 回调页可能在 Steam 客户端的内置浏览器里打开（没有管理端的登录 cookie），所以这一步不要求登录，
// 只接受本机开始的、还没过期也没用过的 state（用一次就删），其他请求一律 expired_login 且不联系网站
func (a *AccountManager) FinishLogin(params map[string]string) error {
	if a.api() == "" {
		return ErrArkfgfDisabled
	}
	if params["state"] == "" || len(params) > accountParamsMax {
		return ErrLoginExpired
	}
	for k, v := range params {
		if len(k) > accountParamMax || len(v) > accountParamMax {
			return ErrLoginExpired
		}
	}
	if !a.consumeLogin(params["state"]) {
		return ErrLoginExpired
	}
	// 失败原因要在核验结束（finishDone）之前记下：设置页看到"不再登录中"时，要么已绑定，要么能看到失败原因
	defer a.finishDone()
	if err := a.finish(params); err != nil {
		a.loginFailed(err)
		return err
	}
	return nil
}

// finish 把回调参数转给网站，成功后保存账号
func (a *AccountManager) finish(params map[string]string) error {
	if a.bound() {
		return ErrAccountBound
	}
	var resp struct {
		Token string `json:"token"`
		accountMeResponse
	}
	status, _, code, err := a.request(http.MethodPost, "/commander/login/finish", "", map[string]any{"params": params}, &resp)
	if err != nil {
		return loginError(status, code, err)
	}
	if !strings.HasPrefix(resp.Token, "acm_") {
		return fmt.Errorf("%w：网站没有返回账号令牌", ErrArkfgfUnavailable)
	}
	if err := a.save(resp.Token, resp.accountMeResponse); err != nil {
		return err
	}
	// 已绑定：其他进行中的登录不再需要
	a.mu.Lock()
	a.logins = nil
	a.loginErr, a.loginMsg = "", ""
	a.mu.Unlock()
	return nil
}

// save 保存登录得到的账号，清掉未绑定期间的发布和联机房状态，马上处理一轮
func (a *AccountManager) save(token string, info accountMeResponse) error {
	arkfgfMu.Lock()
	if a.bound() {
		// 同时完成了两次登录：保留先保存的那个，网站上多出的 Commander 由服主在后台移除
		arkfgfMu.Unlock()
		return ErrAccountBound
	}
	account := models.ArkfgfAccount{
		ID:            accountRowID,
		Token:         token,
		CommanderID:   info.Commander.ID,
		CommanderName: info.Commander.Name,
		SteamID:       info.User.SteamID,
		UserName:      info.User.Name,
		AvatarURL:     info.User.AvatarURL,
		DashboardURL:  info.DashboardURL,
		BoundAt:       a.now(),
	}
	if err := database.DB.Save(&account).Error; err != nil {
		arkfgfMu.Unlock()
		return fmt.Errorf("保存 arkfgf.com 账号失败: %w", err)
	}
	a.mu.Lock()
	a.account = account
	a.loaded = true
	a.mu.Unlock()
	publisher.reset()
	linkManager.reset()
	arkfgfMu.Unlock()
	publisher.wake()
	linkManager.wake()
	utils.Info("已绑定 arkfgf.com 账号", zap.Uint("commander_id", account.CommanderID), zap.String("user", account.UserName))
	return nil
}

// Unbind 解除绑定：DELETE /commander/me（网站删除这个 Commander 发布的服务器和联机房），
// 成功或网站已经不认识这个账号（invalid_account）后清空本地。连不上网站时 force 为 true 才只清本地
func (a *AccountManager) Unbind(force bool) error {
	token := a.token()
	if token == "" {
		return nil
	}
	if a.api() == "" {
		if !force {
			return fmt.Errorf("%w：%v", ErrArkfgfUnavailable, ErrArkfgfDisabled)
		}
	} else {
		status, _, code, err := a.request(http.MethodDelete, "/commander/me", token, nil, nil)
		if err != nil && !isInvalidAccount(status, code) {
			if !force {
				return arkfgfCallError(status, code, err)
			}
			utils.Warn("解除 arkfgf.com 绑定时无法通知网站，只清除本地", zap.Error(err))
		}
	}
	arkfgfMu.Lock()
	defer arkfgfMu.Unlock()
	return a.clearLocked("已解除 arkfgf.com 账号绑定")
}

// invalidate 网站不认识这个账号了（服主在网站后台移除了这个 Commander）：清空本地
func (a *AccountManager) invalidate() {
	arkfgfMu.Lock()
	defer arkfgfMu.Unlock()
	a.invalidateLocked()
}

// invalidateLocked 同 invalidate，调用时持有 arkfgfMu（发布和联机房的一轮处理中）
func (a *AccountManager) invalidateLocked() {
	if a.token() == "" {
		return
	}
	if err := a.clearLocked("arkfgf.com 返回账号无效（可能已在网站上移除了这个 Commander），已解除绑定"); err != nil {
		utils.Warn("清除 arkfgf.com 账号失败", zap.Error(err))
	}
}

// clearLocked 清空本地账号和本地的发布 / 联机房登记，调用时持有 arkfgfMu
func (a *AccountManager) clearLocked(reason string) error {
	if err := database.DB.Where("id = ?", accountRowID).Delete(&models.ArkfgfAccount{}).Error; err != nil {
		return fmt.Errorf("清除账号失败: %w", err)
	}
	a.mu.Lock()
	a.account = models.ArkfgfAccount{}
	a.loaded = true
	a.mu.Unlock()

	if err := database.DB.Unscoped().Model(&models.Server{}).Where("published = ?", true).
		Update("published", false).Error; err != nil {
		utils.Warn("清除发布记录失败", zap.Error(err))
	}
	// 房间资料和旧版令牌保留：重新绑定后登记时旧令牌还能让网站删除旧房间
	if err := database.DB.Model(&models.LinkRoom{}).Where("id = ?", linkRoomID).
		Updates(map[string]any{"room_id": 0, "network_name": "", "network_secret": ""}).Error; err != nil {
		utils.Warn("清除联机房登记失败", zap.Error(err))
	}
	linkManager.resetNode()
	publisher.reset()
	linkManager.reset()
	statsMu.Lock()
	clear(statsCache)
	statsMu.Unlock()
	utils.Info(reason)
	return nil
}

func (a *AccountManager) request(method, path, token string, body, out any) (int, time.Duration, string, error) {
	return arkfgfRequest(a.client, method, a.api()+path, token, body, out, " arkfgf.com ")
}

// isInvalidAccount 网站不认识这个账号令牌（被移除或已解绑），401 invalid_account
func isInvalidAccount(status int, code string) bool {
	return status == http.StatusUnauthorized && code == "invalid_account"
}

// loginError 登录请求的失败：网站的登录错误码对应各自的服务层错误，其他同 arkfgfCallError
func loginError(status int, code string, err error) error {
	switch code {
	case "expired_login":
		return ErrLoginExpired
	case "login_cancelled":
		return ErrLoginCancelled
	case "steam_verify_failed":
		return ErrSteamVerifyFailed
	case "too_many_commanders":
		return ErrTooManyCommanders
	case "invalid_input":
		return fmt.Errorf("%w：%v", ErrArkfgfRejected, err)
	}
	return arkfgfCallError(status, code, err)
}

// arkfgfCallError 把网站请求的失败转换为服务层错误：连不上、限流和 5xx 为 ErrArkfgfUnavailable，其他为 ErrArkfgfRejected
func arkfgfCallError(status int, code string, err error) error {
	if status == 0 || status == http.StatusTooManyRequests || status >= 500 {
		return fmt.Errorf("%w：%v", ErrArkfgfUnavailable, err)
	}
	if code != "" {
		return fmt.Errorf("%w（%s）：%v", ErrArkfgfRejected, code, err)
	}
	return fmt.Errorf("%w：%v", ErrArkfgfRejected, err)
}

// ---- 服务层接口 ----

// GetArkfgfAccount arkfgf.com 账号的绑定状态
func (s *ServerService) GetArkfgfAccount(refresh bool) models.ArkfgfAccountStatus {
	return accountManager.Status(refresh)
}

// StartArkfgfLogin 开始用 Steam 登录绑定，返回 Steam 登录地址
func (s *ServerService) StartArkfgfLogin(name, returnTo string) (string, error) {
	return accountManager.StartLogin(name, returnTo)
}

// FinishArkfgfLogin 用 Steam 回调参数完成绑定（公开接口，只接受本机开始的登录的 state）
func (s *ServerService) FinishArkfgfLogin(params map[string]string) error {
	return accountManager.FinishLogin(params)
}

// UnbindArkfgfAccount 解除绑定。force 为 true 时连不上网站也清除本地
func (s *ServerService) UnbindArkfgfAccount(force bool) error {
	return accountManager.Unbind(force)
}
