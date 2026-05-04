package http

import (
	"github.com/artcodefun/heat-expansion-server/internal/game/application/cqrs"
	"github.com/artcodefun/heat-expansion-server/internal/game/application/ports"
	"github.com/artcodefun/heat-expansion-server/internal/game/interfaces/http/dtos"
	"github.com/artcodefun/heat-expansion-server/internal/game/interfaces/http/handlers"
	"github.com/artcodefun/heat-expansion-server/internal/game/interfaces/http/middleware"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

// Commands groups CQRS command interfaces needed by HTTP handlers.
type Commands struct {
	User      cqrs.UserCommands
	Base      cqrs.BaseCommands
	Building  cqrs.BuildingCommands
	Army      cqrs.ArmyCommands
	Tech      cqrs.TechCommands
	Storage   cqrs.StorageCommands
	Operation cqrs.OperationCommands
	Trade     cqrs.TradeCommands
	Alert     cqrs.AlertCommands
	Diplomacy cqrs.DiplomacyCommands
}

// Queries groups CQRS query interfaces needed by HTTP handlers.
type Queries struct {
	User      cqrs.UserQueries
	Base      cqrs.BaseQueries
	Building  cqrs.BuildingQueries
	Army      cqrs.ArmyQueries
	Tech      cqrs.TechQueries
	Storage   cqrs.StorageQueries
	Trade     cqrs.TradeQueries
	Sector    cqrs.SectorQueries
	Radar     cqrs.RadarQueries
	Operation cqrs.OperationQueries
	Activity  cqrs.ActivityQueries
	Alert     cqrs.AlertQueries
	Diplomacy cqrs.DiplomacyQueries
}

// NewRouter constructs the Gin engine, registers middleware and routes.
func NewRouter(cmd Commands, qry Queries, tokenValidator ports.TokenValidator, tr ports.Translator) *gin.Engine {
	r := gin.Default()
	r.Use(otelgin.Middleware("heat-expansion-game"))
	registerCustomValidators()

	// Initialize handlers at the top for consistency
	userHandler := handlers.NewUserHandler(cmd.User, qry.User, tr)
	baseHandler := handlers.NewBaseHandler(qry.Base, cmd.Base, tr)
	buildingHandler := handlers.NewBuildingHandler(qry.Building, cmd.Building, tr)
	armyHandler := handlers.NewArmyHandler(qry.Army, cmd.Army, tr)
	techHandler := handlers.NewTechHandler(qry.Tech, cmd.Tech, tr)
	storageHandler := handlers.NewStorageHandler(qry.Storage, cmd.Storage, tr)
	tradeHandler := handlers.NewTradeHandler(cmd.Trade, qry.Trade, tr)
	sectorHandler := handlers.NewSectorHandler(qry.Sector, tr)
	radarHandler := handlers.NewRadarHandler(qry.Radar, tr)
	operationHandler := handlers.NewOperationHandler(qry.Operation, cmd.Operation, tr)
	activityHandler := handlers.NewActivityHandler(qry.Activity, tr)
	alertHandler := handlers.NewAlertHandler(qry.Alert, cmd.Alert, tr)
	diplomacyHandler := handlers.NewDiplomacyHandler(qry.Diplomacy, cmd.Diplomacy, tr)

	// Global routes
	r.GET("/health", HealthHandler)

	// Public routes (no auth)
	publicApi := r.Group("/api/v1")
	{
		publicApi.GET("/min-client-version", MinClientVersionHandler)
	}

	// Private routes (auth required)
	api := r.Group("/api/v1")
	api.Use(middleware.Auth(tokenValidator))
	{
		// User
		api.GET("/user/balance", userHandler.GetCrystalBalance)

		// Base
		api.GET("/bases", baseHandler.ListUserBases)
		api.GET("/bases/:baseId/status", baseHandler.GetBaseStatus)

		// Buildings
		buildings := api.Group("/bases/:baseId/buildings")
		{
			buildings.GET("/new", buildingHandler.ListNew)
			buildings.GET("/pending", buildingHandler.ListPending)
			buildings.GET("/in-production", buildingHandler.ListInProduction)
			buildings.GET("/present", buildingHandler.ListPresent)
			buildings.POST("/queue", buildingHandler.Queue)
			buildings.POST("/production/:itemId/speed-up", buildingHandler.SpeedUpProduction)
			buildings.POST("/pending/:itemId/cancel", buildingHandler.CancelPending)
			buildings.DELETE("/present/:itemId", buildingHandler.DeletePresent)
		}

		// Army
		army := api.Group("/bases/:baseId/army")
		{
			army.GET("/new", armyHandler.ListNew)
			army.GET("/pending", armyHandler.ListPending)
			army.GET("/in-production", armyHandler.ListInProduction)
			army.GET("/present", armyHandler.ListPresent)
			army.POST("/queue", armyHandler.Queue)
			army.POST("/production/:itemId/speed-up", armyHandler.SpeedUpProduction)
			army.POST("/pending/:itemId/cancel", armyHandler.CancelPending)
			army.DELETE("/present/:itemId", armyHandler.DeletePresent)
		}

		// Tech
		tech := api.Group("/bases/:baseId/tech")
		{
			tech.GET("/new", techHandler.ListNew)
			tech.GET("/in-progress", techHandler.ListInProgress)
			tech.GET("/done", techHandler.ListDone)
			tech.POST("/queue", techHandler.Queue)
			tech.POST("/production/:itemId/speed-up", techHandler.SpeedUpProduction)
		}

		// Storage
		storage := api.Group("/bases/:baseId/storage")
		{
			storage.GET("/present", storageHandler.ListPresent)
			storage.DELETE("/items/:itemId", storageHandler.DeleteItem)
			storage.POST("/items/:itemId/activate", storageHandler.ActivateBuff)
			storage.POST("/items/:itemId/decrypt", storageHandler.StartIntelDecryption)
			storage.POST("/items/:itemId/restore", storageHandler.StartDamagedItemRestoration)
			storage.POST("/items/:itemId/enable", storageHandler.ActivateArtifact)
			storage.POST("/items/:itemId/disable", storageHandler.DeactivateArtifact)
			storage.POST("/items/:itemId/open", storageHandler.OpenBox)
		}

		// Sector scan reports
		sectors := api.Group("/bases/:baseId/sectors")
		{
			sectors.GET("/scans/near", sectorHandler.GetScansNear)
			sectors.GET("/scans/:id", sectorHandler.GetScanByID)
			sectors.GET("/scans/before", sectorHandler.GetLatestScanBefore)
		}

		// Radar threats
		api.GET("/bases/:baseId/threats", radarHandler.ListIncomingThreats)

		// Operations
		operations := api.Group("/bases/:baseId/operations")
		{
			operations.GET("", operationHandler.ListByBase)
			operations.GET("/active", operationHandler.ListActive)
			operations.GET("/:operationId", operationHandler.GetOperation)
			operations.POST("/:operationId/speed-up", operationHandler.SpeedUp)
			operations.POST("/:operationId/cancel", operationHandler.Cancel)
			operations.POST("", operationHandler.Create)
		}

		// Activities
		activities := api.Group("/bases/:baseId/activities")
		{
			activities.GET("/offense", activityHandler.ListOffense)
			activities.GET("/defense", activityHandler.ListDefense)
			activities.GET("/scan", activityHandler.ListScan)
			activities.GET("/radar", activityHandler.ListRadar)
			activities.GET("/trade", activityHandler.ListTrade)
		}

		// Alerts
		alerts := api.Group("/alerts")
		{
			alerts.GET("", alertHandler.ListActive)
			alerts.GET("/unread-count", alertHandler.GetUnreadCount)
			alerts.POST("/read-all", alertHandler.MarkAllAsRead)
		}

		// Trade
		trade := api.Group("/trade/bases/:baseId")
		{
			trade.GET("/info", tradeHandler.GetInfo)

			operations := trade.Group("/operations")
			{
				operations.GET("", tradeHandler.ListActive)
				operations.POST("", tradeHandler.Create)
				operations.GET("/:operationId", tradeHandler.GetOperation)
				operations.POST("/:operationId/accept", tradeHandler.Accept)
				operations.POST("/:operationId/decline", tradeHandler.Decline)
				operations.POST("/:operationId/cancel", tradeHandler.Cancel)
				operations.POST("/:operationId/speed-up", tradeHandler.SpeedUp)
			}
		}

		// Diplomacy
		diplomacy := api.Group("/diplomacy")
		{
			relationships := diplomacy.Group("/relationships")
			{
				relationships.GET("", diplomacyHandler.ListRelationships)
				relationships.GET("/:userId", diplomacyHandler.GetRelationship)
				relationships.POST("/:userId/declare-war", diplomacyHandler.DeclareWar)
				relationships.POST("/:userId/break-alliance", diplomacyHandler.BreakAlliance)
			}

			messages := diplomacy.Group("/messages")
			{
				messages.GET("/available", diplomacyHandler.ListAvailableInformationalMessages)
				messages.GET("/chats", diplomacyHandler.ListChats)
				messages.GET("/unread-count", diplomacyHandler.GetUnreadCount)
				messages.GET("/chats/:userId", diplomacyHandler.ListChatMessages)
				messages.POST("/chats/:userId/read", diplomacyHandler.MarkChatAsRead)
				messages.POST("", diplomacyHandler.SendInformationalMessage)
			}

			requests := diplomacy.Group("/requests")
			{
				requests.GET("/pending", diplomacyHandler.ListPendingRequests)
				requests.POST("", diplomacyHandler.SendRequest)
				requests.POST("/:requestId/accept", diplomacyHandler.AcceptRequest)
				requests.POST("/:requestId/reject", diplomacyHandler.RejectRequest)
			}
		}
	}

	return r
}

func registerCustomValidators() {
	if validatorEngine, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = validatorEngine.RegisterValidation("army_category", func(fl validator.FieldLevel) bool {
			return dtos.IsValidArmyCategory(fl.Field().String())
		})
		_ = validatorEngine.RegisterValidation("build_category", func(fl validator.FieldLevel) bool {
			return dtos.IsValidBuildCategory(fl.Field().String())
		})
		_ = validatorEngine.RegisterValidation("operation_type", func(fl validator.FieldLevel) bool {
			return dtos.IsValidMilitaryOperationType(fl.Field().String())
		})
		_ = validatorEngine.RegisterValidation("operation_id", func(fl validator.FieldLevel) bool {
			return dtos.IsValidOperationID(fl.Field().String())
		})
		_ = validatorEngine.RegisterValidation("diplomatic_message_content", func(fl validator.FieldLevel) bool {
			return dtos.IsValidUserSendableDiplomaticMessageContent(fl.Field().String())
		})
		_ = validatorEngine.RegisterValidation("diplomatic_request_kind", func(fl validator.FieldLevel) bool {
			return dtos.IsValidDiplomaticRequestKind(fl.Field().String())
		})
		_ = validatorEngine.RegisterValidation("diplomatic_relationship_status", func(fl validator.FieldLevel) bool {
			return dtos.IsValidDiplomaticStatus(fl.Field().String())
		})
		_ = validatorEngine.RegisterValidation("tech_category", func(fl validator.FieldLevel) bool {
			return dtos.IsValidTechCategory(fl.Field().String())
		})
		_ = validatorEngine.RegisterValidation("storage_category", func(fl validator.FieldLevel) bool {
			return dtos.IsValidStorageCategory(fl.Field().String())
		})
	}
}
