package mappers

import (
	"encoding/json"
	"fmt"

	billingevents "github.com/artcodefun/heat-expansion-server/contracts/billing/events"
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

func EncodeIntegrationEvent(ev billingevents.IntegrationEvent) (kind string, payload []byte, err error) {
	payload, err = billingevents.Marshal(ev)
	return ev.Type, payload, err
}

func DecodeIntegrationEvent(kind string, payload []byte) (billingevents.IntegrationEvent, error) {
	return billingevents.Unmarshal(payload)
}
