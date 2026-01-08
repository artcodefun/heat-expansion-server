package handlers

import (
	"net/http"
	"strings"

	"github.com/artcodefun/heat-expansion-api/internal/core/cqrs"
	"github.com/artcodefun/heat-expansion-api/internal/core/domain"
	"github.com/artcodefun/heat-expansion-api/internal/interfaces/http/dtos"
	"github.com/gin-gonic/gin"
)

type OperationHandler struct {
	queries  cqrs.OperationQueries
	commands cqrs.OperationCommands
}

func NewOperationHandler(queries cqrs.OperationQueries, commands cqrs.OperationCommands) *OperationHandler {
	return &OperationHandler{queries: queries, commands: commands}
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

// GetOperation handles GET /operations/:operationId.
func (h *OperationHandler) GetOperation(c *gin.Context) {
	var req dtos.OperationGetRequest
	if !bindRequest(c, &req) {
		return
	}
	ctx := queryCtx(c)
	op, err := h.queries.GetOperation(ctx, req.Uri.OperationID)
	if handleCoreErr(c, err) {
		return
	}
	c.JSON(http.StatusOK, dtos.OperationFromReadModel(op))
}

// ListByBase handles GET /operations/bases/:baseId.
func (h *OperationHandler) ListByBase(c *gin.Context) {
	var req dtos.OperationByBaseRequest
	if !bindRequest(c, &req) {
		return
	}
	ctx := queryCtx(c)
	ops, err := h.queries.ListOperationsByBase(ctx, req.Uri.BaseID)
	if handleCoreErr(c, err) {
		return
	}
	resp := make([]dtos.MilitaryOperationDTO, 0, len(ops))
	for _, item := range ops {
		resp = append(resp, dtos.OperationFromReadModel(item))
	}
	c.JSON(http.StatusOK, resp)
}

// ListActive handles GET /operations/bases/:baseId/active.
func (h *OperationHandler) ListActive(c *gin.Context) {
	var req dtos.OperationActiveByBaseRequest
	if !bindRequest(c, &req) {
		return
	}
	ctx := queryCtx(c)
	ops, err := h.queries.ListActiveOperations(ctx, req.Uri.BaseID)
	if handleCoreErr(c, err) {
		return
	}
	resp := make([]dtos.MilitaryOperationDTO, 0, len(ops))
	for _, item := range ops {
		resp = append(resp, dtos.OperationFromReadModel(item))
	}
	c.JSON(http.StatusOK, resp)
}

// SpeedUp handles POST /operations/:operationId/speed-up.
func (h *OperationHandler) SpeedUp(c *gin.Context) {
	var req dtos.OperationSpeedUpRequest
	if !bindRequest(c, &req) {
		return
	}
	ctx := commandCtx(c)
	if err := h.commands.SpeedUpOperationWithCrystals(ctx, req.Uri.OperationID); handleCoreErr(c, err) {
		return
	}
	c.Status(http.StatusOK)
}

// Create handles POST /operations.
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
	ctx := commandCtx(c)
	// Map DTOs into domain-level deployment requests by item ID.
	deployments := make([]domain.ArmyDeploymentRequest, 0, len(req.Body.Deployed))
	for _, d := range req.Body.Deployed {
		deployments = append(deployments, domain.ArmyDeploymentRequest{
			PresentItemID: d.PresentItemID.Uuid(),
			Count:         d.Count,
		})
	}
	op, err := h.commands.CreateMilitaryOperation(ctx, opType, req.Body.SourceBaseID, *req.Body.TargetX, *req.Body.TargetY, deployments)
	if handleCoreErr(c, err) {
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": op.ID})
}
