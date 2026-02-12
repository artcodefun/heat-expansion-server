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
	"github.com/artcodefun/heat-expansion-server/internal/game/infrastructure/jobs"
	httpapi "github.com/artcodefun/heat-expansion-server/internal/game/interfaces/http"
)

type Module struct {
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

func NewModule() *Module {
	port := os.Getenv("GAME_PORT")
	dbURL := os.Getenv("GAME_DB_URL")
	jwtSecret := os.Getenv("GAME_JWT_SECRET")
	staticBaseURL := os.Getenv("GAME_STATIC_BASE_URL")
	assetsDir := os.Getenv("GAME_ASSETS_DIR")

	if port == "" || dbURL == "" || jwtSecret == "" || staticBaseURL == "" {
		log.Fatal("Missing required environment variables. Please check your .env file.")
	}

	if assetsDir == "" {
		assetsDir = "./assets"
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

	return &Module{
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

func (m *Module) Run() {
	fmt.Printf("Starting server on port %s\n", m.Port)
	fmt.Printf("Connecting to DB: %s\n", m.DBURL)
	fmt.Println("JWT secret configured")
	fmt.Printf("Static base URL: %s\n", m.StaticBaseURL)

	if runner, ok := m.Adapters.Scheduler.(*jobs.DBScheduler); ok {
		go runner.Run(context.Background())
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
