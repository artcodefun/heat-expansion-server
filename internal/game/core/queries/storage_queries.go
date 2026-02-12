package queries

import (
	"github.com/artcodefun/heat-expansion-api/internal/game/core/cqrs"
	"github.com/artcodefun/heat-expansion-api/internal/game/core/cqrs/readmodels"
	"github.com/artcodefun/heat-expansion-api/internal/game/core/ports"
	"github.com/artcodefun/heat-expansion-api/internal/game/core/services"
)

type StorageQueries struct {
	Repo   ports.StorageReadRepository
	Access *services.AccessControlService
}

func NewStorageQueries(repo ports.StorageReadRepository, access *services.AccessControlService) *StorageQueries {
	return &StorageQueries{Repo: repo, Access: access}
}

func (q *StorageQueries) ListPresentStorageItems(ctx cqrs.QueryContext, baseID int, category readmodels.StorageCategory) ([]*readmodels.StorageItemPresent, error) {
	if err := q.Access.EnsureBaseOwnership(ctx.UserID, baseID); err != nil {
		return nil, err
	}
	items, err := q.Repo.ListPresentStorageItems(baseID, category)
	return items, repoErr(err)
}
