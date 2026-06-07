package handler

import (
	"net/http"

	"github.com/Bruce-Sakura/UploadMyself/backend/pkg/messages/dto"
	"github.com/Bruce-Sakura/UploadMyself/backend/pkg/messages/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AgentHandler struct {
	svc service.AgentService
}

func NewAgentHandler(svc service.AgentService) *AgentHandler {
	return &AgentHandler{svc: svc}
}

// Register mounts agent routes under the given group (e.g. /api/v1).
func (h *AgentHandler) Register(rg *gin.RouterGroup) {
	g := rg.Group("/agent")
	g.POST("/chat", h.Chat)
	g.GET("/tools", h.ListTools)
}

func (h *AgentHandler) Chat(c *gin.Context) {
	var req dto.ChatReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.ConversationID == "" {
		req.ConversationID = uuid.NewString()
	}
	resp, err := h.svc.Chat(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *AgentHandler) ListTools(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"tools": h.svc.ListTools()})
}
