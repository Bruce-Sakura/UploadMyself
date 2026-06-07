package handler

import (
	"net/http"

	"github.com/Bruce-Sakura/UploadMyself/backend/pkg/tasks/dto"
	"github.com/Bruce-Sakura/UploadMyself/backend/pkg/tasks/service"
	"github.com/gin-gonic/gin"
)

type TaskHandler struct {
	svc service.TaskService
}

func NewTaskHandler(svc service.TaskService) *TaskHandler {
	return &TaskHandler{svc: svc}
}

// Register mounts task routes under the given group (e.g. /api/v1).
func (h *TaskHandler) Register(rg *gin.RouterGroup) {
	rg.GET("/tasks", h.List)
	rg.GET("/tasks/:id", h.Get)
}

func (h *TaskHandler) Get(c *gin.Context) {
	t, err := h.svc.Get(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, t)
}

func (h *TaskHandler) List(c *gin.Context) {
	ts, err := h.svc.List(c.Request.Context(), dto.ListTaskReq{
		Type:  c.Query("type"),
		RefID: c.Query("ref_id"),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, ts)
}
