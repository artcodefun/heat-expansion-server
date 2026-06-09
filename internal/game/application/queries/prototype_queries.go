package queries

import (
	"context"

	"github.com/artcodefun/heat-expansion-server/internal/game/application/cqrs/readmodels"
	"github.com/artcodefun/heat-expansion-server/internal/game/application/ports"
)

// PrototypeQueries serves raw prototype catalog reads for the private API.
type PrototypeQueries struct {
	ArmyRepo  ports.ArmyPrototypeReadRepository
	BuildRepo ports.BuildPrototypeReadRepository
}

func NewPrototypeQueries(armyRepo ports.ArmyPrototypeReadRepository, buildRepo ports.BuildPrototypeReadRepository) *PrototypeQueries {
	return &PrototypeQueries{ArmyRepo: armyRepo, BuildRepo: buildRepo}
}

func (q *PrototypeQueries) ListArmyPrototypes(ctx context.Context) ([]*readmodels.ArmyItemPrototype, error) {
	protos, err := q.ArmyRepo.ListArmyPrototypes(ctx)
	if err != nil {
		return nil, repoErr(err)
	}
	return protos, nil
}

func (q *PrototypeQueries) GetArmyPrototype(ctx context.Context, id int) (*readmodels.ArmyItemPrototype, error) {
	proto, err := q.ArmyRepo.GetArmyPrototype(ctx, id)
	if err != nil {
		return nil, repoErr(err)
	}
	return proto, nil
}

func (q *PrototypeQueries) ListBuildPrototypes(ctx context.Context) ([]*readmodels.BuildItemPrototype, error) {
	protos, err := q.BuildRepo.ListBuildPrototypes(ctx)
	if err != nil {
		return nil, repoErr(err)
	}
	return protos, nil
}

func (q *PrototypeQueries) GetBuildPrototype(ctx context.Context, id int) (*readmodels.BuildItemPrototype, error) {
	proto, err := q.BuildRepo.GetBuildPrototype(ctx, id)
	if err != nil {
		return nil, repoErr(err)
	}
	return proto, nil
}
