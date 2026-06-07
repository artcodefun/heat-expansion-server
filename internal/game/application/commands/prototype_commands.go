package commands

import (
	"context"

	"github.com/artcodefun/heat-expansion-server/internal/game/application/ports"
	"github.com/artcodefun/heat-expansion-server/internal/game/domain"
)

type PrototypeCommands struct {
	ArmyRepo ports.ArmyPrototypeRepository
	TxMgr    ports.TransactionManager
}

func NewPrototypeCommands(armyRepo ports.ArmyPrototypeRepository, txMgr ports.TransactionManager) *PrototypeCommands {
	return &PrototypeCommands{ArmyRepo: armyRepo, TxMgr: txMgr}
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
