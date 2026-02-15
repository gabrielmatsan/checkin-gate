include .env
export

MIGRATIONS_PATH = ./migrations
.PHONY: setup run dev swagger migrate-up migrate-down migrate-force migrate-create migrate-version

setup:
	go mod download
	go install github.com/air-verse/air@latest
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	go install github.com/swaggo/swag/cmd/swag@latest

swagger:
	$(shell go env GOPATH)/bin/swag init -g cmd/api/main.go -o docs

run:
	go run cmd/api/main.go

dev:
	make swagger
	$(shell go env GOPATH)/bin/air

migrate-up:
	migrate -path $(MIGRATIONS_PATH) -database "$(DATABASE_URL)" up

migrate-down:
	migrate -path $(MIGRATIONS_PATH) -database "$(DATABASE_URL)" down 1

migrate-down-all:
	migrate -path $(MIGRATIONS_PATH) -database "$(DATABASE_URL)" down -all

migrate-force:
	migrate -path $(MIGRATIONS_PATH) -database "$(DATABASE_URL)" force $(version)

migrate-version:
	migrate -path $(MIGRATIONS_PATH) -database "$(DATABASE_URL)" version

migrate-create:
	migrate create -ext sql -dir $(MIGRATIONS_PATH) -seq $(name)
