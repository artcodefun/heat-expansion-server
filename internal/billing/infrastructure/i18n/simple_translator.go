package i18n

import (
	"github.com/artcodefun/heat-expansion-server/internal/billing/application/ports"
	"github.com/artcodefun/heat-expansion-server/internal/billing/infrastructure/i18n/locales"
	platformi18n "github.com/artcodefun/heat-expansion-server/internal/platform/i18n"
)

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
