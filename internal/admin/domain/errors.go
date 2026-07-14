package domain

import "fmt"

type TranslationKey = string

// H is a convenience alias for a string-keyed params map.
type H = map[string]any

// Error is a translatable domain error carrying a translation key and optional params.
type Error struct {
	Key    TranslationKey
	Params H
}

func (e Error) Error() string {
	if len(e.Params) == 0 {
		return e.Key
	}
	return fmt.Sprintf("%s: %v", e.Key, e.Params)
}

// NewError creates a translatable domain error.
func NewError(key TranslationKey, params H) error {
	return Error{Key: key, Params: params}
}
