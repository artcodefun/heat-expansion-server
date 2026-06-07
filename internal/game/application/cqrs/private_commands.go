package cqrs

import (
	"context"

	"github.com/artcodefun/heat-expansion-server/internal/game/domain"
)

// ArmyPrototypeCommands encapsulates privileged writes to the army prototype catalog. No Actor; pre-authorized by caller.
type ArmyPrototypeCommands interface {
	CreateArmyPrototype(ctx context.Context, proto *domain.ArmyItemPrototype) (*domain.ArmyItemPrototype, error)
	UpdateArmyPrototype(ctx context.Context, proto *domain.ArmyItemPrototype) (*domain.ArmyItemPrototype, error)
}
