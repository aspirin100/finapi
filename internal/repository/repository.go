package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/shopspring/decimal"

	"github.com/aspirin100/finapi/internal/entity"
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

func (r *Repository) GetTransactions(ctx context.Context,
	userID uuid.UUID) ([]entity.Transaction, error) {
	return nil, nil
}

func (r *Repository) UpdateBalance(ctx context.Context,
	userID uuid.UUID,
	amount decimal.Decimal) error {
	return nil
}
