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
	gameModule := gameBootstrap.NewModule()

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		authModule.Run(ctx)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		gameModule.Run(ctx)
	}()

	wg.Wait()
	slog.Info("server stopped cleanly")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := otelShutdown(shutdownCtx); err != nil {
		slog.Error("OpenTelemetry shutdown error", "error", err)
	}
}
