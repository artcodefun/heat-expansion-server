# Makefile for Heat Expansion Server
BINARY_NAME=heat-expansion-server
GOARCH?=amd64
GOOS?=linux
CGO_ENABLED?=0

.PHONY: build run test sqlc clean

build:
	go build -o bin/${BINARY_NAME} ./cmd/server

run:
	set -a; source .env; set +a; go run ./cmd/server

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

GAME_MIGRATION_DIR=internal/game/infrastructure/db/migrations
GAME_DB_URL?=postgres://user:password@localhost:5432/heatdb?sslmode=disable
GAME_MIGRATIONS_TABLE?=game_schema_migrations

AUTH_MIGRATION_DIR=internal/auth/infrastructure/db/migrations
AUTH_DB_URL?=postgres://user:password@localhost:5432/heatdb?sslmode=disable
AUTH_MIGRATIONS_TABLE?=auth_schema_migrations

GAME_MIGRATE_DB_URL=$(GAME_DB_URL)$(if $(findstring ?,$(GAME_DB_URL)),&,?)x-migrations-table=$(GAME_MIGRATIONS_TABLE)
AUTH_MIGRATE_DB_URL=$(AUTH_DB_URL)$(if $(findstring ?,$(AUTH_DB_URL)),&,?)x-migrations-table=$(AUTH_MIGRATIONS_TABLE)

migrate-up:
	migrate -path $(GAME_MIGRATION_DIR) -database "$(GAME_MIGRATE_DB_URL)" up
	migrate -path $(AUTH_MIGRATION_DIR) -database "$(AUTH_MIGRATE_DB_URL)" up

migrate-down:
	migrate -path $(AUTH_MIGRATION_DIR) -database "$(AUTH_MIGRATE_DB_URL)" down
	migrate -path $(GAME_MIGRATION_DIR) -database "$(GAME_MIGRATE_DB_URL)" down

migrate-drop:
	migrate -path $(AUTH_MIGRATION_DIR) -database "$(AUTH_MIGRATE_DB_URL)" drop
	migrate -path $(GAME_MIGRATION_DIR) -database "$(GAME_MIGRATE_DB_URL)" drop

game-migrate-create:
	migrate create -ext sql -dir $(GAME_MIGRATION_DIR) -seq $(name)

auth-migrate-create:
	migrate create -ext sql -dir $(AUTH_MIGRATION_DIR) -seq $(name)

sqlc:
	sqlc -f internal/game/infrastructure/sqlc.yaml generate
	sqlc -f internal/auth/infrastructure/sqlc.yaml generate