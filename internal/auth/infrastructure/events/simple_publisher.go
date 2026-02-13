package events

import (
	"github.com/artcodefun/heat-expansion-server/internal/auth/application/ports"
	"github.com/artcodefun/heat-expansion-server/internal/auth/domain"
)

// SimplePublisher is a minimal in-process EventPublisher.
type SimplePublisher struct {
	listener func(domain.DomainEvent) error
}

var _ ports.EventPublisher = (*SimplePublisher)(nil)

func NewSimplePublisher() *SimplePublisher {
	return &SimplePublisher{}
}

// Publish forwards the event to the registered listener, if any.
func (p *SimplePublisher) Publish(event domain.DomainEvent) error {
	if p.listener != nil {
		return p.listener(event)
	}
	return nil
}

// Listen registers a single callback for published events.
func (p *SimplePublisher) Listen(cb func(domain.DomainEvent) error) {
	p.listener = cb
}
