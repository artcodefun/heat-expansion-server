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
- **Billing service**: `internal/billing` — crystal-package purchases via YooKassa. The webhook handler never trusts the request body; it re-queries YooKassa for canonical payment state. On success it emits `CrystalsPurchasedV1`, which the game service consumes to credit crystals (idempotent on `order_id`).
- **Shared contracts**: `contracts/` (Versioned integration event schemas and HTTP OpenAPI contracts)
- **Shared platform**: `internal/platform/` — infrastructure adapters reused across services (RabbitMQ publisher/consumer, JWT token validator, in-process event publisher, i18n translator core). When an adapter is needed by more than one service, it lives here rather than being duplicated.

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
- **Scheduled jobs**
  - Jobs are created via `Scheduler.Schedule()` inside event handlers, which run outside the outbox transaction. The system provides **at-least-once** delivery — duplicate job executions are possible under crash/retry scenarios.
  - **Every job handler must be idempotent.** Use domain state guards so a duplicate invocation is a no-op.
  - Jobs are never created directly in command handlers or outside of event handlers.
  - Self-rescheduling jobs must use **positive-only jitter** when computing the next `executeAt`. Never subtract from the period — firing early can cause legitimate runs to be silently skipped by time-based idempotency checks.
- **Module lifecycle: constructors may compute, only Run may connect**
  - `bootstrap.NewModule()` is pure wiring: env validation, PEM parsing, struct construction. It must not start goroutines. `db.Open` and `db.Ping` are the one accepted exception — they are cheap, context-free, and a bad DB URL or unreachable host should abort the process immediately before anything else starts.
  - All other infrastructure I/O happens inside `Module.Run(ctx)` in a startup phase before the HTTP server starts serving: one-shot setup like `Adapters.Setup(ctx)` (e.g. game's translation load) and job seeding runs first, then broker connections are dialed by worker loops. This keeps startup cancelable by the signal context.
  - Adapters that own connections (e.g. `RabbitMQPublisher`, `RabbitMQConsumer`) expose a blocking `Start(ctx) error`: it dials, reconnects on drops, and releases the connection when ctx is cancelled. They run as worker loops and never connect in their constructors. A failed initial dial is returned as an error and fails the module (fail fast).
  - `Run(ctx) error` returns startup/serve failures instead of calling `log.Fatal`. `cmd/server` runs modules in an `errgroup`: one module's failure cancels the shared context so the others drain gracefully, and the process exits non-zero.
- **Dependencies are never optional**
   - Command handlers, query handlers, and services treat their constructor dependencies as required. Do **not** add `nil` checks (e.g. `if c.Outbox != nil`) around injected ports/services.
   - If something can be `nil`, fix the wiring in `internal/bootstrap` (adapters/services/commands/queries) instead of guarding at the use site.
- **Access control & provisioning**
  - Authorization checks should use `AccessControlService` (from `internal/game/application/services/access_control_service.go`) provided via the game bootstrap wiring.
  - Lazy creation of sectors/bases/locations uses `SectorProvisioningService` (also in `application/services`).
- **HTTP layer**
  - Handlers in `internal/game/interfaces/http/handlers` should call into `bootstrap.Commands`/`bootstrap.Queries`, not directly into repositories or DB.
  - Map domain/CQRS errors to HTTP status codes consistently using existing helpers in `http/dtos` and middleware.
  - HTTP contracts live under `contracts/<service>/http/vN/openapi.yaml`. Whenever an HTTP route, request DTO, response DTO, status code, auth requirement, or public API behavior changes, update the corresponding OpenAPI file in the same change.

- **Integration Events**
  - The shared `IntegrationEvent` envelope lives in `contracts/events/envelope.go` and uses `json.RawMessage` for its payload field. Versioned payload structs live in `contracts/<service>/events/v1/` alongside a typed `EventXxx` string constant (e.g., `EventAccountRegisteredV1 = "auth.account.registered.v1"`).
  - **Producers**: `json.Marshal` the payload struct, then call `events.NewIntegrationEvent(originID, occurredAt, v1.EventXxx, payload)` and pass the result to `outbox.Save`. No `init()` registration or interface implementation needed on the payload type.
  - **Consumers** (bootstrap wiring): `json.Unmarshal(d.Body, &envelope)` to get the envelope, then `switch envelope.Type` on the typed constant, then `json.Unmarshal(envelope.Payload, &typed)` for each known case. The default branch logs a warning and acks the message.
  - Integration events are produced from domain events via an `IntegrationProducer` and stored in an `IntegrationOutbox`.
  - Idempotency is enforced on `(origin_id, kind)` in the integration outbox table.
  - Publishing is performed via RabbitMQ (using a `topic` exchange and event type as routing key) with a console fallback for local dev.

- **Observability & logging**
  - Telemetry is initialised in `cmd/server/telemetry.go`. It wires up OTel TracerProvider, MeterProvider, and LoggerProvider, and bridges the global `slog` logger to OTel via `otelslog`. When `OTEL_EXPORTER_OTLP_ENDPOINT` is empty the whole setup is a no-op.
  - **Always use context-aware `slog` methods** (`slog.InfoContext`, `slog.WarnContext`, `slog.ErrorContext`, `slog.DebugContext`) rather than their context-free variants (`slog.Info`, `slog.Warn`, etc.). The `otelslog` bridge reads the active span from the context and injects `trace_id` and `span_id` into the log record, correlating logs with traces in Grafana. Dropping the context breaks that correlation silently.

- **Internationalization (i18n)**
  - All user-facing strings (errors, notifications, prototype names) must be translatable.
  - **Systemic locales**: Errors, alerts, and world generation text are stored in `internal/game/infrastructure/i18n/locales/*.json` and embedded into the binary via `go:embed`.
  - **Content locales**: Service-specific content (e.g. prototype names, descriptions) is stored in the database and loaded at startup.
  - **Key parity requirement**: Whenever a new translation key is introduced, add its translations immediately in both English and Russian locale files in the same change.
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
  1. Add a payload struct and `EventXxx` string constant in `contracts/<service>/events/v1/`.
  2. In the producer service, `json.Marshal` the payload and call `events.NewIntegrationEvent` with the typed constant, then save via `outbox.Save`.
  3. Add a `case v1.EventXxx:` branch to the relevant consumer's type-switch in `bootstrap/` that unmarshals `envelope.Payload` and dispatches to the handler.
- Keep serialization, DTOs, and DB schemas in infra layers (`dtos/`, `mappers/`, `queries/`), not in domain or application packages.
