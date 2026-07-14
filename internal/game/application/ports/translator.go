package ports

import "context"

// Translation is a single localised string entry.
type Translation struct {
	Key    string
	Locale string
	Value  string
}

// Translator defines the interface for localizing keys and managing the
// content translation catalog.
type Translator interface {
	// T returns the localised string for key in locale, substituting params.
	T(locale, key string, params map[string]any) string

	// LoadAll returns all translation entries.
	LoadAll(ctx context.Context) ([]*Translation, error)

	// Set upserts a single translation entry.
	Set(ctx context.Context, key, locale, value string) error
}
