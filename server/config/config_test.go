package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadJWTSecretGeneratesAndReuses(t *testing.T) {
	t.Setenv("JWT_SECRET", "")
	dir := t.TempDir()
	DBPath = filepath.Join(dir, "ark_server.db")

	first, err := loadJWTSecret()
	if err != nil {
		t.Fatalf("首次生成密钥失败: %v", err)
	}
	if len(first) < minJWTSecretLength {
		t.Fatalf("生成的密钥过短: %d", len(first))
	}

	info, err := os.Stat(filepath.Join(dir, jwtSecretFileName))
	if err != nil {
		t.Fatalf("密钥文件未生成: %v", err)
	}
	if perm := info.Mode().Perm(); perm != 0o600 {
		t.Fatalf("密钥文件权限应为 0600，实际为 %o", perm)
	}

	second, err := loadJWTSecret()
	if err != nil {
		t.Fatalf("读取已有密钥失败: %v", err)
	}
	if string(first) != string(second) {
		t.Fatal("重启后密钥发生变化，已登录会话会失效")
	}
}

func TestLoadJWTSecretFromEnv(t *testing.T) {
	DBPath = filepath.Join(t.TempDir(), "ark_server.db")

	t.Setenv("JWT_SECRET", "Zq8vR2mXc7LpT4wYk9NbH3sJf6GdA1eU")
	secret, err := loadJWTSecret()
	if err != nil {
		t.Fatalf("合法的 JWT_SECRET 被拒绝: %v", err)
	}
	if string(secret) != "Zq8vR2mXc7LpT4wYk9NbH3sJf6GdA1eU" {
		t.Fatal("未使用环境变量中的密钥")
	}

	t.Setenv("JWT_SECRET", "too-short")
	if _, err := loadJWTSecret(); err == nil {
		t.Fatal("过短的 JWT_SECRET 应被拒绝")
	}
}

func TestLoadArkfgfAPIs(t *testing.T) {
	cases := []struct {
		name                  string
		env                   map[string]string
		account, servers, lnk string
	}{
		{"默认", nil, DefaultAccountAPI, DefaultServersAPI, DefaultLinkAPI},
		{"关闭", map[string]string{"ARKFGF_API": "off", "ARKFGF_LINK_API": "http://x/api/v1"}, "", "", ""},
		{"为空", map[string]string{"ARKFGF_API": " "}, "", "", ""},
		{"本地一个地址", map[string]string{"ARKFGF_API": "http://127.0.0.1:8090/api/v1/"},
			"http://127.0.0.1:8090/api/v1", "http://127.0.0.1:8090/api/v1", "http://127.0.0.1:8090/api/v1"},
		{"分别指定", map[string]string{"ARKFGF_API": "http://a/api/v1", "ARKFGF_SERVERS_API": "http://s/api/v1/", "ARKFGF_LINK_API": "http://l/api/v1"},
			"http://a/api/v1", "http://s/api/v1", "http://l/api/v1"},
		{"只改联机房", map[string]string{"ARKFGF_LINK_API": "http://l/api/v1"}, DefaultAccountAPI, DefaultServersAPI, "http://l/api/v1"},
		{"单独的为空时按 ARKFGF_API", map[string]string{"ARKFGF_API": "http://a/api/v1", "ARKFGF_SERVERS_API": ""},
			"http://a/api/v1", "http://a/api/v1", "http://a/api/v1"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			lookup := func(key string) (string, bool) {
				v, ok := tc.env[key]
				return v, ok
			}
			account, servers, link := loadArkfgfAPIs(lookup)
			if account != tc.account || servers != tc.servers || link != tc.lnk {
				t.Fatalf("得到 %q %q %q，应为 %q %q %q", account, servers, link, tc.account, tc.servers, tc.lnk)
			}
		})
	}
}
