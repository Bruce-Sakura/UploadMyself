package service

import (
	"context"

	"github.com/Bruce-Sakura/UploadMyself/backend/pkg/messages/dto"
)

// AgentService is the core conversation engine contract.
type AgentService interface {
	Chat(ctx context.Context, req dto.ChatReq) (*dto.ChatResp, error)
	ListTools() []dto.ToolInfoVO
}
