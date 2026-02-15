package domain

import (
	"fmt"
)

// TranslationKey is a type alias for string used for translatable items.
type TranslationKey = string

// Error represents a translatable domain error.
type Error struct {
	Key    TranslationKey
	Params H
}

func (e Error) Error() string {
	// Fallback description for logging/debugging if not translated
	if len(e.Params) == 0 {
		return e.Key
	}
	return fmt.Sprintf("%s: %v", e.Key, e.Params)
}

// NewError creates a new domain error with key and params.
func NewError(key TranslationKey, params H) error {
	return Error{
		Key:    key,
		Params: params,
	}
}

// Convenience alias for string Hashmap
type H = map[string]any
