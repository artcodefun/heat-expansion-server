package handlers

import (
	"net/http"

	"github.com/artcodefun/heat-expansion-server/internal/game/application/cqrs"
	"github.com/artcodefun/heat-expansion-server/internal/game/interfaces/http/dtos"
	"github.com/gin-gonic/gin"
)

type RadarHandler struct {
	queries cqrs.RadarQueries
}

func NewRadarHandler(queries cqrs.RadarQueries) *RadarHandler {
	return &RadarHandler{queries: queries}
}

// ListIncomingThreats handles GET /bases/:baseId/threats.
func (h *RadarHandler) ListIncomingThreats(c *gin.Context) {
	var req dtos.RadarThreatsListRequest
	if !bindRequest(c, &req) {
		return
	}
	ctx := queryCtx(c)
	threats, err := h.queries.ListIncomingThreats(ctx, req.Uri.BaseID)
	if handleCoreErr(c, err) {
		return
	}
	c.JSON(http.StatusOK, dtos.RadarThreatsFromReadModels(threats))
}
