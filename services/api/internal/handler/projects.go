package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/daikiito-dk/trustyon/services/api/internal/model"
	"github.com/daikiito-dk/trustyon/services/api/internal/repository"
)

type ProjectHandler struct { repo *repository.ProjectRepository }
func NewProjectHandler(repo *repository.ProjectRepository) *ProjectHandler { return &ProjectHandler{repo: repo} }

func (h *ProjectHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method == http.MethodGet {
		items, err := h.repo.List(r.Context()); if err != nil { http.Error(w, `{"error":"database error"}`, 500); return }
		json.NewEncoder(w).Encode(items); return
	}
	if r.Method == http.MethodPost {
		var p model.Project
		if json.NewDecoder(r.Body).Decode(&p) != nil || strings.TrimSpace(p.Name) == "" || strings.TrimSpace(p.Slug) == "" {
			http.Error(w, `{"error":"name and slug are required"}`, 400); return
		}
		if p.Status == "" { p.Status = "active" }
		if err := h.repo.Create(r.Context(), &p); err != nil { http.Error(w, `{"error":"could not create project"}`, 500); return }
		w.WriteHeader(http.StatusCreated); json.NewEncoder(w).Encode(p); return
	}
	w.Header().Set("Allow", "GET, POST"); http.Error(w, `{"error":"method not allowed"}`, 405)
}
