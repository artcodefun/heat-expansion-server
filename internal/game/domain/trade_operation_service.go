package domain

import (
	"github.com/google/uuid"
)

// TradeOperationService coordinates cross-aggregate validations between
// sender base, receiver base, and a trade operation.
type TradeOperationService struct {
	Sender    *UserBaseModel
	Receiver  *UserBaseModel
	Operation *TradeOperation
}

func NewTradeOperationService(sender *UserBaseModel, receiver *UserBaseModel, operation *TradeOperation) *TradeOperationService {
	return &TradeOperationService{
		Sender:    sender,
		Receiver:  receiver,
		Operation: operation,
	}
}

// BuildTradeOperationForCreation validates trade creation inputs from both
// participants and constructs a TradeOperation aggregate without mutating bases.
// The returned operation can then be used by CommitSenderForTradeCreation.
func BuildTradeOperationForCreation(
	sender *UserBaseModel,
	receiver *UserBaseModel,
	offeredResources PriceModel,
	offeredArmyRequests []ArmyDeploymentRequest,
	offeredStorageItemIDs []uuid.UUID,
	requestedResources PriceModel,
	requestedArmyRequests []ArmyDeploymentRequest,
	requestedStorageItemIDs []uuid.UUID,
	transportRequests []ArmyDeploymentRequest,
) (*TradeOperation, error) {
	if sender == nil || receiver == nil {
		return nil, NewError("error.domain.trade.participant_mismatch", nil)
	}

	if err := validateTradeDeploymentAllocationAvailability(sender, transportRequests, offeredArmyRequests); err != nil {
		return nil, err
	}

	var err error
	var readyTransport []DeploymentReadyItem
	if len(transportRequests) > 0 {
		readyTransport, err = sender.GetReadyToDeployArmy(transportRequests)
		if err != nil {
			return nil, err
		}
	}
	var readyOfferedArmy []DeploymentReadyItem
	if len(offeredArmyRequests) > 0 {
		readyOfferedArmy, err = sender.GetReadyToDeployArmy(offeredArmyRequests)
		if err != nil {
			return nil, err
		}
	}

	var readyRequestedArmy []DeploymentReadyItem
	if len(requestedArmyRequests) > 0 {
		readyRequestedArmy, err = receiver.GetReadyToDeployArmy(requestedArmyRequests)
		if err != nil {
			return nil, err
		}
	}

	if err := sender.Stats.CheckResources(offeredResources); err != nil {
		return nil, err
	}
	if err := receiver.Stats.CheckResources(requestedResources); err != nil {
		return nil, err
	}

	offeredStorage, err := buildTradeStoragePayloadByIDs(sender, offeredStorageItemIDs)
	if err != nil {
		return nil, err
	}
	requestedStorage, err := buildTradeStoragePayloadByIDs(receiver, requestedStorageItemIDs)
	if err != nil {
		return nil, err
	}

	offeredPayload, err := NewTradePayload(offeredResources, offeredStorage, buildTradeArmyPayloadFromReady(readyOfferedArmy))
	if err != nil {
		return nil, err
	}
	requestedPayload, err := NewTradePayload(requestedResources, requestedStorage, buildTradeArmyPayloadFromReady(readyRequestedArmy))
	if err != nil {
		return nil, err
	}

	return NewTradeOperation(
		sender.UserID,
		sender.ID,
		receiver.UserID,
		receiver.ID,
		sender.Coordinates,
		receiver.Coordinates,
		offeredPayload,
		requestedPayload,
		MilitaryUnitsFromDeployed(readyTransport),
		sender.ActiveStorageSnaps(),
	)
}

// CommitSenderForTradeCreation applies sender-side synchronous mutations for
// trade creation after the operation aggregate has already been constructed.
// validates offered payload availability, deploys transport/offered units, deducts
// offered resources and removes offered storage items.
func (s *TradeOperationService) CommitSenderForTradeCreation() error {
	if err := s.validateParticipants(); err != nil {
		return err
	}
	if err := s.validatePayloadAvailability(s.Sender, s.Operation.OfferedPayload); err != nil {
		return err
	}
	if err := validateSenderCommitArmyAvailability(s.Sender, s.Operation.TransportUnits, s.Operation.OfferedPayload.Army); err != nil {
		return err
	}

	if err := allocateArmyByUnitPrototype(s.Sender, s.Operation.TransportUnits, s.Operation.ID); err != nil {
		return err
	}
	if err := allocateArmyByPayloadPrototype(s.Sender, s.Operation.OfferedPayload.Army, s.Operation.ID); err != nil {
		return err
	}

	s.Sender.Stats.SubtractResources(s.Operation.OfferedPayload.Resources)
	if err := s.Sender.AllocateTradePayloadStorageToOperation(s.Operation.OfferedPayload, s.Operation.ID); err != nil {
		return err
	}
	return nil
}

