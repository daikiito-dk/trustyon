package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/daikiito-dk/trustyon/services/api/internal/httpx"
)

type GitHubStats struct {
	From         time.Time `json:"from"`
	To           time.Time `json:"to"`
	Commits      int       `json:"commits"`
	PullRequests int       `json:"pullRequests"`
	Issues       int       `json:"issues"`
	Repositories int       `json:"repositories"`
}

type githubSearchResponse struct {
	TotalCount int `json:"total_count"`
}

func NewGitHubStatsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("trustyon_github_token")
		if err != nil || strings.TrimSpace(cookie.Value) == "" {
			httpx.Error(w, http.StatusUnauthorized, "GitHub authentication required")
			return
		}

		to := time.Now().UTC()
		from := to.AddDate(0, 0, -7)
		login, err := githubLogin(r, cookie.Value)
		if err != nil {
			httpx.Error(w, http.StatusBadGateway, err.Error())
			return
		}

		commits, err := githubSearchCount(r, cookie.Value, fmt.Sprintf("author:%s committer-date:%s..%s", login, from.Format("2006-01-02"), to.Format("2006-01-02")), "commits")
		if err != nil {
			httpx.Error(w, http.StatusBadGateway, err.Error())
			return
		}
		prs, err := githubSearchCount(r, cookie.Value, fmt.Sprintf("author:%s created:%s..%s type:pr", login, from.Format("2006-01-02"), to.Format("2006-01-02")), "issues")
		if err != nil {
			httpx.Error(w, http.StatusBadGateway, err.Error())
			return
		}
		issues, err := githubSearchCount(r, cookie.Value, fmt.Sprintf("author:%s created:%s..%s type:issue", login, from.Format("2006-01-02"), to.Format("2006-01-02")), "issues")
		if err != nil {
			httpx.Error(w, http.StatusBadGateway, err.Error())
			return
		}
		repositories, err := githubRepoCount(r, cookie.Value)
		if err != nil {
			httpx.Error(w, http.StatusBadGateway, err.Error())
			return
		}

		httpx.Write(w, http.StatusOK, GitHubStats{From: from, To: to, Commits: commits, PullRequests: prs, Issues: issues, Repositories: repositories})
	}
}

func githubLogin(r *http.Request, token string) (string, error) {
	var user struct{ Login string `json:"login"` }
	if err := githubJSON(r, token, "https://api.github.com/user", &user); err != nil { return "", err }
	if user.Login == "" { return "", fmt.Errorf("GitHub user response did not include login") }
	return user.Login, nil
}

func githubSearchCount(r *http.Request, token, query, endpoint string) (int, error) {
	var result githubSearchResponse
	values := url.Values{"q": {query}}
	if err := githubJSON(r, token, "https://api.github.com/search/"+endpoint+"?"+values.Encode(), &result); err != nil { return 0, err }
	return result.TotalCount, nil
}

func githubRepoCount(r *http.Request, token string) (int, error) {
	req, err := http.NewRequestWithContext(r.Context(), http.MethodGet, "https://api.github.com/user/repos?per_page=100&page=1&affiliation=owner,collaborator,organization_member", nil)
	if err != nil { return 0, fmt.Errorf("could not create GitHub repository request") }
	setGitHubHeaders(req, token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil { return 0, fmt.Errorf("GitHub is unavailable") }
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK { return 0, fmt.Errorf("GitHub returned %d", resp.StatusCode) }
	var repos []json.RawMessage
	if err := json.NewDecoder(resp.Body).Decode(&repos); err != nil { return 0, fmt.Errorf("invalid GitHub repository response") }
	count := len(repos)
	if link := resp.Header.Get("Link"); link != "" && strings.Contains(link, `rel="last"`) {
		for _, part := range strings.Split(link, ",") {
			if !strings.Contains(part, `rel="last"`) { continue }
			start, end := strings.Index(part, "page="), strings.Index(part[start:], ">")
			if start >= 0 && end > 0 {
				end += start
				pageText := part[start+5 : end]
				if page, parseErr := strconv.Atoi(pageText); parseErr == nil { count = (page-1)*100 + len(repos) }
			}
		}
	}
	return count, nil
}

func githubJSON(r *http.Request, token, endpoint string, target any) error {
	req, err := http.NewRequestWithContext(r.Context(), http.MethodGet, endpoint, nil)
	if err != nil { return fmt.Errorf("could not create GitHub request") }
	setGitHubHeaders(req, token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil { return fmt.Errorf("GitHub is unavailable") }
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK { return fmt.Errorf("GitHub returned %d", resp.StatusCode) }
	if err := json.NewDecoder(resp.Body).Decode(target); err != nil { return fmt.Errorf("invalid GitHub response") }
	return nil
}

func setGitHubHeaders(req *http.Request, token string) {
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
}

var _ = os.Getenv
