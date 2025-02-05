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
	docker run -d -p 8090:8080 -e SWAGGER_JSON=/openapi/openapi_v1.yml -v $(CURDIR):/openapi swaggerapi/swagger-ui

.PHONY: run
run:
	docker build -t finapi-img . && \
	docker-compose up -d && \
	make goose-up
