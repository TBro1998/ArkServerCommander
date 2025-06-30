package routes

import (
	"ark-server-manager/controllers"
	"ark-server-manager/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	// 添加健康检查端点
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "message": "服务器运行正常"})
	})

	// API路由组
	api := r.Group("/api")
	{
		// 认证相关路由
		auth := api.Group("/auth")
		{
			auth.GET("/check-init", controllers.CheckInit)
			auth.POST("/init", controllers.InitUser)
			auth.POST("/login", controllers.Login)
		}

		// 需要认证的路由
		protected := api.Group("") // 改为空字符串，避免双斜杠
		protected.Use(middleware.AuthMiddleware())
		{
			protected.GET("/profile", controllers.GetProfile)

			// 服务器管理路由 - 直接在 protected 组下定义
			protected.GET("/servers", controllers.GetServers)
			protected.POST("/servers", controllers.CreateServer)
			protected.GET("/servers/:id", controllers.GetServer)
			protected.GET("/servers/:id/rcon", controllers.GetServerRCON)
			protected.GET("/servers/:id/folder", controllers.GetServerFolderInfo)
			protected.GET("/servers/:id/config", controllers.GetServerConfig)
			protected.GET("/servers/:id/config/files", controllers.ListServerConfigFiles)
			protected.PUT("/servers/:id", controllers.UpdateServer)
			protected.PUT("/servers/:id/config", controllers.UpdateServerConfig)
			protected.DELETE("/servers/:id", controllers.DeleteServer)
			protected.POST("/servers/:id/start", controllers.StartServer)
			protected.POST("/servers/:id/stop", controllers.StopServer)
		}
	}
}
