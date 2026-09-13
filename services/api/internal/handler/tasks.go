package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/daikiito-dk/trustyon/services/api/internal/model"
	"github.com/daikiito-dk/trustyon/services/api/internal/repository"
)

type TaskHandler struct { repo *repository.TaskRepository }

func NewTaskHandler(repo *repository.TaskRepository) *TaskHandler { return &TaskHandler{repo: repo} }

func (h *TaskHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	switch r.Method {
	case http.MethodGet:
		items, err := h.repo.List(r.Context())
		if err != nil { http.Error(w, `{"error":"database error"}`, http.StatusInternalServerError); return }
		json.NewEncoder(w).Encode(items)
	case http.MethodPost:
		var t model.Task
		if json.NewDecoder(r.Body).Decode(&t) != nil || strings.TrimSpace(t.Title) == "" {
			http.Error(w, `{"error":"title is required"}`, http.StatusBadRequest); return
		}
		if t.Status == "" { t.Status = "todo" }
		if t.Priority == "" { t.Priority = "medium" }
		if err := h.repo.Create(r.Context(), &t); err != nil { http.Error(w, `{"error":"could not create task"}`, http.StatusInternalServerError); return }
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(t)
	default:
		w.Header().Set("Allow", "GET, POST")
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
	}
}
