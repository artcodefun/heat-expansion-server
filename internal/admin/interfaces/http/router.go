package http

import (
	"net/http"

	contract "github.com/artcodefun/heat-expansion-server/contracts/admin/http/v1"
	"github.com/artcodefun/heat-expansion-server/internal/admin/application/cqrs"
	"github.com/artcodefun/heat-expansion-server/internal/admin/application/ports"
	"github.com/artcodefun/heat-expansion-server/internal/admin/interfaces/http/handlers"
	"github.com/artcodefun/heat-expansion-server/internal/admin/interfaces/http/middleware"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

// Commands groups CQRS command interfaces needed by HTTP handlers.
type Commands struct {
	Admin       cqrs.AdminCommands
	Prototype   cqrs.PrototypeCommands
	Translation cqrs.TranslationCommands
	Package     cqrs.PackageCommands
}

// Queries groups CQRS query interfaces needed by HTTP handlers.
type Queries struct {
	Admin       cqrs.AdminQueries
	Prototype   cqrs.PrototypeQueries
	Translation cqrs.TranslationQueries
	Package     cqrs.PackageQueries
}

// NewRouter constructs the Gin engine, registers middleware and routes.
func NewRouter(cmd Commands, qry Queries, sessionValidator ports.SessionValidator, tr ports.Translator) *gin.Engine {
	r := gin.Default()
	r.Use(otelgin.Middleware("heat-expansion-admin"))

	adminHandler := handlers.NewAdminHandler(cmd.Admin, qry.Admin, tr)
	protoHandler := handlers.NewPrototypeHandler(cmd.Prototype, qry.Prototype, tr)
	translationHandler := handlers.NewTranslationHandler(cmd.Translation, qry.Translation, tr)
	packageHandler := handlers.NewPackageHandler(cmd.Package, qry.Package, tr)

	// Global routes
	r.GET("/health", HealthHandler)

	// Public routes (no session required)
	publicApi := r.Group("/api/v1")
	{
		publicApi.GET("/openapi.yaml", func(c *gin.Context) {
			c.Data(http.StatusOK, "application/yaml; charset=utf-8", contract.OpenAPI())
		})
		publicApi.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("../openapi.yaml")))

		auth := publicApi.Group("/auth")
		{
			auth.POST("/register", adminHandler.Register)
			auth.POST("/login", adminHandler.Login)
		}
	}

	// Private routes (valid session required)
	api := r.Group("/api/v1")
	api.Use(middleware.Auth(sessionValidator))
	{
		auth := api.Group("/auth")
		{
			auth.POST("/logout", adminHandler.Logout)
			auth.GET("/me", adminHandler.Me)
		}

		// Game prototype catalog
		army := api.Group("/game/prototypes/army")
		{
			army.GET("", protoHandler.ListArmy)
			army.GET("/:id", protoHandler.GetArmy)
			army.POST("", protoHandler.CreateArmy)
			army.PUT("/:id", protoHandler.UpdateArmy)
		}

		build := api.Group("/game/prototypes/build")
		{
			build.GET("", protoHandler.ListBuild)
			build.GET("/:id", protoHandler.GetBuild)
			build.POST("", protoHandler.CreateBuild)
			build.PUT("/:id", protoHandler.UpdateBuild)
		}

		storage := api.Group("/game/prototypes/storage")
		{
			storage.GET("", protoHandler.ListStorage)
			storage.GET("/:id", protoHandler.GetStorage)
			storage.POST("", protoHandler.CreateStorage)
			storage.PUT("/:id", protoHandler.UpdateStorage)
		}

		tech := api.Group("/game/prototypes/tech")
		{
			tech.GET("", protoHandler.ListTech)
			tech.GET("/:id", protoHandler.GetTech)
			tech.POST("", protoHandler.CreateTech)
			tech.PUT("/:id", protoHandler.UpdateTech)
		}

		// Game translation catalog
		translations := api.Group("/game/translations")
		{
			translations.GET("", translationHandler.ListTranslations)
			translations.PUT("", translationHandler.UpsertTranslation)
		}

		// Billing crystal package catalog
		packages := api.Group("/billing/packages")
		{
			packages.GET("", packageHandler.ListPackages)
			packages.GET("/:id", packageHandler.GetPackage)
			packages.POST("", packageHandler.CreatePackage)
			packages.PUT("/:id", packageHandler.UpdatePackage)
		}
	}

	return r
}
