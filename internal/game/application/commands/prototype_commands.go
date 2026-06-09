package commands

import (
	"context"

	"github.com/artcodefun/heat-expansion-server/internal/game/application/ports"
	"github.com/artcodefun/heat-expansion-server/internal/game/domain"
)

type PrototypeCommands struct {
	ArmyRepo  ports.ArmyPrototypeRepository
	BuildRepo ports.BuildPrototypeRepository
	TxMgr     ports.TransactionManager
}

func NewPrototypeCommands(armyRepo ports.ArmyPrototypeRepository, buildRepo ports.BuildPrototypeRepository, txMgr ports.TransactionManager) *PrototypeCommands {
	return &PrototypeCommands{ArmyRepo: armyRepo, BuildRepo: buildRepo, TxMgr: txMgr}
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
