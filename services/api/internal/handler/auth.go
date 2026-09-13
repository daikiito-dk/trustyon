package handler

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type GitHubAuthHandler struct{}

func NewGitHubAuthHandler() *GitHubAuthHandler { return &GitHubAuthHandler{} }

func (h *GitHubAuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	clientID := strings.TrimSpace(os.Getenv("GITHUB_CLIENT_ID"))
	if clientID == "" {
		http.Error(w, "GITHUB_CLIENT_ID is not configured", http.StatusServiceUnavailable)
		return
	}
	stateBytes := make([]byte, 24)
	if _, err := rand.Read(stateBytes); err != nil { http.Error(w, "could not create OAuth state", http.StatusInternalServerError); return }
	state := hex.EncodeToString(stateBytes)
	http.SetCookie(w, &http.Cookie{Name: "trustyon_oauth_state", Value: state, Path: "/", HttpOnly: true, SameSite: http.SameSiteLaxMode, Secure: os.Getenv("TRUSTYON_COOKIE_SECURE") == "true", MaxAge: 600})
	callback := envOr("GITHUB_REDIRECT_URI", "http://localhost:8080/api/auth/github/callback")
	params := url.Values{"client_id": {clientID}, "redirect_uri": {callback}, "scope": {"read:user repo"}, "state": {state}}
	http.Redirect(w, r, "https://github.com/login/oauth/authorize?"+params.Encode(), http.StatusFound)
}

func (h *GitHubAuthHandler) Callback(w http.ResponseWriter, r *http.Request) {
	stateCookie, err := r.Cookie("trustyon_oauth_state")
	if err != nil || stateCookie.Value == "" || stateCookie.Value != r.URL.Query().Get("state") {
		http.Error(w, "invalid OAuth state", http.StatusBadRequest)
		return
	}
	clientID := strings.TrimSpace(os.Getenv("GITHUB_CLIENT_ID"))
	clientSecret := strings.TrimSpace(os.Getenv("GITHUB_CLIENT_SECRET"))
	if clientID == "" || clientSecret == "" { http.Error(w, "GitHub OAuth is not configured", http.StatusServiceUnavailable); return }
	code := strings.TrimSpace(r.URL.Query().Get("code"))
	if code == "" { http.Error(w, "missing OAuth code", http.StatusBadRequest); return }

	form := url.Values{"client_id": {clientID}, "client_secret": {clientSecret}, "code": {code}, "redirect_uri": {envOr("GITHUB_REDIRECT_URI", "http://localhost:8080/api/auth/github/callback")}}
	req, err := http.NewRequestWithContext(r.Context(), http.MethodPost, "https://github.com/login/oauth/access_token", strings.NewReader(form.Encode()))
	if err != nil { http.Error(w, "could not create token request", http.StatusInternalServerError); return }
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := http.DefaultClient.Do(req)
	if err != nil { http.Error(w, "GitHub token exchange failed", http.StatusBadGateway); return }
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK { http.Error(w, fmt.Sprintf("GitHub token exchange returned %d", resp.StatusCode), http.StatusBadGateway); return }
	var token struct { AccessToken string `json:"access_token"`; Error string `json:"error"` }
	if err := json.NewDecoder(resp.Body).Decode(&token); err != nil { http.Error(w, "invalid GitHub token response", http.StatusBadGateway); return }
	if token.AccessToken == "" { http.Error(w, "GitHub did not return an access token", http.StatusBadGateway); return }

	cookie := &http.Cookie{Name: "trustyon_github_token", Value: token.AccessToken, Path: "/", HttpOnly: true, SameSite: http.SameSiteLaxMode, Secure: os.Getenv("TRUSTYON_COOKIE_SECURE") == "true", MaxAge: 30 * 24 * 60 * 60}
	http.SetCookie(w, cookie)
	http.SetCookie(w, &http.Cookie{Name: "trustyon_oauth_state", Value: "", Path: "/", HttpOnly: true, MaxAge: -1})
	http.Redirect(w, r, envOr("TRUSTYON_WEB_URL", "http://localhost:5173"), http.StatusFound)
}

func (h *GitHubAuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("trustyon_github_token")
	if err != nil || cookie.Value == "" { writeAuthJSON(w, http.StatusUnauthorized, map[string]string{"error": "not authenticated"}); return }
	req, err := http.NewRequestWithContext(r.Context(), http.MethodGet, "https://api.github.com/user", nil)
	if err != nil { writeAuthJSON(w, http.StatusInternalServerError, map[string]string{"error": "could not create GitHub request"}); return }
	req.Header.Set("Authorization", "Bearer "+cookie.Value)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	resp, err := http.DefaultClient.Do(req)
	if err != nil { writeAuthJSON(w, http.StatusBadGateway, map[string]string{"error": "GitHub is unavailable"}); return }
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK { writeAuthJSON(w, http.StatusUnauthorized, map[string]string{"error": "GitHub session is invalid"}); return }
	body, err := io.ReadAll(resp.Body)
	if err != nil { writeAuthJSON(w, http.StatusBadGateway, map[string]string{"error": "could not read GitHub response"}); return }
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(body)
}

func (h *GitHubAuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{Name: "trustyon_github_token", Value: "", Path: "/", HttpOnly: true, SameSite: http.SameSiteLaxMode, Secure: os.Getenv("TRUSTYON_COOKIE_SECURE") == "true", MaxAge: -1, Expires: time.Unix(1, 0)})
	w.WriteHeader(http.StatusNoContent)
}

func writeAuthJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func envOr(name, fallback string) string { if value := strings.TrimSpace(os.Getenv(name)); value != "" { return value }; return fallback }
