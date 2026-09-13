CREATE TABLE IF NOT EXISTS github_activities (
    id TEXT PRIMARY KEY,
    type TEXT NOT NULL,
    repo TEXT NOT NULL,
    summary TEXT NOT NULL,
    url TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL,
    fetched_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_github_activities_created_at ON github_activities(created_at DESC);
