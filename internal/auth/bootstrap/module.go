package bootstrap

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/artcodefun/heat-expansion-server/internal/auth/infrastructure/db/repo"
	httpapi "github.com/artcodefun/heat-expansion-server/internal/auth/interfaces/http"
	_ "github.com/lib/pq"
)

type Module struct {
	Port     string
	DB       *sql.DB
	Adapters *Adapters
	Services *AppServices
	Commands *Commands
	Queries  *Queries
}

func NewModule() *Module {
	port := os.Getenv("AUTH_PORT")
	dbURL := os.Getenv("AUTH_DB_URL")
	jwtSecret := os.Getenv("AUTH_JWT_SECRET")

	if port == "" || dbURL == "" || jwtSecret == "" {
		log.Fatal("Missing required auth environment variables (AUTH_PORT, AUTH_DB_URL, AUTH_JWT_SECRET)")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Failed to connect to auth database:", err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal("Auth database is unreachable:", err)
	}
	fmt.Println("Connected to auth database successfully!")

	adapters, err := NewAdapters(db, jwtSecret)
	if err != nil {
		log.Fatal("Failed to initialize auth adapters:", err)
	}

	services := NewAppServices(adapters)
	commands := NewCommands(adapters)
	queries := NewQueries(adapters)

	WireIntegrationEvents(services, adapters.Events)

	return &Module{
		Port:     port,
		DB:       db,
		Adapters: adapters,
		Services: services,
		Commands: commands,
		Queries:  queries,
	}
}

func (m *Module) Run() {
	fmt.Printf("Starting auth service on port %s\n", m.Port)
	ctx := context.Background()

	// Outbox loop
	go func() {
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
					slog.Error("auth outbox dispatch failed (polling)", "error", err.Error())
				}
			case <-signalChan:
				if err := m.Services.Outbox.ProcessBatch(100); err != nil {
					slog.Error("auth outbox dispatch failed (signaled)", "error", err.Error())
				}
			}
		}
	}()

	// Integration Outbox loop
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if err := m.Services.IntegrationOutbox.ProcessBatch(100); err != nil {
					slog.Error("auth integration outbox dispatch failed", "error", err.Error())
				}
			}
		}
	}()

	// HTTP Server
	router := httpapi.NewRouter(
		httpapi.Commands{Account: m.Commands.Account},
		httpapi.Queries{Account: m.Queries.Account},
	)
	server := httpapi.NewServer(router)

	addr := fmt.Sprintf(":%s", m.Port)
	if err := server.Start(addr); err != nil {
		log.Fatalf("auth http server exited: %v", err)
	}
}
