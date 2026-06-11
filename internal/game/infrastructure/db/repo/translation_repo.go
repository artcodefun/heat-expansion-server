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

func (r *TranslationRepo) Upsert(ctx context.Context, key, locale, value string) error {
	if err := r.q.UpsertTranslation(ctx, dbgen.UpsertTranslationParams{
		Key:    key,
		Locale: locale,
		Value:  value,
	}); err != nil {
		return err
	}
	return r.q.NotifyTranslationsChanged(ctx)
}
