package httpx

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWrite(t *testing.T) {
	w := httptest.NewRecorder()
	Write(w, http.StatusCreated, map[string]string{"status": "created"})

	if w.Code != http.StatusCreated {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusCreated)
	}
	if got := w.Header().Get("Content-Type"); got != "application/json" {
		t.Fatalf("content type = %q, want application/json", got)
	}
	if got := w.Body.String(); got != "{\"status\":\"created\"}\n" {
		t.Fatalf("body = %q", got)
	}
}

func TestError(t *testing.T) {
	w := httptest.NewRecorder()
	Error(w, http.StatusBadRequest, "invalid input")

	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
	if got := w.Body.String(); got != "{\"error\":\"invalid input\"}\n" {
		t.Fatalf("body = %q", got)
	}
}
