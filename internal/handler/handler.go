package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/aspirin100/finapi/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

var (
	ErrInvalidFormat = errors.New("invalid user id format")
	ErrEmptyRequest  = errors.New("request body is required")
)

type depositRequestParams struct {
	UserID uuid.UUID       `json:"omitempty"`
	Amount decimal.Decimal `json:"amount"`
}

type TransactionManager interface {
	Deposit(ctx context.Context, userID uuid.UUID, amount decimal.Decimal) error
	GetTransactions(ctx context.Context, userID uuid.UUID) error
	Transfer(ctx context.Context, senderID, receiverID uuid.UUID, amount decimal.Decimal) error
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
		Addr:    hostname + ":" + port,
		Handler: router,
	}

	handler.Server = srv

	return handler
}

func (h *Handler) GetUserTransactions(ctx *gin.Context) {

}

func (h *Handler) Deposit(ctx *gin.Context) {
	params, err := validateDepositRequest(ctx.GetString("userID"), ctx.Request)
	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidFormat):
			ctx.String(http.StatusBadRequest, "wrong user id format")
		case errors.Is(err, ErrEmptyRequest):
			ctx.String(http.StatusBadRequest, "request body is required")
		default:
			ctx.Status(http.StatusInternalServerError)
		}
	}

	err = h.tmanager.Deposit(ctx, params.UserID, params.Amount)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrUserNotFound):
			ctx.String(http.StatusNotFound, "user not found")
		case errors.Is(err, service.ErrNegativeBalance):
			ctx.String(http.StatusBadRequest, "not enough money on account")
		default:
			ctx.Status(http.StatusInternalServerError)
		}
	}

	ctx.Status(http.StatusOK)
}

func (h *Handler) TransferMoney(ctx *gin.Context) {

}

func validateDepositRequest(
	userID string,
	req *http.Request) (*depositRequestParams, error) {
	useridParsed, err := uuid.Parse(userID)
	if err != nil {
		return nil, ErrInvalidFormat
	}

	var params depositRequestParams

	decoder := json.NewDecoder(req.Body)

	err = decoder.Decode(&params)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal request body: %w", err)
	}

	params.UserID = useridParsed

	return &params, nil
}
