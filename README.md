# Heat Expansion API

This project provides a backend for multiplayer video game called Heat Expansion 🔥🔥🔥.

## Example Environment Variables

Create a `.env` file in the project root with the following variables:

```
PORT=8080
DB_URL=postgres://user:password@localhost:5432/heatdb
JWT_SECRET=your_jwt_secret_here
CONTENT_DIR=./assets/content
STATIC_BASE_URL=http://localhost:8080/static
```

## Getting Started

1. Install Go (1.22+ recommended) and PostgreSQL.
2. Create a `.env` file in the project root using the example above and adjust values for your local setup.
3. Apply database migrations (for example, using `make migrate` if you have a Make target, or your usual migration tool pointing at `DB_URL`).
4. Run the API server from the `cmd/api` entrypoint:
   ```bash
   go run ./cmd/api
   ```
5. The server will start on `PORT` (default `8080`) and expose HTTP routes under `/api/v1`.

## Domain Overview

At a high level, this API models a 4X-style strategy game where players manage planetary bases, develop infrastructure, and interact with a shared world.

- **Players and Bases**: Each user controls one or more `UserBaseModel` instances. A base has coordinates on the map, local resource production, storage, buildings, technologies, armies, and a snapshot of stats (attack, defence, space, resources, etc.).
- **Resources and Economy**: Bases track currencies like credits, iron, titanium, and antimatter through `UserBaseStats`. Buildings, armies, and technologies all have `PriceModel` costs and can refund part of their value on deletion. Over time, stats are recalculated based on built structures, active tech, and queued items.
- **Buildings and Construction Queues**: Buildings (from `BuildItemPrototype`) are queued via methods like `AddToBuildQueue`. Items move through `pending` → `in production` → `present`, with production times and optional crystal skip prices. Finished buildings affect base stats (space capacity, attack/defence bonuses, unlocking army categories, etc.).
- **Armies and Units**: Armies are stacks of units based on `ArmyItemPrototype`, categorized by `ArmyCategory` (e.g. regular combat units vs spies). Units exist as `pending`, `in production`, `present`, or `deployed` (`ArmiesDeployed` grouped per operation). Production uses queues (`QueueArmy`, `MoveArmyQueue`) and depends on military buildings and resources.
- **Technologies and Progression**: Techs (`TechItemPrototype`) are researched over time (`TechItemInProgress` → `TechItemDone`). Each tech can unlock other content (buildings/armies) or provide effects like space, attack, defence, or resource bonuses via `TechnologyEffect`. Completed techs gate what prototypes are available to a base.
- **Sectors, Locations, and Exploration**: The world is divided into sectors (`SectorModel`) with `LocationDetails` (name/description/image). Around bases, content generation spawns resource sites (`ResourceLocationModel`) and dangerous locations. These locations hold `LocationResourceStats`, defending units, and defensive structures, acting as objectives for operations and loot.
- **Military Operations and Combat**: Military operations (`MilitaryOperation`) move units between coordinates as attacks or spy missions. Operations progress through phases (pending, outbound, at target, resolving, returning, completed) and simulate combat between `MilitaryUnit` stacks and defensive structures (`DefenseStructure`). Results (`AttackResult`, `SpyResult`) determine surviving units, damaged structures, loot recovered (`PriceModel`), and whether intel (`SectorScanReport`) is produced.
- **Scanning and Intel**: Successful operations or scans create `SectorScanReport` snapshots that record resource estimates, attack/defence strength, and flavor details for a sector, including whether the target was cloaked. Clients can use these reports to drive fog-of-war and scouting UI.
- **Activities and Timeline**: Domain events (`events.go`) and activity items capture important changes like finished builds, produced armies, resolved operations, and created scan reports. The read side exposes these as activity feeds so clients can reconstruct a player’s recent history.

## Architecture

This project uses Hexagonal Architecture (Ports and Adapters), DDD (Domain-driven design) and the CQRS (Command Query Responsibility Segregation) pattern, ensuring maintainability, testability, and clear separation of concerns.

![DDD Hexagon](https://raw.githubusercontent.com/Sairyss/domain-driven-hexagon/refs/heads/master/assets/images/DomainDrivenHexagon.png)

### Main Concepts

- **Core Business Logic** is located in the `internal/core` directory:
   - `domain/`: Domain model and rules: aggregates (e.g. `UserBaseModel`, `MilitaryOperation`), entities, value objects (coordinates, prices, stats, resource/location details), domain events, and most of the core business invariants and calculations.
   - `ports/`: Interfaces (ports) for repositories, schedulers, token providers, and other external dependencies.
   - `commands/`: Command handlers for write operations (mutations) that wrap domain aggregates and enforce access control.
   - `queries/`: Query handlers for read operations that work against read-store projections.
   - `cqrs/`: Shared CQRS definitions (command/query interfaces, contexts, readmodels).

- **Infrastructure** adapters (secondary adapters) live in `internal/infrastructure`:
   - `db/`:
     - `migrations/`: SQL migrations (schema + prototypes + base items).
     - `queries/`: sqlc query definitions for the write side.
     - `gen/`: sqlc-generated Go code for DB access.
     - `repo/`: concrete repository implementations for core ports (user bases, sectors, prototypes, operations, activities, etc.).
     - `dtos/` and `mappers/`: DTOs and mappers between DB models and domain models.
   - `readstore/`:
     - `queries/`: sqlc query definitions for the read side.
     - `gen/`: sqlc-generated read models and querier.
     - `mappers/`: translators from raw rows into `internal/core/cqrs/readmodels`.
     - `repo/`: read repositories exposed to query handlers.
   - `content/`: content generator used to provision sectors / world state from static assets.
   - `events/`: in-memory event publisher implementation used by command handlers.
   - `jobs/`: in-memory scheduler implementation for domain jobs (e.g. operation phase updates, world generation).
   - `security/`: password hashing and token provider implementations.

- **Interfaces** adapters (primary adapters) are in `internal/interfaces`:
   - `http/`:
     - `router.go`: builds the Gin engine, wires middleware, and mounts all HTTP routes.
     - `server.go`: thin wrapper owning the Gin engine lifecycle.
     - `handlers/`: HTTP handlers for users, bases, buildings, armies, tech, storage, sectors, operations, and activities; they translate HTTP into CQRS commands/queries.
     - `dtos/`: request/response DTOs that map CQRS readmodels and domain errors into JSON.
     - `middleware/`: auth middleware that uses the token provider from `ports` to populate CQRS contexts.
   - (Future) additional adapters like `ws/` or `cli/` can plug into the same core via ports.

- **Bootstrap** logic (dependency wiring, aggregators) is in `internal/bootstrap`:
   - `container.go`: wires concrete infrastructure adapters to core ports (repositories, tx manager, scheduler, security, content).
   - `commands.go`, `queries.go`: aggregate all command/query handlers into cohesive structs for the interfaces layer.
   - `app.go`: builds `App` (DB, adapters, commands, queries, HTTP server) and runs the scheduler + Gin server.

- The entry point `cmd/api/main.go` creates an `App` via `bootstrap.NewApp`, optionally seeds some prototype rows for local development, and then calls `App.Run` to start background jobs and the HTTP server.

This structure allows the backend to be easily extended, tested, and maintained, with clear boundaries between business logic and technical details. The CQRS pattern separates read and write concerns, improving scalability and clarity.
