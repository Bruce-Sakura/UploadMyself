package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	App      AppConfig
	DB       DBConfig
	Redis    RedisConfig
	MinIO    MinIOConfig
	Provider ProviderConfig
	Models   ModelsConfig
}

type AppConfig struct {
	Host     string
	Port     int
	Env      string
	LogLevel string
}

type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
}

func (d DBConfig) DSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		d.Host, d.Port, d.User, d.Password, d.Name)
}

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

type MinIOConfig struct {
	Endpoint  string
	AccessKey string
	SecretKey string
	Bucket    string
	UseSSL    bool
}

type ProviderConfig struct {
	Mode string // local | cloud | hybrid

	// 云端 API Keys
	OpenAIKey      string
	AnthropicKey   string
	QwenKey        string
	ElevenLabsKey  string
	FishAudioKey   string
	HeyGenKey      string
	Tripo3DKey     string
}

type ModelsConfig struct {
	Dir            string
	VoiceModelDir  string
	AvatarModelDir string
	BaseModelDir   string
}

func Load() (*Config, error) {
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.AddConfigPath("../")
	viper.AutomaticEnv()

	// 默认值
	viper.SetDefault("APP_HOST", "0.0.0.0")
	viper.SetDefault("APP_PORT", 8000)
	viper.SetDefault("APP_ENV", "development")
	viper.SetDefault("LOG_LEVEL", "info")
	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_PORT", 5432)
	viper.SetDefault("DB_USER", "uploadmyself")
	viper.SetDefault("DB_PASSWORD", "uploadmyself")
	viper.SetDefault("DB_NAME", "uploadmyself")
	viper.SetDefault("REDIS_ADDR", "localhost:6379")
	viper.SetDefault("REDIS_PASSWORD", "")
	viper.SetDefault("REDIS_DB", 0)
	viper.SetDefault("MINIO_ENDPOINT", "localhost:9000")
	viper.SetDefault("MINIO_ACCESS_KEY", "minioadmin")
	viper.SetDefault("MINIO_SECRET_KEY", "minioadmin")
	viper.SetDefault("MINIO_BUCKET", "uploadmyself")
	viper.SetDefault("PROVIDER_MODE", "local")
	viper.SetDefault("MODEL_DIR", "./ml/models")

	// 读取 .env 文件（忽略不存在）
	_ = viper.ReadInConfig()

	cfg := &Config{
		App: AppConfig{
			Host:     viper.GetString("APP_HOST"),
			Port:     viper.GetInt("APP_PORT"),
			Env:      viper.GetString("APP_ENV"),
			LogLevel: viper.GetString("LOG_LEVEL"),
		},
		DB: DBConfig{
			Host:     viper.GetString("DB_HOST"),
			Port:     viper.GetInt("DB_PORT"),
			User:     viper.GetString("DB_USER"),
			Password: viper.GetString("DB_PASSWORD"),
			Name:     viper.GetString("DB_NAME"),
		},
		Redis: RedisConfig{
			Addr:     viper.GetString("REDIS_ADDR"),
			Password: viper.GetString("REDIS_PASSWORD"),
			DB:       viper.GetInt("REDIS_DB"),
		},
		MinIO: MinIOConfig{
			Endpoint:  viper.GetString("MINIO_ENDPOINT"),
			AccessKey: viper.GetString("MINIO_ACCESS_KEY"),
			SecretKey: viper.GetString("MINIO_SECRET_KEY"),
			Bucket:    viper.GetString("MINIO_BUCKET"),
			UseSSL:    viper.GetBool("MINIO_USE_SSL"),
		},
		Provider: ProviderConfig{
			Mode:          viper.GetString("PROVIDER_MODE"),
			OpenAIKey:     viper.GetString("OPENAI_API_KEY"),
			AnthropicKey:  viper.GetString("ANTHROPIC_API_KEY"),
			QwenKey:       viper.GetString("QWEN_API_KEY"),
			ElevenLabsKey: viper.GetString("ELEVENLABS_API_KEY"),
			FishAudioKey:  viper.GetString("FISH_AUDIO_API_KEY"),
			HeyGenKey:     viper.GetString("HEYGEN_API_KEY"),
			Tripo3DKey:    viper.GetString("TRIPO3D_API_KEY"),
		},
		Models: ModelsConfig{
			Dir:            viper.GetString("MODEL_DIR"),
			VoiceModelDir:  viper.GetString("MODEL_DIR") + "/voice",
			AvatarModelDir: viper.GetString("MODEL_DIR") + "/avatar_2d",
			BaseModelDir:   viper.GetString("MODEL_DIR") + "/base",
		},
	}

	return cfg, nil
}
