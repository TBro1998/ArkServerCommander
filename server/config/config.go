package config

import (
	"fmt"
	"os"
	"strings"
)

var (
	JWTSecret  []byte
	DBPath     = "ark_server.db"
	ServerPort = "8080"
)

// 弱密钥黑名单
var weakSecrets = []string{
	"ark-server-commander-secret-key",
	"secret",
	"password",
	"123456",
	"default",
	"changeme",
	"test",
}

func InitConfig() error {
	// JWT密钥必须从环境变量读取
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return fmt.Errorf("JWT_SECRET environment variable is required")
	}

	// 验证密钥长度
	if len(secret) < 32 {
		return fmt.Errorf("JWT_SECRET must be at least 32 characters long (current: %d)", len(secret))
	}

	// 检查是否使用了弱密钥
	secretLower := strings.ToLower(secret)
	for _, weak := range weakSecrets {
		if strings.Contains(secretLower, weak) {
			return fmt.Errorf("JWT_SECRET contains weak/common password pattern: '%s'", weak)
		}
	}

	JWTSecret = []byte(secret)

	// 读取其他配置
	if dbPath := os.Getenv("DB_PATH"); dbPath != "" {
		DBPath = dbPath
	}

	if port := os.Getenv("SERVER_PORT"); port != "" {
		ServerPort = port
	}

	return nil
}
