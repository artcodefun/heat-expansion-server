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
	radar   ports.RadarReadRepository
}

func NewActivityReadRepo(q *gen.Queries, ops ports.OperationReadRepository, sectors ports.SectorReadRepository, radar ports.RadarReadRepository) *ActivityReadRepo {
	return &ActivityReadRepo{q: q, ops: ops, sectors: sectors, radar: radar}
}

func (r *ActivityReadRepo) ListOffenseActivities(baseID int, subtype readmodels.OffenseActivitySubtype, limit int) ([]*readmodels.ActivityItem, error) {
	rows, err := r.q.ListOffenseActivities(context.Background(), gen.ListOffenseActivitiesParams{BaseID: int64(baseID), Column2: string(subtype), Limit: int32(limit)})
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

func (r *ActivityReadRepo) ListDefenseActivities(baseID int, subtype readmodels.DefenseActivitySubtype, limit int) ([]*readmodels.ActivityItem, error) {
	rows, err := r.q.ListDefenseActivities(context.Background(), gen.ListDefenseActivitiesParams{BaseID: int64(baseID), Column2: string(subtype), Limit: int32(limit)})
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

func (r *ActivityReadRepo) ListScanActivities(baseID int, subtype readmodels.ScanActivitySubtype, limit int) ([]*readmodels.ActivityItem, error) {
	rows, err := r.q.ListScanActivities(context.Background(), gen.ListScanActivitiesParams{BaseID: int64(baseID), Column2: string(subtype), Limit: int32(limit)})
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

func (r *ActivityReadRepo) ListRadarActivities(baseID int, limit int) ([]*readmodels.ActivityItem, error) {
	rows, err := r.q.ListRadarActivities(context.Background(), gen.ListRadarActivitiesParams{BaseID: int64(baseID), Limit: int32(limit)})
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

func (r *ActivityReadRepo) ListTradeActivities(baseID int, limit int) ([]*readmodels.ActivityItem, error) {
	rows, err := r.q.ListTradeActivities(context.Background(), gen.ListTradeActivitiesParams{BaseID: int64(baseID), Limit: int32(limit)})
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
	if v.Offense != nil {
		op, err := r.ops.GetOperation(v.Offense.OpID)
		if err != nil && !errors.Is(err, ports.ErrNotFound) {
			return err
		}
		if err == nil {
			v.Offense.Operation = op
		}
	}
	if v.Defense != nil {
		op, err := r.ops.GetOperation(v.Defense.OpID)
		if err != nil && !errors.Is(err, ports.ErrNotFound) {
			return err
		}
		if err == nil {
			v.Defense.Offender = &readmodels.OffenderInfo{
				Type:              op.Type,
				SourceCoordinates: op.SourceCoordinates,
				TargetCoordinates: op.TargetCoordinates,
				ContactDate:       op.OutboundArriveAt,
				Result:            op.Result,
				Units:             op.Units,
				SpyResult:         op.SpyResult,
				AttackResult:      op.AttackResult,
			}
			// Enrich with prior opponent scan for the defender (scan of the offender's source coordinates)
			if op.SourceCoordinates != (readmodels.Vector2i{}) && op.OutboundArriveAt > 0 {
				report, err := r.sectors.GetLatestScanBefore(v.BaseID, op.SourceCoordinates.X, op.SourceCoordinates.Y, op.OutboundArriveAt)
				if err != nil && !errors.Is(err, ports.ErrNotFound) {
					return err
				}
				if err == nil {
					v.Defense.PriorOpponentScan = report
				}
			}
		}
	}
	if v.Scan != nil && v.Scan.ReportID != nil {
		report, err := r.sectors.GetScanReportByID(v.BaseID, *v.Scan.ReportID)
		if err != nil && !errors.Is(err, ports.ErrNotFound) {
			return err
		}
		if err == nil {
			v.Scan.Report = report
		}
	}
	if v.Radar != nil {
		threat, err := r.radar.GetRadarThreat(v.Radar.ThreatID)
		if err != nil && !errors.Is(err, ports.ErrNotFound) {
			return err
		}
		if err == nil {
			v.Radar.Threat = threat
		}
	}
	return nil
}
