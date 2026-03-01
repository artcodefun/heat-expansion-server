package queries

import (
	"context"

	"github.com/artcodefun/heat-expansion-server/internal/game/application/cqrs"
	"github.com/artcodefun/heat-expansion-server/internal/game/application/cqrs/readmodels"
	"github.com/artcodefun/heat-expansion-server/internal/game/application/ports"
	"github.com/artcodefun/heat-expansion-server/internal/game/application/services"
)

type SectorQueries struct {
	Repo   ports.SectorReadRepository
	Access *services.AccessControlService
}

func NewSectorQueries(repo ports.SectorReadRepository, access *services.AccessControlService) *SectorQueries {
	return &SectorQueries{Repo: repo, Access: access}
}
func (q *SectorQueries) GetScansNear(ctx context.Context, actor cqrs.Actor, baseID int, centerX, centerY, radius int) ([]*readmodels.SectorScanReport, error) {
	if err := q.Access.EnsureBaseOwnership(ctx, actor.UserID, baseID); err != nil {
		return nil, err
	}
	reports, err := q.Repo.GetScansNear(ctx, baseID, centerX, centerY, radius)
	return reports, repoErr(err)
}

func (q *SectorQueries) GetScanReportByID(ctx context.Context, actor cqrs.Actor, baseID, id int) (*readmodels.SectorScanReport, error) {
	if err := q.Access.EnsureBaseOwnership(ctx, actor.UserID, baseID); err != nil {
		return nil, err
	}
	report, err := q.Repo.GetScanReportByID(ctx, baseID, id)
	if err != nil {
		return nil, repoErr(err)
	}
	return report, nil
}

func (q *SectorQueries) GetLatestScanBefore(ctx context.Context, actor cqrs.Actor, baseID, x, y int, before int64) (*readmodels.SectorScanReport, error) {
	if err := q.Access.EnsureBaseOwnership(ctx, actor.UserID, baseID); err != nil {
		return nil, err
	}
	report, err := q.Repo.GetLatestScanBefore(ctx, baseID, x, y, before)
	if err != nil {
		return nil, repoErr(err)
	}
	return report, nil
}
