package bootstrap

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	httpapi "github.com/artcodefun/heat-expansion-server/internal/auth/interfaces/http"
	_ "github.com/lib/pq"
)

type Module struct {
	Port  string
	DBURL string

	DB       *sql.DB
	Adapters *Adapters
	Services *AppServices
	Workers  *Workers
	Commands *Commands
	Queries  *Queries
}

func NewModule() *Module {
	rabbitURL := os.Getenv("RABBITMQ_URL")
	port := os.Getenv("AUTH_PORT")
	dbURL := os.Getenv("AUTH_DB_URL")
	jwtSecret := os.Getenv("AUTH_JWT_SECRET")
	intExchange := os.Getenv("AUTH_INTEGRATION_EXCHANGE")

	if port == "" || dbURL == "" || jwtSecret == "" || rabbitURL == "" || intExchange == "" {
		log.Fatal("Missing required auth environment variables (AUTH_PORT, AUTH_DB_URL, AUTH_JWT_SECRET, RABBITMQ_URL, AUTH_INTEGRATION_EXCHANGE)")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Failed to connect to auth database:", err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal("Auth database is unreachable:", err)
	}
	fmt.Println("Connected to auth database successfully!")

	adapters, err := NewAdapters(db, jwtSecret, rabbitURL, intExchange)
	if err != nil {
		log.Fatal("Failed to initialize auth adapters:", err)
	}

	services := NewAppServices(adapters)
	workers := NewWorkers(dbURL, services.Outbox, services.IntegrationOutbox)
	commands := NewCommands(adapters)
	queries := NewQueries(adapters)

	WireIntegrationEvents(services, adapters.Events)

	return &Module{
		Port:     port,
		DBURL:    dbURL,
		DB:       db,
		Adapters: adapters,
		Services: services,
		Workers:  workers,
		Commands: commands,
		Queries:  queries,
	}
}

func (m *Module) Run() {
	fmt.Printf("Starting auth service on port %s\n", m.Port)
	ctx := context.Background()

	// Start background workers
	go m.Workers.DomainOutboxLoop(ctx)
	go m.Workers.IntegrationOutboxLoop(ctx)

	// HTTP Server
	router := httpapi.NewRouter(
		httpapi.Commands{Account: m.Commands.Account},
		httpapi.Queries{Account: m.Queries.Account},
		m.Adapters.Translator,
	)
	server := httpapi.NewServer(router)

	addr := fmt.Sprintf(":%s", m.Port)
	if err := server.Start(addr); err != nil {
		log.Fatalf("auth http server exited: %v", err)
	}
}
