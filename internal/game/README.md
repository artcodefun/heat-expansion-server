# Heat Expansion — Game Service

This directory contains the **game** service inside the Heat Expansion modular monolith.

It follows Hexagonal Architecture (Ports & Adapters), DDD, and CQRS:
- **Domain** contains business rules and invariants.
- **Application** contains use-cases (commands/queries) and application services.
- **Infrastructure** contains secondary adapters (DB/sqlc, readstore, jobs, security, etc.).
- **Interfaces** contains primary adapters (HTTP).

## Layout

- **Domain**: `internal/game/domain`
  - Aggregates/entities/value objects (e.g. base, sector, operations)
  - Domain services and domain events

- **Application**: `internal/game/application`
  - `commands/`: write-side command handlers (mutations)
  - `queries/`: read-side query handlers
  - `cqrs/`: CQRS primitives + readmodels
  - `ports/`: required interfaces (repositories, tx manager, scheduler, token provider, etc.)
  - `services/`: app-level services (access control, provisioning, outbox loop, …)

- **Infrastructure**: `internal/game/infrastructure`
  - `db/`: write-side persistence using sqlc (`migrations/`, `queries/`, `gen/`, `repo/`, `dtos/`, `mappers/`)
  - `readstore/`: read-side persistence using sqlc (`queries/`, `gen/`, `repo/`, `mappers/`)
  - `events/`, `jobs/`, `security/`, `content/`: secondary adapters

- **Interfaces (HTTP)**: `internal/game/interfaces/http`
  - handlers, DTOs, middleware, router/server wiring

- **Bootstrap / Wiring**: `internal/game/bootstrap`
  - wires adapters ↔ ports and exposes aggregated `Commands`/`Queries`

## Development

From repo root:

- Build: `make build`
- Run: `make run`
- Tests: `make test`
- SQLC: `make sqlc`
- Migrations: `make migrate-up` / `make migrate-down`

## Database schema

All tables live in the `game` schema (e.g. `game.users`, `game.user_bases`, …).
The sqlc config uses `rename:` mappings so generated Go structs have clean names (e.g. `User`, not `GameUser`).
