package domain

import (
	"github.com/google/uuid"
)

// DomainEvent is the interface for all domain events.
type DomainEvent interface {
	ID() uuid.UUID
	OccurredAt() int64
}

// BasicEvent carries the ID and timestamp for all domain events.
type BasicEvent struct {
	EventID   uuid.UUID
	Timestamp int64
}

func (e BasicEvent) ID() uuid.UUID {
	return e.EventID
}

func (e BasicEvent) OccurredAt() int64 {
	return e.Timestamp
}

func NewBasicEvent() BasicEvent {
	id := uuid.Must(uuid.NewV7())
	return BasicEvent{
		EventID:   id,
		Timestamp: NowUnix(),
	}
}

// EventProducer records domain events for aggregates/entities.
type EventProducer struct {
	events []DomainEvent
}

func (ep *EventProducer) AddEvent(event DomainEvent) {
	ep.events = append(ep.events, event)
}

func (ep *EventProducer) PullEvents() []DomainEvent {
	events := ep.events
	ep.events = nil
	return events
}

// OrderPaidEvent is emitted when an order is successfully paid.
type OrderPaidEvent struct {
	BasicEvent
	OrderID   uuid.UUID
	UserID    uuid.UUID
	PackageID uuid.UUID
	Crystals  int
}

func NewOrderPaidEvent(orderID, userID, packageID uuid.UUID, crystals int) OrderPaidEvent {
	return OrderPaidEvent{
		BasicEvent: NewBasicEvent(),
		OrderID:    orderID,
		UserID:     userID,
		PackageID:  packageID,
		Crystals:   crystals,
	}
}

// OrderFailedEvent is emitted when an order payment fails.
type OrderFailedEvent struct {
	BasicEvent
	OrderID uuid.UUID
}

func NewOrderFailedEvent(orderID uuid.UUID) OrderFailedEvent {
	return OrderFailedEvent{
		BasicEvent: NewBasicEvent(),
		OrderID:    orderID,
	}
}
