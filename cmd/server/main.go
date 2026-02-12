package main

import "github.com/artcodefun/heat-expansion-server/internal/game/bootstrap"

func main() {
	module := bootstrap.NewModule()
	module.Run()
}
