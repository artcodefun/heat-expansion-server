# Makefile for Heat Expansion API
BINARY_NAME=heat-expansion-api
GOARCH?=amd64
GOOS?=linux
CGO_ENABLED?=0

.PHONY: build run test sqlc clean

build:
	go build -o bin/${BINARY_NAME} ./cmd/api

run:
	set -a; source .env; set +a; go run ./cmd/api

test:
	go test ./...

clean:
	rm -rf bin/

# Check if a .env file exists.
ifneq (,$(wildcard ./.env))
    # Include the .env file as make variables
    include .env
    # Export all variables from the .env file as environment variables
    export
endif

MIGRATION_DIR=internal/infrastructure/db/migrations
DB_URL?=postgres://user:password@localhost:5432/heatdb?sslmode=disable

migrate-up:
	migrate -path $(MIGRATION_DIR) -database "$(DB_URL)" up

migrate-down:
	migrate -path $(MIGRATION_DIR) -database "$(DB_URL)" down

migrate-create:
	migrate create -ext sql -dir $(MIGRATION_DIR) -seq $(name)

sqlc:
	sqlc generate