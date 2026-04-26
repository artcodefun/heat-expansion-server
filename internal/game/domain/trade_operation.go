package domain

import "github.com/google/uuid"

type TradeOperationPhase string

const (
	TradePhasePending   TradeOperationPhase = "PENDING"
	TradePhaseOutbound  TradeOperationPhase = "OUTBOUND"
	TradePhaseArrived   TradeOperationPhase = "ARRIVED"
	TradePhaseReturning TradeOperationPhase = "RETURNING"
	TradePhaseCompleted TradeOperationPhase = "COMPLETED"
)

type TradeOperationResult string

const (
	TradeResultUnknown  TradeOperationResult = "UNKNOWN"
	TradeResultSuccess  TradeOperationResult = "SUCCESS"
	TradeResultCanceled TradeOperationResult = "CANCELED"
	TradeResultDeclined TradeOperationResult = "DECLINED"
	TradeResultExpired  TradeOperationResult = "EXPIRED"
	TradeResultFailure  TradeOperationResult = "FAILURE"
)

type TradeOperation struct {
	EventProducer

	ID   int
	UUID uuid.UUID

	CreatedAt int64

	SenderUserID   uuid.UUID
	SenderBaseID   int
	ReceiverUserID uuid.UUID
	ReceiverBaseID int

	SourceCoordinates Vector2i
	TargetCoordinates Vector2i

	OfferedPayload   TradePayload
	RequestedPayload TradePayload

	TransportUnits []MilitaryUnitSnap
	StorageSnaps   []StorageItemSnap
	TotalModifiers MilitaryModifiers

	ExpiresAt int64

	OutboundDepartAt  int64
	OutboundArriveAt  int64
	ArrivedAtTargetAt int64
	ReturnDepartAt    int64
	ReturnArriveAt    int64
	CompletedAt       int64

	CrystalsSkipPrice int

	Phase  TradeOperationPhase
	Result TradeOperationResult
}

func NewTradeOperation(
	senderUserID uuid.UUID,
	senderBaseID int,
	receiverUserID uuid.UUID,
	receiverBaseID int,
	source Vector2i,
	target Vector2i,
	offered TradePayload,
	requested TradePayload,
	transportUnits []MilitaryUnitSnap,
	storageSnaps []StorageItemSnap,
) (*TradeOperation, error) {
	if senderUserID == receiverUserID || senderBaseID == receiverBaseID {
		return nil, NewError("error.domain.trade.same_participant", nil)
	}
	if source == target {
		return nil, NewError("error.domain.trade.invalid_coordinates", nil)
	}
	if err := validateTradeOperationLegUnits(offered, transportUnits, "error.domain.trade.outbound_units_required"); err != nil {
		return nil, err
	}
	if err := validateTradeOperationLegUnits(requested, transportUnits, "error.domain.trade.return_units_required"); err != nil {
		return nil, err
	}
	mods := MilitaryModifiersFromSnaps(storageSnaps)
	if err := validateTradeOperationLegCapacity(offered, transportUnits, mods.CapacityMul, "error.domain.trade.outbound_capacity_insufficient"); err != nil {
		return nil, err
	}
	if err := validateTradeOperationLegCapacity(requested, transportUnits, mods.CapacityMul, "error.domain.trade.return_capacity_insufficient"); err != nil {
		return nil, err
	}

	now := NowUnix()
	expiresAt := now + 24*60*60

	op := &TradeOperation{
		UUID:              uuid.Must(uuid.NewV7()),
		CreatedAt:         now,
		SenderUserID:      senderUserID,
		SenderBaseID:      senderBaseID,
		ReceiverUserID:    receiverUserID,
		ReceiverBaseID:    receiverBaseID,
		SourceCoordinates: source,
		TargetCoordinates: target,
		OfferedPayload:    offered,
		RequestedPayload:  requested,
		TransportUnits:    cloneUnits(transportUnits),
		StorageSnaps:      cloneStorageSnaps(storageSnaps),
		TotalModifiers:    mods,
		ExpiresAt:         expiresAt,
		Phase:             TradePhasePending,
		Result:            TradeResultUnknown,
	}
	return op, nil
}

// EmitCreatedEvent emits operation-created event after persistent ID assignment.
func (op *TradeOperation) EmitCreatedEvent() {
	op.AddEvent(NewTradeOperationCreatedEvent(op.ID, op.UUID, op.SenderBaseID, op.ReceiverBaseID, op.ExpiresAt))
}

func validateTradeOperationLegCapacity(payload TradePayload, transportUnits []MilitaryUnitSnap, capacityMul float64, insufficientKey TranslationKey) error {
	required := payload.RequiredResourceCapacity()
	providedRaw := payload.ProvidedArmyCapacity() + SumEffectiveCapacity(transportUnits, 1)
	provided := providedRaw * capacityMul
	if provided < required {
		return NewError(insufficientKey, H{"required": required, "capacity": provided})
	}
	return nil
}

func validateTradeOperationLegUnits(payload TradePayload, transportUnits []MilitaryUnitSnap, insufficientKey TranslationKey) error {
	if len(transportUnits) == 0 && len(payload.Army) == 0 {
		return NewError(insufficientKey, nil)
	}
	return nil
}

