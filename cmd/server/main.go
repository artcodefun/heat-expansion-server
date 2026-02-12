package main

import "github.com/artcodefun/heat-expansion-api/internal/game/bootstrap"

func main() {
	module := bootstrap.NewModule()
	module.Run()
}
