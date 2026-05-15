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
	return &RabbitMQConsumer{
		url: url,
	}
}

func (c *RabbitMQConsumer) Subscribe(exchange, queue, routingKey string, handler DeliveryHandler) {
	c.subscriptions = append(c.subscriptions, Subscription{
		Exchange:   exchange,
		Queue:      queue,
		RoutingKey: routingKey,
		Handler:    handler,
	})
}

func (c *RabbitMQConsumer) Start(ctx context.Context) error {
	if err := c.connect(ctx); err != nil {
		return err
	}

	for _, sub := range c.subscriptions {
		go c.consumeLoop(ctx, sub)
	}

	return nil
}

func (c *RabbitMQConsumer) connect(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	conn, err := amqp.Dial(c.url)
	if err != nil {
		return err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return err
	}
	defer ch.Close()

	for _, sub := range c.subscriptions {
		// Declare exchange
		err = ch.ExchangeDeclare(
			sub.Exchange,
			"topic",
			true,
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			conn.Close()
			return err
		}

		// Declare queue
		_, err = ch.QueueDeclare(
			sub.Queue,
			true,
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			conn.Close()
			return err
		}

		// Bind queue
		err = ch.QueueBind(
			sub.Queue,
			sub.RoutingKey,
			sub.Exchange,
			false,
			nil,
		)
		if err != nil {
			conn.Close()
			return err
		}
	}

	c.conn = conn

	// Reconnection
	go func(ctx context.Context) {
		closeChan := conn.NotifyClose(make(chan *amqp.Error))
		select {
		case <-ctx.Done():
			return
		case err := <-closeChan:
			if err == nil {
				return
			}

			slog.WarnContext(ctx, "rabbitmq connection closed; reconnecting",
				"error", err.Error(),
			)
		}

		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(5 * time.Second):
			}

			if err := c.connect(ctx); err == nil {
				slog.InfoContext(ctx, "rabbitmq reconnected")
				return
			} else {
				slog.WarnContext(ctx, "rabbitmq reconnect attempt failed", "error", err.Error())
			}
		}
	}(ctx)

	return nil
}

func (c *RabbitMQConsumer) consumeLoop(ctx context.Context, sub Subscription) {
	for {
		c.mu.RLock()
		conn := c.conn
		c.mu.RUnlock()

		if conn == nil || conn.IsClosed() {
			time.Sleep(1 * time.Second)
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
			time.Sleep(5 * time.Second)
			continue
		}

		msgs, err := ch.Consume(
			sub.Queue,
			"",
			false, // manual ack
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			ch.Close()
			slog.WarnContext(ctx, "failed to register rabbitmq consumer",
				"queue", sub.Queue,
				"exchange", sub.Exchange,
				"routing_key", sub.RoutingKey,
				"error", err.Error(),
			)
			time.Sleep(5 * time.Second)
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
