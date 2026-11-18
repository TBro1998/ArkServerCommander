package servers

import (
	"net/http"

	"ark-server-commander/models"

	"github.com/gin-gonic/gin"
)

// GetLinkMembers 联机房成员和封禁列表
// @Summary 联机房成员
// @Description 代理 arkfgf.com 的 GET /link/rooms/:id/members（房间令牌认证）：当前持有未吊销凭证的玩家和封禁列表
// @Tags 服务器管理
// @Produce json
// @Security Bearer
// @Param id path int true "服务器ID"
// @Success 200 {object} map[string]models.LinkMembers "成员和封禁列表"
// @Failure 409 {object} map[string]string "联机房还没有登记"
// @Failure 502 {object} map[string]string "无法连接 arkfgf.com"
// @Router /servers/{id}/link/members [get]
func GetLinkMembers(c *gin.Context) {
	members, err := serverService.GetLinkMembers(c.GetUint("user_id"), c.Param("id"))
	if err != nil {
		respondServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": members})
}

// KickLinkMember 踢出联机房成员
// @Summary 踢出联机房成员
// @Description 吊销玩家的凭证（下一次心跳时在节点里吊销，3 秒内断开），ban 为 true 时同时封禁
// @Tags 服务器管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "服务器ID"
// @Param request body models.LinkKickRequest true "玩家和是否封禁"
// @Success 200 {object} map[string]string "已踢出"
// @Failure 400 {object} map[string]string "SteamID 无效或网站拒绝"
// @Failure 409 {object} map[string]string "联机房还没有登记"
// @Failure 502 {object} map[string]string "无法连接 arkfgf.com"
// @Router /servers/{id}/link/kick [post]
func KickLinkMember(c *gin.Context) {
	var req models.LinkKickRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondBadRequest(c)
		return
	}
	if err := serverService.KickLinkMember(c.GetUint("user_id"), c.Param("id"), req); err != nil {
		respondServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "已踢出"})
}

// UnbanLinkMember 解除联机房封禁
// @Summary 解除联机房封禁
// @Description 代理 arkfgf.com 的 DELETE /link/rooms/:id/bans/:steam_id
// @Tags 服务器管理
// @Produce json
// @Security Bearer
// @Param id path int true "服务器ID"
// @Param steam_id path string true "玩家的 SteamID64"
// @Success 200 {object} map[string]string "已解封"
// @Failure 400 {object} map[string]string "SteamID 无效或网站拒绝"
// @Failure 409 {object} map[string]string "联机房还没有登记"
// @Failure 502 {object} map[string]string "无法连接 arkfgf.com"
// @Router /servers/{id}/link/bans/{steam_id} [delete]
func UnbanLinkMember(c *gin.Context) {
	if err := serverService.UnbanLinkMember(c.GetUint("user_id"), c.Param("id"), c.Param("steam_id")); err != nil {
		respondServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "已解封"})
}
