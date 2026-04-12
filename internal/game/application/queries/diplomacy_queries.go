package queries

import (
	"context"

	"github.com/artcodefun/heat-expansion-server/internal/game/application/cqrs"
	readmodels "github.com/artcodefun/heat-expansion-server/internal/game/application/cqrs/readmodels"
	"github.com/artcodefun/heat-expansion-server/internal/game/application/ports"
	"github.com/google/uuid"
)

type DiplomacyQueries struct {
	Repo ports.DiplomacyReadRepository
}

func NewDiplomacyQueries(repo ports.DiplomacyReadRepository) *DiplomacyQueries {
	return &DiplomacyQueries{Repo: repo}
}

func (q *DiplomacyQueries) ListRelationships(ctx context.Context, actor cqrs.Actor, status *readmodels.DiplomaticStatus) ([]*readmodels.DiplomaticRelationship, error) {
	if actor.UserID == uuid.Nil {
		return nil, cqrs.ErrForbidden
	}
	items, err := q.Repo.ListRelationships(ctx, actor.UserID, status)
	return items, repoErr(err)
}

func (q *DiplomacyQueries) GetRelationship(ctx context.Context, actor cqrs.Actor, otherUserID uuid.UUID) (*readmodels.DiplomaticRelationship, error) {
	if actor.UserID == uuid.Nil {
		return nil, cqrs.ErrForbidden
	}
	item, err := q.Repo.GetRelationship(ctx, actor.UserID, otherUserID)
	return item, repoErr(err)
}

func (q *DiplomacyQueries) ListChats(ctx context.Context, actor cqrs.Actor) ([]*readmodels.DiplomaticChat, error) {
	if actor.UserID == uuid.Nil {
		return nil, cqrs.ErrForbidden
	}
	items, err := q.Repo.ListChats(ctx, actor.UserID)
	return items, repoErr(err)
}

func (q *DiplomacyQueries) GetUnreadMessagesCount(ctx context.Context, actor cqrs.Actor) (int, error) {
	if actor.UserID == uuid.Nil {
		return 0, cqrs.ErrForbidden
	}
	count, err := q.Repo.GetUnreadMessagesCount(ctx, actor.UserID)
	return count, repoErr(err)
}

func (q *DiplomacyQueries) ListChatMessages(ctx context.Context, actor cqrs.Actor, otherUserID uuid.UUID) ([]*readmodels.DiplomaticMessage, error) {
	if actor.UserID == uuid.Nil {
		return nil, cqrs.ErrForbidden
	}
	items, err := q.Repo.ListChatMessages(ctx, actor.UserID, otherUserID)
	return items, repoErr(err)
}

func (q *DiplomacyQueries) ListPendingRequests(ctx context.Context, actor cqrs.Actor) ([]*readmodels.DiplomaticRequest, error) {
	if actor.UserID == uuid.Nil {
		return nil, cqrs.ErrForbidden
	}
	items, err := q.Repo.ListPendingRequests(ctx, actor.UserID)
	return items, repoErr(err)
}
