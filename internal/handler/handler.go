package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/aspirin100/finapi/internal/entity"
	"github.com/aspirin100/finapi/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

var (
	ErrInvalidFormat  = errors.New("invalid user id format")
	ErrNegativeAmount = errors.New("deposit amount must be positive")
	ErrSameUser       = errors.New("receiver and sender must be different person")
)

type depositRequestParams struct {
	UserID uuid.UUID       `json:"omitempty"`
	Amount decimal.Decimal `json:"amount"`
}

type transferRequestParams struct {
	SenderID   uuid.UUID       `json:"omitempty"`
	ReceiverID uuid.UUID       `json:"receiverID"`
	Amount     decimal.Decimal `json:"amount"`
}

type TransactionManager interface {
	Deposit(ctx context.Context, userID uuid.UUID, amount decimal.Decimal) error
	GetTransactions(ctx context.Context, userID uuid.UUID) ([]entity.Transaction, error)
	Transfer(ctx context.Context, receiverID, senderID uuid.UUID, amount decimal.Decimal) error
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
	router.PATCH("/:userID/transfer", handler.TransferMoney)

	srv := &http.Server{
		Addr:    hostname + ":" + port,
		Handler: router,
	}

	handler.Server = srv

	return handler
}

func (h *Handler) GetUserTransactions(ctx *gin.Context) {
	userIDarsed, err := uuid.Parse(ctx.Param("userID"))
	if err != nil {
		ctx.String(http.StatusNotFound, "user not found")
	}

	transactions, err := h.tmanager.GetTransactions(
		ctx,
		userIDarsed)
	if err != nil {
		responseOnServiceError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, transactions)
}

func (h *Handler) Deposit(ctx *gin.Context) {
	params, err := validateDepositRequest(ctx.Param("userID"), ctx.Request)
	if err != nil {
		responseOnValidationErr(ctx, err)
	}

	err = h.tmanager.Deposit(ctx, params.UserID, params.Amount)
	if err != nil {
		responseOnServiceError(ctx, err)
		return
	}

	ctx.Status(http.StatusOK)
}

func (h *Handler) TransferMoney(ctx *gin.Context) {
	params, err := validateTransferRequest(ctx.Param("userID"), ctx.Request)
	if err != nil {
		responseOnValidationErr(ctx, err)
	}

	err = h.tmanager.Transfer(
		ctx,
		params.ReceiverID,
		params.SenderID,
		params.Amount)
	if err != nil {
		responseOnServiceError(ctx, err)
		return
	}

	ctx.Status(http.StatusOK)
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

	if decimal.Zero.Compare(params.Amount) >= 0 {
		return nil, ErrNegativeAmount
	}

	params.UserID = useridParsed

	return &params, nil
}

func validateTransferRequest(
	userID string,
	req *http.Request) (*transferRequestParams, error) {
	senderIDParsed, err := uuid.Parse(userID)
	if err != nil {
		return nil, ErrInvalidFormat
	}

	var params transferRequestParams

	decoder := json.NewDecoder(req.Body)

	err = decoder.Decode(&params)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal request body: %w", err)
	}

	receiverIDParsed, err := uuid.Parse(params.ReceiverID.String())
	if err != nil {
		return nil, ErrInvalidFormat
	}

	if receiverIDParsed == senderIDParsed {
		return nil, ErrSameUser
	}

	if decimal.Zero.Compare(params.Amount) >= 0 {
		return nil, ErrNegativeAmount
	}

	return &transferRequestParams{
		SenderID:   senderIDParsed,
		ReceiverID: receiverIDParsed,
		Amount:     params.Amount,
	}, nil
}

func responseOnValidationErr(ctx *gin.Context, err error) {
	switch {
	case errors.Is(err, ErrInvalidFormat):
		ctx.String(http.StatusBadRequest, "wrong user id format")
	case errors.Is(err, ErrNegativeAmount):
		ctx.String(http.StatusBadRequest, "amount must be positive")
	case errors.Is(err, ErrSameUser):
		ctx.String(http.StatusBadRequest, "can't transfer money to the same account")
	default:
		ctx.Status(http.StatusInternalServerError)
	}
}

func responseOnServiceError(ctx *gin.Context, err error) {
	switch {
	case errors.Is(err, service.ErrUserNotFound):
		ctx.String(http.StatusNotFound, "user not found")
	case errors.Is(err, service.ErrNegativeBalance):
		ctx.String(http.StatusBadRequest, "not enough money on account")
	default:
		ctx.Status(http.StatusInternalServerError)
	}
}
