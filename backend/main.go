package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Bruce-Sakura/UploadMyself/backend/handler"
	"github.com/Bruce-Sakura/UploadMyself/backend/middleware"
	"github.com/Bruce-Sakura/UploadMyself/backend/model"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func main() {
	// Config
	viper.AutomaticEnv()
	viper.SetDefault("APP_PORT", 8000)
	viper.SetDefault("DB_DSN", "host=localhost user=uploadmyself password=uploadmyself dbname=uploadmyself port=5432 sslmode=disable")
	viper.SetDefault("ML_SCRIPTS_DIR", "../ml/scripts")
	viper.SetDefault("PYTHON_BIN", "python3")

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	// Configure handler-level globals for ML scripts
	handler.MLScriptsDir = viper.GetString("ML_SCRIPTS_DIR")
	handler.PythonBin = viper.GetString("PYTHON_BIN")

	// Ensure uploads directory exists
	uploadsDir := "./uploads"
	if err := os.MkdirAll(uploadsDir, 0755); err != nil {
		log.Fatalf("create uploads dir: %v", err)
	}
	handler.UploadDir = uploadsDir

	// Database
	db, err := model.Connect(viper.GetString("DB_DSN"))
	if err != nil {
		log.Fatalf("db connect: %v", err)
	}
	logger.Info("database connected")

	// Router
	r := gin.Default()
	r.Use(middleware.CORS())

	// Health
	r.GET("/health", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) })

	// API v1
	h := handler.New(db)
	v1 := r.Group("/api/v1")
	{
		// Upload & file serving
		v1.POST("/upload", h.UploadFile)
		v1.GET("/files/:id", h.ServeFile)

		// Skills
		skill := v1.Group("/skills")
		{
			skill.POST("", h.CreateSkill)
			skill.GET("", h.ListSkills)
			skill.GET("/:id", h.GetSkill)
			skill.PUT("/:id", h.UpdateSkill)
			skill.DELETE("/:id", h.DeleteSkill)
			skill.POST("/:id/process", h.ProcessSkill)
		}
		// Voices
		voice := v1.Group("/voices")
		{
			voice.POST("", h.CreateVoice)
			voice.GET("", h.ListVoices)
			voice.GET("/:id", h.GetVoice)
			voice.DELETE("/:id", h.DeleteVoice)
			voice.POST("/:id/train", h.TrainVoice)
			voice.POST("/:id/synthesize", h.SynthesizeVoice)
		}
		// Avatars
		avatar := v1.Group("/avatars")
		{
			avatar.POST("", h.CreateAvatar)
			avatar.GET("", h.ListAvatars)
			avatar.GET("/:id", h.GetAvatar)
			avatar.DELETE("/:id", h.DeleteAvatar)
			avatar.POST("/:id/process", h.ProcessAvatar)
		}
		// Tasks
		task := v1.Group("/tasks")
		{
			task.GET("", h.ListTasks)
			task.GET("/:id", h.GetTask)
		}
	}

	// Serve
	addr := fmt.Sprintf(":%d", viper.GetInt("APP_PORT"))
	go func() {
		logger.Info("starting server", zap.String("addr", addr))
		if err := r.Run(addr); err != nil {
			log.Fatalf("server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("shutting down")
}
