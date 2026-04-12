package database

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

// ConnectPostgres creates a connection pool to PostgreSQL and validates it with a ping.
func ConnectPostgres(poolURL string) (*pgxpool.Pool, error) {
	poolConfig, err := pgxpool.ParseConfig(poolURL)
	if err != nil {
		return nil, fmt.Errorf("unable to parse database url: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %w", err)
	}

	if err := pool.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("unable to ping database: %w", err)
	}

	slog.Info("Successfully connected to PostgreSQL")
	return pool, nil
}
