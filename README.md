# Heat Expansion Server

Heat Expansion is a Go backend for a multiplayer 4X strategy game.

This repository is structured as a **modular monolith**: multiple services live under `internal/`, with low coupling so that services can be extracted into separate deployables later.

## Services

![Bounded Contexts](.github/BoundedContexts.png)

- **Game**: core gameplay domain, CQRS, HTTP API, persistence.
  - Docs: [internal/game/README.md](internal/game/README.md)

- **Auth**: identity and access management, JWT issuance, integration events.
  - Docs: [internal/auth/README.md](internal/auth/README.md)

- **Billing**: billing/subscription-related code (in progress).
  - Location: `internal/billing`

## API Contracts

- Auth HTTP contract: `contracts/auth/http/v1/openapi.yaml`
- Game HTTP contract: `contracts/game/http/v1/openapi.yaml`
- Auth integration events: `contracts/auth/events/`
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

The Game Server listens on `GAME_PORT` (default `8080`) and the Auth Server on `AUTH_PORT` (default `8081`).

## Observability

The server exports traces, metrics, and logs via **OTLP/gRPC** to an OpenTelemetry Collector, which forwards them to Grafana Cloud.

Set `OTEL_EXPORTER_OTLP_ENDPOINT` to your collector's gRPC address (e.g. `localhost:4317`).  
  Leave it **empty** to run in no-op mode — telemetry is fully disabled, which is the default for local dev.

## Internationalization (i18n)

The server supports multi-language responses based on the `Accept-Language` HTTP header. 

- **Systemic Locales**: Embedded in the binary for stability (errors, system alerts).
- **Content Locales**: Loaded from an external directory at runtime. Point the `GAME_I18N_PATH` environment variable to your translation files.