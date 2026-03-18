package repo

import (
	"context"
	"errors"

	"github.com/artcodefun/heat-expansion-server/internal/game/application/cqrs/readmodels"
	"github.com/artcodefun/heat-expansion-server/internal/game/application/ports"
	"github.com/artcodefun/heat-expansion-server/internal/game/infrastructure/readstore/gen"
	"github.com/artcodefun/heat-expansion-server/internal/game/infrastructure/readstore/mappers"
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

func (r *ActivityReadRepo) ListOffenseActivities(ctx context.Context, baseID int, subtype readmodels.OffenseActivitySubtype, limit int) ([]*readmodels.ActivityItem, error) {
	rows, err := r.q.ListOffenseActivities(ctx, gen.ListOffenseActivitiesParams{BaseID: int64(baseID), Column2: string(subtype), Limit: int32(limit)})
	if err != nil {
		return nil, err
	}
	out := make([]*readmodels.ActivityItem, 0, len(rows))
	for _, a := range rows {
		v := mappers.ActivityItemFromModel(a)
		if err := r.enrichActivity(ctx, &v); err != nil {
			return nil, err
		}
		out = append(out, &v)
	}
	return out, nil
}

func (r *ActivityReadRepo) ListDefenseActivities(ctx context.Context, baseID int, subtype readmodels.DefenseActivitySubtype, limit int) ([]*readmodels.ActivityItem, error) {
	rows, err := r.q.ListDefenseActivities(ctx, gen.ListDefenseActivitiesParams{BaseID: int64(baseID), Column2: string(subtype), Limit: int32(limit)})
	if err != nil {
		return nil, err
	}
	out := make([]*readmodels.ActivityItem, 0, len(rows))
	for _, a := range rows {
		v := mappers.ActivityItemFromModel(a)
		if err := r.enrichActivity(ctx, &v); err != nil {
			return nil, err
		}
		out = append(out, &v)
	}
	return out, nil
}

func (r *ActivityReadRepo) ListScanActivities(ctx context.Context, baseID int, subtype readmodels.ScanActivitySubtype, limit int) ([]*readmodels.ActivityItem, error) {
	rows, err := r.q.ListScanActivities(ctx, gen.ListScanActivitiesParams{BaseID: int64(baseID), Column2: string(subtype), Limit: int32(limit)})
	if err != nil {
		return nil, err
	}
	out := make([]*readmodels.ActivityItem, 0, len(rows))
	for _, a := range rows {
		v := mappers.ActivityItemFromModel(a)
		if err := r.enrichActivity(ctx, &v); err != nil {
			return nil, err
		}
		out = append(out, &v)
	}
	return out, nil
}

func (r *ActivityReadRepo) ListRadarActivities(ctx context.Context, baseID int, limit int) ([]*readmodels.ActivityItem, error) {
	rows, err := r.q.ListRadarActivities(ctx, gen.ListRadarActivitiesParams{BaseID: int64(baseID), Limit: int32(limit)})
	if err != nil {
		return nil, err
	}
	out := make([]*readmodels.ActivityItem, 0, len(rows))
	for _, a := range rows {
		v := mappers.ActivityItemFromModel(a)
		if err := r.enrichActivity(ctx, &v); err != nil {
			return nil, err
		}
		out = append(out, &v)
	}
	return out, nil
}

func (r *ActivityReadRepo) ListTradeActivities(ctx context.Context, baseID int, limit int) ([]*readmodels.ActivityItem, error) {
	rows, err := r.q.ListTradeActivities(ctx, gen.ListTradeActivitiesParams{BaseID: int64(baseID), Limit: int32(limit)})
	if err != nil {
		return nil, err
	}
	out := make([]*readmodels.ActivityItem, 0, len(rows))
	for _, a := range rows {
		v := mappers.ActivityItemFromModel(a)
		if err := r.enrichActivity(ctx, &v); err != nil {
			return nil, err
		}
		out = append(out, &v)
	}
	return out, nil
}

func (r *ActivityReadRepo) enrichActivity(ctx context.Context, v *readmodels.ActivityItem) error {
	if v.Offense != nil {
		op, err := r.ops.GetOperation(ctx, v.Offense.OpID)
		if err != nil && !errors.Is(err, ports.ErrNotFound) {
			return err
		}
		if err == nil {
			v.Offense.Operation = op
		}
	}
	if v.Defense != nil {
		op, err := r.ops.GetOperation(ctx, v.Defense.OpID)
		if err != nil && !errors.Is(err, ports.ErrNotFound) {
			return err
		}
		if err == nil {
			v.Defense.Offender = readmodels.NewOffenderInfoFromOperation(op)

			// Enrich with prior opponent scan for the defender (scan of the offender's source coordinates)
			if v.Defense.Offender.SourceCoordinates != nil && op.OutboundArriveAt > 0 {
				report, err := r.sectors.GetLatestScanBefore(ctx, v.BaseID, v.Defense.Offender.SourceCoordinates.X, v.Defense.Offender.SourceCoordinates.Y, op.OutboundArriveAt)
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
		report, err := r.sectors.GetScanReportByID(ctx, v.BaseID, *v.Scan.ReportID)
		if err != nil && !errors.Is(err, ports.ErrNotFound) {
			return err
		}
		if err == nil {
			v.Scan.Report = report
		}
	}
	if v.Radar != nil {
		threat, err := r.radar.GetRadarThreat(ctx, v.Radar.ThreatID)
		if err != nil && !errors.Is(err, ports.ErrNotFound) {
			return err
		}
		if err == nil {
			v.Radar.Threat = threat
		}
	}
	return nil
}
