package service

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/aspirin100/finapi/internal/repository"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

const (
	PostgresDSN = "postgres://postgres:postgres@localhost:5432/finapi?sslmode=disable"
)

var UserIDs []string = []string{
	"3fec06e9-29cc-4ff4-9ae7-fb0e7c757b61",
	"4178f61f-2ff9-4ab5-afa5-f30dc16e6ad9",
}

func initService() (*Service, error) {
	repo, err := repository.NewConnection(context.Background(), PostgresDSN)
	if err != nil {
		return nil, fmt.Errorf("failed to db connect: %w", err)
	}

	return New(repo), nil
}

func TestDeposit(t *testing.T) {
	ctx := context.Background()

	srvc, err := initService()
	if err != nil {
		log.Print(err)
		t.Fail()
	}

	type Params struct {
		UserID uuid.UUID
		Amount decimal.Decimal
	}

	cases := []struct {
		Name        string
		ExpectedErr error
		Request     Params
	}{
		{
			Name:        "ok case",
			ExpectedErr: nil,
			Request: Params{
				UserID: uuid.MustParse(UserIDs[0]),
				Amount: decimal.NewFromFloat(9999999),
			},
		},
		{
			Name:        "user not found case",
			ExpectedErr: ErrUserNotFound,
			Request: Params{
				UserID: uuid.Nil,
				Amount: decimal.NewFromFloat(100),
			},
		},
	}

	for _, tcase := range cases {
		t.Run(tcase.Name, func(t *testing.T) {
			err := srvc.Deposit(ctx,
				tcase.Request.UserID,
				tcase.Request.Amount)

			require.EqualValues(t, tcase.ExpectedErr, err)
		})
	}
}
