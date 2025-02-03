package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"github.com/aspirin100/finapi/internal/entity"
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

func (s *Service) Deposit(userID uuid.UUID, amount decimal.Decimal) error {
	return nil
}

func (s *Service) GetTransactions(userID uuid.UUID) error {
	return nil
}

func (s *Service) Transfer(senderID, receiverID uuid.UUID, amount decimal.Decimal) error {
	return nil
}
