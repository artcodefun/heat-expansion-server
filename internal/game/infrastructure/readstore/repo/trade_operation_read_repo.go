package repo

import (
	"context"
	"database/sql"
	"errors"

	"github.com/artcodefun/heat-expansion-server/internal/game/application/cqrs/readmodels"
	"github.com/artcodefun/heat-expansion-server/internal/game/application/ports"
	"github.com/artcodefun/heat-expansion-server/internal/game/infrastructure/readstore/gen"
	"github.com/artcodefun/heat-expansion-server/internal/game/infrastructure/readstore/mappers"
)

type TradeOperationReadRepo struct {
	q       *gen.Queries
	sectors ports.SectorReadRepository
}

func NewTradeOperationReadRepo(q *gen.Queries, sectors ports.SectorReadRepository) *TradeOperationReadRepo {
	return &TradeOperationReadRepo{q: q, sectors: sectors}
}

func (r *TradeOperationReadRepo) GetTradeOperation(ctx context.Context, operationID int) (*readmodels.TradeOperation, error) {
	row, err := r.q.GetTradeOperation(ctx, int64(operationID))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ports.ErrNotFound
		}
		return nil, err
	}
	armyMap, storageMap, err := r.loadPrototypeMaps(ctx)
	if err != nil {
		return nil, err
	}
	v := mappers.TradeOperationFromModel(row)
	if err := r.enrichTradeOperation(ctx, &v, armyMap, storageMap); err != nil {
		return nil, err
	}
	return &v, nil
}

func (r *TradeOperationReadRepo) ListActiveTradeOperations(ctx context.Context, baseID int) ([]*readmodels.TradeOperation, error) {
	rows, err := r.q.ListActiveTradeOperations(ctx, int64(baseID))
	if err != nil {
		return nil, err
	}
	armyMap, storageMap, err := r.loadPrototypeMaps(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]*readmodels.TradeOperation, 0, len(rows))
	for _, row := range rows {
		v := mappers.TradeOperationFromModel(row)
		if err := r.enrichTradeOperation(ctx, &v, armyMap, storageMap); err != nil {
			return nil, err
		}
		out = append(out, &v)
	}
	return out, nil
}

func (r *TradeOperationReadRepo) enrichTradeOperation(ctx context.Context, op *readmodels.TradeOperation, armyMap map[int]readmodels.ArmyItemPrototype, storageMap map[int]readmodels.StorageItemPrototype) error {
	mappers.EnrichTradeOperationItems(op, armyMap, storageMap)
	if op.CreatedAt <= 0 {
		return nil
	}
	if report, err := r.sectors.GetLatestScanBefore(ctx, op.SenderBaseID, op.TargetCoordinates.X, op.TargetCoordinates.Y, op.CreatedAt); err == nil {
		op.SenderScanOfReceiver = report
	} else if !errors.Is(err, ports.ErrNotFound) {
		return err
	}
	if report, err := r.sectors.GetLatestScanBefore(ctx, op.ReceiverBaseID, op.SourceCoordinates.X, op.SourceCoordinates.Y, op.CreatedAt); err == nil {
		op.ReceiverScanOfSender = report
	} else if !errors.Is(err, ports.ErrNotFound) {
		return err
	}
	return nil
}

func (r *TradeOperationReadRepo) loadPrototypeMaps(ctx context.Context) (map[int]readmodels.ArmyItemPrototype, map[int]readmodels.StorageItemPrototype, error) {
	armyRows, err := r.q.ListArmyPrototypes(ctx)
	if err != nil {
		return nil, nil, err
	}
	storageRows, err := r.q.ListStoragePrototypes(ctx)
	if err != nil {
		return nil, nil, err
	}
	armyMap := make(map[int]readmodels.ArmyItemPrototype, len(armyRows))
	for _, p := range armyRows {
		armyMap[int(p.ID)] = mappers.ArmyPrototypeFromModel(p)
	}
	storageMap := make(map[int]readmodels.StorageItemPrototype, len(storageRows))
	for _, p := range storageRows {
		storageMap[int(p.ID)] = mappers.StoragePrototypeFromModel(p)
	}
	return armyMap, storageMap, nil
}
