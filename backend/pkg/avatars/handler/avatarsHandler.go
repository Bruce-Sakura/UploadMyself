package handler

import (
	"net/http"

	"github.com/Bruce-Sakura/UploadMyself/backend/pkg/avatars/dto"
	"github.com/Bruce-Sakura/UploadMyself/backend/pkg/avatars/service"
	"github.com/gin-gonic/gin"
)

type AvatarHandler struct {
	svc service.AvatarService
}

func NewAvatarHandler(svc service.AvatarService) *AvatarHandler {
	return &AvatarHandler{svc: svc}
}

// Register mounts avatar routes under the given group (e.g. /api/v1).
func (h *AvatarHandler) Register(rg *gin.RouterGroup) {
	g := rg.Group("/avatars")
	g.POST("", h.Create)
	g.GET("", h.List)
	g.GET("/:id", h.Get)
	g.DELETE("/:id", h.Delete)
	g.POST("/:id/process", h.Process)
}

func (h *AvatarHandler) Create(c *gin.Context) {
	var req dto.CreateAvatarReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	a, err := h.svc.Create(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, a)
}

func (h *AvatarHandler) List(c *gin.Context) {
	as, err := h.svc.List(c.Request.Context(), dto.ListAvatarReq{Type: c.Query("type")})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, as)
}

func (h *AvatarHandler) Get(c *gin.Context) {
	a, err := h.svc.Get(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, a)
}

func (h *AvatarHandler) Delete(c *gin.Context) {
	if err := h.svc.Delete(c.Request.Context(), c.Param("id")); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"deleted": true})
}

func (h *AvatarHandler) Process(c *gin.Context) {
	taskID, err := h.svc.Process(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "avatar not found"})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{"id": taskID, "status": "processing"})
}
