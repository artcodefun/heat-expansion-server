package domain

import (
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
	EventProducer
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
	// Optional link to the scanner that produced this report (for traceability/idempotency)
	SourceScannerID *uuid.UUID
	// Optional link to the intel item that was decrypted to produce this report (for traceability/activity logging)
	SourceIntelItemID *uuid.UUID
}

// EmitCreated records a domain event indicating this report has been created/persisted for its owning base.
// Uses r.BaseID as the source base that initiated the scan-producing action (e.g., spy/attack operation).
func (r *SectorScanReport) EmitCreated() {
	// Only emit if we have a valid identifier; repository Create should set r.ID before calling this.
	if r == nil || r.ID <= 0 {
		return
	}
	r.AddEvent(NewScanReportCreatedEvent(r.ID, r.BaseID, r.SourceOperationID))
}

// NewSectorScanReportFromUserBase builds a scan report from a defender user base stats.
func NewSectorScanReportFromUserBase(baseID int, coords Vector2i, base *UserBaseModel, emptyDetails LocationDetails) *SectorScanReport {
	if base == nil {
		// Cloaked: provide fallback empty sector flavor
		return &SectorScanReport{BaseID: baseID, Coordinates: coords, CreatedAt: NowUnix(), IsCloaked: true, Type: LocationTypeUserBase, Details: emptyDetails}
	}
	info := ScanInfo{
		Credits:    base.Stats.Credits,
		Iron:       base.Stats.Iron,
		Titanium:   base.Stats.Titanium,
		Antimatter: base.Stats.Antimatter,
		Defence:    base.Stats.Defence,
		Attack:     base.Stats.Attack,
		Space:      base.Stats.Space,
	}
	return &SectorScanReport{BaseID: baseID, Coordinates: coords, CreatedAt: NowUnix(), Info: info, IsCloaked: false, Type: LocationTypeUserBase, Details: base.LocationDetails}
}

// NewSectorScanReportFromResourceLocation builds a scan report from a resource location snapshot.
func NewSectorScanReportFromResourceLocation(baseID int, coords Vector2i, loc *ResourceLocationModel, emptyDetails LocationDetails) *SectorScanReport {
	if loc == nil {
		return &SectorScanReport{BaseID: baseID, Coordinates: coords, CreatedAt: NowUnix(), IsCloaked: true, Type: LocationTypeResourceful, Details: emptyDetails}
	}
	armySnaps := loc.MaterializeDefenderArmySnapshot()
	structSnaps := loc.MaterializeDefenderStructureSnapshot()
	defence := sumDefence(armySnaps) + sumStructureDefence(structSnaps)
	attack := sumAttack(armySnaps)
	info := ScanInfo{
		Credits:    loc.Resources.Credits,
		Iron:       loc.Resources.Iron,
		Titanium:   loc.Resources.Titanium,
		Antimatter: loc.Resources.Antimatter,
		Defence:    defence,
		Attack:     attack,
		Space:      0,
	}
	return &SectorScanReport{BaseID: baseID, Coordinates: coords, CreatedAt: NowUnix(), Info: info, IsCloaked: false, Type: LocationTypeResourceful, Details: loc.LocationDetails}
}

// NewSectorScanReportFromDangerousLocation builds a scan report from a dangerous location snapshot.
func NewSectorScanReportFromDangerousLocation(baseID int, coords Vector2i, loc *DangerousLocationModel, emptyDetails LocationDetails) *SectorScanReport {
	if loc == nil {
		return &SectorScanReport{BaseID: baseID, Coordinates: coords, CreatedAt: NowUnix(), IsCloaked: true, Type: LocationTypeDangerous, Details: emptyDetails}
	}
	armySnaps := loc.MaterializeDefenderArmySnapshot()
	structSnaps := loc.MaterializeDefenderStructureSnapshot()
	defence := sumDefence(armySnaps) + sumStructureDefence(structSnaps)
	attack := sumAttack(armySnaps)
	info := ScanInfo{
		Credits:    loc.Resources.Credits,
		Iron:       loc.Resources.Iron,
		Titanium:   loc.Resources.Titanium,
		Antimatter: loc.Resources.Antimatter,
		Defence:    defence,
		Attack:     attack,
		Space:      0,
	}
	return &SectorScanReport{BaseID: baseID, Coordinates: coords, CreatedAt: NowUnix(), Info: info, IsCloaked: false, Type: LocationTypeDangerous, Details: loc.LocationDetails}
}

// NewSectorScanReportFromEmptyLocation builds a scan report for an empty sector.
func NewSectorScanReportFromEmptySector(baseID int, coords Vector2i, sector *SectorModel) *SectorScanReport {
	// Sector must exist for empty sector scan reports; caller responsible for ensuring non-nil.
	if sector == nil {
		return nil
	}
	info := ScanInfo{}
	return &SectorScanReport{BaseID: baseID, Coordinates: coords, CreatedAt: NowUnix(), Info: info, IsCloaked: false, Type: LocationTypeEmpty, Details: sector.Details}
}
