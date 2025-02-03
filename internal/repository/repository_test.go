package repository

import (
	"context"
	"log"
	"testing"

	"github.com/google/uuid"
)

const PostgresDSN = "postgres://postgres:postgres@localhost:5432/finapi?sslmode=disable"

func TestUpdateBalance(t *testing.T) {
	ctx := context.Background()

	repo, err := NewConnection(ctx, PostgresDSN)
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
			UserID:      uuid.Nil,
		},
	}

}
