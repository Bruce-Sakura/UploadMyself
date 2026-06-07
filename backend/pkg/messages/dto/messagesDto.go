package dto

import "time"

// ChatReq is the agent chat request.
type ChatReq struct {
	ConversationID string `json:"conversation_id"`
	SkillID        string `json:"skill_id"`
	Message        string `json:"message" binding:"required"`
}

// ToolResultVO is a single executed tool-call result.
type ToolResultVO struct {
	ToolCallID string `json:"tool_call_id"`
	Content    string `json:"content"`
}

// ChatResp is the agent chat response.
type ChatResp struct {
	Reply          string         `json:"reply"`
	ToolCalls      []ToolResultVO `json:"tool_calls,omitempty"`
	ConversationID string         `json:"conversation_id"`
	Timestamp      time.Time      `json:"timestamp"`
}

// ToolInfoVO describes an available agent tool.
type ToolInfoVO struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}
