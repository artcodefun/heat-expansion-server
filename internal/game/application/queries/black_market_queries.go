package queries

import (
	"context"

	"github.com/artcodefun/heat-expansion-server/internal/game/application/cqrs"
	"github.com/artcodefun/heat-expansion-server/internal/game/application/cqrs/readmodels"
	"github.com/artcodefun/heat-expansion-server/internal/game/application/ports"
	"github.com/artcodefun/heat-expansion-server/internal/game/application/services"
	"github.com/artcodefun/heat-expansion-server/internal/game/domain"
)

type BlackMarketQueries struct {
	Offers ports.BlackMarketReadRepository
	Access *services.AccessControlService
}

func NewBlackMarketQueries(
	offers ports.BlackMarketReadRepository,
	access *services.AccessControlService,
) *BlackMarketQueries {
	return &BlackMarketQueries{
		Offers: offers,
		Access: access,
	}
}

func (q *BlackMarketQueries) ListResourceRates(ctx context.Context, actor cqrs.Actor, baseID int) ([]*readmodels.BlackMarketResourceRate, error) {
	if err := q.Access.EnsureBaseOwnership(ctx, actor.UserID, baseID); err != nil {
		return nil, err
	}

	rates := domain.ListBlackMarketResourceRates()
	return readmodels.BlackMarketResourceRateListFromDomain(rates), nil
}

func (q *BlackMarketQueries) ListActiveOffers(ctx context.Context, actor cqrs.Actor, baseID int, kind *domain.BlackMarketOfferKind, limited *bool) ([]*readmodels.BlackMarketOffer, error) {
	if err := q.Access.EnsureBaseOwnership(ctx, actor.UserID, baseID); err != nil {
		return nil, err
	}
	offers, err := q.Offers.ListActiveOffers(ctx, blackMarketOfferKindToReadModel(kind), limited)
	if err != nil {
		return nil, repoErr(err)
	}
	return offers, nil
}

func blackMarketOfferKindToReadModel(kind *domain.BlackMarketOfferKind) *readmodels.BlackMarketOfferKind {
	if kind == nil {
		return nil
	}
	converted := readmodels.BlackMarketOfferKind(*kind)
	return &converted
}
