package service

import (
	"context"
	"mime/multipart"

	"github.com/Bruce-Sakura/UploadMyself/backend/pkg/file_uploads/dto"
)

// FileUploadService is the business-logic contract for file uploads.
type FileUploadService interface {
	// SaveUpload validates and stores a multipart file, persisting its metadata.
	SaveUpload(ctx context.Context, fh *multipart.FileHeader) (*dto.UploadResultVO, error)
	// Get returns stored-file metadata (including the on-disk path) for serving.
	Get(ctx context.Context, id string) (*dto.FileUploadVO, error)
	// ExtractCorpus runs text extraction (PDF/Word/OCR/plain) on an uploaded file.
	ExtractCorpus(ctx context.Context, fh *multipart.FileHeader) (*dto.CorpusResultVO, error)
}
