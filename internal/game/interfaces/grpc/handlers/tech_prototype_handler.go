package handlers

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	gamev1 "github.com/artcodefun/heat-expansion-server/contracts/game/grpc/v1"
	"github.com/artcodefun/heat-expansion-server/internal/game/application/cqrs"
	"github.com/artcodefun/heat-expansion-server/internal/game/application/ports"
	"github.com/artcodefun/heat-expansion-server/internal/game/interfaces/grpc/dtos"
)

// TechPrototypeHandler implements gamev1.TechPrototypeServiceServer.
type TechPrototypeHandler struct {
	gamev1.UnimplementedTechPrototypeServiceServer
	commands   cqrs.TechPrototypeCommands
	queries    cqrs.TechPrototypeQueries
	translator ports.Translator
}

func NewTechPrototypeHandler(
	commands cqrs.TechPrototypeCommands,
	queries cqrs.TechPrototypeQueries,
	translator ports.Translator,
) *TechPrototypeHandler {
	return &TechPrototypeHandler{commands: commands, queries: queries, translator: translator}
}

func (h *TechPrototypeHandler) ListTechPrototypes(ctx context.Context, _ *gamev1.ListTechPrototypesRequest) (*gamev1.ListTechPrototypesResponse, error) {
	protos, err := h.queries.ListTechPrototypes(ctx)
	if err != nil {
		return nil, dtos.StatusFromError(ctx, h.translator, err)
	}
	return &gamev1.ListTechPrototypesResponse{Prototypes: dtos.TechPrototypesToProto(protos)}, nil
}

func (h *TechPrototypeHandler) GetTechPrototype(ctx context.Context, req *gamev1.GetTechPrototypeRequest) (*gamev1.GetTechPrototypeResponse, error) {
	if req.GetId() <= 0 {
		return nil, status.Error(codes.InvalidArgument, "id is required and must be positive")
	}
	proto, err := h.queries.GetTechPrototype(ctx, int(req.GetId()))
	if err != nil {
		return nil, dtos.StatusFromError(ctx, h.translator, err)
	}
	return &gamev1.GetTechPrototypeResponse{Prototype: dtos.TechPrototypeToProto(proto)}, nil
}

func (h *TechPrototypeHandler) CreateTechPrototype(ctx context.Context, req *gamev1.CreateTechPrototypeRequest) (*gamev1.CreateTechPrototypeResponse, error) {
	proto, err := dtos.TechPrototypeFromProto(req.GetPrototype())
	if err != nil {
		return nil, err
	}
	created, err := h.commands.CreateTechPrototype(ctx, proto)
	if err != nil {
		return nil, dtos.StatusFromError(ctx, h.translator, err)
	}
	return &gamev1.CreateTechPrototypeResponse{Prototype: dtos.TechPrototypeDomainToProto(created)}, nil
}

func (h *TechPrototypeHandler) UpdateTechPrototype(ctx context.Context, req *gamev1.UpdateTechPrototypeRequest) (*gamev1.UpdateTechPrototypeResponse, error) {
	proto, err := dtos.TechPrototypeFromProto(req.GetPrototype())
	if err != nil {
		return nil, err
	}
	updated, err := h.commands.UpdateTechPrototype(ctx, proto)
	if err != nil {
		return nil, dtos.StatusFromError(ctx, h.translator, err)
	}
	return &gamev1.UpdateTechPrototypeResponse{Prototype: dtos.TechPrototypeDomainToProto(updated)}, nil
}
