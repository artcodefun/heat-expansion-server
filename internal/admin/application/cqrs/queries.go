package cqrs

import (
	"context"

	"github.com/artcodefun/heat-expansion-server/internal/admin/application/cqrs/readmodels"
	"github.com/artcodefun/heat-expansion-server/internal/admin/application/ports"
	"github.com/google/uuid"
)

// AdminQueries encapsulates admin read-side operations.
type AdminQueries interface {
	GetProfile(ctx context.Context, actor Actor, adminID uuid.UUID) (*readmodels.AdminProfile, error)
}

// PackageQueries exposes full crystal package catalog reads.
type PackageQueries interface {
	ListCrystalPackages(ctx context.Context, actor Actor) ([]*ports.CrystalPackage, error)
	GetCrystalPackage(ctx context.Context, actor Actor, id uuid.UUID) (*ports.CrystalPackage, error)
}

// PrototypeQueries exposes full prototype catalog reads across all four types.
type PrototypeQueries interface {
	ListArmyPrototypes(ctx context.Context, actor Actor) ([]*ports.ArmyPrototype, error)
	GetArmyPrototype(ctx context.Context, actor Actor, id int64) (*ports.ArmyPrototype, error)

	ListBuildPrototypes(ctx context.Context, actor Actor) ([]*ports.BuildPrototype, error)
	GetBuildPrototype(ctx context.Context, actor Actor, id int64) (*ports.BuildPrototype, error)

	ListStoragePrototypes(ctx context.Context, actor Actor) ([]*ports.StoragePrototype, error)
	GetStoragePrototype(ctx context.Context, actor Actor, id int64) (*ports.StoragePrototype, error)

	ListTechPrototypes(ctx context.Context, actor Actor) ([]*ports.TechPrototype, error)
	GetTechPrototype(ctx context.Context, actor Actor, id int64) (*ports.TechPrototype, error)
}

// TranslationQueries exposes read access to the game translation catalog.
type TranslationQueries interface {
	ListTranslations(ctx context.Context, actor Actor) ([]*ports.Translation, error)
}
