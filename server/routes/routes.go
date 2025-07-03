package routes

import (
	"ark-server-manager/controllers"
	"ark-server-manager/middleware"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	// 添加健康检查端点（需要日志）
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "message": "服务器运行正常"})
	})

	// 静态文件服务 - 服务前端文件
	// 检查静态文件目录是否存在
	if _, err := os.Stat("./static"); err == nil {
		// 服务 Nuxt.js 生成的静态资源
		r.Static("/assets", "./static/assets")
		r.Static("/_nuxt", "./static/_nuxt")
		r.StaticFile("/favicon.ico", "./static/favicon.ico")

		// 处理 SPA 路由 - 所有非 API 和静态资源的请求都返回 index.html
		r.NoRoute(func(c *gin.Context) {
			path := c.Request.URL.Path
			// 如果是 API 请求，返回 404
			if len(path) >= 4 && path[:4] == "/api" {
				c.JSON(http.StatusNotFound, gin.H{"error": "API 路由不存在"})
				return
			}
			// 否则返回前端应用的 index.html
			c.File("./static/index.html")
		})
	}

	// API路由组（需要日志）
	api := r.Group("/api")
	api.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("[%s] %s %s %d %s Origin:%s\n",
			param.TimeStamp.Format("2006/01/02 - 15:04:05"),
			param.Method,
			param.Path,
			param.StatusCode,
			param.Latency,
			param.Request.Header.Get("Origin"),
		)
	}))
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
			protected.POST("/servers/:id/rcon/execute", controllers.ExecuteRCONCommand)
			protected.GET("/servers/:id/folder", controllers.GetServerFolderInfo)
			protected.PUT("/servers/:id", controllers.UpdateServer)
			protected.DELETE("/servers/:id", controllers.DeleteServer)
			protected.POST("/servers/:id/start", controllers.StartServer)
			protected.POST("/servers/:id/stop", controllers.StopServer)
		}
	}
}
