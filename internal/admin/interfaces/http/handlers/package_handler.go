package handlers

import (
	"net/http"

	"github.com/artcodefun/heat-expansion-server/internal/admin/application/cqrs"
	"github.com/artcodefun/heat-expansion-server/internal/admin/application/ports"
	"github.com/artcodefun/heat-expansion-server/internal/admin/interfaces/http/dtos"
	"github.com/gin-gonic/gin"
)

// PackageHandler handles CRUD endpoints for billing crystal packages.
type PackageHandler struct {
	commands   cqrs.PackageCommands
	queries    cqrs.PackageQueries
	translator ports.Translator
}

func NewPackageHandler(commands cqrs.PackageCommands, queries cqrs.PackageQueries, translator ports.Translator) *PackageHandler {
	return &PackageHandler{commands: commands, queries: queries, translator: translator}
}

// ListPackages handles GET /api/v1/billing/packages.
func (h *PackageHandler) ListPackages(c *gin.Context) {
	list, err := h.queries.ListCrystalPackages(c.Request.Context(), actor(c))
	if handleCoreErr(c, h.translator, err) {
		return
	}
	out := make([]dtos.CrystalPackageResponse, len(list))
	for i, p := range list {
		out[i] = dtos.CrystalPackageResponseFromModel(p)
	}
	c.JSON(http.StatusOK, out)
}

// GetPackage handles GET /api/v1/billing/packages/:id.
func (h *PackageHandler) GetPackage(c *gin.Context) {
	var uri dtos.GetPackageURI
	if !bindURI(c, &uri) {
		return
	}
	p, err := h.queries.GetCrystalPackage(c.Request.Context(), actor(c), uri.ID.Uuid())
	if handleCoreErr(c, h.translator, err) {
		return
	}
	c.JSON(http.StatusOK, dtos.CrystalPackageResponseFromModel(p))
}

// CreatePackage handles POST /api/v1/billing/packages.
func (h *PackageHandler) CreatePackage(c *gin.Context) {
	var req dtos.CreateCrystalPackageRequest
	if !bindRequest(c, &req) {
		return
	}
	p, err := h.commands.CreateCrystalPackage(c.Request.Context(), actor(c), req.ToModel())
	if handleCoreErr(c, h.translator, err) {
		return
	}
	c.JSON(http.StatusCreated, dtos.CrystalPackageResponseFromModel(p))
}

// UpdatePackage handles PUT /api/v1/billing/packages/:id.
func (h *PackageHandler) UpdatePackage(c *gin.Context) {
	var uri dtos.GetPackageURI
	if !bindURI(c, &uri) {
		return
	}
	var req dtos.UpdateCrystalPackageRequest
	if !bindRequest(c, &req) {
		return
	}
	p, err := h.commands.UpdateCrystalPackage(c.Request.Context(), actor(c), req.ToModel(uri.ID.Uuid()))
	if handleCoreErr(c, h.translator, err) {
		return
	}
	c.JSON(http.StatusOK, dtos.CrystalPackageResponseFromModel(p))
}
