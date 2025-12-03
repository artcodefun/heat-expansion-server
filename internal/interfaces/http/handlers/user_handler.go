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
	var body dtos.RegisterRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}
	ctx := commandCtx(c)
	if err := h.commands.Create(ctx, body.Name, body.Email, body.Password); handleCQRS(c, err) {
		return
	}
	c.Status(http.StatusCreated)
}

func (h *UserHandler) Login(c *gin.Context) {
	var body dtos.LoginRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}
	if body.Email == "" || body.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email and password are required"})
		return
	}
	ctx := commandCtx(c)
	token, err := h.commands.Authenticate(ctx, body.Email, body.Password)
	if handleCQRS(c, err) {
		return
	}
	c.JSON(http.StatusOK, dtos.LoginResponse{Token: token})
}
