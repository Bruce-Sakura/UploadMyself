package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"io"
	"os/exec"
	"strings"
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
	id := c.Param("id")
	var s model.Skill
	if err := h.db.First(&s, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	// 删除关联文件
	outputDir := fmt.Sprintf("%s/%s_output", UploadDir, id)
	os.RemoveAll(outputDir)
	h.db.Where("ref_id = ?", id).Delete(&model.Task{})
	h.db.Delete(&model.Skill{}, "id = ?", id)
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
	id := c.Param("id")
	var v model.Voice
	if err := h.db.First(&v, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	// 删除关联文件
	if v.AudioPath != "" {
		os.Remove(v.AudioPath)
	}
	if v.ModelPath != "" {
		os.Remove(v.ModelPath)
	}
	outputDir := fmt.Sprintf("%s/%s_output", UploadDir, id)
	os.RemoveAll(outputDir)
	h.db.Where("ref_id = ?", id).Delete(&model.Task{})
	h.db.Delete(&model.Voice{}, "id = ?", id)
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
	id := c.Param("id")
	var a model.Avatar
	if err := h.db.First(&a, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	// 删除所有关联文件
	if a.PhotoPath != "" {
		os.Remove(a.PhotoPath)
	}
	if a.OutputPath != "" {
		os.Remove(a.OutputPath)
	}
	outputDir := fmt.Sprintf("%s/%s_output", UploadDir, id)
	os.RemoveAll(outputDir)
	h.db.Where("ref_id = ?", id).Delete(&model.Task{})
	h.db.Delete(&model.Avatar{}, "id = ?", id)
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
		

		// Write config JSON for voice_clone_train.py
		configPath := fmt.Sprintf("%s/%s_train_config.json", UploadDir, voiceID)
		configData := fmt.Sprintf(`{"voice_id":"%s","audio_path":"%s","name":"%s"}`, voiceID, v.AudioPath, v.Name)
		os.WriteFile(configPath, []byte(configData), 0644)

		script := fmt.Sprintf("%s/voice_clone_train.py", MLScriptsDir)
		cmd := exec.CommandContext(ctx, PythonBin, script,
			"--config", configPath,
		)

		out, err := cmd.CombinedOutput()
		if err != nil {
			UpdateTaskStatus(h.db, task.ID, "failed", 0, fmt.Sprintf("train error: %v\noutput: %s", err, string(out)))
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
	

	outputPath := fmt.Sprintf("%s/%s_synth.wav", UploadDir, voiceID)

	script := fmt.Sprintf("%s/voice_synthesize.py", MLScriptsDir)
	cmd := exec.CommandContext(ctx, PythonBin, script,
		"--voice-id", voiceID,
		"--text", body.Text,
		"--output", outputPath,
		"--model-dir", UploadDir,
	)

	out, err := cmd.CombinedOutput()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("synthesize error: %v\noutput: %s", err, string(out))})
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

  ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
  defer cancel()
		

		// 用 MiMo 生成 SKILL.md
		prompt := fmt.Sprintf(`你是一个思维框架分析专家。分析以下用户的文本语料，生成一个完整的 SKILL.md 文件。

用户名称: %s

文本语料:
%s

请生成 SKILL.md，包含：
1. 身份卡（用第一人称写50字自我介绍）
2. 核心智模型（3-5个，含证据和应用）
3. 决策启发式（3-5条）
4. 表达DNA（句式/词汇/幽默风格）
5. 诚实边界（局限性）

用 Markdown 格式输出。`, s.Name, s.Corpus)

		llmResp, llmErr := h.agent.LLMChat(ctx, prompt)
		if llmErr != nil {
			UpdateTaskStatus(h.db, task.ID, "failed", 0, fmt.Sprintf("LLM error: %v", llmErr))
			h.db.Model(&s).Updates(map[string]interface{}{"status": "failed"})
			return
		}

		h.db.Model(&s).Update("result", llmResp)
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

		

		// CharacterGen: photo → 4-view → 3D mesh → GLB
		outputDir := fmt.Sprintf("%s/%s_output", UploadDir, avatarID)
		mlServiceURL := os.Getenv("ML_SERVICE_URL")
		if mlServiceURL == "" {
			mlServiceURL = "http://host.docker.internal:8001"
		}

		// 调用 ML 服务
		reqBody, _ := json.Marshal(map[string]interface{}{
			"input_path": a.PhotoPath,
			"output_dir": outputDir,
			"seed":       2333,
			"timestep":   40,
		})
		resp, err := http.Post(mlServiceURL+"/generate-avatar", "application/json", strings.NewReader(string(reqBody)))
		if err != nil {
			UpdateTaskStatus(h.db, task.ID, "failed", 0, fmt.Sprintf("ML service error: %v", err))
			h.db.Model(&a).Update("status", "failed")
			return
		}
		defer resp.Body.Close()

		respBody, _ := io.ReadAll(resp.Body)
		var result map[string]interface{}
		json.Unmarshal(respBody, &result)

		if resp.StatusCode != 200 {
			errMsg := "unknown error"
			if e, ok := result["error"]; ok {
				errMsg = fmt.Sprintf("%v", e)
			}
			UpdateTaskStatus(h.db, task.ID, "failed", 0, fmt.Sprintf("ML error: %s", errMsg))
			h.db.Model(&a).Update("status", "failed")
			return
		}

		// Convert paths to relative URLs
		for _, key := range []string{"cartoon_image", "skeleton_image", "animation_data"} {
			if v, ok := result[key]; ok {
				result[key] = strings.Replace(fmt.Sprintf("%v", v), UploadDir, "uploads", 1)
			}
		}
		// Store full JSON as result (contains cartoon + skeleton + animation paths)
		animJSON, _ := json.Marshal(result)
		h.db.Model(&a).Update("result", string(animJSON))
		// Store cartoon as output_path for quick access
		if cartoon, ok := result["cartoon_image"]; ok {
			h.db.Model(&a).Update("output_path", fmt.Sprintf("%v", cartoon))
		}

		h.db.Model(&a).Update("status", "done")
		UpdateTaskStatus(h.db, task.ID, "done", 100, "")
	}()

	c.JSON(http.StatusAccepted, gin.H{"task_id": task.ID, "status": "processing"})
}
