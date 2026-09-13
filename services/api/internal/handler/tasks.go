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

type TaskHandler struct { repo *repository.TaskRepository }

func NewTaskHandler(repo *repository.TaskRepository) *TaskHandler { return &TaskHandler{repo: repo} }

func (h *TaskHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	idText := strings.TrimPrefix(r.URL.Path, "/api/tasks/")
	if idText != r.URL.Path {
		id, err := strconv.ParseInt(strings.Trim(idText, "/"), 10, 64)
		if err != nil || id < 1 { httpx.Error(w, http.StatusBadRequest, "invalid task id"); return }
		switch r.Method {
		case http.MethodPatch:
			var t model.Task
			if json.NewDecoder(r.Body).Decode(&t) != nil || strings.TrimSpace(t.Title) == "" {
				httpx.Error(w, http.StatusBadRequest, "title is required"); return
			}
			if t.Status == "" { t.Status = "todo" }
			if t.Priority == "" { t.Priority = "medium" }
			if err := h.repo.Update(r.Context(), id, &t); err != nil { httpx.Error(w, http.StatusNotFound, "task not found"); return }
			httpx.Write(w, http.StatusOK, t)
		case http.MethodDelete:
			if err := h.repo.Delete(r.Context(), id); err != nil { httpx.Error(w, http.StatusInternalServerError, "could not delete task"); return }
			w.WriteHeader(http.StatusNoContent)
		default:
			w.Header().Set("Allow", "PATCH, DELETE")
			httpx.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		}
		return
	}

	switch r.Method {
	case http.MethodGet:
		items, err := h.repo.List(r.Context())
		if err != nil { httpx.Error(w, http.StatusInternalServerError, "database error"); return }
		httpx.Write(w, http.StatusOK, items)
	case http.MethodPost:
		var t model.Task
		if json.NewDecoder(r.Body).Decode(&t) != nil || strings.TrimSpace(t.Title) == "" {
			httpx.Error(w, http.StatusBadRequest, "title is required"); return
		}
		if t.Status == "" { t.Status = "todo" }
		if t.Priority == "" { t.Priority = "medium" }
		if err := h.repo.Create(r.Context(), &t); err != nil { httpx.Error(w, http.StatusInternalServerError, "could not create task"); return }
		httpx.Write(w, http.StatusCreated, t)
	default:
		w.Header().Set("Allow", "GET, POST")
		httpx.Error(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}
