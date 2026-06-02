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

	httpapi "github.com/artcodefun/heat-expansion-server/internal/billing/interfaces/http"
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
	port := os.Getenv("BILLING_PORT")
	dbURL := os.Getenv("BILLING_DB_URL")
	jwtPublicKeyPEM := os.Getenv("AUTH_JWT_PUBLIC_KEY")
	intExchange := os.Getenv("BILLING_INTEGRATION_EXCHANGE")
	authIntExchange := os.Getenv("AUTH_INTEGRATION_EXCHANGE")
	rabbitURL := os.Getenv("RABBITMQ_URL")
	yookassaShopID := os.Getenv("YOOKASSA_SHOP_ID")
	yookassaSecretKey := os.Getenv("YOOKASSA_SECRET_KEY")

	if port == "" || dbURL == "" || jwtPublicKeyPEM == "" || intExchange == "" || authIntExchange == "" || rabbitURL == "" {
		log.Fatal("Missing required billing environment variables (BILLING_PORT, BILLING_DB_URL, AUTH_JWT_PUBLIC_KEY, BILLING_INTEGRATION_EXCHANGE, AUTH_INTEGRATION_EXCHANGE, RABBITMQ_URL)")
	}

	if yookassaShopID == "" || yookassaSecretKey == "" {
		log.Fatal("Missing required YooKassa environment variables (YOOKASSA_SHOP_ID, YOOKASSA_SECRET_KEY)")
	}

	db, err := otelsql.Open("postgres", dbURL,
		otelsql.WithAttributes(semconv.DBSystemPostgreSQL),
		otelsql.WithSpanOptions(otelsql.SpanOptions{Ping: false}),
	)
	if err != nil {
		log.Fatal("Failed to connect to billing database:", err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal("Billing database is unreachable:", err)
	}
	slog.Info("connected to billing database")

	intPublisher := rabbitmq.NewRabbitMQPublisher(rabbitURL, intExchange)

	adapters, err := NewAdapters(db, jwtPublicKeyPEM, intPublisher, yookassaShopID, yookassaSecretKey)
	if err != nil {
		log.Fatal("Failed to initialize billing adapters:", err)
	}

	consumer := rabbitmq.NewRabbitMQConsumer(rabbitURL)

	appServices := NewAppServices(adapters)
	workers := NewWorkers(dbURL, appServices.Outbox, appServices.IntegrationOutbox, consumer, intPublisher)
	commands := NewCommands(adapters)
	queries := NewQueries(adapters)

	WireIntegrationEvents(appServices, adapters.Events)
	WireConsumerIntegrationEvents(commands, consumer, authIntExchange, "billing.auth.integration.events")

	addr := fmt.Sprintf(":%s", port)
	router := httpapi.NewRouter(
		httpapi.Commands{Order: commands.Order},
		httpapi.Queries{Package: queries.Package, Order: queries.Order},
		adapters.Tokens,
		adapters.Translator,
	)
	server := httpapi.NewServer(router, addr)

	return &Module{
		Port:       port,
		DBURL:      dbURL,
		DB:         db,
		Adapters:   adapters,
		Services:   appServices,
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
	slog.InfoContext(ctx, "starting billing service", "port", m.Port)

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
		slog.InfoContext(ctx, "billing module: shutdown signal received, draining...")
	case err := <-httpErr:
		runErr = fmt.Errorf("billing http server failed: %w", err)
		slog.ErrorContext(ctx, "billing module: http server failed, draining...", "error", err)
		cancel()
	case err := <-workerErr:
		runErr = fmt.Errorf("billing workers failed: %w", err)
		slog.ErrorContext(ctx, "billing module: worker failed, draining...", "error", err)
		cancel()
	}

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer shutdownCancel()
	if err := m.HTTPServer.Shutdown(shutdownCtx); err != nil {
		slog.ErrorContext(ctx, "billing http server shutdown error", "error", err)
	}

	if err := m.Workers.Wait(); err != nil && runErr == nil {
		runErr = fmt.Errorf("billing workers failed: %w", err)
	}

	if err := m.DB.Close(); err != nil {
		slog.ErrorContext(ctx, "billing database close error", "error", err)
	}

	slog.InfoContext(ctx, "billing module: stopped")
	return runErr
}
