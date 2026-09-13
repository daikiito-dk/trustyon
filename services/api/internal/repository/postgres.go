package repository

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPool(ctx context.Context) (*pgxpool.Pool, error) {
	url := os.Getenv("DATABASE_URL")
	if url == "" {
		url = "postgres://trustyon:trustyon@localhost:5432/trustyon?sslmode=disable"
	}
	return pgxpool.New(ctx, url)
}
