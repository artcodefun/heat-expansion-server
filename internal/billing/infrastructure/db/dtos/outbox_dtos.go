package dtos

import (
	"github.com/artcodefun/heat-expansion-server/internal/billing/domain"
	"github.com/google/uuid"
)

type OrderPaidEventDTO struct {
	ID         uuid.UUID `json:"id"`
	OccurredAt int64     `json:"occurred_at"`
	OrderID    uuid.UUID `json:"order_id"`
	UserID     uuid.UUID `json:"user_id"`
	PackageID  uuid.UUID `json:"package_id"`
	Crystals   int       `json:"crystals"`
}

func OrderPaidEventDTOFromDomain(e domain.OrderPaidEvent) OrderPaidEventDTO {
	return OrderPaidEventDTO{
		ID:         e.ID(),
		OccurredAt: e.OccurredAt(),
		OrderID:    e.OrderID,
		UserID:     e.UserID,
		PackageID:  e.PackageID,
		Crystals:   e.Crystals,
	}
}

func OrderPaidEventFromDTO(d OrderPaidEventDTO) domain.OrderPaidEvent {
	return domain.OrderPaidEvent{
		BasicEvent: domain.BasicEvent{EventID: d.ID, Timestamp: d.OccurredAt},
		OrderID:    d.OrderID,
		UserID:     d.UserID,
		PackageID:  d.PackageID,
		Crystals:   d.Crystals,
	}
}

type OrderFailedEventDTO struct {
	ID         uuid.UUID `json:"id"`
	OccurredAt int64     `json:"occurred_at"`
	OrderID    uuid.UUID `json:"order_id"`
}

func OrderFailedEventDTOFromDomain(e domain.OrderFailedEvent) OrderFailedEventDTO {
	return OrderFailedEventDTO{
		ID:         e.ID(),
		OccurredAt: e.OccurredAt(),
		OrderID:    e.OrderID,
	}
}

func OrderFailedEventFromDTO(d OrderFailedEventDTO) domain.OrderFailedEvent {
	return domain.OrderFailedEvent{
		BasicEvent: domain.BasicEvent{EventID: d.ID, Timestamp: d.OccurredAt},
		OrderID:    d.OrderID,
	}
}
