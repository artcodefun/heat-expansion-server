package handlers

import (
	"net/http"
	"strings"

	"github.com/artcodefun/heat-expansion-server/internal/admin/application/cqrs"
	"github.com/artcodefun/heat-expansion-server/internal/admin/application/ports"
	"github.com/artcodefun/heat-expansion-server/internal/admin/interfaces/http/dtos"
	"github.com/gin-gonic/gin"
)

// AdminHandler handles admin authentication endpoints.
type AdminHandler struct {
	commands   cqrs.AdminCommands
	queries    cqrs.AdminQueries
	translator ports.Translator
}

func NewAdminHandler(commands cqrs.AdminCommands, queries cqrs.AdminQueries, translator ports.Translator) *AdminHandler {
	return &AdminHandler{commands: commands, queries: queries, translator: translator}
}

// Register handles POST /api/v1/auth/register.
func (h *AdminHandler) Register(c *gin.Context) {
	var req dtos.RegisterRequest
	if !bindRequest(c, &req) {
		return
	}
	token, err := h.commands.Register(c.Request.Context(), actor(c), req.Username, req.InviteToken, req.Password)
	if handleCoreErr(c, h.translator, err) {
		return
	}
	c.JSON(http.StatusOK, dtos.SessionResponse{Token: token})
}

// Login handles POST /api/v1/auth/login.
func (h *AdminHandler) Login(c *gin.Context) {
	var req dtos.LoginRequest
	if !bindRequest(c, &req) {
		return
	}
	token, err := h.commands.Login(c.Request.Context(), actor(c), req.Username, req.Password)
	if handleCoreErr(c, h.translator, err) {
		return
	}
	c.JSON(http.StatusOK, dtos.SessionResponse{Token: token})
}

// Logout handles POST /api/v1/auth/logout.
func (h *AdminHandler) Logout(c *gin.Context) {
	auth := c.GetHeader("Authorization")
	token := strings.TrimPrefix(auth, "Bearer ")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing bearer token"})
		return
	}
	if handleCoreErr(c, h.translator, h.commands.Logout(c.Request.Context(), actor(c), token)) {
		return
	}
	c.Status(http.StatusNoContent)
}

// Me handles GET /api/v1/auth/me.
func (h *AdminHandler) Me(c *gin.Context) {
	a := actor(c)
	profile, err := h.queries.GetProfile(c.Request.Context(), a, a.AdminID)
	if handleCoreErr(c, h.translator, err) {
		return
	}
	c.JSON(http.StatusOK, dtos.ProfileResponseFromModel(profile))
}
