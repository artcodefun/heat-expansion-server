package ports

type Translator interface {
	T(locale, key string, params map[string]any) string
}
