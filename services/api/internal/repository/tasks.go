package repository

import (
	"context"

	"github.com/daikiito-dk/trustyon/services/api/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TaskRepository struct { db *pgxpool.Pool }

func NewTaskRepository(db *pgxpool.Pool) *TaskRepository { return &TaskRepository{db: db} }

func (r *TaskRepository) List(ctx context.Context) ([]model.Task, error) {
	rows, err := r.db.Query(ctx, `SELECT id,project_id,title,description,status,priority,due_date,created_at,updated_at FROM tasks ORDER BY created_at DESC`)
	if err != nil { return nil, err }
	defer rows.Close()
	items := make([]model.Task, 0)
	for rows.Next() {
		var t model.Task
		if err := rows.Scan(&t.ID,&t.ProjectID,&t.Title,&t.Description,&t.Status,&t.Priority,&t.DueDate,&t.CreatedAt,&t.UpdatedAt); err != nil { return nil, err }
		items = append(items, t)
	}
	return items, rows.Err()
}

func (r *TaskRepository) Create(ctx context.Context, t *model.Task) error {
	return r.db.QueryRow(ctx, `INSERT INTO tasks(project_id,title,description,status,priority,due_date) VALUES($1,$2,$3,$4,$5,$6) RETURNING id,created_at,updated_at`, t.ProjectID,t.Title,t.Description,t.Status,t.Priority,t.DueDate).Scan(&t.ID,&t.CreatedAt,&t.UpdatedAt)
}
