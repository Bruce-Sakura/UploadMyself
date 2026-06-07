package service

import (
	"context"

	"github.com/Bruce-Sakura/UploadMyself/backend/pkg/avatars/dto"
)

// AvatarService is the business-logic contract for avatars.
type AvatarService interface {
	Create(ctx context.Context, req dto.CreateAvatarReq) (*dto.AvatarVO, error)
	Get(ctx context.Context, id string) (*dto.AvatarVO, error)
	List(ctx context.Context, req dto.ListAvatarReq) ([]dto.AvatarVO, error)
	Delete(ctx context.Context, id string) error
	// Process kicks off async 2D/3D generation and returns the tracking task ID.
	Process(ctx context.Context, id string) (taskID string, err error)
}
