package mapper

import (
	"context"

	"github.com/Bruce-Sakura/UploadMyself/backend/pkg/tasks/entity"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// TaskMapper is the data-access layer for the tasks table.
type TaskMapper struct {
	pool *pgxpool.Pool
}

func NewTaskMapper(pool *pgxpool.Pool) *TaskMapper {
	return &TaskMapper{pool: pool}
}

func (m *TaskMapper) Insert(ctx context.Context, t *entity.Task) error {
	_, err := m.pool.Exec(ctx,
		`INSERT INTO tasks (id, type, ref_id, status, progress, error)
		 VALUES ($1, $2, $3, $4, $5, $6)`,
		t.ID, t.Type, t.RefID, t.Status, t.Progress, t.Error)
	return err
}

func (m *TaskMapper) GetByID(ctx context.Context, id string) (*entity.Task, error) {
	row := m.pool.QueryRow(ctx,
		`SELECT id, type, ref_id, status, progress, error, created_at, updated_at
		 FROM tasks WHERE id = $1`, id)
	var t entity.Task
	if err := row.Scan(&t.ID, &t.Type, &t.RefID, &t.Status, &t.Progress, &t.Error, &t.CreatedAt, &t.UpdatedAt); err != nil {
		return nil, err
	}
	return &t, nil
}

// UpdateStatus updates status/progress and (optionally) the error message.
func (m *TaskMapper) UpdateStatus(ctx context.Context, id, status string, progress int, errMsg string) error {
	_, err := m.pool.Exec(ctx,
		`UPDATE tasks SET status = $2, progress = $3,
		 error = CASE WHEN $4 <> '' THEN $4 ELSE error END,
		 updated_at = now() WHERE id = $1`,
		id, status, progress, errMsg)
	return err
}

func (m *TaskMapper) List(ctx context.Context, typ, refID string) ([]entity.Task, error) {
	rows, err := m.pool.Query(ctx,
		`SELECT id, type, ref_id, status, progress, error, created_at, updated_at
		 FROM tasks
		 WHERE ($1 = '' OR type = $1) AND ($2 = '' OR ref_id = $2)
		 ORDER BY created_at DESC`, typ, refID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks, err := pgx.CollectRows(rows, func(r pgx.CollectableRow) (entity.Task, error) {
		var t entity.Task
		err := r.Scan(&t.ID, &t.Type, &t.RefID, &t.Status, &t.Progress, &t.Error, &t.CreatedAt, &t.UpdatedAt)
		return t, err
	})
	return tasks, err
}
