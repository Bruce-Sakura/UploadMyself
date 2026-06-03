package models

import (
	"time"
)

// ==================== 通用模型 ====================

type TaskStatus string

const (
	TaskPending   TaskStatus = "pending"
	TaskRunning   TaskStatus = "running"
	TaskCompleted TaskStatus = "completed"
	TaskFailed    TaskStatus = "failed"
)

// Task 异步任务
type Task struct {
	ID        string     `json:"id"`
	Type      string     `json:"type"` // skill | voice | avatar_2d | avatar_3d | distill
	Status    TaskStatus `json:"status"`
	Progress  int        `json:"progress"` // 0-100
	Result    string     `json:"result,omitempty"`
	Error     string     `json:"error,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

// ==================== Skill ====================

type Skill struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CorpusFile string   `json:"corpus_file"`
	Status    TaskStatus `json:"status"`
	SkillMD   string    `json:"skill_md,omitempty"` // 生成的 SKILL.md 内容
	Models    int       `json:"models"`             // 心智模型数量
	CreatedAt time.Time `json:"created_at"`
}

// ==================== Voice ====================

type Voice struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	AudioFile   string    `json:"audio_file"`
	Status      TaskStatus `json:"status"`
	ModelPath   string    `json:"model_path,omitempty"`
	SampleRate  int       `json:"sample_rate"`
	DurationSec float64   `json:"duration_sec"`
	CreatedAt   time.Time `json:"created_at"`
}

// ==================== Avatar ====================

type Avatar struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	PhotoFile string    `json:"photo_file"`
	Type      string    `json:"type"` // 2d | 3d
	Status    TaskStatus `json:"status"`
	Style     string    `json:"style,omitempty"`
	ModelPath string    `json:"model_path,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

// ==================== Distill ====================

type DistillJob struct {
	ID            string    `json:"id"`
	TeacherModel  string    `json:"teacher_model"`
	StudentModel  string    `json:"student_model"`
	TaskType      string    `json:"task_type"` // llm | voice | avatar_2d
	Status        TaskStatus `json:"status"`
	Epoch         int       `json:"epoch"`
	TotalEpochs   int       `json:"total_epochs"`
	CurrentLoss   float64   `json:"current_loss"`
	CreatedAt     time.Time `json:"created_at"`
}
