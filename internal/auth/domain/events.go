package domain

import (
	"time"

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
		Timestamp: time.Now().Unix(),
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

// AccountRegisteredEvent is emitted when a new account is created.
type AccountRegisteredEvent struct {
	BasicEvent
	AccountID uuid.UUID
	Name      string
	Email     string
}

func NewAccountRegisteredEvent(accountID uuid.UUID, name, email string) AccountRegisteredEvent {
	return AccountRegisteredEvent{
		BasicEvent: NewBasicEvent(),
		AccountID:  accountID,
		Name:       name,
		Email:      email,
	}
}
