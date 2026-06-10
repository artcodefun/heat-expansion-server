package cqrs

import (
	"context"

	"github.com/artcodefun/heat-expansion-server/internal/game/application/cqrs/readmodels"
)

// ArmyPrototypeQueries exposes raw army prototype catalog reads, without per-base availability filtering.
type ArmyPrototypeQueries interface {
	ListArmyPrototypes(ctx context.Context) ([]*readmodels.ArmyItemPrototype, error)
	GetArmyPrototype(ctx context.Context, id int) (*readmodels.ArmyItemPrototype, error)
}

// BuildPrototypeQueries exposes raw build prototype catalog reads, without per-base availability filtering.
type BuildPrototypeQueries interface {
	ListBuildPrototypes(ctx context.Context) ([]*readmodels.BuildItemPrototype, error)
	GetBuildPrototype(ctx context.Context, id int) (*readmodels.BuildItemPrototype, error)
}

// StoragePrototypeQueries exposes raw storage prototype catalog reads, without per-base availability filtering.
type StoragePrototypeQueries interface {
	ListStoragePrototypes(ctx context.Context) ([]*readmodels.StorageItemPrototype, error)
	GetStoragePrototype(ctx context.Context, id int) (*readmodels.StorageItemPrototype, error)
}

// TechPrototypeQueries exposes raw tech prototype catalog reads, without per-base availability filtering.
type TechPrototypeQueries interface {
	ListTechPrototypes(ctx context.Context) ([]*readmodels.TechItemPrototype, error)
	GetTechPrototype(ctx context.Context, id int) (*readmodels.TechItemPrototype, error)
}
