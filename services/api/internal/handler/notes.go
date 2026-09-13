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

type NoteHandler struct { repo *repository.NoteRepository }

func NewNoteHandler(repo *repository.NoteRepository) *NoteHandler { return &NoteHandler{repo: repo} }

func (h *NoteHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	idText := strings.TrimPrefix(r.URL.Path, "/api/notes/")
	if idText != r.URL.Path {
		id, err := strconv.ParseInt(strings.Trim(idText, "/"), 10, 64)
		if err != nil || id < 1 { httpx.Error(w, http.StatusBadRequest, "invalid note id"); return }
		switch r.Method {
		case http.MethodPatch:
			var n model.Note
			if json.NewDecoder(r.Body).Decode(&n) != nil || strings.TrimSpace(n.Title) == "" {
				httpx.Error(w, http.StatusBadRequest, "title is required"); return
			}
			if err := h.repo.Update(r.Context(), id, &n); err != nil { httpx.Error(w, http.StatusNotFound, "note not found"); return }
			httpx.Write(w, http.StatusOK, n)
		case http.MethodDelete:
			if err := h.repo.Delete(r.Context(), id); err != nil { httpx.Error(w, http.StatusInternalServerError, "could not delete note"); return }
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
		var n model.Note
		if json.NewDecoder(r.Body).Decode(&n) != nil || strings.TrimSpace(n.Title) == "" {
			httpx.Error(w, http.StatusBadRequest, "title is required"); return
		}
		if err := h.repo.Create(r.Context(), &n); err != nil { httpx.Error(w, http.StatusInternalServerError, "could not create note"); return }
		httpx.Write(w, http.StatusCreated, n)
	default:
		w.Header().Set("Allow", "GET, POST")
		httpx.Error(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}
