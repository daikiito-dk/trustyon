package repository

import (
	"context"

	"github.com/daikiito-dk/trustyon/services/api/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProjectRepository struct { db *pgxpool.Pool }

func NewProjectRepository(db *pgxpool.Pool) *ProjectRepository { return &ProjectRepository{db: db} }

func (r *ProjectRepository) List(ctx context.Context) ([]model.Project, error) {
	rows, err := r.db.Query(ctx, `SELECT id,name,slug,description,status,repo_url,created_at,updated_at FROM projects ORDER BY created_at DESC`)
	if err != nil { return nil, err }
	defer rows.Close()
	items := make([]model.Project, 0)
	for rows.Next() {
		var p model.Project
		if err := rows.Scan(&p.ID,&p.Name,&p.Slug,&p.Description,&p.Status,&p.RepoURL,&p.CreatedAt,&p.UpdatedAt); err != nil { return nil, err }
		items = append(items, p)
	}
	return items, rows.Err()
}

func (r *ProjectRepository) Create(ctx context.Context, p *model.Project) error {
	return r.db.QueryRow(ctx, `INSERT INTO projects(name,slug,description,status,repo_url) VALUES($1,$2,$3,$4,$5) RETURNING id,created_at,updated_at`, p.Name,p.Slug,p.Description,p.Status,p.RepoURL).Scan(&p.ID,&p.CreatedAt,&p.UpdatedAt)
}
