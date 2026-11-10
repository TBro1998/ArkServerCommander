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
