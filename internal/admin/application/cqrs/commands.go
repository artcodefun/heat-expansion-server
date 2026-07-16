package cqrs

import (
	"context"

	"github.com/artcodefun/heat-expansion-server/internal/admin/application/ports"
)

// AdminCommands encapsulates admin authentication mutations.
type AdminCommands interface {
	// Register completes first-time setup using the invite token, sets a password,
	// and issues a session immediately. Returns the raw session token.
	Register(ctx context.Context, actor Actor, username, inviteToken, newPassword string) (string, error)
	// Login verifies credentials and issues a new session. Returns the raw session token.
	Login(ctx context.Context, actor Actor, username, password string) (string, error)
	// Logout revokes the session identified by the given bearer token.
	Logout(ctx context.Context, actor Actor, token string) error
}

// PackageCommands encapsulates privileged writes to the billing crystal package catalog.
type PackageCommands interface {
	CreateCrystalPackage(ctx context.Context, actor Actor, p *ports.CrystalPackage) (*ports.CrystalPackage, error)
	UpdateCrystalPackage(ctx context.Context, actor Actor, p *ports.CrystalPackage) (*ports.CrystalPackage, error)
}

// PrototypeCommands encapsulates privileged writes to game prototype catalogs.
// Calls are pre-authorized by the private gRPC key; Actor is carried for
// auditability but not used for access control.
type PrototypeCommands interface {
	CreateArmyPrototype(ctx context.Context, actor Actor, p *ports.ArmyPrototype) (*ports.ArmyPrototype, error)
	UpdateArmyPrototype(ctx context.Context, actor Actor, p *ports.ArmyPrototype) (*ports.ArmyPrototype, error)

	CreateBuildPrototype(ctx context.Context, actor Actor, p *ports.BuildPrototype) (*ports.BuildPrototype, error)
	UpdateBuildPrototype(ctx context.Context, actor Actor, p *ports.BuildPrototype) (*ports.BuildPrototype, error)

	CreateStoragePrototype(ctx context.Context, actor Actor, p *ports.StoragePrototype) (*ports.StoragePrototype, error)
	UpdateStoragePrototype(ctx context.Context, actor Actor, p *ports.StoragePrototype) (*ports.StoragePrototype, error)

	CreateTechPrototype(ctx context.Context, actor Actor, p *ports.TechPrototype) (*ports.TechPrototype, error)
	UpdateTechPrototype(ctx context.Context, actor Actor, p *ports.TechPrototype) (*ports.TechPrototype, error)
}

// TranslationCommands encapsulates privileged writes to the game translation catalog.
type TranslationCommands interface {
	UpsertTranslation(ctx context.Context, actor Actor, locale, key, value string) (*ports.Translation, error)
}
