package repository

import (
	"context"

	"github.com/daikiito-dk/trustyon/services/api/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProjectRepository struct { db *pgxpool.Pool }

func NewProjectRepository(db *pgxpool.Pool) *ProjectRepository { return &ProjectRepository{db: db} }

const projectColumns = `id,name,slug,description,status,repo_url,github_repo_id,github_full_name,github_synced_at,created_at,updated_at`

func (r *ProjectRepository) List(ctx context.Context) ([]model.Project, error) {
	rows, err := r.db.Query(ctx, `SELECT `+projectColumns+` FROM projects ORDER BY created_at DESC`)
	if err != nil { return nil, err }
	defer rows.Close()
	items := make([]model.Project, 0)
	for rows.Next() {
		var p model.Project
		if err := rows.Scan(&p.ID,&p.Name,&p.Slug,&p.Description,&p.Status,&p.RepoURL,&p.GitHubRepoID,&p.GitHubFullName,&p.GitHubSyncedAt,&p.CreatedAt,&p.UpdatedAt); err != nil { return nil, err }
		items = append(items, p)
	}
	return items, rows.Err()
}

func (r *ProjectRepository) Create(ctx context.Context, p *model.Project) error {
	return r.db.QueryRow(ctx, `INSERT INTO projects(name,slug,description,status,repo_url,github_repo_id,github_full_name,github_synced_at) VALUES($1,$2,$3,$4,$5,$6,$7,$8) RETURNING `+projectColumns, p.Name,p.Slug,p.Description,p.Status,p.RepoURL,p.GitHubRepoID,p.GitHubFullName,p.GitHubSyncedAt).Scan(&p.ID,&p.Name,&p.Slug,&p.Description,&p.Status,&p.RepoURL,&p.GitHubRepoID,&p.GitHubFullName,&p.GitHubSyncedAt,&p.CreatedAt,&p.UpdatedAt)
}

func (r *ProjectRepository) Update(ctx context.Context, id int64, p *model.Project) error {
	return r.db.QueryRow(ctx, `UPDATE projects SET name=$1,slug=$2,description=$3,status=$4,repo_url=$5,github_repo_id=$6,github_full_name=$7,github_synced_at=$8,updated_at=NOW() WHERE id=$9 RETURNING `+projectColumns, p.Name,p.Slug,p.Description,p.Status,p.RepoURL,p.GitHubRepoID,p.GitHubFullName,p.GitHubSyncedAt,id).Scan(&p.ID,&p.Name,&p.Slug,&p.Description,&p.Status,&p.RepoURL,&p.GitHubRepoID,&p.GitHubFullName,&p.GitHubSyncedAt,&p.CreatedAt,&p.UpdatedAt)
}

func (r *ProjectRepository) Delete(ctx context.Context, id int64) error {
	_, err := r.db.Exec(ctx, `DELETE FROM projects WHERE id=$1`, id)
	return err
}

func (r *ProjectRepository) FindByGitHubRepoID(ctx context.Context, repoID int64) (*model.Project, error) {
	var p model.Project
	err := r.db.QueryRow(ctx, `SELECT `+projectColumns+` FROM projects WHERE github_repo_id=$1`, repoID).Scan(&p.ID,&p.Name,&p.Slug,&p.Description,&p.Status,&p.RepoURL,&p.GitHubRepoID,&p.GitHubFullName,&p.GitHubSyncedAt,&p.CreatedAt,&p.UpdatedAt)
	if err != nil { return nil, err }
	return &p, nil
}

func (r *ProjectRepository) SyncGitHub(ctx context.Context, repoID int64, fullName, name, description, repoURL string, now interface{}) (*model.Project, error) {
	var p model.Project
	err := r.db.QueryRow(ctx, `
		INSERT INTO projects(name,slug,description,status,repo_url,github_repo_id,github_full_name,github_synced_at)
		VALUES($1,$2,$3,'active',$4,$5,$6,$7)
		ON CONFLICT (github_repo_id) DO UPDATE SET name=EXCLUDED.name,description=EXCLUDED.description,repo_url=EXCLUDED.repo_url,github_full_name=EXCLUDED.github_full_name,github_synced_at=EXCLUDED.github_synced_at,updated_at=NOW()
		RETURNING `+projectColumns, name, slugify(name), description, repoURL, repoID, fullName, now).Scan(&p.ID,&p.Name,&p.Slug,&p.Description,&p.Status,&p.RepoURL,&p.GitHubRepoID,&p.GitHubFullName,&p.GitHubSyncedAt,&p.CreatedAt,&p.UpdatedAt)
	if err != nil { return nil, err }
	return &p, nil
}

func slugify(value string) string {
	out := make([]byte, 0, len(value))
	for _, c := range value {
		if c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c >= '0' && c <= '9' { if c >= 'A' && c <= 'Z' { c += 'a' - 'A' }; out = append(out, byte(c)) } else if len(out) > 0 && out[len(out)-1] != '-' { out = append(out, '-') }
	}
	for len(out) > 0 && out[len(out)-1] == '-' { out = out[:len(out)-1] }
	if len(out) == 0 { return "github-project" }
	return string(out)
}
