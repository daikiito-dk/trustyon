package model

import "time"

type Project struct {
	ID             int64      `json:"id"`
	Name           string     `json:"name"`
	Slug           string     `json:"slug"`
	Description    string     `json:"description"`
	Status         string     `json:"status"`
	RepoURL        string     `json:"repoUrl"`
	GitHubRepoID   *int64     `json:"githubRepoId,omitempty"`
	GitHubFullName string     `json:"githubFullName,omitempty"`
	GitHubSyncedAt *time.Time `json:"githubSyncedAt,omitempty"`
	CreatedAt      time.Time  `json:"createdAt"`
	UpdatedAt      time.Time  `json:"updatedAt"`
}

type Task struct {
	ID          int64      `json:"id"`
	ProjectID   *int64     `json:"projectId,omitempty"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      string     `json:"status"`
	Priority    string     `json:"priority"`
	DueDate     *time.Time `json:"dueDate,omitempty"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
}

type Note struct {
	ID        int64     `json:"id"`
	ProjectID *int64    `json:"projectId,omitempty"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
