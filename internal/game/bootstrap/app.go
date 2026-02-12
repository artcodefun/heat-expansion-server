package bootstrap

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"time"

	"database/sql"

	_ "github.com/lib/pq"

	"github.com/artcodefun/heat-expansion-api/internal/game/infrastructure/db/repo"
	"github.com/artcodefun/heat-expansion-api/internal/game/infrastructure/jobs"
	httpapi "github.com/artcodefun/heat-expansion-api/internal/game/interfaces/http"
)

type App struct {
	// GameLogic *core.GameLogic
	Port          string
	DBURL         string
	JWTSecret     string
	LogLevel      string
	StaticBaseURL string
	AssetsDir     string

	DB         *sql.DB
	Adapters   *Adapters
	Services   *AppServices
	Commands   *Commands
	Queries    *Queries
	HTTPServer *httpapi.Server
}

func NewApp() *App {
	// Read environment variables
	port := os.Getenv("PORT")
	dbURL := os.Getenv("DB_URL")
	jwtSecret := os.Getenv("JWT_SECRET")
	staticBaseURL := os.Getenv("STATIC_BASE_URL")
	assetsDir := os.Getenv("ASSETS_DIR")

	if port == "" || dbURL == "" || jwtSecret == "" || staticBaseURL == "" {
		log.Fatal("Missing required environment variables. Please check your .env file.")
	}

	if assetsDir == "" {
		assetsDir = "./assets" // Default value
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal("Database is unreachable:", err)
	}
	fmt.Println("Connected to database successfully!")

	adapters, err := NewAdapters(db, jwtSecret, staticBaseURL)
	if err != nil {
		log.Fatal("Failed to initialize adapters:", err)
	}
	services := NewAppServices(adapters)
	commands := NewCommands(adapters, services)
	queries := NewQueries(adapters, services)
	// Wire event subscriptions and job dispatcher for commands
	WireCommandEvents(commands, adapters.Events)
	WireCommandSchedulerHandler(commands, adapters.Scheduler)

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
	router := httpapi.NewRouter(httpCommands, httpQueries, adapters.Tokens, assetsDir)
	httpServer := httpapi.NewServer(router)

	return &App{
		// GameLogic: gameLogic,
		Port:          port,
		DBURL:         dbURL,
		JWTSecret:     jwtSecret,
		StaticBaseURL: staticBaseURL,
		AssetsDir:     assetsDir,
		DB:            db,
		Adapters:      adapters,
		Services:      services,
		Commands:      commands,
		Queries:       queries,
		HTTPServer:    httpServer,
	}
}

func (a *App) Run() {
	fmt.Printf("Starting server on port %s\n", a.Port)
	fmt.Printf("Connecting to DB: %s\n", a.DBURL)
	fmt.Println("JWT secret configured")
	fmt.Printf("Static base URL: %s\n", a.StaticBaseURL)

	// Start scheduler loop
	if runner, ok := a.Adapters.Scheduler.(*jobs.DBScheduler); ok {
		go runner.Run(context.Background())
	}

	// Start outbox dispatcher loop
	go func() {
		ctx := context.Background()
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		var signalChan <-chan struct{}
		if txMgr, ok := a.Adapters.TxMgr.(*repo.DBTxManager); ok {
			signalChan = txMgr.CommitSignal
		}

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if err := a.Services.Outbox.ProcessBatch(100); err != nil {
					slog.Error("outbox dispatch failed (polling)", "error", err.Error())
				}
			case <-signalChan:
				if err := a.Services.Outbox.ProcessBatch(100); err != nil {
					slog.Error("outbox dispatch failed (signaled)", "error", err.Error())
				}
			}
		}
	}()

	// Seed initial world generation job (after short delay) if not already present.
	// This is now handled per-base via UserBaseCreatedEvent and historical migrations.

	addr := fmt.Sprintf(":%s", a.Port)
	fmt.Printf("Listening on %s\n", addr)
	if err := a.HTTPServer.Start(addr); err != nil {
		log.Fatalf("http server exited: %v", err)
	}
}
