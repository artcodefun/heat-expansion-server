package queries

import (
	"context"

	"github.com/artcodefun/heat-expansion-server/internal/game/application/cqrs"
	"github.com/artcodefun/heat-expansion-server/internal/game/application/cqrs/readmodels"
	"github.com/artcodefun/heat-expansion-server/internal/game/application/services"
	"github.com/artcodefun/heat-expansion-server/internal/game/domain"
)

type BlackMarketQueries struct {
	Access *services.AccessControlService
}

func NewBlackMarketQueries(access *services.AccessControlService) *BlackMarketQueries {
	return &BlackMarketQueries{Access: access}
}

func (q *BlackMarketQueries) ListResourceRates(ctx context.Context, actor cqrs.Actor, baseID int) ([]*readmodels.BlackMarketResourceRate, error) {
	if err := q.Access.EnsureBaseOwnership(ctx, actor.UserID, baseID); err != nil {
		return nil, err
	}

	rates := domain.ListBlackMarketResourceRates()
	return readmodels.BlackMarketResourceRateListFromDomain(rates), nil
}
