package service

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"github.com/aspirin100/finapi/internal/entity"
)

type UserManager interface {
	UpdateBalance(userID uuid.UUID, amount decimal.Decimal) error
	GetTransactions(userID uuid.UUID) ([]entity.Transaction, error)
}

type Service struct {
	userManager UserManager
}

func New(userManager UserManager) *Service {
	return &Service{
		userManager: userManager,
	}
}