// AcceptAndCommitReceiver applies receiver-side synchronous mutations for accepting
// a trade: validates requested payload availability, deploys requested army,
// deducts requested resources/storage and transitions operation to OUTBOUND.
func (s *TradeOperationService) AcceptAndCommitReceiver() error {
	if err := s.validateParticipants(); err != nil {
		return err
	}
	if err := s.validatePayloadAvailability(s.Receiver, s.Operation.RequestedPayload); err != nil {
		return err
	}

	if err := allocateArmyByPayloadPrototype(s.Receiver, s.Operation.RequestedPayload.Army, s.Operation.ID); err != nil {
		return err
	}

	s.Receiver.Stats.SubtractResources(s.Operation.RequestedPayload.Resources)
	if err := s.Receiver.AllocateTradePayloadStorageToOperation(s.Operation.RequestedPayload, s.Operation.ID); err != nil {
		return err
	}

	if err := s.Operation.Accept(); err != nil {
		return err
	}

	return nil
}

// CancelAndReleaseReceiverIfCommitted cancels the trade operation and releases
// receiver commitments only if they were previously committed (accepted/outbound).
// Sender commitments are released later when return-arrived handling runs.
func (s *TradeOperationService) CancelAndReleaseReceiverIfCommitted() error {
	if err := s.validateParticipants(); err != nil {
		return err
	}

	prevPhase := s.Operation.Phase
	if err := s.Operation.CancelByInitiator(); err != nil {
		return err
	}

	if prevPhase != TradePhaseOutbound {
		return nil
	}

	return rollbackPayloadCommitments(s.Receiver, s.Operation.RequestedPayload, s.Operation.ID)
}

// FinalizeSenderAfterReturn releases sender-side reserved assets at return completion.
// On successful trade, sender receives the requested payload's resources and storage items.
// On non-success terminal outcomes (DECLINED, EXPIRED, CANCELED), sender's own offered
// resources and storage items are restored because the convoy never delivered them.
func (s *TradeOperationService) FinalizeSenderAfterReturn() error {
	if err := s.validateParticipants(); err != nil {
		return err
	}
	if s.Operation.Phase != TradePhaseCompleted {
		return NewError("error.domain.trade.invalid_phase", H{"phase": s.Operation.Phase})
	}

	s.Sender.ReturnAllDeployedFromOperation(OperationKindTrade, s.Operation.ID)
	s.Sender.ReturnAllDeployedStorageFromOperation(OperationKindTrade, s.Operation.ID)
	switch s.Operation.Result {
	case TradeResultSuccess:
		s.Sender.CreditLoot(s.Operation.RequestedPayload.Resources)
	case TradeResultDeclined, TradeResultExpired, TradeResultCanceled:
		s.Sender.CreditLoot(s.Operation.OfferedPayload.Resources)
	}

	return nil
}

