package mapper

import (
	"context"
	"database/sql"

	"github.com/Bruce-Sakura/UploadMyself/backend/pkg/file_uploads/entity"
)

// FileUploadMapper is the data-access layer for the file_uploads table.
type FileUploadMapper struct {
	db *sql.DB
}

func NewFileUploadMapper(db *sql.DB) *FileUploadMapper {
	return &FileUploadMapper{db: db}
}

func (m *FileUploadMapper) Insert(ctx context.Context, f *entity.FileUpload) error {
	_, err := m.db.ExecContext(ctx,
		`INSERT INTO file_uploads (id, original_name, stored_path, size, mime_type)
		 VALUES (?, ?, ?, ?, ?)`,
		f.ID, f.OriginalName, f.StoredPath, f.Size, f.MimeType)
	return err
}

func (m *FileUploadMapper) GetByID(ctx context.Context, id string) (*entity.FileUpload, error) {
	row := m.db.QueryRowContext(ctx,
		`SELECT id, original_name, stored_path, size, mime_type, created_at
		 FROM file_uploads WHERE id = ?`, id)
	var f entity.FileUpload
	if err := row.Scan(&f.ID, &f.OriginalName, &f.StoredPath, &f.Size, &f.MimeType, &f.CreatedAt); err != nil {
		return nil, err
	}
	return &f, nil
}
