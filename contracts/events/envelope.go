package events

import (
	"encoding/json"

	"github.com/google/uuid"
)

// IntegrationEvent is a generic envelope for all integration events.
type IntegrationEvent struct {
	ID         uuid.UUID       `json:"id"`
	Type       string          `json:"type"`
	OriginID   uuid.UUID       `json:"origin_id"`
	OccurredAt int64           `json:"occurred_at"`
	Payload    json.RawMessage `json:"payload"`
}

// NewIntegrationEvent creates a new IntegrationEvent with a UUID v7 ID.
func NewIntegrationEvent(originID uuid.UUID, occurredAt int64, eventType string, payload json.RawMessage) IntegrationEvent {
	id, err := uuid.NewV7()
	if err != nil {
		id = uuid.New()
	}
	return IntegrationEvent{
		ID:         id,
		Type:       eventType,
		OriginID:   originID,
		OccurredAt: occurredAt,
		Payload:    payload,
	}
}