// ProcessArrivalAndStartReturn executes the goods swap at ARRIVED phase:
// delivers offered payload to receiver, transfers requested army to sender's
// convoy for the return leg, then starts the return journey.
// Must be called after UpdatePhaseBasedOnTime has transitioned op to ARRIVED.
func (s *TradeOperationService) ProcessArrivalAndStartReturn() error {
	if err := s.validateParticipants(); err != nil {
		return err
	}
	if s.Operation.Phase != TradePhaseArrived {
		return NewError("error.domain.trade.invalid_phase", H{"phase": s.Operation.Phase})
	}

	// Move offered army from sender convoy to receiver convoy.
	offeredStacks, err := s.Sender.RemoveTradeDeployedArmyByPayload(s.Operation.OfferedPayload.Army, s.Operation.ID)
	if err != nil {
		return err
	}
	s.Receiver.AddTradeDeployedArmyStacks(offeredStacks, s.Operation.ID)

	// Move requested army from receiver convoy to sender convoy for the return leg.
	requestedStacks, err := s.Receiver.RemoveTradeDeployedArmyByPayload(s.Operation.RequestedPayload.Army, s.Operation.ID)
	if err != nil {
		return err
	}
	s.Sender.AddTradeDeployedArmyStacks(requestedStacks, s.Operation.ID)

	// Move offered storage from sender convoy to receiver convoy.
	offeredStorage, err := s.Sender.RemoveTradeDeployedStorageByPayload(s.Operation.OfferedPayload.Storage, s.Operation.ID)
	if err != nil {
		return err
	}
	s.Receiver.AddTradeDeployedStorageItems(offeredStorage, s.Operation.ID)

	// Move requested storage from receiver convoy to sender convoy for the return leg.
	requestedStorage, err := s.Receiver.RemoveTradeDeployedStorageByPayload(s.Operation.RequestedPayload.Storage, s.Operation.ID)
	if err != nil {
		return err
	}
	s.Sender.AddTradeDeployedStorageItems(requestedStorage, s.Operation.ID)

	// Receiver's deployed convoy units become present at the destination.
	s.Receiver.CreditLoot(s.Operation.OfferedPayload.Resources)
	s.Receiver.ReturnAllDeployedFromOperation(OperationKindTrade, s.Operation.ID)
	s.Receiver.ReturnAllDeployedStorageFromOperation(OperationKindTrade, s.Operation.ID)

	return s.Operation.StartReturn()
}

func (s *TradeOperationService) validateParticipants() error {
	if s.Sender == nil || s.Receiver == nil || s.Operation == nil {
		return NewError("error.domain.trade.participant_mismatch", nil)
	}
	if s.Sender.UserID != s.Operation.SenderUserID || s.Sender.ID != s.Operation.SenderBaseID {
		return NewError("error.domain.trade.participant_mismatch", nil)
	}
	if s.Receiver.UserID != s.Operation.ReceiverUserID || s.Receiver.ID != s.Operation.ReceiverBaseID {
		return NewError("error.domain.trade.participant_mismatch", nil)
	}
	return nil
}

func (s *TradeOperationService) validatePayloadAvailability(base *UserBaseModel, payload TradePayload) error {
	if err := base.Stats.CheckResources(payload.Resources); err != nil {
		return err
	}

	storageByID := make(map[uuid.UUID]StorageItemPresent, len(base.StorageItemsPresent))
	for _, it := range base.StorageItemsPresent {
		storageByID[it.ID] = it
	}
	for _, st := range payload.Storage {
		item, ok := storageByID[st.ItemID]
		if !ok || item.IsActive {
			return NewError("error.domain.trade.storage_item_not_tradeable", H{"item_id": st.ItemID})
		}
	}

	presentByProto := make(map[int]int)
	for _, ap := range base.ArmiesPresent {
		presentByProto[ap.Prototype.ID] += ap.Count
	}
	for _, ar := range payload.Army {
		available := presentByProto[ar.PrototypeID]
		if available < ar.Count {
			return NewError("error.domain.trade.insufficient_army_units", H{"prototype_id": ar.PrototypeID, "required": ar.Count, "available": available})
		}
	}

	return nil
}

func rollbackPayloadCommitments(base *UserBaseModel, payload TradePayload, operationID int) error {
	base.CreditLoot(payload.Resources)
	base.ReturnAllDeployedFromOperation(OperationKindTrade, operationID)
	base.ReturnAllDeployedStorageFromOperation(OperationKindTrade, operationID)
	return nil
}

func allocateArmyByPayloadPrototype(base *UserBaseModel, army []TradeArmyItemSnap, operationID int) error {
	for _, req := range army {
		remaining := req.Count
		for remaining > 0 {
			idx := -1
			for i, p := range base.ArmiesPresent {
				if p.Prototype.ID == req.PrototypeID && p.Count > 0 {
					idx = i
					break
				}
			}
			if idx == -1 {
				return NewError("error.domain.trade.insufficient_army_units", H{"prototype_id": req.PrototypeID, "required": req.Count, "available": req.Count - remaining})
			}

			take := remaining
			if base.ArmiesPresent[idx].Count < take {
				take = base.ArmiesPresent[idx].Count
			}

			if _, err := base.AllocateArmyToOperation(ArmyDeploymentRequest{PresentItemID: base.ArmiesPresent[idx].ID, Count: take}, OperationKindTrade, operationID); err != nil {
				return err
			}
			remaining -= take
		}
	}
	return nil
}

