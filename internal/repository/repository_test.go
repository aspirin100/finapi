package repository_test

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

func TestUpdateBalance(t *testing.T) {
	ctx := context.Background()

	repo, err := repository.NewConnection(ctx, PostgresDSN)
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
				Amount: decimal.NewFromFloat(1),
			},
		},
		{
			Name:        "user not found",
			ExpectedErr: repository.ErrUserNotFound,
			Request: Params{
				UserID: uuid.Nil,
				Amount: decimal.NewFromFloat(1111),
			},
		},
		{
			Name:        "not enough money case",
			ExpectedErr: repository.ErrNegativeBalance,
			Request: Params{
				UserID: uuid.MustParse(UserIDs[0]),
				Amount: decimal.NewFromFloat(-100000),
			},
		},
	}

	//wg := sync.WaitGroup{}

	// for i := 0; i < 10000; i++ {
	// 	wg.Add(1)
	// 	go func() {
	// 		defer wg.Done()

	// 		ctx, cor, err := repo.BeginTx(ctx)
	// 		if err != nil {
	// 			log.Println(err)
	// 			t.Fail()
	// 		}

	// 		defer func() {
	// 			err = cor(err)
	// 			if err != nil {
	// 				log.Println(err)
	// 				t.Fail()
	// 			}
	// 		}()

	// 		_, err = repo.UpdateBalance(ctx,
	// 			cases[0].Request.UserID,
	// 			cases[0].Request.Amount)
	// 		if err != nil {
	// 			log.Println(err)
	// 		}

	// 	}()
	// }

	// wg.Wait()

	for _, tcase := range cases {
		t.Run(tcase.Name, func(t *testing.T) {
			balance, err := repo.UpdateBalance(ctx, tcase.Request.UserID, tcase.Request.Amount)

			require.EqualValues(t, tcase.ExpectedErr, err)
			log.Println("current balance:", balance)
		})
	}
}

func TestSaveTransaction(t *testing.T) {
	ctx := context.Background()

	repo, err := repository.NewConnection(ctx, PostgresDSN)
	if err != nil {
		log.Print(err)
		t.Fail()
	}

	type Params struct {
		ReceiverID uuid.UUID
		SenderID   uuid.UUID
		Amount     decimal.Decimal
		Operation  string
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
				SenderID:   uuid.MustParse(UserIDs[0]),
				Amount:     decimal.NewFromFloat(1),
				Operation:  "deposit",
			},
		},
		{
			Name:        "user not found",
			ExpectedErr: repository.ErrUserNotFound,
			Request: Params{
				ReceiverID: uuid.Nil,
				SenderID:   uuid.Nil,
				Amount:     decimal.NewFromFloat(10000),
				Operation:  "transfer",
			},
		},
	}

	// wg := sync.WaitGroup{}

	// for i := 0; i < 10000; i++ {
	// 	wg.Add(1)
	// 	go func() {
	// 		defer wg.Done()
	// 		_, err := repo.SaveTransaction(
	// 			ctx,
	// 			cases[0].Request.ReceiverID,
	// 			cases[0].Request.SenderID,
	// 			cases[0].Request.Amount,
	// 			cases[0].Request.Operation)
	// 		if err != nil {
	// 			log.Print(err)
	// 			t.Fail()
	// 		}
	// 	}()
	// }

	// wg.Wait()

	for _, tcase := range cases {
		t.Run(tcase.Name, func(t *testing.T) {
			tx, err := repo.SaveTransaction(
				ctx,
				tcase.Request.ReceiverID,
				tcase.Request.SenderID,
				tcase.Request.Amount,
				tcase.Request.Operation)

			require.EqualValues(t, tcase.ExpectedErr, err)
			fmt.Println(tx)
		})
	}
}

func TestGetTransactions(t *testing.T) {
	ctx := context.Background()

	repo, err := repository.NewConnection(ctx, PostgresDSN)
	if err != nil {
		log.Print(err)
		t.Fail()
	}

	cases := []struct {
		Name        string
		ExpectedErr error
		UserID      uuid.UUID
	}{
		{
			Name:        "ok case",
			ExpectedErr: nil,
			UserID:      uuid.MustParse(UserIDs[0]),
		},
	}

	for _, tcase := range cases {
		t.Run(tcase.Name, func(t *testing.T) {
			list, err := repo.GetTransactions(ctx, tcase.UserID)

			fmt.Println(list)

			require.EqualValues(t, tcase.ExpectedErr, err)
		})
	}
}
