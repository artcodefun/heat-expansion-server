package repo

import (
	"context"
	"database/sql"
	"errors"

	"github.com/artcodefun/heat-expansion-server/internal/game/application/ports"
	"github.com/artcodefun/heat-expansion-server/internal/game/domain"
	"github.com/artcodefun/heat-expansion-server/internal/game/infrastructure/db/gen"
	"github.com/artcodefun/heat-expansion-server/internal/game/infrastructure/db/mappers"
	"github.com/google/uuid"
)

type UserBaseRepo struct {
	q *gen.Queries
}

func NewUserBaseRepo(q *gen.Queries) *UserBaseRepo { return &UserBaseRepo{q: q} }

func (r *UserBaseRepo) Tx(tx ports.Transaction) ports.UserBaseRepository {
	if sqlTx, ok := tx.(*sql.Tx); ok {
		return &UserBaseRepo{q: r.q.WithTx(sqlTx)}
	}
	return r
}

func (r *UserBaseRepo) Create(ctx context.Context, base *domain.UserBaseModel) error {
	row, err := r.q.CreateBase(ctx, mappers.InsertBaseParamsFromDomain(base))
	if err != nil {
		return err
	}
	created := mappers.UserBaseFromDB(row)
	base.ID = created.ID
	base.Stats = created.Stats
	base.Name = created.Name
	base.Description = created.Description
	base.ImageURL = created.ImageURL

	return r.persistAllItems(ctx, base)
}

func (r *UserBaseRepo) FindByID(ctx context.Context, id int) (*domain.UserBaseModel, error) {
	row, err := r.q.GetBaseByID(ctx, int64(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ports.ErrNotFound
		}
		return nil, err
	}
	base := mappers.UserBaseFromDB(row)
	_ = r.hydrateBase(ctx, base)
	return base, nil
}

func (r *UserBaseRepo) FindByIDForUpdate(ctx context.Context, id int) (*domain.UserBaseModel, error) {
	row, err := r.q.GetBaseByIDForUpdate(ctx, int64(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ports.ErrNotFound
		}
		return nil, err
	}
	base := mappers.UserBaseFromDB(row)
	_ = r.hydrateBase(ctx, base)
	return base, nil
}

// Update replaces all per-item rows (army/build/tech/storage) with the current aggregate state and updates base stats.
func (r *UserBaseRepo) Update(ctx context.Context, base *domain.UserBaseModel) error {
	// Update base stats and metadata first.
	if _, err := r.q.UpdateBase(ctx, mappers.UpdateBaseParamsFromDomain(base)); err != nil {
		return err
	}
	return r.persistAllItems(ctx, base)
}

func (r *UserBaseRepo) Delete(ctx context.Context, id int) error {
	return r.q.DeleteBase(ctx, int64(id))
}

func (r *UserBaseRepo) FindByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.UserBaseModel, error) {
	rows, err := r.q.ListBasesByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	out := make([]*domain.UserBaseModel, 0, len(rows))
	for _, b := range rows {
		base := mappers.UserBaseFromDB(b)
		_ = r.hydrateBase(ctx, base)
		out = append(out, base)
	}
	return out, nil
}

func (r *UserBaseRepo) FindByCoordinates(ctx context.Context, x, y int) (*domain.UserBaseModel, error) {
	row, err := r.q.GetBaseByCoordinates(ctx, gen.GetBaseByCoordinatesParams{SectorX: int32(x), SectorY: int32(y)})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ports.ErrNotFound
		}
		return nil, err
	}
	base := mappers.UserBaseFromDB(row)
	_ = r.hydrateBase(ctx, base)
	return base, nil
}

func (r *UserBaseRepo) FindByCoordinatesForUpdate(ctx context.Context, x, y int) (*domain.UserBaseModel, error) {
	row, err := r.q.GetBaseByCoordinatesForUpdate(ctx, gen.GetBaseByCoordinatesForUpdateParams{SectorX: int32(x), SectorY: int32(y)})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ports.ErrNotFound
		}
		return nil, err
	}
	base := mappers.UserBaseFromDB(row)
	_ = r.hydrateBase(ctx, base)
	return base, nil
}

func (r *UserBaseRepo) FindClosest(ctx context.Context, x, y int) (*domain.UserBaseModel, error) {
	row, err := r.q.FindClosestBase(ctx, gen.FindClosestBaseParams{X: int32(x), Y: int32(y)})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ports.ErrNotFound
		}
		return nil, err
	}
	base := mappers.UserBaseFromDB(row)
	_ = r.hydrateBase(ctx, base)
	return base, nil
}

func (r *UserBaseRepo) FindAll(ctx context.Context) ([]*domain.UserBaseModel, error) {
	rows, err := r.q.ListAllBases(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]*domain.UserBaseModel, 0, len(rows))
	for _, b := range rows {
		base := mappers.UserBaseFromDB(b)
		out = append(out, base)
	}
	return out, nil
}

