package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type Repository struct {
	DB *pgx.Conn
}

func NewConnection(ctx context.Context, postgresDSN string) (*Repository, error) {
	conn, err := pgx.Connect(ctx, postgresDSN)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres: %w", err)
	}

	return &Repository{
		DB: conn,
	}, nil
}
