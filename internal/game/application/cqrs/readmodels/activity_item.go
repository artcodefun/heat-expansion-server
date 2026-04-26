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
	ID        uuid.UUID
	Kind      ActivityKind
	CreatedAt int64
	BaseID    int

	Offense *OffenseActivity
	Defense *DefenseActivity
	Scan    *ScanActivity
	Radar   *RadarActivity
	Trade   *TradeActivity
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

	Operation *MilitaryOperation
}

// OffenderInfo provides a restricted view of an attacking operation for the defender.
type OffenderInfo struct {
	Type              MilitaryOperationType
	SourceCoordinates *Vector2i
	TargetCoordinates Vector2i
	ContactDate       int64
	Result            MilitaryOperationResult
	Units             []MilitaryUnitSnap
	StorageSnaps      []StorageItemSnap
	TotalModifiers    MilitaryModifiers
	SpyResult         *SpyResult
	AttackResult      *AttackResult
}

// NewOffenderInfoFromOperation creates a restricted view of an operation for the defender,
// applying business rules for information masking (e.g. hiding source coordinates for successful stealthy spies).
func NewOffenderInfoFromOperation(op *MilitaryOperation) *OffenderInfo {
	if op == nil {
		return nil
	}

	var sourceCoords *Vector2i
	// Hide source coordinates for successful SPY operations where the defender didn't detect them via spies
	if op.Type == MilitaryOperationTypeSpy && op.SpyResult != nil && op.SpyResult.Outcome == SpyOutcomeReportProduced {
		sourceCoords = nil
	} else {
		// Create a copy of coordinates to point to
		coords := op.SourceCoordinates
		sourceCoords = &coords
	}

	return &OffenderInfo{
		Type:              op.Type,
		SourceCoordinates: sourceCoords,
		TargetCoordinates: op.TargetCoordinates,
		ContactDate:       op.OutboundArriveAt,
		Result:            op.Result,
		Units:             op.Units,
		StorageSnaps:      op.StorageSnaps,
		TotalModifiers:    op.TotalModifiers,
		SpyResult:         op.SpyResult,
		AttackResult:      op.AttackResult,
	}
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

	Offender          *OffenderInfo
	PriorOpponentScan *SectorScanReport
}

// ScanActivitySubtype specifies the subtype of a scan-related activity.
type ScanActivitySubtype string

const (
	ScanActivitySubtypeReportProduced       ScanActivitySubtype = "REPORT_PRODUCED"
	ScanActivitySubtypeExternalScanDetected ScanActivitySubtype = "EXTERNAL_SCAN_DETECTED"
)

// ScanInterceptInfo represents detected hostile scanning with optional triangulation.
type ScanInterceptInfo struct {
	ScannedCoordinates     Vector2i
	ScanPenetratedCloaking bool
	PossibleSource         *Vector2i
	UncertaintyRadius      int
}

// ScanActivity wraps scan-related events into the activity stream.
type ScanActivity struct {
	Subtype   ScanActivitySubtype
	ReportID  *int
	Intercept *ScanInterceptInfo

	Report *SectorScanReport
}

// RadarActivity represents a link to a stateful radar threat.
type RadarActivity struct {
	ThreatID uuid.UUID
	Threat   *RadarThreat
}

// TradeActivity summarizes a trade lifecycle event.
type TradeActivity struct {
	OpID int

	Operation *TradeOperation
}
