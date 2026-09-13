package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/daikiito-dk/trustyon/services/api/internal/httpx"
	"github.com/daikiito-dk/trustyon/services/api/internal/model"
	"github.com/daikiito-dk/trustyon/services/api/internal/repository"
)

type ProjectHandler struct { repo *repository.ProjectRepository }
func NewProjectHandler(repo *repository.ProjectRepository) *ProjectHandler { return &ProjectHandler{repo: repo} }

func (h *ProjectHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	idText := strings.TrimPrefix(r.URL.Path, "/api/projects/")
	if idText != r.URL.Path {
		id, err := strconv.ParseInt(strings.Trim(idText, "/"), 10, 64)
		if err != nil || id < 1 { httpx.Error(w, http.StatusBadRequest, "invalid project id"); return }
		switch r.Method {
		case http.MethodPatch:
			var p model.Project
			if json.NewDecoder(r.Body).Decode(&p) != nil || strings.TrimSpace(p.Name) == "" || strings.TrimSpace(p.Slug) == "" { httpx.Error(w, http.StatusBadRequest, "name and slug are required"); return }
			if p.Status == "" { p.Status = "active" }
			if err := h.repo.Update(r.Context(), id, &p); err != nil { httpx.Error(w, http.StatusNotFound, "project not found"); return }
			httpx.Write(w, http.StatusOK, p)
		case http.MethodDelete:
			if err := h.repo.Delete(r.Context(), id); err != nil { httpx.Error(w, http.StatusInternalServerError, "could not delete project"); return }
			w.WriteHeader(http.StatusNoContent)
		default:
			w.Header().Set("Allow", "PATCH, DELETE")
			httpx.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		}
		return
	}

	switch r.Method {
	case http.MethodGet:
		items, err := h.repo.List(r.Context()); if err != nil { httpx.Error(w, http.StatusInternalServerError, "database error"); return }
		httpx.Write(w, http.StatusOK, items)
	case http.MethodPost:
		var p model.Project
		if json.NewDecoder(r.Body).Decode(&p) != nil || strings.TrimSpace(p.Name) == "" || strings.TrimSpace(p.Slug) == "" { httpx.Error(w, http.StatusBadRequest, "name and slug are required"); return }
		if p.Status == "" { p.Status = "active" }
		if err := h.repo.Create(r.Context(), &p); err != nil { httpx.Error(w, http.StatusConflict, "could not create project"); return }
		httpx.Write(w, http.StatusCreated, p)
	default:
		w.Header().Set("Allow", "GET, POST")
		httpx.Error(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}
