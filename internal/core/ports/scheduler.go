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

// DeleteExpiredBuffJob is a serializable domain job to delete an expired buff item.
type DeleteExpiredBuffJob struct {
	BaseID int
	ItemID uuid.UUID
}

// UpdateMilitaryOperationJob asks the system to advance an operation's phase based on time
// (e.g., arrival at target, arrival back at source). Safe to enqueue multiple times.
type UpdateMilitaryOperationJob struct {
	OperationID int
}

// SpawnNearbyLocationsJob triggers spawning of resourceful/dangerous locations near a random user base.
// The job handler is responsible for rescheduling itself.
type SpawnNearbyLocationsJob struct{}

// RadarScanJob asks the system to perform a radar scan for a specific radar building of a base.
// The job is idempotent: if the building no longer exists, it should no-op and not reschedule.
type RadarScanJob struct {
	BaseID     int
	BuildingID uuid.UUID
}

// Scheduler defines the interface for scheduling domain actions at a future time.
type Scheduler interface {
	// Schedule schedules a domain job (payload struct) to be executed at the specified Unix timestamp.
	Schedule(job interface{}, executeAt int64) error
}
