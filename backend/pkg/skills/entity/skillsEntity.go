package entity

import "time"

// Skill maps the skills table — thinking-framework clone.
type Skill struct {
	ID        string
	Name      string
	Corpus    string
	Status    string // pending | processing | done | failed
	Result    string // generated SKILL.md
	CreatedAt time.Time
	UpdatedAt time.Time
}
