# Copilot Instructions for Heat Expansion Server

This repo is a Go backend for the Heat Expansion strategy game. It uses Hexagonal Architecture, DDD, and CQRS. Follow these project-specific guidelines when editing or generating code.

## Architecture & Layout
- This repo is a **modular monolith**: each service lives under `internal/<service>`.
- **Game service (primary)**: `internal/game`
  - `domain/`: Aggregates and domain logic (e.g. `UserBase`, `MilitaryOperation`, events, value objects).
  - `application/`: CQRS + ports + application services
    - `commands/`: Write-side command handlers.
    - `queries/`: Read-side query handlers.
    - `cqrs/`: CQRS contract definitions and readmodels.
    - `ports/`: Interfaces for repositories, schedulers, token providers, event publishing, transactions, etc.
    - `services/`: App-level services (access control, provisioning, outbox loop, ...).
  - `infrastructure/`: Secondary adapters (DB/sqlc, readstore, jobs, events, security, content, ...).
  - `interfaces/http`: Primary adapter (HTTP handlers/DTOs/middleware/router).
  - `bootstrap/`: Dependency wiring for the game service.
- **Auth service**: `internal/auth` (IAM, JWT, Integration producers)
- **Shared contracts**: `contracts/` (Versioned integration event schemas)
- **Other services (WIP)**: `internal/billing`.

## Key Patterns & Conventions
The patterns and conventions below apply to the **Game** service (`internal/game`) unless stated otherwise.

- **Commands vs. queries**
  - Commands mutate state and live in `internal/game/application/commands`. They run inside `TransactionManager.WithTx`, use repositories from `application/ports`, and interact with aggregates in `domain/`.
  - Queries do not mutate state and live in `internal/game/application/queries`. They depend on read repositories from `internal/game/infrastructure/readstore` and use shared services like access control from `application/services`.
- **Repositories & transactions**
  - Repositories are declared as ports (e.g. `UserBaseRepository`, `SectorRepository`) and implemented in `internal/game/infrastructure/db/repo` using sqlc-generated `gen` packages.
  - Use `Tx(tx)` on repositories and outbox interfaces when working inside a transaction; do not create new DB connections directly in core or handlers.
- **Domain events & outbox**
  - Aggregates emit domain events via `EventProducer` in `internal/game/domain`.
  - Command handlers do **not** publish directly; they call `OutboxEventRepository.Save(events)` inside the transaction.
  - `OutboxService` (in `internal/game/application/services/outbox_service.go`) runs in a background loop from `App.Run` and pulls from the `domain_events` table to publish via `EventPublisher`.
  - When adding new events, update outbox DTOs/mappers in `internal/game/infrastructure/db/dtos` and `mappers` rather than encoding from domain types directly in handlers.
 - **Dependencies are never optional**
   - Command handlers, query handlers, and services treat their constructor dependencies as required. Do **not** add `nil` checks (e.g. `if c.Outbox != nil`) around injected ports/services.
   - If something can be `nil`, fix the wiring in `internal/bootstrap` (adapters/services/commands/queries) instead of guarding at the use site.
- **Access control & provisioning**
  - Authorization checks should use `AccessControlService` (from `internal/game/application/services/access_control_service.go`) provided via the game bootstrap wiring.
  - Lazy creation of sectors/bases/locations uses `SectorProvisioningService` (also in `application/services`).
- **HTTP layer**
  - Handlers in `internal/game/interfaces/http/handlers` should call into `bootstrap.Commands`/`bootstrap.Queries`, not directly into repositories or DB.
  - Map domain/CQRS errors to HTTP status codes consistently using existing helpers in `http/dtos` and middleware.

- **Integration Events**
  - External events live in `contracts/`. They consist of a generic `IntegrationEvent` envelope and versioned payloads (e.g., `contracts/auth/v1/AccountRegisteredV1`).
  - Use `RegisterPayload` in `init()` functions to enable polymorphic unmarshaling.
  - Integration events are produced from domain events via an `IntegrationProducer` and stored in an `IntegrationOutbox`.
  - Idempotency is enforced on `(origin_id, event_type)` in the integration outbox table.
  - Publishing is performed via RabbitMQ (using a `topic` exchange and event type as routing key) with a console fallback for local dev.

- **Internationalization (i18n)**
  - All user-facing strings (errors, notifications, prototype names) must be translatable.
  - **Systemic locales**: Errors, alerts, and world generation text are stored in `internal/game/infrastructure/i18n/locales/*.json` and embedded into the binary via `go:embed`.
  - **Content locales**: Prototype-specific data (army names, descriptions) is generated into an external directory and loaded at runtime via the `GAME_I18N_PATH` environment variable.
  - **Domain Errors**: Use `domain.NewError(key, params)` in domain logic. Never use `fmt.Errorf` with hardcoded English strings.
  - **Application Errors**: Use `cqrs.NewAppError(kind, key)` or `cqrs.NewAppErrorWithParams` for high-level application failures.
  - **Presentation Layer**: DTO mappers and HTTP handlers must use `ports.Translator` to resolve keys into final strings using the `locale` from the `Accept-Language` header.

## Workflows
- **Build**: `make build` builds `./cmd/server` into `bin/heat-expansion-server`.
- **Run locally**: `make run` (loads `.env`, then `go run ./cmd/server`). Ensure DB is running and `GAME_DB_URL`/`AUTH_DB_URL` are set.
- **Tests**: `make test` or `go test ./...` from repo root.
- **Migrations**: use `make migrate-up` / `make migrate-down` (requires `migrate` CLI). Game SQL files live in `internal/game/infrastructure/db/migrations`.
- **Generated files**: never edit generated files directly. This includes `internal/game/infrastructure/readstore/gen/**` and any other generator output.
- **SQLC queries**: never write ad hoc SQL strings inside Go files for database access. If a change needs a new or updated query, edit only `internal/game/infrastructure/db/queries/*.sql` or `internal/game/infrastructure/readstore/queries/*.sql`, then run `make sqlc`.

## How to Extend Safely
- When introducing new domain behavior, prefer adding methods to aggregates in `internal/game/domain` and invoking them from command handlers, rather than mutating models directly in handlers.
- For new write-side features, add/extend ports in `internal/game/application/ports`, implement them in `internal/game/infrastructure/db/repo`, wire them in `internal/game/bootstrap/adapters.go`, and inject via `Commands`/services.
- For new read-side endpoints, add read models and queries in `internal/game/infrastructure/readstore`, then wire new query facades in `internal/game/application/queries` and expose via HTTP handlers.
- For new integration events:
  1. Define the payload in `contracts/`.
  2. Implement `IntegrationEventType() string` and register the factory.
  3. Create/update an `IntegrationProducer` to map domain events to the new contract.
- Keep serialization, DTOs, and DB schemas in infra layers (`dtos/`, `mappers/`, `queries/`), not in domain or application packages.
