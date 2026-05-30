package repo

import (
	"log"
	"time"

	"github.com/lib/pq"
)

// PostgresListener listens to a Postgres NOTIFY channel and signals a Go channel.
type PostgresListener struct {
	Events chan struct{}
}

func NewPostgresListener(dsn string, channelName string) *PostgresListener {
	events := make(chan struct{}, 1)
	reportProblem := func(ev pq.ListenerEventType, err error) {
		if err != nil {
			log.Printf("Postgres listener error: %v", err)
		}
	}

	listener := pq.NewListener(dsn, 10*time.Second, time.Minute, reportProblem)
	if err := listener.Listen(channelName); err != nil {
		log.Fatalf("Failed to listen on %s: %v", channelName, err)
	}

	go func() {
		for {
			select {
			case <-listener.Notify:
				select {
				case events <- struct{}{}:
				default:
				}
			case <-time.After(time.Minute):
				if err := listener.Ping(); err != nil {
					log.Printf("Postgres listener ping error: %v", err)
				}
			}
		}
	}()

	return &PostgresListener{Events: events}
}
