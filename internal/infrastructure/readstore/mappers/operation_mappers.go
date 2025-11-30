package mappers

import (
	"github.com/artcodefun/heat-expansion-api/internal/core/cqrs/readmodels"
	"github.com/artcodefun/heat-expansion-api/internal/infrastructure/readstore/gen"
)

func OperationFromModel(m gen.MilitaryOperation) readmodels.MilitaryOperation {
	return readmodels.MilitaryOperation{
		ID:                int(m.ID),
		Type:              readmodels.MilitaryOperationType(m.Type),
		OwnerUserID:       int(m.OwnerUserID),
		SourceBaseID:      int(m.SourceBaseID),
		SourceCoordinates: readmodels.Vector2i{X: int(m.SourceX), Y: int(m.SourceY)},
		TargetCoordinates: readmodels.Vector2i{X: int(m.TargetX), Y: int(m.TargetY)},
		OutboundDepartAt:  m.OutboundDepartAt,
		OutboundArriveAt:  m.OutboundArriveAt,
		ReturnDepartAt:    m.ReturnDepartAt,
		ReturnArriveAt:    m.ReturnArriveAt,
		CompletedAt:       m.CompletedAt,
		Phase:             readmodels.MilitaryOperationPhase(m.Phase),
		Result:            readmodels.MilitaryOperationResult(m.Result),
		Units:             []readmodels.MilitaryUnit{},
		SpyResult:         nil,
		AttackResult:      nil,
	}
}
