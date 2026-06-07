package service

import (
	"context"

	"github.com/Bruce-Sakura/UploadMyself/backend/pkg/voices/dto"
)

// VoiceService is the business-logic contract for voices.
type VoiceService interface {
	Create(ctx context.Context, req dto.CreateVoiceReq) (*dto.VoiceVO, error)
	Get(ctx context.Context, id string) (*dto.VoiceVO, error)
	List(ctx context.Context) ([]dto.VoiceVO, error)
	Delete(ctx context.Context, id string) error
	// Train kicks off async voice model training; returns the tracking task ID.
	Train(ctx context.Context, id string) (taskID string, err error)
	// Synthesize runs TTS synchronously and returns the generated audio path.
	Synthesize(ctx context.Context, id, text string) (audioPath string, err error)
}
