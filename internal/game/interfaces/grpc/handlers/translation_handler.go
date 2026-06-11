package handlers

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	gamev1 "github.com/artcodefun/heat-expansion-server/contracts/game/grpc/v1"
	"github.com/artcodefun/heat-expansion-server/internal/game/application/ports"
	"github.com/artcodefun/heat-expansion-server/internal/game/interfaces/grpc/dtos"
)

// TranslationHandler implements gamev1.TranslationServiceServer.
type TranslationHandler struct {
	gamev1.UnimplementedTranslationServiceServer
	translator ports.Translator
}

func NewTranslationHandler(translator ports.Translator) *TranslationHandler {
	return &TranslationHandler{translator: translator}
}

func (h *TranslationHandler) UpsertTranslation(ctx context.Context, req *gamev1.UpsertTranslationRequest) (*gamev1.UpsertTranslationResponse, error) {
	e := req.GetEntry()
	if e == nil {
		return nil, status.Error(codes.InvalidArgument, "entry is required")
	}
	if e.Key == "" {
		return nil, status.Error(codes.InvalidArgument, "entry.key is required")
	}
	if err := dtos.ValidateLocale(e.Locale); err != nil {
		return nil, err
	}

	if err := h.translator.Set(ctx, e.Key, e.Locale, e.Value); err != nil {
		return nil, dtos.StatusFromError(ctx, h.translator, err)
	}
	return &gamev1.UpsertTranslationResponse{
		Entry: &gamev1.TranslationEntry{Key: e.Key, Locale: e.Locale, Value: e.Value},
	}, nil
}

func (h *TranslationHandler) ListTranslations(ctx context.Context, _ *gamev1.ListTranslationsRequest) (*gamev1.ListTranslationsResponse, error) {
	entries, err := h.translator.LoadAll(ctx)
	if err != nil {
		return nil, dtos.StatusFromError(ctx, h.translator, err)
	}
	return &gamev1.ListTranslationsResponse{Entries: dtos.TranslationEntriesToProto(entries)}, nil
}
