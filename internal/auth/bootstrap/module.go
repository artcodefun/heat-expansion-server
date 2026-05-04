package bootstrap

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"os"
	"sync"
	"time"

	"github.com/XSAM/otelsql"
	_ "github.com/lib/pq"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"

	httpapi "github.com/artcodefun/heat-expansion-server/internal/auth/interfaces/http"
)

type Module struct {
	Port  string
	DBURL string

	DB         *sql.DB
	Adapters   *Adapters
	Services   *AppServices
	Workers    *Workers
	Commands   *Commands
	Queries    *Queries
	HTTPServer *httpapi.Server
}

func NewModule() *Module {
	rabbitURL := os.Getenv("RABBITMQ_URL")
	port := os.Getenv("AUTH_PORT")
	dbURL := os.Getenv("AUTH_DB_URL")
	jwtSecret := os.Getenv("AUTH_JWT_SECRET")
	intExchange := os.Getenv("AUTH_INTEGRATION_EXCHANGE")

	if port == "" || dbURL == "" || jwtSecret == "" || rabbitURL == "" || intExchange == "" {
		log.Fatal("Missing required auth environment variables (AUTH_PORT, AUTH_DB_URL, AUTH_JWT_SECRET, RABBITMQ_URL, AUTH_INTEGRATION_EXCHANGE)")
	}

	db, err := otelsql.Open("postgres", dbURL,
		otelsql.WithAttributes(semconv.DBSystemPostgreSQL),
		otelsql.WithSpanOptions(otelsql.SpanOptions{Ping: false}),
	)
	if err != nil {
		log.Fatal("Failed to connect to auth database:", err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal("Auth database is unreachable:", err)
	}
	slog.Info("connected to auth database")

	adapters, err := NewAdapters(db, jwtSecret, rabbitURL, intExchange)
	if err != nil {
		log.Fatal("Failed to initialize auth adapters:", err)
	}

	services := NewAppServices(adapters)
	workers := NewWorkers(dbURL, services.Outbox, services.IntegrationOutbox)
	commands := NewCommands(adapters)
	queries := NewQueries(adapters)

	WireIntegrationEvents(services, adapters.Events)

	addr := fmt.Sprintf(":%s", port)
	router := httpapi.NewRouter(
		httpapi.Commands{Account: commands.Account},
		httpapi.Queries{Account: queries.Account},
		adapters.Translator,
	)
	server := httpapi.NewServer(router, addr)

	return &Module{
		Port:       port,
		DBURL:      dbURL,
		DB:         db,
		Adapters:   adapters,
		Services:   services,
		Workers:    workers,
		Commands:   commands,
		Queries:    queries,
		HTTPServer: server,
	}
}

func (m *Module) Run(ctx context.Context) {
	slog.Info("starting auth service", "port", m.Port)

	var wg sync.WaitGroup
	wg.Add(2)
	go func() { defer wg.Done(); m.Workers.DomainOutboxLoop(ctx) }()
	go func() { defer wg.Done(); m.Workers.IntegrationOutboxLoop(ctx) }()

	go func() {
		if err := m.HTTPServer.Start(); err != nil {
			slog.Error("auth http server stopped", "error", err)
		}
	}()

	<-ctx.Done()
	slog.Info("auth module: shutdown signal received, draining...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := m.HTTPServer.Shutdown(shutdownCtx); err != nil {
		slog.Error("auth http server shutdown error", "error", err)
	}

	wg.Wait()
	slog.Info("auth module: stopped")
}
