package grpc

import (
	"google.golang.org/grpc"

	billingv1 "github.com/artcodefun/heat-expansion-server/contracts/billing/grpc/v1"
	"github.com/artcodefun/heat-expansion-server/internal/billing/application/cqrs"
	"github.com/artcodefun/heat-expansion-server/internal/billing/application/ports"
	"github.com/artcodefun/heat-expansion-server/internal/billing/interfaces/grpc/handlers"
	"github.com/artcodefun/heat-expansion-server/internal/billing/interfaces/grpc/interceptor"
)

// Commands groups CQRS command interfaces needed by gRPC handlers.
type Commands struct {
	Package cqrs.CrystalPackageCommands
}

// Queries groups CQRS query interfaces needed by gRPC handlers.
type Queries struct {
	Package cqrs.CrystalPackageQueries
}

type Router struct {
	srv *grpc.Server
}

// NewRouter builds the configured gRPC server: installs interceptors and registers services.
func NewRouter(cmd Commands, qry Queries, internalKey string, tr ports.Translator) Router {
	srv := grpc.NewServer(
		grpc.ChainUnaryInterceptor(interceptor.KeyAuth(internalKey)),
	)
	billingv1.RegisterCrystalPackageServiceServer(srv, handlers.NewPackageHandler(cmd.Package, qry.Package, tr))
	return Router{srv}
}
