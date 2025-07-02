package main

import (
	"fmt"
	"log"

	"ark-server-manager/config"
	"ark-server-manager/database"
	"ark-server-manager/routes"
	"ark-server-manager/utils"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "ark-server-manager/docs" // 导入生成的docs包
)

// @title ARK服务器管理器 API
// @version 1.0
// @description 基于Gin+Gorm的ARK服务器管理系统API文档
// @termsOfService http://swagger.io/terms/

// @contact.name API支持
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description JWT token in the format "Bearer {token}"

func main() {
	// 初始化配置
	config.InitConfig()

	// 初始化数据库
	database.InitDB()

	// 为现有服务器初始化Docker容器和卷
	if err := utils.InitializeDockerForExistingServers(); err != nil {
		log.Printf("Warning: Failed to initialize Docker for existing servers: %v", err)
	}

	// 同步服务器状态与Docker容器状态
	if err := utils.SyncServerStatusWithDocker(); err != nil {
		log.Printf("Warning: Failed to sync server status with Docker: %v", err)
	}

	// 创建Gin实例
	r := gin.Default()

	// 配置可信任的代理（安全设置）
	r.SetTrustedProxies(config.TrustedProxies)

	// 最简单的CORS解决方案 - 允许所有来源（仅开发环境）
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization, X-Requested-With")
		c.Header("Access-Control-Allow-Credentials", "true")

		// 处理预检请求
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// 添加请求日志中间件
	r.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("[%s] %s %s %d %s Origin:%s\n",
			param.TimeStamp.Format("2006/01/02 - 15:04:05"),
			param.Method,
			param.Path,
			param.StatusCode,
			param.Latency,
			param.Request.Header.Get("Origin"),
		)
	}))

	// 注册路由
	routes.RegisterRoutes(r)

	// 添加Swagger路由
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 启动服务器
	log.Println("=========================================")
	log.Println("🚀 ARK服务器管理器后端启动成功")
	log.Println("📍 服务器地址: http://localhost:8080")
	log.Println("📚 API文档: http://localhost:8080/swagger/index.html")
	log.Println("🔗 健康检查: http://localhost:8080/health")
	log.Println("🌐 CORS: 已启用（允许所有来源）")
	log.Printf("🐳 Docker容器化ARK服务器管理")
	log.Println("=========================================")

	r.Run(":8080")
}
