package mcpserver

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"ark-server-commander/database"
	"ark-server-commander/models"
	"ark-server-commander/service/apitoken"

	"github.com/glebarez/sqlite"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(filepath.Join(t.TempDir(), "test.db")), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err := db.AutoMigrate(&models.User{}, &models.Server{}, &models.Schedule{}, &models.APIToken{}); err != nil {
		t.Fatal(err)
	}
	if err := db.Create(&models.User{Username: "admin", Password: "x"}).Error; err != nil {
		t.Fatal(err)
	}

	old := database.DB
	database.DB = db
	t.Cleanup(func() {
		database.DB = old
		if sqlDB, err := db.DB(); err == nil {
			sqlDB.Close()
		}
	})
}

type bearerTransport struct{ token string }

func (b bearerTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req = req.Clone(req.Context())
	req.Header.Set("Authorization", "Bearer "+b.token)
	return http.DefaultTransport.RoundTrip(req)
}

func connect(t *testing.T, url, token string) *mcp.ClientSession {
	t.Helper()
	client := mcp.NewClient(&mcp.Implementation{Name: "test", Version: "1"}, nil)
	session, err := client.Connect(context.Background(), &mcp.StreamableClientTransport{
		Endpoint:   url,
		HTTPClient: &http.Client{Transport: bearerTransport{token}},
		MaxRetries: -1,
	}, nil)
	if err != nil {
		t.Fatalf("连接 MCP 失败: %v", err)
	}
	t.Cleanup(func() { session.Close() })
	return session
}

func toolNames(t *testing.T, session *mcp.ClientSession) []string {
	t.Helper()
	var names []string
	for tool, err := range session.Tools(context.Background(), nil) {
		if err != nil {
			t.Fatal(err)
		}
		names = append(names, tool.Name)
	}
	return names
}

func resultText(res *mcp.CallToolResult) string {
	var sb strings.Builder
	for _, c := range res.Content {
		if tc, ok := c.(*mcp.TextContent); ok {
			sb.WriteString(tc.Text)
		}
	}
	return sb.String()
}

func TestMCPRejectsMissingOrInvalidToken(t *testing.T) {
	setupTestDB(t)
	srv := httptest.NewServer(Handler())
	defer srv.Close()

	body := `{"jsonrpc":"2.0","id":1,"method":"tools/list"}`
	for _, auth := range []string{"", "Bearer nope", "Bearer asc_wrong", "Basic abc"} {
		req, _ := http.NewRequest(http.MethodPost, srv.URL, strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json, text/event-stream")
		if auth != "" {
			req.Header.Set("Authorization", auth)
		}
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		res.Body.Close()
		if res.StatusCode != http.StatusUnauthorized {
			t.Errorf("Authorization=%q 应返回 401，实际 %d", auth, res.StatusCode)
		}
	}
}

func TestMCPReadOnlyTokenOnlySeesQueryTools(t *testing.T) {
	setupTestDB(t)
	srv := httptest.NewServer(Handler())
	defer srv.Close()

	full, err := apitoken.Create(1, "full", false)
	if err != nil {
		t.Fatal(err)
	}
	ro, err := apitoken.Create(1, "ro", true)
	if err != nil {
		t.Fatal(err)
	}

	fullTools := toolNames(t, connect(t, srv.URL, full.Token))
	roSession := connect(t, srv.URL, ro.Token)
	roTools := toolNames(t, roSession)

	for _, name := range []string{"list_servers", "start_server", "delete_server", "get_admin_password", "set_config_values", "rcon_command"} {
		if !slices.Contains(fullTools, name) {
			t.Errorf("完整权限令牌缺少工具 %s", name)
		}
	}
	for _, name := range []string{"list_servers", "get_server", "get_logs", "list_backups", "suggest_ports"} {
		if !slices.Contains(roTools, name) {
			t.Errorf("只读令牌缺少查询工具 %s", name)
		}
	}
	for _, name := range []string{"start_server", "delete_server", "get_admin_password", "rcon_command", "update_image"} {
		if slices.Contains(roTools, name) {
			t.Errorf("只读令牌不应看到 %s", name)
		}
	}
	if len(roTools) >= len(fullTools) {
		t.Errorf("只读工具数 %d 应少于完整工具数 %d", len(roTools), len(fullTools))
	}

	// 只读令牌直接调用修改类工具也必须失败
	res, err := roSession.CallTool(context.Background(), &mcp.CallToolParams{
		Name: "delete_server", Arguments: map[string]any{"server_id": 1},
	})
	if err == nil && (res == nil || !res.IsError) {
		t.Fatal("只读令牌调用 delete_server 应该失败")
	}
}

func TestMCPToolCalls(t *testing.T) {
	setupTestDB(t)
	srv := httptest.NewServer(Handler())
	defer srv.Close()

	full, err := apitoken.Create(1, "full", false)
	if err != nil {
		t.Fatal(err)
	}
	session := connect(t, srv.URL, full.Token)
	ctx := context.Background()

	// 端口建议只读数据库，不需要 Docker
	res, err := session.CallTool(ctx, &mcp.CallToolParams{Name: "suggest_ports", Arguments: map[string]any{}})
	if err != nil {
		t.Fatal(err)
	}
	if res.IsError {
		t.Fatalf("suggest_ports 失败: %s", resultText(res))
	}
	var ports struct {
		Port int `json:"port"`
	}
	if err := json.Unmarshal([]byte(resultText(res)), &ports); err != nil || ports.Port != 7777 {
		t.Fatalf("suggest_ports 结果不对: %s (%v)", resultText(res), err)
	}

	// 服务层错误带错误码返回
	res, err = session.CallTool(ctx, &mcp.CallToolParams{Name: "set_config_values", Arguments: map[string]any{
		"server_id": 1, "file": "Engine.ini", "section": "ServerSettings", "values": map[string]string{"A": "1"},
	}})
	if err != nil {
		t.Fatal(err)
	}
	if !res.IsError || !strings.Contains(resultText(res), "[invalid_config_value]") {
		t.Fatalf("set_config_values 应返回 invalid_config_value，实际 %s", resultText(res))
	}

	res, err = session.CallTool(ctx, &mcp.CallToolParams{Name: "list_schedules", Arguments: map[string]any{"server_id": 42}})
	if err != nil {
		t.Fatal(err)
	}
	if !res.IsError || !strings.Contains(resultText(res), "[server_not_found]") {
		t.Fatalf("不存在的服务器应返回 server_not_found，实际 %s", resultText(res))
	}

	// 缺少必填参数时由 SDK 按输入 schema 拒绝
	res, err = session.CallTool(ctx, &mcp.CallToolParams{Name: "list_schedules", Arguments: map[string]any{}})
	if err == nil && (res == nil || !res.IsError) {
		t.Fatal("缺少 server_id 应该失败")
	}
}
