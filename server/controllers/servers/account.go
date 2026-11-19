package servers

import (
	"net/http"

	"ark-server-commander/models"

	"github.com/gin-gonic/gin"
)

// arkfgf.com 账号绑定（在 Commander 里直接用 Steam 登录）：发布到收录站和 P2P 联机房都需要先绑定。
// 账号令牌只存在后端，不返回给浏览器

// GetArkfgfAccount arkfgf.com 账号的绑定状态
// @Summary arkfgf.com 账号
// @Description 是否启用（ARKFGF_API 不为空）、是否已绑定、账号信息（Steam 昵称和头像、Commander 在网站上的名称，不含令牌）、服主后台地址，以及进行中的 Steam 登录（login：是否在等待、过期时间、最近一次完成登录失败的原因）。refresh=true 且已绑定时先向网站查询最新的账号信息，网站返回账号无效时视为已解绑
// @Tags arkfgf.com 账号
// @Produce json
// @Security Bearer
// @Param refresh query bool false "向网站查询最新的账号信息"
// @Success 200 {object} map[string]models.ArkfgfAccountStatus "账号状态"
// @Router /arkfgf/account [get]
func GetArkfgfAccount(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"data": serverService.GetArkfgfAccount(c.Query("refresh") == "true")})
}

// StartArkfgfLogin 开始用 Steam 登录绑定
// @Summary 用 Steam 登录绑定 arkfgf.com 账号
// @Description 转给网站的 POST /commander/login，返回 Steam 登录地址，记下其中的 state（10 分钟有效，最多同时 5 个）。前端用 steam://openurl/ 在 Steam 客户端里打开它（或在浏览器里跳转）；登录后 Steam 跳回 return_to（Commander 的 /arkfgf-callback 页面），由该页面调用 /arkfgf/account/login/finish 完成绑定
// @Tags arkfgf.com 账号
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body models.ArkfgfLoginRequest true "Commander 名称和回调地址"
// @Success 200 {object} map[string]models.ArkfgfLoginResponse "Steam 登录地址"
// @Failure 400 {object} map[string]string "名称或回调地址无效，或网站拒绝"
// @Failure 409 {object} map[string]string "已经绑定，或 ARKFGF_API 为空"
// @Failure 502 {object} map[string]string "无法连接 arkfgf.com"
// @Router /arkfgf/account/login [post]
func StartArkfgfLogin(c *gin.Context) {
	var req models.ArkfgfLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondBadRequest(c)
		return
	}
	steamURL, err := serverService.StartArkfgfLogin(req.Name, req.ReturnTo)
	if err != nil {
		respondServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": models.ArkfgfLoginResponse{SteamURL: steamURL}})
}

// FinishArkfgfLogin 完成 Steam 登录绑定
// @Summary 完成 Steam 登录绑定（公开）
// @Description 不需要登录：回调页可能在 Steam 客户端的内置浏览器里打开，那里没有管理端的登录 cookie。只接受本机 /arkfgf/account/login 开始的、10 分钟内没用过的 state（用一次就删），否则 400 expired_login 且不联系网站。通过后把 Steam 跳回回调页时 URL 上的全部 query 参数（state、openid.*）转给网站的 POST /commander/login/finish，网站核验通过后保存账号令牌，发布和联机房马上开始。失败：expired_login（过期或已用过）、login_cancelled（在 Steam 取消）、steam_verify_failed（核验失败）为 400，too_many_commanders（账号绑定数已满）为 409；失败原因也会出现在 GET /arkfgf/account 的 login 里
// @Tags arkfgf.com 账号
// @Accept json
// @Produce json
// @Param request body models.ArkfgfLoginFinishRequest true "回调参数"
// @Success 200 {object} map[string]string "已绑定"
// @Failure 400 {object} map[string]string "登录过期、取消或核验失败"
// @Failure 409 {object} map[string]string "已经绑定、账号绑定数已满，或 ARKFGF_API 为空"
// @Failure 502 {object} map[string]string "无法连接 arkfgf.com"
// @Router /arkfgf/account/login/finish [post]
func FinishArkfgfLogin(c *gin.Context) {
	// 公开接口，限制请求体大小（Steam 回调参数只有一两 KB）
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 64<<10)
	var req models.ArkfgfLoginFinishRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondBadRequest(c)
		return
	}
	if err := serverService.FinishArkfgfLogin(req.Params); err != nil {
		respondServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "已绑定"})
}

// UnbindArkfgfAccount 解除绑定
// @Summary 解除 arkfgf.com 账号绑定
// @Description 调用网站的 DELETE /commander/me（删除这个 Commander 发布的服务器和联机房），成功后清空本地账号和本地的发布、联机房登记（房间资料保留）。连不上网站时返回 502，force=true 时只清本地（网站上的记录由服主在服主后台移除）
// @Tags arkfgf.com 账号
// @Produce json
// @Security Bearer
// @Param force query bool false "连不上网站时也清除本地"
// @Success 200 {object} map[string]string "已解除绑定"
// @Failure 400 {object} map[string]string "网站拒绝"
// @Failure 502 {object} map[string]string "无法连接 arkfgf.com"
// @Router /arkfgf/account [delete]
func UnbindArkfgfAccount(c *gin.Context) {
	if err := serverService.UnbindArkfgfAccount(c.Query("force") == "true"); err != nil {
		respondServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "已解除绑定"})
}
