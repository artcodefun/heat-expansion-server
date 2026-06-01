package bootstrap

import (
	"context"
	"crypto/ecdsa"
	"crypto/x509"
	"database/sql"
	"encoding/pem"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/XSAM/otelsql"
	_ "github.com/lib/pq"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"

	infraevents "github.com/artcodefun/heat-expansion-server/internal/billing/infrastructure/events"
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
	// IntegrationPublisher is the one adapter that starts background work (a
	// broker connection and reconnect goroutine) the moment it is built, so the
	// module constructs and owns it explicitly rather than hiding it in Adapters.
	IntegrationPublisher *infraevents.RabbitMQPublisher
	// Consumer projects inbound auth integration events into billing's local
	// users table. Unlike the publisher it dials lazily in Start(ctx) and tears
	// itself down on ctx cancellation, so the module owns it but does not close it.
	Consumer *infraevents.RabbitMQConsumer
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

	jwtPublicKey, err := parseECPublicKey(jwtPublicKeyPEM)
	if err != nil {
		log.Fatal("Failed to parse AUTH_JWT_PUBLIC_KEY:", err)
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

	intPublisher, err := infraevents.NewRabbitMQPublisher(rabbitURL, intExchange)
	if err != nil {
		log.Fatal("Failed to initialize billing RabbitMQ publisher:", err)
	}

	adapters, err := NewAdapters(db, jwtPublicKey, intPublisher, yookassaShopID, yookassaSecretKey)
	if err != nil {
		log.Fatal("Failed to initialize billing adapters:", err)
	}

	consumer := infraevents.NewRabbitMQConsumer(rabbitURL)

	appServices := NewAppServices(adapters)
	workers := NewWorkers(dbURL, appServices.Outbox, appServices.IntegrationOutbox, consumer)
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

		IntegrationPublisher: intPublisher,
		Consumer:             consumer,
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

	// Workers have drained, so no goroutine is still publishing or querying.
	// Release infrastructure connections so the broker and DB can reclaim
	// resources promptly instead of waiting for timeouts.
	if err := m.IntegrationPublisher.Close(); err != nil {
		slog.Error("billing integration publisher close error", "error", err)
	}
	if err := m.DB.Close(); err != nil {
		slog.Error("billing database close error", "error", err)
	}

	slog.Info("billing module: stopped")
}

func parseECPublicKey(pemStr string) (*ecdsa.PublicKey, error) {
	pemStr = strings.ReplaceAll(pemStr, `\n`, "\n")
	block, _ := pem.Decode([]byte(pemStr))
	if block == nil {
		return nil, errors.New("failed to decode PEM block for EC public key")
	}
	key, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	ecKey, ok := key.(*ecdsa.PublicKey)
	if !ok {
		return nil, errors.New("key is not an ECDSA public key")
	}
	if ecKey.Curve == nil || ecKey.Curve.Params().Name != "P-256" {
		return nil, errors.New("ECDSA public key must use P-256 curve for ES256")
	}
	return ecKey, nil
}
