package servers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// maxPluginUpload 插件压缩包上传大小上限
const maxPluginUpload = 200 << 20

// ListPlugins 插件列表
// @Summary ArkApi 插件列表
// @Description 列出插件卷里的插件和已停用的插件（PluginInfo.json 中的名称、版本等），以及 ArkApi 版本
// @Tags 插件管理
// @Produce json
// @Security Bearer
// @Param id path int true "服务器ID"
// @Success 200 {object} map[string]server.PluginsOverview "插件列表"
// @Router /servers/{id}/plugins [get]
func ListPlugins(c *gin.Context) {
	overview, err := serverService.ListPlugins(c.GetUint("user_id"), c.Param("id"))
	if err != nil {
		respondServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": overview})
}

// InstallPlugin 上传安装插件
// @Summary 上传安装插件
// @Description 上传插件 zip（multipart 字段 file）。压缩包里每个 <名称>/<名称>.dll 所在目录安装为一个插件；识别不出时用 name 指定插件名称。覆盖安装时保留原 config.json。只能在服务器停止时操作
// @Tags 插件管理
// @Accept multipart/form-data
// @Produce json
// @Security Bearer
// @Param id path int true "服务器ID"
// @Param file formData file true "插件 zip"
// @Param name formData string false "插件名称（压缩包结构识别不出时填写）"
// @Success 200 {object} map[string][]string "安装的插件名称"
// @Failure 400 {object} map[string]string "压缩包无效或服务器未停止"
// @Router /servers/{id}/plugins [post]
func InstallPlugin(c *gin.Context) {
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxPluginUpload)
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请上传不超过 200 MB 的 zip 文件", "code": "bad_request"})
		return
	}
	f, err := file.Open()
	if err != nil {
		respondBadRequest(c)
		return
	}
	defer f.Close()
	names, err := serverService.InstallPlugins(c.GetUint("user_id"), c.Param("id"), f, file.Size, c.PostForm("name"))
	if err != nil {
		respondServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "插件已安装", "data": names})
}

// InstallPluginFromURLRequest 从下载地址安装插件
type InstallPluginFromURLRequest struct {
	URL  string `json:"url" binding:"required"`
	Name string `json:"name"`
}

// InstallPluginFromURL 从下载地址安装插件
// @Summary 从下载地址安装插件
// @Description 后端下载 zip（http/https，最多 200 MB）后按上传安装的规则安装；识别不出插件目录时用 name 指定插件名称。覆盖安装时保留原 config.json。只能在服务器停止时操作
// @Tags 插件管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "服务器ID"
// @Param body body InstallPluginFromURLRequest true "下载地址和可选的插件名称"
// @Success 200 {object} map[string][]string "安装的插件名称"
// @Failure 400 {object} map[string]string "地址无效、压缩包无效或服务器未停止"
// @Failure 502 {object} map[string]string "下载失败"
// @Router /servers/{id}/plugins/fetch [post]
func InstallPluginFromURL(c *gin.Context) {
	var req InstallPluginFromURLRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondBadRequest(c)
		return
	}
	names, err := serverService.InstallPluginsFromURL(c.GetUint("user_id"), c.Param("id"), req.URL, req.Name)
	if err != nil {
		respondServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "插件已安装", "data": names})
}

// SetPluginStateRequest 启用或停用插件
type SetPluginStateRequest struct {
	Enabled bool `json:"enabled"`
}

// SetPluginState 启用 / 停用插件
// @Summary 启用或停用插件
// @Description 停用的插件移出插件卷，ArkApi 不再加载。只能在服务器停止时操作
// @Tags 插件管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "服务器ID"
// @Param name path string true "插件名称"
// @Param body body SetPluginStateRequest true "是否启用"
// @Success 200 {object} map[string]string "已更新"
// @Router /servers/{id}/plugins/{name} [put]
func SetPluginState(c *gin.Context) {
	var req SetPluginStateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondBadRequest(c)
		return
	}
	if err := serverService.SetPluginEnabled(c.GetUint("user_id"), c.Param("id"), c.Param("name"), req.Enabled); err != nil {
		respondServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "插件状态已更新"})
}

// DeletePlugin 删除插件
// @Summary 删除插件
// @Description 删除插件目录（包括配置）。只能在服务器停止时操作
// @Tags 插件管理
// @Produce json
// @Security Bearer
// @Param id path int true "服务器ID"
// @Param name path string true "插件名称"
// @Success 200 {object} map[string]string "已删除"
// @Router /servers/{id}/plugins/{name} [delete]
func DeletePlugin(c *gin.Context) {
	if err := serverService.DeletePlugin(c.GetUint("user_id"), c.Param("id"), c.Param("name")); err != nil {
		respondServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "插件已删除"})
}

// PluginConfigRequest 配置文件内容
type PluginConfigRequest struct {
	Content string `json:"content"`
}

// GetPluginConfig 读取插件配置
// @Summary 读取插件 config.json
// @Tags 插件管理
// @Produce json
// @Security Bearer
// @Param id path int true "服务器ID"
// @Param name path string true "插件名称"
// @Success 200 {object} map[string]string "配置内容，文件不存在时为空"
// @Router /servers/{id}/plugins/{name}/config [get]
func GetPluginConfig(c *gin.Context) {
	getConfig(c, c.Param("name"))
}

// SavePluginConfig 保存插件配置
// @Summary 保存插件 config.json
// @Description 内容必须是合法 JSON。运行中也可以保存，插件重新加载或服务器重启后生效
// @Tags 插件管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "服务器ID"
// @Param name path string true "插件名称"
// @Param body body PluginConfigRequest true "配置内容"
// @Success 200 {object} map[string]string "已保存"
// @Router /servers/{id}/plugins/{name}/config [put]
func SavePluginConfig(c *gin.Context) {
	saveConfig(c, c.Param("name"))
}

