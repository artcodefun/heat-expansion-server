package events

import (
	"github.com/artcodefun/heat-expansion-server/internal/game/application/ports"
	"github.com/artcodefun/heat-expansion-server/internal/game/domain"
)

// SimplePublisher is a minimal in-process EventPublisher used for tests and dev.
// It forwards each published event to a single listener if one is registered.
// Not safe for multi-process usage; intended to be wired once at startup.
type SimplePublisher struct {
	listener func(domain.DomainEvent) error
}

var _ ports.EventPublisher = (*SimplePublisher)(nil)

func NewSimplePublisher() *SimplePublisher {
	return &SimplePublisher{}
}

// Publish forwards the event to the registered listener, if any.
func (p *SimplePublisher) Publish(event domain.DomainEvent) error {
	listener := p.listener
	if listener != nil {
		return listener(event)
	}
	return nil
}

// Listen registers a single callback for published events. It returns an unsubscribe function.
func (p *SimplePublisher) Listen(cb func(domain.DomainEvent) error) (unsubscribe func()) {
	p.listener = cb
	return func() {
		p.listener = nil
	}
}
