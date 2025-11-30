package readmodels

// LocationType represents the occupant classification at given coordinates.
// Derived from presence of user base / resource / dangerous location; persisted here temporarily.
type LocationType string

const (
	LocationTypeUserBase    LocationType = "BASE"
	LocationTypeEmpty       LocationType = "EMPTY"
	LocationTypeResourceful LocationType = "RESOURCEFUL"
	LocationTypeDangerous   LocationType = "DANGEROUS"
)

// Vector2i represents coordinates (x, y).
type Vector2i struct {
	X int
	Y int
}

// SectorModel represents static sector info.
type SectorModel struct {
	Coordinates Vector2i
	// Merged empty-location flavor fields for deterministic generation/persistence
	Details LocationDetails
}

type LocationDetails struct {
	Name        string
	Description string
	ImageURL    string
}

// ScanInfo represents scan information for a sector.
type ScanInfo struct {
	Credits    int
	Iron       int
	Titanium   int
	Antimatter int
	Defence    int
	Attack     int
	Space      int
}

// SectorScanReport captures a user's scan snapshot of a sector at a moment in time.
// Produced after successful military operations (attack/spy) or other scan actions.
type SectorScanReport struct {
	ID          int
	BaseID      int
	Coordinates Vector2i
	CreatedAt   int64

	Details LocationDetails
	Type    LocationType

	Info ScanInfo
	// If true, occupant intel was cloaked; only fallback empty sector details provided.
	IsCloaked bool
	// Optional link to the operation that produced this report (for traceability/idempotency)
	SourceOperationID int
}
