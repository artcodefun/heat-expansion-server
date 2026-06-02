package i18n

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
	"sync"
)

// Translator is a thread-safe, locale-aware string lookup with parameter substitution.
// Embed it in service-specific translator structs and call LoadFromJsonFiles to populate.
type Translator struct {
	mu      sync.RWMutex
	bundles map[string]map[string]string
}

// LoadFromJsonFiles walks files and loads every *.json file as a locale bundle.
// The locale is inferred from the filename stem (e.g. "en.json" → "en", "messages-ru.json" → "ru").
func (t *Translator) LoadFromJsonFiles(files fs.FS) error {
	return fs.WalkDir(files, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || filepath.Ext(path) != ".json" {
			return nil
		}

		data, err := fs.ReadFile(files, path)
		if err != nil {
			return fmt.Errorf("read i18n file %s: %w", path, err)
		}

		var bundle map[string]string
		if err := json.Unmarshal(data, &bundle); err != nil {
			return fmt.Errorf("parse i18n file %s: %w", path, err)
		}

		stem := strings.TrimSuffix(filepath.Base(path), ".json")
		parts := strings.Split(stem, "-")
		locale := parts[len(parts)-1]

		t.mu.Lock()
		for k, v := range bundle {
			t.setLocked(locale, k, v)
		}
		t.mu.Unlock()

		return nil
	})
}

// AddBundle merges entries into the given locale, overwriting existing keys.
// Useful for loading content translations from a database on top of systemic ones.
func (t *Translator) AddBundle(locale string, entries map[string]string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	for k, v := range entries {
		t.setLocked(locale, k, v)
	}
}

// T returns the localized string for key in locale, with params substituted.
// Falls back to English if the requested locale is missing a key, and returns
// the key itself if neither locale has it.
func (t *Translator) T(locale, key string, params map[string]any) string {
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

func (t *Translator) setLocked(locale, key, value string) {
	locale = strings.ToLower(locale)
	if t.bundles == nil {
		t.bundles = make(map[string]map[string]string)
	}
	if _, ok := t.bundles[locale]; !ok {
		t.bundles[locale] = make(map[string]string)
	}
	t.bundles[locale][key] = value
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
