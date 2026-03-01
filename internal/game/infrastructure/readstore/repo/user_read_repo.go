package repo

import (
	"context"
	"database/sql"
	"errors"

	"github.com/artcodefun/heat-expansion-server/internal/game/application/cqrs/readmodels"
	"github.com/artcodefun/heat-expansion-server/internal/game/application/ports"
	"github.com/artcodefun/heat-expansion-server/internal/game/infrastructure/readstore/gen"
	"github.com/artcodefun/heat-expansion-server/internal/game/infrastructure/readstore/mappers"
	"github.com/google/uuid"
)

type UserReadRepo struct{ q *gen.Queries }

func NewUserReadRepo(q *gen.Queries) *UserReadRepo { return &UserReadRepo{q: q} }

func (r *UserReadRepo) GetUserProfile(ctx context.Context, userID uuid.UUID) (*readmodels.User, error) {
	row, err := r.q.GetUserProfile(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ports.ErrNotFound
		}
		return nil, err
	}
	model := mappers.UserFromModel(row)
	return &model, nil
}
