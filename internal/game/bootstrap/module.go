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

	grpcapi "github.com/artcodefun/heat-expansion-server/internal/game/interfaces/grpc"
	httpapi "github.com/artcodefun/heat-expansion-server/internal/game/interfaces/http"
	"github.com/artcodefun/heat-expansion-server/internal/platform/rabbitmq"
)

type Module struct {
	Port            string
	DBURL           string
	LogLevel        string
	StaticBaseURL   string
	RabbitURL       string
	AuthExchange    string
	BillingExchange string

	DB         *sql.DB
	Adapters   *Adapters
	Services   *AppServices
	Workers    *Workers
	Commands   *Commands
	Queries    *Queries
	HTTPServer *httpapi.Server
	GRPCServer *grpcapi.Server
}

func NewModule() *Module {
	port := os.Getenv("GAME_PORT")
	dbURL := os.Getenv("GAME_DB_URL")
	jwtPublicKeyPEM := os.Getenv("AUTH_JWT_PUBLIC_KEY")
	staticBaseURL := os.Getenv("GAME_STATIC_BASE_URL")
	rabbitURL := os.Getenv("RABBITMQ_URL")
	authExchange := os.Getenv("AUTH_INTEGRATION_EXCHANGE")
	billingExchange := os.Getenv("BILLING_INTEGRATION_EXCHANGE")
	grpcPort := os.Getenv("GAME_GRPC_PORT")
	grpcKey := os.Getenv("GAME_GRPC_KEY")

	if port == "" || dbURL == "" || jwtPublicKeyPEM == "" || staticBaseURL == "" || rabbitURL == "" || authExchange == "" || billingExchange == "" || grpcPort == "" || grpcKey == "" {
		log.Fatal("Missing required environment variables. Please check your .env file.")
	}

	db, err := otelsql.Open("postgres", dbURL,
		otelsql.WithAttributes(semconv.DBSystemPostgreSQL),
		otelsql.WithSpanOptions(otelsql.SpanOptions{Ping: false}),
	)
	if err != nil {
		log.Fatal("Failed to connect to game database:", err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal("Game database is unreachable:", err)
	}
	slog.Info("connected to game database")

	adapters, err := NewAdapters(db, staticBaseURL, jwtPublicKeyPEM)
	if err != nil {
		log.Fatal("Failed to initialize adapters:", err)
	}

	services := NewAppServices(adapters)
	commands := NewCommands(adapters, services)
	queries := NewQueries(adapters, services)

	WireCommandEvents(commands, adapters.Events)
	WireCommandSchedulerHandler(commands, adapters.Scheduler)

	consumer := rabbitmq.NewRabbitMQConsumer(rabbitURL)
	WireCommandIntegrationEvents(commands, consumer, authExchange, "game.auth.integration.events", billingExchange, "game.billing.integration.events")
	workers := NewWorkers(dbURL, services.Outbox, adapters.Scheduler, consumer)

	httpCommands := httpapi.Commands{
		User:        commands.User,
		BlackMarket: commands.BlackMarket,
		Base:        commands.Base,
		Building:    commands.Building,
		Army:        commands.Army,
		Tech:        commands.Tech,
		Storage:     commands.Storage,
		Operation:   commands.Operation,
		Trade:       commands.Trade,
		Alert:       commands.Alert,
		Diplomacy:   commands.Diplomacy,
	}
	httpQueries := httpapi.Queries{
		User:        queries.User,
		BlackMarket: queries.BlackMarket,
		Base:        queries.Base,
		Building:    queries.Building,
		Army:        queries.Army,
		Tech:        queries.Tech,
		Storage:     queries.Storage,
		Trade:       queries.Trade,
		Sector:      queries.Sector,
		Radar:       queries.Radar,
		Operation:   queries.Operation,
		Activity:    queries.Activity,
		Alert:       queries.Alert,
		Diplomacy:   queries.Diplomacy,
	}
	addr := fmt.Sprintf(":%s", port)
	router := httpapi.NewRouter(httpCommands, httpQueries, adapters.Tokens, adapters.Translator)
	httpServer := httpapi.NewServer(router, addr)

	grpcAddr := fmt.Sprintf(":%s", grpcPort)
	grpcCommands := grpcapi.Commands{ArmyPrototype: commands.Prototype}
	grpcQueries := grpcapi.Queries{ArmyPrototype: queries.Prototype}
	grpcRouter := grpcapi.NewRouter(grpcCommands, grpcQueries, grpcKey, adapters.Translator)
	grpcServer := grpcapi.NewServer(grpcRouter, grpcAddr)

	return &Module{
		Port:            port,
		DBURL:           dbURL,
		StaticBaseURL:   staticBaseURL,
		RabbitURL:       rabbitURL,
		AuthExchange:    authExchange,
		BillingExchange: billingExchange,
		DB:              db,
		Adapters:        adapters,
		Services:        services,
		Workers:         workers,
		Commands:        commands,
		Queries:         queries,
		HTTPServer:      httpServer,
		GRPCServer:      grpcServer,
	}
}

// Run connects the module's infrastructure, serves until the context is
// cancelled or a fatal error occurs, then drains and releases resources.
// Construction does no I/O; everything cancelable happens here.
func (m *Module) Run(ctx context.Context) error {
	slog.InfoContext(ctx, "starting game service", "port", m.Port)

	// Module-local cancel: a fatal HTTP error must stop the workers (and a
	// fatal worker error must stop the HTTP server), otherwise the drain
	// below would block forever.
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	if err := m.Adapters.Setup(ctx); err != nil {
		return fmt.Errorf("game adapters setup failed: %w", err)
	}

	if err := seedPeriodicJobs(ctx, m.DB); err != nil {
		return fmt.Errorf("failed to seed periodic jobs: %w", err)
	}

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

	grpcErr := make(chan error, 1)
	go func() {
		if err := m.GRPCServer.Start(); err != nil {
			grpcErr <- err
		}
	}()

	var runErr error
	select {
	case <-ctx.Done():
		slog.InfoContext(ctx, "game module: shutdown signal received, draining...")
	case err := <-httpErr:
		runErr = fmt.Errorf("game http server failed: %w", err)
		slog.ErrorContext(ctx, "game module: http server failed, draining...", "error", err)
		cancel()
	case err := <-grpcErr:
		runErr = fmt.Errorf("game grpc server failed: %w", err)
		slog.ErrorContext(ctx, "game module: grpc server failed, draining...", "error", err)
		cancel()
	case err := <-workerErr:
		runErr = fmt.Errorf("game workers failed: %w", err)
		slog.ErrorContext(ctx, "game module: worker failed, draining...", "error", err)
		cancel()
	}

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer shutdownCancel()
	if err := m.HTTPServer.Shutdown(shutdownCtx); err != nil {
		slog.ErrorContext(ctx, "game http server shutdown error", "error", err)
	}
	if err := m.GRPCServer.Shutdown(shutdownCtx); err != nil {
		slog.ErrorContext(ctx, "game grpc server shutdown error", "error", err)
	}

	if err := m.Workers.Wait(); err != nil && runErr == nil {
		runErr = fmt.Errorf("game workers failed: %w", err)
	}

	if err := m.DB.Close(); err != nil {
		slog.ErrorContext(ctx, "game database close error", "error", err)
	}

	slog.InfoContext(ctx, "game module: stopped")
	return runErr
}
