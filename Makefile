POSTGRES_DSN := "postgres://postgres:postgres@localhost:5432/finapi?sslmode=disable"

.PHONY: run
run:
	go run ./cmd/server/main.go

DEFAULT_GOAL: run

goose-create:
	go run github.com/pressly/goose/v3/cmd/goose@latest \
	-dir ./internal/repository/migrations create test_users_add postgres 

goose-up:
	go run github.com/pressly/goose/v3/cmd/goose@latest \
	-dir ./internal/repository/migrations postgres $(POSTGRES_DSN) up

goose-down:
	go run github.com/pressly/goose/v3/cmd/goose@latest \
	-dir ./internal/repository/migrations postgres $(POSTGRES_DSN) down
