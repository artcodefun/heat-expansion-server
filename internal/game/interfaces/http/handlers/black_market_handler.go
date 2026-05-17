package handlers

import (
	"net/http"

	"github.com/artcodefun/heat-expansion-server/internal/game/application/cqrs"
	"github.com/artcodefun/heat-expansion-server/internal/game/application/ports"
	"github.com/artcodefun/heat-expansion-server/internal/game/interfaces/http/dtos"
	"github.com/gin-gonic/gin"
)

type BlackMarketHandler struct {
	translator ports.Translator
	commands   cqrs.BlackMarketCommands
	queries    cqrs.BlackMarketQueries
}

// NewBlackMarketHandler constructs a handler for Black Market endpoints.
func NewBlackMarketHandler(commands cqrs.BlackMarketCommands, queries cqrs.BlackMarketQueries, translator ports.Translator) *BlackMarketHandler {
	return &BlackMarketHandler{commands: commands, queries: queries, translator: translator}
}

// ListResourceRates handles GET /bases/:baseId/black-market/resources and returns the current Black Market exchange rates.
func (h *BlackMarketHandler) ListResourceRates(c *gin.Context) {
	var req dtos.BlackMarketResourceRatesRequest
	if !bindRequest(c, &req) {
		return
	}

	items, err := h.queries.ListResourceRates(c.Request.Context(), actor(c), req.Uri.BaseID)
	if handleCoreErr(c, h.translator, err) {
		return
	}

	c.JSON(http.StatusOK, dtos.BlackMarketResourceRateDTOListFromReadModels(items))
}

// PurchaseResources handles POST /bases/:baseId/black-market/resources/purchase and spends crystals for resources.
func (h *BlackMarketHandler) PurchaseResources(c *gin.Context) {
	var req dtos.BlackMarketResourcesPurchaseRequest
	if !bindRequest(c, &req) {
		return
	}

	err := h.commands.PurchaseResources(
		c.Request.Context(),
		actor(c),
		req.Uri.BaseID,
		dtos.ResourceTypeFromDTO(req.Body.ResourceType),
		req.Body.Crystals,
	)
	if handleCoreErr(c, h.translator, err) {
		return
	}

	c.Status(http.StatusOK)
}
