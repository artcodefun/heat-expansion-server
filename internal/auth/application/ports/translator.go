package ports

// Translator defines the interface for localizing keys in Auth service.
type Translator interface {
	T(locale, key string, params map[string]any) string
}
