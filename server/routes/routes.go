package routes

import (
	"ark-server-commander/controllers/auth"
	"ark-server-commander/controllers/images"
	"ark-server-commander/controllers/servers"
	"ark-server-commander/middleware"
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
		// 服务 Next.js 生成的静态资源
		r.Static("/_next", "./static/_next")
		r.Static("/public", "./static/public")
		r.StaticFile("/favicon.ico", "./static/public/favicon.ico")

		// 处理 SPA 路由 - 所有非 API 和静态资源的请求都返回 index.html
		r.NoRoute(func(c *gin.Context) {
			path := c.Request.URL.Path
			// 如果是 API 请求，返回 404
			if len(path) >= 5 && path[:5] == "/api/" {
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
		authRoutes := api.Group("/auth")
		{
			authRoutes.GET("/check-init", auth.CheckInit)
			authRoutes.POST("/init", auth.InitUser)
			authRoutes.POST("/login", auth.Login)
		}

		// 需要认证的路由
		protected := api.Group("") // 改为空字符串，避免双斜杠
		protected.Use(middleware.AuthMiddleware())
		{
			protected.GET("/profile", auth.GetProfile)

			// 服务器管理路由
			serverRoutes := protected.Group("/servers")
			{
				serverRoutes.GET("", servers.GetServers)
				serverRoutes.POST("", servers.CreateServer)
				serverRoutes.GET("/:id", servers.GetServer)
				serverRoutes.PUT("/:id", servers.UpdateServer)
				serverRoutes.DELETE("/:id", servers.DeleteServer)
				serverRoutes.POST("/:id/start", servers.StartServer)
				serverRoutes.POST("/:id/stop", servers.StopServer)
				serverRoutes.POST("/:id/recreate", servers.RecreateContainer)
				serverRoutes.GET("/:id/rcon", servers.GetServerRCON)
			}

			// 镜像管理路由
			imageRoutes := protected.Group("/images")
			{
				imageRoutes.GET("/status", images.GetImageStatus)
				imageRoutes.POST("/pull", images.PullImage)
				imageRoutes.GET("/check-updates", images.CheckImageUpdates)
				imageRoutes.POST("/update", images.UpdateImage)
				imageRoutes.GET("/affected", images.GetAffectedServers)
			}
		}
	}
}
