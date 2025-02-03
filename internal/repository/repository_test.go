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
	UserID      = "3fec06e9-29cc-4ff4-9ae7-fb0e7c757b61"
)

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
				UserID: uuid.MustParse(UserID),
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
	}

	for _, tcase := range cases {
		t.Run(tcase.Name, func(t *testing.T) {
			err := repo.UpdateBalance(ctx, tcase.Request.UserID, tcase.Request.Amount)

			require.EqualValues(t, tcase.ExpectedErr, err)
		})
	}
}
