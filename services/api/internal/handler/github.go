package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/daikiito-dk/trustyon/services/api/internal/httpx"
	"github.com/daikiito-dk/trustyon/services/api/internal/repository"
)

type GitHubActivity struct {
	ID        string    `json:"id"`
	Type      string    `json:"type"`
	Repo      string    `json:"repo"`
	Summary   string    `json:"summary"`
	URL       string    `json:"url,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
}

type githubEvent struct {
	ID     string `json:"id"`
	Type   string `json:"type"`
	Repo   struct { Name string `json:"name"` } `json:"repo"`
	Created string `json:"created_at"`
	Payload struct {
		Action      string `json:"action"`
		Commits     []struct { Message string `json:"message"` } `json:"commits"`
		Issue       struct { Title string `json:"title"` } `json:"issue"`
		PullRequest struct { Title string `json:"title"` } `json:"pull_request"`
	} `json:"payload"`
}

func NewGitHubActivityHandler(repo *repository.GitHubActivityRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username := strings.TrimSpace(os.Getenv("GITHUB_USERNAME"))
		if username == "" { username = "daikiito-dk" }

		req, err := http.NewRequestWithContext(r.Context(), http.MethodGet, "https://api.github.com/users/"+username+"/events/public?per_page=100", nil)
		if err != nil { httpx.Error(w, http.StatusInternalServerError, "could not create GitHub request"); return }
		req.Header.Set("Accept", "application/vnd.github+json")
		req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
		resp, err := http.DefaultClient.Do(req)
		if err != nil { httpx.Error(w, http.StatusBadGateway, "GitHub is unavailable"); return }
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK { httpx.Error(w, http.StatusBadGateway, fmt.Sprintf("GitHub returned %d", resp.StatusCode)); return }

		var events []githubEvent
		if err := json.NewDecoder(resp.Body).Decode(&events); err != nil { httpx.Error(w, http.StatusBadGateway, "invalid GitHub response"); return }

		for _, event := range events {
			created, err := time.Parse(time.RFC3339, event.Created)
			if err != nil { continue }
			summary := event.Type
			switch event.Type {
			case "PushEvent":
				count := len(event.Payload.Commits)
				summary = fmt.Sprintf("pushed %d commit%s", count, plural(count))
				if count > 0 && strings.TrimSpace(event.Payload.Commits[0].Message) != "" { summary += ": " + strings.Split(event.Payload.Commits[0].Message, "\n")[0] }
			case "IssuesEvent":
				summary = fmt.Sprintf("%s issue: %s", event.Payload.Action, event.Payload.Issue.Title)
			case "PullRequestEvent":
				summary = fmt.Sprintf("%s PR: %s", event.Payload.Action, event.Payload.PullRequest.Title)
			case "CreateEvent":
				summary = "created a repository reference"
			case "DeleteEvent":
				summary = "deleted a repository reference"
			case "WatchEvent":
				summary = "starred a repository"
			case "ForkEvent":
				summary = "forked a repository"
			}
			if err := repo.Upsert(r.Context(), event.ID, event.Type, event.Repo.Name, summary, "", created); err != nil {
				httpx.Error(w, http.StatusInternalServerError, "could not save GitHub activity")
				return
			}
		}

		rows, err := repo.List(r.Context(), 100)
		if err != nil { httpx.Error(w, http.StatusInternalServerError, "could not load GitHub activity"); return }
		activities := make([]GitHubActivity, 0, len(rows))
		for _, row := range rows {
			activities = append(activities, GitHubActivity{ID: row.ID, Type: row.Type, Repo: row.Repo, Summary: row.Summary, URL: row.URL, CreatedAt: row.CreatedAt})
		}
		httpx.Write(w, http.StatusOK, activities)
	}
}

func plural(n int) string { if n == 1 { return "" }; return "s" }
