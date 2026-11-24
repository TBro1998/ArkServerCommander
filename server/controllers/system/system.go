package system

import (
	"net/http"

	"ark-server-commander/controllers/servers"
	"ark-server-commander/service/server"

	"github.com/gin-gonic/gin"
)

var serverService = server.NewServerService()

// GetOverview 首页概览
// @Summary 首页概览
// @Description 服务器数量、运行和可连接数量、在线人数、CPU/内存合计、各服务器磁盘占用、公网地址
// @Tags 系统
// @Produce json
// @Security Bearer
// @Param refresh query bool false "强制重新计算磁盘占用（较慢）"
// @Success 200 {object} map[string]server.Overview "概览"
// @Router /system/overview [get]
func GetOverview(c *gin.Context) {
	ov, err := serverService.GetOverview(c.GetUint("user_id"), c.Query("refresh") == "true")
	if err != nil {
		servers.RespondServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": ov})
}

// GetInfo 系统信息
// @Summary 系统信息
// @Description 玩家连接用的公网地址（PUBLIC_HOST 或自动探测）
// @Tags 系统
// @Produce json
// @Security Bearer
// @Success 200 {object} map[string]string "系统信息"
// @Router /system/info [get]
func GetInfo(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"data": gin.H{"public_host": server.PublicHost()}})
}

// GetVersion 版本信息
// @Summary 版本信息
// @Description 当前版本和 GitHub 上的最新发布版本（结果缓存 6 小时，查询失败 10 分钟后重试）。当前为开发版本（dev）时 has_update 始终为 false。self_update 是能否一键更新、更新进度和最近一次更新的结果
// @Tags 系统
// @Produce json
// @Security Bearer
// @Param refresh query bool false "重新查询最新版本"
// @Success 200 {object} map[string]server.VersionInfo "版本信息"
// @Router /system/version [get]
func GetVersion(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"data": serverService.GetVersionInfo(c.Query("refresh") == "true")})
}

// SelfUpdateRequest 更新管理端请求
type SelfUpdateRequest struct {
	Version string `json:"version" binding:"max=64"` // 界面上看到的最新版本，只用于记录，可以为空
}

// StartSelfUpdate 更新管理端
// @Summary 更新管理端
// @Description 下载管理端的新镜像，再由一个临时的更新容器用新镜像重建管理端容器（失败时恢复旧容器），立即返回。进度和结果看 GET /system/version 的 self_update。只支持用 tbro98/arkservercommander:latest 运行、挂载了 /var/run/docker.sock 的容器；有服务器正在启停、正在备份恢复或下载镜像时返回 409
// @Tags 系统
// @Accept json
// @Produce json
// @Security Bearer
// @Param body body SelfUpdateRequest false "最新版本"
// @Success 202 {object} map[string]string "已开始更新"
// @Failure 409 {object} map[string]string "不支持或正在进行"
// @Router /system/update [post]
func StartSelfUpdate(c *gin.Context) {
	var req SelfUpdateRequest
	if c.Request.ContentLength > 0 {
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
			return
		}
	}
	if err := serverService.StartSelfUpdate(req.Version); err != nil {
		servers.RespondServiceError(c, err)
		return
	}
	c.JSON(http.StatusAccepted, gin.H{"message": "已开始更新管理端"})
}

// DeleteCluster 删除未使用的集群数据
// @Summary 删除未使用的集群数据
// @Description 删除集群卷 ase-cluster-<id> 和数据目录 clusters/<id>（跨服传送的角色和物品）。仍有服务器设置该集群 ID，或卷仍被容器挂载时返回 409
// @Tags 系统
// @Produce json
// @Security Bearer
// @Param clusterId path string true "集群ID"
// @Success 200 {object} map[string]string "已删除"
// @Failure 409 {object} map[string]string "集群仍在使用中"
// @Router /system/clusters/{clusterId} [delete]
func DeleteCluster(c *gin.Context) {
	if err := serverService.DeleteCluster(c.Param("clusterId")); err != nil {
		servers.RespondServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "集群数据已删除"})
}

// SuggestPorts 推荐可用端口
// @Summary 推荐可用端口
// @Description 返回第一组与现有服务器不冲突的游戏端口、查询端口和 RCON 端口
// @Tags 系统
// @Produce json
// @Security Bearer
// @Success 200 {object} map[string]server.SuggestedPorts "端口"
// @Router /system/suggest-ports [get]
func SuggestPorts(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"data": serverService.SuggestPorts()})
}

// GetServerDefaults 新建服务器的默认配置
// @Summary 新建服务器的默认配置
// @Description 新建服务器时写入的默认 GameUserSettings.ini、Game.ini 和启动参数，新建页以它为基础编辑
// @Tags 系统
// @Produce json
// @Security Bearer
// @Success 200 {object} map[string]server.ServerDefaults "默认配置"
// @Router /system/server-defaults [get]
func GetServerDefaults(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"data": serverService.ServerDefaults()})
}

// GetMods 查询创意工坊 Mod 信息
// @Summary 查询创意工坊 Mod 信息
// @Description 根据 Mod ID 查询名称等信息（Steam 公开接口，结果缓存 6 小时）
// @Tags 系统
// @Produce json
// @Security Bearer
// @Param ids query string true "逗号分隔的 Mod ID"
// @Success 200 {object} map[string][]server.ModInfo "Mod 信息"
// @Router /mods [get]
func GetMods(c *gin.Context) {
	mods, err := serverService.GetModInfo(c.Query("ids"))
	if err != nil {
		servers.RespondServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": mods})
}
