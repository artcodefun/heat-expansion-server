package readmodels

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

	Operation         *MilitaryOperation
	PriorOpponentScan *SectorScanReport
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

	Operation         *MilitaryOperation
	PriorOpponentScan *SectorScanReport
}

// ScanActivity wraps a SectorScanReport into the activity stream.
type ScanActivity struct {
	ReportID int

	Report *SectorScanReport
}

// RadarActivity represents a link to a stateful radar threat.
type RadarActivity struct {
	ThreatID uuid.UUID
	Threat   *RadarThreat
}
