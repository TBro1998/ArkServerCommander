package servers

import (
	"errors"
	"net/http"

	"ark-server-commander/service/server"

	"github.com/gin-gonic/gin"
)

// RespondServiceError 将服务层错误映射为 HTTP 状态码并写入响应。
// 响应包含 error（中文详细信息）和 code（错误码，未知错误没有），前端按语言选择显示哪个
func RespondServiceError(c *gin.Context, err error) {
	status := http.StatusInternalServerError
	switch {
	case errors.Is(err, server.ErrServerNotFound),
		errors.Is(err, server.ErrBackupNotFound),
		errors.Is(err, server.ErrScheduleNotFound),
		errors.Is(err, server.ErrPluginNotFound),
		errors.Is(err, server.ErrArkApiNotInstalled),
		errors.Is(err, server.ErrArkApiNotFound):
		status = http.StatusNotFound
	case errors.Is(err, server.ErrInvalidServerID),
		errors.Is(err, server.ErrIdentifierExists),
		errors.Is(err, server.ErrServerRunning),
		errors.Is(err, server.ErrAlreadyRunning),
		errors.Is(err, server.ErrAlreadyStarting),
		errors.Is(err, server.ErrAlreadyStopped),
		errors.Is(err, server.ErrAlreadyStopping),
		errors.Is(err, server.ErrInvalidPort),
		errors.Is(err, server.ErrInvalidClusterID),
		errors.Is(err, server.ErrInvalidModIDs),
		errors.Is(err, server.ErrServerNotRunning),
		errors.Is(err, server.ErrServerNotStopped),
		errors.Is(err, server.ErrInvalidCommand),
		errors.Is(err, server.ErrInvalidBackup),
		errors.Is(err, server.ErrInvalidSchedule),
		errors.Is(err, server.ErrInvalidPlugin),
		errors.Is(err, server.ErrInvalidPluginConfig),
		errors.Is(err, server.ErrInvalidArkApi),
		errors.Is(err, server.ErrInvalidConfigValue),
		errors.Is(err, server.ErrArkApiBuiltin),
		errors.Is(err, server.ErrImageNotAllowed):
		status = http.StatusBadRequest
	case errors.Is(err, server.ErrPortConflict),
		errors.Is(err, server.ErrServerBusy),
		errors.Is(err, server.ErrBackupBusy),
		errors.Is(err, server.ErrClusterInUse),
		errors.Is(err, server.ErrImagePulling),
		errors.Is(err, server.ErrRconUnavailable):
		status = http.StatusConflict
	case errors.Is(err, server.ErrRconAuthFailed),
		errors.Is(err, server.ErrArkApiDownload):
		status = http.StatusBadGateway
	}
	body := gin.H{"error": err.Error()}
	if code := server.ErrorCode(err); code != "" {
		body["code"] = code
	}
	c.JSON(status, body)
}

// respondServiceError 包内旧名称，保持现有调用不变
func respondServiceError(c *gin.Context, err error) { RespondServiceError(c, err) }

// respondBadRequest 请求参数错误
func respondBadRequest(c *gin.Context) {
	c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误", "code": "bad_request"})
}
