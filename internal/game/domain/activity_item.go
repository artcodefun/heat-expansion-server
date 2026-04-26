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
	EventProducer

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

// ScanActivitySubtype specifies the subtype of a scan-related activity.
type ScanActivitySubtype string

const (
	ScanActivitySubtypeReportProduced       ScanActivitySubtype = "REPORT_PRODUCED"
	ScanActivitySubtypeExternalScanDetected ScanActivitySubtype = "EXTERNAL_SCAN_DETECTED"
)

// ScanInterceptInfo represents detected hostile scanning with optional triangulation.
type ScanInterceptInfo struct {
	// ScannedCoordinates is the sector that was targeted by the attacker scan.
	ScannedCoordinates Vector2i

	// ScanPenetratedCloaking indicates whether the attacker scan successfully
	// obtained accurate intel despite defender cloaking.
	ScanPenetratedCloaking bool

	// PossibleSource is a vague estimate of scan origin (optional).
	PossibleSource *Vector2i

	// UncertaintyRadius defines how imprecise PossibleSource is (0 => unknown / not triangulated).
	UncertaintyRadius int
}

// ScanActivity wraps scan-related events into the activity stream.
type ScanActivity struct {
	Subtype   ScanActivitySubtype
	ReportID  *int
	Intercept *ScanInterceptInfo
}

// RadarActivity represents a detected incoming hostility (future wiring).
type RadarActivity struct {
	ThreatID uuid.UUID
}

// TradeActivity summarizes a trade creation activity.
type TradeActivity struct {
	OpID int
}

// Helpers to build ActivityItem from domain entities

func newEmptyActivity(userID uuid.UUID, kind ActivityKind, baseID int, subtype string) ActivityItem {
	id := uuid.Must(uuid.NewV7())
	item := ActivityItem{
		ID:        id,
		Kind:      kind,
		CreatedAt: NowUnix(),
		BaseID:    baseID,
	}
	item.AddEvent(NewActivityCreatedEvent(id, userID, baseID, kind, subtype))
	return item
}

func NewActivityFromOffenseOperation(baseID int, op *MilitaryOperation) ActivityItem {
	subtype := func() OffenseActivitySubtype {
		if op.Type == MilitaryOperationTypeSpy {
			return OffenseActivitySubtypeSpy
		}
		return OffenseActivitySubtypeAttack
	}()

	item := newEmptyActivity(op.OwnerUserID, ActivityKindOffense, baseID, string(subtype))
	item.Offense = &OffenseActivity{
		OpID:    op.ID,
		Subtype: subtype,
	}
	return item
}

func NewActivityFromDefenseOperation(userID uuid.UUID, baseID int, op *MilitaryOperation) ActivityItem {
	subtype := func() DefenseActivitySubtype {
		if op.Type == MilitaryOperationTypeSpy {
			return DefenseActivitySubtypeSpy
		}
		return DefenseActivitySubtypeAttack
	}()

	item := newEmptyActivity(userID, ActivityKindDefense, baseID, string(subtype))
	item.Defense = &DefenseActivity{
		OpID:    op.ID,
		Subtype: subtype,
	}
	return item
}

func NewActivityFromScan(userID uuid.UUID, baseID int, r *SectorScanReport) ActivityItem {
	rid := r.ID
	subtype := ScanActivitySubtypeReportProduced
	item := newEmptyActivity(userID, ActivityKindScan, baseID, string(subtype))

	item.Scan = &ScanActivity{
		Subtype:  subtype,
		ReportID: &rid,
	}
	return item
}

func NewActivityFromScanIntercept(
	userID uuid.UUID,
	baseID int,
	info ScanInterceptInfo,
) ActivityItem {
	subtype := ScanActivitySubtypeExternalScanDetected
	item := newEmptyActivity(userID, ActivityKindScan, baseID, string(subtype))

	item.Scan = &ScanActivity{
		Subtype:   subtype,
		Intercept: &info,
	}
	return item
}

func NewActivityFromRadarThreat(userID uuid.UUID, t *RadarThreat) ActivityItem {
	item := newEmptyActivity(userID, ActivityKindRadar, t.OwnerBaseID, "")
	item.Radar = &RadarActivity{
		ThreatID: t.ID,
	}
	return item
}

func NewActivityFromTradeOperation(userID uuid.UUID, baseID int, opID int) ActivityItem {
	item := newEmptyActivity(userID, ActivityKindTrade, baseID, "")
	item.Trade = &TradeActivity{OpID: opID}
	return item
}
