package repo

import (
	"context"
	"database/sql"
	"sort"
	"time"

	"github.com/artcodefun/heat-expansion-server/internal/game/application/cqrs/readmodels"
	"github.com/artcodefun/heat-expansion-server/internal/game/infrastructure/readstore/gen"
	"github.com/artcodefun/heat-expansion-server/internal/game/infrastructure/readstore/mappers"
)

type BlackMarketReadRepo struct{ q *gen.Queries }

func NewBlackMarketReadRepo(q *gen.Queries) *BlackMarketReadRepo { return &BlackMarketReadRepo{q: q} }

func blackMarketNowParam(now int64) sql.NullInt64 {
	return sql.NullInt64{Int64: now, Valid: true}
}

func blackMarketLimitedParam(limited *bool) sql.NullBool {
	if limited == nil {
		return sql.NullBool{}
	}
	return sql.NullBool{Bool: *limited, Valid: true}
}

func (r *BlackMarketReadRepo) ListActiveOffers(ctx context.Context, kind *readmodels.BlackMarketOfferKind, limited *bool) ([]*readmodels.BlackMarketOffer, error) {
	now := time.Now().Unix()
	if kind != nil {
		switch *kind {
		case readmodels.BlackMarketOfferKindBuilding:
			return r.listBuildingOffers(ctx, now, limited)
		case readmodels.BlackMarketOfferKindArmy:
			return r.listArmyOffers(ctx, now, limited)
		case readmodels.BlackMarketOfferKindStorage:
			return r.listStorageOffers(ctx, now, limited)
		default:
			return []*readmodels.BlackMarketOffer{}, nil
		}
	}

	buildings, err := r.q.ListActiveBlackMarketBuildingOffers(ctx, gen.ListActiveBlackMarketBuildingOffersParams{
		Limited: blackMarketLimitedParam(limited),
		Now:     blackMarketNowParam(now),
	})
	if err != nil {
		return nil, err
	}
	armies, err := r.q.ListActiveBlackMarketArmyOffers(ctx, gen.ListActiveBlackMarketArmyOffersParams{
		Limited: blackMarketLimitedParam(limited),
		Now:     blackMarketNowParam(now),
	})
	if err != nil {
		return nil, err
	}
	storages, err := r.q.ListActiveBlackMarketStorageOffers(ctx, gen.ListActiveBlackMarketStorageOffersParams{
		Limited: blackMarketLimitedParam(limited),
		Now:     blackMarketNowParam(now),
	})
	if err != nil {
		return nil, err
	}

	merged := make([]*readmodels.BlackMarketOffer, 0, len(buildings)+len(armies)+len(storages))
	for _, row := range buildings {
		item := mappers.BlackMarketBuildingOfferFromRow(row)
		merged = append(merged, &item)
	}
	for _, row := range armies {
		item := mappers.BlackMarketArmyOfferFromRow(row)
		merged = append(merged, &item)
	}
	for _, row := range storages {
		item := mappers.BlackMarketStorageOfferFromRow(row)
		merged = append(merged, &item)
	}

	sort.SliceStable(merged, func(i, j int) bool {
		if merged[i].Priority == merged[j].Priority {
			return merged[i].ID < merged[j].ID
		}
		return merged[i].Priority > merged[j].Priority
	})

	out := make([]*readmodels.BlackMarketOffer, 0, len(merged))
	for _, item := range merged {
		out = append(out, item)
	}
	return out, nil
}

func (r *BlackMarketReadRepo) listBuildingOffers(ctx context.Context, now int64, limited *bool) ([]*readmodels.BlackMarketOffer, error) {
	rows, err := r.q.ListActiveBlackMarketBuildingOffers(ctx, gen.ListActiveBlackMarketBuildingOffersParams{
		Limited: blackMarketLimitedParam(limited),
		Now:     blackMarketNowParam(now),
	})
	if err != nil {
		return nil, err
	}
	out := make([]*readmodels.BlackMarketOffer, 0, len(rows))
	for _, row := range rows {
		item := mappers.BlackMarketBuildingOfferFromRow(row)
		out = append(out, &item)
	}
	return out, nil
}

func (r *BlackMarketReadRepo) listArmyOffers(ctx context.Context, now int64, limited *bool) ([]*readmodels.BlackMarketOffer, error) {
	rows, err := r.q.ListActiveBlackMarketArmyOffers(ctx, gen.ListActiveBlackMarketArmyOffersParams{
		Limited: blackMarketLimitedParam(limited),
		Now:     blackMarketNowParam(now),
	})
	if err != nil {
		return nil, err
	}
	out := make([]*readmodels.BlackMarketOffer, 0, len(rows))
	for _, row := range rows {
		item := mappers.BlackMarketArmyOfferFromRow(row)
		out = append(out, &item)
	}
	return out, nil
}

func (r *BlackMarketReadRepo) listStorageOffers(ctx context.Context, now int64, limited *bool) ([]*readmodels.BlackMarketOffer, error) {
	rows, err := r.q.ListActiveBlackMarketStorageOffers(ctx, gen.ListActiveBlackMarketStorageOffersParams{
		Limited: blackMarketLimitedParam(limited),
		Now:     blackMarketNowParam(now),
	})
	if err != nil {
		return nil, err
	}
	out := make([]*readmodels.BlackMarketOffer, 0, len(rows))
	for _, row := range rows {
		item := mappers.BlackMarketStorageOfferFromRow(row)
		out = append(out, &item)
	}
	return out, nil
}
