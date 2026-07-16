package ports

import (
	"context"
)

// GamePrivateClient abstracts outbound calls to the game module's private gRPC API.
type GamePrivateClient interface {
	ListArmyPrototypes(ctx context.Context) ([]*ArmyPrototype, error)
	GetArmyPrototype(ctx context.Context, id int64) (*ArmyPrototype, error)
	CreateArmyPrototype(ctx context.Context, p *ArmyPrototype) (*ArmyPrototype, error)
	UpdateArmyPrototype(ctx context.Context, p *ArmyPrototype) (*ArmyPrototype, error)

	ListBuildPrototypes(ctx context.Context) ([]*BuildPrototype, error)
	GetBuildPrototype(ctx context.Context, id int64) (*BuildPrototype, error)
	CreateBuildPrototype(ctx context.Context, p *BuildPrototype) (*BuildPrototype, error)
	UpdateBuildPrototype(ctx context.Context, p *BuildPrototype) (*BuildPrototype, error)

	ListStoragePrototypes(ctx context.Context) ([]*StoragePrototype, error)
	GetStoragePrototype(ctx context.Context, id int64) (*StoragePrototype, error)
	CreateStoragePrototype(ctx context.Context, p *StoragePrototype) (*StoragePrototype, error)
	UpdateStoragePrototype(ctx context.Context, p *StoragePrototype) (*StoragePrototype, error)

	ListTechPrototypes(ctx context.Context) ([]*TechPrototype, error)
	GetTechPrototype(ctx context.Context, id int64) (*TechPrototype, error)
	CreateTechPrototype(ctx context.Context, p *TechPrototype) (*TechPrototype, error)
	UpdateTechPrototype(ctx context.Context, p *TechPrototype) (*TechPrototype, error)

	UpsertTranslation(ctx context.Context, locale, key, value string) (*Translation, error)
	ListTranslations(ctx context.Context) ([]*Translation, error)
}
