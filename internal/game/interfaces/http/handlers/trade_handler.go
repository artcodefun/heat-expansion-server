package handlers

import (
	"net/http"

	"github.com/artcodefun/heat-expansion-server/internal/game/application/cqrs"
	"github.com/artcodefun/heat-expansion-server/internal/game/application/ports"
	"github.com/artcodefun/heat-expansion-server/internal/game/domain"
	"github.com/artcodefun/heat-expansion-server/internal/game/interfaces/http/dtos"
	"github.com/gin-gonic/gin"
)

type TradeHandler struct {
	commands   cqrs.TradeCommands
	queries    cqrs.TradeQueries
	translator ports.Translator
}

func NewTradeHandler(commands cqrs.TradeCommands, queries cqrs.TradeQueries, translator ports.Translator) *TradeHandler {
	return &TradeHandler{commands: commands, queries: queries, translator: translator}
}

// GetInfo handles GET /trade/bases/:baseId/info.
func (h *TradeHandler) GetInfo(c *gin.Context) {
	var req dtos.TradeInfoRequest
	if !bindRequest(c, &req) {
		return
	}

	info, err := h.queries.GetTradeInfo(c.Request.Context(), actor(c), *req.Query.TargetX, *req.Query.TargetY)
	if handleCoreErr(c, h.translator, err) {
		return
	}

	c.JSON(http.StatusOK, dtos.TradeInfoFromReadModel(info, h.translator, getLocale(c)))
}

// Create handles POST /trade/bases/:baseId/operations.
func (h *TradeHandler) Create(c *gin.Context) {
	var req dtos.TradeOperationCreateRequest
	if !bindRequest(c, &req) {
		return
	}

	created, err := h.commands.CreateTradeOperation(
		c.Request.Context(),
		actor(c),
		req.Uri.BaseID,
		*req.Body.TargetX,
		*req.Body.TargetY,
		domain.PriceModel(dtos.PriceModelFromDTO(req.Body.OfferedResources)),
		dtos.ArmyDeploymentRequestsFromDTOs(req.Body.OfferedArmy),
		dtos.UUIDs(req.Body.OfferedStorageItemIDs),
		domain.PriceModel(dtos.PriceModelFromDTO(req.Body.RequestedResources)),
		dtos.ArmyDeploymentRequestsFromDTOs(req.Body.RequestedArmy),
		dtos.UUIDs(req.Body.RequestedStorageItemIDs),
		dtos.ArmyDeploymentRequestsFromDTOs(req.Body.TransportUnits),
	)
	if handleCoreErr(c, h.translator, err) {
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": created.ID})
}

// Accept handles POST /trade/bases/:baseId/operations/:operationId/accept.
func (h *TradeHandler) Accept(c *gin.Context) {
	var req dtos.TradeBaseOperationIDRequest
	if !bindRequest(c, &req) {
		return
	}

	if err := h.commands.AcceptTradeOperation(c.Request.Context(), actor(c), req.Uri.OperationID); handleCoreErr(c, h.translator, err) {
		return
	}
	c.Status(http.StatusOK)
}

// Decline handles POST /trade/bases/:baseId/operations/:operationId/decline.
func (h *TradeHandler) Decline(c *gin.Context) {
	var req dtos.TradeBaseOperationIDRequest
	if !bindRequest(c, &req) {
		return
	}

	if err := h.commands.DeclineTradeOperation(c.Request.Context(), actor(c), req.Uri.OperationID); handleCoreErr(c, h.translator, err) {
		return
	}
	c.Status(http.StatusOK)
}

// Cancel handles POST /trade/bases/:baseId/operations/:operationId/cancel.
func (h *TradeHandler) Cancel(c *gin.Context) {
	var req dtos.TradeBaseOperationIDRequest
	if !bindRequest(c, &req) {
		return
	}

	if err := h.commands.CancelTradeOperationByInitiator(c.Request.Context(), actor(c), req.Uri.OperationID); handleCoreErr(c, h.translator, err) {
		return
	}
	c.Status(http.StatusOK)
}

// SpeedUp handles POST /trade/bases/:baseId/operations/:operationId/speed-up.
func (h *TradeHandler) SpeedUp(c *gin.Context) {
	var req dtos.TradeBaseOperationIDRequest
	if !bindRequest(c, &req) {
		return
	}

	if err := h.commands.SpeedUpTradeOperationWithCrystals(c.Request.Context(), actor(c), req.Uri.OperationID); handleCoreErr(c, h.translator, err) {
		return
	}
	c.Status(http.StatusOK)
}

// GetOperation handles GET /trade/bases/:baseId/operations/:operationId.
func (h *TradeHandler) GetOperation(c *gin.Context) {
	var req dtos.TradeBaseOperationIDRequest
	if !bindRequest(c, &req) {
		return
	}

	op, err := h.queries.GetTradeOperation(c.Request.Context(), actor(c), req.Uri.BaseID, req.Uri.OperationID)
	if handleCoreErr(c, h.translator, err) {
		return
	}

	c.JSON(http.StatusOK, dtos.TradeOperationFromReadModel(op, h.translator, getLocale(c)))
}

// ListActive handles GET /trade/bases/:baseId/operations.
func (h *TradeHandler) ListActive(c *gin.Context) {
	var req dtos.TradeBaseRequest
	if !bindRequest(c, &req) {
		return
	}

	items, err := h.queries.ListActiveTradeOperations(c.Request.Context(), actor(c), req.Uri.BaseID)
	if handleCoreErr(c, h.translator, err) {
		return
	}

	c.JSON(http.StatusOK, dtos.TradeOperationsFromReadModels(items, h.translator, getLocale(c)))
}
