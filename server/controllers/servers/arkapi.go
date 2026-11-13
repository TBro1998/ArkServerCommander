package servers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// maxArkApiUpload ArkApi 压缩包上传大小上限
const maxArkApiUpload = 100 << 20

// ListArkApiVersions ArkApi 版本库
// @Summary ArkApi 版本列表
// @Description 内置版本、已下载或上传到数据目录 arkapi/ 的版本，以及 GitHub（Michidu/ARK-Server-API）上发布的版本，从新到旧。GitHub 结果缓存 6 小时，refresh=true 时重新查询
// @Tags 插件管理
// @Produce json
// @Security Bearer
// @Param refresh query bool false "重新查询 GitHub"
// @Success 200 {object} map[string]server.ArkApiVersions "版本列表"
// @Router /arkapi/versions [get]
func ListArkApiVersions(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"data": serverService.ListArkApiVersions(c.Query("refresh") == "true")})
}

// UploadArkApi 上传 ArkApi 压缩包
// @Summary 上传 ArkApi
// @Description 上传 ArkApi zip（multipart 字段 file）到版本库，压缩包里要有 version.dll。版本号从文件名（如 ArkApi_3.56.zip）识别，识别不出时用 version 指定；同版本会被覆盖
// @Tags 插件管理
// @Accept multipart/form-data
// @Produce json
// @Security Bearer
// @Param file formData file true "ArkApi zip"
// @Param version formData string false "版本号"
// @Success 200 {object} map[string]string "版本号"
// @Failure 400 {object} map[string]string "压缩包无效"
// @Router /arkapi/versions [post]
func UploadArkApi(c *gin.Context) {
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxArkApiUpload+1<<20)
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请上传不超过 100 MB 的 zip 文件", "code": "bad_request"})
		return
	}
	f, err := file.Open()
	if err != nil {
		respondBadRequest(c)
		return
	}
	defer f.Close()
	version, err := serverService.UploadArkApi(f, file.Size, file.Filename, c.PostForm("version"))
	if err != nil {
		respondServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "已上传", "data": version})
}

// DeleteArkApiVersion 从版本库删除
// @Summary 删除 ArkApi 版本
// @Description 删除版本库里下载或上传的压缩包，不影响已经安装到服务器的文件；内置版本不能删除
// @Tags 插件管理
// @Produce json
// @Security Bearer
// @Param version path string true "版本号"
// @Success 200 {object} map[string]string "已删除"
// @Router /arkapi/versions/{version} [delete]
func DeleteArkApiVersion(c *gin.Context) {
	if err := serverService.DeleteArkApiVersion(c.Param("version")); err != nil {
		respondServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "已删除"})
}

// InstallArkApiRequest 安装 ArkApi
type InstallArkApiRequest struct {
	Version string `json:"version"` // 为空时安装最新版本
}

// InstallArkApi 安装或更新服务器的 ArkApi
// @Summary 安装 ArkApi
// @Description 把版本库里的指定版本安装到服务器的 ShooterGame/Binaries/Win64（只在 GitHub 上的版本先下载），保留已有的 config.json。只能在服务器停止时操作；version 为空时安装最新版本
// @Tags 插件管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "服务器ID"
// @Param body body InstallArkApiRequest true "版本"
// @Success 200 {object} map[string]string "安装的版本"
// @Failure 400 {object} map[string]string "服务器未停止"
// @Router /servers/{id}/arkapi/install [post]
func InstallArkApi(c *gin.Context) {
	var req InstallArkApiRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondBadRequest(c)
		return
	}
	version, err := serverService.InstallArkApi(c.GetUint("user_id"), c.Param("id"), req.Version)
	if err != nil {
		respondServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "ArkApi 已安装", "data": version})
}
