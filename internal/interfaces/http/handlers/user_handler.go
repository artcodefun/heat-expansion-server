package handlers

import (
	"net/http"

	"github.com/artcodefun/heat-expansion-api/internal/core/cqrs"
	"github.com/artcodefun/heat-expansion-api/internal/interfaces/http/dtos"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	commands cqrs.UserCommands
}

func NewUserHandler(commands cqrs.UserCommands) *UserHandler {
	return &UserHandler{commands: commands}
}

func (h *UserHandler) Register(c *gin.Context) {
	var req dtos.UserRegisterRequest
	if !bindRequest(c, &req) {
		return
	}
	ctx := commandCtx(c)
	if err := h.commands.Create(ctx, req.Body.Name, req.Body.Email, req.Body.Password); handleCoreErr(c, err) {
		return
	}
	c.Status(http.StatusCreated)
}

func (h *UserHandler) Login(c *gin.Context) {
	var req dtos.UserLoginRequest
	if !bindRequest(c, &req) {
		return
	}
	if req.Body.Email == "" || req.Body.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email and password are required"})
		return
	}
	ctx := commandCtx(c)
	token, err := h.commands.Authenticate(ctx, req.Body.Email, req.Body.Password)
	if handleCoreErr(c, err) {
		return
	}
	c.JSON(http.StatusOK, dtos.LoginResponse{Token: token})
}
