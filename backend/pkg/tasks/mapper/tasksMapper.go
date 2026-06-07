package mapper

import (
	"context"
	"database/sql"

	"github.com/Bruce-Sakura/UploadMyself/backend/pkg/tasks/entity"
)

// TaskMapper is the data-access layer for the tasks table.
type TaskMapper struct {
	db *sql.DB
}

func NewTaskMapper(db *sql.DB) *TaskMapper {
	return &TaskMapper{db: db}
}

func (m *TaskMapper) Insert(ctx context.Context, t *entity.Task) error {
	_, err := m.db.ExecContext(ctx,
		`INSERT INTO tasks (id, type, ref_id, status, progress, error)
		 VALUES (?, ?, ?, ?, ?, ?)`,
		t.ID, t.Type, t.RefID, t.Status, t.Progress, t.Error)
	return err
}

func (m *TaskMapper) GetByID(ctx context.Context, id string) (*entity.Task, error) {
	row := m.db.QueryRowContext(ctx,
		`SELECT id, type, ref_id, status, progress, error, created_at, updated_at
		 FROM tasks WHERE id = ?`, id)
	var t entity.Task
	if err := row.Scan(&t.ID, &t.Type, &t.RefID, &t.Status, &t.Progress, &t.Error, &t.CreatedAt, &t.UpdatedAt); err != nil {
		return nil, err
	}
	return &t, nil
}

// UpdateStatus updates status/progress and (optionally) the error message.
func (m *TaskMapper) UpdateStatus(ctx context.Context, id, status string, progress int, errMsg string) error {
	_, err := m.db.ExecContext(ctx,
		`UPDATE tasks SET status = ?, progress = ?,
		 error = CASE WHEN ? <> '' THEN ? ELSE error END,
		 updated_at = CURRENT_TIMESTAMP WHERE id = ?`,
		status, progress, errMsg, errMsg, id)
	return err
}

func (m *TaskMapper) List(ctx context.Context, typ, refID string) ([]entity.Task, error) {
	rows, err := m.db.QueryContext(ctx,
		`SELECT id, type, ref_id, status, progress, error, created_at, updated_at
		 FROM tasks
		 WHERE (? = '' OR type = ?) AND (? = '' OR ref_id = ?)
		 ORDER BY created_at DESC`, typ, typ, refID, refID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []entity.Task
	for rows.Next() {
		var t entity.Task
		if err := rows.Scan(&t.ID, &t.Type, &t.RefID, &t.Status, &t.Progress, &t.Error, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	return tasks, rows.Err()
}
