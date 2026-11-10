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

	// ListenHost 监听地址，默认仅本机回环；容器内由 docker-entrypoint.sh 指定
	ListenHost = "127.0.0.1"

	// CORSAllowedOrigins 允许跨域的来源白名单，为空时仅同源访问
	CORSAllowedOrigins []string
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

	if host := os.Getenv("LISTEN_HOST"); host != "" {
		ListenHost = host
	}

	if origins := os.Getenv("CORS_ALLOWED_ORIGINS"); origins != "" {
		CORSAllowedOrigins = strings.Split(origins, ",")
	}

	return nil
}
