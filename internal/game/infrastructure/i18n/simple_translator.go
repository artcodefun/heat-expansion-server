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

// LoadFromRepo fetches all content translations from the repository and merges them
// into the bundle, taking precedence over systemic translations for the same key+locale.
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

var _ ports.Translator = (*SimpleTranslator)(nil)
