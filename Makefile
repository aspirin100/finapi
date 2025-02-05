POSTGRES_DSN := "postgres://postgres:postgres@localhost:5432/finapi?sslmode=disable"

build:
	mkdir -p bin && \
	go build -o ./bin/finapi-server ./cmd/finapi/main.go

goose-create:
	go run github.com/pressly/goose/v3/cmd/goose@latest \
	-dir ./internal/repository/migrations create test_users_add postgres 

goose-up:
	go run github.com/pressly/goose/v3/cmd/goose@latest \
	-dir ./internal/repository/migrations postgres $(POSTGRES_DSN) up

goose-down:
	go run github.com/pressly/goose/v3/cmd/goose@latest \
	-dir ./internal/repository/migrations postgres $(POSTGRES_DSN) down

lint:
	go run github.com/golangci/golangci-lint/cmd/golangci-lint@latest run

cover:
		go test -race -coverprofile=coverage.out ./... && \
		go tool cover -html=coverage.out && \
		rm coverage.out

swagger:
	docker run -d -p 8090:8080 -e SWAGGER_JSON=/openapi/openapi_v1.yml -v $(CURDIR)/docs:/openapi swaggerapi/swagger-ui

migrations-up:
	go run ./cmd/migrator/main.go \
	--migrations-path ./internal/repository/migrations \
	--dsn postgres://postgres:postgres@:5432/finapi?sslmode=disable

postgres-run:
	docker run -d \
	-e POSTGRES_USER="postgres" \
    -e POSTGRES_PASSWORD="postgres" \
    -e POSTGRES_DB="finapi" \
	-p 5432:5432 \
	--network finapi-local-net \
	--name postgres \
	postgres:latest  \

server-run:
	docker run -d \
	-e FINAPI_HOSTNAME=0.0.0.0 \
    -e FINAPI_PORT=8080 \
	-e FINAPI_POSTGRES_DSN=postgres://postgres:postgres@postgres:5432/finapi?sslmode=disable \
    -e FINAPI_DB_TIMEOUT=5s \
	-p 8080:8080 \
	--network finapi-local-net \
	finapi-img

.PHONY: run
run:
	docker network create finapi-local-net && \
	make postgres-run && \
	make migrations-up && \
	docker build -t finapi-img . && \
	make server-run \
	
