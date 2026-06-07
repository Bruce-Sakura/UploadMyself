package dto

import "time"

// FileUploadVO is the stored-file representation (used for serving).
type FileUploadVO struct {
	ID           string    `json:"id"`
	OriginalName string    `json:"original_name"`
	StoredPath   string    `json:"stored_path"`
	Size         int64     `json:"size"`
	MimeType     string    `json:"mime_type"`
	CreatedAt    time.Time `json:"created_at"`
}

// UploadResultVO is returned after a successful upload.
type UploadResultVO struct {
	ID       string `json:"id"`
	Filename string `json:"filename"`
	Path     string `json:"path"`
	Size     int64  `json:"size"`
}

// CorpusResultVO is the extracted-text result from a corpus file.
type CorpusResultVO struct {
	Text   string `json:"text"`
	Method string `json:"method"`
	Name   string `json:"name"`
}
