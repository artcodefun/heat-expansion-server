package ports

// import (
// 	"github.com/artcodefun/heat-expansion-api/internal/core/domain"
// 	"github.com/google/uuid"
// )

// // OutboxStatus represents delivery state of an outbox record.
// type OutboxStatus string

// const (
// 	OutboxPending OutboxStatus = "PENDING"
// 	OutboxSent    OutboxStatus = "SENT"
// 	OutboxFailed  OutboxStatus = "FAILED"
// )

// // OutboxRecord carries a domain event to be delivered asynchronously.
// // Serialization concerns (type strings, JSON, timestamps) are an infrastructure detail.
// type OutboxRecord struct {
// 	ID     uuid.UUID
// 	Event  domain.DomainEvent
// 	Status OutboxStatus
// }

// // OutboxRepository persists and retrieves outbox records.
// // Implementations may assign IDs on insert if missing and must populate Event when fetching.
// type OutboxRepository interface {
// 	Enqueue(records []OutboxRecord) error
// 	FetchPending(limit int) ([]OutboxRecord, error)
// 	MarkSent(id uuid.UUID) error
// 	MarkFailed(id uuid.UUID) error
// }
