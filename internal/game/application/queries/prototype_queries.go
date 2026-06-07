package queries

import (
	"context"

	"github.com/artcodefun/heat-expansion-server/internal/game/application/cqrs/readmodels"
	"github.com/artcodefun/heat-expansion-server/internal/game/application/ports"
)

// PrototypeQueries serves raw prototype catalog reads for the private API.
type PrototypeQueries struct {
	ArmyRepo ports.ArmyPrototypeReadRepository
}

func NewPrototypeQueries(armyRepo ports.ArmyPrototypeReadRepository) *PrototypeQueries {
	return &PrototypeQueries{ArmyRepo: armyRepo}
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
