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
	SaveTransaction(ctx context.Context, senderID, receiverID uuid.UUID, amount decimal.Decimal) error
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
		switch {
		case errors.Is(err, repository.ErrUserNotFound):
			return ErrUserNotFound
		default:
			return fmt.Errorf("failed to update balance: %w", err)
		}
	}

	return nil
}

func (s *Service) GetTransactions(ctx context.Context, userID uuid.UUID) error {
	return nil
}

func (s *Service) Transfer(ctx context.Context, senderID, receiverID uuid.UUID, amount decimal.Decimal) error {
	return nil
}
