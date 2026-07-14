package handlers

import (
	"context"

	gamev1 "github.com/artcodefun/heat-expansion-server/contracts/game/grpc/v1"
	"github.com/artcodefun/heat-expansion-server/internal/game/application/cqrs"
	"github.com/artcodefun/heat-expansion-server/internal/game/application/ports"
	"github.com/artcodefun/heat-expansion-server/internal/game/interfaces/grpc/dtos"
)

// BuildPrototypeHandler implements gamev1.BuildPrototypeServiceServer.
type BuildPrototypeHandler struct {
	gamev1.UnimplementedBuildPrototypeServiceServer
	commands   cqrs.BuildPrototypeCommands
	queries    cqrs.BuildPrototypeQueries
	translator ports.Translator
}

func NewBuildPrototypeHandler(commands cqrs.BuildPrototypeCommands, queries cqrs.BuildPrototypeQueries, translator ports.Translator) *BuildPrototypeHandler {
	return &BuildPrototypeHandler{commands: commands, queries: queries, translator: translator}
}

func (h *BuildPrototypeHandler) ListBuildPrototypes(ctx context.Context, _ *gamev1.ListBuildPrototypesRequest) (*gamev1.ListBuildPrototypesResponse, error) {
	protos, err := h.queries.ListBuildPrototypes(ctx)
	if err != nil {
		return nil, dtos.StatusFromError(ctx, h.translator, err)
	}
	return &gamev1.ListBuildPrototypesResponse{Prototypes: dtos.BuildPrototypesToProto(protos)}, nil
}

func (h *BuildPrototypeHandler) GetBuildPrototype(ctx context.Context, req *gamev1.GetBuildPrototypeRequest) (*gamev1.GetBuildPrototypeResponse, error) {
	proto, err := h.queries.GetBuildPrototype(ctx, int(req.GetId()))
	if err != nil {
		return nil, dtos.StatusFromError(ctx, h.translator, err)
	}
	return &gamev1.GetBuildPrototypeResponse{Prototype: dtos.BuildPrototypeToProto(proto)}, nil
}

func (h *BuildPrototypeHandler) CreateBuildPrototype(ctx context.Context, req *gamev1.CreateBuildPrototypeRequest) (*gamev1.CreateBuildPrototypeResponse, error) {
	proto, err := dtos.BuildPrototypeFromProto(req.GetPrototype())
	if err != nil {
		return nil, err // already a plain InvalidArgument status error
	}
	created, err := h.commands.CreateBuildPrototype(ctx, proto)
	if err != nil {
		return nil, dtos.StatusFromError(ctx, h.translator, err)
	}
	return &gamev1.CreateBuildPrototypeResponse{Prototype: dtos.BuildPrototypeDomainToProto(created)}, nil
}

func (h *BuildPrototypeHandler) UpdateBuildPrototype(ctx context.Context, req *gamev1.UpdateBuildPrototypeRequest) (*gamev1.UpdateBuildPrototypeResponse, error) {
	proto, err := dtos.BuildPrototypeFromProto(req.GetPrototype())
	if err != nil {
		return nil, err // already a plain InvalidArgument status error
	}
	updated, err := h.commands.UpdateBuildPrototype(ctx, proto)
	if err != nil {
		return nil, dtos.StatusFromError(ctx, h.translator, err)
	}
	return &gamev1.UpdateBuildPrototypeResponse{Prototype: dtos.BuildPrototypeDomainToProto(updated)}, nil
}
