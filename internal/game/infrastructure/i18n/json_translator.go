package i18n

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/artcodefun/heat-expansion-server/internal/game/application/ports"
)

type JSONTranslator struct {
	bundles map[string]map[string]string
	mu      sync.RWMutex
}

func NewJSONTranslator() *JSONTranslator {
	return &JSONTranslator{
		bundles: make(map[string]map[string]string),
	}
}

// LoadFromFS loads translations from a filesystem (primarily for embedded files).
func (t *JSONTranslator) LoadFromFS(f fs.FS, dirPath string) error {
	return fs.WalkDir(f, dirPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || filepath.Ext(path) != ".json" {
			return nil
		}

		data, err := fs.ReadFile(f, path)
		if err != nil {
			return fmt.Errorf("failed to read i18n file %s: %w", path, err)
		}

		return t.loadData(path, data)
	})
}

// LoadFromDir loads translations from a standard directory.
func (t *JSONTranslator) LoadFromDir(dirPath string) error {
	files, err := os.ReadDir(dirPath)
	if err != nil {
		return fmt.Errorf("failed to read i18n directory: %w", err)
	}

	for _, file := range files {
		if file.IsDir() || filepath.Ext(file.Name()) != ".json" {
			continue
		}

		filePath := filepath.Join(dirPath, file.Name())
		data, err := os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("failed to read i18n file %s: %w", filePath, err)
		}

		if err := t.loadData(file.Name(), data); err != nil {
			return err
		}
	}

	return nil
}

func (t *JSONTranslator) loadData(fileName string, data []byte) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	var bundle map[string]string
	if err := json.Unmarshal(data, &bundle); err != nil {
		return fmt.Errorf("failed to parse i18n data from %s: %w", fileName, err)
	}

	localePart := strings.TrimSuffix(filepath.Base(fileName), ".json")
	parts := strings.Split(localePart, "-")
	locale := strings.ToLower(parts[len(parts)-1])

	if _, ok := t.bundles[locale]; !ok {
		t.bundles[locale] = make(map[string]string)
	}

	for k, v := range bundle {
		t.bundles[locale][k] = v
	}

	return nil
}

func (t *JSONTranslator) T(locale, key string, params map[string]any) string {
	t.mu.RLock()
	defer t.mu.RUnlock()

	locale = strings.ToLower(locale)
	bundle, ok := t.bundles[locale]
	if !ok {
		bundle, ok = t.bundles["en"]
		if !ok {
			return gogoTemplate(key, params)
		}
	}

	msg, ok := bundle[key]
	if !ok {
		if locale != "en" {
			if enBundle, ok := t.bundles["en"]; ok {
				if enMsg, ok := enBundle[key]; ok {
					return gogoTemplate(enMsg, params)
				}
			}
		}
		return gogoTemplate(key, params)
	}

	return gogoTemplate(msg, params)
}

func gogoTemplate(tpl string, params map[string]any) string {
	if len(params) == 0 {
		return tpl
	}
	result := tpl
	for k, v := range params {
		placeholder := fmt.Sprintf("{%s}", k)
		result = strings.ReplaceAll(result, placeholder, fmt.Sprintf("%v", v))
	}
	return result
}

var _ ports.Translator = (*JSONTranslator)(nil)
