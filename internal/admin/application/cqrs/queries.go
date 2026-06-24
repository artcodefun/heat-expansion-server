package cqrs

import (
	"context"

	"github.com/artcodefun/heat-expansion-server/internal/admin/application/cqrs/readmodels"
	"github.com/google/uuid"
)

// AdminQueries encapsulates admin read-side operations.
type AdminQueries interface {
	GetProfile(ctx context.Context, actor Actor, adminID uuid.UUID) (*readmodels.AdminProfile, error)
}

// PackageQueries exposes full crystal package catalog reads.
type PackageQueries interface {
	ListCrystalPackages(ctx context.Context, actor Actor) ([]*readmodels.CrystalPackage, error)
	GetCrystalPackage(ctx context.Context, actor Actor, id uuid.UUID) (*readmodels.CrystalPackage, error)
}

// PrototypeQueries exposes full prototype catalog reads across all four types.
type PrototypeQueries interface {
	ListArmyPrototypes(ctx context.Context, actor Actor) ([]*readmodels.ArmyPrototype, error)
	GetArmyPrototype(ctx context.Context, actor Actor, id int64) (*readmodels.ArmyPrototype, error)

	ListBuildPrototypes(ctx context.Context, actor Actor) ([]*readmodels.BuildPrototype, error)
	GetBuildPrototype(ctx context.Context, actor Actor, id int64) (*readmodels.BuildPrototype, error)

	ListStoragePrototypes(ctx context.Context, actor Actor) ([]*readmodels.StoragePrototype, error)
	GetStoragePrototype(ctx context.Context, actor Actor, id int64) (*readmodels.StoragePrototype, error)

	ListTechPrototypes(ctx context.Context, actor Actor) ([]*readmodels.TechPrototype, error)
	GetTechPrototype(ctx context.Context, actor Actor, id int64) (*readmodels.TechPrototype, error)
}

// TranslationQueries exposes read access to the game translation catalog.
type TranslationQueries interface {
	ListTranslations(ctx context.Context, actor Actor) ([]*readmodels.Translation, error)
}