func allocateArmyByUnitPrototype(base *UserBaseModel, units []MilitaryUnitSnap, operationID int) error {
	for _, req := range units {
		remaining := req.Count
		for remaining > 0 {
			idx := -1
			for i, p := range base.ArmiesPresent {
				if p.Prototype.ID == req.PrototypeID && p.Count > 0 {
					idx = i
					break
				}
			}
			if idx == -1 {
				return NewError("error.domain.trade.insufficient_army_units", H{"prototype_id": req.PrototypeID, "required": req.Count, "available": req.Count - remaining})
			}

			take := remaining
			if base.ArmiesPresent[idx].Count < take {
				take = base.ArmiesPresent[idx].Count
			}

			if _, err := base.AllocateArmyToOperation(ArmyDeploymentRequest{PresentItemID: base.ArmiesPresent[idx].ID, Count: take}, OperationKindTrade, operationID); err != nil {
				return err
			}
			remaining -= take
		}
	}
	return nil
}

func validateSenderCommitArmyAvailability(base *UserBaseModel, transport []MilitaryUnitSnap, offered []TradeArmyItemSnap) error {
	requiredByProto := make(map[int]int)
	for _, u := range transport {
		requiredByProto[u.PrototypeID] += u.Count
	}
	for _, a := range offered {
		requiredByProto[a.PrototypeID] += a.Count
	}

	availableByProto := make(map[int]int)
	for _, p := range base.ArmiesPresent {
		availableByProto[p.Prototype.ID] += p.Count
	}

	for protoID, required := range requiredByProto {
		available := availableByProto[protoID]
		if required > available {
			return NewError("error.domain.trade.insufficient_army_units", H{"prototype_id": protoID, "required": required, "available": available})
		}
	}

	return nil
}

func buildTradeStoragePayloadByIDs(base *UserBaseModel, itemIDs []uuid.UUID) ([]TradeStorageItemSnap, error) {
	if len(itemIDs) == 0 {
		return nil, nil
	}

	byID := make(map[uuid.UUID]StorageItemPresent, len(base.StorageItemsPresent))
	for _, it := range base.StorageItemsPresent {
		byID[it.ID] = it
	}

	out := make([]TradeStorageItemSnap, 0, len(itemIDs))
	for _, itemID := range itemIDs {
		item, ok := byID[itemID]
		if !ok || item.IsActive {
			return nil, NewError("error.domain.trade.storage_item_not_tradeable", H{"item_id": itemID})
		}
		out = append(out, TradeStorageItemSnap{
			ItemID:      item.ID,
			PrototypeID: item.Prototype.ID,
			Category:    item.Prototype.Category,
		})
	}

	return out, nil
}

func buildTradeArmyPayloadFromReady(ready []DeploymentReadyItem) []TradeArmyItemSnap {
	if len(ready) == 0 {
		return nil
	}

	byProto := make(map[int]TradeArmyItemSnap)
	for _, r := range ready {
		existing, ok := byProto[r.Prototype.ID]
		if !ok {
			byProto[r.Prototype.ID] = TradeArmyItemSnap{
				PrototypeID: r.Prototype.ID,
				Count:       r.Count,
				Capacity:    r.Prototype.Capacity,
			}
			continue
		}
		existing.Count += r.Count
		byProto[r.Prototype.ID] = existing
	}

	out := make([]TradeArmyItemSnap, 0, len(byProto))
	for _, snap := range byProto {
		out = append(out, snap)
	}
	return out
}

// validateTradeDeploymentAllocationAvailability ensures the total amount requested
// from each present army stack across transport and offered groups does not exceed
// what is currently available in sender present armies.
func validateTradeDeploymentAllocationAvailability(base *UserBaseModel, transportRequests []ArmyDeploymentRequest, offeredRequests []ArmyDeploymentRequest) error {
	combined := make([]ArmyDeploymentRequest, 0, len(transportRequests)+len(offeredRequests))
	combined = append(combined, transportRequests...)
	combined = append(combined, offeredRequests...)
	_, err := base.GetReadyToDeployArmy(combined)
	return err
}
