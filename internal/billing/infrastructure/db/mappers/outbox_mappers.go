package mappers

import (
	"encoding/json"
	"fmt"

	"github.com/artcodefun/heat-expansion-server/contracts/events"
	"github.com/artcodefun/heat-expansion-server/internal/billing/domain"
	"github.com/artcodefun/heat-expansion-server/internal/billing/infrastructure/db/dtos"
)

const (
	evKindOrderPaid   = "OrderPaid"
	evKindOrderFailed = "OrderFailed"
)

func EncodeDomainEvent(ev domain.DomainEvent) (kind string, payload []byte, err error) {
	switch e := ev.(type) {
	case domain.OrderPaidEvent:
		payload, err = json.Marshal(dtos.OrderPaidEventDTOFromDomain(e))
		return evKindOrderPaid, payload, err
	case domain.OrderFailedEvent:
		payload, err = json.Marshal(dtos.OrderFailedEventDTOFromDomain(e))
		return evKindOrderFailed, payload, err
	default:
		return "", nil, fmt.Errorf("unsupported domain event type: %T", ev)
	}
}

func DecodeDomainEvent(kind string, payload []byte) (domain.DomainEvent, error) {
	switch kind {
	case evKindOrderPaid:
		var d dtos.OrderPaidEventDTO
		if err := json.Unmarshal(payload, &d); err != nil {
			return nil, err
		}
		return dtos.OrderPaidEventFromDTO(d), nil
	case evKindOrderFailed:
		var d dtos.OrderFailedEventDTO
		if err := json.Unmarshal(payload, &d); err != nil {
			return nil, err
		}
		return dtos.OrderFailedEventFromDTO(d), nil
	default:
		return nil, fmt.Errorf("unknown domain event kind: %s", kind)
	}
}

func EncodeIntegrationEvent(ev events.IntegrationEvent) (kind string, payload []byte, err error) {
	payload, err = json.Marshal(ev)
	return ev.Type, payload, err
}

func DecodeIntegrationEvent(_ string, payload []byte) (events.IntegrationEvent, error) {
	var ev events.IntegrationEvent
	if err := json.Unmarshal(payload, &ev); err != nil {
		return events.IntegrationEvent{}, err
	}
	return ev, nil
}
