package events

import (
	"sync"

	"github.com/artcodefun/heat-expansion-api/internal/core/domain"
	"github.com/artcodefun/heat-expansion-api/internal/core/ports"
)

// InMemoryPublisher is a simple in-process event publisher useful for tests and dev.
// It stores all published events in memory and optionally notifies subscribers.
// Not safe for multi-process usage. Safe for concurrent goroutines.
type InMemoryPublisher struct {
	mu          sync.Mutex
	subscribers map[int]func(domain.DomainEvent)
	nextSubID   int
	events      []domain.DomainEvent
}

var _ ports.EventPublisher = (*InMemoryPublisher)(nil)

func NewInMemoryPublisher() *InMemoryPublisher {
	return &InMemoryPublisher{
		subscribers: make(map[int]func(domain.DomainEvent)),
		events:      make([]domain.DomainEvent, 0, 128),
	}
}

// Publish records the event and notifies subscribers.
func (p *InMemoryPublisher) Publish(event domain.DomainEvent) error {
	p.mu.Lock()
	p.events = append(p.events, event)
	subs := make([]func(domain.DomainEvent), 0, len(p.subscribers))
	for _, cb := range p.subscribers {
		subs = append(subs, cb)
	}
	p.mu.Unlock()
	// Notify without holding the lock to avoid deadlocks
	for _, cb := range subs {
		cb(event)
	}
	return nil
}

// Subscribe registers a callback for published events. It returns an unsubscribe function.
func (p *InMemoryPublisher) Subscribe(cb func(domain.DomainEvent)) (unsubscribe func()) {
	p.mu.Lock()
	id := p.nextSubID
	p.nextSubID++
	p.subscribers[id] = cb
	p.mu.Unlock()
	return func() {
		p.mu.Lock()
		delete(p.subscribers, id)
		p.mu.Unlock()
	}
}

// Events returns a snapshot of all events published so far.
func (p *InMemoryPublisher) Events() []domain.DomainEvent {
	p.mu.Lock()
	defer p.mu.Unlock()
	out := make([]domain.DomainEvent, len(p.events))
	copy(out, p.events)
	return out
}

// Drain returns all recorded events and clears the buffer.
func (p *InMemoryPublisher) Drain() []domain.DomainEvent {
	p.mu.Lock()
	defer p.mu.Unlock()
	out := p.events
	p.events = nil
	return out
}
