package readmodels

import (
	"github.com/artcodefun/heat-expansion-server/internal/game/domain"
	"github.com/google/uuid"
)

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
	Name        domain.TranslationKey
	Description domain.TranslationKey
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

type ScanReportSourceType string

const (
	ScanReportSourceUnknown          ScanReportSourceType = "UNKNOWN"
	ScanReportSourceOperation        ScanReportSourceType = "OPERATION"
	ScanReportSourceScanner          ScanReportSourceType = "SCANNER"
	ScanReportSourceIntel            ScanReportSourceType = "INTEL"
	ScanReportSourceDiplomaticReveal ScanReportSourceType = "DIPLOMATIC_REVEAL"
)

// SectorOwner represents the public owner data attached to a scanned base location.
type SectorOwner struct {
	ID   uuid.UUID
	Name string
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
	Owner   *SectorOwner

	Info ScanInfo
	// If true, occupant intel was cloaked; only fallback empty sector details provided.
	IsCloaked  bool
	SourceType ScanReportSourceType
	SourceID   *uuid.UUID
}
