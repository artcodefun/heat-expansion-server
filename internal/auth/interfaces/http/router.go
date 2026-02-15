package http

import (
	"github.com/artcodefun/heat-expansion-server/internal/auth/application/cqrs"
	"github.com/artcodefun/heat-expansion-server/internal/auth/application/ports"
	"github.com/artcodefun/heat-expansion-server/internal/auth/interfaces/http/handlers"
	"github.com/gin-gonic/gin"
)

// Commands groups CQRS command interfaces needed by HTTP handlers.
type Commands struct {
	Account cqrs.AccountCommands
}

// Queries groups CQRS query interfaces needed by HTTP handlers.
type Queries struct {
	Account cqrs.AccountQueries
}

func NewRouter(cmd Commands, qry Queries, tr ports.Translator) *gin.Engine {
	r := gin.Default()

	handler := handlers.NewAccountHandler(cmd.Account, qry.Account, tr)

	r.GET("/health", HealthHandler)

	v1 := r.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/register", handler.Register)
			auth.POST("/login", handler.Login)
		}
	}

	return r
}
