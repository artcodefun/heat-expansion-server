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

	"github.com/artcodefun/heat-expansion-server/internal/auth/infrastructure/events"
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
	// IntegrationPublisher is the one adapter that starts background work (a
	// broker connection and reconnect goroutine) the moment it is built, so the
	// module constructs and owns it explicitly rather than hiding it in Adapters.
	IntegrationPublisher *events.RabbitMQPublisher
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

	jwtPrivateKey, err := parseECPrivateKey(jwtPrivateKeyPEM)
	if err != nil {
		log.Fatal("Failed to parse AUTH_JWT_PRIVATE_KEY:", err)
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

	intPublisher, err := events.NewRabbitMQPublisher(rabbitURL, intExchange)
	if err != nil {
		log.Fatal("Failed to initialize auth RabbitMQ publisher:", err)
	}

	adapters, err := NewAdapters(db, jwtPrivateKey, intPublisher, smtpCfg)
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

		IntegrationPublisher: intPublisher,
	}
}

func (m *Module) Run(ctx context.Context) {
	slog.Info("starting auth service", "port", m.Port)

	m.Workers.Start(ctx)

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

	m.Workers.Wait()

	// Workers have drained, so no goroutine is still publishing or querying.
	// Release infrastructure connections so the broker and DB can reclaim
	// resources promptly instead of waiting for timeouts.
	if err := m.IntegrationPublisher.Close(); err != nil {
		slog.Error("auth integration publisher close error", "error", err)
	}
	if err := m.DB.Close(); err != nil {
		slog.Error("auth database close error", "error", err)
	}

	slog.Info("auth module: stopped")
}

func parseECPrivateKey(pemStr string) (*ecdsa.PrivateKey, error) {
	pemStr = strings.ReplaceAll(pemStr, `\n`, "\n")
	block, _ := pem.Decode([]byte(pemStr))
	if block == nil {
		return nil, errors.New("failed to decode PEM block for EC private key")
	}
	return x509.ParseECPrivateKey(block.Bytes)
}
