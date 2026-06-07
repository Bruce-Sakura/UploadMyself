package mapper

import (
	"context"
	"database/sql"

	"github.com/Bruce-Sakura/UploadMyself/backend/pkg/skills/entity"
)

// SkillMapper is the data-access layer for the skills table.
type SkillMapper struct {
	db *sql.DB
}

func NewSkillMapper(db *sql.DB) *SkillMapper {
	return &SkillMapper{db: db}
}

const skillCols = `id, name, corpus, status, result, created_at, updated_at`

func scanSkill(s interface {
	Scan(dest ...any) error
}) (entity.Skill, error) {
	var sk entity.Skill
	err := s.Scan(&sk.ID, &sk.Name, &sk.Corpus, &sk.Status, &sk.Result, &sk.CreatedAt, &sk.UpdatedAt)
	return sk, err
}

func (m *SkillMapper) Insert(ctx context.Context, s *entity.Skill) error {
	_, err := m.db.ExecContext(ctx,
		`INSERT INTO skills (id, name, corpus, status) VALUES (?, ?, ?, ?)`,
		s.ID, s.Name, s.Corpus, s.Status)
	return err
}

func (m *SkillMapper) GetByID(ctx context.Context, id string) (*entity.Skill, error) {
	sk, err := scanSkill(m.db.QueryRowContext(ctx, `SELECT `+skillCols+` FROM skills WHERE id = ?`, id))
	if err != nil {
		return nil, err
	}
	return &sk, nil
}

func (m *SkillMapper) List(ctx context.Context) ([]entity.Skill, error) {
	rows, err := m.db.QueryContext(ctx, `SELECT `+skillCols+` FROM skills ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []entity.Skill
	for rows.Next() {
		sk, err := scanSkill(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, sk)
	}
	return out, rows.Err()
}

// Update applies non-nil fields (COALESCE keeps current value for NULL args).
// database/sql 不支持 *string 参数，需先解引用为 any（nil → SQL NULL）。
func (m *SkillMapper) Update(ctx context.Context, id string, name, corpus, status, result *string) error {
	_, err := m.db.ExecContext(ctx,
		`UPDATE skills SET
		   name = COALESCE(?, name),
		   corpus = COALESCE(?, corpus),
		   status = COALESCE(?, status),
		   result = COALESCE(?, result),
		   updated_at = CURRENT_TIMESTAMP
		 WHERE id = ?`,
		deref(name), deref(corpus), deref(status), deref(result), id)
	return err
}

func deref(p *string) any {
	if p == nil {
		return nil
	}
	return *p
}

func (m *SkillMapper) Delete(ctx context.Context, id string) error {
	_, err := m.db.ExecContext(ctx, `DELETE FROM skills WHERE id = ?`, id)
	return err
}
