package ports

// Translator resolves a translation key + params into a localised string.
type Translator interface {
	T(locale, key string, params map[string]any) string
}
