package api

import (
	"github.com/gin-gonic/gin"
	"github.com/Bruce-Sakura/UploadMyself/backend/config"
	"github.com/Bruce-Sakura/UploadMyself/backend/services"
)

// NewRouter 创建并配置所有路由
func NewRouter(cfg *config.Config, svc *services.Services) *gin.Engine {
	if cfg.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	// CORS 中间件
	r.Use(corsMiddleware())

	// 健康检查
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"name":    "UploadMyself",
			"version": "0.1.0",
			"docs":    "/docs",
		})
	})
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// API v1
	v1 := r.Group("/api/v1")
	{
		// 思维框架克隆
		skill := v1.Group("/skill")
		{
			skill.POST("/create", skillHandler.Create)
			skill.GET("/:id/result", skillHandler.Result)
			skill.GET("/:id/download", skillHandler.Download)
		}

		// 语音克隆
		voice := v1.Group("/voice")
		{
			voice.POST("/upload", voiceHandler.Upload)
			voice.POST("/train", voiceHandler.Train)
			voice.POST("/synthesize", voiceHandler.Synthesize)
			voice.GET("/:id/samples", voiceHandler.Samples)
		}

		// 2D 虚拟形象
		avatar2d := v1.Group("/avatar/2d")
		{
			avatar2d.POST("/upload", avatar2DHandler.Upload)
			avatar2d.POST("/generate", avatar2DHandler.Generate)
			avatar2d.POST("/animate", avatar2DHandler.Animate)
		}

		// 3D 虚拟形象
		avatar3d := v1.Group("/avatar/3d")
		{
			avatar3d.POST("/upload", avatar3DHandler.Upload)
			avatar3d.POST("/reconstruct", avatar3DHandler.Reconstruct)
			avatar3d.GET("/:id/model", avatar3DHandler.Model)
			avatar3d.GET("/:id/preview", avatar3DHandler.Preview)
		}

		// 模型蒸馏
		distill := v1.Group("/distill")
		{
			distill.POST("/start", distillHandler.Start)
			distill.GET("/:id/status", distillHandler.Status)
			distill.GET("/:id/metrics", distillHandler.Metrics)
		}

		// 任务管理
		tasks := v1.Group("/tasks")
		{
			tasks.GET("/:id", taskHandler.Status)
		}
	}

	return r
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

// 临时 handler 占位 — 后续拆分到独立文件
var (
	skillHandler    = &SkillHandler{}
	voiceHandler    = &VoiceHandler{}
	avatar2DHandler = &Avatar2DHandler{}
	avatar3DHandler = &Avatar3DHandler{}
	distillHandler  = &DistillHandler{}
	taskHandler     = &TaskHandler{}
)
