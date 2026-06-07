package entity

import "time"

// Voice maps the voices table — voice clone.
type Voice struct {
	ID           string
	Name         string
	AudioPath    string
	Duration     float64 // seconds
	ModelPath    string
	RefAudioPath string
	Status       string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
