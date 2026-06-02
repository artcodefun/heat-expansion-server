package events

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"sync"
	"time"

	billingevents "github.com/artcodefun/heat-expansion-server/contracts/billing/events"
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQPublisher struct {
	url      string
	exchange string

	mu   sync.RWMutex
	conn *amqp.Connection
}

func NewRabbitMQPublisher(url string, exchange string) *RabbitMQPublisher {
	return &RabbitMQPublisher{
		url:      url,
		exchange: exchange,
	}
}

// Start dials the broker, declares the exchange, and blocks until ctx is
// cancelled. It reconnects automatically on connection drops. On cancellation
// it closes the connection so the broker can reclaim resources promptly, then
// returns nil. Returns a non-nil error only if the initial dial fails.
func (p *RabbitMQPublisher) Start(ctx context.Context) error {
	conn, err := dialAndSetupPublisher(p.url, p.exchange)
	if err != nil {
		return err
	}

	p.mu.Lock()
	p.conn = conn
	p.mu.Unlock()

	current := conn
	for {
		closeChan := current.NotifyClose(make(chan *amqp.Error, 1))
		select {
		case <-ctx.Done():
			p.mu.Lock()
			p.conn.Close()
			p.conn = nil
			p.mu.Unlock()
			return nil
		case amqpErr := <-closeChan:
			if amqpErr != nil {
				slog.WarnContext(ctx, "rabbitmq publisher connection lost; reconnecting", "error", amqpErr.Error())
			}
		}

		for {
			select {
			case <-ctx.Done():
				return nil
			case <-time.After(5 * time.Second):
			}
			newConn, err := dialAndSetupPublisher(p.url, p.exchange)
			if err != nil {
				slog.WarnContext(ctx, "rabbitmq publisher reconnect attempt failed", "error", err.Error())
				continue
			}
			p.mu.Lock()
			p.conn = newConn
			p.mu.Unlock()
			current = newConn
			slog.InfoContext(ctx, "rabbitmq publisher reconnected")
			break
		}
	}
}

func dialAndSetupPublisher(url, exchange string) (*amqp.Connection, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}
	defer ch.Close()

	if err := ch.ExchangeDeclare(
		exchange,
		"topic",
		true,  // durable
		false, // auto-deleted
		false, // internal
		false, // no-wait
		nil,
	); err != nil {
		conn.Close()
		return nil, err
	}

	return conn, nil
}

func (p *RabbitMQPublisher) Publish(ctx context.Context, event billingevents.IntegrationEvent) error {
	p.mu.RLock()
	conn := p.conn
	p.mu.RUnlock()

	if conn == nil || conn.IsClosed() {
		return fmt.Errorf("rabbitmq connection is not available")
	}

	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open channel: %w", err)
	}
	defer ch.Close()

	if err := ch.Confirm(false); err != nil {
		return fmt.Errorf("failed to enable publisher confirms: %w", err)
	}

	confirms := ch.NotifyPublish(make(chan amqp.Confirmation, 1))

	body, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if err := ch.PublishWithContext(ctx,
		p.exchange,
		event.Type,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Body:         body,
			MessageId:    event.ID.String(),
			Timestamp:    time.Unix(event.OccurredAt, 0),
		},
	); err != nil {
		return fmt.Errorf("failed to publish: %w", err)
	}

	select {
	case confirm := <-confirms:
		if confirm.Ack {
			return nil
		}
		return fmt.Errorf("publisher confirm NACK received")
	case <-ctx.Done():
		return fmt.Errorf("publisher confirm timeout: %w", ctx.Err())
	}
}
