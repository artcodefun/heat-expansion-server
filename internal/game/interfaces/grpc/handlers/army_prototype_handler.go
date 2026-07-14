package handlers

import (
	"context"

	gamev1 "github.com/artcodefun/heat-expansion-server/contracts/game/grpc/v1"
	"github.com/artcodefun/heat-expansion-server/internal/game/application/cqrs"
	"github.com/artcodefun/heat-expansion-server/internal/game/application/ports"
	"github.com/artcodefun/heat-expansion-server/internal/game/interfaces/grpc/dtos"
)

// ArmyPrototypeHandler implements gamev1.ArmyPrototypeServiceServer.
type ArmyPrototypeHandler struct {
	gamev1.UnimplementedArmyPrototypeServiceServer
	commands   cqrs.ArmyPrototypeCommands
	queries    cqrs.ArmyPrototypeQueries
	translator ports.Translator
}

func NewArmyPrototypeHandler(commands cqrs.ArmyPrototypeCommands, queries cqrs.ArmyPrototypeQueries, translator ports.Translator) *ArmyPrototypeHandler {
	return &ArmyPrototypeHandler{commands: commands, queries: queries, translator: translator}
}

func (h *ArmyPrototypeHandler) ListArmyPrototypes(ctx context.Context, _ *gamev1.ListArmyPrototypesRequest) (*gamev1.ListArmyPrototypesResponse, error) {
	protos, err := h.queries.ListArmyPrototypes(ctx)
	if err != nil {
		return nil, dtos.StatusFromError(ctx, h.translator, err)
	}
	return &gamev1.ListArmyPrototypesResponse{Prototypes: dtos.ArmyPrototypesToProto(protos)}, nil
}

func (h *ArmyPrototypeHandler) GetArmyPrototype(ctx context.Context, req *gamev1.GetArmyPrototypeRequest) (*gamev1.GetArmyPrototypeResponse, error) {
	proto, err := h.queries.GetArmyPrototype(ctx, int(req.GetId()))
	if err != nil {
		return nil, dtos.StatusFromError(ctx, h.translator, err)
	}
	return &gamev1.GetArmyPrototypeResponse{Prototype: dtos.ArmyPrototypeToProto(proto)}, nil
}

func (h *ArmyPrototypeHandler) CreateArmyPrototype(ctx context.Context, req *gamev1.CreateArmyPrototypeRequest) (*gamev1.CreateArmyPrototypeResponse, error) {
	proto, err := dtos.ArmyPrototypeFromProto(req.GetPrototype())
	if err != nil {
		return nil, err // already a plain InvalidArgument status error
	}
	created, err := h.commands.CreateArmyPrototype(ctx, proto)
	if err != nil {
		return nil, dtos.StatusFromError(ctx, h.translator, err)
	}
	return &gamev1.CreateArmyPrototypeResponse{Prototype: dtos.ArmyPrototypeDomainToProto(created)}, nil
}

func (h *ArmyPrototypeHandler) UpdateArmyPrototype(ctx context.Context, req *gamev1.UpdateArmyPrototypeRequest) (*gamev1.UpdateArmyPrototypeResponse, error) {
	proto, err := dtos.ArmyPrototypeFromProto(req.GetPrototype())
	if err != nil {
		return nil, err // already a plain InvalidArgument status error
	}
	updated, err := h.commands.UpdateArmyPrototype(ctx, proto)
	if err != nil {
		return nil, dtos.StatusFromError(ctx, h.translator, err)
	}
	return &gamev1.UpdateArmyPrototypeResponse{Prototype: dtos.ArmyPrototypeDomainToProto(updated)}, nil
}
