package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"github.com/aspirin100/finapi/internal/entity"
	"github.com/aspirin100/finapi/internal/repository"
)

var (
	ErrUserNotFound    = errors.New("user not found")
	ErrNegativeBalance = errors.New("not enough money on balance")
)

const (
	operationTransfer = "transfer"
	operationDeposit  = "deposit"
)

type UserManager interface {
	UpdateBalance(ctx context.Context, userID uuid.UUID, amount decimal.Decimal) (*decimal.Decimal, error)
	GetTransactions(ctx context.Context, userID uuid.UUID) ([]entity.Transaction, error)
	SaveTransaction(ctx context.Context,
		receiverID,
		senderID uuid.UUID,
		amount decimal.Decimal,
		operation string) (*entity.Transaction, error)
	BeginTx(ctx context.Context) (context.Context, repository.CommitOrRollback, error)
}

type Service struct {
	userManager UserManager
}

func New(userManager UserManager) *Service {
	return &Service{
		userManager: userManager,
	}
}

func (s *Service) Deposit(ctx context.Context, userID uuid.UUID, amount decimal.Decimal) (*decimal.Decimal, error) {
	ctx, commitOrRollback, err := s.userManager.BeginTx(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin bd transaction: %w", err)
	}

	defer func() {
		errTx := commitOrRollback(&err)
		if errTx != nil {
			err = errTx
		}
	}()

	currentBalance, err := s.userManager.UpdateBalance(ctx, userID, amount)
	if err != nil {
		return nil, responseOnRepoError(err)
	}

	_, err = s.userManager.SaveTransaction(ctx, userID, userID, amount, operationDeposit)
	if err != nil {
		return nil, responseOnRepoError(err)
	}

	return currentBalance, nil
}

func (s *Service) GetTransactions(ctx context.Context,
	userID uuid.UUID) ([]entity.Transaction, error) {
	transactions, err := s.userManager.GetTransactions(ctx, userID)
	if err != nil {
		return nil, responseOnRepoError(err)
	}

	return transactions, nil
}

func (s *Service) Transfer(ctx context.Context, receiverID, senderID uuid.UUID, amount decimal.Decimal) (*entity.Transaction, error) {
	ctx, commitOrRollback, err := s.userManager.BeginTx(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin db transaction: %w", err)
	}

	defer func() {
		errTx := commitOrRollback(&err)
		if errTx != nil {
			err = errTx
		}
	}()

	// sender balance update
	_, err = s.userManager.UpdateBalance(
		ctx,
		senderID,
		decimal.Zero.Sub(amount))
	if err != nil {
		return nil, responseOnRepoError(err)
	}
	// receiver balance update
	_, err = s.userManager.UpdateBalance(
		ctx,
		receiverID,
		amount)
	if err != nil {
		return nil, responseOnRepoError(err)
	}

	transaction, err := s.userManager.SaveTransaction(ctx, receiverID, senderID, amount, operationTransfer)
	if err != nil {
		return nil, responseOnRepoError(err)
	}

	return transaction, nil
}

func responseOnRepoError(err error) error {
	switch {
	case errors.Is(err, repository.ErrNegativeBalance):
		return ErrNegativeBalance
	case errors.Is(err, repository.ErrUserNotFound):
		return ErrUserNotFound
	default:
		return fmt.Errorf("repository fail: %w", err)
	}
}
