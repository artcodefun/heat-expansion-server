package queries

import (
	"context"

	"github.com/artcodefun/heat-expansion-server/internal/game/application/cqrs"
	"github.com/artcodefun/heat-expansion-server/internal/game/application/cqrs/readmodels"
	"github.com/artcodefun/heat-expansion-server/internal/game/application/ports"
	"github.com/artcodefun/heat-expansion-server/internal/game/application/services"
)

type StorageQueries struct {
	Repo   ports.StorageReadRepository
	Access *services.AccessControlService
}

func NewStorageQueries(repo ports.StorageReadRepository, access *services.AccessControlService) *StorageQueries {
	return &StorageQueries{Repo: repo, Access: access}
}

func (q *StorageQueries) ListPresentStorageItems(ctx context.Context, actor cqrs.Actor, baseID int, category readmodels.StorageCategory) ([]*readmodels.StorageItemPresent, error) {
	if err := q.Access.EnsureBaseOwnership(ctx, actor.UserID, baseID); err != nil {
		return nil, err
	}
	items, err := q.Repo.ListPresentStorageItems(ctx, baseID, category)
	return items, repoErr(err)
}
