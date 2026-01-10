# Copilot Instructions for Heat Expansion API

This repo is a Go backend for the Heat Expansion strategy game. It uses Hexagonal Architecture, DDD, and CQRS. Follow these project-specific guidelines when editing or generating code.

## Architecture & Layout
- **Core domain & CQRS**: `internal/core`
  - `domain/`: Aggregates and domain logic (e.g. `UserBaseModel`, `MilitaryOperation`, events, value objects).
  - `ports/`: Interfaces for repositories, schedulers, token providers, event publishing, transactions, etc.
  - `commands/`: Write-side command handlers. They load aggregates via ports, apply domain methods, and persist changes.
  - `queries/`: Read-side query handlers working on read models from `internal/infrastructure/readstore`.
  - `cqrs/`: Common CQRS primitives and readmodels.
- **Infrastructure adapters**: `internal/infrastructure`
  - `db/`: Write-side persistence using `sqlc` (`queries/`, `gen/`, `repo/`, `dtos/`, `mappers/`).
  - `readstore/`: Read-side persistence and mappers to CQRS readmodels.
  - `events/`: Event publisher implementation.
  - `jobs/`: Scheduler implementation for delayed jobs.
  - `content/`: Content generator used by provisioning services.
  - `security/`: Hashing and token implementations.
- **Interfaces (HTTP)**: `internal/interfaces/http`
  - Handlers translate HTTP to CQRS commands/queries; DTOs live in `dtos/`; router/server wiring lives here.
- **Bootstrap**: `internal/bootstrap`
  - `adapters.go`: Wires infrastructure to ports (repositories, scheduler, tx manager, tokens, etc.).
  - `app_services.go`: Aggregates app-level services (e.g. `AppServices` with provisioning, access control, outbox).
  - `commands.go` / `queries.go`: Aggregated command/query structs created from `Adapters` and `AppServices`.
  - `app.go`: Builds and runs the `App` (DB, adapters, services, commands, queries, HTTP server, background loops).
  - Entry point: `cmd/api/main.go` just builds `App` and calls `Run`.

## Key Patterns & Conventions
- **Commands vs. queries**
  - Commands mutate state and live in `internal/core/commands`. They run inside `TransactionManager.WithTx`, use repositories from `ports`, and interact with aggregates in `domain/`.
  - Queries do not mutate state and live in `internal/core/queries`. They depend on read repositories from `internal/infrastructure/readstore` and use shared services like access control from `AppServices`.
- **Repositories & transactions**
  - Repositories are declared as ports (e.g. `UserBaseRepository`, `SectorRepository`) and implemented in `internal/infrastructure/db/repo` using sqlc-generated `gen` packages.
  - Use `Tx(tx)` on repositories and outbox interfaces when working inside a transaction; do not create new DB connections directly in core or handlers.
- **Domain events & outbox**
  - Aggregates emit domain events via `EventProducer` in `internal/core/domain`.
  - Command handlers do **not** publish directly; they call `OutboxEventRepository.Save(events)` inside the transaction.
  - `OutboxService` (in `internal/core/services/outbox_service.go`) runs in a background loop from `App.Run` and pulls from the `domain_events` table to publish via `EventPublisher`.
  - When adding new events, update outbox DTOs/mappers in `internal/infrastructure/db/dtos` and `mappers` rather than encoding from domain types directly in handlers.
 - **Dependencies are never optional**
   - Command handlers, query handlers, and services treat their constructor dependencies as required. Do **not** add `nil` checks (e.g. `if c.Outbox != nil`) around injected ports/services.
   - If something can be `nil`, fix the wiring in `internal/bootstrap` (adapters/services/commands/queries) instead of guarding at the use site.
- **Access control & provisioning**
  - Authorization checks should use `AccessControlService` (from `internal/core/services/access_control_service.go`) provided via `AppServices.Access`, not ad-hoc user/base checks scattered through commands/queries.
  - Lazy creation of sectors/bases/locations uses `SectorProvisioningService` (also in `services`) via `AppServices.Provisioner`.
- **HTTP layer**
  - Handlers in `internal/interfaces/http/handlers` should call into `bootstrap.Commands`/`bootstrap.Queries`, not directly into repositories or DB.
  - Map domain/CQRS errors to HTTP status codes consistently using existing helpers in `http/dtos` and middleware.

## Workflows
- **Build**: `make build` builds `./cmd/api` into `bin/heat-expansion-api`.
- **Run locally**: `make run` (loads `.env`, then `go run ./cmd/api`). Ensure DB is running and `DB_URL` is set.
- **Tests**: `make test` or `go test ./...` from repo root.
- **Migrations**: use `make migrate-up` / `make migrate-down` (requires `migrate` CLI). SQL files live in `internal/infrastructure/db/migrations`.
- **sqlc**: when changing DB queries, edit `internal/infrastructure/db/queries/*.sql` or readstore queries, then run `make sqlc`.

## How to Extend Safely
- When introducing new domain behavior, prefer adding methods to aggregates in `internal/core/domain` and invoking them from command handlers, rather than mutating models directly in handlers.
- For new write-side features, add/extend ports in `internal/core/ports`, implement them in `internal/infrastructure/db/repo`, wire them in `internal/bootstrap/adapters.go`, and inject via `Commands`/`AppServices`.
- For new read-side endpoints, add read models and queries in `internal/infrastructure/readstore`, then wire new query facades in `internal/core/queries` and expose via HTTP handlers.
- Keep serialization, DTOs, and DB schemas in infra layers (`dtos/`, `mappers/`, `queries/`), not in domain or command/query packages.
