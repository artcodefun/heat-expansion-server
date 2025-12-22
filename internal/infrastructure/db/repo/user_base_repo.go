package repo

import (
	"context"
	"database/sql"
	"errors"

	"github.com/artcodefun/heat-expansion-api/internal/core/domain"
	"github.com/artcodefun/heat-expansion-api/internal/core/ports"
	"github.com/artcodefun/heat-expansion-api/internal/infrastructure/db/gen"
	"github.com/artcodefun/heat-expansion-api/internal/infrastructure/db/mappers"
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

func (r *UserBaseRepo) Create(base *domain.UserBaseModel) error {
	row, err := r.q.CreateBase(context.Background(), mappers.InsertBaseParamsFromDomain(base))
	if err != nil {
		return err
	}
	created := mappers.UserBaseFromDB(row)
	base.ID = created.ID
	base.Stats = created.Stats
	base.Name = created.Name
	base.Description = created.Description
	base.ImageURL = created.ImageURL
	return nil
}

func (r *UserBaseRepo) FindByID(id int) (*domain.UserBaseModel, error) {
	row, err := r.q.GetBaseByID(context.Background(), int64(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ports.ErrNotFound
		}
		return nil, err
	}
	base := mappers.UserBaseFromDB(row)
	_ = r.hydrateBase(base)
	return base, nil
}

func (r *UserBaseRepo) FindByIDForUpdate(id int) (*domain.UserBaseModel, error) {
	row, err := r.q.GetBaseByIDForUpdate(context.Background(), int64(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ports.ErrNotFound
		}
		return nil, err
	}
	base := mappers.UserBaseFromDB(row)
	_ = r.hydrateBase(base)
	return base, nil
}

// Update replaces all per-item rows (army/build/tech/storage) with the current aggregate state and updates base stats.
func (r *UserBaseRepo) Update(base *domain.UserBaseModel) error {
	// Update base stats and metadata first.
	if _, err := r.q.UpdateBase(context.Background(), mappers.UpdateBaseParamsFromDomain(base)); err != nil {
		return err
	}
	if err := r.persistArmyItems(base); err != nil {
		return err
	}
	if err := r.persistBuildItems(base); err != nil {
		return err
	}
	if err := r.persistTechItems(base); err != nil {
		return err
	}
	if err := r.persistStorageItems(base); err != nil {
		return err
	}
	return nil
}

func (r *UserBaseRepo) Delete(id int) error {
	return r.q.DeleteBase(context.Background(), int64(id))
}

func (r *UserBaseRepo) FindByUserID(userID int) ([]*domain.UserBaseModel, error) {
	rows, err := r.q.ListBasesByUserID(context.Background(), int64(userID))
	if err != nil {
		return nil, err
	}
	out := make([]*domain.UserBaseModel, 0, len(rows))
	for _, b := range rows {
		base := mappers.UserBaseFromDB(b)
		_ = r.hydrateBase(base)
		out = append(out, base)
	}
	return out, nil
}

func (r *UserBaseRepo) FindByCoordinates(x, y int) (*domain.UserBaseModel, error) {
	row, err := r.q.GetBaseByCoordinates(context.Background(), gen.GetBaseByCoordinatesParams{SectorX: int32(x), SectorY: int32(y)})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ports.ErrNotFound
		}
		return nil, err
	}
	base := mappers.UserBaseFromDB(row)
	_ = r.hydrateBase(base)
	return base, nil
}

func (r *UserBaseRepo) FindByCoordinatesForUpdate(x, y int) (*domain.UserBaseModel, error) {
	row, err := r.q.GetBaseByCoordinatesForUpdate(context.Background(), gen.GetBaseByCoordinatesForUpdateParams{SectorX: int32(x), SectorY: int32(y)})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ports.ErrNotFound
		}
		return nil, err
	}
	base := mappers.UserBaseFromDB(row)
	_ = r.hydrateBase(base)
	return base, nil
}

func (r *UserBaseRepo) FindAll() ([]*domain.UserBaseModel, error) {
	rows, err := r.q.ListAllBases(context.Background())
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
func (r *UserBaseRepo) GetOwnerID(baseID int) (int, error) {
	row, err := r.q.GetBaseByID(context.Background(), int64(baseID))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ports.ErrNotFound
		}
		return 0, err
	}
	return int(row.UserID), nil
}

// PersistAggregate updates the base stats row and fully replaces all item rows with the current aggregate collections.
// This is a bulk-replacement write path keeping persistence logic simple and avoiding per-row diffing.

// --- Internal per-table persistence helpers ---

func (r *UserBaseRepo) persistArmyItems(base *domain.UserBaseModel) error {
	ctx := context.Background()
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

func (r *UserBaseRepo) persistBuildItems(base *domain.UserBaseModel) error {
	ctx := context.Background()
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

func (r *UserBaseRepo) persistTechItems(base *domain.UserBaseModel) error {
	ctx := context.Background()
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

func (r *UserBaseRepo) persistStorageItems(base *domain.UserBaseModel) error {
	ctx := context.Background()
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
func (r *UserBaseRepo) hydrateBase(base *domain.UserBaseModel) error {
	ctx := context.Background()

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
	armyProtos, err := armyProtoRepo.FindAllPrototypes()
	if err != nil {
		return err
	}
	buildProtoRepo := NewBuildPrototypeRepo(r.q)
	buildProtos, err := buildProtoRepo.FindAllPrototypes()
	if err != nil {
		return err
	}
	techProtoRepo := NewTechPrototypeRepo(r.q)
	techProtos, err := techProtoRepo.FindAllPrototypes()
	if err != nil {
		return err
	}
	storageProtoRepo := NewStoragePrototypeRepo(r.q)
	storageProtos, err := storageProtoRepo.FindAllPrototypes()
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
