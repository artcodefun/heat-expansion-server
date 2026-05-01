package main

import (
	"context"
	"log/slog"
	"os/signal"
	"sync"
	"syscall"

	authBootstrap "github.com/artcodefun/heat-expansion-server/internal/auth/bootstrap"
	gameBootstrap "github.com/artcodefun/heat-expansion-server/internal/game/bootstrap"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

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
}
