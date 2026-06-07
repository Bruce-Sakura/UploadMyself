package impl

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/Bruce-Sakura/UploadMyself/backend/pkg/voices/dto"
	"github.com/Bruce-Sakura/UploadMyself/backend/pkg/voices/entity"
	"github.com/Bruce-Sakura/UploadMyself/backend/pkg/voices/mapper"
	"github.com/Bruce-Sakura/UploadMyself/backend/pkg/voices/service"
	taskdto "github.com/Bruce-Sakura/UploadMyself/backend/pkg/tasks/dto"
	taskservice "github.com/Bruce-Sakura/UploadMyself/backend/pkg/tasks/service"
	"github.com/google/uuid"
)

// Config holds runtime settings for voice ML scripts.
type Config struct {
	MLScriptsDir string
	PythonBin    string
	UploadDir    string
}

type VoiceServiceImpl struct {
	mapper  *mapper.VoiceMapper
	taskSvc taskservice.TaskService
	cfg     Config
}

func NewVoiceService(m *mapper.VoiceMapper, taskSvc taskservice.TaskService, cfg Config) service.VoiceService {
	return &VoiceServiceImpl{mapper: m, taskSvc: taskSvc, cfg: cfg}
}

func (s *VoiceServiceImpl) Create(ctx context.Context, req dto.CreateVoiceReq) (*dto.VoiceVO, error) {
	v := &entity.Voice{
		ID:        uuid.NewString(),
		Name:      req.Name,
		AudioPath: req.AudioPath,
		Duration:  req.Duration,
		Status:    "pending",
	}
	if err := s.mapper.Insert(ctx, v); err != nil {
		return nil, err
	}
	return s.Get(ctx, v.ID)
}

func (s *VoiceServiceImpl) Get(ctx context.Context, id string) (*dto.VoiceVO, error) {
	v, err := s.mapper.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return toVO(v), nil
}

func (s *VoiceServiceImpl) List(ctx context.Context) ([]dto.VoiceVO, error) {
	vs, err := s.mapper.List(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]dto.VoiceVO, 0, len(vs))
	for i := range vs {
		out = append(out, *toVO(&vs[i]))
	}
	return out, nil
}

func (s *VoiceServiceImpl) Delete(ctx context.Context, id string) error {
	return s.mapper.Delete(ctx, id)
}

func (s *VoiceServiceImpl) Train(ctx context.Context, id string) (string, error) {
	v, err := s.mapper.GetByID(ctx, id)
	if err != nil {
		return "", err
	}
	task, err := s.taskSvc.Create(ctx, taskdto.CreateTaskReq{Type: "voice_train", RefID: id})
	if err != nil {
		return "", err
	}
	_ = s.mapper.UpdateStatus(ctx, id, "training")

	go s.runTrain(task.ID, v.ID, v.Name, v.AudioPath)
	return task.ID, nil
}

func (s *VoiceServiceImpl) runTrain(taskID, voiceID, name, audioPath string) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()

	_ = s.taskSvc.UpdateStatus(ctx, taskID, "running", 10, "")

	configPath := fmt.Sprintf("%s/%s_train_config.json", s.cfg.UploadDir, voiceID)
	configData := fmt.Sprintf(`{"voice_id":"%s","audio_path":"%s","name":"%s"}`, voiceID, audioPath, name)
	_ = os.WriteFile(configPath, []byte(configData), 0644)

	script := fmt.Sprintf("%s/voice_clone_train.py", s.cfg.MLScriptsDir)
	cmd := exec.CommandContext(ctx, s.cfg.PythonBin, script, "--config", configPath)
	out, err := cmd.CombinedOutput()
	if err != nil {
		_ = s.mapper.UpdateStatus(ctx, voiceID, "failed")
		_ = s.taskSvc.UpdateStatus(ctx, taskID, "failed", 0, fmt.Sprintf("train error: %v\noutput: %s", err, string(out)))
		return
	}

	var result map[string]any
	if json.Unmarshal(out, &result) == nil {
		if mp, ok := result["model_path"]; ok {
			_ = s.mapper.UpdateModelPath(ctx, voiceID, fmt.Sprintf("%v", mp))
		}
	}

	_ = s.mapper.UpdateStatus(ctx, voiceID, "done")
	_ = s.taskSvc.UpdateStatus(ctx, taskID, "done", 100, "")
}

func (s *VoiceServiceImpl) Synthesize(ctx context.Context, id, text string) (string, error) {
	if _, err := s.mapper.GetByID(ctx, id); err != nil {
		return "", err
	}

	cctx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()

	outputPath := fmt.Sprintf("%s/%s_synth.wav", s.cfg.UploadDir, id)
	script := fmt.Sprintf("%s/voice_synthesize.py", s.cfg.MLScriptsDir)
	cmd := exec.CommandContext(cctx, s.cfg.PythonBin, script,
		"--voice-id", id,
		"--text", text,
		"--output", outputPath,
		"--model-dir", s.cfg.UploadDir,
	)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("synthesize error: %v\noutput: %s", err, string(out))
	}

	var result map[string]any
	if err := json.Unmarshal(out, &result); err != nil {
		return "", fmt.Errorf("invalid script output")
	}
	audioPath, _ := result["audio_path"].(string)
	return audioPath, nil
}

func toVO(v *entity.Voice) *dto.VoiceVO {
	return &dto.VoiceVO{
		ID:           v.ID,
		Name:         v.Name,
		AudioPath:    v.AudioPath,
		Duration:     v.Duration,
		ModelPath:    v.ModelPath,
		RefAudioPath: v.RefAudioPath,
		Status:       v.Status,
		CreatedAt:    v.CreatedAt,
		UpdatedAt:    v.UpdatedAt,
	}
}
