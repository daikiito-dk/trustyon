package handler

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/daikiito-dk/trustyon/services/api/internal/httpx"
	"github.com/daikiito-dk/trustyon/services/api/internal/repository"
)

type GitHubRepo struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	FullName    string `json:"fullName"`
	Description string `json:"description"`
	Private     bool   `json:"private"`
	HTMLURL     string `json:"htmlUrl"`
}

func NewGitHubSyncHandler(projects *repository.ProjectRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("trustyon_github_token")
		if err != nil || strings.TrimSpace(cookie.Value) == "" { httpx.Error(w, http.StatusUnauthorized, "GitHub authentication required"); return }
		req, err := http.NewRequestWithContext(r.Context(), http.MethodGet, "https://api.github.com/user/repos?per_page=100&sort=updated", nil)
		if err != nil { httpx.Error(w, http.StatusInternalServerError, "could not create GitHub request"); return }
		req.Header.Set("Authorization", "Bearer "+cookie.Value); req.Header.Set("Accept", "application/vnd.github+json"); req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
		resp, err := http.DefaultClient.Do(req); if err != nil { httpx.Error(w, http.StatusBadGateway, "GitHub is unavailable"); return }; defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK { httpx.Error(w, http.StatusBadGateway, "GitHub repository request failed"); return }
		var repos []GitHubRepo
		if err := json.NewDecoder(resp.Body).Decode(&repos); err != nil { httpx.Error(w, http.StatusBadGateway, "invalid GitHub repository response"); return }
		synced := make([]any, 0, len(repos)); now := time.Now().UTC()
		for _, repo := range repos {
			project, err := projects.SyncGitHub(r.Context(), repo.ID, repo.FullName, repo.Name, repo.Description, repo.HTMLURL, now)
			if err != nil { httpx.Error(w, http.StatusInternalServerError, "could not sync GitHub projects"); return }
			synced = append(synced, project)
		}
		httpx.Write(w, http.StatusOK, synced)
	}
}
