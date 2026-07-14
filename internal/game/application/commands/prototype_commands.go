package commands

import (
	"context"

	"github.com/artcodefun/heat-expansion-server/internal/game/application/ports"
	"github.com/artcodefun/heat-expansion-server/internal/game/domain"
)

type PrototypeCommands struct {
	ArmyRepo    ports.ArmyPrototypeRepository
	BuildRepo   ports.BuildPrototypeRepository
	StorageRepo ports.StoragePrototypeRepository
	TechRepo    ports.TechPrototypeRepository
	TxMgr       ports.TransactionManager
}

func NewPrototypeCommands(armyRepo ports.ArmyPrototypeRepository, buildRepo ports.BuildPrototypeRepository, storageRepo ports.StoragePrototypeRepository, techRepo ports.TechPrototypeRepository, txMgr ports.TransactionManager) *PrototypeCommands {
	return &PrototypeCommands{ArmyRepo: armyRepo, BuildRepo: buildRepo, StorageRepo: storageRepo, TechRepo: techRepo, TxMgr: txMgr}
}

// CreateArmyPrototype inserts a new army prototype using the caller-supplied ID
// (prototypes are ordered and grouped into id ranges by type) and returns it.
// A clash with an existing ID surfaces as a conflict error.
func (c *PrototypeCommands) CreateArmyPrototype(ctx context.Context, proto *domain.ArmyItemPrototype) (*domain.ArmyItemPrototype, error) {
	err := c.TxMgr.WithTx(ctx, func(tx ports.Transaction) error {
		return c.ArmyRepo.Tx(tx).CreatePrototype(ctx, proto)
	})
	if err != nil {
		return nil, repoErr(err)
	}
	return proto, nil
}

// UpdateArmyPrototype overwrites an existing army prototype identified by
// proto.ID, returning the updated aggregate.
func (c *PrototypeCommands) UpdateArmyPrototype(ctx context.Context, proto *domain.ArmyItemPrototype) (*domain.ArmyItemPrototype, error) {
	err := c.TxMgr.WithTx(ctx, func(tx ports.Transaction) error {
		return c.ArmyRepo.Tx(tx).UpdatePrototype(ctx, proto)
	})
	if err != nil {
		return nil, repoErr(err)
	}
	return proto, nil
}

// CreateBuildPrototype inserts a new build prototype using the caller-supplied ID
// (prototypes are ordered and grouped into id ranges by type) and returns it.
func (c *PrototypeCommands) CreateBuildPrototype(ctx context.Context, proto *domain.BuildItemPrototype) (*domain.BuildItemPrototype, error) {
	err := c.TxMgr.WithTx(ctx, func(tx ports.Transaction) error {
		return c.BuildRepo.Tx(tx).CreatePrototype(ctx, proto)
	})
	if err != nil {
		return nil, repoErr(err)
	}
	return proto, nil
}

// UpdateBuildPrototype overwrites an existing build prototype identified by
// proto.ID, returning the updated aggregate.
func (c *PrototypeCommands) UpdateBuildPrototype(ctx context.Context, proto *domain.BuildItemPrototype) (*domain.BuildItemPrototype, error) {
	err := c.TxMgr.WithTx(ctx, func(tx ports.Transaction) error {
		return c.BuildRepo.Tx(tx).UpdatePrototype(ctx, proto)
	})
	if err != nil {
		return nil, repoErr(err)
	}
	return proto, nil
}

// CreateStoragePrototype inserts a new storage prototype using the caller-supplied ID
// (prototypes are ordered and grouped into id ranges by type) and returns it.
func (c *PrototypeCommands) CreateStoragePrototype(ctx context.Context, proto *domain.StorageItemPrototype) (*domain.StorageItemPrototype, error) {
	err := c.TxMgr.WithTx(ctx, func(tx ports.Transaction) error {
		return c.StorageRepo.Tx(tx).CreatePrototype(ctx, proto)
	})
	if err != nil {
		return nil, repoErr(err)
	}
	return proto, nil
}

// UpdateStoragePrototype overwrites an existing storage prototype identified by
// proto.ID, returning the updated aggregate.
func (c *PrototypeCommands) UpdateStoragePrototype(ctx context.Context, proto *domain.StorageItemPrototype) (*domain.StorageItemPrototype, error) {
	err := c.TxMgr.WithTx(ctx, func(tx ports.Transaction) error {
		return c.StorageRepo.Tx(tx).UpdatePrototype(ctx, proto)
	})
	if err != nil {
		return nil, repoErr(err)
	}
	return proto, nil
}

// CreateTechPrototype inserts a new tech prototype using the caller-supplied ID
// (prototypes are ordered and grouped into id ranges by type) and returns it.
func (c *PrototypeCommands) CreateTechPrototype(ctx context.Context, proto *domain.TechItemPrototype) (*domain.TechItemPrototype, error) {
	err := c.TxMgr.WithTx(ctx, func(tx ports.Transaction) error {
		return c.TechRepo.Tx(tx).CreatePrototype(ctx, proto)
	})
	if err != nil {
		return nil, repoErr(err)
	}
	return proto, nil
}

// UpdateTechPrototype overwrites an existing tech prototype identified by
// proto.ID, returning the updated aggregate.
func (c *PrototypeCommands) UpdateTechPrototype(ctx context.Context, proto *domain.TechItemPrototype) (*domain.TechItemPrototype, error) {
	err := c.TxMgr.WithTx(ctx, func(tx ports.Transaction) error {
		return c.TechRepo.Tx(tx).UpdatePrototype(ctx, proto)
	})
	if err != nil {
		return nil, repoErr(err)
	}
	return proto, nil
}