// GetArkApiConfig 读取 ArkApi 配置
// @Summary 读取 ArkApi 本体的 config.json
// @Tags 插件管理
// @Produce json
// @Security Bearer
// @Param id path int true "服务器ID"
// @Success 200 {object} map[string]string "配置内容"
// @Failure 404 {object} map[string]string "ArkApi 还没有安装"
// @Router /servers/{id}/arkapi/config [get]
func GetArkApiConfig(c *gin.Context) {
	getConfig(c, "")
}

// SaveArkApiConfig 保存 ArkApi 配置
// @Summary 保存 ArkApi 本体的 config.json
// @Description 内容必须是合法 JSON，重启服务器后生效。入口脚本不会用镜像里的默认配置覆盖它
// @Tags 插件管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "服务器ID"
// @Param body body PluginConfigRequest true "配置内容"
// @Success 200 {object} map[string]string "已保存"
// @Router /servers/{id}/arkapi/config [put]
func SaveArkApiConfig(c *gin.Context) {
	saveConfig(c, "")
}

// ListPluginLogs 插件日志文件列表
// @Summary ArkApi 和插件的日志文件列表
// @Description 列出服务端文件里 ShooterGame/Binaries/Win64/logs 下的 .log 文件，按修改时间从新到旧，最多 200 个（total 是总数）。目录不存在时为空。服务器停止后也能看
// @Tags 插件管理
// @Produce json
// @Security Bearer
// @Param id path int true "服务器ID"
// @Success 200 {object} map[string]server.PluginLogFiles "日志文件列表"
// @Router /servers/{id}/arkapi/logs [get]
func ListPluginLogs(c *gin.Context) {
	files, err := serverService.ListPluginLogs(c.GetUint("user_id"), c.Param("id"))
	if err != nil {
		respondServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": files})
}

// GetPluginLog 读取插件日志
// @Summary 读取一个 ArkApi / 插件日志文件
// @Description 返回 ShooterGame/Binaries/Win64/logs 下这个 .log 文件最后 tail 行（默认 1000，最多 5000），UTF-8 / UTF-16 自动识别
// @Tags 插件管理
// @Produce json
// @Security Bearer
// @Param id path int true "服务器ID"
// @Param file path string true "日志文件名（来自日志文件列表）"
// @Param tail query int false "行数"
// @Success 200 {object} map[string]server.GameLog "日志内容"
// @Failure 404 {object} map[string]string "日志文件不存在"
// @Router /servers/{id}/arkapi/logs/{file} [get]
func GetPluginLog(c *gin.Context) {
	tail, _ := strconv.Atoi(c.DefaultQuery("tail", "1000"))
	if tail <= 0 || tail > 5000 {
		tail = 1000
	}
	log, err := serverService.GetPluginLog(c.GetUint("user_id"), c.Param("id"), c.Param("file"), tail)
	if err != nil {
		respondServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": log})
}

func getConfig(c *gin.Context, name string) {
	content, err := serverService.GetPluginConfig(c.GetUint("user_id"), c.Param("id"), name)
	if err != nil {
		respondServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": gin.H{"content": content}})
}

func saveConfig(c *gin.Context, name string) {
	var req PluginConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondBadRequest(c)
		return
	}
	if err := serverService.SavePluginConfig(c.GetUint("user_id"), c.Param("id"), name, req.Content); err != nil {
		respondServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "配置已保存"})
}

// GetPluginCatalog 插件目录
// @Summary ArkApi 插件目录
// @Description arkfgf.com 整理的 ark-server-api.com 上的 ASE 插件（官方、免费、付费），带 download_url 的可以一键安装。缓存 1 小时，refresh=true 时重新获取。ARKFGF_API 关闭时返回 409 arkfgf_disabled
// @Tags 插件管理
// @Produce json
// @Security Bearer
// @Param refresh query bool false "重新获取"
// @Success 200 {object} map[string]server.PluginCatalog "插件目录"
// @Failure 502 {object} map[string]string "暂时无法获取插件目录"
// @Router /plugins/catalog [get]
func GetPluginCatalog(c *gin.Context) {
	catalog, err := serverService.GetPluginCatalog(c.Query("refresh") == "true")
	if err != nil {
		respondServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": catalog})
}

// InstallCatalogPluginRequest 从插件目录安装
type InstallCatalogPluginRequest struct {
	ResourceID uint `json:"resource_id" binding:"required"`
}

// InstallCatalogPlugin 从插件目录安装或更新插件
// @Summary 从插件目录安装插件
// @Description 按目录里的资源 ID 从 ark-server-api.com 下载并安装（已安装时更新，保留原 config.json 和启用状态）。只能在服务器停止时操作
// @Tags 插件管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "服务器ID"
// @Param body body InstallCatalogPluginRequest true "资源 ID"
// @Success 200 {object} map[string][]string "安装的插件名称"
// @Failure 400 {object} map[string]string "不能一键安装、压缩包无效或服务器未停止"
// @Failure 404 {object} map[string]string "目录里没有这个插件"
// @Failure 502 {object} map[string]string "获取目录或下载失败"
// @Router /servers/{id}/plugins/catalog [post]
func InstallCatalogPlugin(c *gin.Context) {
	var req InstallCatalogPluginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondBadRequest(c)
		return
	}
	names, err := serverService.InstallCatalogPlugin(c.GetUint("user_id"), c.Param("id"), req.ResourceID)
	if err != nil {
		respondServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "插件已安装", "data": names})
}
