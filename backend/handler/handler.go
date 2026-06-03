package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"time"

	"github.com/Bruce-Sakura/UploadMyself/backend/agent"
	"github.com/Bruce-Sakura/UploadMyself/backend/model"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// MLScriptsDir and PythonBin are configured from main via env vars.
var (
	MLScriptsDir = "../ml/scripts"
	PythonBin    = "python3"
)

type Handler struct {
	db    *gorm.DB
	agent *agent.Agent
}

func New(db *gorm.DB, agt *agent.Agent) *Handler {
	return &Handler{db: db, agent: agt}
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

// ==================== Voice Processing ====================

// TrainVoice triggers voice_clone_train.py for the given voice ID.
func (h *Handler) TrainVoice(c *gin.Context) {
	voiceID := c.Param("id")
	var v model.Voice
	if err := h.db.First(&v, "id = ?", voiceID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "voice not found"})
		return
	}

	// Create task
	task, err := CreateTask(h.db, "voice_train", voiceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create task"})
		return
	}

	// Update voice status
	h.db.Model(&v).Update("status", "training")

	// Run async
	go func() {
		UpdateTaskStatus(h.db, task.ID, "running", 10, "")

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
		defer cancel()

		script := fmt.Sprintf("%s/voice_clone_train.py", MLScriptsDir)
		cmd := exec.CommandContext(ctx, PythonBin, script,
			"--voice-id", voiceID,
			"--audio-path", v.AudioPath,
			"--name", v.Name,
		)

		out, err := cmd.Output()
		if err != nil {
			UpdateTaskStatus(h.db, task.ID, "failed", 0, fmt.Sprintf("train error: %v", err))
			h.db.Model(&v).Update("status", "failed")
			return
		}

		// Parse JSON stdout
		var result map[string]interface{}
		if json.Unmarshal(out, &result) == nil {
			if mp, ok := result["model_path"]; ok {
				h.db.Model(&v).Update("model_path", fmt.Sprintf("%v", mp))
			}
		}

		h.db.Model(&v).Update("status", "done")
		UpdateTaskStatus(h.db, task.ID, "done", 100, "")
	}()

	c.JSON(http.StatusAccepted, gin.H{"task_id": task.ID, "status": "training"})
}

// SynthesizeVoice calls voice_synthesize.py and returns audio.
func (h *Handler) SynthesizeVoice(c *gin.Context) {
	voiceID := c.Param("id")
	var v model.Voice
	if err := h.db.First(&v, "id = ?", voiceID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "voice not found"})
		return
	}

	var body struct {
		Text string `json:"text" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	script := fmt.Sprintf("%s/voice_synthesize.py", MLScriptsDir)
	cmd := exec.CommandContext(ctx, PythonBin, script,
		"--voice-id", voiceID,
		"--model-path", v.ModelPath,
		"--ref-audio", v.RefAudioPath,
		"--text", body.Text,
	)

	out, err := cmd.Output()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("synthesize error: %v", err)})
		return
	}

	var result map[string]interface{}
	if err := json.Unmarshal(out, &result); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid script output"})
		return
	}

	audioPath, _ := result["audio_path"]
	c.JSON(http.StatusOK, gin.H{"audio_path": audioPath, "result": result})
}

// ==================== Skill Processing ====================

// ProcessSkill triggers analyze_corpus.py for the given skill.
func (h *Handler) ProcessSkill(c *gin.Context) {
	skillID := c.Param("id")
	var s model.Skill
	if err := h.db.First(&s, "id = ?", skillID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "skill not found"})
		return
	}

	task, err := CreateTask(h.db, "skill_process", skillID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create task"})
		return
	}

	h.db.Model(&s).Update("status", "processing")

	go func() {
		UpdateTaskStatus(h.db, task.ID, "running", 10, "")

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
		defer cancel()

		script := fmt.Sprintf("%s/analyze_corpus.py", MLScriptsDir)
		cmd := exec.CommandContext(ctx, PythonBin, script,
			"--skill-id", skillID,
			"--corpus", s.Corpus,
			"--name", s.Name,
		)

		out, err := cmd.Output()
		if err != nil {
			UpdateTaskStatus(h.db, task.ID, "failed", 0, fmt.Sprintf("process error: %v", err))
			h.db.Model(&s).Updates(map[string]interface{}{"status": "failed"})
			return
		}

		var result map[string]interface{}
		if json.Unmarshal(out, &result) == nil {
			if sk, ok := result["skill_md"]; ok {
				h.db.Model(&s).Update("result", fmt.Sprintf("%v", sk))
			}
		}

		h.db.Model(&s).Updates(map[string]interface{}{"status": "done"})
		UpdateTaskStatus(h.db, task.ID, "done", 100, "")
	}()

	c.JSON(http.StatusAccepted, gin.H{"task_id": task.ID, "status": "processing"})
}

// ==================== Avatar Processing ====================

// ProcessAvatar triggers detect_face.py for the given avatar.
func (h *Handler) ProcessAvatar(c *gin.Context) {
	avatarID := c.Param("id")
	a := model.Avatar{}
	if err := h.db.First(&a, "id = ?", avatarID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "avatar not found"})
		return
	}

	task, err := CreateTask(h.db, "avatar_process", avatarID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create task"})
		return
	}

	h.db.Model(&a).Update("status", "processing")

	go func() {
		UpdateTaskStatus(h.db, task.ID, "running", 10, "")

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
		defer cancel()

		script := fmt.Sprintf("%s/detect_face.py", MLScriptsDir)
		cmd := exec.CommandContext(ctx, PythonBin, script,
			"--avatar-id", avatarID,
			"--photo-path", a.PhotoPath,
			"--type", a.Type,
			"--style", a.Style,
		)

		out, err := cmd.Output()
		if err != nil {
			UpdateTaskStatus(h.db, task.ID, "failed", 0, fmt.Sprintf("process error: %v", err))
			h.db.Model(&a).Update("status", "failed")
			return
		}

		var result map[string]interface{}
		if json.Unmarshal(out, &result) == nil {
			if op, ok := result["output_path"]; ok {
				h.db.Model(&a).Update("output_path", fmt.Sprintf("%v", op))
			}
			if r, ok := result["result"]; ok {
				h.db.Model(&a).Update("result", fmt.Sprintf("%v", r))
			}
		}

		h.db.Model(&a).Update("status", "done")
		UpdateTaskStatus(h.db, task.ID, "done", 100, "")
	}()

	c.JSON(http.StatusAccepted, gin.H{"task_id": task.ID, "status": "processing"})
}
