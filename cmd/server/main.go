package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"

	authBootstrap "github.com/artcodefun/heat-expansion-server/internal/auth/bootstrap"
	billingBootstrap "github.com/artcodefun/heat-expansion-server/internal/billing/bootstrap"
	gameBootstrap "github.com/artcodefun/heat-expansion-server/internal/game/bootstrap"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	otelShutdown, err := initTelemetry(ctx)
	if err != nil {
		slog.Error("failed to initialize OpenTelemetry", "error", err)
		os.Exit(1)
	}

	// Construction is pure wiring: env validation and dependency graphs only.
	// All infrastructure I/O (DB pings, broker connections, content loads)
	// happens inside each module's Run under the cancellable signal context.
	authModule := authBootstrap.NewModule()
	billingModule := billingBootstrap.NewModule()
	gameModule := gameBootstrap.NewModule()

	runErr := runModules(ctx, authModule, billingModule, gameModule)
	if runErr != nil {
		slog.ErrorContext(ctx, "server stopped with error", "error", runErr)
	} else {
		slog.InfoContext(ctx, "server stopped cleanly")
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := otelShutdown(shutdownCtx); err != nil {
		slog.Error("OpenTelemetry shutdown error", "error", err)
	}

	if runErr != nil {
		os.Exit(1)
	}
}

type Module interface {
	Run(ctx context.Context) error
}

// runModules runs every module until the shared context is cancelled or a
// module fails. A single failure cancels the group context so the remaining
// modules drain gracefully; the first error is returned.
func runModules(ctx context.Context, modules ...Module) error {
	g, ctx := errgroup.WithContext(ctx)
	for _, m := range modules {
		g.Go(func() error { return m.Run(ctx) })
	}
	return g.Wait()
}
