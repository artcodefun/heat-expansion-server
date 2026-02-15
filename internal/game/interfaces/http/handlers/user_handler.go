package handlers

import (
	"net/http"

	"github.com/artcodefun/heat-expansion-server/internal/game/application/cqrs"
	"github.com/artcodefun/heat-expansion-server/internal/game/application/ports"
	"github.com/artcodefun/heat-expansion-server/internal/game/interfaces/http/dtos"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserHandler struct {
	commands   cqrs.UserCommands
	queries    cqrs.UserQueries
	translator ports.Translator
}

func NewUserHandler(commands cqrs.UserCommands, queries cqrs.UserQueries, translator ports.Translator) *UserHandler {
	return &UserHandler{commands: commands, queries: queries, translator: translator}
}

// GetCrystalBalance handles GET /user/balance and returns the authenticated user's crystal balance.
func (h *UserHandler) GetCrystalBalance(c *gin.Context) {
	ctx := queryCtx(c)
	if ctx.UserID == uuid.Nil {
		// Should not normally happen with Auth middleware, but guard just in case.
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	user, err := h.queries.GetUserProfile(ctx, ctx.UserID)
	if handleCoreErr(c, h.translator, err) {
		return
	}

	c.JSON(http.StatusOK, dtos.UserCrystalBalanceResponse{Crystals: user.Crystals})
}
