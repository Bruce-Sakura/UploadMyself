package handler

import (
	"net/http"

	"github.com/Bruce-Sakura/UploadMyself/backend/pkg/skills/dto"
	"github.com/Bruce-Sakura/UploadMyself/backend/pkg/skills/service"
	"github.com/gin-gonic/gin"
)

type SkillHandler struct {
	svc service.SkillService
}

func NewSkillHandler(svc service.SkillService) *SkillHandler {
	return &SkillHandler{svc: svc}
}

// Register mounts skill routes under the given group (e.g. /api/v1).
func (h *SkillHandler) Register(rg *gin.RouterGroup) {
	g := rg.Group("/skills")
	g.POST("", h.Create)
	g.POST("/import", h.Import)
	g.GET("", h.List)
	g.GET("/:id", h.Get)
	g.PUT("/:id", h.Update)
	g.DELETE("/:id", h.Delete)
	g.POST("/:id/process", h.Process)
}

func (h *SkillHandler) Create(c *gin.Context) {
	var req dto.CreateSkillReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	sk, err := h.svc.Create(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, sk)
}

func (h *SkillHandler) List(c *gin.Context) {
	sks, err := h.svc.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, sks)
}

func (h *SkillHandler) Get(c *gin.Context) {
	sk, err := h.svc.Get(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, sk)
}

func (h *SkillHandler) Update(c *gin.Context) {
	var req dto.UpdateSkillReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	sk, err := h.svc.Update(c.Request.Context(), c.Param("id"), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, sk)
}

func (h *SkillHandler) Delete(c *gin.Context) {
	if err := h.svc.Delete(c.Request.Context(), c.Param("id")); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"deleted": true})
}

func (h *SkillHandler) Import(c *gin.Context) {
	var req dto.ImportSkillReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	sk, err := h.svc.Import(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, sk)
}

func (h *SkillHandler) Process(c *gin.Context) {
	taskID, err := h.svc.Process(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "skill not found"})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{"id": taskID, "status": "processing"})
}
