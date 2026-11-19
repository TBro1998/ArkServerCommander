package routes

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"ark-server-commander/utils"

	"github.com/gin-gonic/gin"
)

// 完成 arkfgf.com 登录的接口不要求 JWT（回调页可能在 Steam 客户端的内置浏览器里打开），
// 但只接受本机开始的登录的 state：随便提交的 state 返回 400 expired_login，不是 401
func TestArkfgfLoginFinishIsPublic(t *testing.T) {
	utils.InitLogger()
	gin.SetMode(gin.TestMode)
	r := gin.New()
	RegisterRoutes(r)

	do := func(method, path, body string) (*httptest.ResponseRecorder, map[string]any) {
		req := httptest.NewRequest(method, path, strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		var out map[string]any
		_ = json.Unmarshal(w.Body.Bytes(), &out)
		return w, out
	}

	w, out := do(http.MethodPost, "/api/arkfgf/account/login/finish", `{"params":{"state":"acs_forged","openid.mode":"id_res"}}`)
	if w.Code != http.StatusBadRequest || out["code"] != "expired_login" {
		t.Fatalf("finish without a local state: %d %v", w.Code, out)
	}
	w, _ = do(http.MethodPost, "/api/arkfgf/account/login/finish", `not json`)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("bad body: %d", w.Code)
	}

	// 其他账号接口仍然要求登录
	for _, rt := range [][2]string{{http.MethodGet, "/api/arkfgf/account"}, {http.MethodPost, "/api/arkfgf/account/login"}, {http.MethodDelete, "/api/arkfgf/account"}} {
		if w, _ := do(rt[0], rt[1], `{}`); w.Code != http.StatusUnauthorized {
			t.Fatalf("%s %s without JWT: %d", rt[0], rt[1], w.Code)
		}
	}
}
