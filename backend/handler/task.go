package handler

import (
	"net/http"

	"github.com/Bruce-Sakura/UploadMyself/backend/model"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// CreateTask creates a new async task record and returns it.
func CreateTask(db *gorm.DB, taskType, refID string) (*model.Task, error) {
	t := model.Task{
		ID:       uuid.New().String(),
		Type:     taskType,
		RefID:    refID,
		Status:   "pending",
		Progress: 0,
	}
	if err := db.Create(&t).Error; err != nil {
		return nil, err
	}
	return &t, nil
}

// UpdateTaskStatus updates a task's status, progress, and optional error message.
func UpdateTaskStatus(db *gorm.DB, taskID, status string, progress int, errMsg string) error {
	updates := map[string]interface{}{
		"status":   status,
		"progress": progress,
	}
	if errMsg != "" {
		updates["error"] = errMsg
	}
	return db.Model(&model.Task{}).Where("id = ?", taskID).Updates(updates).Error
}

// ListTasks returns tasks, optionally filtered by type query param.
func (h *Handler) ListTasks(c *gin.Context) {
	var tasks []model.Task
	q := h.db.Order("created_at desc")
	if t := c.Query("type"); t != "" {
		q = q.Where("type = ?", t)
	}
	if refID := c.Query("ref_id"); refID != "" {
		q = q.Where("ref_id = ?", refID)
	}
	q.Find(&tasks)
	c.JSON(http.StatusOK, tasks)
}
