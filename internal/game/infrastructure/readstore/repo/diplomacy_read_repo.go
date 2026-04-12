package repo

import (
	"context"
	"database/sql"
	"errors"

	readmodels "github.com/artcodefun/heat-expansion-server/internal/game/application/cqrs/readmodels"
	"github.com/artcodefun/heat-expansion-server/internal/game/application/ports"
	"github.com/artcodefun/heat-expansion-server/internal/game/infrastructure/readstore/gen"
	"github.com/artcodefun/heat-expansion-server/internal/game/infrastructure/readstore/mappers"
	"github.com/google/uuid"
)

type DiplomacyReadRepo struct {
	q     *gen.Queries
	bases ports.BaseReadRepository
}

func NewDiplomacyReadRepo(q *gen.Queries, bases ports.BaseReadRepository) *DiplomacyReadRepo {
	return &DiplomacyReadRepo{q: q, bases: bases}
}

func (r *DiplomacyReadRepo) ListRelationships(ctx context.Context, userID uuid.UUID, status *readmodels.DiplomaticStatus) ([]*readmodels.DiplomaticRelationship, error) {
	statusFilter := ""
	if status != nil {
		statusFilter = string(*status)
	}
	rows, err := r.q.ListDiplomaticRelationships(ctx, gen.ListDiplomaticRelationshipsParams{CurrentUserID: userID, Status: statusFilter})
	if err != nil {
		return nil, err
	}
	return mappers.DiplomaticRelationshipsFromModels(rows), nil
}

func (r *DiplomacyReadRepo) GetRelationship(ctx context.Context, userID, otherUserID uuid.UUID) (*readmodels.DiplomaticRelationship, error) {
	row, err := r.q.GetDiplomaticRelationship(ctx, gen.GetDiplomaticRelationshipParams{CurrentUserID: userID, OtherUserID: otherUserID})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ports.ErrNotFound
		}
		return nil, err
	}
	return mappers.OneDiplomaticRelationshipFromModel(row), nil
}

func (r *DiplomacyReadRepo) ListChats(ctx context.Context, userID uuid.UUID) ([]*readmodels.DiplomaticChat, error) {
	rows, err := r.q.ListDiplomaticChats(ctx, userID)
	if err != nil {
		return nil, err
	}
	out := mappers.DiplomaticChatsFromRows(rows)
	requestCache := make(map[uuid.UUID]*readmodels.DiplomaticRequest)
	for _, item := range out {
		if item.LastMessage == nil {
			continue
		}
		if err := r.enrichMessage(ctx, item.LastMessage, requestCache); err != nil {
			return nil, err
		}
	}
	return out, nil
}

func (r *DiplomacyReadRepo) GetUnreadMessagesCount(ctx context.Context, userID uuid.UUID) (int, error) {
	count, err := r.q.CountUnreadDiplomaticMessagesByUser(ctx, userID)
	if err != nil {
		return 0, err
	}
	return int(count), nil
}

func (r *DiplomacyReadRepo) ListChatMessages(ctx context.Context, userID, otherUserID uuid.UUID) ([]*readmodels.DiplomaticMessage, error) {
	rows, err := r.q.ListDiplomaticMessagesByChat(ctx, gen.ListDiplomaticMessagesByChatParams{CurrentUserID: userID, OtherUserID: otherUserID})
	if err != nil {
		return nil, err
	}
	out := mappers.DiplomaticMessagesFromChatRows(rows)
	requestCache := make(map[uuid.UUID]*readmodels.DiplomaticRequest)
	for _, item := range out {
		if err := r.enrichMessage(ctx, item, requestCache); err != nil {
			return nil, err
		}
	}
	return out, nil
}

func (r *DiplomacyReadRepo) ListPendingRequests(ctx context.Context, userID uuid.UUID) ([]*readmodels.DiplomaticRequest, error) {
	rows, err := r.q.ListPendingDiplomaticRequests(ctx, userID)
	if err != nil {
		return nil, err
	}
	return mappers.DiplomaticRequestsFromPendingRows(rows), nil
}

func (r *DiplomacyReadRepo) enrichMessage(ctx context.Context, item *readmodels.DiplomaticMessage, requestCache map[uuid.UUID]*readmodels.DiplomaticRequest) error {
	if item == nil {
		return nil
	}
	if item.SenderBaseID != nil {
		base, err := r.bases.GetBase(ctx, *item.SenderBaseID)
		if err != nil {
			if !errors.Is(err, ports.ErrNotFound) {
				return err
			}
		} else {
			item.SenderBase = &readmodels.BaseLocationData{
				Coordinates: base.Coordinates,
				Details:     base.LocationDetails,
			}
		}
	}
	if item.RequestID == nil {
		return nil
	}
	requestID := *item.RequestID
	if request, ok := requestCache[requestID]; ok {
		item.Request = request
		return nil
	}
	row, err := r.q.GetDiplomaticRequest(ctx, requestID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}
		return err
	}
	request := mappers.DiplomaticRequestFromGetRow(row)
	requestCache[requestID] = request
	item.Request = request
	return nil
}
