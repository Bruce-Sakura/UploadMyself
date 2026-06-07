package service

import (
	"context"

	"github.com/Bruce-Sakura/UploadMyself/backend/pkg/skills/dto"
)

// SkillService is the business-logic contract for skills.
type SkillService interface {
	Create(ctx context.Context, req dto.CreateSkillReq) (*dto.SkillVO, error)
	Get(ctx context.Context, id string) (*dto.SkillVO, error)
	List(ctx context.Context) ([]dto.SkillVO, error)
	Update(ctx context.Context, id string, req dto.UpdateSkillReq) (*dto.SkillVO, error)
	Delete(ctx context.Context, id string) error
	// Process generates the SKILL.md via LLM asynchronously; returns the task ID.
	Process(ctx context.Context, id string) (taskID string, err error)
	// Import downloads a ready-made SKILL.md (URL/GitHub) or stores inline content
	// as a new skill package on disk.
	Import(ctx context.Context, req dto.ImportSkillReq) (*dto.SkillVO, error)
}
