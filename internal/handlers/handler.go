package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type TransactionManager interface {
	Deposit(userID uuid.UUID, amount decimal.Decimal) error
	GetTransactions(userID uuid.UUID) error
	Transfer(senderID, receiverID uuid.UUID, amount decimal.Decimal) error
}

type Handler struct {
	Server   *http.Server
	tmanager TransactionManager
}

func New(hostname, port string, tmanager TransactionManager) *Handler {
	handler := &Handler{
		tmanager: tmanager,
	}

	router := gin.Default()
	
	router.GET("/:userID/transactions", handler.GetUserTransactions)
	router.PATCH("/:userID/balance", handler.Deposit)
	router.PATCH("/transfer", handler.TransferMoney)

	srv := &http.Server{
		Addr: hostname + ":" + port,
		Handler: router,
	}

	handler.Server = srv

	return handler
}

func (h *Handler) GetUserTransactions(ctx *gin.Context) {

}

func (h *Handler) Deposit(ctx *gin.Context) {

}

func (h *Handler) TransferMoney(ctx *gin.Context) {

}
