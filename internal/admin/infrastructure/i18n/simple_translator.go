package i18n

import (
	"github.com/artcodefun/heat-expansion-server/internal/admin/application/ports"
	"github.com/artcodefun/heat-expansion-server/internal/admin/infrastructure/i18n/locales"
	platformi18n "github.com/artcodefun/heat-expansion-server/internal/platform/i18n"
)

// SimpleTranslator is an embed-only translator for the admin service.
// Admin strings are all systemic (errors, status messages) so no DB-backed
// content loading is needed.
type SimpleTranslator struct {
	platformi18n.Translator
}

func NewSimpleTranslator() (*SimpleTranslator, error) {
	t := &SimpleTranslator{}
	if err := t.LoadFromJsonFiles(locales.Files); err != nil {
		return nil, err
	}
	return t, nil
}

var _ ports.Translator = (*SimpleTranslator)(nil)
