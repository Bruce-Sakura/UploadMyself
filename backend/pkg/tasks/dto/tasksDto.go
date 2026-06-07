package dto

import "time"

// TaskVO is the task representation returned to clients.
type TaskVO struct {
	ID        string    `json:"id"`
	Type      string    `json:"type"`
	RefID     string    `json:"ref_id"`
	Status    string    `json:"status"`
	Progress  int       `json:"progress"`
	Error     string    `json:"error"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateTaskReq creates a new async task.
type CreateTaskReq struct {
	Type  string `json:"type"`
	RefID string `json:"ref_id"`
}

// ListTaskReq filters the task list.
type ListTaskReq struct {
	Type  string `json:"type"`
	RefID string `json:"ref_id"`
}
