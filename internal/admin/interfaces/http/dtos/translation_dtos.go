package dtos

import "github.com/artcodefun/heat-expansion-server/internal/admin/application/ports"

// UpsertTranslationRequest is the body for PUT /api/v1/game/translations.
type UpsertTranslationRequest struct {
	Key    string `json:"key"    binding:"required"`
	Locale string `json:"locale" binding:"required"`
	Value  string `json:"value"`
}

// TranslationResponse is the JSON representation of a single translation entry.
type TranslationResponse struct {
	Key    string `json:"key"`
	Locale string `json:"locale"`
	Value  string `json:"value"`
}

func TranslationResponseFromModel(m *ports.Translation) TranslationResponse {
	return TranslationResponse{Key: m.Key, Locale: m.Locale, Value: m.Value}
}
