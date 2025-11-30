package readmodels

// ActivityKind enumerates the kinds of activities the system can present.
type ActivityKind string

const (
	ActivityKindMilitary ActivityKind = "MILITARY"
	ActivityKindScan     ActivityKind = "SCAN"
	ActivityKindRadar    ActivityKind = "RADAR"
	ActivityKindTrade    ActivityKind = "TRADE"
)

// ActivityItem is a domain projection used by the Activities use case.
type ActivityItem struct {
	ID        int
	Kind      ActivityKind
	CreatedAt int64
	BaseID    int

	Operation *OperationActivity
	Scan      *ScanActivity
	Radar     *RadarActivity
}

// MilitaryActivitySubtype specifies the subtype within the MILITARY kind.
type MilitaryActivitySubtype string

const (
	MilitaryActivitySubtypeAttack  MilitaryActivitySubtype = "ATTACK"
	MilitaryActivitySubtypeSpy     MilitaryActivitySubtype = "SPY"
	MilitaryActivitySubtypeDefense MilitaryActivitySubtype = "DEFENSE"
)

// OperationRole indicates the viewer's role relative to an operation.
type OperationRole string

const (
	OperationRoleAttacker OperationRole = "ATTACKER"
	OperationRoleDefender OperationRole = "DEFENDER"
)

// OperationActivity summarizes an operation for activities list.
type OperationActivity struct {
	OpID    int
	Subtype MilitaryActivitySubtype
	Role    OperationRole
}

// ScanActivity wraps a SectorScanReport into the activity stream.
type ScanActivity struct {
	ReportID int
}

// RadarActivity represents a detected incoming hostility (future wiring).
type RadarActivity struct {
	OpID              int
	DetectedAt        int64
	EtaAtBase         int64
	SourceCoordinates Vector2i
	TargetCoordinates Vector2i
	Threat            Threat
}

type Threat struct {
	Attack  int
	Defence int
}
