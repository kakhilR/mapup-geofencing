package repository

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func ConnectDB() (*pgxpool.Pool, error) {
	connString := os.Getenv("DATABASE_URL")

	if connString == "" {
		connString = "postgres://mapup:mapup_password@localhost:5432/mapup"
	}

	pool, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		return nil, fmt.Errorf("failed to create database pool: %w", err)
	}

	if err := pool.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return pool, nil
}