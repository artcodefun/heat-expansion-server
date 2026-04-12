package mappers

import (
	"github.com/artcodefun/heat-expansion-server/internal/game/domain"
	"github.com/artcodefun/heat-expansion-server/internal/game/infrastructure/db/gen"
)

func InsertDiplomaticRelationshipParamsFromDomain(relationship *domain.DiplomaticRelationship) gen.InsertDiplomaticRelationshipParams {
	return gen.InsertDiplomaticRelationshipParams{
		ID:                       relationship.ID,
		UserAID:                  relationship.UserAID,
		UserBID:                  relationship.UserBID,
		Status:                   string(relationship.Status),
		ChangedByUserID:          relationship.ChangedByUserID,
		ChangedAt:                relationship.ChangedAt,
		WarDeclaredAt:            int64PtrToNullInt64(relationship.WarDeclaredAt),
		WarAttacksAllowedAt:      int64PtrToNullInt64(relationship.WarAttacksAllowedAt),
		NeutralityProtectedUntil: int64PtrToNullInt64(relationship.NeutralityProtectedUntil),
	}
}

func UpdateDiplomaticRelationshipParamsFromDomain(relationship *domain.DiplomaticRelationship) gen.UpdateDiplomaticRelationshipParams {
	return gen.UpdateDiplomaticRelationshipParams{
		UserAID:                  relationship.UserAID,
		UserBID:                  relationship.UserBID,
		Status:                   string(relationship.Status),
		ChangedByUserID:          relationship.ChangedByUserID,
		ChangedAt:                relationship.ChangedAt,
		WarDeclaredAt:            int64PtrToNullInt64(relationship.WarDeclaredAt),
		WarAttacksAllowedAt:      int64PtrToNullInt64(relationship.WarAttacksAllowedAt),
		NeutralityProtectedUntil: int64PtrToNullInt64(relationship.NeutralityProtectedUntil),
	}
}

func DiplomaticRelationshipFromDB(row gen.GameDiplomaticRelationship) *domain.DiplomaticRelationship {
	return &domain.DiplomaticRelationship{
		ID:                       row.ID,
		UserAID:                  row.UserAID,
		UserBID:                  row.UserBID,
		Status:                   domain.DiplomaticStatus(row.Status),
		ChangedByUserID:          row.ChangedByUserID,
		ChangedAt:                row.ChangedAt,
		WarDeclaredAt:            nullInt64ToInt64Ptr(row.WarDeclaredAt),
		WarAttacksAllowedAt:      nullInt64ToInt64Ptr(row.WarAttacksAllowedAt),
		NeutralityProtectedUntil: nullInt64ToInt64Ptr(row.NeutralityProtectedUntil),
	}
}

func InsertDiplomaticRequestParamsFromDomain(request *domain.DiplomaticRequest) gen.InsertDiplomaticRequestParams {
	return gen.InsertDiplomaticRequestParams{
		ID:             request.ID,
		SenderUserID:   request.SenderUserID,
		ReceiverUserID: request.ReceiverUserID,
		SenderBaseID:   nullableBaseID(request.SenderBaseID),
		ReceiverBaseID: nullableBaseID(request.ReceiverBaseID),
		Kind:           string(request.Kind),
		Status:         string(request.Status),
		CreatedAt:      request.CreatedAt,
		ResolvedAt:     int64PtrToNullInt64(request.ResolvedAt),
		ExpiresAt:      request.ExpiresAt,
	}
}

func UpdateDiplomaticRequestParamsFromDomain(request *domain.DiplomaticRequest) gen.UpdateDiplomaticRequestParams {
	return gen.UpdateDiplomaticRequestParams{
		ID:         request.ID,
		Status:     string(request.Status),
		ResolvedAt: int64PtrToNullInt64(request.ResolvedAt),
	}
}

func DiplomaticRequestFromDB(row gen.GameDiplomaticRequest) *domain.DiplomaticRequest {
	return &domain.DiplomaticRequest{
		ID:             row.ID,
		SenderUserID:   row.SenderUserID,
		ReceiverUserID: row.ReceiverUserID,
		SenderBaseID:   nullInt64ToBaseIDPtr(row.SenderBaseID),
		ReceiverBaseID: nullInt64ToBaseIDPtr(row.ReceiverBaseID),
		Kind:           domain.DiplomaticRequestKind(row.Kind),
		Status:         domain.DiplomaticRequestStatus(row.Status),
		CreatedAt:      row.CreatedAt,
		ResolvedAt:     nullInt64ToInt64Ptr(row.ResolvedAt),
		ExpiresAt:      row.ExpiresAt,
	}
}

func InsertDiplomaticMessageParamsFromDomain(message *domain.DiplomaticMessage) gen.InsertDiplomaticMessageParams {
	return gen.InsertDiplomaticMessageParams{
		ID:               message.ID,
		SenderUserID:     message.SenderUserID,
		ReceiverUserID:   message.ReceiverUserID,
		SenderBaseID:     nullableBaseID(message.SenderBaseID),
		ReceiverBaseID:   nullableBaseID(message.ReceiverBaseID),
		RequestID:        nullableUUID(message.RequestID),
		ReplyToMessageID: nullableUUID(message.ReplyToMessageID),
		IsRead:           message.IsRead,
		Content:          message.Content,
		CreatedAt:        message.CreatedAt,
	}
}

func DiplomaticMessageFromDB(row gen.GameDiplomaticMessage) *domain.DiplomaticMessage {
	return &domain.DiplomaticMessage{
		ID:               row.ID,
		SenderUserID:     row.SenderUserID,
		ReceiverUserID:   row.ReceiverUserID,
		SenderBaseID:     nullInt64ToBaseIDPtr(row.SenderBaseID),
		ReceiverBaseID:   nullInt64ToBaseIDPtr(row.ReceiverBaseID),
		RequestID:        nullUUIDToUUIDPtr(row.RequestID),
		ReplyToMessageID: nullUUIDToUUIDPtr(row.ReplyToMessageID),
		IsRead:           row.IsRead,
		Content:          row.Content,
		CreatedAt:        row.CreatedAt,
	}
}
