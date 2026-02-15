package ports

// Translator defines the interface for localizing keys.
type Translator interface {
	T(locale, key string, params map[string]any) string
}
