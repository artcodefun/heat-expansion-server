package i18n

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
	"sync"

	"github.com/artcodefun/heat-expansion-server/internal/auth/application/ports"
	"github.com/artcodefun/heat-expansion-server/internal/auth/infrastructure/i18n/locales"
)

type SimpleTranslator struct {
	bundles map[string]map[string]string
	mu      sync.RWMutex
}

func NewSimpleTranslator() (*SimpleTranslator, error) {
	t := &SimpleTranslator{
		bundles: make(map[string]map[string]string),
	}
	if err := t.loadSystem(); err != nil {
		return nil, err
	}
	return t, nil
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
		locale := strings.ToLower(parts[len(parts)-1])

		t.mu.Lock()
		if _, ok := t.bundles[locale]; !ok {
			t.bundles[locale] = make(map[string]string)
		}
		for k, v := range bundle {
			t.bundles[locale][k] = v
		}
		t.mu.Unlock()

		return nil
	})
}

func (t *SimpleTranslator) T(locale, key string, params map[string]any) string {
	t.mu.RLock()
	defer t.mu.RUnlock()

	locale = strings.ToLower(locale)
	bundle, ok := t.bundles[locale]
	if !ok {
		bundle = t.bundles["en"]
	}

	val, ok := bundle[key]
	if !ok {
		return key
	}

	for k, v := range params {
		val = strings.ReplaceAll(val, fmt.Sprintf("{%s}", k), fmt.Sprintf("%v", v))
	}

	return val
}

var _ ports.Translator = (*SimpleTranslator)(nil)
