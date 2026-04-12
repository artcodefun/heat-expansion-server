package readmodels

import "github.com/google/uuid"

type DiplomaticStatus string

const (
	DiplomaticStatusNeutral DiplomaticStatus = "NEUTRAL"
	DiplomaticStatusAllied  DiplomaticStatus = "ALLIED"
	DiplomaticStatusWar     DiplomaticStatus = "WAR"
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

type DiplomaticMessageContent string

const (
	DiplomaticMessageContentCoalitionProposal     DiplomaticMessageContent = "diplomacy.message.coalition_proposal.content"
	DiplomaticMessageContentCeasefireProposal     DiplomaticMessageContent = "diplomacy.message.ceasefire_proposal.content"
	DiplomaticMessageContentCoalitionBreakNotice  DiplomaticMessageContent = "diplomacy.message.coalition_break_notice.content"
	DiplomaticMessageContentCoalitionAcceptance   DiplomaticMessageContent = "diplomacy.message.coalition_acceptance.content"
	DiplomaticMessageContentCoalitionRejection    DiplomaticMessageContent = "diplomacy.message.coalition_rejection.content"
	DiplomaticMessageContentWarDeclaration        DiplomaticMessageContent = "diplomacy.message.war_declaration.content"
	DiplomaticMessageContentCeasefireAcceptance   DiplomaticMessageContent = "diplomacy.message.ceasefire_acceptance.content"
	DiplomaticMessageContentCeasefireRejection    DiplomaticMessageContent = "diplomacy.message.ceasefire_rejection.content"
	DiplomaticMessageContentGreetingFriendly      DiplomaticMessageContent = "diplomacy.message.greeting_friendly_hello.content"
	DiplomaticMessageContentGreetingFormal        DiplomaticMessageContent = "diplomacy.message.greeting_formal_introduction.content"
	DiplomaticMessageContentGreetingAdmiration    DiplomaticMessageContent = "diplomacy.message.greeting_sign_of_admiration.content"
	DiplomaticMessageContentWarningKeepOut        DiplomaticMessageContent = "diplomacy.message.warning_keep_out_of_our_sector.content"
	DiplomaticMessageContentWarningDoNotInterfere DiplomaticMessageContent = "diplomacy.message.warning_do_not_interfere.content"
	DiplomaticMessageContentWarningGeneralThreat  DiplomaticMessageContent = "diplomacy.message.warning_general_threat.content"
)

type BaseLocationData struct {
	Coordinates Vector2i
	Details     LocationDetails
}

type DiplomaticChat struct {
	OtherUserID   uuid.UUID
	OtherUserName string
	LastMessage   *DiplomaticMessage
	UnreadCount   int
}

type DiplomaticRelationship struct {
	OtherUserID              uuid.UUID
	OtherUserName            string
	Status                   DiplomaticStatus
	ChangedByUserID          uuid.UUID
	ChangedAt                int64
	WarDeclaredAt            *int64
	WarAttacksAllowedAt      *int64
	NeutralityProtectedUntil *int64
}

type DiplomaticMessage struct {
	ID               uuid.UUID
	SenderUserID     uuid.UUID
	SenderUserName   string
	ReceiverUserID   uuid.UUID
	ReceiverUserName string
	SenderBaseID     *int
	SenderBase       *BaseLocationData
	ReceiverBaseID   *int
	RequestID        *uuid.UUID
	Request          *DiplomaticRequest
	ReplyToMessageID *uuid.UUID
	IsRead           bool
	Content          DiplomaticMessageContent
	CreatedAt        int64
}

type DiplomaticRequest struct {
	ID               uuid.UUID
	SenderUserID     uuid.UUID
	SenderUserName   string
	ReceiverUserID   uuid.UUID
	ReceiverUserName string
	SenderBaseID     *int
	ReceiverBaseID   *int
	Kind             DiplomaticRequestKind
	Status           DiplomaticRequestStatus
	CreatedAt        int64
	ResolvedAt       *int64
	ExpiresAt        int64
}
