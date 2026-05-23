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

	"github.com/artcodefun/heat-expansion-server/internal/game/infrastructure/events"
	httpapi "github.com/artcodefun/heat-expansion-server/internal/game/interfaces/http"
)

type Module struct {
	Port          string
	DBURL         string
	JWTSecret     string
	LogLevel      string
	StaticBaseURL string
	RabbitURL     string
	AuthExchange  string

	DB         *sql.DB
	Adapters   *Adapters
	Services   *AppServices
	Workers    *Workers
	Commands   *Commands
	Queries    *Queries
	HTTPServer *httpapi.Server
	Consumer   *events.RabbitMQConsumer
}

func NewModule() *Module {
	port := os.Getenv("GAME_PORT")
	dbURL := os.Getenv("GAME_DB_URL")
	jwtSecret := os.Getenv("GAME_JWT_SECRET")
	staticBaseURL := os.Getenv("GAME_STATIC_BASE_URL")
	i18nPath := os.Getenv("GAME_I18N_PATH")
	rabbitURL := os.Getenv("RABBITMQ_URL")
	authExchange := os.Getenv("AUTH_INTEGRATION_EXCHANGE")

	if port == "" || dbURL == "" || jwtSecret == "" || staticBaseURL == "" || rabbitURL == "" || authExchange == "" {
		log.Fatal("Missing required environment variables. Please check your .env file.")
	}

	db, err := otelsql.Open("postgres", dbURL,
		otelsql.WithAttributes(semconv.DBSystemPostgreSQL),
		otelsql.WithSpanOptions(otelsql.SpanOptions{Ping: false}),
	)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal("Database is unreachable:", err)
	}
	slog.Info("connected to game database")

	adapters, err := NewAdapters(db, staticBaseURL, jwtSecret, i18nPath)
	if err != nil {
		log.Fatal("Failed to initialize adapters:", err)
	}

	services := NewAppServices(adapters)
	commands := NewCommands(adapters, services)
	queries := NewQueries(adapters, services)

	WireCommandEvents(commands, adapters.Events)
	WireCommandSchedulerHandler(commands, adapters.Scheduler)

	consumer := events.NewRabbitMQConsumer(rabbitURL)
	WireCommandIntegrationEvents(commands, consumer, authExchange, "game.auth.integration.events")
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

	return &Module{
		Port:          port,
		DBURL:         dbURL,
		JWTSecret:     jwtSecret,
		StaticBaseURL: staticBaseURL,
		RabbitURL:     rabbitURL,
		AuthExchange:  authExchange,
		DB:            db,
		Adapters:      adapters,
		Services:      services,
		Workers:       workers,
		Commands:      commands,
		Queries:       queries,
		HTTPServer:    httpServer,
		Consumer:      consumer,
	}
}

func (m *Module) Run(ctx context.Context) {
	slog.Info("starting game service", "port", m.Port)

	if err := seedPeriodicJobs(ctx, m.DB); err != nil {
		log.Fatal("Failed to seed periodic jobs:", err)
	}

	m.Workers.Start(ctx)

	go func() {
		if err := m.HTTPServer.Start(); err != nil {
			slog.Error("game http server stopped", "error", err)
		}
	}()

	<-ctx.Done()
	slog.Info("game module: shutdown signal received, draining...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := m.HTTPServer.Shutdown(shutdownCtx); err != nil {
		slog.Error("game http server shutdown error", "error", err)
	}

	m.Workers.Wait()
	slog.Info("game module: stopped")
}
