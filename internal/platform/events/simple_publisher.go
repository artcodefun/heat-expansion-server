package events

import "context"

// SimplePublisher is a minimal in-process event publisher.
// It forwards each published event to a single registered listener.
// Not safe for multi-process usage; intended to be wired once at startup.
type SimplePublisher[E any] struct {
	listener func(context.Context, E) error
}

func NewSimplePublisher[E any]() *SimplePublisher[E] {
	return &SimplePublisher[E]{}
}

// Publish forwards the event to the registered listener, if any.
func (p *SimplePublisher[E]) Publish(ctx context.Context, event E) error {
	if p.listener != nil {
		return p.listener(ctx, event)
	}
	return nil
}

// Listen registers a callback for published events and returns an unsubscribe func.
func (p *SimplePublisher[E]) Listen(cb func(context.Context, E) error) func() {
	p.listener = cb
	return func() { p.listener = nil }
}
