package bootstrap

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"os"
	"time"

	_ "github.com/lib/pq"

	"github.com/artcodefun/heat-expansion-server/internal/game/infrastructure/db/repo"
	"github.com/artcodefun/heat-expansion-server/internal/game/infrastructure/events"
	"github.com/artcodefun/heat-expansion-server/internal/game/infrastructure/jobs"
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

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal("Database is unreachable:", err)
	}
	fmt.Println("Connected to database successfully!")

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

	httpCommands := httpapi.Commands{
		User:      commands.User,
		Base:      commands.Base,
		Building:  commands.Building,
		Army:      commands.Army,
		Tech:      commands.Tech,
		Storage:   commands.Storage,
		Operation: commands.Operation,
		Alert:     commands.Alert,
	}
	httpQueries := httpapi.Queries{
		User:      queries.User,
		Base:      queries.Base,
		Building:  queries.Building,
		Army:      queries.Army,
		Tech:      queries.Tech,
		Storage:   queries.Storage,
		Sector:    queries.Sector,
		Radar:     queries.Radar,
		Operation: queries.Operation,
		Activity:  queries.Activity,
		Alert:     queries.Alert,
	}
	router := httpapi.NewRouter(httpCommands, httpQueries, adapters.Tokens, adapters.Translator)
	httpServer := httpapi.NewServer(router)

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
		Commands:      commands,
		Queries:       queries,
		HTTPServer:    httpServer,
		Consumer:      consumer,
	}
}

func (m *Module) Run() {
	fmt.Printf("Starting server on port %s\n", m.Port)
	fmt.Printf("Connecting to DB: %s\n", m.DBURL)
	fmt.Println("JWT secret configured")
	fmt.Printf("Static base URL: %s\n", m.StaticBaseURL)
	fmt.Printf("RabbitMQ URL: %s\n", m.RabbitURL)
	fmt.Printf("Auth Exchange: %s\n", m.AuthExchange)

	if runner, ok := m.Adapters.Scheduler.(*jobs.DBScheduler); ok {
		go runner.Run(context.Background())
	}

	if err := m.Consumer.Start(context.Background()); err != nil {
		log.Printf("Warning: Failed to start RabbitMQ consumer: %v", err)
	}

	go func() {
		ctx := context.Background()
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		var signalChan <-chan struct{}
		if txMgr, ok := m.Adapters.TxMgr.(*repo.DBTxManager); ok {
			signalChan = txMgr.CommitSignal
		}

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if err := m.Services.Outbox.ProcessBatch(100); err != nil {
					slog.Error("outbox dispatch failed (polling)", "error", err.Error())
				}
			case <-signalChan:
				if err := m.Services.Outbox.ProcessBatch(100); err != nil {
					slog.Error("outbox dispatch failed (signaled)", "error", err.Error())
				}
			}
		}
	}()

	addr := fmt.Sprintf(":%s", m.Port)
	fmt.Printf("Listening on %s\n", addr)
	if err := m.HTTPServer.Start(addr); err != nil {
		log.Fatalf("http server exited: %v", err)
	}
}
