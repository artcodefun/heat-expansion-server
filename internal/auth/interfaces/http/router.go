package http

import (
	"net/http"

	contract "github.com/artcodefun/heat-expansion-server/contracts/auth/http/v1"
	"github.com/artcodefun/heat-expansion-server/internal/auth/application/cqrs"
	"github.com/artcodefun/heat-expansion-server/internal/auth/application/ports"
	"github.com/artcodefun/heat-expansion-server/internal/auth/interfaces/http/handlers"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
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
	r.Use(otelgin.Middleware("heat-expansion-auth"))

	handler := handlers.NewAccountHandler(cmd.Account, qry.Account, tr)

	r.GET("/health", HealthHandler)

	v1 := r.Group("/api/v1")
	{
		v1.GET("/openapi.yaml", func(c *gin.Context) {
			c.Data(http.StatusOK, "application/yaml; charset=utf-8", contract.OpenAPI())
		})
		v1.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("../openapi.yaml")))

		v1.POST("/register", handler.Register)
		v1.POST("/login", handler.Login)
		v1.POST("/password-reset/request", handler.RequestPasswordReset)
		v1.POST("/password-reset/confirm", handler.ResetPassword)
	}

	return r
}
