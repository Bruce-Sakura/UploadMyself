package dto

import "time"

// AvatarVO is the avatar representation returned to clients.
type AvatarVO struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	Type       string    `json:"type"`
	PhotoPath  string    `json:"photo_path"`
	Style      string    `json:"style"`
	Status     string    `json:"status"`
	Result     string    `json:"result"`
	OutputPath string    `json:"output_path"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// CreateAvatarReq creates a new avatar record.
type CreateAvatarReq struct {
	Name      string `json:"name" binding:"required"`
	Type      string `json:"type" binding:"required"` // 2d | 3d
	PhotoPath string `json:"photo_path"`
	Style     string `json:"style"`
}

// ListAvatarReq filters the avatar list.
type ListAvatarReq struct {
	Type string `json:"type"`
}
