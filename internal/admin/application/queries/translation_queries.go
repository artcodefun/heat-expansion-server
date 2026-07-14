package queries

import (
	"context"

	"github.com/artcodefun/heat-expansion-server/internal/admin/application/cqrs"
	"github.com/artcodefun/heat-expansion-server/internal/admin/application/cqrs/readmodels"
	"github.com/artcodefun/heat-expansion-server/internal/admin/application/ports"
)

// TranslationQueries delegates translation reads to the game private gRPC API.
type TranslationQueries struct {
	game ports.GamePrivateClient
}

func NewTranslationQueries(game ports.GamePrivateClient) *TranslationQueries {
	return &TranslationQueries{game: game}
}

func (q *TranslationQueries) ListTranslations(ctx context.Context, _ cqrs.Actor) ([]*readmodels.Translation, error) {
	list, err := q.game.ListTranslations(ctx)
	return list, clientErr(err)
}
