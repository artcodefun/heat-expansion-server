package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

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

	authModule := authBootstrap.NewModule()
	billingModule := billingBootstrap.NewModule()
	gameModule := gameBootstrap.NewModule()

	runModules(ctx, authModule, billingModule, gameModule)
	slog.Info("server stopped cleanly")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := otelShutdown(shutdownCtx); err != nil {
		slog.Error("OpenTelemetry shutdown error", "error", err)
	}
}

type Module interface {
	Run(ctx context.Context)
}

func runModules(ctx context.Context, modules ...Module) {
	var wg sync.WaitGroup
	for _, m := range modules {
		wg.Go(func() { m.Run(ctx) })
	}
	wg.Wait()
}
