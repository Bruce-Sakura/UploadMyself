package mapper

import (
	"context"

	"github.com/Bruce-Sakura/UploadMyself/backend/pkg/avatars/entity"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// AvatarMapper is the data-access layer for the avatars table.
type AvatarMapper struct {
	pool *pgxpool.Pool
}

func NewAvatarMapper(pool *pgxpool.Pool) *AvatarMapper {
	return &AvatarMapper{pool: pool}
}

const avatarCols = `id, name, type, photo_path, style, status, result, output_path, created_at, updated_at`

func scanAvatar(r pgx.Row) (entity.Avatar, error) {
	var a entity.Avatar
	err := r.Scan(&a.ID, &a.Name, &a.Type, &a.PhotoPath, &a.Style, &a.Status, &a.Result, &a.OutputPath, &a.CreatedAt, &a.UpdatedAt)
	return a, err
}

func (m *AvatarMapper) Insert(ctx context.Context, a *entity.Avatar) error {
	_, err := m.pool.Exec(ctx,
		`INSERT INTO avatars (id, name, type, photo_path, style, status)
		 VALUES ($1, $2, $3, $4, $5, $6)`,
		a.ID, a.Name, a.Type, a.PhotoPath, a.Style, a.Status)
	return err
}

func (m *AvatarMapper) GetByID(ctx context.Context, id string) (*entity.Avatar, error) {
	a, err := scanAvatar(m.pool.QueryRow(ctx, `SELECT `+avatarCols+` FROM avatars WHERE id = $1`, id))
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (m *AvatarMapper) List(ctx context.Context, typ string) ([]entity.Avatar, error) {
	rows, err := m.pool.Query(ctx,
		`SELECT `+avatarCols+` FROM avatars
		 WHERE ($1 = '' OR type = $1)
		 ORDER BY created_at DESC`, typ)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return pgx.CollectRows(rows, func(r pgx.CollectableRow) (entity.Avatar, error) {
		return scanAvatar(r)
	})
}

func (m *AvatarMapper) UpdateStatus(ctx context.Context, id, status string) error {
	_, err := m.pool.Exec(ctx, `UPDATE avatars SET status = $2, updated_at = now() WHERE id = $1`, id, status)
	return err
}

// UpdateResult stores the generation output (result JSON + quick-access path) and status.
func (m *AvatarMapper) UpdateResult(ctx context.Context, id, result, outputPath, status string) error {
	_, err := m.pool.Exec(ctx,
		`UPDATE avatars SET result = $2, output_path = $3, status = $4, updated_at = now() WHERE id = $1`,
		id, result, outputPath, status)
	return err
}

func (m *AvatarMapper) Delete(ctx context.Context, id string) error {
	_, err := m.pool.Exec(ctx, `DELETE FROM avatars WHERE id = $1`, id)
	return err
}
