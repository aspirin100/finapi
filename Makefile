POSTGRES_DSN := "postgres://postgres:postgres@localhost:5432/finapi?sslmode=disable"

.PHONY: run
run:
	make run ./cmd/server/main.go

DEFAULT_GOAL: run

goose-create:
	go run github.com/pressly/goose/v3/cmd/goose@latest \
	-dir ./internal/repository/migrations create init_users postgres 

goose-up:
	go run github.com/pressly/goose/v3/cmd/goose@latest \
	-dir ./internal/repository/migrations postgres $(POSTGRES_DSN) up

