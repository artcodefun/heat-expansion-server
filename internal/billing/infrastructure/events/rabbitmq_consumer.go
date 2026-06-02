package events

import (
	"context"
	"log/slog"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type DeliveryHandler func(ctx context.Context, d amqp.Delivery) error

type Subscription struct {
	Exchange   string
	Queue      string
	RoutingKey string
	Handler    DeliveryHandler
}

type RabbitMQConsumer struct {
	url string

	mu            sync.RWMutex
	conn          *amqp.Connection
	subscriptions []Subscription
}

func NewRabbitMQConsumer(url string) *RabbitMQConsumer {
	return &RabbitMQConsumer{url: url}
}

func (c *RabbitMQConsumer) Subscribe(exchange, queue, routingKey string, handler DeliveryHandler) {
	c.subscriptions = append(c.subscriptions, Subscription{
		Exchange:   exchange,
		Queue:      queue,
		RoutingKey: routingKey,
		Handler:    handler,
	})
}

// Start dials the broker, declares exchanges/queues/bindings, launches consume
// loops, and blocks until ctx is cancelled. It reconnects automatically on
// connection drops. On cancellation it waits for the consume loops to finish
// their in-flight deliveries (so acks still reach the broker), closes the
// connection so the broker can reclaim resources promptly, then returns nil.
// Returns a non-nil error only if the initial dial fails.
func (c *RabbitMQConsumer) Start(ctx context.Context) error {
	conn, err := dialAndSetupConsumer(c.url, c.subscriptions)
	if err != nil {
		return err
	}

	c.mu.Lock()
	c.conn = conn
	c.mu.Unlock()

	var loops sync.WaitGroup
	for _, sub := range c.subscriptions {
		loops.Go(func() { c.consumeLoop(ctx, sub) })
	}

	drain := func() {
		loops.Wait()
		c.mu.Lock()
		if c.conn != nil {
			c.conn.Close()
			c.conn = nil
		}
		c.mu.Unlock()
	}

	current := conn
	for {
		closeChan := current.NotifyClose(make(chan *amqp.Error, 1))
		select {
		case <-ctx.Done():
			drain()
			return nil
		case amqpErr := <-closeChan:
			if amqpErr != nil {
				slog.WarnContext(ctx, "rabbitmq consumer connection lost; reconnecting", "error", amqpErr.Error())
			}
		}

		for {
			select {
			case <-ctx.Done():
				drain()
				return nil
			case <-time.After(5 * time.Second):
			}
			newConn, err := dialAndSetupConsumer(c.url, c.subscriptions)
			if err != nil {
				slog.WarnContext(ctx, "rabbitmq consumer reconnect attempt failed", "error", err.Error())
				continue
			}
			c.mu.Lock()
			c.conn = newConn
			c.mu.Unlock()
			current = newConn
			slog.InfoContext(ctx, "rabbitmq consumer reconnected")
			break
		}
	}
}

func dialAndSetupConsumer(url string, subscriptions []Subscription) (*amqp.Connection, error) {
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

	for _, sub := range subscriptions {
		if err := ch.ExchangeDeclare(sub.Exchange, "topic", true, false, false, false, nil); err != nil {
			conn.Close()
			return nil, err
		}
		if _, err := ch.QueueDeclare(sub.Queue, true, false, false, false, nil); err != nil {
			conn.Close()
			return nil, err
		}
		if err := ch.QueueBind(sub.Queue, sub.RoutingKey, sub.Exchange, false, nil); err != nil {
			conn.Close()
			return nil, err
		}
	}

	return conn, nil
}

func (c *RabbitMQConsumer) consumeLoop(ctx context.Context, sub Subscription) {
	for {
		c.mu.RLock()
		conn := c.conn
		c.mu.RUnlock()

		if conn == nil || conn.IsClosed() {
			select {
			case <-ctx.Done():
				return
			case <-time.After(1 * time.Second):
			}
			continue
		}

		ch, err := conn.Channel()
		if err != nil {
			slog.WarnContext(ctx, "failed to open rabbitmq channel",
				"queue", sub.Queue,
				"exchange", sub.Exchange,
				"routing_key", sub.RoutingKey,
				"error", err.Error(),
			)
			select {
			case <-ctx.Done():
				return
			case <-time.After(5 * time.Second):
			}
			continue
		}

		msgs, err := ch.Consume(sub.Queue, "", false, false, false, false, nil)
		if err != nil {
			ch.Close()
			slog.WarnContext(ctx, "failed to register rabbitmq consumer",
				"queue", sub.Queue,
				"exchange", sub.Exchange,
				"routing_key", sub.RoutingKey,
				"error", err.Error(),
			)
			select {
			case <-ctx.Done():
				return
			case <-time.After(5 * time.Second):
			}
			continue
		}

		slog.InfoContext(ctx, "rabbitmq consumer loop started",
			"queue", sub.Queue,
			"exchange", sub.Exchange,
			"routing_key", sub.RoutingKey,
		)

		done := false
		for !done {
			select {
			case <-ctx.Done():
				ch.Close()
				return
			case d, ok := <-msgs:
				if !ok {
					slog.WarnContext(ctx, "rabbitmq delivery channel closed; restarting consumer loop",
						"queue", sub.Queue,
						"exchange", sub.Exchange,
						"routing_key", sub.RoutingKey,
					)
					done = true
					break
				}
				if err := sub.Handler(ctx, d); err != nil {
					slog.WarnContext(ctx, "rabbitmq delivery handler failed",
						"queue", sub.Queue,
						"exchange", sub.Exchange,
						"routing_key", sub.RoutingKey,
						"delivery_tag", d.DeliveryTag,
						"error", err.Error(),
					)
					d.Nack(false, true)
				} else {
					d.Ack(false)
				}
			}
		}
		ch.Close()
	}
}
