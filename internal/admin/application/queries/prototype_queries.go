package queries

import (
	"context"

	"github.com/artcodefun/heat-expansion-server/internal/admin/application/cqrs"
	"github.com/artcodefun/heat-expansion-server/internal/admin/application/cqrs/readmodels"
	"github.com/artcodefun/heat-expansion-server/internal/admin/application/ports"
)

// PrototypeQueries delegates prototype reads to the game private gRPC API.
type PrototypeQueries struct {
	game ports.GamePrivateClient
}

func NewPrototypeQueries(game ports.GamePrivateClient) *PrototypeQueries {
	return &PrototypeQueries{game: game}
}

func (q *PrototypeQueries) ListArmyPrototypes(ctx context.Context, _ cqrs.Actor) ([]*readmodels.ArmyPrototype, error) {
	list, err := q.game.ListArmyPrototypes(ctx)
	return list, clientErr(err)
}

func (q *PrototypeQueries) GetArmyPrototype(ctx context.Context, _ cqrs.Actor, id int64) (*readmodels.ArmyPrototype, error) {
	p, err := q.game.GetArmyPrototype(ctx, id)
	return p, clientErr(err)
}

func (q *PrototypeQueries) ListBuildPrototypes(ctx context.Context, _ cqrs.Actor) ([]*readmodels.BuildPrototype, error) {
	list, err := q.game.ListBuildPrototypes(ctx)
	return list, clientErr(err)
}

func (q *PrototypeQueries) GetBuildPrototype(ctx context.Context, _ cqrs.Actor, id int64) (*readmodels.BuildPrototype, error) {
	p, err := q.game.GetBuildPrototype(ctx, id)
	return p, clientErr(err)
}

func (q *PrototypeQueries) ListStoragePrototypes(ctx context.Context, _ cqrs.Actor) ([]*readmodels.StoragePrototype, error) {
	list, err := q.game.ListStoragePrototypes(ctx)
	return list, clientErr(err)
}

func (q *PrototypeQueries) GetStoragePrototype(ctx context.Context, _ cqrs.Actor, id int64) (*readmodels.StoragePrototype, error) {
	p, err := q.game.GetStoragePrototype(ctx, id)
	return p, clientErr(err)
}

func (q *PrototypeQueries) ListTechPrototypes(ctx context.Context, _ cqrs.Actor) ([]*readmodels.TechPrototype, error) {
	list, err := q.game.ListTechPrototypes(ctx)
	return list, clientErr(err)
}

func (q *PrototypeQueries) GetTechPrototype(ctx context.Context, _ cqrs.Actor, id int64) (*readmodels.TechPrototype, error) {
	p, err := q.game.GetTechPrototype(ctx, id)
	return p, clientErr(err)
}
