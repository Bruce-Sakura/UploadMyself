package impl

import (
	"context"
	"fmt"
	"time"

	"github.com/Bruce-Sakura/UploadMyself/backend/internal/llm"
	"github.com/Bruce-Sakura/UploadMyself/backend/pkg/messages/dto"
	"github.com/Bruce-Sakura/UploadMyself/backend/pkg/messages/mapper"
	"github.com/Bruce-Sakura/UploadMyself/backend/pkg/messages/service"
	skillservice "github.com/Bruce-Sakura/UploadMyself/backend/pkg/skills/service"
)

const defaultSystemPrompt = "你是一个有帮助的AI助手。"

type AgentServiceImpl struct {
	mapper   *mapper.MessageMapper
	llm      *llm.Client
	skillSvc skillservice.SkillService
	tools    *toolRegistry
}

// NewAgentService wires the message mapper, LLM client, and skill service
// (to load SKILL.md as the system prompt) into the agent engine.
func NewAgentService(m *mapper.MessageMapper, llmClient *llm.Client, skillSvc skillservice.SkillService) service.AgentService {
	return &AgentServiceImpl{
		mapper:   m,
		llm:      llmClient,
		skillSvc: skillSvc,
		tools:    newToolRegistry(),
	}
}

func (s *AgentServiceImpl) ListTools() []dto.ToolInfoVO {
	return s.tools.info()
}

func (s *AgentServiceImpl) Chat(ctx context.Context, req dto.ChatReq) (*dto.ChatResp, error) {
	systemPrompt := s.loadSkillPrompt(ctx, req.SkillID)

	history, err := s.loadHistory(ctx, req.ConversationID)
	if err != nil {
		return nil, fmt.Errorf("load history: %w", err)
	}

	messages := buildMessages(systemPrompt, history, req.Message)

	reply, toolCalls, err := s.llm.Chat(ctx, messages, s.tools.toolDefs())
	if err != nil {
		return nil, fmt.Errorf("llm chat: %w", err)
	}

	var toolResults []dto.ToolResultVO
	for len(toolCalls) > 0 {
		for _, tc := range toolCalls {
			res := s.tools.execute(ctx, tc)
			toolResults = append(toolResults, res)
			messages = append(messages, llm.Message{Role: "tool", ToolCallID: tc.ID, Content: res.Content})
		}
		reply, toolCalls, err = s.llm.Chat(ctx, messages, s.tools.toolDefs())
		if err != nil {
			return nil, fmt.Errorf("llm chat after tools: %w", err)
		}
	}

	s.saveMessages(ctx, req.ConversationID, req.Message, reply)

	return &dto.ChatResp{
		Reply:          reply,
		ToolCalls:      toolResults,
		ConversationID: req.ConversationID,
		Timestamp:      time.Now(),
	}, nil
}

// loadSkillPrompt fetches the SKILL.md from the skill service; falls back to a
// default assistant prompt when no skill is selected or it isn't ready.
func (s *AgentServiceImpl) loadSkillPrompt(ctx context.Context, skillID string) string {
	if skillID == "" {
		return defaultSystemPrompt
	}
	sk, err := s.skillSvc.Get(ctx, skillID)
	if err != nil || sk.Result == "" {
		return defaultSystemPrompt
	}
	return sk.Result
}

func (s *AgentServiceImpl) loadHistory(ctx context.Context, convID string) ([]llm.Message, error) {
	if convID == "" {
		return nil, nil
	}
	msgs, err := s.mapper.LoadHistory(ctx, convID, 20)
	if err != nil {
		return nil, err
	}
	out := make([]llm.Message, 0, len(msgs))
	for _, m := range msgs {
		out = append(out, llm.Message{Role: m.Role, Content: m.Content})
	}
	return out, nil
}

func (s *AgentServiceImpl) saveMessages(ctx context.Context, convID, userMsg, assistantMsg string) {
	if convID == "" {
		return
	}
	_ = s.mapper.Insert(ctx, convID, "user", userMsg)
	_ = s.mapper.Insert(ctx, convID, "assistant", assistantMsg)
}

func buildMessages(systemPrompt string, history []llm.Message, userMsg string) []llm.Message {
	var msgs []llm.Message
	if systemPrompt != "" {
		msgs = append(msgs, llm.Message{Role: "system", Content: systemPrompt})
	}
	msgs = append(msgs, history...)
	msgs = append(msgs, llm.Message{Role: "user", Content: userMsg})
	return msgs
}
