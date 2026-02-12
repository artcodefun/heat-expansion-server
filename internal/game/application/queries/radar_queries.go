package queries

import (
	"github.com/artcodefun/heat-expansion-api/internal/game/application/cqrs"
	"github.com/artcodefun/heat-expansion-api/internal/game/application/cqrs/readmodels"
	"github.com/artcodefun/heat-expansion-api/internal/game/application/ports"
	"github.com/artcodefun/heat-expansion-api/internal/game/application/services"
)

type RadarQueries struct {
	Repo   ports.RadarReadRepository
	Access *services.AccessControlService
}

func NewRadarQueries(repo ports.RadarReadRepository, access *services.AccessControlService) *RadarQueries {
	return &RadarQueries{Repo: repo, Access: access}
}

func (q *RadarQueries) ListIncomingThreats(ctx cqrs.QueryContext, baseID int) ([]*readmodels.RadarThreat, error) {
	if err := q.Access.EnsureBaseOwnership(ctx.UserID, baseID); err != nil {
		return nil, err
	}
	threats, err := q.Repo.ListIncomingThreats(baseID)
	return threats, repoErr(err)
}
