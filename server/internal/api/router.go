package api

import (
	"github.com/agenthub/server/internal/config"
	"github.com/agenthub/server/internal/storage"
	"github.com/gin-gonic/gin"
)

// NewRouter 创建路由
func NewRouter(cfg *config.Config, store *storage.Storage) *gin.Engine {
	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(CORSMiddleware())

	// 创建处理器
	h := NewHandler(cfg, store)

	// API 版本
	v1 := r.Group("/api/v1")
	{
		// 健康检查
		v1.GET("/health", h.Health)

		// 认证
		auth := v1.Group("/auth")
		{
			auth.POST("/register", h.Register)
			auth.POST("/login", h.Login)
			auth.POST("/refresh", h.RefreshToken)
			auth.POST("/logout", AuthMiddleware(cfg), h.Logout)
		}

		// 用户
		users := v1.Group("/users")
		{
			users.GET("/:username", h.GetUser)
			users.GET("/:username/agents", h.GetUserAgents)
			users.PUT("/me", AuthMiddleware(cfg), h.UpdateProfile)
		}

		// 智能体
		agents := v1.Group("/agents")
		{
			agents.GET("", h.ListAgents)
			agents.GET("/:namespace/:name", h.GetAgent)
			agents.GET("/:namespace/:name/versions", h.ListVersions)
			agents.GET("/:namespace/:name/versions/:version", h.GetVersion)
			agents.GET("/:namespace/:name/files/*path", h.GetFile)

			// 需要认证
			agents.POST("", AuthMiddleware(cfg), h.CreateAgent)
			agents.PUT("/:namespace/:name", AuthMiddleware(cfg), h.UpdateAgent)
			agents.DELETE("/:namespace/:name", AuthMiddleware(cfg), h.DeleteAgent)
			agents.POST("/:namespace/:name/versions", AuthMiddleware(cfg), h.PublishVersion)
			agents.POST("/:namespace/:name/like", AuthMiddleware(cfg), h.LikeAgent)
			agents.DELETE("/:namespace/:name/like", AuthMiddleware(cfg), h.UnlikeAgent)
		}

		// 搜索
		v1.GET("/search", h.Search)

		// 分类
		v1.GET("/categories", h.ListCategories)

		// 热门/推荐
		v1.GET("/trending", h.GetTrending)
		v1.GET("/featured", h.GetFeatured)

		// API Keys
		keys := v1.Group("/keys", AuthMiddleware(cfg))
		{
			keys.GET("", h.ListAPIKeys)
			keys.POST("", h.CreateAPIKey)
			keys.DELETE("/:id", h.DeleteAPIKey)
		}
	}

	// Agent 调用接口 (需要 API Key)
	invoke := r.Group("/invoke", APIKeyMiddleware(store))
	{
		invoke.POST("/:namespace/:name", h.InvokeAgent)
		invoke.POST("/:namespace/:name/stream", h.InvokeAgentStream)
	}

	return r
}

// CORSMiddleware CORS 中间件
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Authorization, Accept, X-Requested-With, X-API-Key")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
