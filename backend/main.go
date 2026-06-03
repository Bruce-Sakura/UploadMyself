package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Bruce-Sakura/UploadMyself/backend/config"
	"github.com/Bruce-Sakura/UploadMyself/backend/api"
	"github.com/Bruce-Sakura/UploadMyself/backend/services"
	"go.uber.org/zap"
)

func main() {
	// 加载配置
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 初始化日志
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	zap.ReplaceGlobals(logger)

	// 初始化服务
	svc, err := services.New(cfg)
	if err != nil {
		log.Fatalf("初始化服务失败: %v", err)
	}
	defer svc.Close()

	// 创建路由
	router := api.NewRouter(cfg, svc)

	// 启动服务
	addr := fmt.Sprintf("%s:%d", cfg.App.Host, cfg.App.Port)
	go func() {
		logger.Info("启动服务", zap.String("addr", addr))
		if err := router.Run(addr); err != nil {
			log.Fatalf("服务启动失败: %v", err)
		}
	}()

	// 优雅退出
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("服务正在关闭...")
}
