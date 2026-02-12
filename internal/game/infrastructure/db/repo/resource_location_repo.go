package repo

import (
	"context"
	"database/sql"
	"errors"

	"github.com/artcodefun/heat-expansion-server/internal/game/application/ports"
	"github.com/artcodefun/heat-expansion-server/internal/game/domain"
	"github.com/artcodefun/heat-expansion-server/internal/game/infrastructure/db/gen"
	"github.com/artcodefun/heat-expansion-server/internal/game/infrastructure/db/mappers"
)

type ResourceLocationRepo struct {
	q              *gen.Queries
	armyProtoRepo  ports.ArmyPrototypeRepository
	buildProtoRepo ports.BuildPrototypeRepository
}

func NewResourceLocationRepo(q *gen.Queries, armyProtoRepo ports.ArmyPrototypeRepository, buildProtoRepo ports.BuildPrototypeRepository) *ResourceLocationRepo {
	return &ResourceLocationRepo{q: q, armyProtoRepo: armyProtoRepo, buildProtoRepo: buildProtoRepo}
}

func (r *ResourceLocationRepo) Tx(tx ports.Transaction) ports.ResourceLocationRepository {
	if sqlTx, ok := tx.(*sql.Tx); ok {
		return &ResourceLocationRepo{
			q:              r.q.WithTx(sqlTx),
			armyProtoRepo:  r.armyProtoRepo.Tx(tx),
			buildProtoRepo: r.buildProtoRepo.Tx(tx),
		}
	}
	return r
}

func (r *ResourceLocationRepo) Create(loc *domain.ResourceLocationModel) error {
	id, err := r.q.InsertResourceLocation(context.Background(), mappers.InsertResourceLocationParamsFromDomain(loc))
	if err != nil {
		return err
	}
	loc.ID = int(id)
	return nil
}

func (r *ResourceLocationRepo) FindByID(id int) (*domain.ResourceLocationModel, error) {
	ctx := context.Background()
	row, err := r.q.GetResourceLocationByID(ctx, int64(id))
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
	loc := mappers.ResourceLocationFromDB(row, armyProtos, buildProtos)
	return loc, nil
}

func (r *ResourceLocationRepo) FindByCoordinates(x, y int) (*domain.ResourceLocationModel, error) {
	ctx := context.Background()
	row, err := r.q.GetResourceLocationBySector(ctx, gen.GetResourceLocationBySectorParams{SectorX: int32(x), SectorY: int32(y)})
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
	loc := mappers.ResourceLocationFromDB(row, armyProtos, buildProtos)
	return loc, nil
}

func (r *ResourceLocationRepo) FindByCoordinatesForUpdate(x, y int) (*domain.ResourceLocationModel, error) {
	ctx := context.Background()
	row, err := r.q.GetResourceLocationBySectorForUpdate(ctx, gen.GetResourceLocationBySectorForUpdateParams{SectorX: int32(x), SectorY: int32(y)})
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
	loc := mappers.ResourceLocationFromDB(row, armyProtos, buildProtos)
	return loc, nil
}

func (r *ResourceLocationRepo) FindClosest(x, y int) (*domain.ResourceLocationModel, error) {
	ctx := context.Background()
	row, err := r.q.FindClosestResourceLocation(ctx, gen.FindClosestResourceLocationParams{X: int32(x), Y: int32(y)})
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
	loc := mappers.ResourceLocationFromDB(row, armyProtos, buildProtos)
	return loc, nil
}

func (r *ResourceLocationRepo) Update(loc *domain.ResourceLocationModel) error {
	return r.q.UpdateResourceLocation(context.Background(), mappers.UpdateResourceLocationParamsFromDomain(loc))
}

func (r *ResourceLocationRepo) Delete(id int) error {
	return r.q.DeleteResourceLocation(context.Background(), int64(id))
}

func (r *ResourceLocationRepo) DeleteByCoordinates(x, y int) error {
	return r.q.DeleteResourceLocationBySector(context.Background(), gen.DeleteResourceLocationBySectorParams{
		SectorX: int32(x),
		SectorY: int32(y),
	})
}

func (r *ResourceLocationRepo) loadPrototypes(ctx context.Context) (map[int]*domain.ArmyItemPrototype, map[int]*domain.BuildItemPrototype, error) {
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
