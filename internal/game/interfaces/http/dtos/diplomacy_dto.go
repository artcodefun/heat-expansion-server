package dtos

import (
	readmodels "github.com/artcodefun/heat-expansion-server/internal/game/application/cqrs/readmodels"
	"github.com/artcodefun/heat-expansion-server/internal/game/application/ports"
	"github.com/google/uuid"
)

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
	DiplomaticMessageContentGreetingFriendly      DiplomaticMessageContent = "diplomacy.message.greeting_friendly_hello.content"
	DiplomaticMessageContentGreetingFormal        DiplomaticMessageContent = "diplomacy.message.greeting_formal_introduction.content"
	DiplomaticMessageContentGreetingAdmiration    DiplomaticMessageContent = "diplomacy.message.greeting_sign_of_admiration.content"
	DiplomaticMessageContentWarningKeepOut        DiplomaticMessageContent = "diplomacy.message.warning_keep_out_of_our_sector.content"
	DiplomaticMessageContentWarningDoNotInterfere DiplomaticMessageContent = "diplomacy.message.warning_do_not_interfere.content"
	DiplomaticMessageContentWarningGeneralThreat  DiplomaticMessageContent = "diplomacy.message.warning_general_threat.content"
)

type DiplomaticRelationshipDTO struct {
	OtherUserID              uuid.UUID        `json:"otherUserId"`
	OtherUserName            string           `json:"otherUserName"`
	Status                   DiplomaticStatus `json:"status"`
	ChangedByUserID          uuid.UUID        `json:"changedByUserId"`
	ChangedAt                int64            `json:"changedAt"`
	WarDeclaredAt            *int64           `json:"warDeclaredAt,omitempty"`
	WarAttacksAllowedAt      *int64           `json:"warAttacksAllowedAt,omitempty"`
	NeutralityProtectedUntil *int64           `json:"neutralityProtectedUntil,omitempty"`
}

type DiplomaticMessageDTO struct {
	ID               uuid.UUID             `json:"id"`
	SenderUserID     uuid.UUID             `json:"senderUserId"`
	SenderUserName   string                `json:"senderUserName"`
	ReceiverUserID   uuid.UUID             `json:"receiverUserId"`
	ReceiverUserName string                `json:"receiverUserName"`
	SenderBaseID     *int                  `json:"senderBaseId,omitempty"`
	SenderBase       *BaseLocationDataDTO  `json:"senderBase,omitempty"`
	ReceiverBaseID   *int                  `json:"receiverBaseId,omitempty"`
	RequestID        *uuid.UUID            `json:"requestId,omitempty"`
	Request          *DiplomaticRequestDTO `json:"request,omitempty"`
	ReplyToMessageID *uuid.UUID            `json:"replyToMessageId,omitempty"`
	IsRead           bool                  `json:"isRead"`
	Content          string                `json:"content"`
	CreatedAt        int64                 `json:"createdAt"`
}

type DiplomaticChatDTO struct {
	OtherUserID   uuid.UUID             `json:"otherUserId"`
	OtherUserName string                `json:"otherUserName"`
	LastMessage   *DiplomaticMessageDTO `json:"lastMessage,omitempty"`
	UnreadCount   int                   `json:"unreadCount"`
}

type BaseLocationDataDTO struct {
	Coordinates Vector2iDTO `json:"coordinates"`
	Name        string      `json:"name"`
	Description string      `json:"description,omitempty"`
	ImageURL    string      `json:"imageUrl,omitempty"`
}

type DiplomaticRequestDTO struct {
	ID               uuid.UUID               `json:"id"`
	SenderUserID     uuid.UUID               `json:"senderUserId"`
	SenderUserName   string                  `json:"senderUserName"`
	ReceiverUserID   uuid.UUID               `json:"receiverUserId"`
	ReceiverUserName string                  `json:"receiverUserName"`
	SenderBaseID     *int                    `json:"senderBaseId,omitempty"`
	ReceiverBaseID   *int                    `json:"receiverBaseId,omitempty"`
	Kind             DiplomaticRequestKind   `json:"kind"`
	Status           DiplomaticRequestStatus `json:"status"`
	CreatedAt        int64                   `json:"createdAt"`
	ResolvedAt       *int64                  `json:"resolvedAt,omitempty"`
	ExpiresAt        int64                   `json:"expiresAt"`
}

func DiplomaticStatusPtrFromDTO(value *DiplomaticStatus) *readmodels.DiplomaticStatus {
	if value == nil {
		return nil
	}
	status := readmodels.DiplomaticStatus(*value)
	return &status
}

func DiplomaticRelationshipFromReadModel(item *readmodels.DiplomaticRelationship) DiplomaticRelationshipDTO {
	return DiplomaticRelationshipDTO{
		OtherUserID:              item.OtherUserID,
		OtherUserName:            item.OtherUserName,
		Status:                   DiplomaticStatus(item.Status),
		ChangedByUserID:          item.ChangedByUserID,
		ChangedAt:                item.ChangedAt,
		WarDeclaredAt:            item.WarDeclaredAt,
		WarAttacksAllowedAt:      item.WarAttacksAllowedAt,
		NeutralityProtectedUntil: item.NeutralityProtectedUntil,
	}
}

