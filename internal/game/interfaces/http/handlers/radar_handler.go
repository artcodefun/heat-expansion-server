package handlers

import (
	"net/http"

	"github.com/artcodefun/heat-expansion-server/internal/game/application/cqrs"
	"github.com/artcodefun/heat-expansion-server/internal/game/application/ports"
	"github.com/artcodefun/heat-expansion-server/internal/game/interfaces/http/dtos"
	"github.com/gin-gonic/gin"
)

type RadarHandler struct {
	queries    cqrs.RadarQueries
	translator ports.Translator
}

func NewRadarHandler(queries cqrs.RadarQueries, translator ports.Translator) *RadarHandler {
	return &RadarHandler{
		queries:    queries,
		translator: translator,
	}
}

// ListIncomingThreats handles GET /bases/:baseId/threats.
func (h *RadarHandler) ListIncomingThreats(c *gin.Context) {
	var req dtos.RadarThreatsListRequest
	if !bindRequest(c, &req) {
		return
	}
	actor := actor(c)
	threats, err := h.queries.ListIncomingThreats(c.Request.Context(), actor, req.Uri.BaseID)
	if handleCoreErr(c, h.translator, err) {
		return
	}
	c.JSON(http.StatusOK, dtos.RadarThreatsFromReadModels(threats, h.translator, getLocale(c)))
}
