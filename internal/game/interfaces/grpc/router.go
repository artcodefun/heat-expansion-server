package grpc

import (
	"google.golang.org/grpc"

	gamev1 "github.com/artcodefun/heat-expansion-server/contracts/game/grpc/v1"
	"github.com/artcodefun/heat-expansion-server/internal/game/application/cqrs"
	"github.com/artcodefun/heat-expansion-server/internal/game/application/ports"
	"github.com/artcodefun/heat-expansion-server/internal/game/interfaces/grpc/handlers"
	"github.com/artcodefun/heat-expansion-server/internal/game/interfaces/grpc/interceptor"
)

// Commands groups CQRS command interfaces needed by gRPC handlers.
type Commands struct {
	ArmyPrototype cqrs.ArmyPrototypeCommands
}

// Queries groups CQRS query interfaces needed by gRPC handlers.
type Queries struct {
	ArmyPrototype cqrs.ArmyPrototypeQueries
}

type Router struct {
	srv *grpc.Server
}

// NewRouter builds the configured gRPC server: installs interceptors and registers services.
func NewRouter(cmd Commands, qry Queries, internalKey string, tr ports.Translator) Router {
	srv := grpc.NewServer(
		grpc.ChainUnaryInterceptor(interceptor.KeyAuth(internalKey)),
	)
	gamev1.RegisterArmyPrototypeServiceServer(srv, handlers.NewArmyPrototypeHandler(cmd.ArmyPrototype, qry.ArmyPrototype, tr))
	return Router{srv}
}
