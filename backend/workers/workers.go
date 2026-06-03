package workers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/hibiken/asynq"
	"github.com/Bruce-Sakura/UploadMyself/backend/config"
	"go.uber.org/zap"
)

const (
	QueueSkill   = "skill"
	QueueVoice   = "voice"
	QueueAvatar  = "avatar"
	QueueDistill = "distill"
)

// TaskType 定义任务类型
const (
	TaskTypeSkillCreate   = "skill:create"
	TaskTypeVoiceTrain    = "voice:train"
	TaskTypeVoiceSynth    = "voice:synthesize"
	TaskTypeAvatar2DGen   = "avatar:2d:generate"
	TaskTypeAvatar2DAnim  = "avatar:2d:animate"
	TaskTypeAvatar3DRecon = "avatar:3d:reconstruct"
	TaskTypeDistill       = "distill:train"
)

// NewAsynqServer 创建异步任务服务器
func NewAsynqServer(cfg *config.Config) *asynq.Server {
	return asynq.NewServer(
		asynq.RedisClientOpt{
			Addr:     cfg.Redis.Addr,
			Password: cfg.Redis.Password,
			DB:       cfg.Redis.DB,
		},
		asynq.Config{
			Concurrency: 4,
			Queues: map[string]int{
				QueueSkill:   3,
				QueueVoice:   3,
				QueueAvatar:  3,
				QueueDistill: 1,
			},
		},
	)
}

// NewAsynqClient 创建异步任务客户端
func NewAsynqClient(cfg *config.Config) *asynq.Client {
	return asynq.NewClient(asynq.RedisClientOpt{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})
}

// EnqueueSkillCreate 提交 Skill 生成任务
func EnqueueSkillCreate(client *asynq.Client, payload map[string]interface{}) (string, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("序列化 payload 失败: %w", err)
	}
	task := asynq.NewTask(TaskTypeSkillCreate, data)
	info, err := client.Enqueue(task, asynq.Queue(QueueSkill))
	if err != nil {
		return "", fmt.Errorf("提交任务失败: %w", err)
	}
	return info.ID, nil
}

// EnqueueVoiceTrain 提交语音训练任务
func EnqueueVoiceTrain(client *asynq.Client, payload map[string]interface{}) (string, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	task := asynq.NewTask(TaskTypeVoiceTrain, data)
	info, err := client.Enqueue(task, asynq.Queue(QueueVoice))
	if err != nil {
		return "", err
	}
	return info.ID, nil
}

// RegisterHandlers 注册所有任务处理器
func RegisterHandlers(mux *asynq.ServeMux, logger *zap.Logger) {
	mux.HandleFunc(TaskTypeSkillCreate, handleSkillCreate)
	mux.HandleFunc(TaskTypeVoiceTrain, handleVoiceTrain)
	mux.HandleFunc(TaskTypeVoiceSynth, handleVoiceSynth)
	mux.HandleFunc(TaskTypeAvatar2DGen, handleAvatar2DGen)
	mux.HandleFunc(TaskTypeAvatar2DAnim, handleAvatar2DAnim)
	mux.HandleFunc(TaskTypeAvatar3DRecon, handleAvatar3DRecon)
	mux.HandleFunc(TaskTypeDistill, handleDistill)
}

func handleSkillCreate(ctx context.Context, t *asynq.Task) error {
	var payload map[string]interface{}
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return err
	}
	log.Printf("[Skill] 开始生成 Skill: %v", payload["name"])
	// TODO: 调用 skill_engine
	return nil
}

func handleVoiceTrain(ctx context.Context, t *asynq.Task) error {
	var payload map[string]interface{}
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return err
	}
	log.Printf("[Voice] 开始训练声音模型: %v", payload["voice_id"])
	// TODO: 调用 voice_engine
	return nil
}

func handleVoiceSynth(ctx context.Context, t *asynq.Task) error {
	log.Printf("[Voice] 开始语音合成")
	// TODO: 调用 voice_engine
	return nil
}

func handleAvatar2DGen(ctx context.Context, t *asynq.Task) error {
	log.Printf("[Avatar2D] 开始生成 2D 形象")
	// TODO: 调用 avatar_engine
	return nil
}

func handleAvatar2DAnim(ctx context.Context, t *asynq.Task) error {
	log.Printf("[Avatar2D] 开始驱动动画")
	// TODO: 调用 avatar_engine
	return nil
}

func handleAvatar3DRecon(ctx context.Context, t *asynq.Task) error {
	log.Printf("[Avatar3D] 开始 3D 重建")
	// TODO: 调用 avatar_engine
	return nil
}

func handleDistill(ctx context.Context, t *asynq.Task) error {
	log.Printf("[Distill] 开始模型蒸馏")
	// TODO: 调用 distill_engine
	return nil
}
