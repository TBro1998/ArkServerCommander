package auth

import (
	"errors"
	"net/http"
	"strconv"

	"ark-server-commander/models"
	"ark-server-commander/service/apitoken"

	"github.com/gin-gonic/gin"
)

// ListAPITokens 列出 MCP 访问令牌
// @Summary 列出 MCP 访问令牌
// @Description 列出当前用户的 MCP 访问令牌（不含令牌明文）
// @Tags 用户
// @Produce json
// @Security Bearer
// @Success 200 {object} map[string][]models.APIToken "令牌列表"
// @Router /tokens [get]
func ListAPITokens(c *gin.Context) {
	tokens, err := apitoken.List(c.GetUint("user_id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "读取令牌失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": tokens})
}

// CreateAPIToken 创建 MCP 访问令牌
// @Summary 创建 MCP 访问令牌
// @Description 生成一个长期有效的令牌，供 Agent 通过 /mcp 调用。响应中的 token 是明文，只返回这一次。read_only 为 true 时只能调用查询类工具
// @Tags 用户
// @Accept json
// @Produce json
// @Security Bearer
// @Param body body models.APITokenRequest true "令牌名称和权限"
// @Success 201 {object} map[string]models.APITokenCreated "新令牌"
// @Failure 400 {object} map[string]string "请求错误"
// @Router /tokens [post]
func CreateAPIToken(c *gin.Context) {
	var req models.APITokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误", "code": "bad_request"})
		return
	}
	created, err := apitoken.Create(c.GetUint("user_id"), req.Name, req.ReadOnly)
	switch {
	case errors.Is(err, apitoken.ErrInvalidTokenName):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "code": "invalid_token_name"})
	case errors.Is(err, apitoken.ErrTooManyTokens):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "code": "too_many_tokens"})
	case err != nil:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建令牌失败"})
	default:
		c.JSON(http.StatusCreated, gin.H{"data": created})
	}
}

// DeleteAPIToken 撤销 MCP 访问令牌
// @Summary 撤销 MCP 访问令牌
// @Description 删除令牌，使用它的 Agent 立即失去访问权限
// @Tags 用户
// @Produce json
// @Security Bearer
// @Param id path int true "令牌ID"
// @Success 200 {object} map[string]string "已撤销"
// @Failure 404 {object} map[string]string "令牌不存在"
// @Router /tokens/{id} [delete]
func DeleteAPIToken(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": apitoken.ErrTokenNotFound.Error(), "code": "token_not_found"})
		return
	}
	err = apitoken.Delete(c.GetUint("user_id"), uint(id))
	switch {
	case errors.Is(err, apitoken.ErrTokenNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error(), "code": "token_not_found"})
	case err != nil:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "撤销令牌失败"})
	default:
		c.JSON(http.StatusOK, gin.H{"message": "令牌已撤销"})
	}
}
