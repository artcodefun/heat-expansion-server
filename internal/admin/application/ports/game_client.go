package ports

import (
	"context"

	"github.com/artcodefun/heat-expansion-server/internal/admin/application/cqrs/readmodels"
)

// GamePrivateClient abstracts outbound calls to the game module's private gRPC API.
type GamePrivateClient interface {
	ListArmyPrototypes(ctx context.Context) ([]*readmodels.ArmyPrototype, error)
	GetArmyPrototype(ctx context.Context, id int64) (*readmodels.ArmyPrototype, error)
	CreateArmyPrototype(ctx context.Context, p *readmodels.ArmyPrototype) (*readmodels.ArmyPrototype, error)
	UpdateArmyPrototype(ctx context.Context, p *readmodels.ArmyPrototype) (*readmodels.ArmyPrototype, error)

	ListBuildPrototypes(ctx context.Context) ([]*readmodels.BuildPrototype, error)
	GetBuildPrototype(ctx context.Context, id int64) (*readmodels.BuildPrototype, error)
	CreateBuildPrototype(ctx context.Context, p *readmodels.BuildPrototype) (*readmodels.BuildPrototype, error)
	UpdateBuildPrototype(ctx context.Context, p *readmodels.BuildPrototype) (*readmodels.BuildPrototype, error)

	ListStoragePrototypes(ctx context.Context) ([]*readmodels.StoragePrototype, error)
	GetStoragePrototype(ctx context.Context, id int64) (*readmodels.StoragePrototype, error)
	CreateStoragePrototype(ctx context.Context, p *readmodels.StoragePrototype) (*readmodels.StoragePrototype, error)
	UpdateStoragePrototype(ctx context.Context, p *readmodels.StoragePrototype) (*readmodels.StoragePrototype, error)

	ListTechPrototypes(ctx context.Context) ([]*readmodels.TechPrototype, error)
	GetTechPrototype(ctx context.Context, id int64) (*readmodels.TechPrototype, error)
	CreateTechPrototype(ctx context.Context, p *readmodels.TechPrototype) (*readmodels.TechPrototype, error)
	UpdateTechPrototype(ctx context.Context, p *readmodels.TechPrototype) (*readmodels.TechPrototype, error)

	UpsertTranslation(ctx context.Context, locale, key, value string) (*readmodels.Translation, error)
	ListTranslations(ctx context.Context) ([]*readmodels.Translation, error)
}
