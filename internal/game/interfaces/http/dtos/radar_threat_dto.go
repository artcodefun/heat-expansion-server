package dtos

import (
	"github.com/artcodefun/heat-expansion-api/internal/game/core/cqrs/readmodels"
)

type ThreatType string

const (
	ThreatTypeAttack ThreatType = "ATTACK"
	ThreatTypeSpy    ThreatType = "SPY"
)

type ThreatStatus string

const (
	ThreatStatusArriving ThreatStatus = "ARRIVING"
	ThreatStatusLost     ThreatStatus = "LOST"
	ThreatStatusArrived  ThreatStatus = "ARRIVED"
)

type RadarThreatDTO struct {
	ID                  string       `json:"id"`
	OperationID         int          `json:"operationId"`
	OwnerBaseID         int          `json:"ownerBaseId"`
	DetectedAt          int64        `json:"detectedAt"`
	DetectedCoordinates Vector2iDTO  `json:"detectedCoordinates"`
	SourceCoordinates   Vector2iDTO  `json:"source"`
	TargetCoordinates   Vector2iDTO  `json:"target"`
	EstimatedArrivalAt  int64        `json:"estimatedArrivalAt"`
	ArrivalAt           *int64       `json:"arrivalAt"`
	Type                ThreatType   `json:"type"`
	Status              ThreatStatus `json:"status"`
	Attack              int          `json:"attack"`
	Speed               int          `json:"speed"`
	Stealth             int          `json:"stealth"`
	Capacity            int          `json:"capacity"`
}

func RadarThreatFromReadModel(t *readmodels.RadarThreat) RadarThreatDTO {
	return RadarThreatDTO{
		ID:                  t.ID.String(),
		OperationID:         t.OperationID,
		OwnerBaseID:         t.OwnerBaseID,
		DetectedAt:          t.DetectedAt,
		DetectedCoordinates: Vector2iFromReadModel(t.DetectedCoordinates),
		SourceCoordinates:   Vector2iFromReadModel(t.SourceCoordinates),
		TargetCoordinates:   Vector2iFromReadModel(t.TargetCoordinates),
		EstimatedArrivalAt:  t.EstimatedArrivalAt,
		ArrivalAt:           t.ArrivalAt,
		Type:                ThreatType(t.Type),
		Status:              ThreatStatus(t.Status),
		Attack:              t.Attack,
		Speed:               t.Speed,
		Stealth:             t.Stealth,
		Capacity:            t.Capacity,
	}
}

func RadarThreatsFromReadModels(threats []*readmodels.RadarThreat) []RadarThreatDTO {
	out := make([]RadarThreatDTO, 0, len(threats))
	for _, t := range threats {
		out = append(out, RadarThreatFromReadModel(t))
	}
	return out
}
