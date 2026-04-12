package domain

import (
	"time"

	"github.com/google/uuid"
)

type DiplomaticRequestKind string

const (
	DiplomaticRequestKindCoalitionProposal DiplomaticRequestKind = "COALITION_PROPOSAL"
	DiplomaticRequestKindCeasefireProposal DiplomaticRequestKind = "CEASEFIRE_PROPOSAL"
)

type DiplomaticRequestStatus string

const (
	DiplomaticRequestStatusPending  DiplomaticRequestStatus = "PENDING"
	DiplomaticRequestStatusAccepted DiplomaticRequestStatus = "ACCEPTED"
	DiplomaticRequestStatusRejected DiplomaticRequestStatus = "REJECTED"
	DiplomaticRequestStatusExpired  DiplomaticRequestStatus = "EXPIRED"
)

const DiplomaticRequestExpiry = 24 * time.Hour

type DiplomaticRequest struct {
	EventProducer

	ID             uuid.UUID
	SenderUserID   uuid.UUID
	ReceiverUserID uuid.UUID
	SenderBaseID   *int
	ReceiverBaseID *int
	Kind           DiplomaticRequestKind
	Status         DiplomaticRequestStatus
	CreatedAt      int64
	ResolvedAt     *int64
	ExpiresAt      int64
}

func IsDiplomaticRequestKind(kind DiplomaticRequestKind) bool {
	switch kind {
	case DiplomaticRequestKindCoalitionProposal, DiplomaticRequestKindCeasefireProposal:
		return true
	default:
		return false
	}
}

func NewDiplomaticRequest(senderUserID, receiverUserID uuid.UUID, senderBaseID, receiverBaseID *int, kind DiplomaticRequestKind) (*DiplomaticRequest, error) {
	if err := ValidateDiplomaticParticipants(senderUserID, receiverUserID); err != nil {
		return nil, err
	}
	if !IsDiplomaticRequestKind(kind) {
		return nil, NewError("error.domain.diplomacy.invalid_request_kind", nil)
	}
	now := NowUnix()
	request := &DiplomaticRequest{
		ID:             uuid.Must(uuid.NewV7()),
		SenderUserID:   senderUserID,
		ReceiverUserID: receiverUserID,
		SenderBaseID:   senderBaseID,
		ReceiverBaseID: receiverBaseID,
		Kind:           kind,
		Status:         DiplomaticRequestStatusPending,
		CreatedAt:      now,
		ExpiresAt:      now + int64(DiplomaticRequestExpiry.Seconds()),
	}
	request.AddEvent(NewDiplomaticRequestCreatedEvent(request.ID, request.SenderUserID, request.ReceiverUserID, request.ReceiverBaseID, request.Kind))
	return request, nil
}

func (r *DiplomaticRequest) ValidateAgainstRelationship(rel *DiplomaticRelationship, pendingExists bool) error {
	if pendingExists {
		return NewError("error.domain.diplomacy.proposal_already_pending", nil)
	}
	switch r.Kind {
	case DiplomaticRequestKindCoalitionProposal:
		return rel.CanCreateAllianceRequest()
	case DiplomaticRequestKindCeasefireProposal:
		return rel.CanCreateCeasefireRequest()
	default:
		return NewError("error.domain.diplomacy.invalid_request_kind", nil)
	}
}

func (r *DiplomaticRequest) Accept(actorUserID uuid.UUID, rel *DiplomaticRelationship) (TranslationKey, error) {
	if err := r.ensurePendingForReceiver(actorUserID); err != nil {
		return "", err
	}
	if rel == nil {
		return "", NewError("error.domain.diplomacy.invalid_transition", nil)
	}
	switch r.Kind {
	case DiplomaticRequestKindCoalitionProposal:
		rel.FormAlliance(actorUserID)
	case DiplomaticRequestKindCeasefireProposal:
		rel.AcceptCeasefire(actorUserID)
	default:
		return "", NewError("error.domain.diplomacy.invalid_request_kind", nil)
	}
	now := NowUnix()
	r.Status = DiplomaticRequestStatusAccepted
	r.ResolvedAt = &now
	return r.acceptanceMessageContent(), nil
}

func (r *DiplomaticRequest) Reject(actorUserID uuid.UUID) (TranslationKey, error) {
	if err := r.ensurePendingForReceiver(actorUserID); err != nil {
		return "", err
	}
	now := NowUnix()
	r.Status = DiplomaticRequestStatusRejected
	r.ResolvedAt = &now
	return r.rejectionMessageContent(), nil
}

func (r *DiplomaticRequest) CanExpire() bool {
	if r == nil {
		return false
	}
	if r.Status != DiplomaticRequestStatusPending {
		return false
	}
	return NowUnix() >= r.ExpiresAt
}

func (r *DiplomaticRequest) Expire() {
	now := NowUnix()
	r.Status = DiplomaticRequestStatusExpired
	r.ResolvedAt = &now
}

func (r *DiplomaticRequest) ensurePendingForReceiver(actorUserID uuid.UUID) error {
	if r == nil || actorUserID != r.ReceiverUserID {
		return NewError("error.domain.diplomacy.invalid_transition", nil)
	}
	if r.Status != DiplomaticRequestStatusPending {
		return NewError("error.domain.diplomacy.invalid_transition", nil)
	}
	if NowUnix() >= r.ExpiresAt {
		return NewError("error.domain.diplomacy.invalid_transition", nil)
	}
	return nil
}

func (r *DiplomaticRequest) MessageContent() TranslationKey {
	if r == nil {
		return ""
	}
	switch r.Kind {
	case DiplomaticRequestKindCoalitionProposal:
		return DiplomaticMessageContentCoalitionProposal
	case DiplomaticRequestKindCeasefireProposal:
		return DiplomaticMessageContentCeasefireProposal
	default:
		return ""
	}
}

func (r *DiplomaticRequest) acceptanceMessageContent() TranslationKey {
	switch r.Kind {
	case DiplomaticRequestKindCoalitionProposal:
		return DiplomaticMessageContentCoalitionAcceptance
	case DiplomaticRequestKindCeasefireProposal:
		return DiplomaticMessageContentCeasefireAcceptance
	default:
		return ""
	}
}

func (r *DiplomaticRequest) rejectionMessageContent() TranslationKey {
	switch r.Kind {
	case DiplomaticRequestKindCoalitionProposal:
		return DiplomaticMessageContentCoalitionRejection
	case DiplomaticRequestKindCeasefireProposal:
		return DiplomaticMessageContentCeasefireRejection
	default:
		return ""
	}
}
