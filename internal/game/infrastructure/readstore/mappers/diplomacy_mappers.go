package mappers

import (
	readmodels "github.com/artcodefun/heat-expansion-server/internal/game/application/cqrs/readmodels"
	"github.com/artcodefun/heat-expansion-server/internal/game/infrastructure/readstore/gen"
)

func DiplomaticRelationshipFromModel(row gen.ListDiplomaticRelationshipsRow) *readmodels.DiplomaticRelationship {
	return &readmodels.DiplomaticRelationship{
		OtherUserID:              interfaceUUID(row.OtherUserID),
		OtherUserName:            row.OtherUserName,
		Status:                   readmodels.DiplomaticStatus(row.Status),
		ChangedByUserID:          row.ChangedByUserID,
		ChangedAt:                row.ChangedAt,
		WarDeclaredAt:            nullInt64ToInt64Ptr(row.WarDeclaredAt),
		WarAttacksAllowedAt:      nullInt64ToInt64Ptr(row.WarAttacksAllowedAt),
		NeutralityProtectedUntil: nullInt64ToInt64Ptr(row.NeutralityProtectedUntil),
	}
}

func DiplomaticRelationshipsFromModels(rows []gen.ListDiplomaticRelationshipsRow) []*readmodels.DiplomaticRelationship {
	out := make([]*readmodels.DiplomaticRelationship, 0, len(rows))
	for _, row := range rows {
		out = append(out, DiplomaticRelationshipFromModel(row))
	}
	return out
}

func OneDiplomaticRelationshipFromModel(row gen.GetDiplomaticRelationshipRow) *readmodels.DiplomaticRelationship {
	return &readmodels.DiplomaticRelationship{
		OtherUserID:              interfaceUUID(row.OtherUserID),
		OtherUserName:            row.OtherUserName,
		Status:                   readmodels.DiplomaticStatus(row.Status),
		ChangedByUserID:          row.ChangedByUserID,
		ChangedAt:                row.ChangedAt,
		WarDeclaredAt:            nullInt64ToInt64Ptr(row.WarDeclaredAt),
		WarAttacksAllowedAt:      nullInt64ToInt64Ptr(row.WarAttacksAllowedAt),
		NeutralityProtectedUntil: nullInt64ToInt64Ptr(row.NeutralityProtectedUntil),
	}
}

func DiplomaticMessageFromChatRow(row gen.ListDiplomaticMessagesByChatRow) *readmodels.DiplomaticMessage {
	return &readmodels.DiplomaticMessage{
		ID:               row.ID,
		SenderUserID:     row.SenderUserID,
		SenderUserName:   row.SenderUserName,
		ReceiverUserID:   row.ReceiverUserID,
		ReceiverUserName: row.ReceiverUserName,
		SenderBaseID:     nullBaseIDPtr(row.SenderBaseID),
		ReceiverBaseID:   nullBaseIDPtr(row.ReceiverBaseID),
		RequestID:        nullUUIDPtr(row.RequestID),
		ReplyToMessageID: nullUUIDPtr(row.ReplyToMessageID),
		IsRead:           row.IsRead,
		Content:          readmodels.DiplomaticMessageContent(row.Content),
		CreatedAt:        row.CreatedAt,
	}
}

func DiplomaticChatFromRow(row gen.ListDiplomaticChatsRow) *readmodels.DiplomaticChat {
	return &readmodels.DiplomaticChat{
		OtherUserID:   interfaceUUID(row.OtherUserID),
		OtherUserName: row.OtherUserName,
		LastMessage: &readmodels.DiplomaticMessage{
			ID:               row.ID,
			SenderUserID:     row.SenderUserID,
			SenderUserName:   row.SenderUserName,
			ReceiverUserID:   row.ReceiverUserID,
			ReceiverUserName: row.ReceiverUserName,
			SenderBaseID:     nullBaseIDPtr(row.SenderBaseID),
			ReceiverBaseID:   nullBaseIDPtr(row.ReceiverBaseID),
			RequestID:        nullUUIDPtr(row.RequestID),
			ReplyToMessageID: nullUUIDPtr(row.ReplyToMessageID),
			IsRead:           row.IsRead,
			Content:          readmodels.DiplomaticMessageContent(row.Content),
			CreatedAt:        row.CreatedAt,
		},
		UnreadCount: int(row.UnreadCount),
	}
}

func DiplomaticChatsFromRows(rows []gen.ListDiplomaticChatsRow) []*readmodels.DiplomaticChat {
	out := make([]*readmodels.DiplomaticChat, 0, len(rows))
	for _, row := range rows {
		out = append(out, DiplomaticChatFromRow(row))
	}
	return out
}

func DiplomaticMessagesFromChatRows(rows []gen.ListDiplomaticMessagesByChatRow) []*readmodels.DiplomaticMessage {
	out := make([]*readmodels.DiplomaticMessage, 0, len(rows))
	for _, row := range rows {
		out = append(out, DiplomaticMessageFromChatRow(row))
	}
	return out
}

func DiplomaticRequestFromPendingRow(row gen.ListPendingDiplomaticRequestsRow) *readmodels.DiplomaticRequest {
	return &readmodels.DiplomaticRequest{
		ID:               row.ID,
		SenderUserID:     row.SenderUserID,
		SenderUserName:   row.SenderUserName,
		ReceiverUserID:   row.ReceiverUserID,
		ReceiverUserName: row.ReceiverUserName,
		SenderBaseID:     nullBaseIDPtr(row.SenderBaseID),
		ReceiverBaseID:   nullBaseIDPtr(row.ReceiverBaseID),
		Kind:             readmodels.DiplomaticRequestKind(row.Kind),
		Status:           readmodels.DiplomaticRequestStatus(row.Status),
		CreatedAt:        row.CreatedAt,
		ResolvedAt:       nullInt64ToInt64Ptr(row.ResolvedAt),
		ExpiresAt:        row.ExpiresAt,
	}
}

func DiplomaticRequestFromGetRow(row gen.GetDiplomaticRequestRow) *readmodels.DiplomaticRequest {
	return &readmodels.DiplomaticRequest{
		ID:               row.ID,
		SenderUserID:     row.SenderUserID,
		SenderUserName:   row.SenderUserName,
		ReceiverUserID:   row.ReceiverUserID,
		ReceiverUserName: row.ReceiverUserName,
		SenderBaseID:     nullBaseIDPtr(row.SenderBaseID),
		ReceiverBaseID:   nullBaseIDPtr(row.ReceiverBaseID),
		Kind:             readmodels.DiplomaticRequestKind(row.Kind),
		Status:           readmodels.DiplomaticRequestStatus(row.Status),
		CreatedAt:        row.CreatedAt,
		ResolvedAt:       nullInt64ToInt64Ptr(row.ResolvedAt),
		ExpiresAt:        row.ExpiresAt,
	}
}

func DiplomaticRequestsFromPendingRows(rows []gen.ListPendingDiplomaticRequestsRow) []*readmodels.DiplomaticRequest {
	out := make([]*readmodels.DiplomaticRequest, 0, len(rows))
	for _, row := range rows {
		out = append(out, DiplomaticRequestFromPendingRow(row))
	}
	return out
}
