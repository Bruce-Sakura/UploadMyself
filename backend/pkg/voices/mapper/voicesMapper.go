package mapper

import (
	"context"
	"database/sql"

	"github.com/Bruce-Sakura/UploadMyself/backend/pkg/voices/entity"
)

// VoiceMapper is the data-access layer for the voices table.
type VoiceMapper struct {
	db *sql.DB
}

func NewVoiceMapper(db *sql.DB) *VoiceMapper {
	return &VoiceMapper{db: db}
}

const voiceCols = `id, name, audio_path, duration, model_path, ref_audio_path, status, created_at, updated_at`

func scanVoice(s interface {
	Scan(dest ...any) error
}) (entity.Voice, error) {
	var v entity.Voice
	err := s.Scan(&v.ID, &v.Name, &v.AudioPath, &v.Duration, &v.ModelPath, &v.RefAudioPath, &v.Status, &v.CreatedAt, &v.UpdatedAt)
	return v, err
}

func (m *VoiceMapper) Insert(ctx context.Context, v *entity.Voice) error {
	_, err := m.db.ExecContext(ctx,
		`INSERT INTO voices (id, name, audio_path, duration, status) VALUES (?, ?, ?, ?, ?)`,
		v.ID, v.Name, v.AudioPath, v.Duration, v.Status)
	return err
}

func (m *VoiceMapper) GetByID(ctx context.Context, id string) (*entity.Voice, error) {
	v, err := scanVoice(m.db.QueryRowContext(ctx, `SELECT `+voiceCols+` FROM voices WHERE id = ?`, id))
	if err != nil {
		return nil, err
	}
	return &v, nil
}

func (m *VoiceMapper) List(ctx context.Context) ([]entity.Voice, error) {
	rows, err := m.db.QueryContext(ctx, `SELECT `+voiceCols+` FROM voices ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []entity.Voice
	for rows.Next() {
		v, err := scanVoice(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, v)
	}
	return out, rows.Err()
}

func (m *VoiceMapper) UpdateStatus(ctx context.Context, id, status string) error {
	_, err := m.db.ExecContext(ctx, `UPDATE voices SET status = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`, status, id)
	return err
}

func (m *VoiceMapper) UpdateModelPath(ctx context.Context, id, modelPath string) error {
	_, err := m.db.ExecContext(ctx, `UPDATE voices SET model_path = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`, modelPath, id)
	return err
}

func (m *VoiceMapper) Delete(ctx context.Context, id string) error {
	_, err := m.db.ExecContext(ctx, `DELETE FROM voices WHERE id = ?`, id)
	return err
}
