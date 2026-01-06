package repo

import (
	"context"
	"errors"

	"github.com/artcodefun/heat-expansion-api/internal/core/cqrs/readmodels"
	"github.com/artcodefun/heat-expansion-api/internal/core/ports"
	"github.com/artcodefun/heat-expansion-api/internal/infrastructure/readstore/gen"
	"github.com/artcodefun/heat-expansion-api/internal/infrastructure/readstore/mappers"
)

type ActivityReadRepo struct {
	q       *gen.Queries
	ops     ports.OperationReadRepository
	sectors ports.SectorReadRepository
}

func NewActivityReadRepo(q *gen.Queries, ops ports.OperationReadRepository, sectors ports.SectorReadRepository) *ActivityReadRepo {
	return &ActivityReadRepo{q: q, ops: ops, sectors: sectors}
}

func (r *ActivityReadRepo) ListActivities(baseID int, limit int) ([]*readmodels.ActivityItem, error) {
	rows, err := r.q.ListActivities(context.Background(), gen.ListActivitiesParams{BaseID: int64(baseID), Limit: int32(limit)})
	if err != nil {
		return nil, err
	}
	out := make([]*readmodels.ActivityItem, 0, len(rows))
	for _, a := range rows {
		v := mappers.ActivityItemFromModel(a)
		if err := r.enrichActivity(&v); err != nil {
			return nil, err
		}
		out = append(out, &v)
	}
	return out, nil
}

func (r *ActivityReadRepo) ListActivitiesByKind(baseID int, kind readmodels.ActivityKind, limit int) ([]*readmodels.ActivityItem, error) {
	rows, err := r.q.ListActivitiesByKind(context.Background(), gen.ListActivitiesByKindParams{BaseID: int64(baseID), Kind: string(kind), Limit: int32(limit)})
	if err != nil {
		return nil, err
	}
	out := make([]*readmodels.ActivityItem, 0, len(rows))
	for _, a := range rows {
		v := mappers.ActivityItemFromModel(a)
		if err := r.enrichActivity(&v); err != nil {
			return nil, err
		}
		out = append(out, &v)
	}
	return out, nil
}

func (r *ActivityReadRepo) enrichActivity(v *readmodels.ActivityItem) error {
	if v.Operation != nil {
		op, err := r.ops.GetOperation(v.Operation.OpID)
		if err != nil && !errors.Is(err, ports.ErrNotFound) {
			return err
		}
		if err == nil {
			v.Operation.Operation = op
			// Enrich with prior opponent scan if coordinates and timeline are available.
			var target readmodels.Vector2i
			switch v.Operation.Role {
			case readmodels.OperationRoleAttacker:
				target = op.TargetCoordinates
			case readmodels.OperationRoleDefender:
				target = op.SourceCoordinates
			}
			if target != (readmodels.Vector2i{}) && op.OutboundDepartAt > 0 {
				report, err := r.sectors.GetLatestScanBefore(v.BaseID, target.X, target.Y, op.OutboundDepartAt)
				if err != nil && !errors.Is(err, ports.ErrNotFound) {
					return err
				}
				if err == nil {
					v.Operation.PriorOpponentScan = report
				}
			}
		}
	}
	if v.Scan != nil {
		report, err := r.sectors.GetScanReportByID(v.BaseID, v.Scan.ReportID)
		if err != nil && !errors.Is(err, ports.ErrNotFound) {
			return err
		}
		if err == nil {
			v.Scan.Report = report
		}
	}
	return nil
}
