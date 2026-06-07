package dto

import "time"

// VoiceVO is the voice representation returned to clients.
type VoiceVO struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	AudioPath    string    `json:"audio_path"`
	Duration     float64   `json:"duration"`
	ModelPath    string    `json:"model_path"`
	RefAudioPath string    `json:"ref_audio_path"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// CreateVoiceReq creates a new voice record.
type CreateVoiceReq struct {
	Name      string  `json:"name" binding:"required"`
	AudioPath string  `json:"audio_path"`
	Duration  float64 `json:"duration"`
}

// SynthesizeReq requests TTS synthesis from a trained voice.
type SynthesizeReq struct {
	Text string `json:"text" binding:"required"`
}
