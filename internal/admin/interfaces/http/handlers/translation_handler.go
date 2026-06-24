package handlers

import (
	"net/http"

	"github.com/artcodefun/heat-expansion-server/internal/admin/application/cqrs"
	"github.com/artcodefun/heat-expansion-server/internal/admin/application/ports"
	"github.com/artcodefun/heat-expansion-server/internal/admin/interfaces/http/dtos"
	"github.com/gin-gonic/gin"
)

// TranslationHandler handles CRUD endpoints for game translations.
type TranslationHandler struct {
	commands   cqrs.TranslationCommands
	queries    cqrs.TranslationQueries
	translator ports.Translator
}

func NewTranslationHandler(commands cqrs.TranslationCommands, queries cqrs.TranslationQueries, translator ports.Translator) *TranslationHandler {
	return &TranslationHandler{commands: commands, queries: queries, translator: translator}
}

// ListTranslations handles GET /api/v1/game/translations.
func (h *TranslationHandler) ListTranslations(c *gin.Context) {
	list, err := h.queries.ListTranslations(c.Request.Context(), actor(c))
	if handleCoreErr(c, h.translator, err) {
		return
	}
	out := make([]dtos.TranslationResponse, len(list))
	for i, t := range list {
		out[i] = dtos.TranslationResponseFromModel(t)
	}
	c.JSON(http.StatusOK, out)
}

// UpsertTranslation handles PUT /api/v1/game/translations.
func (h *TranslationHandler) UpsertTranslation(c *gin.Context) {
	var req dtos.UpsertTranslationRequest
	if !bindRequest(c, &req) {
		return
	}
	t, err := h.commands.UpsertTranslation(c.Request.Context(), actor(c), req.Locale, req.Key, req.Value)
	if handleCoreErr(c, h.translator, err) {
		return
	}
	c.JSON(http.StatusOK, dtos.TranslationResponseFromModel(t))
}
