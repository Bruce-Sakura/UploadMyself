package mapper

import (
	"context"

	"github.com/Bruce-Sakura/UploadMyself/backend/pkg/file_uploads/entity"
	"github.com/jackc/pgx/v5/pgxpool"
)

// FileUploadMapper is the data-access layer for the file_uploads table.
type FileUploadMapper struct {
	pool *pgxpool.Pool
}

func NewFileUploadMapper(pool *pgxpool.Pool) *FileUploadMapper {
	return &FileUploadMapper{pool: pool}
}

func (m *FileUploadMapper) Insert(ctx context.Context, f *entity.FileUpload) error {
	_, err := m.pool.Exec(ctx,
		`INSERT INTO file_uploads (id, original_name, stored_path, size, mime_type)
		 VALUES ($1, $2, $3, $4, $5)`,
		f.ID, f.OriginalName, f.StoredPath, f.Size, f.MimeType)
	return err
}

func (m *FileUploadMapper) GetByID(ctx context.Context, id string) (*entity.FileUpload, error) {
	row := m.pool.QueryRow(ctx,
		`SELECT id, original_name, stored_path, size, mime_type, created_at
		 FROM file_uploads WHERE id = $1`, id)
	var f entity.FileUpload
	if err := row.Scan(&f.ID, &f.OriginalName, &f.StoredPath, &f.Size, &f.MimeType, &f.CreatedAt); err != nil {
		return nil, err
	}
	return &f, nil
}
