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

func (h *AccountHandler) Register(c *gin.Context) {
	var req dtos.AccountRegisterRequest
	if !bindRequest(c, &req) {
		return
	}
	ctx := commandCtx(c)
	if err := h.commands.RegisterAccount(ctx, req.Name, req.Email, req.Password); handleCoreErr(c, h.translator, err) {
		return
	}
	c.Status(http.StatusCreated)
}

func (h *AccountHandler) Login(c *gin.Context) {
	var req dtos.AccountLoginRequest
	if !bindRequest(c, &req) {
		return
	}
	ctx := commandCtx(c)
	token, err := h.commands.Login(ctx, req.Email, req.Password)
	if handleCoreErr(c, h.translator, err) {
		return
	}
	c.JSON(http.StatusOK, dtos.LoginResponse{Token: token})
}
