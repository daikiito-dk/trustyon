package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/daikiito-dk/trustyon/services/api/internal/model"
	"github.com/daikiito-dk/trustyon/services/api/internal/repository"
)

type NoteHandler struct { repo *repository.NoteRepository }

func NewNoteHandler(repo *repository.NoteRepository) *NoteHandler { return &NoteHandler{repo: repo} }

func (h *NoteHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	switch r.Method {
	case http.MethodGet:
		items, err := h.repo.List(r.Context())
		if err != nil { http.Error(w, `{"error":"database error"}`, http.StatusInternalServerError); return }
		json.NewEncoder(w).Encode(items)
	case http.MethodPost:
		var n model.Note
		if json.NewDecoder(r.Body).Decode(&n) != nil || strings.TrimSpace(n.Title) == "" {
			http.Error(w, `{"error":"title is required"}`, http.StatusBadRequest); return
		}
		if err := h.repo.Create(r.Context(), &n); err != nil { http.Error(w, `{"error":"could not create note"}`, http.StatusInternalServerError); return }
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(n)
	default:
		w.Header().Set("Allow", "GET, POST")
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
	}
}
