package handlers

import (
	"net/http"

	"github.com/artcodefun/heat-expansion-server/internal/billing/application"
	"github.com/artcodefun/heat-expansion-server/internal/billing/application/ports"
	"github.com/artcodefun/heat-expansion-server/internal/billing/interfaces/http/dtos"
	"github.com/gin-gonic/gin"
)

type PackageHandler struct {
	queries    application.PackageQueries
	translator ports.Translator
}

func NewPackageHandler(queries application.PackageQueries, translator ports.Translator) *PackageHandler {
	return &PackageHandler{queries: queries, translator: translator}
}

// ListPackages handles GET /packages.
func (h *PackageHandler) ListPackages(c *gin.Context) {
	pkgs, err := h.queries.ListPackages(c.Request.Context())
	if handleCoreErr(c, h.translator, err) {
		return
	}
	c.JSON(http.StatusOK, dtos.CrystalPackageResponsesFromReadModels(pkgs))
}
