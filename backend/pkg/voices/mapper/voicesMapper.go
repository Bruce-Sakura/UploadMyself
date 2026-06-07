package mapper

import (
	"context"

	"github.com/Bruce-Sakura/UploadMyself/backend/pkg/voices/entity"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// VoiceMapper is the data-access layer for the voices table.
type VoiceMapper struct {
	pool *pgxpool.Pool
}

func NewVoiceMapper(pool *pgxpool.Pool) *VoiceMapper {
	return &VoiceMapper{pool: pool}
}

const voiceCols = `id, name, audio_path, duration, model_path, ref_audio_path, status, created_at, updated_at`

func scanVoice(r pgx.Row) (entity.Voice, error) {
	var v entity.Voice
	err := r.Scan(&v.ID, &v.Name, &v.AudioPath, &v.Duration, &v.ModelPath, &v.RefAudioPath, &v.Status, &v.CreatedAt, &v.UpdatedAt)
	return v, err
}

func (m *VoiceMapper) Insert(ctx context.Context, v *entity.Voice) error {
	_, err := m.pool.Exec(ctx,
		`INSERT INTO voices (id, name, audio_path, duration, status) VALUES ($1, $2, $3, $4, $5)`,
		v.ID, v.Name, v.AudioPath, v.Duration, v.Status)
	return err
}

func (m *VoiceMapper) GetByID(ctx context.Context, id string) (*entity.Voice, error) {
	v, err := scanVoice(m.pool.QueryRow(ctx, `SELECT `+voiceCols+` FROM voices WHERE id = $1`, id))
	if err != nil {
		return nil, err
	}
	return &v, nil
}

func (m *VoiceMapper) List(ctx context.Context) ([]entity.Voice, error) {
	rows, err := m.pool.Query(ctx, `SELECT `+voiceCols+` FROM voices ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return pgx.CollectRows(rows, func(r pgx.CollectableRow) (entity.Voice, error) {
		return scanVoice(r)
	})
}

func (m *VoiceMapper) UpdateStatus(ctx context.Context, id, status string) error {
	_, err := m.pool.Exec(ctx, `UPDATE voices SET status = $2, updated_at = now() WHERE id = $1`, id, status)
	return err
}

func (m *VoiceMapper) UpdateModelPath(ctx context.Context, id, modelPath string) error {
	_, err := m.pool.Exec(ctx, `UPDATE voices SET model_path = $2, updated_at = now() WHERE id = $1`, id, modelPath)
	return err
}

func (m *VoiceMapper) Delete(ctx context.Context, id string) error {
	_, err := m.pool.Exec(ctx, `DELETE FROM voices WHERE id = $1`, id)
	return err
}
