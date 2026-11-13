package servers

import (
	"net/http"

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
