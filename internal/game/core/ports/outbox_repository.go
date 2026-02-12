package ports

import "github.com/artcodefun/heat-expansion-api/internal/game/core/domain"

// OutboxEventRecord represents a persisted domain event entry used by a
// transactional outbox. The concrete event payload is held as a typed
// domain.DomainEvent; infrastructure is responsible for encoding/decoding it
// when storing to the underlying database.
//
// All timestamps are Unix seconds.
type OutboxEventRecord struct {
	ID          int64
	Event       domain.DomainEvent
	CreatedAt   int64
	Published   bool
	PublishedAt int64
}

// OutboxEventRepository provides persistence for domain event outbox records.
// Implementations are expected to be used within a TransactionManager.WithTx
// scope via the Tx(tx) method.
type OutboxEventRepository interface {
	// Save persists a batch of domain events into the outbox. Implementations
	// are responsible for translating typed events into concrete records.
	Save(events []domain.DomainEvent) error
	// ClaimUnpublished returns a batch of not-yet-published events ordered by ID
	// up to the provided limit, using database-level locking so that multiple
	// workers can safely process the outbox in parallel.
	ClaimUnpublished(limit int) ([]*OutboxEventRecord, error)
	// MarkPublished marks an event as published at the given timestamp.
	MarkPublished(id int64, publishedAt int64) error
	// Tx binds the repository to a concrete transaction implementation.
	Tx(tx Transaction) OutboxEventRepository
}
