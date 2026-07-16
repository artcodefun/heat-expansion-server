package commands

import (
	"context"

	"github.com/artcodefun/heat-expansion-server/internal/admin/application/cqrs"
	"github.com/artcodefun/heat-expansion-server/internal/admin/application/ports"
)

// PrototypeCommands delegates prototype writes to the game private gRPC API.
type PrototypeCommands struct {
	game ports.GamePrivateClient
}

func NewPrototypeCommands(game ports.GamePrivateClient) *PrototypeCommands {
	return &PrototypeCommands{game: game}
}

func (c *PrototypeCommands) CreateArmyPrototype(ctx context.Context, _ cqrs.Actor, p *ports.ArmyPrototype) (*ports.ArmyPrototype, error) {
	created, err := c.game.CreateArmyPrototype(ctx, p)
	return created, clientErr(err)
}

func (c *PrototypeCommands) UpdateArmyPrototype(ctx context.Context, _ cqrs.Actor, p *ports.ArmyPrototype) (*ports.ArmyPrototype, error) {
	updated, err := c.game.UpdateArmyPrototype(ctx, p)
	return updated, clientErr(err)
}

func (c *PrototypeCommands) CreateBuildPrototype(ctx context.Context, _ cqrs.Actor, p *ports.BuildPrototype) (*ports.BuildPrototype, error) {
	created, err := c.game.CreateBuildPrototype(ctx, p)
	return created, clientErr(err)
}

func (c *PrototypeCommands) UpdateBuildPrototype(ctx context.Context, _ cqrs.Actor, p *ports.BuildPrototype) (*ports.BuildPrototype, error) {
	updated, err := c.game.UpdateBuildPrototype(ctx, p)
	return updated, clientErr(err)
}

func (c *PrototypeCommands) CreateStoragePrototype(ctx context.Context, _ cqrs.Actor, p *ports.StoragePrototype) (*ports.StoragePrototype, error) {
	created, err := c.game.CreateStoragePrototype(ctx, p)
	return created, clientErr(err)
}

func (c *PrototypeCommands) UpdateStoragePrototype(ctx context.Context, _ cqrs.Actor, p *ports.StoragePrototype) (*ports.StoragePrototype, error) {
	updated, err := c.game.UpdateStoragePrototype(ctx, p)
	return updated, clientErr(err)
}

func (c *PrototypeCommands) CreateTechPrototype(ctx context.Context, _ cqrs.Actor, p *ports.TechPrototype) (*ports.TechPrototype, error) {
	created, err := c.game.CreateTechPrototype(ctx, p)
	return created, clientErr(err)
}

func (c *PrototypeCommands) UpdateTechPrototype(ctx context.Context, _ cqrs.Actor, p *ports.TechPrototype) (*ports.TechPrototype, error) {
	updated, err := c.game.UpdateTechPrototype(ctx, p)
	return updated, clientErr(err)
}
