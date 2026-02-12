package ports

import "github.com/google/uuid"

// MoveBuildQueueJob is a serializable domain job to move the build queue forward.
type MoveBuildQueueJob struct {
	BaseID int
}

// MoveArmyQueueJob is a serializable domain job to move the army queue forward.
type MoveArmyQueueJob struct {
	BaseID int
}

// MoveTechQueueJob is a serializable domain job to move the tech queue forward.
type MoveTechQueueJob struct {
	BaseID int
}

// DeleteExpiredBuffJob is a serializable domain job to cleanup an expired buff.
type DeleteExpiredBuffJob struct {
	BaseID int
	ItemID uuid.UUID
}

// RestoreDamagedItemJob is a serializable domain job to cleanup a finished restoration.
type RestoreDamagedItemJob struct {
	BaseID int
	ItemID uuid.UUID
}

// DecryptIntelItemJob is a serializable domain job to cleanup a finished decryption.
type DecryptIntelItemJob struct {
	BaseID int
	ItemID uuid.UUID
}

// UpdateMilitaryOperationJob asks the system to advance an operation's phase based on time
// (e.g., arrival at target, arrival back at source). Safe to enqueue multiple times.
type UpdateMilitaryOperationJob struct {
	OperationID int
}

// SpawnNearbyLocationsJob triggers spawning of resourceful/dangerous locations near a specific user base.
// The job handler is responsible for rescheduling itself for that specific base.
type SpawnNearbyLocationsJob struct {
	BaseID int
}

// IntelligenceScanJob asks the system to perform a periodic scan for a scanner building.
type IntelligenceScanJob struct {
	BaseID     int
	BuildingID uuid.UUID
}

// IntelligenceRadarJob asks the system to detect incoming threats for a radar building.
type IntelligenceRadarJob struct {
	BaseID      int
	OperationID int
}

// SchadulableJob represents a generic job that can be scheduled for execution.
// This type can be used to define any job or task that needs to be managed by a scheduler.
type SchadulableJob any

// Scheduler defines the interface for scheduling domain actions at a future time.
type Scheduler interface {
	// Schedule schedules a domain job (payload struct) to be executed at the specified Unix timestamp.
	Schedule(job SchadulableJob, executeAt int64) error
}
