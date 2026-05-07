package mappers

import (
	"encoding/json"
	"fmt"

	authevents "github.com/artcodefun/heat-expansion-server/contracts/auth/events"
	"github.com/artcodefun/heat-expansion-server/internal/auth/domain"
	"github.com/artcodefun/heat-expansion-server/internal/auth/infrastructure/db/dtos"
)

const (
	evKindAccountRegistered = "AccountRegistered"
)

// EncodeDomainEvent converts a typed DomainEvent into a kind label and JSON payload.
func EncodeDomainEvent(ev domain.DomainEvent) (kind string, payload []byte, err error) {
	switch e := ev.(type) {
	case domain.AccountRegisteredEvent:
		payload, err = json.Marshal(dtos.AccountRegisteredEventDTOFromDomain(e))
		return evKindAccountRegistered, payload, err
	case *domain.AccountRegisteredEvent:
		payload, err = json.Marshal(dtos.AccountRegisteredEventDTOFromDomain(*e))
		return evKindAccountRegistered, payload, err
	default:
		return "", nil, fmt.Errorf("unsupported domain event type: %T", ev)
	}
}

// DecodeDomainEvent reconstructs a DomainEvent from its kind label and JSON payload.
func DecodeDomainEvent(kind string, payload []byte) (domain.DomainEvent, error) {
	switch kind {
	case evKindAccountRegistered:
		var d dtos.AccountRegisteredEventDTO
		if err := json.Unmarshal(payload, &d); err != nil {
			return nil, err
		}
		return dtos.AccountRegisteredEventFromDTO(d), nil
	default:
		return nil, fmt.Errorf("unknown domain event kind: %s", kind)
	}
}

// EncodeIntegrationEvent converts a typed IntegrationEvent into a kind label and JSON payload.
func EncodeIntegrationEvent(ev authevents.IntegrationEvent) (kind string, payload []byte, err error) {
	payload, err = authevents.Marshal(ev)
	return ev.Type, payload, err
}

// DecodeIntegrationEvent reconstructs an IntegrationEvent from its kind label and JSON payload.
func DecodeIntegrationEvent(kind string, payload []byte) (authevents.IntegrationEvent, error) {
	return authevents.Unmarshal(payload)
}
