package repo

import (
	"context"
	"database/sql"
	"errors"

	"github.com/artcodefun/heat-expansion-api/internal/game/core/domain"
	"github.com/artcodefun/heat-expansion-api/internal/game/core/ports"
	"github.com/artcodefun/heat-expansion-api/internal/game/infrastructure/db/gen"
	"github.com/artcodefun/heat-expansion-api/internal/game/infrastructure/db/mappers"
)

type DangerousLocationRepo struct {
	q              *gen.Queries
	armyProtoRepo  ports.ArmyPrototypeRepository
	buildProtoRepo ports.BuildPrototypeRepository
}

func NewDangerousLocationRepo(q *gen.Queries, armyProtoRepo ports.ArmyPrototypeRepository, buildProtoRepo ports.BuildPrototypeRepository) *DangerousLocationRepo {
	return &DangerousLocationRepo{q: q, armyProtoRepo: armyProtoRepo, buildProtoRepo: buildProtoRepo}
}

func (r *DangerousLocationRepo) Tx(tx ports.Transaction) ports.DangerousLocationRepository {
	if sqlTx, ok := tx.(*sql.Tx); ok {
		return &DangerousLocationRepo{
			q:              r.q.WithTx(sqlTx),
			armyProtoRepo:  r.armyProtoRepo.Tx(tx),
			buildProtoRepo: r.buildProtoRepo.Tx(tx),
		}
	}
	return r
}

func (r *DangerousLocationRepo) Create(loc *domain.DangerousLocationModel) error {
	id, err := r.q.InsertDangerousLocation(context.Background(), mappers.InsertDangerousLocationParamsFromDomain(loc))
	if err != nil {
		return err
	}
	loc.ID = int(id)
	return nil
}

func (r *DangerousLocationRepo) FindByID(id int) (*domain.DangerousLocationModel, error) {
	ctx := context.Background()
	row, err := r.q.GetDangerousLocationByID(ctx, int64(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ports.ErrNotFound
		}
		return nil, err
	}
	armyProtos, buildProtos, err := r.loadPrototypes(ctx)
	if err != nil {
		return nil, err
	}
	loc := mappers.DangerousLocationFromDB(row, armyProtos, buildProtos)
	return loc, nil
}

func (r *DangerousLocationRepo) FindByCoordinates(x, y int) (*domain.DangerousLocationModel, error) {
	ctx := context.Background()
	row, err := r.q.GetDangerousLocationBySector(ctx, gen.GetDangerousLocationBySectorParams{SectorX: int32(x), SectorY: int32(y)})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ports.ErrNotFound
		}
		return nil, err
	}
	armyProtos, buildProtos, err := r.loadPrototypes(ctx)
	if err != nil {
		return nil, err
	}
	loc := mappers.DangerousLocationFromDB(row, armyProtos, buildProtos)
	return loc, nil
}

func (r *DangerousLocationRepo) FindByCoordinatesForUpdate(x, y int) (*domain.DangerousLocationModel, error) {
	ctx := context.Background()
	row, err := r.q.GetDangerousLocationBySectorForUpdate(ctx, gen.GetDangerousLocationBySectorForUpdateParams{SectorX: int32(x), SectorY: int32(y)})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ports.ErrNotFound
		}
		return nil, err
	}
	armyProtos, buildProtos, err := r.loadPrototypes(ctx)
	if err != nil {
		return nil, err
	}
	loc := mappers.DangerousLocationFromDB(row, armyProtos, buildProtos)
	return loc, nil
}

func (r *DangerousLocationRepo) FindClosest(x, y int) (*domain.DangerousLocationModel, error) {
	ctx := context.Background()
	row, err := r.q.FindClosestDangerousLocation(ctx, gen.FindClosestDangerousLocationParams{X: int32(x), Y: int32(y)})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ports.ErrNotFound
		}
		return nil, err
	}
	armyProtos, buildProtos, err := r.loadPrototypes(ctx)
	if err != nil {
		return nil, err
	}
	loc := mappers.DangerousLocationFromDB(row, armyProtos, buildProtos)
	return loc, nil
}

func (r *DangerousLocationRepo) Update(loc *domain.DangerousLocationModel) error {
	return r.q.UpdateDangerousLocation(context.Background(), mappers.UpdateDangerousLocationParamsFromDomain(loc))
}

func (r *DangerousLocationRepo) Delete(id int) error {
	return r.q.DeleteDangerousLocation(context.Background(), int64(id))
}

func (r *DangerousLocationRepo) DeleteByCoordinates(x, y int) error {
	return r.q.DeleteDangerousLocationBySector(context.Background(), gen.DeleteDangerousLocationBySectorParams{
		SectorX: int32(x),
		SectorY: int32(y),
	})
}

func (r *DangerousLocationRepo) loadPrototypes(ctx context.Context) (map[int]*domain.ArmyItemPrototype, map[int]*domain.BuildItemPrototype, error) {
	// For now, load all of them as they are relatively few. In a larger game, we'd fetch only needed ones.
	armyList, err := r.armyProtoRepo.FindAllPrototypes()
	if err != nil {
		return nil, nil, err
	}
	buildList, err := r.buildProtoRepo.FindAllPrototypes()
	if err != nil {
		return nil, nil, err
	}
	armyMap := make(map[int]*domain.ArmyItemPrototype, len(armyList))
	for _, p := range armyList {
		armyMap[p.ID] = p
	}
	buildMap := make(map[int]*domain.BuildItemPrototype, len(buildList))
	for _, p := range buildList {
		buildMap[p.ID] = p
	}
	return armyMap, buildMap, nil
}
