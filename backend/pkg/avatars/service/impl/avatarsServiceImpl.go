package impl

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/Bruce-Sakura/UploadMyself/backend/pkg/avatars/dto"
	"github.com/Bruce-Sakura/UploadMyself/backend/pkg/avatars/entity"
	"github.com/Bruce-Sakura/UploadMyself/backend/pkg/avatars/mapper"
	"github.com/Bruce-Sakura/UploadMyself/backend/pkg/avatars/service"
	taskdto "github.com/Bruce-Sakura/UploadMyself/backend/pkg/tasks/dto"
	taskservice "github.com/Bruce-Sakura/UploadMyself/backend/pkg/tasks/service"
	"github.com/google/uuid"
)

// Config holds runtime settings for avatar generation.
type Config struct {
	MLServiceURL string // e.g. http://host.docker.internal:8001
	UploadDir    string // e.g. ./uploads
}

type AvatarServiceImpl struct {
	mapper  *mapper.AvatarMapper
	taskSvc taskservice.TaskService
	cfg     Config
}

func NewAvatarService(m *mapper.AvatarMapper, taskSvc taskservice.TaskService, cfg Config) service.AvatarService {
	return &AvatarServiceImpl{mapper: m, taskSvc: taskSvc, cfg: cfg}
}

func (s *AvatarServiceImpl) Create(ctx context.Context, req dto.CreateAvatarReq) (*dto.AvatarVO, error) {
	a := &entity.Avatar{
		ID:        uuid.NewString(),
		Name:      req.Name,
		Type:      req.Type,
		PhotoPath: req.PhotoPath,
		Style:     req.Style,
		Status:    "pending",
	}
	if err := s.mapper.Insert(ctx, a); err != nil {
		return nil, err
	}
	return s.Get(ctx, a.ID)
}

func (s *AvatarServiceImpl) Get(ctx context.Context, id string) (*dto.AvatarVO, error) {
	a, err := s.mapper.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return toVO(a), nil
}

func (s *AvatarServiceImpl) List(ctx context.Context, req dto.ListAvatarReq) ([]dto.AvatarVO, error) {
	as, err := s.mapper.List(ctx, req.Type)
	if err != nil {
		return nil, err
	}
	out := make([]dto.AvatarVO, 0, len(as))
	for i := range as {
		out = append(out, *toVO(&as[i]))
	}
	return out, nil
}

func (s *AvatarServiceImpl) Delete(ctx context.Context, id string) error {
	return s.mapper.Delete(ctx, id)
}

// Process triggers async 2D/3D generation through the GPU ML service.
func (s *AvatarServiceImpl) Process(ctx context.Context, id string) (string, error) {
	a, err := s.mapper.GetByID(ctx, id)
	if err != nil {
		return "", err
	}

	task, err := s.taskSvc.Create(ctx, taskdto.CreateTaskReq{Type: "avatar_process", RefID: id})
	if err != nil {
		return "", err
	}
	_ = s.mapper.UpdateStatus(ctx, id, "processing")

	go s.runGeneration(task.ID, a.ID, a.Type, a.PhotoPath)

	return task.ID, nil
}

func (s *AvatarServiceImpl) runGeneration(taskID, avatarID, avatarType, photoPath string) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()

	_ = s.taskSvc.UpdateStatus(ctx, taskID, "running", 10, "")

	outputDir := fmt.Sprintf("%s/%s_output", s.cfg.UploadDir, avatarID)

	// mode tells the ML service whether to run the full 3D reconstruction.
	// 2D skips the expensive 3D recon so it returns quickly (avoids poll timeout).
	mode := "3d"
	if avatarType == "2d" {
		mode = "2d"
	}

	reqBody, _ := json.Marshal(map[string]any{
		"input_path": photoPath,
		"output_dir": outputDir,
		"mode":       mode,
		"seed":       2333,
		"timestep":   40,
	})

	resp, err := http.Post(s.cfg.MLServiceURL+"/generate-avatar", "application/json", bytes.NewReader(reqBody))
	if err != nil {
		s.fail(ctx, taskID, avatarID, fmt.Sprintf("ML service error: %v", err))
		return
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	var result map[string]any
	_ = json.Unmarshal(respBody, &result)

	if resp.StatusCode != http.StatusOK {
		msg := "unknown error"
		if e, ok := result["error"]; ok {
			msg = fmt.Sprintf("%v", e)
		}
		s.fail(ctx, taskID, avatarID, fmt.Sprintf("ML error: %s", msg))
		return
	}

	s.normalizePaths(result)

	// output_path: 3D uses the GLB model; 2D uses the cartoon image.
	outputPath := ""
	if avatarType == "3d" {
		if v, ok := result["glb_model"].(string); ok {
			outputPath = v
		}
	}
	if outputPath == "" {
		if v, ok := result["cartoon_image"].(string); ok {
			outputPath = v
		}
	}

	resultJSON, _ := json.Marshal(result)
	if err := s.mapper.UpdateResult(ctx, avatarID, string(resultJSON), outputPath, "done"); err != nil {
		s.fail(ctx, taskID, avatarID, fmt.Sprintf("save result: %v", err))
		return
	}
	_ = s.taskSvc.UpdateStatus(ctx, taskID, "done", 100, "")
}

// normalizePaths rewrites absolute/relative ML paths to servable /uploads URLs.
func (s *AvatarServiceImpl) normalizePaths(result map[string]any) {
	for _, key := range []string{"cartoon_image", "glb_model", "obj_model"} {
		if v, ok := result[key].(string); ok {
			result[key] = strings.Replace(v, s.cfg.UploadDir, "uploads", 1)
		}
	}
	if views, ok := result["views"].([]any); ok {
		for i, v := range views {
			if str, ok := v.(string); ok {
				views[i] = strings.Replace(str, s.cfg.UploadDir, "uploads", 1)
			}
		}
	}
}

func (s *AvatarServiceImpl) fail(ctx context.Context, taskID, avatarID, msg string) {
	_ = s.taskSvc.UpdateStatus(ctx, taskID, "failed", 0, msg)
	_ = s.mapper.UpdateStatus(ctx, avatarID, "failed")
}

func toVO(a *entity.Avatar) *dto.AvatarVO {
	return &dto.AvatarVO{
		ID:         a.ID,
		Name:       a.Name,
		Type:       a.Type,
		PhotoPath:  a.PhotoPath,
		Style:      a.Style,
		Status:     a.Status,
		Result:     a.Result,
		OutputPath: a.OutputPath,
		CreatedAt:  a.CreatedAt,
		UpdatedAt:  a.UpdatedAt,
	}
}
