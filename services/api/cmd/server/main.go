package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/daikiito-dk/trustyon/services/api/internal/handler"
	"github.com/daikiito-dk/trustyon/services/api/internal/repository"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db, err := repository.NewPool(ctx)
	if err != nil { log.Fatal(err) }
	defer db.Close()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if err := db.Ping(r.Context()); err != nil { w.WriteHeader(http.StatusServiceUnavailable); json.NewEncoder(w).Encode(map[string]string{"status":"degraded"}); return }
		json.NewEncoder(w).Encode(map[string]string{"status":"ok"})
	})
	mux.HandleFunc("GET /api/github/activity", handler.GitHubActivityHandler)

	projects := repository.NewProjectRepository(db)
	tasks := repository.NewTaskRepository(db)
	notes := repository.NewNoteRepository(db)
	projectHandler := handler.NewProjectHandler(projects)
	taskHandler := handler.NewTaskHandler(tasks)
	noteHandler := handler.NewNoteHandler(notes)
	mux.Handle("/api/projects", projectHandler)
	mux.Handle("/api/projects/", projectHandler)
	mux.Handle("/api/tasks", taskHandler)
	mux.Handle("/api/tasks/", taskHandler)
	mux.Handle("/api/notes", noteHandler)
	mux.Handle("/api/notes/", noteHandler)

	server := &http.Server{Addr: ":8080", Handler: logging(mux)}
	log.Println("Trustyon API listening on :8080")
	log.Fatal(server.ListenAndServe())
}

func logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s", r.Method, r.URL.Path, time.Since(start))
	})
}