func DiplomaticRelationshipsFromReadModels(items []*readmodels.DiplomaticRelationship) []DiplomaticRelationshipDTO {
	out := make([]DiplomaticRelationshipDTO, 0, len(items))
	for _, item := range items {
		out = append(out, DiplomaticRelationshipFromReadModel(item))
	}
	return out
}

func DiplomaticMessageFromReadModel(item *readmodels.DiplomaticMessage, tr ports.Translator, locale string) DiplomaticMessageDTO {
	dto := DiplomaticMessageDTO{
		ID:               item.ID,
		SenderUserID:     item.SenderUserID,
		SenderUserName:   item.SenderUserName,
		ReceiverUserID:   item.ReceiverUserID,
		ReceiverUserName: item.ReceiverUserName,
		SenderBaseID:     item.SenderBaseID,
		ReceiverBaseID:   item.ReceiverBaseID,
		RequestID:        item.RequestID,
		ReplyToMessageID: item.ReplyToMessageID,
		IsRead:           item.IsRead,
		Content:          tr.T(locale, string(item.Content), nil),
		CreatedAt:        item.CreatedAt,
	}
	if item.Request != nil {
		request := DiplomaticRequestFromReadModel(item.Request)
		dto.Request = &request
	}
	if item.SenderBase != nil {
		dto.SenderBase = &BaseLocationDataDTO{
			Coordinates: Vector2iFromReadModel(item.SenderBase.Coordinates),
			Name:        tr.T(locale, item.SenderBase.Details.Name, nil),
			Description: tr.T(locale, item.SenderBase.Details.Description, nil),
			ImageURL:    item.SenderBase.Details.ImageURL,
		}
	}
	return dto
}

func DiplomaticMessagesFromReadModels(items []*readmodels.DiplomaticMessage, tr ports.Translator, locale string) []DiplomaticMessageDTO {
	out := make([]DiplomaticMessageDTO, 0, len(items))
	for _, item := range items {
		out = append(out, DiplomaticMessageFromReadModel(item, tr, locale))
	}
	return out
}

func DiplomaticChatFromReadModel(item *readmodels.DiplomaticChat, tr ports.Translator, locale string) DiplomaticChatDTO {
	dto := DiplomaticChatDTO{
		OtherUserID:   item.OtherUserID,
		OtherUserName: item.OtherUserName,
		UnreadCount:   item.UnreadCount,
	}
	if item.LastMessage != nil {
		m := DiplomaticMessageFromReadModel(item.LastMessage, tr, locale)
		dto.LastMessage = &m
	}
	return dto
}

func DiplomaticChatsFromReadModels(items []*readmodels.DiplomaticChat, tr ports.Translator, locale string) []DiplomaticChatDTO {
	out := make([]DiplomaticChatDTO, 0, len(items))
	for _, item := range items {
		out = append(out, DiplomaticChatFromReadModel(item, tr, locale))
	}
	return out
}

func DiplomaticRequestFromReadModel(item *readmodels.DiplomaticRequest) DiplomaticRequestDTO {
	return DiplomaticRequestDTO{
		ID:               item.ID,
		SenderUserID:     item.SenderUserID,
		SenderUserName:   item.SenderUserName,
		ReceiverUserID:   item.ReceiverUserID,
		ReceiverUserName: item.ReceiverUserName,
		SenderBaseID:     item.SenderBaseID,
		ReceiverBaseID:   item.ReceiverBaseID,
		Kind:             DiplomaticRequestKind(item.Kind),
		Status:           DiplomaticRequestStatus(item.Status),
		CreatedAt:        item.CreatedAt,
		ResolvedAt:       item.ResolvedAt,
		ExpiresAt:        item.ExpiresAt,
	}
}

func DiplomaticRequestsFromReadModels(items []*readmodels.DiplomaticRequest) []DiplomaticRequestDTO {
	out := make([]DiplomaticRequestDTO, 0, len(items))
	for _, item := range items {
		out = append(out, DiplomaticRequestFromReadModel(item))
	}
	return out
}

func IsValidDiplomaticStatus(value string) bool {
	switch DiplomaticStatus(value) {
	case DiplomaticStatusNeutral, DiplomaticStatusAllied, DiplomaticStatusWar:
		return true
	default:
		return false
	}
}

func IsValidUserSendableDiplomaticMessageContent(value string) bool {
	switch DiplomaticMessageContent(value) {
	case DiplomaticMessageContentGreetingFriendly,
		DiplomaticMessageContentGreetingFormal,
		DiplomaticMessageContentGreetingAdmiration,
		DiplomaticMessageContentWarningKeepOut,
		DiplomaticMessageContentWarningDoNotInterfere,
		DiplomaticMessageContentWarningGeneralThreat:
		return true
	default:
		return false
	}
}

func IsValidDiplomaticRequestKind(value string) bool {
	switch DiplomaticRequestKind(value) {
	case DiplomaticRequestKindCoalitionProposal, DiplomaticRequestKindCeasefireProposal:
		return true
	default:
		return false
	}
}
