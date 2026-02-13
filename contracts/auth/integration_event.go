package auth

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/google/uuid"
)

// IntegrationEventPayload is an interface that all versioned payloads must implement.
type IntegrationEventPayload interface {
	IntegrationEventType() string
}

// IntegrationEvent is a generic envelope for all integration events.
type IntegrationEvent struct {
	ID         uuid.UUID               `json:"id"`
	Type       string                  `json:"type"`
	OriginID   uuid.UUID               `json:"origin_id"`
	OccurredAt int64                   `json:"occurred_at"`
	Payload    IntegrationEventPayload `json:"payload"`
}

// NewIntegrationEvent creates a new IntegrationEvent with a UUID v7 ID.
func NewIntegrationEvent(originID uuid.UUID, occurredAt int64, payload IntegrationEventPayload) IntegrationEvent {
	id, err := uuid.NewV7()
	if err != nil {
		// Fallback to v4 if entropy fails
		id = uuid.New()
	}

	return IntegrationEvent{
		ID:         id,
		Type:       payload.IntegrationEventType(),
		OriginID:   originID,
		OccurredAt: occurredAt,
		Payload:    payload,
	}
}

// Marshal converts the IntegrationEvent to a JSON byte slice.
func Marshal(event IntegrationEvent) ([]byte, error) {
	return json.Marshal(event)
}

// PayloadFactory is a function that returns a new instance of a payload.
type PayloadFactory func() IntegrationEventPayload

var (
	payloadRegistry   = make(map[string]PayloadFactory)
	payloadRegistryMu sync.RWMutex
)

// RegisterPayload registers a factory for a given event type string.
func RegisterPayload(eventType string, factory PayloadFactory) {
	if eventType == "" {
		panic("auth.RegisterPayload: eventType is empty")
	}
	if factory == nil {
		panic("auth.RegisterPayload: factory is nil")
	}

	payloadRegistryMu.Lock()
	defer payloadRegistryMu.Unlock()

	if _, exists := payloadRegistry[eventType]; exists {
		panic(fmt.Sprintf("auth.RegisterPayload: duplicate registration for type %q", eventType))
	}

	payloadRegistry[eventType] = factory
}

// UnknownPayload is a fallback for events not registered in the registry.
type UnknownPayload json.RawMessage

func (u UnknownPayload) IntegrationEventType() string { return "unknown" }
func (u UnknownPayload) MarshalJSON() ([]byte, error) { return json.RawMessage(u).MarshalJSON() }

// Unmarshal parses the JSON-encoded data and stores the result in an IntegrationEvent.
func Unmarshal(data []byte) (IntegrationEvent, error) {
	// First, unmarshal metadata to find the type
	var envelope struct {
		ID         uuid.UUID       `json:"id"`
		Type       string          `json:"type"`
		OriginID   uuid.UUID       `json:"origin_id"`
		OccurredAt int64           `json:"occurred_at"`
		Payload    json.RawMessage `json:"payload"`
	}

	if err := json.Unmarshal(data, &envelope); err != nil {
		return IntegrationEvent{}, err
	}

	payloadRegistryMu.RLock()
	factory, ok := payloadRegistry[envelope.Type]
	payloadRegistryMu.RUnlock()

	if !ok {
		// If unknown type, return with UnknownPayload as payload
		return IntegrationEvent{
			ID:         envelope.ID,
			Type:       envelope.Type,
			OriginID:   envelope.OriginID,
			OccurredAt: envelope.OccurredAt,
			Payload:    UnknownPayload(envelope.Payload),
		}, nil
	}

	payloadInstance := factory()
	if err := json.Unmarshal(envelope.Payload, payloadInstance); err != nil {
		return IntegrationEvent{}, fmt.Errorf("failed to unmarshal payload of type %s: %w", envelope.Type, err)
	}

	if payloadInstance.IntegrationEventType() != envelope.Type {
		return IntegrationEvent{}, fmt.Errorf(
			"payload type mismatch: envelope=%q payload=%q",
			envelope.Type,
			payloadInstance.IntegrationEventType(),
		)
	}

	return IntegrationEvent{
		ID:         envelope.ID,
		Type:       envelope.Type,
		OriginID:   envelope.OriginID,
		OccurredAt: envelope.OccurredAt,
		Payload:    payloadInstance,
	}, nil
}
