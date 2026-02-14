package events

import (
	"context"
	"log"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type DeliveryHandler func(d amqp.Delivery) error

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
	if err := c.connect(); err != nil {
		return err
	}

	for _, sub := range c.subscriptions {
		go c.consumeLoop(ctx, sub)
	}

	return nil
}

func (c *RabbitMQConsumer) connect() error {
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
	go func() {
		closeChan := conn.NotifyClose(make(chan *amqp.Error))
		<-closeChan
		log.Println("RabbitMQ connection closed, reconnecting...")
		for {
			time.Sleep(5 * time.Second)
			if err := c.connect(); err == nil {
				log.Println("RabbitMQ reconnected")
				return
			}
		}
	}()

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
			log.Printf("Failed to open channel for queue %s: %v", sub.Queue, err)
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
			log.Printf("Failed to register consumer for queue %s: %v", sub.Queue, err)
			time.Sleep(5 * time.Second)
			continue
		}

		log.Printf("Started consuming from queue: %s", sub.Queue)

		done := false
		for !done {
			select {
			case <-ctx.Done():
				ch.Close()
				return
			case d, ok := <-msgs:
				if !ok {
					done = true
					break
				}
				if err := sub.Handler(d); err != nil {
					log.Printf("Error handling delivery from queue %s: %v", sub.Queue, err)
					d.Nack(false, true)
				} else {
					d.Ack(false)
				}
			}
		}
		ch.Close()
	}
}
