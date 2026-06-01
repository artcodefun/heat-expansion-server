package i18n

import (
	"context"
	"encoding/json"
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
	"sync"

	"github.com/artcodefun/heat-expansion-server/internal/game/application/ports"
	repo "github.com/artcodefun/heat-expansion-server/internal/game/infrastructure/db/repo"
	"github.com/artcodefun/heat-expansion-server/internal/game/infrastructure/i18n/locales"
)

type SimpleTranslator struct {
	bundles map[string]map[string]string
	mu      sync.RWMutex
	repo    *repo.TranslationRepo
}

func NewSimpleTranslator(r *repo.TranslationRepo) (*SimpleTranslator, error) {
	t := &SimpleTranslator{
		bundles: make(map[string]map[string]string),
		repo:    r,
	}
	if err := t.loadSystem(); err != nil {
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
	t.mu.Lock()
	defer t.mu.Unlock()
	for _, e := range entries {
		t.setLocked(e.Locale, e.Key, e.Value)
	}
	return nil
}

func (t *SimpleTranslator) loadSystem() error {
	return fs.WalkDir(locales.Files, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || filepath.Ext(path) != ".json" {
			return nil
		}

		data, err := fs.ReadFile(locales.Files, path)
		if err != nil {
			return fmt.Errorf("failed to read i18n file %s: %w", path, err)
		}

		var bundle map[string]string
		if err := json.Unmarshal(data, &bundle); err != nil {
			return fmt.Errorf("failed to parse i18n file %s: %w", path, err)
		}

		localePart := strings.TrimSuffix(filepath.Base(path), ".json")
		parts := strings.Split(localePart, "-")
		locale := parts[len(parts)-1]

		t.mu.Lock()
		for k, v := range bundle {
			t.setLocked(locale, k, v)
		}
		t.mu.Unlock()

		return nil
	})
}

func (t *SimpleTranslator) setLocked(locale, key, value string) {
	locale = strings.ToLower(locale)
	if _, ok := t.bundles[locale]; !ok {
		t.bundles[locale] = make(map[string]string)
	}
	t.bundles[locale][key] = value
}

func (t *SimpleTranslator) T(locale, key string, params map[string]any) string {
	t.mu.RLock()
	defer t.mu.RUnlock()

	locale = strings.ToLower(locale)
	bundle, ok := t.bundles[locale]
	if !ok {
		bundle, ok = t.bundles["en"]
		if !ok {
			return applyParams(key, params)
		}
	}

	msg, ok := bundle[key]
	if !ok {
		if locale != "en" {
			if enBundle, ok := t.bundles["en"]; ok {
				if enMsg, ok := enBundle[key]; ok {
					return applyParams(enMsg, params)
				}
			}
		}
		return applyParams(key, params)
	}

	return applyParams(msg, params)
}

func applyParams(tpl string, params map[string]any) string {
	if len(params) == 0 {
		return tpl
	}
	result := tpl
	for k, v := range params {
		result = strings.ReplaceAll(result, fmt.Sprintf("{%s}", k), fmt.Sprintf("%v", v))
	}
	return result
}

var _ ports.Translator = (*SimpleTranslator)(nil)
