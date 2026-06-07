package main

import (
	"context"
	_ "embed"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/Bruce-Sakura/UploadMyself/backend/internal/llm"
	"github.com/Bruce-Sakura/UploadMyself/backend/middleware"

	avatarhandler "github.com/Bruce-Sakura/UploadMyself/backend/pkg/avatars/handler"
	avatarmapper "github.com/Bruce-Sakura/UploadMyself/backend/pkg/avatars/mapper"
	avatarimpl "github.com/Bruce-Sakura/UploadMyself/backend/pkg/avatars/service/impl"

	filehandler "github.com/Bruce-Sakura/UploadMyself/backend/pkg/file_uploads/handler"
	filemapper "github.com/Bruce-Sakura/UploadMyself/backend/pkg/file_uploads/mapper"
	fileimpl "github.com/Bruce-Sakura/UploadMyself/backend/pkg/file_uploads/service/impl"

	msghandler "github.com/Bruce-Sakura/UploadMyself/backend/pkg/messages/handler"
	msgmapper "github.com/Bruce-Sakura/UploadMyself/backend/pkg/messages/mapper"
	msgimpl "github.com/Bruce-Sakura/UploadMyself/backend/pkg/messages/service/impl"

	skillhandler "github.com/Bruce-Sakura/UploadMyself/backend/pkg/skills/handler"
	skillmapper "github.com/Bruce-Sakura/UploadMyself/backend/pkg/skills/mapper"
	skillimpl "github.com/Bruce-Sakura/UploadMyself/backend/pkg/skills/service/impl"

	taskhandler "github.com/Bruce-Sakura/UploadMyself/backend/pkg/tasks/handler"
	taskmapper "github.com/Bruce-Sakura/UploadMyself/backend/pkg/tasks/mapper"
	taskimpl "github.com/Bruce-Sakura/UploadMyself/backend/pkg/tasks/service/impl"

	voicehandler "github.com/Bruce-Sakura/UploadMyself/backend/pkg/voices/handler"
	voicemapper "github.com/Bruce-Sakura/UploadMyself/backend/pkg/voices/mapper"
	voiceimpl "github.com/Bruce-Sakura/UploadMyself/backend/pkg/voices/service/impl"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

//go:embed migrations/001_init.sql
var initSQL string

func main() {
	// ---- Config ----
	viper.AutomaticEnv()
	viper.SetDefault("APP_PORT", 8000)
	viper.SetDefault("DB_DSN", "host=localhost user=uploadmyself password=uploadmyself dbname=uploadmyself port=5432 sslmode=disable")
	viper.SetDefault("ML_SCRIPTS_DIR", "../ml/scripts")
	viper.SetDefault("PYTHON_BIN", "python3")
	viper.SetDefault("ML_SERVICE_URL", "http://host.docker.internal:8001")
	viper.SetDefault("SKILLS_DIR", "./skills")
	viper.SetDefault("LLM_API_KEY", "")
	viper.SetDefault("LLM_BASE_URL", "https://api.openai.com/v1")
	viper.SetDefault("LLM_MODEL", "mimo-v2.5-pro")

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	uploadsDir := "./uploads"
	if err := os.MkdirAll(uploadsDir, 0755); err != nil {
		log.Fatalf("create uploads dir: %v", err)
	}
	skillsDir := viper.GetString("SKILLS_DIR")
	if err := os.MkdirAll(skillsDir, 0755); err != nil {
		log.Fatalf("create skills dir: %v", err)
	}

	// ---- Database (pgxpool) ----
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, viper.GetString("DB_DSN"))
	if err != nil {
		log.Fatalf("db connect: %v", err)
	}
	defer pool.Close()
	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("db ping: %v", err)
	}
	if err := runMigrations(ctx, pool); err != nil {
		log.Fatalf("migrate: %v", err)
	}
	logger.Info("database connected & migrated")

	// ---- LLM client ----
	llmClient := llm.New(
		viper.GetString("LLM_API_KEY"),
		viper.GetString("LLM_BASE_URL"),
		viper.GetString("LLM_MODEL"),
	)

	mlScriptsDir := viper.GetString("ML_SCRIPTS_DIR")
	pythonBin := viper.GetString("PYTHON_BIN")
	mlServiceURL := viper.GetString("ML_SERVICE_URL")

	// ---- Dependency injection: mapper -> service -> handler ----
	taskSvc := taskimpl.NewTaskService(taskmapper.NewTaskMapper(pool))
	taskH := taskhandler.NewTaskHandler(taskSvc)

	skillSvc := skillimpl.NewSkillService(skillmapper.NewSkillMapper(pool), taskSvc, llmClient, skillimpl.Config{
		SkillsDir: skillsDir,
	})
	skillH := skillhandler.NewSkillHandler(skillSvc)

	voiceSvc := voiceimpl.NewVoiceService(voicemapper.NewVoiceMapper(pool), taskSvc, voiceimpl.Config{
		MLScriptsDir: mlScriptsDir,
		PythonBin:    pythonBin,
		UploadDir:    uploadsDir,
	})
	voiceH := voicehandler.NewVoiceHandler(voiceSvc)

	avatarSvc := avatarimpl.NewAvatarService(avatarmapper.NewAvatarMapper(pool), taskSvc, avatarimpl.Config{
		MLServiceURL: mlServiceURL,
		UploadDir:    uploadsDir,
	})
	avatarH := avatarhandler.NewAvatarHandler(avatarSvc)

	fileSvc := fileimpl.NewFileUploadService(filemapper.NewFileUploadMapper(pool), fileimpl.Config{
		UploadDir:    uploadsDir,
		MLScriptsDir: mlScriptsDir,
		PythonBin:    pythonBin,
	})
	fileH := filehandler.NewFileUploadHandler(fileSvc)

	agentSvc := msgimpl.NewAgentService(msgmapper.NewMessageMapper(pool), llmClient, skillSvc)
	agentH := msghandler.NewAgentHandler(agentSvc)

	// ---- Router ----
	r := gin.Default()
	r.Use(middleware.CORS())
	r.GET("/health", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) })
	r.Static("/uploads", uploadsDir)

	v1 := r.Group("/api/v1")
	taskH.Register(v1)
	skillH.Register(v1)
	voiceH.Register(v1)
	avatarH.Register(v1)
	fileH.Register(v1)
	agentH.Register(v1)

	// ---- Serve ----
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

// runMigrations executes the embedded DDL statement-by-statement.
// (pgx's extended protocol rejects multi-statement strings, so we split on ';'.)
func runMigrations(ctx context.Context, pool *pgxpool.Pool) error {
	for _, stmt := range strings.Split(initSQL, ";") {
		if strings.TrimSpace(stmt) == "" {
			continue
		}
		if _, err := pool.Exec(ctx, stmt); err != nil {
			return fmt.Errorf("%w\nstatement: %s", err, strings.TrimSpace(stmt))
		}
	}
	return nil
}
