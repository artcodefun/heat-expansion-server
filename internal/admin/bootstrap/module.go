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

	httpapi "github.com/artcodefun/heat-expansion-server/internal/admin/interfaces/http"
)

// Module wires and runs the admin service.
type Module struct {
	Port  string
	DBURL string

	DB         *sql.DB
	Adapters   *Adapters
	Commands   *Commands
	Queries    *Queries
	HTTPServer *httpapi.Server
}

func NewModule() *Module {
	port := os.Getenv("ADMIN_PORT")
	dbURL := os.Getenv("ADMIN_DB_URL")
	gameAddr := os.Getenv("GAME_GRPC_ADDR")
	gameKey := os.Getenv("GAME_GRPC_KEY")
	billingAddr := os.Getenv("BILLING_GRPC_ADDR")
	billingKey := os.Getenv("BILLING_GRPC_KEY")

	if port == "" || dbURL == "" || gameAddr == "" || gameKey == "" || billingAddr == "" || billingKey == "" {
		log.Fatal("Missing required environment variables: ADMIN_PORT, ADMIN_DB_URL, GAME_GRPC_ADDR, GAME_GRPC_KEY, BILLING_GRPC_ADDR, BILLING_GRPC_KEY")
	}

	db, err := otelsql.Open("postgres", dbURL,
		otelsql.WithAttributes(semconv.DBSystemPostgreSQL),
		otelsql.WithSpanOptions(otelsql.SpanOptions{Ping: false}),
	)
	if err != nil {
		log.Fatal("Failed to connect to admin database:", err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal("Admin database is unreachable:", err)
	}
	slog.Info("connected to admin database")

	adapters, err := NewAdapters(db, gameAddr, gameKey, billingAddr, billingKey)
	if err != nil {
		log.Fatal("Failed to initialise admin adapters:", err)
	}

	cmds := NewCommands(adapters)
	qrys := NewQueries(adapters)

	addr := fmt.Sprintf(":%s", port)
	router := httpapi.NewRouter(
		httpapi.Commands{
			Admin:       cmds.Admin,
			Prototype:   cmds.Prototype,
			Translation: cmds.Translation,
			Package:     cmds.Package,
		},
		httpapi.Queries{
			Admin:       qrys.Admin,
			Prototype:   qrys.Prototype,
			Translation: qrys.Translation,
			Package:     qrys.Package,
		},
		adapters.SessionValidator,
		adapters.Translator,
	)
	httpServer := httpapi.NewServer(router, addr)

	return &Module{
		Port:       port,
		DBURL:      dbURL,
		DB:         db,
		Adapters:   adapters,
		Commands:   cmds,
		Queries:    qrys,
		HTTPServer: httpServer,
	}
}

// Run serves until the context is cancelled or a fatal error occurs, then
// drains and releases resources. Construction does no I/O; all cancelable
// infrastructure starts here.
func (m *Module) Run(ctx context.Context) error {
	slog.InfoContext(ctx, "starting admin service", "port", m.Port)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	httpErr := make(chan error, 1)
	go func() {
		if err := m.HTTPServer.Start(); err != nil {
			httpErr <- err
		}
	}()

	var runErr error
	select {
	case <-ctx.Done():
		slog.InfoContext(ctx, "admin module: shutdown signal received, draining...")
	case err := <-httpErr:
		runErr = fmt.Errorf("admin http server failed: %w", err)
		slog.ErrorContext(ctx, "admin module: http server failed", "error", err)
		cancel()
	}

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer shutdownCancel()
	if err := m.HTTPServer.Shutdown(shutdownCtx); err != nil {
		slog.ErrorContext(ctx, "admin http server shutdown error", "error", err)
	}

	if err := m.DB.Close(); err != nil {
		slog.ErrorContext(ctx, "admin database close error", "error", err)
	}

	slog.InfoContext(ctx, "admin module: stopped")
	return runErr
}
