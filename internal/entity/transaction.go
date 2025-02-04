package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Transaction struct {
	ID         uuid.UUID       `json:"id"`
	SenderID   uuid.UUID       `json:"senderID"`   //nolint:tagliatelle
	ReceiverID uuid.UUID       `json:"receiverID"` //nolint:tagliatelle
	Amount     decimal.Decimal `json:"amount"`
	Operation  string          `json:"operation"`
	CreatedAt  time.Time       `json:"createdAt"`
}