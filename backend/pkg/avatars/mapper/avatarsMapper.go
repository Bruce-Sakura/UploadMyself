package mapper

import (
	"context"
	"database/sql"

	"github.com/Bruce-Sakura/UploadMyself/backend/pkg/avatars/entity"
)

// AvatarMapper is the data-access layer for the avatars table.
type AvatarMapper struct {
	db *sql.DB
}

func NewAvatarMapper(db *sql.DB) *AvatarMapper {
	return &AvatarMapper{db: db}
}

const avatarCols = `id, name, type, photo_path, style, status, result, output_path, created_at, updated_at`

func scanAvatar(s interface {
	Scan(dest ...any) error
}) (entity.Avatar, error) {
	var a entity.Avatar
	err := s.Scan(&a.ID, &a.Name, &a.Type, &a.PhotoPath, &a.Style, &a.Status, &a.Result, &a.OutputPath, &a.CreatedAt, &a.UpdatedAt)
	return a, err
}

func (m *AvatarMapper) Insert(ctx context.Context, a *entity.Avatar) error {
	_, err := m.db.ExecContext(ctx,
		`INSERT INTO avatars (id, name, type, photo_path, style, status)
		 VALUES (?, ?, ?, ?, ?, ?)`,
		a.ID, a.Name, a.Type, a.PhotoPath, a.Style, a.Status)
	return err
}

func (m *AvatarMapper) GetByID(ctx context.Context, id string) (*entity.Avatar, error) {
	a, err := scanAvatar(m.db.QueryRowContext(ctx, `SELECT `+avatarCols+` FROM avatars WHERE id = ?`, id))
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (m *AvatarMapper) List(ctx context.Context, typ string) ([]entity.Avatar, error) {
	rows, err := m.db.QueryContext(ctx,
		`SELECT `+avatarCols+` FROM avatars
		 WHERE (? = '' OR type = ?)
		 ORDER BY created_at DESC`, typ, typ)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []entity.Avatar
	for rows.Next() {
		a, err := scanAvatar(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, a)
	}
	return out, rows.Err()
}

func (m *AvatarMapper) UpdateStatus(ctx context.Context, id, status string) error {
	_, err := m.db.ExecContext(ctx, `UPDATE avatars SET status = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`, status, id)
	return err
}

// UpdateResult stores the generation output (result JSON + quick-access path) and status.
func (m *AvatarMapper) UpdateResult(ctx context.Context, id, result, outputPath, status string) error {
	_, err := m.db.ExecContext(ctx,
		`UPDATE avatars SET result = ?, output_path = ?, status = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`,
		result, outputPath, status, id)
	return err
}

func (m *AvatarMapper) Delete(ctx context.Context, id string) error {
	_, err := m.db.ExecContext(ctx, `DELETE FROM avatars WHERE id = ?`, id)
	return err
}
