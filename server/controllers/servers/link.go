package servers

import (
	"net/http"

	"ark-server-commander/models"

	"github.com/gin-gonic/gin"
)

// 联机房：每个管理端一个房间，所有开启了 link 的服务器都在这个房间里

// GetLinkRoom 联机房资料和状态
// @Summary 联机房资料和状态
// @Description 房间名称、收录资料（设置的原始值，为空时登记时用第一台服务器的）、状态和房间里的服务器
// @Tags 联机房
// @Produce json
// @Security Bearer
// @Success 200 {object} map[string]models.LinkRoomInfo "联机房"
// @Router /link/room [get]
func GetLinkRoom(c *gin.Context) {
	info, err := serverService.GetLinkRoom()
	if err != nil {
		respondServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": info})
}

// UpdateLinkRoom 修改联机房资料
// @Summary 修改联机房资料
// @Description 房间名称和收录资料（简介、标签、QQ 群、Discord），下一轮（15 秒内）重新登记时生效
// @Tags 联机房
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body models.LinkRoomSettings true "房间资料"
// @Success 200 {object} map[string]models.LinkRoomInfo "联机房"
// @Failure 400 {object} map[string]string "资料无效"
// @Router /link/room [put]
func UpdateLinkRoom(c *gin.Context) {
	var req models.LinkRoomSettings
	if err := c.ShouldBindJSON(&req); err != nil {
		respondBadRequest(c)
		return
	}
	info, err := serverService.UpdateLinkRoom(req)
	if err != nil {
		respondServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": info})
}

// GetLinkMembers 联机房成员和封禁列表
// @Summary 联机房成员
// @Description 代理 arkfgf.com 的 GET /link/rooms/:id/members（账号令牌认证）：当前持有未吊销凭证的玩家和封禁列表
// @Tags 联机房
// @Produce json
// @Security Bearer
// @Success 200 {object} map[string]models.LinkMembers "成员和封禁列表"
// @Failure 409 {object} map[string]string "没有绑定 arkfgf.com 账号，或联机房还没有登记"
// @Failure 502 {object} map[string]string "无法连接 arkfgf.com"
// @Router /link/members [get]
func GetLinkMembers(c *gin.Context) {
	members, err := serverService.GetLinkMembers()
	if err != nil {
		respondServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": members})
}

// KickLinkMember 踢出联机房成员
// @Summary 踢出联机房成员
// @Description 吊销玩家的凭证（下一次心跳时在节点里吊销，3 秒内断开，房间里所有服务器都连不上），ban 为 true 时同时封禁
// @Tags 联机房
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body models.LinkKickRequest true "玩家和是否封禁"
// @Success 200 {object} map[string]string "已踢出"
// @Failure 400 {object} map[string]string "SteamID 无效或网站拒绝"
// @Failure 409 {object} map[string]string "没有绑定 arkfgf.com 账号，或联机房还没有登记"
// @Failure 502 {object} map[string]string "无法连接 arkfgf.com"
// @Router /link/kick [post]
func KickLinkMember(c *gin.Context) {
	var req models.LinkKickRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondBadRequest(c)
		return
	}
	if err := serverService.KickLinkMember(req); err != nil {
		respondServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "已踢出"})
}

// UnbanLinkMember 解除联机房封禁
// @Summary 解除联机房封禁
// @Description 代理 arkfgf.com 的 DELETE /link/rooms/:id/bans/:steam_id
// @Tags 联机房
// @Produce json
// @Security Bearer
// @Param steam_id path string true "玩家的 SteamID64"
// @Success 200 {object} map[string]string "已解封"
// @Failure 400 {object} map[string]string "SteamID 无效或网站拒绝"
// @Failure 409 {object} map[string]string "没有绑定 arkfgf.com 账号，或联机房还没有登记"
// @Failure 502 {object} map[string]string "无法连接 arkfgf.com"
// @Router /link/bans/{steam_id} [delete]
func UnbanLinkMember(c *gin.Context) {
	if err := serverService.UnbanLinkMember(c.Param("steam_id")); err != nil {
		respondServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "已解封"})
}
