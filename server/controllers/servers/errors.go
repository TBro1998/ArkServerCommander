package servers

import (
	"errors"
	"net/http"

	"ark-server-commander/service/server"

	"github.com/gin-gonic/gin"
)

// respondServiceError 将服务层错误映射为 HTTP 状态码并写入响应
func respondServiceError(c *gin.Context, err error) {
	status := http.StatusInternalServerError
	switch {
	case errors.Is(err, server.ErrServerNotFound):
		status = http.StatusNotFound
	case errors.Is(err, server.ErrInvalidServerID),
		errors.Is(err, server.ErrIdentifierExists),
		errors.Is(err, server.ErrServerRunning),
		errors.Is(err, server.ErrAlreadyRunning),
		errors.Is(err, server.ErrAlreadyStarting),
		errors.Is(err, server.ErrAlreadyStopped),
		errors.Is(err, server.ErrAlreadyStopping):
		status = http.StatusBadRequest
	}
	c.JSON(status, gin.H{"error": err.Error()})
}
