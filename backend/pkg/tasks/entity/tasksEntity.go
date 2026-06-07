package entity

import "time"

// Task maps the tasks table — async task tracking.
type Task struct {
	ID        string
	Type      string // skill_process | voice_train | avatar_process
	RefID     string
	Status    string // pending | running | done | failed
	Progress  int    // 0-100
	Error     string
	CreatedAt time.Time
	UpdatedAt time.Time
}
