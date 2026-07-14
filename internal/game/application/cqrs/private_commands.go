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

// BuildPrototypeCommands encapsulates privileged writes to the build prototype catalog. No Actor; pre-authorized by caller.
type BuildPrototypeCommands interface {
	CreateBuildPrototype(ctx context.Context, proto *domain.BuildItemPrototype) (*domain.BuildItemPrototype, error)
	UpdateBuildPrototype(ctx context.Context, proto *domain.BuildItemPrototype) (*domain.BuildItemPrototype, error)
}

// StoragePrototypeCommands encapsulates privileged writes to the storage prototype catalog. No Actor; pre-authorized by caller.
type StoragePrototypeCommands interface {
	CreateStoragePrototype(ctx context.Context, proto *domain.StorageItemPrototype) (*domain.StorageItemPrototype, error)
	UpdateStoragePrototype(ctx context.Context, proto *domain.StorageItemPrototype) (*domain.StorageItemPrototype, error)
}

// TechPrototypeCommands encapsulates privileged writes to the tech prototype catalog. No Actor; pre-authorized by caller.
type TechPrototypeCommands interface {
	CreateTechPrototype(ctx context.Context, proto *domain.TechItemPrototype) (*domain.TechItemPrototype, error)
	UpdateTechPrototype(ctx context.Context, proto *domain.TechItemPrototype) (*domain.TechItemPrototype, error)
}
