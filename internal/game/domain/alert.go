package domain

import (
	"time"

	"github.com/google/uuid"
)

type AlertKind string

const (
	AlertKindCombat     AlertKind = "COMBAT"
	AlertKindIntel      AlertKind = "INTEL"
	AlertKindSystem     AlertKind = "SYSTEM"
	AlertKindDiplomatic AlertKind = "DIPLOMATIC"
	AlertKindTrade      AlertKind = "TRADE"
)

type Alert struct {
	ID         uuid.UUID
	UserID     uuid.UUID
	BaseID     *int
	ActivityID *uuid.UUID
	Kind       AlertKind
	Title      TranslationKey
	Content    TranslationKey
	IsRead     bool
	CreatedAt  int64
	ExpiresAt  int64
}

const defaultAlertTTL = 72 * time.Hour

func NewAlert(userID uuid.UUID, baseID *int, activityID *uuid.UUID, kind AlertKind, title, content TranslationKey, ttl time.Duration) *Alert {
	now := NowUnix()
	return &Alert{
		ID:         uuid.Must(uuid.NewV7()),
		UserID:     userID,
		BaseID:     baseID,
		ActivityID: activityID,
		Kind:       kind,
		Title:      title,
		Content:    content,
		IsRead:     false,
		CreatedAt:  now,
		ExpiresAt:  now + int64(ttl.Seconds()),
	}
}

func NewActivityAlert(event ActivityCreatedEvent) (*Alert, bool) {
	var kind AlertKind
	var title, content TranslationKey

	switch event.Kind {
	case ActivityKindDefense:
		kind = AlertKindCombat
		title = "alert.combat.attack.title"
		if event.Subtype == string(DefenseActivitySubtypeSpy) {
			content = "alert.combat.spy.content"
		} else {
			content = "alert.combat.attack.content"
		}
	case ActivityKindScan:
		if event.Subtype != string(ScanActivitySubtypeExternalScanDetected) {
			return nil, false
		}
		kind = AlertKindIntel
		title = "alert.intel.scan.title"
		content = "alert.intel.scan.content"
	case ActivityKindRadar:
		kind = AlertKindIntel
		title = "alert.intel.threat.title"
		content = "alert.intel.threat.content"
	default:
		return nil, false
	}

	baseID := event.BaseID
	return NewAlert(event.UserID, &baseID, &event.ActivityID, kind, title, content, defaultAlertTTL), true
}

type TradeAlertKind string

const (
	TradeAlertKindCreated   TradeAlertKind = "CREATED"
	TradeAlertKindAccepted  TradeAlertKind = "ACCEPTED"
	TradeAlertKindDeclined  TradeAlertKind = "DECLINED"
	TradeAlertKindCancelled TradeAlertKind = "CANCELLED"
	TradeAlertKindExpired   TradeAlertKind = "EXPIRED"
	TradeAlertKindCompleted TradeAlertKind = "COMPLETED"
)

func NewTradeAlert(op *TradeOperation, isSender bool, kind TradeAlertKind) *Alert {
	if op == nil {
		return nil
	}

	var userID uuid.UUID
	var baseID int
	if isSender {
		userID = op.SenderUserID
		baseID = op.SenderBaseID
	} else {
		userID = op.ReceiverUserID
		baseID = op.ReceiverBaseID
	}

	var title, content TranslationKey
	switch kind {
	case TradeAlertKindCreated:
		title = "alert.trade.created.title"
		content = "alert.trade.created.content"
	case TradeAlertKindAccepted:
		title = "alert.trade.accepted.title"
		content = "alert.trade.accepted.content"
	case TradeAlertKindDeclined:
		title = "alert.trade.declined.title"
		content = "alert.trade.declined.content"
	case TradeAlertKindCancelled:
		title = "alert.trade.cancelled.title"
		content = "alert.trade.cancelled.content"
	case TradeAlertKindExpired:
		title = "alert.trade.expired.title"
		content = "alert.trade.expired.content"
	case TradeAlertKindCompleted:
		title = "alert.trade.completed.title"
		content = "alert.trade.completed.content"
	default:
		return nil
	}

	return NewAlert(userID, &baseID, nil, AlertKindTrade, title, content, defaultAlertTTL)
}

func NewDiplomaticMessageAlert(event DiplomaticMessageSentEvent) (*Alert, bool) {
	var title, content TranslationKey

	switch event.Content {
	case DiplomaticMessageContentWarDeclaration:
		title = "alert.diplomacy.war_declared.title"
		content = "alert.diplomacy.war_declared.content"
	case DiplomaticMessageContentCoalitionAcceptance:
		title = "alert.diplomacy.coalition_accepted.title"
		content = "alert.diplomacy.coalition_accepted.content"
	case DiplomaticMessageContentCoalitionRejection:
		title = "alert.diplomacy.coalition_rejected.title"
		content = "alert.diplomacy.coalition_rejected.content"
	case DiplomaticMessageContentCeasefireAcceptance:
		title = "alert.diplomacy.ceasefire_accepted.title"
		content = "alert.diplomacy.ceasefire_accepted.content"
	case DiplomaticMessageContentCeasefireRejection:
		title = "alert.diplomacy.ceasefire_rejected.title"
		content = "alert.diplomacy.ceasefire_rejected.content"
	case DiplomaticMessageContentCoalitionBreakNotice:
		title = "alert.diplomacy.coalition_broken.title"
		content = "alert.diplomacy.coalition_broken.content"
	default:
		if !IsUserSendableDiplomaticMessageContent(event.Content) {
			return nil, false
		}
		title = "alert.diplomacy.informational.title"
		content = "alert.diplomacy.informational.content"
	}

	return NewAlert(event.ReceiverUserID, event.ReceiverBaseID, nil, AlertKindDiplomatic, title, content, defaultAlertTTL), true
}

func NewDiplomaticRequestAlert(event DiplomaticRequestCreatedEvent) (*Alert, bool) {
	var title, content TranslationKey

	switch event.Kind {
	case DiplomaticRequestKindCoalitionProposal:
		title = "alert.diplomacy.coalition_proposal.title"
		content = "alert.diplomacy.coalition_proposal.content"
	case DiplomaticRequestKindCeasefireProposal:
		title = "alert.diplomacy.ceasefire_proposal.title"
		content = "alert.diplomacy.ceasefire_proposal.content"
	default:
		return nil, false
	}

	return NewAlert(event.ReceiverUserID, event.ReceiverBaseID, nil, AlertKindDiplomatic, title, content, defaultAlertTTL), true
}
