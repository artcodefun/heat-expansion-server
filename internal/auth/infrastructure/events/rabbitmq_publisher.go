package events

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/artcodefun/heat-expansion-server/contracts/auth"
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQPublisher struct {
	url      string
	exchange string

	mu   sync.RWMutex
	conn *amqp.Connection
}

func NewRabbitMQPublisher(url string, exchange string) (*RabbitMQPublisher, error) {
	p := &RabbitMQPublisher{
		url:      url,
		exchange: exchange,
	}

	if err := p.connect(); err != nil {
		return nil, err
	}

	return p, nil
}

func (p *RabbitMQPublisher) connect() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	conn, err := amqp.Dial(p.url)
	if err != nil {
		return err
	}

	// Ensure exchange exists
	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return err
	}
	defer ch.Close()

	err = ch.ExchangeDeclare(
		p.exchange,
		"topic",
		true,  // durable
		false, // auto-deleted
		false, // internal
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		conn.Close()
		return err
	}

	p.conn = conn

	// Reconnection goroutine
	go func() {
		closeChan := conn.NotifyClose(make(chan *amqp.Error))
		<-closeChan

		// Connection closed, try to reconnect until successful
		for {
			time.Sleep(5 * time.Second)
			if err := p.connect(); err == nil {
				return
			}
		}
	}()

	return nil
}

func (p *RabbitMQPublisher) Publish(event auth.IntegrationEvent) error {
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

	// Enable publisher confirms
	if err := ch.Confirm(false); err != nil {
		return fmt.Errorf("failed to enable publisher confirms: %w", err)
	}

	confirms := ch.NotifyPublish(make(chan amqp.Confirmation, 1))

	body, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = ch.PublishWithContext(ctx,
		p.exchange, // exchange
		event.Type, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Body:         body,
			MessageId:    event.ID.String(),
			Timestamp:    time.Unix(event.OccurredAt, 0),
		},
	)
	if err != nil {
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

func (p *RabbitMQPublisher) Close() error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.conn != nil {
		return p.conn.Close()
	}
	return nil
}
