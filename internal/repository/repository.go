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

const (
	defaultTransactionsCount = 100
)

var (
	ErrUserNotFound    = errors.New("user not found")
	ErrNegativeBalance = errors.New("not enough money on balance")
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
			errRollback := tx.Rollback(ctx)
			if errRollback != nil {
				return errors.Join(*err, errRollback)
			}

			return *err
		}

		errCommit := tx.Commit(ctx)
		if errCommit != nil {
			return fmt.Errorf("failed to commit transaction: %w", errCommit)
		}

		return nil
	}, nil
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
		var pgErr *pgconn.PgError

		switch {
		case errors.As(err, &pgErr) && pgErr.Code == "23514":
			return ErrNegativeBalance
		default:
			return fmt.Errorf("failed to update user's balance: %w", err)
		}
	} else if comm.RowsAffected() == 0 {
		return ErrUserNotFound
	}

	return nil
}

func (r *Repository) SaveTransaction(ctx context.Context,
	receiverID, senderID uuid.UUID,
	amount decimal.Decimal) error {
	var ex executor = r.DB

	// checks if current operation is in transaction
	tx, ok := ctx.Value(txContextKey).(*pgx.Tx)
	if ok {
		ex = *tx
	}

	transactionID := uuid.New()

	_, err := ex.Exec(ctx,
		NewTransactionQuery,
		transactionID,
		receiverID,
		senderID,
		amount)
	if err != nil {
		var pgErr *pgconn.PgError

		switch {
		case errors.As(err, &pgErr) && pgErr.Code == "23503":
			return ErrUserNotFound
		default:
			return fmt.Errorf("failed to save transaction: %w", err)
		}
	}

	return nil
}

func (r *Repository) GetTransactions(ctx context.Context,
	userID uuid.UUID) ([]entity.Transaction, error) {
	rows, err := r.DB.Query(
		ctx,
		GetTransactionsQuery,
		userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get users's transactions: %w", err)
	}

	transactions := make([]entity.Transaction, 0, defaultTransactionsCount)

	for i := 0; rows.Next(); i++ {
		transactions = append(transactions, entity.Transaction{})

		rows.Scan(
			&transactions[i].ID,
			&transactions[i].ReceiverID,
			&transactions[i].SenderID,
			&transactions[i].Amount,
			&transactions[i].CreatedAt,
		)
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("error during read transaction rows: %w", err)
	}

	// for memory optimization
	result := make([]entity.Transaction, len(transactions))

	copy(result, transactions)

	return result, nil
}

const (
	UpdateBalanceQuery  = `update bank_accounts set balance = (balance + $2) where userID = $1`
	NewTransactionQuery = `insert into transactions(id, receiverID, senderID, amount)
	values ($1, $2, $3, $4)`
	GetTransactionsQuery = `select id, receiverID, senderID, amount, createdAt
	from transactions
	where receiverID = $1 OR senderID = $1`
	GetTransactionsCountQuery = `select count`
)
