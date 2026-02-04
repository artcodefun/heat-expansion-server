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

type OperationReadRepo struct {
	q       *gen.Queries
	sectors ports.SectorReadRepository
}

func NewOperationReadRepo(rq *gen.Queries, sectors ports.SectorReadRepository) *OperationReadRepo {
	return &OperationReadRepo{q: rq, sectors: sectors}
}

func (r *OperationReadRepo) GetOperation(opID int) (*readmodels.MilitaryOperation, error) {
	row, err := r.q.GetOperation(context.Background(), int64(opID))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ports.ErrNotFound
		}
		return nil, err
	}
	v := mappers.OperationFromModel(row)
	if err := r.enrichOperation(&v, nil, nil, nil); err != nil {
		return nil, err
	}
	return &v, nil
}

func (r *OperationReadRepo) ListOperationsByBase(baseID int) ([]*readmodels.MilitaryOperation, error) {
	rows, err := r.q.ListOperationsByBase(context.Background(), int64(baseID))
	if err != nil {
		return nil, err
	}
	armyMap, buildMap, storageMap, err := r.loadPrototypeMaps()
	if err != nil {
		return nil, err
	}
	out := make([]*readmodels.MilitaryOperation, 0, len(rows))
	for _, row := range rows {
		v := mappers.OperationFromModel(row)
		if err := r.enrichOperation(&v, armyMap, buildMap, storageMap); err != nil {
			return nil, err
		}
		out = append(out, &v)
	}
	return out, nil
}

func (r *OperationReadRepo) ListActiveOperations(baseID int) ([]*readmodels.MilitaryOperation, error) {
	rows, err := r.q.ListActiveOperations(context.Background(), int64(baseID))
	if err != nil {
		return nil, err
	}
	armyMap, buildMap, storageMap, err := r.loadPrototypeMaps()
	if err != nil {
		return nil, err
	}
	out := make([]*readmodels.MilitaryOperation, 0, len(rows))
	for _, row := range rows {
		v := mappers.OperationFromModel(row)
		if err := r.enrichOperation(&v, armyMap, buildMap, storageMap); err != nil {
			return nil, err
		}
		out = append(out, &v)
	}
	return out, nil
}

// enrichOperation enriches a single operation.
func (r *OperationReadRepo) enrichOperation(v *readmodels.MilitaryOperation, army map[int]mappers.ArmyPrototypeSnapshot, build map[int]mappers.BuildPrototypeSnapshot, storage map[int]readmodels.StorageItemPrototype) error {
	// 1. Resolve prototypes if not provided
	if army == nil || build == nil || storage == nil {
		var err error
		army, build, storage, err = r.loadPrototypeMaps()
		if err != nil {
			return err
		}
	}

	// 2. Enrich with prototypes
	mappers.EnrichOperationUnitsAndStructures(v, army, build, storage)

	// 3. Fetch produced scan report if any
	reportRow, err := r.q.GetScanReportByOperationID(context.Background(), sql.NullInt64{Int64: int64(v.ID), Valid: true})
	if err == nil {
		report := mappers.SectorScanReportFromModel(reportRow)
		v.ProducedScanReport = &report
	} else if !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	// 4. Fetch prior scan report if coordinates and timeline are available.
	if v.TargetCoordinates != (readmodels.Vector2i{}) && v.OutboundDepartAt > 0 {
		report, err := r.sectors.GetLatestScanBefore(v.SourceBaseID, v.TargetCoordinates.X, v.TargetCoordinates.Y, v.OutboundDepartAt)
		if err == nil {
			v.PriorScanReport = report
		} else if !errors.Is(err, ports.ErrNotFound) {
			return err
		}
	}
	return nil
}

// loadPrototypeMaps fetches all army/build prototypes for read-store and indexes them by ID.
func (r *OperationReadRepo) loadPrototypeMaps() (map[int]mappers.ArmyPrototypeSnapshot, map[int]mappers.BuildPrototypeSnapshot, map[int]readmodels.StorageItemPrototype, error) {
	armyRows, err := r.q.ListArmyPrototypes(context.Background())
	if err != nil {
		return nil, nil, nil, err
	}
	buildRows, err := r.q.ListBuildPrototypes(context.Background())
	if err != nil {
		return nil, nil, nil, err
	}
	storageRows, err := r.q.ListStoragePrototypes(context.Background())
	if err != nil {
		return nil, nil, nil, err
	}
	armyMap := make(map[int]mappers.ArmyPrototypeSnapshot, len(armyRows))
	for _, p := range armyRows {
		name := p.Name
		imageURL := ""
		if p.ImageUrl.Valid {
			imageURL = p.ImageUrl.String
		}
		armyMap[int(p.ID)] = mappers.ArmyPrototypeSnapshot{Name: name, ImageURL: imageURL, Space: int(p.Space)}
	}
	buildMap := make(map[int]mappers.BuildPrototypeSnapshot, len(buildRows))
	for _, p := range buildRows {
		name := p.Name
		imageURL := ""
		if p.ImageUrl.Valid {
			imageURL = p.ImageUrl.String
		}
		buildMap[int(p.ID)] = mappers.BuildPrototypeSnapshot{Name: name, ImageURL: imageURL, Space: int(p.Space)}
	}
	storageMap := make(map[int]readmodels.StorageItemPrototype, len(storageRows))
	for _, p := range storageRows {
		storageMap[int(p.ID)] = mappers.StoragePrototypeFromModel(p)
	}
	return armyMap, buildMap, storageMap, nil
}
