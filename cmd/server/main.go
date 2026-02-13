package main

import (
	authBootstrap "github.com/artcodefun/heat-expansion-server/internal/auth/bootstrap"
	gameBootstrap "github.com/artcodefun/heat-expansion-server/internal/game/bootstrap"
)

func main() {
	authModule := authBootstrap.NewModule()
	go authModule.Run()

	gameModule := gameBootstrap.NewModule()
	gameModule.Run()
}
