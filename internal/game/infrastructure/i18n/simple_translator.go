package i18n

import (
	"context"
	"fmt"

	"github.com/artcodefun/heat-expansion-server/internal/game/application/ports"
	repo "github.com/artcodefun/heat-expansion-server/internal/game/infrastructure/db/repo"
	"github.com/artcodefun/heat-expansion-server/internal/game/infrastructure/i18n/locales"
	platformi18n "github.com/artcodefun/heat-expansion-server/internal/platform/i18n"
)

type SimpleTranslator struct {
	platformi18n.Translator
	repo *repo.TranslationRepo
}

func NewSimpleTranslator(r *repo.TranslationRepo) (*SimpleTranslator, error) {
	t := &SimpleTranslator{repo: r}
	if err := t.LoadFromJsonFiles(locales.Files); err != nil {
		return nil, err
	}
	return t, nil
}

// LoadFromRepo loads content translations from the repository, overriding any
// systemic translations for the same key+locale.
func (t *SimpleTranslator) LoadFromRepo(ctx context.Context) error {
	entries, err := t.repo.GetAll(ctx)
	if err != nil {
		return fmt.Errorf("failed to load content translations: %w", err)
	}
	byLocale := make(map[string]map[string]string)
	for _, e := range entries {
		if byLocale[e.Locale] == nil {
			byLocale[e.Locale] = make(map[string]string)
		}
		byLocale[e.Locale][e.Key] = e.Value
	}
	for locale, bundle := range byLocale {
		t.AddBundle(locale, bundle)
	}
	return nil
}

// LoadAll implements ports.Translator.
func (t *SimpleTranslator) LoadAll(ctx context.Context) ([]*ports.Translation, error) {
	rows, err := t.repo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to load translations: %w", err)
	}
	out := make([]*ports.Translation, len(rows))
	for i, r := range rows {
		out[i] = &ports.Translation{Key: r.Key, Locale: r.Locale, Value: r.Value}
	}
	return out, nil
}

// Set implements ports.Translator.
func (t *SimpleTranslator) Set(ctx context.Context, key, locale, value string) error {
	return t.repo.Upsert(ctx, key, locale, value)
}

var _ ports.Translator = (*SimpleTranslator)(nil)
