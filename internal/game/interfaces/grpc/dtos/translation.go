package dtos

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	gamev1 "github.com/artcodefun/heat-expansion-server/contracts/game/grpc/v1"
	"github.com/artcodefun/heat-expansion-server/internal/game/application/ports"
)

// supportedLocales is the allowlist validated at the gRPC boundary.
var supportedLocales = map[string]struct{}{
	"en": {},
	"ru": {},
}

// ValidateLocale returns an InvalidArgument error if locale is not supported.
func ValidateLocale(locale string) error {
	if _, ok := supportedLocales[locale]; !ok {
		return status.Errorf(codes.InvalidArgument, "unsupported locale %q: must be one of en, ru", locale)
	}
	return nil
}

// TranslationEntryToProto maps a ports.Translation to its wire shape.
func TranslationEntryToProto(t *ports.Translation) *gamev1.TranslationEntry {
	if t == nil {
		return nil
	}
	return &gamev1.TranslationEntry{Key: t.Key, Locale: t.Locale, Value: t.Value}
}

// TranslationEntriesToProto maps a slice of ports.Translation to wire shape.
func TranslationEntriesToProto(ts []*ports.Translation) []*gamev1.TranslationEntry {
	out := make([]*gamev1.TranslationEntry, len(ts))
	for i, t := range ts {
		out[i] = TranslationEntryToProto(t)
	}
	return out
}
