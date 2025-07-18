package main

import (
	"log"

	"ark-server-commander/config"
	"ark-server-commander/database"
	"ark-server-commander/routes"
	"ark-server-commander/service/docker_manager"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "ark-server-commander/docs" // 导入生成的docs包
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

	// 检查Docker环境
	if err := docker_manager.CheckDockerStatus(); err != nil {
		log.Fatalf("Docker环境检查失败: %v\n请确保Docker已安装并运行", err)
	}
	log.Println("✅ Docker环境检查通过")

	// 获取Docker管理器单例实例
	dockerManager, err := docker_manager.GetDockerManager()
	if err != nil {
		log.Fatalf("获取Docker管理器失败: %v", err)
	}
	defer docker_manager.CloseDockerManager()

	// 验证必要的Docker镜像是否存在（不自动下载）
	log.Println("🔍 检查必要的Docker镜像...")
	missingImages, err := dockerManager.ValidateRequiredImages()
	if err != nil {
		log.Printf("⚠️  镜像检查失败: %v", err)
	} else if len(missingImages) > 0 {
		log.Printf("⚠️  缺失镜像: %v\n请手动下载镜像后再启动服务器", missingImages)
	} else {
		log.Println("✅ 所有必要镜像已存在")
	}

	// 为现有服务器初始化Docker卷和配置文件（不创建容器）
	if err := docker_manager.InitializeDockerForExistingServers(); err != nil {
		log.Printf("Warning: Failed to initialize Docker volumes and config files for existing servers: %v", err)
	}

	// 同步服务器状态与Docker容器状态
	if err := docker_manager.SyncServerStatusWithDocker(); err != nil {
		log.Printf("Warning: Failed to sync server status with Docker: %v", err)
	}

	// 创建Gin实例
	r := gin.Default()

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
	log.Println("🔄 Docker镜像正在后台检查中...")
	log.Println("📋 Docker卷和配置文件初始化完成")
	log.Println("📋 服务器状态同步完成")
	log.Println("=========================================")

	r.Run(":8080")
}
