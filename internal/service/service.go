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

type UserManager interface {
	UpdateBalance(ctx context.Context, userID uuid.UUID, amount decimal.Decimal) error
	GetTransactions(ctx context.Context, userID uuid.UUID) ([]entity.Transaction, error)
	SaveTransaction(ctx context.Context, receiverID, senderID uuid.UUID, amount decimal.Decimal) error
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

func (s *Service) Deposit(ctx context.Context, userID uuid.UUID, amount decimal.Decimal) error {
	err := s.userManager.UpdateBalance(ctx, userID, amount)
	if err != nil {
		responseOnRepoError(err)
	}

	return nil
}

func (s *Service) GetTransactions(ctx context.Context, userID uuid.UUID) error {
	return nil
}

func (s *Service) Transfer(ctx context.Context, receiverID, senderID uuid.UUID, amount decimal.Decimal) error {
	ctx, commitOrRollback, err := s.userManager.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin db transaction: %w", err)
	}

	defer func() {
		errTx := commitOrRollback(&err)
		if errTx != nil {
			err = errTx
		}
	}()

	// sender balance update
	err = s.userManager.UpdateBalance(
		ctx,
		senderID,
		decimal.Zero.Sub(amount))
	if err != nil {
		return responseOnRepoError(err)
	}
	// receiver balance update
	err = s.userManager.UpdateBalance(
		ctx,
		receiverID,
		amount)
	if err != nil {
		return responseOnRepoError(err)
	}

	err = s.userManager.SaveTransaction(ctx, receiverID, senderID, amount)
	if err != nil {
		return responseOnRepoError(err)
	}

	return nil
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
