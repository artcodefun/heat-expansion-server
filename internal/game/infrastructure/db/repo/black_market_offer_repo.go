package repo

import (
	"context"
	"database/sql"
	"errors"

	"github.com/artcodefun/heat-expansion-server/internal/game/application/ports"
	"github.com/artcodefun/heat-expansion-server/internal/game/domain"
	"github.com/artcodefun/heat-expansion-server/internal/game/infrastructure/db/gen"
	"github.com/artcodefun/heat-expansion-server/internal/game/infrastructure/db/mappers"
)

type BlackMarketOfferRepo struct {
	q *gen.Queries
}

func NewBlackMarketOfferRepo(q *gen.Queries) *BlackMarketOfferRepo {
	return &BlackMarketOfferRepo{q: q}
}

func (r *BlackMarketOfferRepo) Tx(tx ports.Transaction) ports.BlackMarketOfferRepository {
	if sqlTx, ok := tx.(*sql.Tx); ok {
		return &BlackMarketOfferRepo{q: r.q.WithTx(sqlTx)}
	}
	return r
}

func (r *BlackMarketOfferRepo) Create(ctx context.Context, offer *domain.BlackMarketOffer) error {
	row, err := r.q.InsertBlackMarketOffer(ctx, mappers.InsertBlackMarketOfferParamsFromDomain(offer))
	if err != nil {
		return err
	}
	created := mappers.BlackMarketOfferFromDB(row)
	offer.ID = created.ID
	return nil
}

func (r *BlackMarketOfferRepo) Update(ctx context.Context, offer *domain.BlackMarketOffer) error {
	_, err := r.q.UpdateBlackMarketOffer(ctx, mappers.UpdateBlackMarketOfferParamsFromDomain(offer))
	return err
}

func (r *BlackMarketOfferRepo) FindByID(ctx context.Context, id int64) (*domain.BlackMarketOffer, error) {
	row, err := r.q.GetBlackMarketOfferByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ports.ErrNotFound
		}
		return nil, err
	}
	return mappers.BlackMarketOfferFromDB(row), nil
}

func (r *BlackMarketOfferRepo) FindByIDForUpdate(ctx context.Context, id int64) (*domain.BlackMarketOffer, error) {
	row, err := r.q.GetBlackMarketOfferByIDForUpdate(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ports.ErrNotFound
		}
		return nil, err
	}
	return mappers.BlackMarketOfferFromDB(row), nil
}

func (r *BlackMarketOfferRepo) ListActiveLimitedOffers(ctx context.Context, now int64) ([]*domain.BlackMarketOffer, error) {
	rows, err := r.q.ListActiveLimitedBlackMarketOffers(ctx, sql.NullInt64{Int64: now, Valid: true})
	if err != nil {
		return nil, err
	}
	return mappers.BlackMarketOffersFromDB(rows), nil
}

func (r *BlackMarketOfferRepo) ListExpiredLimitedOffers(ctx context.Context, now int64) ([]*domain.BlackMarketOffer, error) {
	rows, err := r.q.ListExpiredLimitedBlackMarketOffers(ctx, sql.NullInt64{Int64: now, Valid: true})
	if err != nil {
		return nil, err
	}
	return mappers.BlackMarketOffersFromDB(rows), nil
}
