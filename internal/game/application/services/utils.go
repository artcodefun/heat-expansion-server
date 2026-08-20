package services

import (
	"errors"

	"github.com/artcodefun/heat-expansion-server/internal/game/application"
	"github.com/artcodefun/heat-expansion-server/internal/game/application/ports"
)

func repoErr(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, ports.ErrNotFound) {
		return application.ErrNotFound
	}
	return err
}
