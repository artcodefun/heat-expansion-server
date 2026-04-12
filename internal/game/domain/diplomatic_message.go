package domain

import (
	"slices"

	"github.com/google/uuid"
)

const (
	DiplomaticMessageContentCoalitionProposal     TranslationKey = "diplomacy.message.coalition_proposal.content"
	DiplomaticMessageContentCeasefireProposal     TranslationKey = "diplomacy.message.ceasefire_proposal.content"
	DiplomaticMessageContentCoalitionBreakNotice  TranslationKey = "diplomacy.message.coalition_break_notice.content"
	DiplomaticMessageContentCoalitionAcceptance   TranslationKey = "diplomacy.message.coalition_acceptance.content"
	DiplomaticMessageContentCoalitionRejection    TranslationKey = "diplomacy.message.coalition_rejection.content"
	DiplomaticMessageContentWarDeclaration        TranslationKey = "diplomacy.message.war_declaration.content"
	DiplomaticMessageContentCeasefireAcceptance   TranslationKey = "diplomacy.message.ceasefire_acceptance.content"
	DiplomaticMessageContentCeasefireRejection    TranslationKey = "diplomacy.message.ceasefire_rejection.content"
	DiplomaticMessageContentGreetingFriendly      TranslationKey = "diplomacy.message.greeting_friendly_hello.content"
	DiplomaticMessageContentGreetingFormal        TranslationKey = "diplomacy.message.greeting_formal_introduction.content"
	DiplomaticMessageContentGreetingAdmiration    TranslationKey = "diplomacy.message.greeting_sign_of_admiration.content"
	DiplomaticMessageContentWarningKeepOut        TranslationKey = "diplomacy.message.warning_keep_out_of_our_sector.content"
	DiplomaticMessageContentWarningDoNotInterfere TranslationKey = "diplomacy.message.warning_do_not_interfere.content"
	DiplomaticMessageContentWarningGeneralThreat  TranslationKey = "diplomacy.message.warning_general_threat.content"
)

var DiplomaticGreetingMessageContents = []TranslationKey{
	DiplomaticMessageContentGreetingFriendly,
	DiplomaticMessageContentGreetingFormal,
	DiplomaticMessageContentGreetingAdmiration,
}

var DiplomaticWarningMessageContents = []TranslationKey{
	DiplomaticMessageContentWarningKeepOut,
	DiplomaticMessageContentWarningDoNotInterfere,
	DiplomaticMessageContentWarningGeneralThreat,
}

var systemDiplomaticMessageContents = []TranslationKey{
	DiplomaticMessageContentCoalitionProposal,
	DiplomaticMessageContentCeasefireProposal,
	DiplomaticMessageContentWarDeclaration,
	DiplomaticMessageContentCoalitionBreakNotice,
	DiplomaticMessageContentCoalitionAcceptance,
	DiplomaticMessageContentCoalitionRejection,
	DiplomaticMessageContentCeasefireAcceptance,
	DiplomaticMessageContentCeasefireRejection,
}

type DiplomaticMessage struct {
	EventProducer

	ID               uuid.UUID
	SenderUserID     uuid.UUID
	ReceiverUserID   uuid.UUID
	SenderBaseID     *int
	ReceiverBaseID   *int
	RequestID        *uuid.UUID
	ReplyToMessageID *uuid.UUID
	IsRead           bool
	Content          TranslationKey
	CreatedAt        int64
}

func ValidateDiplomaticParticipants(senderUserID, receiverUserID uuid.UUID) error {
	if senderUserID == uuid.Nil || receiverUserID == uuid.Nil {
		return NewError("error.domain.diplomacy.invalid_participants", nil)
	}
	if senderUserID == receiverUserID {
		return NewError("error.domain.diplomacy.self_targeting_forbidden", nil)
	}
	return nil
}

func NewDiplomaticMessage(senderUserID, receiverUserID uuid.UUID, senderBaseID, receiverBaseID *int, requestID, replyToMessageID *uuid.UUID, content TranslationKey) (*DiplomaticMessage, error) {
	if err := ValidateDiplomaticParticipants(senderUserID, receiverUserID); err != nil {
		return nil, err
	}
	if !isKnownDiplomaticMessageContent(content) {
		return nil, NewError("error.domain.diplomacy.invalid_message_kind", nil)
	}
	msg := &DiplomaticMessage{
		ID:               uuid.Must(uuid.NewV7()),
		SenderUserID:     senderUserID,
		ReceiverUserID:   receiverUserID,
		SenderBaseID:     senderBaseID,
		ReceiverBaseID:   receiverBaseID,
		RequestID:        requestID,
		ReplyToMessageID: replyToMessageID,
		IsRead:           false,
		Content:          content,
		CreatedAt:        NowUnix(),
	}
	msg.AddEvent(NewDiplomaticMessageSentEvent(msg.ID, msg.SenderUserID, msg.ReceiverUserID, msg.ReceiverBaseID, msg.Content))
	return msg, nil
}

func InformationalDiplomaticMessageContents() []TranslationKey {
	return slices.Concat(DiplomaticGreetingMessageContents, DiplomaticWarningMessageContents)
}

func IsUserSendableDiplomaticMessageContent(content TranslationKey) bool {
	return slices.Contains(DiplomaticGreetingMessageContents, content) || slices.Contains(DiplomaticWarningMessageContents, content)
}

func isKnownDiplomaticMessageContent(content TranslationKey) bool {
	return IsUserSendableDiplomaticMessageContent(content) || slices.Contains(systemDiplomaticMessageContents, content)
}
