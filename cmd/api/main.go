package main

import (
	"github.com/artcodefun/heat-expansion-api/internal/game/bootstrap"
)

func main() {
	app := bootstrap.NewApp()
	// createTestPrototypes(app)
	app.Run()
}
