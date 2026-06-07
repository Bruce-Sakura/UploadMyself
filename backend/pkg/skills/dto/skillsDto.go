package dto

import "time"

// SkillVO is the skill representation returned to clients.
type SkillVO struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Corpus    string    `json:"corpus"`
	Status    string    `json:"status"`
	Result    string    `json:"result"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateSkillReq creates a new skill from a text corpus.
type CreateSkillReq struct {
	Name   string `json:"name" binding:"required"`
	Corpus string `json:"corpus" binding:"required"`
}

// UpdateSkillReq partially updates a skill (nil = unchanged).
type UpdateSkillReq struct {
	Name   *string `json:"name"`
	Corpus *string `json:"corpus"`
	Status *string `json:"status"`
	Result *string `json:"result"`
}
