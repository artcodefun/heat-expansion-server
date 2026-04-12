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

type DiplomaticRelationshipRepo struct {
	q *gen.Queries
}

func NewDiplomaticRelationshipRepo(q *gen.Queries) *DiplomaticRelationshipRepo {
	return &DiplomaticRelationshipRepo{q: q}
}

func (r *DiplomaticRelationshipRepo) Tx(tx ports.Transaction) ports.DiplomaticRelationshipRepository {
	if sqlTx, ok := tx.(*sql.Tx); ok {
		return &DiplomaticRelationshipRepo{q: r.q.WithTx(sqlTx)}
	}
	return r
}

func (r *DiplomaticRelationshipRepo) Create(ctx context.Context, relationship *domain.DiplomaticRelationship) error {
	return r.q.InsertDiplomaticRelationship(ctx, mappers.InsertDiplomaticRelationshipParamsFromDomain(relationship))
}

func (r *DiplomaticRelationshipRepo) Update(ctx context.Context, relationship *domain.DiplomaticRelationship) error {
	return r.q.UpdateDiplomaticRelationship(ctx, mappers.UpdateDiplomaticRelationshipParamsFromDomain(relationship))
}

func (r *DiplomaticRelationshipRepo) FindBetweenUsers(ctx context.Context, userAID, userBID uuid.UUID) (*domain.DiplomaticRelationship, error) {
	a, b := canonicalDiplomaticPair(userAID, userBID)
	row, err := r.q.GetDiplomaticRelationship(ctx, gen.GetDiplomaticRelationshipParams{UserAID: a, UserBID: b})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ports.ErrNotFound
		}
		return nil, err
	}
	return mappers.DiplomaticRelationshipFromDB(row), nil
}

func (r *DiplomaticRelationshipRepo) FindBetweenUsersForUpdate(ctx context.Context, userAID, userBID uuid.UUID) (*domain.DiplomaticRelationship, error) {
	a, b := canonicalDiplomaticPair(userAID, userBID)
	row, err := r.q.GetDiplomaticRelationshipForUpdate(ctx, gen.GetDiplomaticRelationshipForUpdateParams{UserAID: a, UserBID: b})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ports.ErrNotFound
		}
		return nil, err
	}
	return mappers.DiplomaticRelationshipFromDB(row), nil
}

func canonicalDiplomaticPair(a, b uuid.UUID) (uuid.UUID, uuid.UUID) {
	if a.String() <= b.String() {
		return a, b
	}
	return b, a
}
