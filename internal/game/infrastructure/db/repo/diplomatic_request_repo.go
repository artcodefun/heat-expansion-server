package repo

import (
	"context"
	"database/sql"
	"errors"

	"github.com/artcodefun/heat-expansion-server/internal/game/application/ports"
	"github.com/artcodefun/heat-expansion-server/internal/game/domain"
	"github.com/artcodefun/heat-expansion-server/internal/game/infrastructure/db/gen"
	"github.com/artcodefun/heat-expansion-server/internal/game/infrastructure/db/mappers"
	"github.com/google/uuid"
)

type DiplomaticRequestRepo struct {
	q *gen.Queries
}

func NewDiplomaticRequestRepo(q *gen.Queries) *DiplomaticRequestRepo {
	return &DiplomaticRequestRepo{q: q}
}

func (r *DiplomaticRequestRepo) Tx(tx ports.Transaction) ports.DiplomaticRequestRepository {
	if sqlTx, ok := tx.(*sql.Tx); ok {
		return &DiplomaticRequestRepo{q: r.q.WithTx(sqlTx)}
	}
	return r
}

func (r *DiplomaticRequestRepo) Create(ctx context.Context, request *domain.DiplomaticRequest) error {
	return r.q.InsertDiplomaticRequest(ctx, mappers.InsertDiplomaticRequestParamsFromDomain(request))
}

func (r *DiplomaticRequestRepo) Update(ctx context.Context, request *domain.DiplomaticRequest) error {
	return r.q.UpdateDiplomaticRequest(ctx, mappers.UpdateDiplomaticRequestParamsFromDomain(request))
}

func (r *DiplomaticRequestRepo) FindByID(ctx context.Context, id uuid.UUID) (*domain.DiplomaticRequest, error) {
	row, err := r.q.GetDiplomaticRequest(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ports.ErrNotFound
		}
		return nil, err
	}
	return mappers.DiplomaticRequestFromDB(row), nil
}

func (r *DiplomaticRequestRepo) FindByIDForUpdate(ctx context.Context, id uuid.UUID) (*domain.DiplomaticRequest, error) {
	row, err := r.q.GetDiplomaticRequestForUpdate(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ports.ErrNotFound
		}
		return nil, err
	}
	return mappers.DiplomaticRequestFromDB(row), nil
}

func (r *DiplomaticRequestRepo) ExistsPendingByKind(ctx context.Context, userAID, userBID uuid.UUID, kind domain.DiplomaticRequestKind) (bool, error) {
	return r.q.ExistsPendingDiplomaticRequestByKind(ctx, gen.ExistsPendingDiplomaticRequestByKindParams{
		UserAID: userAID,
		UserBID: userBID,
		Kind:    string(kind),
	})
}
