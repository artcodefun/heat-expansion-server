package handlers

import (
	"net/http"
	"strings"

	"github.com/artcodefun/heat-expansion-server/internal/game/application/cqrs"
	"github.com/artcodefun/heat-expansion-server/internal/game/application/ports"
	"github.com/artcodefun/heat-expansion-server/internal/game/domain"
	"github.com/artcodefun/heat-expansion-server/internal/game/interfaces/http/dtos"
	"github.com/gin-gonic/gin"
)

type OperationHandler struct {
	queries    cqrs.OperationQueries
	commands   cqrs.OperationCommands
	translator ports.Translator
}

func NewOperationHandler(queries cqrs.OperationQueries, commands cqrs.OperationCommands, translator ports.Translator) *OperationHandler {
	return &OperationHandler{
		queries:    queries,
		commands:   commands,
		translator: translator,
	}
}

func parseOperationType(raw string) (domain.MilitaryOperationType, bool) {
	switch strings.ToUpper(raw) {
	case string(domain.MilitaryOperationTypeAttack):
		return domain.MilitaryOperationTypeAttack, true
	case string(domain.MilitaryOperationTypeSpy):
		return domain.MilitaryOperationTypeSpy, true
	case string(domain.MilitaryOperationTypeOccupation):
		return domain.MilitaryOperationTypeOccupation, true
	default:
		return "", false
	}
}

// GetOperation handles GET /bases/:baseId/operations/:operationId.
func (h *OperationHandler) GetOperation(c *gin.Context) {
	var req dtos.OperationGetRequest
	if !bindRequest(c, &req) {
		return
	}
	actor := actor(c)
	op, err := h.queries.GetOperation(c.Request.Context(), actor, req.Uri.OperationID)
	if handleCoreErr(c, h.translator, err) {
		return
	}
	c.JSON(http.StatusOK, dtos.OperationFromReadModel(op, h.translator, getLocale(c)))
}

// ListByBase handles GET /bases/:baseId/operations.
func (h *OperationHandler) ListByBase(c *gin.Context) {
	var req dtos.OperationByBaseRequest
	if !bindRequest(c, &req) {
		return
	}
	actor := actor(c)
	ops, err := h.queries.ListOperationsByBase(c.Request.Context(), actor, req.Uri.BaseID)
	if handleCoreErr(c, h.translator, err) {
		return
	}
	c.JSON(http.StatusOK, dtos.OperationsFromReadModels(ops, h.translator, getLocale(c)))
}

// ListActive handles GET /bases/:baseId/operations/active.
func (h *OperationHandler) ListActive(c *gin.Context) {
	var req dtos.OperationActiveByBaseRequest
	if !bindRequest(c, &req) {
		return
	}
	actor := actor(c)
	ops, err := h.queries.ListActiveOperations(c.Request.Context(), actor, req.Uri.BaseID)
	if handleCoreErr(c, h.translator, err) {
		return
	}
	c.JSON(http.StatusOK, dtos.OperationsFromReadModels(ops, h.translator, getLocale(c)))
}

// SpeedUp handles POST /bases/:baseId/operations/:operationId/speed-up.
func (h *OperationHandler) SpeedUp(c *gin.Context) {
	var req dtos.OperationSpeedUpRequest
	if !bindRequest(c, &req) {
		return
	}
	actor := actor(c)
	if err := h.commands.SpeedUpOperationWithCrystals(c.Request.Context(), actor, req.Uri.OperationID); handleCoreErr(c, h.translator, err) {
		return
	}
	c.Status(http.StatusOK)
}

// Cancel handles POST /bases/:baseId/operations/:operationId/cancel.
func (h *OperationHandler) Cancel(c *gin.Context) {
	var req dtos.OperationCancelRequest
	if !bindRequest(c, &req) {
		return
	}
	actor := actor(c)
	if err := h.commands.CancelMilitaryOperation(c.Request.Context(), actor, req.Uri.OperationID); handleCoreErr(c, h.translator, err) {
		return
	}
	c.Status(http.StatusOK)
}

// Create handles POST /bases/:baseId/operations.
func (h *OperationHandler) Create(c *gin.Context) {
	var req dtos.OperationCreateRequest
	if !bindRequest(c, &req) {
		return
	}
	opType, ok := parseOperationType(string(req.Body.Type))
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid type"})
		return
	}
	actor := actor(c)
	// Map DTOs into domain-level deployment requests by item ID.
	deployments := make([]domain.ArmyDeploymentRequest, 0, len(req.Body.Deployed))
	for _, d := range req.Body.Deployed {
		deployments = append(deployments, domain.ArmyDeploymentRequest{
			PresentItemID: d.PresentItemID.Uuid(),
			Count:         d.Count,
		})
	}
	op, err := h.commands.CreateMilitaryOperation(c.Request.Context(), actor, opType, req.Uri.BaseID, *req.Body.TargetX, *req.Body.TargetY, deployments)
	if handleCoreErr(c, h.translator, err) {
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": op.ID})
}
