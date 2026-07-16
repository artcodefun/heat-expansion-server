package commands

import (
	"context"

	"github.com/artcodefun/heat-expansion-server/internal/admin/application/cqrs"
	"github.com/artcodefun/heat-expansion-server/internal/admin/application/ports"
)

// TranslationCommands delegates translation writes to the game private gRPC API.
type TranslationCommands struct {
	game ports.GamePrivateClient
}

func NewTranslationCommands(game ports.GamePrivateClient) *TranslationCommands {
	return &TranslationCommands{game: game}
}

func (c *TranslationCommands) UpsertTranslation(ctx context.Context, _ cqrs.Actor, locale, key, value string) (*ports.Translation, error) {
	created, err := c.game.UpsertTranslation(ctx, locale, key, value)
	return created, clientErr(err)
}
