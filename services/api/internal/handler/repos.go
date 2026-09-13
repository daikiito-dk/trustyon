package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/daikiito-dk/trustyon/services/api/internal/httpx"
)

type GitHubRepo struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	FullName    string `json:"fullName"`
	Description string `json:"description"`
	URL         string `json:"url"`
	Private     bool   `json:"private"`
	Language    string `json:"language,omitempty"`
	Stars       int    `json:"stars"`
	UpdatedAt   string `json:"updatedAt"`
}

type githubRepoResponse struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	FullName    string `json:"full_name"`
	Description string `json:"description"`
	HTMLURL     string `json:"html_url"`
	Private     bool   `json:"private"`
	Language    string `json:"language"`
	Stars       int    `json:"stargazers_count"`
	UpdatedAt   string `json:"updated_at"`
}

func NewGitHubRepoHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := ""
		if cookie, err := r.Cookie("trustyon_github_token"); err == nil {
			token = cookie.Value
		}
		username := strings.TrimSpace(os.Getenv("GITHUB_USERNAME"))
		endpoint := "https://api.github.com/users/" + username + "/repos?per_page=100&sort=updated"
		if token != "" {
			endpoint = "https://api.github.com/user/repos?per_page=100&sort=updated&affiliation=owner,collaborator,organization_member"
		} else if username == "" {
			httpx.Error(w, http.StatusBadRequest, "GITHUB_USERNAME is not configured")
			return
		}

		req, err := http.NewRequestWithContext(r.Context(), http.MethodGet, endpoint, nil)
		if err != nil {
			httpx.Error(w, http.StatusInternalServerError, "could not create GitHub request")
			return
		}
		req.Header.Set("Accept", "application/vnd.github+json")
		req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
		if token != "" {
			req.Header.Set("Authorization", "Bearer "+token)
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			httpx.Error(w, http.StatusBadGateway, "GitHub is unavailable")
			return
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			httpx.Error(w, http.StatusBadGateway, fmt.Sprintf("GitHub returned %d", resp.StatusCode))
			return
		}

		var source []githubRepoResponse
		if err := json.NewDecoder(resp.Body).Decode(&source); err != nil {
			httpx.Error(w, http.StatusBadGateway, "invalid GitHub response")
			return
		}
		repos := make([]GitHubRepo, 0, len(source))
		for _, repo := range source {
			repos = append(repos, GitHubRepo{ID: repo.ID, Name: repo.Name, FullName: repo.FullName, Description: repo.Description, URL: repo.HTMLURL, Private: repo.Private, Language: repo.Language, Stars: repo.Stars, UpdatedAt: repo.UpdatedAt})
		}
		httpx.Write(w, http.StatusOK, repos)
	}
}
