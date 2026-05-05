package main

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Opens a PostgreSQL connection pool using the provided connection string
func connectDB(ctx context.Context, connStr string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		return nil, err
	}
	return pool, nil
}
