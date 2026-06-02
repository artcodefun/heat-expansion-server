package bootstrap

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/XSAM/otelsql"
	_ "github.com/lib/pq"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"

	httpapi "github.com/artcodefun/heat-expansion-server/internal/auth/interfaces/http"
	"github.com/artcodefun/heat-expansion-server/internal/platform/rabbitmq"
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
	jwtPrivateKeyPEM := os.Getenv("AUTH_JWT_PRIVATE_KEY")
	intExchange := os.Getenv("AUTH_INTEGRATION_EXCHANGE")

	smtpCfg := SMTPConfig{
		Host:     os.Getenv("AUTH_SMTP_HOST"),
		User:     os.Getenv("AUTH_SMTP_USER"),
		Password: os.Getenv("AUTH_SMTP_PASSWORD"),
		From:     os.Getenv("AUTH_SMTP_FROM"),
	}

	if port == "" || dbURL == "" || jwtPrivateKeyPEM == "" || rabbitURL == "" || intExchange == "" {
		log.Fatal("Missing required auth environment variables (AUTH_PORT, AUTH_DB_URL, AUTH_JWT_PRIVATE_KEY, RABBITMQ_URL, AUTH_INTEGRATION_EXCHANGE)")
	}

	if smtpCfg.Host == "" || smtpCfg.User == "" || smtpCfg.Password == "" || smtpCfg.From == "" {
		log.Fatal("Missing required SMTP environment variables (AUTH_SMTP_HOST, AUTH_SMTP_USER, AUTH_SMTP_PASSWORD, AUTH_SMTP_FROM)")
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

	intPublisher := rabbitmq.NewRabbitMQPublisher(rabbitURL, intExchange)

	adapters, err := NewAdapters(db, jwtPrivateKeyPEM, intPublisher, smtpCfg)
	if err != nil {
		log.Fatal("Failed to initialize auth adapters:", err)
	}

	services := NewAppServices(adapters)
	workers := NewWorkers(dbURL, services.Outbox, services.IntegrationOutbox, intPublisher)
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

// Run connects the module's infrastructure, serves until the context is
// cancelled or a fatal error occurs, then drains and releases resources.
// Construction does no I/O; everything cancelable happens here.
func (m *Module) Run(ctx context.Context) error {
	slog.InfoContext(ctx, "starting auth service", "port", m.Port)

	// Module-local cancel: a fatal HTTP error must stop the workers (and a
	// fatal worker error must stop the HTTP server), otherwise the drain
	// below would block forever.
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	m.Workers.Start(ctx)

	// Workers only fail on startup (e.g. the broker is unreachable); surface
	// that like an HTTP serve failure so the module drains and reports it.
	workerErr := make(chan error, 1)
	go func() {
		if err := m.Workers.Wait(); err != nil {
			workerErr <- err
		}
	}()

	httpErr := make(chan error, 1)
	go func() {
		if err := m.HTTPServer.Start(); err != nil {
			httpErr <- err
		}
	}()

	var runErr error
	select {
	case <-ctx.Done():
		slog.InfoContext(ctx, "auth module: shutdown signal received, draining...")
	case err := <-httpErr:
		runErr = fmt.Errorf("auth http server failed: %w", err)
		slog.ErrorContext(ctx, "auth module: http server failed, draining...", "error", err)
		cancel()
	case err := <-workerErr:
		runErr = fmt.Errorf("auth workers failed: %w", err)
		slog.ErrorContext(ctx, "auth module: worker failed, draining...", "error", err)
		cancel()
	}

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer shutdownCancel()
	if err := m.HTTPServer.Shutdown(shutdownCtx); err != nil {
		slog.ErrorContext(ctx, "auth http server shutdown error", "error", err)
	}

	if err := m.Workers.Wait(); err != nil && runErr == nil {
		runErr = fmt.Errorf("auth workers failed: %w", err)
	}

	if err := m.DB.Close(); err != nil {
		slog.ErrorContext(ctx, "auth database close error", "error", err)
	}

	slog.InfoContext(ctx, "auth module: stopped")
	return runErr
}
