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
	jwtSecret := os.Getenv("BILLING_JWT_SECRET")
	intExchange := os.Getenv("BILLING_INTEGRATION_EXCHANGE")
	rabbitURL := os.Getenv("RABBITMQ_URL")
	yookassaShopID := os.Getenv("YOOKASSA_SHOP_ID")
	yookassaSecretKey := os.Getenv("YOOKASSA_SECRET_KEY")

	if port == "" || dbURL == "" || jwtSecret == "" || intExchange == "" || rabbitURL == "" {
		log.Fatal("Missing required billing environment variables (BILLING_PORT, BILLING_DB_URL, BILLING_JWT_SECRET, BILLING_INTEGRATION_EXCHANGE, RABBITMQ_URL)")
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

	adapters, err := NewAdapters(db, jwtSecret, rabbitURL, intExchange, yookassaShopID, yookassaSecretKey)
	if err != nil {
		log.Fatal("Failed to initialize billing adapters:", err)
	}

	appServices := NewAppServices(adapters)
	workers := NewWorkers(dbURL, appServices.Outbox, appServices.IntegrationOutbox)
	commands := NewCommands(adapters)
	queries := NewQueries(adapters)

	WireIntegrationEvents(appServices, adapters.Events)

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

func (m *Module) Run(ctx context.Context) {
	slog.Info("starting billing service", "port", m.Port)

	m.Workers.Start(ctx)

	go func() {
		if err := m.HTTPServer.Start(); err != nil {
			slog.Error("billing http server stopped", "error", err)
		}
	}()

	<-ctx.Done()
	slog.Info("billing module: shutdown signal received, draining...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := m.HTTPServer.Shutdown(shutdownCtx); err != nil {
		slog.Error("billing http server shutdown error", "error", err)
	}

	m.Workers.Wait()
	slog.Info("billing module: stopped")
}
