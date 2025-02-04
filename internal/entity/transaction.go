package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Transaction struct {
	ID         uuid.UUID       `json:"id"`
	SenderID   uuid.UUID       `json:"senderID"`
	ReceiverID uuid.UUID       `json:"receiverID"`
	Amount     decimal.Decimal `json:"amount"`
	Operation  string          `json:"operation"`
	CreatedAt  time.Time       `json:"createdAt"`
}
