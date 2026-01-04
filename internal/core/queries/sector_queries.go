package queries

import (
	"github.com/artcodefun/heat-expansion-api/internal/core/cqrs"
	"github.com/artcodefun/heat-expansion-api/internal/core/cqrs/readmodels"
	"github.com/artcodefun/heat-expansion-api/internal/core/ports"
	"github.com/artcodefun/heat-expansion-api/internal/core/services"
)

type SectorQueries struct {
	Repo   ports.SectorReadRepository
	Access *services.AccessControlService
}

func NewSectorQueries(repo ports.SectorReadRepository, access *services.AccessControlService) *SectorQueries {
	return &SectorQueries{Repo: repo, Access: access}
}
func (q *SectorQueries) GetScansNear(ctx cqrs.QueryContext, baseID int, centerX, centerY, radius int) ([]*readmodels.SectorScanReport, error) {
	if err := q.Access.EnsureBaseOwnership(ctx.UserID, baseID); err != nil {
		return nil, err
	}
	reports, err := q.Repo.GetScansNear(baseID, centerX, centerY, radius)
	return reports, repoErr(err)
}

func (q *SectorQueries) GetScanReportByID(ctx cqrs.QueryContext, baseID, id int) (*readmodels.SectorScanReport, error) {
	if err := q.Access.EnsureBaseOwnership(ctx.UserID, baseID); err != nil {
		return nil, err
	}
	report, err := q.Repo.GetScanReportByID(baseID, id)
	if err != nil {
		return nil, repoErr(err)
	}
	return report, nil
}

func (q *SectorQueries) GetLatestScanBefore(ctx cqrs.QueryContext, baseID, x, y int, before int64) (*readmodels.SectorScanReport, error) {
	if err := q.Access.EnsureBaseOwnership(ctx.UserID, baseID); err != nil {
		return nil, err
	}
	report, err := q.Repo.GetLatestScanBefore(baseID, x, y, before)
	if err != nil {
		return nil, repoErr(err)
	}
	return report, nil
}
