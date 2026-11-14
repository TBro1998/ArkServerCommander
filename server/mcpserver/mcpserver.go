// Package mcpserver 通过 MCP（Model Context Protocol）把服务器管理功能开放给 Agent。
//
// 传输使用无状态的 Streamable HTTP，挂在后端的 /api/mcp，对外由 Next 的 /mcp 路由转发。
// 认证使用设置页生成的 API 令牌（Authorization: Bearer asc_...），不接受登录用的 JWT。
// 工具直接调用 service 层，和 REST 接口共用同一套业务逻辑和校验。
// 只读令牌拿到的是只注册了查询类工具的 MCP 服务，看不到也调用不了修改类工具
package mcpserver

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"slices"
	"strconv"

	"ark-server-commander/service/apitoken"
	"ark-server-commander/service/server"
	"ark-server-commander/utils"

	"github.com/modelcontextprotocol/go-sdk/auth"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"go.uber.org/zap"
)

const (
	scopeRead  = "read"
	scopeWrite = "write"
)

const instructions = `ARK Server Commander manages ARK: Survival Evolved dedicated servers, each running in its own Docker container.

- Start with list_servers to get server IDs, status and live runtime info. Use list_servers (optionally with server_id) to poll status; get_server is heavier (it reads config files from disk) and should only be used when you need the config contents or launch arguments.
- start_server, stop_server, restart_server, recreate_container and update_game are asynchronous: they return as soon as the operation has started. Poll list_servers until status becomes "running" or "stopped"; if an operation failed, last_error explains why.
- A running container is not necessarily joinable. runtime.phase shows progress: pulling_image, copying, downloading (first start downloads about 20 GB), installing_mods, launching, ready. Players can connect only when the phase is "ready".
- Config edits (update_server, set_config_values) made while a server is not stopped are stored as pending (config_pending=true) and written to disk on the next start, because ARK overwrites GameUserSettings.ini when it shuts down. Restart the server to apply them.
- Players connect at <public_host>:<port>; get_system_info returns public_host.
- Errors are returned as "[error_code] message"; messages are in Chinese.`

var service = server.NewServerService()

// Handler 返回 MCP 的 HTTP 处理器（已包含令牌认证）
func Handler() http.Handler {
	full := newServer(false)
	readOnly := newServer(true)

	h := mcp.NewStreamableHTTPHandler(func(r *http.Request) *mcp.Server {
		info := auth.TokenInfoFromContext(r.Context())
		if info == nil {
			return nil
		}
		if slices.Contains(info.Scopes, scopeWrite) {
			return full
		}
		return readOnly
	}, &mcp.StreamableHTTPOptions{
		// 无状态 + JSON 响应：每个请求独立处理，经过 Next 转发也不需要维持会话
		Stateless:    true,
		JSONResponse: true,
	})

	return auth.RequireBearerToken(verifyToken, &auth.RequireBearerTokenOptions{
		// 令牌长期有效，撤销后立即失效，没有过期时间
		AllowMissingExpiration: true,
	})(h)
}

func verifyToken(_ context.Context, token string, _ *http.Request) (*auth.TokenInfo, error) {
	record, err := apitoken.Verify(token)
	if errors.Is(err, apitoken.ErrInvalidToken) {
		return nil, auth.ErrInvalidToken
	}
	if err != nil {
		utils.Warn("校验 MCP 令牌失败", zap.Error(err))
		return nil, err
	}
	scopes := []string{scopeRead}
	if !record.ReadOnly {
		scopes = append(scopes, scopeWrite)
	}
	return &auth.TokenInfo{
		UserID: strconv.FormatUint(uint64(record.UserID), 10),
		Scopes: scopes,
		Extra:  map[string]any{"token_id": record.ID},
	}, nil
}

func newServer(readOnly bool) *mcp.Server {
	s := mcp.NewServer(&mcp.Implementation{
		Name:    "ark-server-commander",
		Title:   "ARK Server Commander",
		Version: utils.Version,
	}, &mcp.ServerOptions{Instructions: instructions})

	r := &registry{server: s, readOnly: readOnly}
	registerServerTools(r)
	registerOperationTools(r)
	registerPluginTools(r)
	registerSystemTools(r)
	return s
}

// ---------- 工具注册 ----------

type registry struct {
	server   *mcp.Server
	readOnly bool
}

// handlerFunc 工具的业务处理函数：返回值序列化为 JSON 作为工具结果，错误作为工具错误返回给 Agent
type handlerFunc[In any] func(ctx context.Context, userID uint, in In) (any, error)

type noInput struct{}

func ptr[T any](v T) *T { return &v }

// query 注册查询类工具，只读令牌也能使用
func query[In any](r *registry, name, title, description string, h handlerFunc[In]) {
	add(r, &mcp.Tool{
		Name:        name,
		Description: description,
		Annotations: &mcp.ToolAnnotations{Title: title, ReadOnlyHint: true, OpenWorldHint: ptr(false)},
	}, h)
}

// secret 注册返回敏感信息（如管理员密码）的查询工具，只读令牌不能使用
func secret[In any](r *registry, name, title, description string, h handlerFunc[In]) {
	if r.readOnly {
		return
	}
	query(r, name, title, description, h)
}

type actionKind int

const (
	additive    actionKind = iota // 只新增或改变状态，不丢数据
	idempotent                    // 重复调用没有额外影响
	destructive                   // 会删除或覆盖数据
)

// action 注册修改类工具，只读令牌不能使用
func action[In any](r *registry, kind actionKind, name, title, description string, h handlerFunc[In]) {
	if r.readOnly {
		return
	}
	add(r, &mcp.Tool{
		Name:        name,
		Description: description,
		Annotations: &mcp.ToolAnnotations{
			Title:           title,
			DestructiveHint: ptr(kind == destructive),
			IdempotentHint:  kind == idempotent,
			OpenWorldHint:   ptr(false),
		},
	}, h)
}

func add[In any](r *registry, tool *mcp.Tool, h handlerFunc[In]) {
	mcp.AddTool(r.server, tool, func(ctx context.Context, req *mcp.CallToolRequest, in In) (*mcp.CallToolResult, any, error) {
		userID, err := userIDFrom(req)
		if err != nil {
			return nil, nil, err
		}
		out, err := h(ctx, userID, in)
		if err != nil {
			return nil, nil, toolError(err)
		}
		if out == nil {
			out = map[string]any{"ok": true}
		}
		return nil, out, nil
	})
}

func userIDFrom(req *mcp.CallToolRequest) (uint, error) {
	if req.Extra == nil || req.Extra.TokenInfo == nil {
		return 0, errors.New("unauthenticated")
	}
	id, err := strconv.ParseUint(req.Extra.TokenInfo.UserID, 10, 32)
	if err != nil {
		return 0, errors.New("unauthenticated")
	}
	return uint(id), nil
}

// toolError 在错误信息前加上错误码，Agent 可以据此判断原因（信息本身是中文）
func toolError(err error) error {
	if code := server.ErrorCode(err); code != "" {
		return fmt.Errorf("[%s] %s", code, err.Error())
	}
	return err
}

func sid(id uint) string { return strconv.FormatUint(uint64(id), 10) }

// message 异步操作或无返回值操作的结果
func message(text string) map[string]any { return map[string]any{"ok": true, "message": text} }
