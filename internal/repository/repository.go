package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/shopspring/decimal"

	"github.com/aspirin100/finapi/internal/entity"
)

var (
	ErrUserNotFound = errors.New("user not found")
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
	_, err := r.DB.Exec(ctx,
		UpdateBalanceQuery, userID, amount)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrUserNotFound
		}

		return fmt.Errorf("failed to update user's balance: %w", err)
	}


	return nil
}

const (
	UpdateBalanceQuery   = `update table bank_accounts set balance = balance + $2 where userID = $1`
	GetTransactionsQuery = `select (id, receiverID, senderID, amount, createdAt)
	from transactions
	where receiverID = $1 OR senderID = $1`
)
