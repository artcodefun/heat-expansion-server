# Heat Expansion API

Heat Expansion is a Go backend for a multiplayer 4X strategy game.

This repository is structured as a **modular monolith**: multiple services live under `internal/`, with low coupling so that services can be extracted into separate deployables later.

## Services

- **Game**: core gameplay domain, CQRS, HTTP API, persistence
  - Docs: [internal/game/README.md](internal/game/README.md)

- **Auth**: authentication-related code (in progress)
  - Location: `internal/auth`

- **Billing**: billing/subscription-related code (in progress)
  - Location: `internal/billing`

## Getting started

1. Install Go and PostgreSQL.
2. Create a `.env` file (see `.env.example`).
3. Apply migrations and run the API:
   - `make migrate-up`
   - `make run`

The Game API listens on `GAME_PORT` (default `8080`).
