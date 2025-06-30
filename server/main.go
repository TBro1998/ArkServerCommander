package main

import (
	"log"

	"ark-server-manager/config"
	"ark-server-manager/database"
	"ark-server-manager/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化配置
	config.InitConfig()

	// 初始化数据库
	database.InitDB()

	// 创建Gin实例
	r := gin.Default()

	// 配置CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// 注册路由
	routes.RegisterRoutes(r)

	// 启动服务器
	log.Println("服务器启动在端口 :8080")
	r.Run(":8080")
}
