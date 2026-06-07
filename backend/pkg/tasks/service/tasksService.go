package service

import (
	"context"

	"github.com/Bruce-Sakura/UploadMyself/backend/pkg/tasks/dto"
)

// TaskService is the business-logic contract for async tasks.
// Other modules (avatars, skills, voices) depend on this to create/track tasks.
type TaskService interface {
	Create(ctx context.Context, req dto.CreateTaskReq) (*dto.TaskVO, error)
	Get(ctx context.Context, id string) (*dto.TaskVO, error)
	UpdateStatus(ctx context.Context, id, status string, progress int, errMsg string) error
	List(ctx context.Context, req dto.ListTaskReq) ([]dto.TaskVO, error)
}
