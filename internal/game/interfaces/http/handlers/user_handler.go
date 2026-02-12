package handlers

import (
	"net/http"

	"github.com/artcodefun/heat-expansion-server/internal/game/application/cqrs"
	"github.com/artcodefun/heat-expansion-server/internal/game/interfaces/http/dtos"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	commands cqrs.UserCommands
	queries  cqrs.UserQueries
}

func NewUserHandler(commands cqrs.UserCommands, queries cqrs.UserQueries) *UserHandler {
	return &UserHandler{commands: commands, queries: queries}
}

// Register handles POST /auth/register.
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

// Login handles POST /auth/login.
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

// GetCrystalBalance handles GET /user/balance and returns the authenticated user's crystal balance.
func (h *UserHandler) GetCrystalBalance(c *gin.Context) {
	ctx := queryCtx(c)
	if ctx.UserID == 0 {
		// Should not normally happen with Auth middleware, but guard just in case.
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	user, err := h.queries.GetUserProfile(ctx, ctx.UserID)
	if handleCoreErr(c, err) {
		return
	}

	c.JSON(http.StatusOK, dtos.UserCrystalBalanceResponse{Crystals: user.Crystals})
}
