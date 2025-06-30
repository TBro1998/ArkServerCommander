package config

import (
	"os"
	"strings"
)

var (
	JWTSecret      = []byte("ark-server-manager-secret-key")
	DBPath         = "ark_server.db"
	ServerPort     = "8080"
	TrustedProxies = []string{"127.0.0.1", "::1"} // 默认只信任本地地址
)

func InitConfig() {
	// 从环境变量读取配置，如果没有则使用默认值
	if secret := os.Getenv("JWT_SECRET"); secret != "" {
		JWTSecret = []byte(secret)
	}

	if dbPath := os.Getenv("DB_PATH"); dbPath != "" {
		DBPath = dbPath
	}

	if port := os.Getenv("SERVER_PORT"); port != "" {
		ServerPort = port
	}

	// 配置可信任的代理地址
	if proxies := os.Getenv("TRUSTED_PROXIES"); proxies != "" {
		TrustedProxies = strings.Split(proxies, ",")
		// 去除空格
		for i, proxy := range TrustedProxies {
			TrustedProxies[i] = strings.TrimSpace(proxy)
		}
	}
}
