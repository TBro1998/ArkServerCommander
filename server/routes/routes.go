package routes

import (
	"ark-server-commander/controllers/auth"
	"ark-server-commander/controllers/images"
	"ark-server-commander/controllers/servers"
	"ark-server-commander/controllers/system"
	"ark-server-commander/mcpserver"
	"ark-server-commander/middleware"
	"fmt"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	// 添加健康检查端点（需要日志）
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "message": "服务器运行正常"})
	})

	// 前端由独立的 Next.js 服务提供，后端只负责 /api

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

		// MCP：供 Agent 调用，使用 API 令牌认证（不接受登录 JWT），对外由 Next 的 /mcp 转发。
		// 无状态模式只处理 POST，GET / DELETE 返回 405
		api.Any("/mcp", gin.WrapH(mcpserver.Handler()))

		// 需要认证的路由
		protected := api.Group("") // 改为空字符串，避免双斜杠
		protected.Use(middleware.AuthMiddleware())
		{
			protected.GET("/profile", auth.GetProfile)
			protected.PUT("/auth/password", auth.ChangePassword)

			// MCP 访问令牌：只能用登录 JWT 管理，令牌本身不能创建或撤销令牌
			protected.GET("/tokens", auth.ListAPITokens)
			protected.POST("/tokens", auth.CreateAPIToken)
			protected.DELETE("/tokens/:id", auth.DeleteAPIToken)

			// 系统信息
			protected.GET("/system/overview", system.GetOverview)
			protected.GET("/system/info", system.GetInfo)
			protected.GET("/system/version", system.GetVersion)
			protected.DELETE("/system/clusters/:clusterId", system.DeleteCluster)
			protected.GET("/system/suggest-ports", system.SuggestPorts)
			protected.GET("/mods", system.GetMods)

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
				serverRoutes.POST("/:id/restart", servers.RestartServer)
				serverRoutes.POST("/:id/update-game", servers.UpdateGame)
				serverRoutes.GET("/:id/mods", servers.ListServerMods)
				serverRoutes.GET("/:id/logs", servers.GetLogs)
				serverRoutes.GET("/:id/logs/stream", servers.StreamLogs)
				serverRoutes.GET("/:id/logs/game", servers.GetGameLog)
				serverRoutes.POST("/:id/rcon/command", servers.RconCommand)
				serverRoutes.GET("/:id/players", servers.ListPlayers)

				// 存档备份
				serverRoutes.GET("/:id/backups", servers.ListBackups)
				serverRoutes.POST("/:id/backups", servers.CreateBackup)
				serverRoutes.POST("/:id/backups/upload", servers.UploadBackup)
				serverRoutes.GET("/:id/backups/:name", servers.DownloadBackup)
				serverRoutes.DELETE("/:id/backups/:name", servers.DeleteBackup)
				serverRoutes.POST("/:id/backups/:name/restore", servers.RestoreBackup)

				// 定时任务
				serverRoutes.GET("/:id/schedules", servers.ListSchedules)
				serverRoutes.POST("/:id/schedules", servers.CreateSchedule)
				serverRoutes.PUT("/:id/schedules/:scheduleId", servers.UpdateSchedule)
				serverRoutes.DELETE("/:id/schedules/:scheduleId", servers.DeleteSchedule)
				serverRoutes.POST("/:id/schedules/:scheduleId/run", servers.RunSchedule)

				// ArkApi 插件
				serverRoutes.GET("/:id/plugins", servers.ListPlugins)
				serverRoutes.POST("/:id/plugins", servers.InstallPlugin)
				serverRoutes.PUT("/:id/plugins/:name", servers.SetPluginState)
				serverRoutes.DELETE("/:id/plugins/:name", servers.DeletePlugin)
				serverRoutes.GET("/:id/plugins/:name/config", servers.GetPluginConfig)
				serverRoutes.PUT("/:id/plugins/:name/config", servers.SavePluginConfig)
				serverRoutes.GET("/:id/arkapi/config", servers.GetArkApiConfig)
				serverRoutes.PUT("/:id/arkapi/config", servers.SaveArkApiConfig)
				serverRoutes.POST("/:id/arkapi/install", servers.InstallArkApi)
			}

			// ArkApi 版本库（所有服务器共用）
			arkApiRoutes := protected.Group("/arkapi")
			{
				arkApiRoutes.GET("/versions", servers.ListArkApiVersions)
				arkApiRoutes.POST("/versions", servers.UploadArkApi)
				arkApiRoutes.DELETE("/versions/:version", servers.DeleteArkApiVersion)
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
