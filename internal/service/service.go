package service

import (
	"github.com/shopspring/decimal"
	"github.com/google/uuid"

	"github.com/aspirin100/finapi/internal/repository"
)

type BankManager interface {
	Deposit(userID uuid.UUID, amount decimal.Decimal) error
	GetTransactions(userID uuid.UUID) 
}

type Service struct {
	Storage *repository.Repository
}
