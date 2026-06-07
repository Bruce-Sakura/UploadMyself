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

// ImportSkillReq imports a ready-made SKILL.md from a URL (GitHub/raw/hub) or
// inline content. Exactly one of URL / Content should be provided.
type ImportSkillReq struct {
	Name    string `json:"name"`
	URL     string `json:"url"`
	Content string `json:"content"`
}
