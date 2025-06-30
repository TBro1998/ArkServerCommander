package config

import (
	"os"
)

var (
	JWTSecret  = []byte("ark-server-manager-secret-key")
	DBPath     = "ark_server.db"
	ServerPort = "8080"
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
}
