package domain

import "github.com/google/uuid"

// ActivityKind enumerates the kinds of activities the system can present.
type ActivityKind string

const (
	ActivityKindOffense ActivityKind = "OFFENSE"
	ActivityKindDefense ActivityKind = "DEFENSE"
	ActivityKindScan    ActivityKind = "SCAN"
	ActivityKindRadar   ActivityKind = "RADAR"
	ActivityKindTrade   ActivityKind = "TRADE"
)

// ActivityItem is a domain projection used by the Activities use case.
type ActivityItem struct {
	ID        int
	Kind      ActivityKind
	CreatedAt int64
	BaseID    int

	Offense *OffenseActivity
	Defense *DefenseActivity
	Scan    *ScanActivity
	Radar   *RadarActivity
}

// OffenseActivitySubtype specifies the subtype of an offensive activity.
type OffenseActivitySubtype string

const (
	OffenseActivitySubtypeAttack OffenseActivitySubtype = "ATTACK"
	OffenseActivitySubtypeSpy    OffenseActivitySubtype = "SPY"
)

// OffenseActivity summarizes an offensive mission.
type OffenseActivity struct {
	OpID    int
	Subtype OffenseActivitySubtype
}

// DefenseActivitySubtype specifies the subtype of a defensive activity.
type DefenseActivitySubtype string

const (
	DefenseActivitySubtypeAttack DefenseActivitySubtype = "ATTACK"
	DefenseActivitySubtypeSpy    DefenseActivitySubtype = "SPY"
)

// DefenseActivity summarizes a defensive engagement.
type DefenseActivity struct {
	OpID    int
	Subtype DefenseActivitySubtype
}

// ScanActivity wraps a SectorScanReport into the activity stream.
type ScanActivity struct {
	ReportID int
}

// RadarActivity represents a detected incoming hostility (future wiring).
type RadarActivity struct {
	ThreatID uuid.UUID
}

// Helpers to build ActivityItem from domain entities

func NewActivityFromOffenseOperation(baseID int, op *MilitaryOperation) ActivityItem {
	return ActivityItem{
		ID:        0,
		Kind:      ActivityKindOffense,
		CreatedAt: NowUnix(),
		BaseID:    baseID,
		Offense: &OffenseActivity{
			OpID: op.ID,
			Subtype: func() OffenseActivitySubtype {
				if op.Type == MilitaryOperationTypeSpy {
					return OffenseActivitySubtypeSpy
				}
				return OffenseActivitySubtypeAttack
			}(),
		},
	}
}

func NewActivityFromDefenseOperation(baseID int, op *MilitaryOperation) ActivityItem {
	return ActivityItem{
		ID:        0,
		Kind:      ActivityKindDefense,
		CreatedAt: NowUnix(),
		BaseID:    baseID,
		Defense: &DefenseActivity{
			OpID: op.ID,
			Subtype: func() DefenseActivitySubtype {
				if op.Type == MilitaryOperationTypeSpy {
					return DefenseActivitySubtypeSpy
				}
				return DefenseActivitySubtypeAttack
			}(),
		},
	}
}

func NewActivityFromScan(baseID int, r *SectorScanReport) ActivityItem {
	return ActivityItem{
		ID:        0, // assigned by persistence layer
		Kind:      ActivityKindScan,
		CreatedAt: NowUnix(),
		BaseID:    baseID,
		// Category removed
		Scan: &ScanActivity{ReportID: r.ID},
	}
}

func NewActivityFromRadarThreat(t *RadarThreat) ActivityItem {
	return ActivityItem{
		ID:        0,
		Kind:      ActivityKindRadar,
		CreatedAt: NowUnix(),
		BaseID:    t.OwnerBaseID,
		Radar: &RadarActivity{
			ThreatID: t.ID,
		},
	}
}
