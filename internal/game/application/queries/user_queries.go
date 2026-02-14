package queries

import (
	"github.com/artcodefun/heat-expansion-server/internal/game/application/cqrs"
	"github.com/artcodefun/heat-expansion-server/internal/game/application/cqrs/readmodels"
	"github.com/artcodefun/heat-expansion-server/internal/game/application/ports"
	"github.com/google/uuid"
)

type UserQueries struct{ Repo ports.UserReadRepository }

func NewUserQueries(repo ports.UserReadRepository) *UserQueries { return &UserQueries{Repo: repo} }

func (q *UserQueries) GetUserProfile(_ cqrs.QueryContext, userID uuid.UUID) (*readmodels.User, error) {
	user, err := q.Repo.GetUserProfile(userID)
	return user, repoErr(err)
}
