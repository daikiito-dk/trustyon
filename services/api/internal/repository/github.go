package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type GitHubActivityRepository struct { db *pgxpool.Pool }

func NewGitHubActivityRepository(db *pgxpool.Pool) *GitHubActivityRepository {
	return &GitHubActivityRepository{db: db}
}

func (r *GitHubActivityRepository) Upsert(ctx context.Context, id, eventType, repo, summary, url string, createdAt time.Time) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO github_activities(id,type,repo,summary,url,created_at)
		VALUES($1,$2,$3,$4,$5,$6)
		ON CONFLICT (id) DO UPDATE SET type=EXCLUDED.type,repo=EXCLUDED.repo,summary=EXCLUDED.summary,url=EXCLUDED.url,created_at=EXCLUDED.created_at,fetched_at=NOW()
	`, id, eventType, repo, summary, url, createdAt)
	return err
}

type GitHubActivityRow struct {
	ID        string
	Type      string
	Repo      string
	Summary   string
	URL       string
	CreatedAt time.Time
}

func (r *GitHubActivityRepository) List(ctx context.Context, limit int) ([]GitHubActivityRow, error) {
	rows, err := r.db.Query(ctx, `SELECT id,type,repo,summary,url,created_at FROM github_activities ORDER BY created_at DESC LIMIT $1`, limit)
	if err != nil { return nil, err }
	defer rows.Close()
	items := make([]GitHubActivityRow, 0)
	for rows.Next() {
		var item GitHubActivityRow
		if err := rows.Scan(&item.ID,&item.Type,&item.Repo,&item.Summary,&item.URL,&item.CreatedAt); err != nil { return nil, err }
		items = append(items, item)
	}
	return items, rows.Err()
}
