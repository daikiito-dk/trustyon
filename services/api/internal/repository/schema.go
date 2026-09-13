package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

// EnsureSchema applies idempotent schema changes introduced after the initial
// MVP migration. This keeps existing development databases compatible with
// newer application binaries while SQL files remain the source of truth for
// fresh PostgreSQL volumes.
func EnsureSchema(ctx context.Context, db *pgxpool.Pool) error {
	statements := []string{
		`CREATE TABLE IF NOT EXISTS github_activities (
			id TEXT PRIMARY KEY,
			type TEXT NOT NULL,
			repo TEXT NOT NULL,
			summary TEXT NOT NULL,
			url TEXT NOT NULL DEFAULT '',
			created_at TIMESTAMPTZ NOT NULL,
			fetched_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,
		`CREATE INDEX IF NOT EXISTS idx_github_activities_created_at ON github_activities(created_at DESC)`,
		`ALTER TABLE projects ADD COLUMN IF NOT EXISTS github_repo_id BIGINT`,
		`ALTER TABLE projects ADD COLUMN IF NOT EXISTS github_full_name TEXT NOT NULL DEFAULT ''`,
		`ALTER TABLE projects ADD COLUMN IF NOT EXISTS github_synced_at TIMESTAMPTZ`,
		`CREATE UNIQUE INDEX IF NOT EXISTS idx_projects_github_repo_id ON projects(github_repo_id) WHERE github_repo_id IS NOT NULL`,
	}
	for _, statement := range statements {
		if _, err := db.Exec(ctx, statement); err != nil {
			return err
		}
	}
	return nil
}
