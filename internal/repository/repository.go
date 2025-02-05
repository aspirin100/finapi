package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"github.com/shopspring/decimal"

	"github.com/aspirin100/finapi/internal/entity"
	"github.com/aspirin100/finapi/internal/repository/migrations"
)

const (
	defaultTransactionsCount = 100
)

var (
	ErrUserNotFound    = errors.New("user not found")
	ErrNegativeBalance = errors.New("not enough money on balance")
)

type Repository struct {
	DB *pgxpool.Pool
}

type executor interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
}

func NewConnection(ctx context.Context, postgresDSN string) (*Repository, error) {
	conn, err := pgxpool.New(ctx, postgresDSN)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres: %w", err)
	}

	return &Repository{
		DB: conn,
	}, nil
}

func (r *Repository) UpMigrations(driver, dsn string) error {
	db, err := sql.Open(driver, dsn)
	if err != nil {
		return fmt.Errorf("open database error: %w", err)
	}

	goose.SetBaseFS(migrations.Migrations)

	err = goose.Up(db, ".")
	if err != nil {
		return fmt.Errorf("migrations up error: %w", err)
	}

	return nil
}

type ctxKey struct{}
type CommitOrRollback func(err error) error

var txContextKey = ctxKey{}

func (r *Repository) BeginTx(ctx context.Context) (context.Context, CommitOrRollback, error) {
	tx, err := r.DB.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.ReadCommitted,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to start transaction: %w", err)
	}

	return context.WithValue(ctx, txContextKey, tx), func(err error) error {
		if err != nil {
			errRollback := tx.Rollback(ctx)
			if errRollback != nil {
				return errors.Join(err, errRollback)
			}

			return err
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
	amount decimal.Decimal) (*decimal.Decimal, error) {
	ex := r.checkTx(ctx)

	rows, err := ex.Query(ctx,
		UpdateBalanceQuery, userID, amount)
	if err != nil {
		return nil, fmt.Errorf("failed to update user's balance: %w", err)
	}

	var currentBalance decimal.Decimal

	for rows.Next() {
		err = rows.Scan(&currentBalance)
		if err != nil {
			return nil, fmt.Errorf("failed to read user's balance from db: %w", err)
		}
	}

	err = rows.Err()
	if err != nil {
		var pgErr *pgconn.PgError

		switch {
		case errors.As(err, &pgErr) && pgErr.Code == "23514":
			return nil, ErrNegativeBalance
		default:
			return nil, fmt.Errorf("update balance query fail: %w", err)
		}
	}

	if rows.CommandTag().RowsAffected() == 0 {
		return nil, ErrUserNotFound
	}

	return &currentBalance, nil
}

func (r *Repository) SaveTransaction(ctx context.Context,
	receiverID, senderID uuid.UUID,
	amount decimal.Decimal,
	operation string) (*entity.Transaction, error) {
	ex := r.checkTx(ctx)

	transactionID := uuid.New()

	rows, err := ex.Query(ctx,
		NewTransactionQuery,
		transactionID,
		receiverID,
		senderID,
		amount,
		operation)
	if err != nil {
		return nil, fmt.Errorf("save transaction query error: %w", err)
	}

	var transaction entity.Transaction

	for rows.Next() {
		err = rows.Scan(&transaction.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("transactions scanning fail: %w", err)
		}
	}

	err = rows.Err()
	if err != nil {
		var pgErr *pgconn.PgError

		switch {
		case errors.As(err, &pgErr) && pgErr.Code == "23503":
			return nil, ErrUserNotFound
		default:
			return nil, fmt.Errorf("failed to save transaction: %w", err)
		}
	}

	transaction.ID = transactionID
	transaction.ReceiverID = receiverID
	transaction.SenderID = senderID
	transaction.Operation = operation
	transaction.Amount = amount

	return &transaction, nil
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

		err = rows.Scan(
			&transactions[i].ID,
			&transactions[i].ReceiverID,
			&transactions[i].SenderID,
			&transactions[i].Amount,
			&transactions[i].Operation,
			&transactions[i].CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scanning error: %w", err)
		}
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("error during read transactions: %w", err)
	}

	// for memory optimization
	result := make([]entity.Transaction, len(transactions))

	copy(result, transactions)

	return result, nil
}

func (r *Repository) checkTx(ctx context.Context) executor {
	var ex executor = r.DB

	// checks if current operation is in transaction
	tx, ok := ctx.Value(txContextKey).(pgx.Tx)
	if ok {
		ex = tx
	}

	return ex
}

const (
	UpdateBalanceQuery = `update bank_accounts set balance = (balance + $2) where userID = $1 returning balance`
	// UpdateBalanceQuery = `
	// UPDATE bank_accounts
	// SET balance = t.new_balance
	// FROM (
	// 	SELECT userid, balance + $2 AS new_balance
	// 	FROM bank_accounts
	// 	WHERE userid = $1
	// 	FOR UPDATE
	// ) AS t
	// WHERE bank_accounts.userid = t.userid
	// RETURNING bank_accounts.balance;`
	NewTransactionQuery = `insert into transactions(id, receiverID, senderID, amount, operation)
	values ($1, $2, $3, $4, $5)
	returning createdAt`
	GetTransactionsQuery = `select
	id, receiverID, senderID, amount, operation, createdAt
	from transactions
	where receiverID = $1 OR senderID = $1
	order by createdAt
	limit 10`
)
