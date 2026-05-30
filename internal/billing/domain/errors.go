package domain

import "fmt"

type TranslationKey = string
type H = map[string]any

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

func NewError(key TranslationKey, params H) error {
	return Error{Key: key, Params: params}
}
