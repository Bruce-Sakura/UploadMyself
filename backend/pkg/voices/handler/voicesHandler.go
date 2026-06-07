package handler

import (
	"net/http"

	"github.com/Bruce-Sakura/UploadMyself/backend/pkg/voices/dto"
	"github.com/Bruce-Sakura/UploadMyself/backend/pkg/voices/service"
	"github.com/gin-gonic/gin"
)

type VoiceHandler struct {
	svc service.VoiceService
}

func NewVoiceHandler(svc service.VoiceService) *VoiceHandler {
	return &VoiceHandler{svc: svc}
}

// Register mounts voice routes under the given group (e.g. /api/v1).
func (h *VoiceHandler) Register(rg *gin.RouterGroup) {
	g := rg.Group("/voices")
	g.POST("", h.Create)
	g.GET("", h.List)
	g.GET("/:id", h.Get)
	g.DELETE("/:id", h.Delete)
	g.POST("/:id/train", h.Train)
	g.POST("/:id/synthesize", h.Synthesize)
}

func (h *VoiceHandler) Create(c *gin.Context) {
	var req dto.CreateVoiceReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	v, err := h.svc.Create(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, v)
}

func (h *VoiceHandler) List(c *gin.Context) {
	vs, err := h.svc.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, vs)
}

func (h *VoiceHandler) Get(c *gin.Context) {
	v, err := h.svc.Get(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, v)
}

func (h *VoiceHandler) Delete(c *gin.Context) {
	if err := h.svc.Delete(c.Request.Context(), c.Param("id")); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"deleted": true})
}

func (h *VoiceHandler) Train(c *gin.Context) {
	taskID, err := h.svc.Train(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "voice not found"})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{"id": taskID, "status": "training"})
}

func (h *VoiceHandler) Synthesize(c *gin.Context) {
	var req dto.SynthesizeReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	audioPath, err := h.svc.Synthesize(c.Request.Context(), c.Param("id"), req.Text)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"audio_path": audioPath})
}
