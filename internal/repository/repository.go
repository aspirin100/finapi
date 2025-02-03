package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/shopspring/decimal"

	"github.com/aspirin100/finapi/internal/entity"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type Repository struct {
	DB *pgx.Conn
}

type executor interface {
	Exec(ctx context.Context, query string, args ...any) (pgconn.CommandTag, error)
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

type ctxKey struct{}
type CommitOrRollback func(err *error) error

var txContextKey = ctxKey{}

func (r *Repository) BeginTx(ctx context.Context) (context.Context, CommitOrRollback, error) {
	tx, err := r.DB.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.ReadCommitted,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to start transaction: %w", err)
	}

	return context.WithValue(ctx, txContextKey, tx), func(err *error) error {
		if *err != nil {
			if errRollback := tx.Rollback(ctx); errRollback != nil {
				return errors.Join(*err, errRollback)
			}

			return *err
		}

		if errCommit := tx.Commit(ctx); errCommit != nil {
			return fmt.Errorf("failed to commit transaction: %w", errCommit)
		}

		return nil
	}, nil
}

func (r *Repository) GetTransactions(ctx context.Context,
	userID uuid.UUID) ([]entity.Transaction, error) {
	return nil, nil
}

func (r *Repository) UpdateBalance(ctx context.Context,
	userID uuid.UUID,
	amount decimal.Decimal) error {
	var ex executor = r.DB

	// checks if current operation is in transaction
	tx, ok := ctx.Value(txContextKey).(*pgx.Tx)
	if ok {
		ex = *tx
	}

	comm, err := ex.Exec(ctx,
		UpdateBalanceQuery, userID, amount)
	if err != nil {
		return fmt.Errorf("failed to update user's balance: %w", err)
	} else if comm.RowsAffected() == 0 {
		return ErrUserNotFound
	}

	return nil
}

const (
	UpdateBalanceQuery   = `update bank_accounts set balance = (balance + $2) where userID = $1`
	GetTransactionsQuery = `select (id, receiverID, senderID, amount, createdAt)
	from transactions
	where receiverID = $1 OR senderID = $1`
)
