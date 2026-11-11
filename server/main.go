package main

import (
	"net"

	"ark-server-commander/config"
	"ark-server-commander/database"
	"ark-server-commander/middleware"
	"ark-server-commander/routes"
	"ark-server-commander/service/docker_manager"
	"ark-server-commander/service/server"
	"ark-server-commander/utils"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"

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
	// 初始化日志系统
	if err := utils.InitLogger(); err != nil {
		panic("日志系统初始化失败: " + err.Error())
	}
	defer utils.Sync()

	// 初始化配置
	if err := config.InitConfig(); err != nil {
		utils.Fatal("配置初始化失败", zap.Error(err))
	}

	// 初始化数据库
	database.InitDB()

	// 检查Docker环境
	if err := docker_manager.CheckDockerStatus(); err != nil {
		utils.Fatal("Docker环境检查失败，请确保Docker已安装并运行", zap.Error(err))
	}
	utils.Info("Docker环境检查通过")

	// 获取Docker管理器单例实例
	_, err := docker_manager.GetDockerManager()
	if err != nil {
		utils.Fatal("获取Docker管理器失败", zap.Error(err))
	}
	defer docker_manager.CloseDockerManager()

	// 确定数据目录在宿主机上的路径，ARK 服务器的卷会绑定到其中
	if err := docker_manager.InitStorage(); err != nil {
		utils.Fatal("初始化数据目录失败", zap.Error(err))
	}

	// 上次退出时未完成的启停操作，按容器实际状态校正
	if err := server.NewServerService().ReconcileStatuses(); err != nil {
		utils.Warn("校正服务器状态失败", zap.Error(err))
	}

	// 后台采集运行阶段、在线人数和资源占用；定时任务调度
	server.StartRuntimeMonitor()
	server.StartScheduler(server.NewServerService())

	// 创建Gin实例
	r := gin.Default()

	// CORS：默认仅同源，通过 CORS_ALLOWED_ORIGINS 配置白名单
	r.Use(middleware.CORS(config.CORSAllowedOrigins))

	// 注册路由
	routes.RegisterRoutes(r)

	// 添加Swagger路由
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 启动服务器
	utils.Info("=========================================")
	utils.Info("🚀 ARK服务器管理器后端启动成功")
	addr := net.JoinHostPort(config.ListenHost, config.ServerPort)
	utils.Infof("📍 服务器地址: http://%s", addr)
	utils.Infof("📚 API文档: http://%s/swagger/index.html", addr)
	utils.Infof("🔗 健康检查: http://%s/health", addr)
	utils.Info("🌐 CORS: 仅允许 CORS_ALLOWED_ORIGINS 白名单来源，未配置时仅同源")
	utils.Info("🐳 Docker容器化ARK服务器管理")
	utils.Info("🔄 Docker镜像正在后台检查中...")
	utils.Info("📋 Docker卷和配置文件初始化完成")
	utils.Info("📋 服务器状态同步完成")
	utils.Info("=========================================")

	if err := r.Run(addr); err != nil {
		utils.Fatal("服务器启动失败", zap.Error(err))
	}
}
