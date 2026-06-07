package impl

import (
	"context"

	"github.com/Bruce-Sakura/UploadMyself/backend/pkg/tasks/dto"
	"github.com/Bruce-Sakura/UploadMyself/backend/pkg/tasks/entity"
	"github.com/Bruce-Sakura/UploadMyself/backend/pkg/tasks/mapper"
	"github.com/Bruce-Sakura/UploadMyself/backend/pkg/tasks/service"
	"github.com/google/uuid"
)

type TaskServiceImpl struct {
	mapper *mapper.TaskMapper
}

// NewTaskService wires the task mapper into a TaskService implementation.
func NewTaskService(m *mapper.TaskMapper) service.TaskService {
	return &TaskServiceImpl{mapper: m}
}

func (s *TaskServiceImpl) Create(ctx context.Context, req dto.CreateTaskReq) (*dto.TaskVO, error) {
	t := &entity.Task{
		ID:       uuid.NewString(),
		Type:     req.Type,
		RefID:    req.RefID,
		Status:   "pending",
		Progress: 0,
	}
	if err := s.mapper.Insert(ctx, t); err != nil {
		return nil, err
	}
	return s.Get(ctx, t.ID)
}

func (s *TaskServiceImpl) Get(ctx context.Context, id string) (*dto.TaskVO, error) {
	t, err := s.mapper.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return toVO(t), nil
}

func (s *TaskServiceImpl) UpdateStatus(ctx context.Context, id, status string, progress int, errMsg string) error {
	return s.mapper.UpdateStatus(ctx, id, status, progress, errMsg)
}

func (s *TaskServiceImpl) List(ctx context.Context, req dto.ListTaskReq) ([]dto.TaskVO, error) {
	ts, err := s.mapper.List(ctx, req.Type, req.RefID)
	if err != nil {
		return nil, err
	}
	out := make([]dto.TaskVO, 0, len(ts))
	for i := range ts {
		out = append(out, *toVO(&ts[i]))
	}
	return out, nil
}

func toVO(t *entity.Task) *dto.TaskVO {
	return &dto.TaskVO{
		ID:        t.ID,
		Type:      t.Type,
		RefID:     t.RefID,
		Status:    t.Status,
		Progress:  t.Progress,
		Error:     t.Error,
		CreatedAt: t.CreatedAt,
		UpdatedAt: t.UpdatedAt,
	}
}
