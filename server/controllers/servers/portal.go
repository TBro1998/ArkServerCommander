package servers

import (
	"net/http"

	"ark-server-commander/models"

	"github.com/gin-gonic/gin"
)

// 玩家站：每个管理端一个，玩家用浏览器或联机助手查看公告、服务器状态和商店

// GetPortal 玩家站设置和状态
// @Summary 玩家站设置和状态
// @Description 设置、玩家站容器的状态和访问地址、每台服务器是否显示、有没有商店
// @Tags 玩家站
// @Produce json
// @Security Bearer
// @Success 200 {object} map[string]models.PortalInfo "玩家站"
// @Router /portal [get]
func GetPortal(c *gin.Context) {
	info, err := serverService.GetPortal()
	if err != nil {
		respondServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": info})
}

// UpdatePortal 修改玩家站设置
// @Summary 修改玩家站设置
// @Description 开关、端口、访问地址、站点资料、隐藏的服务器，整体替换，马上生效
// @Tags 玩家站
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body models.PortalSettingsRequest true "设置"
// @Success 200 {object} map[string]models.PortalInfo "玩家站"
// @Failure 400 {object} map[string]string "设置无效"
// @Failure 409 {object} map[string]string "端口冲突"
// @Router /portal [put]
func UpdatePortal(c *gin.Context) {
	var req models.PortalSettingsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondBadRequest(c)
		return
	}
	info, err := serverService.UpdatePortal(req)
	if err != nil {
		respondServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": info})
}

// ListAnnouncements 玩家站公告
// @Summary 玩家站公告
// @Description 从新到旧
// @Tags 玩家站
// @Produce json
// @Security Bearer
// @Success 200 {object} map[string][]models.PortalAnnouncement "公告"
// @Router /portal/announcements [get]
func ListAnnouncements(c *gin.Context) {
	list, err := serverService.ListAnnouncements()
	if err != nil {
		respondServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": list})
}

// CreateAnnouncement 发布公告
// @Summary 发布公告
// @Tags 玩家站
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body models.PortalAnnouncementRequest true "公告"
// @Success 200 {object} map[string]models.PortalAnnouncement "公告"
// @Failure 400 {object} map[string]string "公告无效"
// @Router /portal/announcements [post]
func CreateAnnouncement(c *gin.Context) {
	var req models.PortalAnnouncementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondBadRequest(c)
		return
	}
	a, err := serverService.CreateAnnouncement(req)
	if err != nil {
		respondServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": a})
}

// UpdateAnnouncement 修改公告
// @Summary 修改公告
// @Tags 玩家站
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "公告 ID"
// @Param request body models.PortalAnnouncementRequest true "公告"
// @Success 200 {object} map[string]models.PortalAnnouncement "公告"
// @Failure 400 {object} map[string]string "公告无效"
// @Failure 404 {object} map[string]string "公告不存在"
// @Router /portal/announcements/{id} [put]
func UpdateAnnouncement(c *gin.Context) {
	var req models.PortalAnnouncementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondBadRequest(c)
		return
	}
	a, err := serverService.UpdateAnnouncement(c.Param("id"), req)
	if err != nil {
		respondServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": a})
}

// DeleteAnnouncement 删除公告
// @Summary 删除公告
// @Tags 玩家站
// @Produce json
// @Security Bearer
// @Param id path int true "公告 ID"
// @Success 200 {object} map[string]string "已删除"
// @Failure 404 {object} map[string]string "公告不存在"
// @Router /portal/announcements/{id} [delete]
func DeleteAnnouncement(c *gin.Context) {
	if err := serverService.DeleteAnnouncement(c.Param("id")); err != nil {
		respondServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "已删除"})
}
