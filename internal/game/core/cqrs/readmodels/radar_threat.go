package readmodels

import "github.com/google/uuid"

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

type RadarThreat struct {
	ID                  uuid.UUID
	OperationID         int
	OwnerBaseID         int
	DetectedAt          int64
	DetectedCoordinates Vector2i
	SourceCoordinates   Vector2i
	TargetCoordinates   Vector2i
	EstimatedArrivalAt  int64
	ArrivalAt           *int64
	Type                ThreatType
	Status              ThreatStatus
	Attack              int
	Speed               int
	Stealth             int
	Capacity            int
}
