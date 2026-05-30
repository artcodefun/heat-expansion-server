# Heat Expansion — Billing Service

This directory contains the **billing** service inside the Heat Expansion modular monolith. It handles crystal package purchases and payment processing via YooKassa.

## Domain Overview

- **Crystal Packages**: Read-only value objects loaded from DB representing purchasable tiers (name, crystal count, price in minor currency units, currency code). Managed externally via direct SQL or a future admin panel.
- **Purchase Orders**: The core aggregate `PurchaseOrder` tracks the full lifecycle of a payment — `PENDING` on creation, `PAID` on successful webhook confirmation, `FAILED` on rejection. Each order is tied to a user ID (from JWT), a package, and a YooKassa payment.
- **YooKassa Integration**: Payments are created via the YooKassa REST API using the order ID as the idempotency key. Confirmation arrives via webhook, which transitions the order to `PAID` or `FAILED`.
- **Events**: On `PAID`, emits a `CrystalsPurchasedV1` integration event that the Game service consumes to credit crystals to the player's balance.

## Architecture

This service uses Hexagonal Architecture (Ports and Adapters), DDD, and CQRS — the same pattern as the Auth and Game services.

### Key Layers

- **Domain**: `internal/billing/domain`
  - `PurchaseOrder` aggregate, `CrystalPackage` value object, domain events (`OrderPaidEvent`, `OrderFailedEvent`).
- **Application**: `internal/billing/application`
  - `commands/`: `CreateOrderCommand`, `ConfirmPaymentCommand` (idempotent — no-op if already `PAID`).
  - `queries/`: `ListPackagesQuery`, `GetOrderQuery` (ownership-checked).
  - `cqrs/`: CQRS contract definitions, read models, error kinds.
  - `ports/`: Interfaces for repositories, payment gateway, outbox, transaction manager, token validator.
  - `services/`: Outbox loop, integration outbox loop, integration producer.
- **Infrastructure**: `internal/billing/infrastructure`
  - `db/`: Persistence using sqlc (`migrations/`, `queries/`, `gen/`, `repo/`, `dtos/`, `mappers/`).
  - `payment/`: YooKassa gateway adapter.
  - `events/`: In-process publisher and RabbitMQ publisher.
  - `security/`: HS256 JWT validator.
  - `i18n/`: JSON translator and embedded locale files.
- **Interfaces**: `internal/billing/interfaces/http`
  - Primary adapters (HTTP handlers, DTOs, router, auth middleware).
- **Bootstrap / Wiring**: `internal/billing/bootstrap`
  - Dependency injection and wiring of infrastructure adapters to application ports.

## HTTP API

Full OpenAPI spec: [`contracts/billing/http/v1/openapi.yaml`](../../contracts/billing/http/v1/openapi.yaml)

## Development

From repo root:

- Build: `make build`
- Run: `make run`
- Migrations: `make migrate-up` (billing migrations run as part of the shared target)
- SQLC: `make sqlc` (regenerates `internal/billing/infrastructure/db/gen/`)
- Create migration: `make billing-migrate-create name=<migration_name>`

## Database Schema

All tables live in the `billing` schema: `billing.crystal_packages`, `billing.purchase_orders`, `billing.domain_events`, `billing.integration_events`.