// GetOwnerID returns the owning user ID for a base.
func (r *UserBaseRepo) GetOwnerID(ctx context.Context, baseID int) (uuid.UUID, error) {
	row, err := r.q.GetBaseByID(ctx, int64(baseID))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return uuid.Nil, ports.ErrNotFound
		}
		return uuid.Nil, err
	}
	return row.UserID, nil
}

// persistAllItems updates the base stats row and fully replaces all item rows with the current aggregate collections.
func (r *UserBaseRepo) persistAllItems(ctx context.Context, base *domain.UserBaseModel) error {
	if err := r.persistArmyItems(ctx, base); err != nil {
		return err
	}
	if err := r.persistBuildItems(ctx, base); err != nil {
		return err
	}
	if err := r.persistTechItems(ctx, base); err != nil {
		return err
	}
	if err := r.persistStorageItems(ctx, base); err != nil {
		return err
	}
	return nil
}

// --- Internal per-table persistence helpers ---

func (r *UserBaseRepo) persistArmyItems(ctx context.Context, base *domain.UserBaseModel) error {
	if err := r.q.DeleteBaseArmyItemsByBase(ctx, int64(base.ID)); err != nil {
		return err
	}
	params := mappers.DehydrateArmyItems(base)
	for _, p := range params {
		if _, err := r.q.InsertBaseArmyItem(ctx, p); err != nil {
			return err
		}
	}
	return nil
}

func (r *UserBaseRepo) persistBuildItems(ctx context.Context, base *domain.UserBaseModel) error {
	if err := r.q.DeleteBaseBuildItemsByBase(ctx, int64(base.ID)); err != nil {
		return err
	}
	params := mappers.DehydrateBuildItems(base)
	for _, p := range params {
		if _, err := r.q.InsertBaseBuildItem(ctx, p); err != nil {
			return err
		}
	}
	return nil
}

func (r *UserBaseRepo) persistTechItems(ctx context.Context, base *domain.UserBaseModel) error {
	if err := r.q.DeleteBaseTechItemsByBase(ctx, int64(base.ID)); err != nil {
		return err
	}
	params := mappers.DehydrateTechItems(base)
	for _, p := range params {
		if _, err := r.q.InsertBaseTechItem(ctx, p); err != nil {
			return err
		}
	}
	return nil
}

func (r *UserBaseRepo) persistStorageItems(ctx context.Context, base *domain.UserBaseModel) error {
	if err := r.q.DeleteBaseStorageItemsByBase(ctx, int64(base.ID)); err != nil {
		return err
	}
	params := mappers.DehydrateStorageItems(base)
	for _, p := range params {
		if _, err := r.q.InsertBaseStorageItem(ctx, p); err != nil {
			return err
		}
	}
	return nil
}

// hydrateBase loads all item rows and prototypes, and populates domain collections.
func (r *UserBaseRepo) hydrateBase(ctx context.Context, base *domain.UserBaseModel) error {

	// Load item rows
	armyRows, err := r.q.ListBaseArmyItems(ctx, int64(base.ID))
	if err != nil {
		return err
	}
	buildRows, err := r.q.ListBaseBuildItems(ctx, int64(base.ID))
	if err != nil {
		return err
	}
	techRows, err := r.q.ListBaseTechItems(ctx, int64(base.ID))
	if err != nil {
		return err
	}
	storageRows, err := r.q.ListBaseStorageItems(ctx, int64(base.ID))
	if err != nil {
		return err
	}

	// Load prototypes via existing read-only repos
	armyProtoRepo := NewArmyPrototypeRepo(r.q)
	armyProtos, err := armyProtoRepo.FindAllPrototypes(ctx)
	if err != nil {
		return err
	}
	buildProtoRepo := NewBuildPrototypeRepo(r.q)
	buildProtos, err := buildProtoRepo.FindAllPrototypes(ctx)
	if err != nil {
		return err
	}
	techProtoRepo := NewTechPrototypeRepo(r.q)
	techProtos, err := techProtoRepo.FindAllPrototypes(ctx)
	if err != nil {
		return err
	}
	storageProtoRepo := NewStoragePrototypeRepo(r.q)
	storageProtos, err := storageProtoRepo.FindAllPrototypes(ctx)
	if err != nil {
		return err
	}

	// Index by ID
	armyMap := make(map[int]*domain.ArmyItemPrototype, len(armyProtos))
	for _, p := range armyProtos {
		armyMap[p.ID] = p
	}
	buildMap := make(map[int]*domain.BuildItemPrototype, len(buildProtos))
	for _, p := range buildProtos {
		buildMap[p.ID] = p
	}
	techMap := make(map[int]*domain.TechItemPrototype, len(techProtos))
	for _, p := range techProtos {
		techMap[p.ID] = p
	}
	storageMap := make(map[int]*domain.StorageItemPrototype, len(storageProtos))
	for _, p := range storageProtos {
		storageMap[p.ID] = p
	}

	// Populate domain collections
	mappers.HydrateArmyItems(base, armyRows, armyMap)
	mappers.HydrateBuildItems(base, buildRows, buildMap)
	mappers.HydrateTechItems(base, techRows, techMap)
	mappers.HydrateStorageItems(base, storageRows, storageMap)

	return nil
}