func (op *TradeOperation) Accept() error {
	if op.Phase != TradePhasePending {
		return NewError("error.domain.trade.invalid_phase", H{"phase": op.Phase})
	}
	now := NowUnix()
	travelSeconds := computeTravelSecondsBetween(op.SourceCoordinates, op.TargetCoordinates, op.TransportUnits, op.TotalModifiers)
	op.OutboundDepartAt = now
	op.OutboundArriveAt = now + travelSeconds
	op.CrystalsSkipPrice = max(1, int(travelSeconds/60))
	op.Phase = TradePhaseOutbound
	op.AddEvent(NewTradeOperationAcceptedEvent(op.ID, op.OutboundArriveAt))
	op.AddEvent(NewTradeOperationOutboundEvent(op.ID, op.OutboundArriveAt))
	return nil
}

func (op *TradeOperation) Decline() error {
	if op.Phase != TradePhasePending {
		return NewError("error.domain.trade.invalid_phase", H{"phase": op.Phase})
	}
	now := NowUnix()
	op.Phase = TradePhaseCompleted
	op.Result = TradeResultDeclined
	op.CompletedAt = now
	op.AddEvent(NewTradeOperationDeclinedEvent(op.ID))
	op.AddEvent(NewTradeOperationReturnArrivedEvent(op.ID))
	return nil
}

// ExpireIfPending is idempotent for scheduler usage: only PENDING transitions to terminal completion.
func (op *TradeOperation) ExpireIfPending() bool {
	if op.Phase != TradePhasePending {
		return false
	}
	now := NowUnix()
	op.Phase = TradePhaseCompleted
	op.Result = TradeResultExpired
	op.CompletedAt = now
	op.AddEvent(NewTradeOperationExpiredEvent(op.ID))
	op.AddEvent(NewTradeOperationReturnArrivedEvent(op.ID))
	return true
}

// CancelByInitiator is allowed only in PENDING and OUTBOUND.
func (op *TradeOperation) CancelByInitiator() error {
	if op.Phase != TradePhasePending && op.Phase != TradePhaseOutbound {
		return NewError("error.domain.trade.invalid_phase", H{"phase": op.Phase})
	}
	if op.Phase == TradePhasePending {
		op.Phase = TradePhaseCompleted
		op.Result = TradeResultCanceled
		op.CompletedAt = NowUnix()
		op.AddEvent(NewTradeOperationCancelledByInitiatorEvent(op.ID))
		op.AddEvent(NewTradeOperationReturnArrivedEvent(op.ID))
		return nil
	}

	cancelAt := NowUnix()
	op.Result = TradeResultCanceled
	current := lerpCoordinates(op.SourceCoordinates, op.TargetCoordinates, op.OutboundDepartAt, op.OutboundArriveAt, cancelAt)
	op.startReturnLegFrom(current, cancelAt)
	op.AddEvent(NewTradeOperationCancelledByInitiatorEvent(op.ID))
	return nil
}

func (op *TradeOperation) OnArrive() error {
	if op.Phase != TradePhaseOutbound {
		return NewError("error.domain.trade.invalid_phase", H{"phase": op.Phase})
	}
	op.ArrivedAtTargetAt = NowUnix()
	op.Phase = TradePhaseArrived
	op.AddEvent(NewTradeOperationArrivedEvent(op.ID))
	return nil
}

func (op *TradeOperation) StartReturn() error {
	if op.Phase != TradePhaseArrived {
		return NewError("error.domain.trade.invalid_phase", H{"phase": op.Phase})
	}
	op.startReturnLegFrom(op.TargetCoordinates, NowUnix())
	return nil
}

func (op *TradeOperation) OnReturnArrive() error {
	if op.Phase != TradePhaseReturning {
		return NewError("error.domain.trade.invalid_phase", H{"phase": op.Phase})
	}
	now := NowUnix()
	op.CompletedAt = now
	op.Phase = TradePhaseCompleted
	op.AddEvent(NewTradeOperationReturnArrivedEvent(op.ID))
	if op.Result == TradeResultCanceled {
		return nil
	}
	op.Result = TradeResultSuccess
	return nil
}

// UpdatePhaseBasedOnTime advances the trade operation phase based on current time.
// It is idempotent and safe to call multiple times.
func (op *TradeOperation) UpdatePhaseBasedOnTime() {
	now := NowUnix()
	switch op.Phase {
	case TradePhaseOutbound:
		if op.OutboundArriveAt > 0 && now >= op.OutboundArriveAt {
			_ = op.OnArrive()
		}
	case TradePhaseReturning:
		if op.ReturnArriveAt > 0 && now >= op.ReturnArriveAt {
			_ = op.OnReturnArrive()
		}
	default:
		// no-op for other phases
	}
}

func (op *TradeOperation) startReturnLegFrom(from Vector2i, now int64) {
	travelSeconds := computeTravelSecondsBetween(from, op.SourceCoordinates, op.TransportUnits, op.TotalModifiers)
	op.ReturnDepartAt = now
	op.ReturnArriveAt = now + travelSeconds
	op.CrystalsSkipPrice = max(1, int(travelSeconds/60))
	op.Phase = TradePhaseReturning
	op.AddEvent(NewTradeOperationReturningEvent(op.ID, op.ReturnArriveAt))
}
