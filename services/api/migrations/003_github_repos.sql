ALTER TABLE projects ADD COLUMN IF NOT EXISTS github_repo_id BIGINT;
ALTER TABLE projects ADD COLUMN IF NOT EXISTS github_full_name TEXT NOT NULL DEFAULT '';
ALTER TABLE projects ADD COLUMN IF NOT EXISTS github_synced_at TIMESTAMPTZ;

CREATE UNIQUE INDEX IF NOT EXISTS idx_projects_github_repo_id
    ON projects(github_repo_id)
    WHERE github_repo_id IS NOT NULL;
