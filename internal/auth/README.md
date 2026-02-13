# Heat Expansion — Auth Service

This directory contains the **authentication** service inside the Heat Expansion modular monolith. It handles user identity, account management, and security.

## Domain Overview

At a high level, the Auth service models user access and security credentials.

- **Accounts**: The core aggregate `Account` represents a user's identity. It stores basic profile information (name, email) and a secure `password_hash`.
- **Registration**: Handles new account creation, ensuring email uniqueness and initializing security credentials.
- **Authentication**: Validates user credentials and issues secure JWT tokens for session management.
- **Events**: Emits `AccountRegisteredEvent` when a new account is created, which is then projected as an integration event for other services (like the Game service) to consume.

## Architecture

This service uses Hexagonal Architecture (Ports and Adapters), DDD (Domain-driven design) and the CQRS (Command Query Responsibility Segregation) pattern.

### Key Layers

- **Domain**: `internal/auth/domain`
  - Business rules, `Account` aggregate, and domain events (e.g., `AccountRegisteredEvent`).
- **Application**: `internal/auth/application`
  - `commands/`: Write-side command handlers for registration and login.
  - `cqrs/`: CQRS contract definitions.
  - `ports/`: Interfaces for repositories, token providers, and password hashers.
  - `services/`: App-level services like the outbox loop.
- **Infrastructure**: `internal/auth/infrastructure`
  - `db/`: Persistence using sqlc (`migrations/`, `repo/`).
  - `security/`: Implementations for JWT tokens and bcrypt hashing.
- **Interfaces**: `internal/auth/interfaces/http`
  - Primary adapters (HTTP handlers, DTOs, and router).
- **Bootstrap / Wiring**: `internal/auth/bootstrap`
  - Dependency injection and wiring of concrete infrastructure adapters to application ports.

## Development

From repo root:

- Build: `make build`
- Run: `make run`
- SQLC: `make sqlc` (inside `internal/auth/infrastructure`)

## Database schema

All tables live in the `auth` schema (e.g. `auth.users`, `auth.domain_events`, …).
