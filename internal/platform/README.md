# Platform

Shared infrastructure adapters used across multiple services. Anything that would otherwise be duplicated between services lives here.

## Packages

- **`rabbitmq/`** — `RabbitMQPublisher` and `RabbitMQConsumer`. Both reconnect automatically on connection drops and expose a blocking `Start(ctx)` that fits the module lifecycle convention.
- **`events/`** — `SimplePublisher[E]`, a generic in-process event publisher. Services embed it via `SimplePublisher[domain.DomainEvent]` and wire listeners at bootstrap.
- **`security/`** — `SimpleTokenValidator`, an ES256 JWT validator that parses a PEM public key and verifies token signature and expiry; and `BcryptHasher`, a bcrypt-based password hasher (hash + verify) at the library's default cost.
- **`i18n/`** — `Translator`, the core locale bundle engine. Services embed it in their own `SimpleTranslator` and call `LoadFromJsonFiles` to populate systemic keys; services that also need DB-backed content translations implement `LoadFromRepo` on top.

## Adding new shared adapters

Move an adapter here when it is needed by more than one service and contains no service-specific logic. Keep service-specific wiring (bootstrap files, port interface checks) in the service itself.
