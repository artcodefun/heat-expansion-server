# Heat Expansion — Admin Service

This directory contains the **administration** service inside the Heat Expansion modular monolith. It provides a back-office API for operators to manage game content and billing configuration across the other modules.

## Domain Overview

At a high level, the Admin service models operator identity and access. It is **self-contained**: it does not depend on the auth service or JWTs.

- **Admin**: The core aggregate `Admin` represents an operator identity. Admins are seeded with a username and a one-time `invite_token`. Registration is completed by calling `Register`, which sets the `password_hash` and clears the invite token so it cannot be reused.
- **Session**: A server-side session record (`admin.sessions`) tied to an authenticated admin. Sessions use opaque, cryptographically-random bearer tokens (no JWT). Validation is a token lookup plus expiry check.
- **Authentication**: Username/password verified with bcrypt; a session token is issued on successful registration or login and revoked on logout.

## Architecture

This service uses Hexagonal Architecture (Ports and Adapters), DDD (Domain-driven design) and the CQRS (Command Query Responsibility Segregation) pattern.

### Key Layers

- **Domain**: `internal/admin/domain`
  - Business rules, the `Admin` aggregate, the `Session` record, and translatable domain errors.
- **Application**: `internal/admin/application`
  - `commands/`: Write-side command handlers for registration, login, logout, and CRUD operations for prototypes, translations, and crystal packages.
  - `queries/`: Read-side query handlers (admin profile, prototypes, translations, crystal packages).
  - `cqrs/`: CQRS contract definitions and read models.
  - `ports/`: Interfaces for repositories, read repositories, password hasher, session token generator, session validator, transaction manager, translator, and outbound gRPC clients (`GamePrivateClient`, `BillingPrivateClient`).
- **Infrastructure**: `internal/admin/infrastructure`
  - `db/`: Write-side persistence using sqlc (`migrations/`, `queries/`, `gen/`, `repo/`).
  - `readstore/`: Read-side projections using a separate sqlc generation (`queries/`, `gen/`, `repo/`, `mappers/`).
  - `grpcclient/`: Outbound gRPC clients for game and billing private APIs (lazy-dial with static key authentication).
  - `security/`: Session token generation (bcrypt password hashing is shared via `internal/platform/security`) & validation.
  - `i18n/`: Embedded translator for systemic admin strings (errors).
- **Interfaces**: `internal/admin/interfaces/http`
  - Primary adapters (HTTP handlers, DTOs, router, and the bearer-token auth middleware).
- **Bootstrap / Wiring**: `internal/admin/bootstrap`
  - Dependency injection and wiring of concrete infrastructure adapters to application ports.

## HTTP API

Full OpenAPI spec: [`contracts/admin/http/v1/openapi.yaml`](../../contracts/admin/http/v1/openapi.yaml)

- **Identity** (`/api/v1/auth`): `POST /register`, `POST /login`, `POST /logout`, `GET /me`
- **Prototypes** (`/api/v1/game/prototypes/{kind}`): CRUD for army, build, storage, and tech prototypes
- **Translations** (`/api/v1/game/translations`): list and upsert translation entries
- **Crystal Packages** (`/api/v1/billing/packages`): CRUD for purchasable crystal packages

The spec is served at runtime at `/api/v1/openapi.yaml` with Swagger UI at `/api/v1/docs/`.

## Development

From repo root:

- Build: `make build`
- Run: `make run`
- Migrations: `make migrate-up` (admin migrations run as part of the shared target)
- SQLC: `make sqlc` (regenerates `internal/admin/infrastructure/db/gen/` and `readstore/gen/`)
- Create migration: `make admin-migrate-create name=<migration_name>`

## Database Schema

All tables live in the `admin` schema: `admin.admins`, `admin.sessions`.
