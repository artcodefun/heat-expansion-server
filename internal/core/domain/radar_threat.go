package domain

import (
	"github.com/google/uuid"
)

type ThreatStatus string

const (
	ThreatStatusArriving ThreatStatus = "ARRIVING"
	ThreatStatusLost     ThreatStatus = "LOST"
	ThreatStatusArrived  ThreatStatus = "ARRIVED"
)

type ThreatType string

const (
	ThreatTypeAttack ThreatType = "ATTACK"
	ThreatTypeSpy    ThreatType = "SPY"
)

// RadarThreat is an aggregate representing an incoming threat detected by radars.
type RadarThreat struct {
	EventProducer
	ID                 uuid.UUID
	OperationID        int
	OwnerBaseID        int
	DetectedAt         int64
	SourceCoordinates  Vector2i
	TargetCoordinates  Vector2i
	EstimatedArrivalAt int64
	ArrivalAt          *int64
	Type               ThreatType
	Status             ThreatStatus
	Attack             int
	Speed              int
	Stealth            int
	Capacity           int
}

func NewRadarThreat(op *MilitaryOperation, ownerBaseID int) *RadarThreat {
	rt := &RadarThreat{
		ID:                 uuid.Must(uuid.NewV7()),
		OperationID:        op.ID,
		OwnerBaseID:        ownerBaseID,
		DetectedAt:         NowUnix(),
		SourceCoordinates:  op.SourceCoordinates,
		TargetCoordinates:  op.TargetCoordinates,
		EstimatedArrivalAt: op.OutboundArriveAt,
		ArrivalAt:          nil,
		Type: func() ThreatType {
			if op.Type == MilitaryOperationTypeSpy {
				return ThreatTypeSpy
			}
			return ThreatTypeAttack
		}(),
		Status:   ThreatStatusArriving,
		Attack:   op.TotalAttack(),
		Speed:    op.TotalSpeed(),
		Stealth:  op.TotalStealth(),
		Capacity: op.TotalCapacity(),
	}

	rt.AddEvent(NewRadarThreatDetectedEvent(rt.ID, rt.OwnerBaseID, rt.OperationID))

	return rt
}

func (t *RadarThreat) UpdateArrivalAt(newArrivalAt int64) {
	if t.Status != ThreatStatusArriving {
		return
	}
	t.EstimatedArrivalAt = newArrivalAt
}

func (t *RadarThreat) MarkLost() {
	if t.Status != ThreatStatusArriving {
		return
	}
	t.Status = ThreatStatusLost
}

func (t *RadarThreat) MarkArrived() {
	if t.Status != ThreatStatusArriving {
		return
	}
	t.Status = ThreatStatusArrived
	now := NowUnix()
	t.ArrivalAt = &now
}
