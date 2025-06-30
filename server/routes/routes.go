package routes

import (
	"ark-server-manager/controllers"
	"ark-server-manager/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
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
		protected := api.Group("/")
		protected.Use(middleware.AuthMiddleware())
		{
			protected.GET("/profile", controllers.GetProfile)

			// 服务器管理路由
			servers := protected.Group("/servers")
			{
				servers.GET("/", controllers.GetServers)
				servers.POST("/", controllers.CreateServer)
				servers.GET("/:id", controllers.GetServer)
				servers.PUT("/:id", controllers.UpdateServer)
				servers.DELETE("/:id", controllers.DeleteServer)
			}
		}
	}
}
