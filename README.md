# Heat Expansion Server

Heat Expansion Server is a Go backend for a multiplayer 4X strategy game.

This repository is structured as a **modular monolith**: multiple services live under `internal/`, with low coupling so that services can be extracted into separate deployables later.

## Services

![Bounded Contexts](.github/BoundedContexts.png)

- **Game**: core gameplay domain, CQRS, HTTP API, persistence.
  - Docs: [internal/game/README.md](internal/game/README.md)

- **Auth**: identity and access management, JWT issuance, integration events.
  - Docs: [internal/auth/README.md](internal/auth/README.md)

- **Billing**: crystal package purchases and YooKassa payment processing.
  - Docs: [internal/billing/README.md](internal/billing/README.md)

- **Admin**: back-office API for game content and billing configuration. Self-contained session-based auth (no JWT).
  - Docs: [internal/admin/README.md](internal/admin/README.md)

- **Platform** (`internal/platform/`): shared infrastructure adapters used across services.
  - Docs: [internal/platform/README.md](internal/platform/README.md)

## API Contracts

- Auth HTTP contract: `contracts/auth/http/v1/openapi.yaml`
- Game HTTP contract: `contracts/game/http/v1/openapi.yaml`
- Billing HTTP contract: `contracts/billing/http/v1/openapi.yaml`
- Admin HTTP contract: `contracts/admin/http/v1/openapi.yaml`
- Game private gRPC contract: `contracts/game/grpc/v1/admin.proto`
- Billing private gRPC contract: `contracts/billing/grpc/v1/admin.proto`
- Integration event envelope: `contracts/events/envelope.go`
- Auth event payloads: `contracts/auth/events/v1/`
- Billing event payloads: `contracts/billing/events/v1/`
- Swagger UI is served by each service at `/api/v1/docs`, backed by the versioned OpenAPI document at `/api/v1/openapi.yaml`.

## Getting started

1. Install Go, PostgreSQL, and RabbitMQ.
2. Create a `.env` file (see `.env.example`).
3. Apply migrations and run the Server:
   - `make migrate-up`
   - `make run`

Alternatively, use Docker Compose:
```bash
docker-compose up --build
```

The Game Server listens on `GAME_PORT` (default `8080`), the Auth Server on `AUTH_PORT` (default `8081`), the Billing Server on `BILLING_PORT` (default `8082`), and the Admin Server on `ADMIN_PORT` (default `8083`).

## Observability

The server exports traces, metrics, and logs via **OTLP/gRPC** to an OpenTelemetry Collector, which forwards them to Grafana Cloud.

Set `OTEL_EXPORTER_OTLP_ENDPOINT` to your collector's gRPC address (e.g. `localhost:4317`).  
  Leave it **empty** to run in no-op mode — telemetry is fully disabled, which is the default for local dev.

## Authentication

Tokens are signed by the **Auth service** using **ES256** (ECDSA P-256) and verified by the Game and Billing services using only the corresponding public key. This means only Auth can issue tokens — other services can verify them but cannot forge new ones.

Tokens carry a `sub` claim (account UUID) and expire after 72 hours.

**Key pair setup** — generate once and distribute via environment variables:

```bash
openssl ecparam -name prime256v1 -genkey -noout -out ec.key
openssl ec -in ec.key -pubout -out ec.pub
```

PEM files are multi-line, but `.env` files require single-line values. Collapse each key into one line with literal `\n` between parts:

```bash
awk 'NF {printf "%s\\n",$0}' ec.key   # paste into AUTH_JWT_PRIVATE_KEY=
awk 'NF {printf "%s\\n",$0}' ec.pub   # paste into AUTH_JWT_PUBLIC_KEY=
```

`AUTH_JWT_PRIVATE_KEY` is used only by the Auth service. `AUTH_JWT_PUBLIC_KEY` is shared with Game and Billing — they receive no access to the private key.

## Internationalization (i18n)

The server supports multi-language responses based on the `Accept-Language` HTTP header. 

- **Systemic Locales**: Embedded in the binary for stability (errors, system alerts).
- **Content Locales**: Service-specific content (e.g. prototype names, descriptions) is stored in the database and loaded at startup.

## License

This project is licensed under the GNU Affero General Public License v3.0.

See the LICENSE file for details.