package entity

import "time"

// Avatar maps the avatars table — 2D/3D virtual avatar.
type Avatar struct {
	ID         string
	Name       string
	Type       string // 2d | 3d
	PhotoPath  string
	Style      string // realistic | cartoon | anime
	Status     string // pending | processing | done | failed
	Result     string // output JSON (cartoon/views/glb/obj URLs)
	OutputPath string // quick-access output file
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
