package domain

import (
	"time"

	"github.com/google/uuid"
)

const (
	DiplomaticWarWarmupDuration           = 24 * time.Hour
	DiplomaticCeasefireProtectionDuration = 24 * time.Hour
)

type DiplomaticStatus string

const (
	DiplomaticStatusUnknown DiplomaticStatus = "UNKNOWN"
	DiplomaticStatusNeutral DiplomaticStatus = "NEUTRAL"
	DiplomaticStatusAllied  DiplomaticStatus = "ALLIED"
	DiplomaticStatusWar     DiplomaticStatus = "WAR"
)

type DiplomaticRelationship struct {
	EventProducer
	ID                       uuid.UUID
	UserAID                  uuid.UUID
	UserBID                  uuid.UUID
	Status                   DiplomaticStatus
	ChangedByUserID          uuid.UUID
	ChangedAt                int64
	WarDeclaredAt            *int64
	WarAttacksAllowedAt      *int64
	NeutralityProtectedUntil *int64
}

func canonicalDiplomaticPair(a, b uuid.UUID) (uuid.UUID, uuid.UUID, error) {
	if a == uuid.Nil || b == uuid.Nil || a == b {
		return uuid.Nil, uuid.Nil, NewError("error.domain.diplomacy.invalid_pair", nil)
	}
	if a.String() <= b.String() {
		return a, b, nil
	}
	return b, a, nil
}

func NewUnknownRelationship(user1, user2 uuid.UUID) (*DiplomaticRelationship, error) {
	a, b, err := canonicalDiplomaticPair(user1, user2)
	if err != nil {
		return nil, err
	}
	return &DiplomaticRelationship{
		ID:              uuid.Must(uuid.NewV7()),
		UserAID:         a,
		UserBID:         b,
		Status:          DiplomaticStatusUnknown,
		ChangedByUserID: uuid.Nil,
		ChangedAt:       0,
	}, nil
}

func (r *DiplomaticRelationship) IsUnknown() bool {
	return r.Status == DiplomaticStatusUnknown
}

func (r *DiplomaticRelationship) Involves(userID uuid.UUID) bool {
	return r.UserAID == userID || r.UserBID == userID
}

func (r *DiplomaticRelationship) OtherUser(userID uuid.UUID) uuid.UUID {
	if r.UserAID == userID {
		return r.UserBID
	}
	if r.UserBID == userID {
		return r.UserAID
	}
	return uuid.Nil
}

func (r *DiplomaticRelationship) DeclareWar(changedBy uuid.UUID) error {
	if err := r.CanDeclareWar(); err != nil {
		return err
	}
	now := NowUnix()
	allowedAt := now + int64(DiplomaticWarWarmupDuration.Seconds())
	r.Status = DiplomaticStatusWar
	r.ChangedByUserID = changedBy
	r.ChangedAt = now
	r.WarDeclaredAt = &now
	r.WarAttacksAllowedAt = &allowedAt
	r.NeutralityProtectedUntil = nil
	return nil
}

func (r *DiplomaticRelationship) CanDeclareWar() error {
	if r.Status != DiplomaticStatusNeutral {
		return NewError("error.domain.diplomacy.invalid_transition", nil)
	}
	if r.NeutralityProtectedUntil != nil && NowUnix() < *r.NeutralityProtectedUntil {
		return NewError("error.domain.diplomacy.invalid_transition", nil)
	}
	return nil
}

func (r *DiplomaticRelationship) EscalateToWar(changedBy uuid.UUID) error {
	if !r.IsUnknown() {
		return NewError("error.domain.diplomacy.invalid_transition", nil)
	}
	now := NowUnix()
	r.Status = DiplomaticStatusWar
	r.ChangedByUserID = changedBy
	r.ChangedAt = now
	r.WarDeclaredAt = &now
	r.WarAttacksAllowedAt = &now
	r.NeutralityProtectedUntil = nil
	r.AddEvent(NewDiplomaticRelationshipCreatedEvent(r.ID, r.UserAID, r.UserBID, r.Status, changedBy))
	return nil
}

func (r *DiplomaticRelationship) EstablishContact(changedBy uuid.UUID) error {
	if !r.IsUnknown() {
		return NewError("error.domain.diplomacy.invalid_transition", nil)
	}
	now := NowUnix()
	r.Status = DiplomaticStatusNeutral
	r.ChangedByUserID = changedBy
	r.ChangedAt = now
	r.WarDeclaredAt = nil
	r.WarAttacksAllowedAt = nil
	r.NeutralityProtectedUntil = nil
	r.AddEvent(NewDiplomaticRelationshipCreatedEvent(r.ID, r.UserAID, r.UserBID, r.Status, changedBy))
	return nil
}

func (r *DiplomaticRelationship) FormAlliance(changedBy uuid.UUID) {
	now := NowUnix()
	r.Status = DiplomaticStatusAllied
	r.ChangedByUserID = changedBy
	r.ChangedAt = now
	r.WarDeclaredAt = nil
	r.WarAttacksAllowedAt = nil
	r.NeutralityProtectedUntil = nil
}

func (r *DiplomaticRelationship) CanBreakAlliance() error {
	if r.Status != DiplomaticStatusAllied {
		return NewError("error.domain.diplomacy.invalid_transition", nil)
	}
	return nil
}

func (r *DiplomaticRelationship) BreakAlliance(changedBy uuid.UUID) error {
	if err := r.CanBreakAlliance(); err != nil {
		return err
	}
	now := NowUnix()
	protectedUntil := now + int64(DiplomaticCeasefireProtectionDuration.Seconds())
	r.Status = DiplomaticStatusNeutral
	r.ChangedByUserID = changedBy
	r.ChangedAt = now
	r.WarDeclaredAt = nil
	r.WarAttacksAllowedAt = nil
	r.NeutralityProtectedUntil = &protectedUntil
	return nil
}

func (r *DiplomaticRelationship) AcceptCeasefire(changedBy uuid.UUID) {
	now := NowUnix()
	protectedUntil := now + int64(DiplomaticCeasefireProtectionDuration.Seconds())
	r.Status = DiplomaticStatusNeutral
	r.ChangedByUserID = changedBy
	r.ChangedAt = now
	r.WarDeclaredAt = nil
	r.WarAttacksAllowedAt = nil
	r.NeutralityProtectedUntil = &protectedUntil
}

func (r *DiplomaticRelationship) CanPerformAttackOperation() error {
	if r.IsUnknown() {
		return nil
	}
	if r.Status != DiplomaticStatusWar {
		return NewError("error.domain.diplomacy.attack_requires_war", nil)
	}
	if r.WarAttacksAllowedAt != nil && NowUnix() < *r.WarAttacksAllowedAt {
		return NewError("error.domain.diplomacy.attack_warmup_active", nil)
	}
	return nil
}

func (r *DiplomaticRelationship) CanCreateAllianceRequest() error {
	if r.IsUnknown() {
		return nil
	}
	if r.Status != DiplomaticStatusNeutral {
		return NewError("error.domain.diplomacy.invalid_transition", nil)
	}
	return nil
}

func (r *DiplomaticRelationship) CanCreateCeasefireRequest() error {
	if r.Status != DiplomaticStatusWar {
		return NewError("error.domain.diplomacy.invalid_transition", nil)
	}
	return nil
}
