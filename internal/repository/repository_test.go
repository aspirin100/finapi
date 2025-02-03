package repository

import (
	"context"
	"log"
	"testing"

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

func TestUpdateBalance(t *testing.T) {
	ctx := context.Background()

	repo, err := NewConnection(ctx, PostgresDSN)
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
				Amount: decimal.NewFromFloat(10000),
			},
		},
		{
			Name:        "user not found",
			ExpectedErr: ErrUserNotFound,
			Request: Params{
				UserID: uuid.Nil,
				Amount: decimal.NewFromFloat(1111),
			},
		},
		{
			Name:        "not enough money case",
			ExpectedErr: ErrNegativeBalance,
			Request: Params{
				UserID: uuid.MustParse(UserIDs[0]),
				Amount: decimal.NewFromFloat(-100000),
			},
		},
	}

	for _, tcase := range cases {
		t.Run(tcase.Name, func(t *testing.T) {
			err := repo.UpdateBalance(ctx, tcase.Request.UserID, tcase.Request.Amount)

			require.EqualValues(t, tcase.ExpectedErr, err)
		})
	}
}

func TestSaveTransaction(t *testing.T) {
	ctx := context.Background()

	repo, err := NewConnection(ctx, PostgresDSN)
	if err != nil {
		log.Print(err)
		t.Fail()
	}

	type Params struct {
		ReceiverID uuid.UUID
		SenderID   uuid.UUID
		Amount     decimal.Decimal
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
				ReceiverID: uuid.MustParse(UserIDs[0]),
				SenderID:   uuid.MustParse(UserIDs[1]),
				Amount:     decimal.NewFromFloat(10000),
			},
		},
		{
			Name:        "user not found",
			ExpectedErr: ErrUserNotFound,
			Request: Params{
				ReceiverID: uuid.Nil,
				SenderID:   uuid.Nil,
				Amount:     decimal.NewFromFloat(10000),
			},
		},
	}

	for _, tcase := range cases {
		t.Run(tcase.Name, func(t *testing.T) {
			err := repo.SaveTransaction(
				ctx,
				tcase.Request.ReceiverID,
				tcase.Request.SenderID,
				tcase.Request.Amount)

			require.EqualValues(t, tcase.ExpectedErr, err)
		})
	}
}
