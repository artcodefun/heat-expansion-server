package handlers

import (
	"net/http"

	"github.com/artcodefun/heat-expansion-server/internal/auth/application/cqrs"
	"github.com/artcodefun/heat-expansion-server/internal/auth/application/ports"
	"github.com/artcodefun/heat-expansion-server/internal/auth/interfaces/http/dtos"
	"github.com/gin-gonic/gin"
)

type AccountHandler struct {
	commands   cqrs.AccountCommands
	queries    cqrs.AccountQueries
	translator ports.Translator
}

func NewAccountHandler(commands cqrs.AccountCommands, queries cqrs.AccountQueries, translator ports.Translator) *AccountHandler {
	return &AccountHandler{commands: commands, queries: queries, translator: translator}
}

// Register handles POST /api/v1/register.
func (h *AccountHandler) Register(c *gin.Context) {
	var req dtos.AccountRegisterRequest
	if !bindRequest(c, &req) {
		return
	}
	actor := actor(c)
	if err := h.commands.RegisterAccount(c.Request.Context(), actor, req.Name, req.Email, req.Password); handleCoreErr(c, h.translator, err) {
		return
	}
	c.Status(http.StatusCreated)
}

// Login handles POST /api/v1/login.
func (h *AccountHandler) Login(c *gin.Context) {
	var req dtos.AccountLoginRequest
	if !bindRequest(c, &req) {
		return
	}
	actor := actor(c)
	token, err := h.commands.Login(c.Request.Context(), actor, req.Email, req.Password)
	if handleCoreErr(c, h.translator, err) {
		return
	}
	c.JSON(http.StatusOK, dtos.LoginResponse{Token: token})
}

// RequestPasswordReset handles POST /api/v1/password-reset/request.
func (h *AccountHandler) RequestPasswordReset(c *gin.Context) {
	var req dtos.PasswordResetRequest
	if !bindRequest(c, &req) {
		return
	}
	actor := actor(c)
	if err := h.commands.RequestPasswordReset(c.Request.Context(), actor, req.Email); handleCoreErr(c, h.translator, err) {
		return
	}
	c.Status(http.StatusNoContent)
}

// ResetPassword handles POST /api/v1/password-reset/confirm.
func (h *AccountHandler) ResetPassword(c *gin.Context) {
	var req dtos.PasswordResetConfirmRequest
	if !bindRequest(c, &req) {
		return
	}
	actor := actor(c)
	if err := h.commands.ResetPassword(c.Request.Context(), actor, req.Email, req.Token, req.NewPassword); handleCoreErr(c, h.translator, err) {
		return
	}
	c.Status(http.StatusNoContent)
}
