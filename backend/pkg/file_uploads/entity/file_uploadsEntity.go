package entity

import "time"

// FileUpload maps the file_uploads table — uploaded file metadata.
type FileUpload struct {
	ID           string
	OriginalName string
	StoredPath   string
	Size         int64
	MimeType     string
	CreatedAt    time.Time
}
