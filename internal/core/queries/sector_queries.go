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

func (q *SectorQueries) GetSector(_ cqrs.QueryContext, x, y int) (*readmodels.SectorModel, error) {
	sector, err := q.Repo.GetSector(x, y)
	return sector, repoErr(err)
}
func (q *SectorQueries) GetLatestScans(ctx cqrs.QueryContext, baseID int) ([]*readmodels.SectorScanReport, error) {
	if err := q.Access.EnsureBaseOwnership(ctx.UserID, baseID); err != nil {
		return nil, err
	}
	reports, err := q.Repo.GetLatestScans(baseID)
	return reports, repoErr(err)
}
func (q *SectorQueries) GetScansNear(ctx cqrs.QueryContext, baseID int, centerX, centerY, radius int) ([]*readmodels.SectorScanReport, error) {
	if err := q.Access.EnsureBaseOwnership(ctx.UserID, baseID); err != nil {
		return nil, err
	}
	reports, err := q.Repo.GetScansNear(baseID, centerX, centerY, radius)
	return reports, repoErr(err)
}

func (q *SectorQueries) ListOccupiedCoordinates(_ cqrs.QueryContext) ([]readmodels.Vector2i, error) {
	coords, err := q.Repo.ListOccupiedCoordinates()
	return coords, repoErr(err)
}
func (q *SectorQueries) ListSectorsInRadius(_ cqrs.QueryContext, centerX, centerY, radius int) ([]*readmodels.SectorModel, error) {
	sectors, err := q.Repo.ListSectorsInRadius(centerX, centerY, radius)
	return sectors, repoErr(err)
}
