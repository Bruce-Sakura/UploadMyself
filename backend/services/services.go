package services

import (
	"context"
	"fmt"

	"github.com/Bruce-Sakura/UploadMyself/backend/config"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

// Services 聚合所有服务依赖
type Services struct {
	Config  *config.Config
	Redis   *redis.Client
	Logger  *zap.Logger
	// TODO: DB, MinIO, ModelRegistry 等
}

func New(cfg *config.Config) (*Services, error) {
	logger := zap.L()

	// 连接 Redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	ctx := context.Background()
	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("连接 Redis 失败: %w", err)
	}
	logger.Info("Redis 连接成功")

	// TODO: 连接 PostgreSQL
	// TODO: 连接 MinIO
	// TODO: 初始化 ModelRegistry

	return &Services{
		Config: cfg,
		Redis:  rdb,
		Logger: logger,
	}, nil
}

func (s *Services) Close() {
	if s.Redis != nil {
		s.Redis.Close()
	}
	// TODO: 关闭 DB, MinIO 等
}
