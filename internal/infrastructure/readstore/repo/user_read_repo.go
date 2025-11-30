package repo

import (
	"context"
	"database/sql"
	"errors"

	"github.com/artcodefun/heat-expansion-api/internal/core/cqrs/readmodels"
	"github.com/artcodefun/heat-expansion-api/internal/core/ports"
	"github.com/artcodefun/heat-expansion-api/internal/infrastructure/readstore/gen"
	"github.com/artcodefun/heat-expansion-api/internal/infrastructure/readstore/mappers"
)

type UserReadRepo struct{ q *gen.Queries }

func NewUserReadRepo(q *gen.Queries) *UserReadRepo { return &UserReadRepo{q: q} }

func (r *UserReadRepo) GetUserProfile(userID int) (*readmodels.User, error) {
	row, err := r.q.GetUserProfile(context.Background(), int64(userID))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ports.ErrNotFound
		}
		return nil, err
	}
	model := mappers.UserFromModel(row)
	return &model, nil
}
