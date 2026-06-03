package handler

import (
	"net/http"

	"github.com/Bruce-Sakura/UploadMyself/backend/agent"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// AgentChat 核心对话端点
func (h *Handler) AgentChat(c *gin.Context) {
	var req agent.ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 自动生成会话 ID
	if req.ConversationID == "" {
		req.ConversationID = uuid.New().String()
	}

	resp, err := h.agent.Chat(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"reply":           resp.Reply,
		"tool_calls":      resp.ToolCalls,
		"conversation_id": req.ConversationID,
		"timestamp":       resp.Timestamp,
	})
}

// ListTools 列出可用工具
func (h *Handler) ListTools(c *gin.Context) {
	// TODO: 从 agent 获取工具列表
	c.JSON(http.StatusOK, gin.H{
		"tools": []gin.H{
			{"name": "shell", "description": "执行 Shell 命令"},
			{"name": "file_read", "description": "读取文件"},
			{"name": "file_write", "description": "写入文件"},
			{"name": "web_search", "description": "搜索互联网"},
		},
	})
}
