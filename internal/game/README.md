# Heat Expansion — Game Service

This directory contains the **game** service inside the Heat Expansion modular monolith. It handles the core gameplay mechanics, player progression, world exploration, and military operations.

## Domain Overview

At a high level, the Game service models a 4X-style strategy game where players manage planetary bases, develop infrastructure, and interact with a shared world.

- **Players and Bases**: Each user controls one or more `UserBaseModel` instances. A base has coordinates on the map, local resource production, storage, buildings, technologies, armies, and a snapshot of stats (attack, defence, space, resources, etc.).
- **Resources and Economy**: Bases track currencies like credits, iron, titanium, and antimatter through `UserBaseStats`. Buildings, armies, and technologies all have `PriceModel` costs and can refund part of their value on deletion. Over time, stats are recalculated based on built structures, active tech, and queued items.
- **Black Market**: Players can spend crystals to buy special resources, buildings, armies, and storage items that are only available from black market offers. Offers can be limited-time, and the read side exposes active offers so the client can browse and purchase them.
- **Buildings and Construction Queues**: Buildings (from `BuildItemPrototype`) are queued via methods like `AddToBuildQueue`. Items move through `pending` → `in production` → `present`, with production times and optional crystal skip prices. Finished buildings affect base stats (space capacity, attack/defence bonuses, unlocking army categories, etc.).
- **Armies and Units**: Armies are stacks of units based on `ArmyItemPrototype`, categorized by `ArmyCategory` (e.g. regular combat units vs spies). Units exist as `pending`, `in production`, `present`, or `deployed` (`ArmiesDeployed` grouped per operation). Production uses queues (`QueueArmy`, `MoveArmyQueue`) and depends on military buildings and resources.
- **Technologies and Progression**: Techs (`TechItemPrototype`) are researched over time (`TechItemInProgress` → `TechItemDone`). Each tech can unlock other content (buildings/armies) or provide effects like space, attack, defence, or resource bonuses via `TechnologyEffect`. Completed techs gate what prototypes are available to a base.
- **Sectors, Locations, and Exploration**: The world is divided into sectors (`SectorModel`) with `LocationDetails` (name/description/image). Around bases, content generation spawns resource sites (`ResourceLocationModel`) and dangerous locations. These locations hold `LocationResourceStats`, defending units, and defensive structures, acting as objectives for operations and loot.
- **Military Operations and Combat**: Military operations (`MilitaryOperation`) move units between coordinates as attacks or spy missions. Operations progress through phases (pending, outbound, at target, resolving, returning, completed) and simulate combat between `MilitaryUnit` stacks and defensive structures (`DefenseStructure`). Results (`AttackResult`, `SpyResult`) determine surviving units, damaged structures, loot recovered (`PriceModel`), and whether intel (`SectorScanReport`) is produced.
- **Diplomacy and Relationships**: Each pair of players maintains a `DiplomaticRelationship` with one of four statuses: `UNKNOWN` (no prior contact), `NEUTRAL`, `ALLIED`, or `WAR`. Statuses are connected to gameplay abilities e.g. attacks are only available to players at war, trades only available to allies.  Players negotiate through `DiplomaticRequest` proposals (coalition or ceasefire), which the recipient can accept, reject, or let expire within 24 hours. A `DiplomaticMessage` inbox records system notifications (war declarations, proposal outcomes) alongside player-sendable greetings and warnings.
- **Trade Operations**: Allied players can exchange resources, storage items, and armies through a `TradeOperation`. Each operation carries a `TradePayload` in both directions — an offered payload outbound and a requested payload on the return leg. Once accepted by receiver, the operation moves through phases (`PENDING → OUTBOUND → ARRIVED → RETURNING → COMPLETED`); the initiator can cancel mid-flight, in which case the convoy turns back from its current interpolated position. Pending offers expire after some time if the receiver does not respond.
- **Scanning and Intel**: Successful operations or scans create `SectorScanReport` snapshots that record resource estimates, attack/defence strength, and flavor details for a sector, including whether the target was cloaked. Clients can use these reports to drive fog-of-war and scouting UI.
- **Activities and Timeline**: Domain events (`events.go`) and activity items capture important changes like finished builds, produced armies, resolved operations, and created scan reports. The read side exposes these as activity feeds so clients can reconstruct a player’s recent history.

## Architecture

This service uses Hexagonal Architecture (Ports and Adapters), DDD (Domain-driven design) and the CQRS (Command Query Responsibility Segregation) pattern.

![DDD Hexagon](../../.github/DomainDrivenHexagon.png)

### Key Layers

- **Domain**: `internal/game/domain`
  - Business rules, aggregates (e.g. `UserBaseModel`, `MilitaryOperation`), entities, value objects, domain events, and core invariants.
- **Application**: `internal/game/application`
  - `commands/`: Write-side command handlers that wrap domain aggregates and enforce access control.
  - `queries/`: Read-side query handlers that work against read-store projections.
  - `cqrs/`: CQRS contract definitions and readmodels.
  - `ports/`: Interfaces for repositories, schedulers, and secondary adapters.
  - `services/`: App-level services like access control, provisioning, and the outbox loop.
- **Infrastructure**: `internal/game/infrastructure`
  - `db/`: Write-side persistence using sqlc (`migrations/`, `repo/`, etc.).
  - `readstore/`: Read-side persistence and cache for queries.
  - `i18n/`: Localization engine. Supports a hybrid model where **systemic** keys (errors) are embedded in the binary, and **content** keys (prototypes) are loaded from the `game.translations` table at startup via `TranslationRepo`.
  - Secondary adapters for `events/`, `jobs/`, `security/`, and `content/`.
- **Interfaces**: `internal/game/interfaces/http`
  - Primary adapters (HTTP handlers, DTOs, middleware, and router).
- **Bootstrap / Wiring**: `internal/game/bootstrap`
  - Dependency injection and wiring of concrete infrastructure adapters to application ports.

## HTTP API

Full OpenAPI spec: [`contracts/game/http/v1/openapi.yaml`](../../contracts/game/http/v1/openapi.yaml)

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
