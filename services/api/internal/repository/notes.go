package repository

import (
	"context"

	"github.com/daikiito-dk/trustyon/services/api/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type NoteRepository struct { db *pgxpool.Pool }

func NewNoteRepository(db *pgxpool.Pool) *NoteRepository { return &NoteRepository{db: db} }

func (r *NoteRepository) List(ctx context.Context) ([]model.Note, error) {
	rows, err := r.db.Query(ctx, `SELECT id,project_id,title,body,created_at,updated_at FROM notes ORDER BY updated_at DESC`)
	if err != nil { return nil, err }
	defer rows.Close()
	items := make([]model.Note, 0)
	for rows.Next() {
		var n model.Note
		if err := rows.Scan(&n.ID,&n.ProjectID,&n.Title,&n.Body,&n.CreatedAt,&n.UpdatedAt); err != nil { return nil, err }
		items = append(items, n)
	}
	return items, rows.Err()
}

func (r *NoteRepository) Create(ctx context.Context, n *model.Note) error {
	return r.db.QueryRow(ctx, `INSERT INTO notes(project_id,title,body) VALUES($1,$2,$3) RETURNING id,created_at,updated_at`, n.ProjectID,n.Title,n.Body).Scan(&n.ID,&n.CreatedAt,&n.UpdatedAt)
}

func (r *NoteRepository) Update(ctx context.Context, id int64, n *model.Note) error {
	return r.db.QueryRow(ctx, `UPDATE notes SET project_id=$1,title=$2,body=$3,updated_at=NOW() WHERE id=$4 RETURNING id,created_at,updated_at`, n.ProjectID,n.Title,n.Body,id).Scan(&n.ID,&n.CreatedAt,&n.UpdatedAt)
}

func (r *NoteRepository) Delete(ctx context.Context, id int64) error {
	_, err := r.db.Exec(ctx, `DELETE FROM notes WHERE id=$1`, id)
	return err
}
