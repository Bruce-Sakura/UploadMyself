package mapper

import (
	"context"

	"github.com/Bruce-Sakura/UploadMyself/backend/pkg/skills/entity"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// SkillMapper is the data-access layer for the skills table.
type SkillMapper struct {
	pool *pgxpool.Pool
}

func NewSkillMapper(pool *pgxpool.Pool) *SkillMapper {
	return &SkillMapper{pool: pool}
}

const skillCols = `id, name, corpus, status, result, created_at, updated_at`

func scanSkill(r pgx.Row) (entity.Skill, error) {
	var s entity.Skill
	err := r.Scan(&s.ID, &s.Name, &s.Corpus, &s.Status, &s.Result, &s.CreatedAt, &s.UpdatedAt)
	return s, err
}

func (m *SkillMapper) Insert(ctx context.Context, s *entity.Skill) error {
	_, err := m.pool.Exec(ctx,
		`INSERT INTO skills (id, name, corpus, status) VALUES ($1, $2, $3, $4)`,
		s.ID, s.Name, s.Corpus, s.Status)
	return err
}

func (m *SkillMapper) GetByID(ctx context.Context, id string) (*entity.Skill, error) {
	s, err := scanSkill(m.pool.QueryRow(ctx, `SELECT `+skillCols+` FROM skills WHERE id = $1`, id))
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (m *SkillMapper) List(ctx context.Context) ([]entity.Skill, error) {
	rows, err := m.pool.Query(ctx, `SELECT `+skillCols+` FROM skills ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return pgx.CollectRows(rows, func(r pgx.CollectableRow) (entity.Skill, error) {
		return scanSkill(r)
	})
}

// Update applies non-nil fields (COALESCE keeps current value for NULL args).
func (m *SkillMapper) Update(ctx context.Context, id string, name, corpus, status, result *string) error {
	_, err := m.pool.Exec(ctx,
		`UPDATE skills SET
		   name = COALESCE($2, name),
		   corpus = COALESCE($3, corpus),
		   status = COALESCE($4, status),
		   result = COALESCE($5, result),
		   updated_at = now()
		 WHERE id = $1`,
		id, name, corpus, status, result)
	return err
}

func (m *SkillMapper) Delete(ctx context.Context, id string) error {
	_, err := m.pool.Exec(ctx, `DELETE FROM skills WHERE id = $1`, id)
	return err
}
