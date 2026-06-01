package repo

import (
	"context"

	dbgen "github.com/artcodefun/heat-expansion-server/internal/game/infrastructure/db/gen"
)

type TranslationRepo struct {
	q *dbgen.Queries
}

func NewTranslationRepo(q *dbgen.Queries) *TranslationRepo {
	return &TranslationRepo{q: q}
}

func (r *TranslationRepo) GetAll(ctx context.Context) ([]dbgen.Translation, error) {
	return r.q.GetAllTranslations(ctx)
}
