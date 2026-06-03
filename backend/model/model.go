package model

import (
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Connect opens DB and auto-migrates all models.
func Connect(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	if err := db.AutoMigrate(&Skill{}, &Voice{}, &Avatar{}, &Task{}, &FileUpload{}); err != nil {
		return nil, err
	}
	return db, nil
}

// Skill — thinking framework clone
type Skill struct {
	ID        string         `json:"id" gorm:"primaryKey"`
	Name      string         `json:"name" gorm:"size:128;not null"`
	Corpus    string         `json:"corpus" gorm:"type:text"`
	Status    string         `json:"status" gorm:"size:32;default:pending"` // pending | processing | done | failed
	Result    string         `json:"result" gorm:"type:text"`               // generated SKILL.md
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

// Voice — voice clone
type Voice struct {
	ID          string         `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"size:128;not null"`
	AudioPath   string         `json:"audio_path" gorm:"size:512"`
	Duration    float64        `json:"duration"` // seconds
	ModelPath   string         `json:"model_path" gorm:"size:512"`   // trained model path
	RefAudioPath string        `json:"ref_audio_path" gorm:"size:512"` // reference audio for synthesis
	Status      string         `json:"status" gorm:"size:32;default:pending"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

// Avatar — 2D/3D virtual avatar
type Avatar struct {
	ID         string         `json:"id" gorm:"primaryKey"`
	Name       string         `json:"name" gorm:"size:128;not null"`
	Type       string         `json:"type" gorm:"size:8;not null"` // 2d | 3d
	PhotoPath  string         `json:"photo_path" gorm:"size:512"`
	Style      string         `json:"style" gorm:"size:32"` // realistic | cartoon | anime
	Status     string         `json:"status" gorm:"size:32;default:pending"`
	Result     string         `json:"result" gorm:"type:text"`   // output file path
	OutputPath string         `json:"output_path" gorm:"size:512"` // generated output file
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
}

// Task — async task tracking
type Task struct {
	ID        string    `json:"id" gorm:"primaryKey"`
	Type      string    `json:"type" gorm:"size:32;not null"` // skill | voice | avatar_2d | avatar_3d
	RefID     string    `json:"ref_id" gorm:"size:64"`       // related entity ID
	Status    string    `json:"status" gorm:"size:32;default:pending"`
	Progress  int       `json:"progress"` // 0-100
	Error     string    `json:"error" gorm:"type:text"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// FileUpload — uploaded file metadata
type FileUpload struct {
	ID           string    `json:"id" gorm:"primaryKey"`
	OriginalName string    `json:"original_name" gorm:"size:256;not null"`
	StoredPath   string    `json:"stored_path" gorm:"size:512;not null"`
	Size         int64     `json:"size"`
	MimeType     string    `json:"mime_type" gorm:"size:128"`
	CreatedAt    time.Time `json:"created_at"`
}
