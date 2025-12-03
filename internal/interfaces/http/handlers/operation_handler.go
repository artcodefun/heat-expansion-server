package handlers

import (
	"net/http"
	"strings"

	"github.com/artcodefun/heat-expansion-api/internal/core/cqrs"
	"github.com/artcodefun/heat-expansion-api/internal/core/domain"
	"github.com/artcodefun/heat-expansion-api/internal/interfaces/http/dtos"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

func (h *OperationHandler) GetOperation(c *gin.Context) {
	var uri struct {
		OperationID int `uri:"operationId" binding:"required,min=1"`
	}
	if err := c.ShouldBindUri(&uri); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid operationId"})
		return
	}
	ctx := queryCtx(c)
	op, err := h.queries.GetOperation(ctx, uri.OperationID)
	if handleCQRS(c, err) {
		return
	}
	c.JSON(http.StatusOK, dtos.OperationFromReadModel(op))
}

func (h *OperationHandler) ListByBase(c *gin.Context) {
	var uri dtos.SectorBaseURI
	if err := c.ShouldBindUri(&uri); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid baseId"})
		return
	}
	ctx := queryCtx(c)
	ops, err := h.queries.ListOperationsByBase(ctx, uri.BaseID)
	if handleCQRS(c, err) {
		return
	}
	resp := make([]dtos.MilitaryOperationDTO, 0, len(ops))
	for _, item := range ops {
		resp = append(resp, dtos.OperationFromReadModel(item))
	}
	c.JSON(http.StatusOK, resp)
}

func (h *OperationHandler) ListActive(c *gin.Context) {
	var uri dtos.SectorBaseURI
	if err := c.ShouldBindUri(&uri); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid baseId"})
		return
	}
	ctx := queryCtx(c)
	ops, err := h.queries.ListActiveOperations(ctx, uri.BaseID)
	if handleCQRS(c, err) {
		return
	}
	resp := make([]dtos.MilitaryOperationDTO, 0, len(ops))
	for _, item := range ops {
		resp = append(resp, dtos.OperationFromReadModel(item))
	}
	c.JSON(http.StatusOK, resp)
}

func (h *OperationHandler) Create(c *gin.Context) {
	var body dtos.CreateOperationRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}
	opType, ok := parseOperationType(string(body.Type))
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid type"})
		return
	}
	ctx := commandCtx(c)
	// Map DTOs into domain-level deployment requests by item ID.
	deployments := make([]domain.ArmyDeploymentRequest, 0, len(body.Deployed))
	for _, d := range body.Deployed {
		itemID, err := uuid.Parse(d.PresentItemID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid itemId"})
			return
		}
		deployments = append(deployments, domain.ArmyDeploymentRequest{
			PresentItemID: itemID,
			Count:         d.Count,
		})
	}
	op, err := h.commands.CreateMilitaryOperation(ctx, opType, body.SourceBaseID, body.TargetX, body.TargetY, deployments)
	if handleCQRS(c, err) {
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": op.ID})
}
