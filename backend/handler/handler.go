package handler

import (
	"net/http"

	"github.com/Bruce-Sakura/UploadMyself/backend/model"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Handler struct {
	db *gorm.DB
}

func New(db *gorm.DB) *Handler {
	return &Handler{db: db}
}

// ==================== Skills ====================

func (h *Handler) CreateSkill(c *gin.Context) {
	var body struct {
		Name   string `json:"name" binding:"required"`
		Corpus string `json:"corpus" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	s := model.Skill{
		ID:     uuid.New().String(),
		Name:   body.Name,
		Corpus: body.Corpus,
		Status: "pending",
	}
	h.db.Create(&s)
	c.JSON(http.StatusCreated, s)
}

func (h *Handler) ListSkills(c *gin.Context) {
	var skills []model.Skill
	h.db.Order("created_at desc").Find(&skills)
	c.JSON(http.StatusOK, skills)
}

func (h *Handler) GetSkill(c *gin.Context) {
	var s model.Skill
	if err := h.db.First(&s, "id = ?", c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, s)
}

func (h *Handler) UpdateSkill(c *gin.Context) {
	var s model.Skill
	if err := h.db.First(&s, "id = ?", c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	var body struct {
		Name   *string `json:"name"`
		Corpus *string `json:"corpus"`
		Status *string `json:"status"`
		Result *string `json:"result"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	updates := map[string]interface{}{}
	if body.Name != nil {
		updates["name"] = *body.Name
	}
	if body.Corpus != nil {
		updates["corpus"] = *body.Corpus
	}
	if body.Status != nil {
		updates["status"] = *body.Status
	}
	if body.Result != nil {
		updates["result"] = *body.Result
	}
	h.db.Model(&s).Updates(updates)
	c.JSON(http.StatusOK, s)
}

func (h *Handler) DeleteSkill(c *gin.Context) {
	if err := h.db.Delete(&model.Skill{}, "id = ?", c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"deleted": true})
}

// ==================== Voices ====================

func (h *Handler) CreateVoice(c *gin.Context) {
	var body struct {
		Name      string  `json:"name" binding:"required"`
		AudioPath string  `json:"audio_path"`
		Duration  float64 `json:"duration"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	v := model.Voice{
		ID:        uuid.New().String(),
		Name:      body.Name,
		AudioPath: body.AudioPath,
		Duration:  body.Duration,
		Status:    "pending",
	}
	h.db.Create(&v)
	c.JSON(http.StatusCreated, v)
}

func (h *Handler) ListVoices(c *gin.Context) {
	var voices []model.Voice
	h.db.Order("created_at desc").Find(&voices)
	c.JSON(http.StatusOK, voices)
}

func (h *Handler) GetVoice(c *gin.Context) {
	var v model.Voice
	if err := h.db.First(&v, "id = ?", c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, v)
}

func (h *Handler) DeleteVoice(c *gin.Context) {
	if err := h.db.Delete(&model.Voice{}, "id = ?", c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"deleted": true})
}

// ==================== Avatars ====================

func (h *Handler) CreateAvatar(c *gin.Context) {
	var body struct {
		Name      string `json:"name" binding:"required"`
		Type      string `json:"type" binding:"required"` // 2d | 3d
		PhotoPath string `json:"photo_path"`
		Style     string `json:"style"` // realistic | cartoon | anime
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	a := model.Avatar{
		ID:        uuid.New().String(),
		Name:      body.Name,
		Type:      body.Type,
		PhotoPath: body.PhotoPath,
		Style:     body.Style,
		Status:    "pending",
	}
	h.db.Create(&a)
	c.JSON(http.StatusCreated, a)
}

func (h *Handler) ListAvatars(c *gin.Context) {
	var avatars []model.Avatar
	q := h.db.Order("created_at desc")
	if t := c.Query("type"); t != "" {
		q = q.Where("type = ?", t)
	}
	q.Find(&avatars)
	c.JSON(http.StatusOK, avatars)
}

func (h *Handler) GetAvatar(c *gin.Context) {
	var a model.Avatar
	if err := h.db.First(&a, "id = ?", c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, a)
}

func (h *Handler) DeleteAvatar(c *gin.Context) {
	if err := h.db.Delete(&model.Avatar{}, "id = ?", c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"deleted": true})
}

// ==================== Tasks ====================

func (h *Handler) GetTask(c *gin.Context) {
	var t model.Task
	if err := h.db.First(&t, "id = ?", c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, t)
}
