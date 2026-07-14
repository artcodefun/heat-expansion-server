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

// StoragePrototypeHandler implements gamev1.StoragePrototypeServiceServer.
type StoragePrototypeHandler struct {
	gamev1.UnimplementedStoragePrototypeServiceServer
	commands   cqrs.StoragePrototypeCommands
	queries    cqrs.StoragePrototypeQueries
	translator ports.Translator
}

func NewStoragePrototypeHandler(
	commands cqrs.StoragePrototypeCommands,
	queries cqrs.StoragePrototypeQueries,
	translator ports.Translator,
) *StoragePrototypeHandler {
	return &StoragePrototypeHandler{commands: commands, queries: queries, translator: translator}
}

func (h *StoragePrototypeHandler) ListStoragePrototypes(ctx context.Context, _ *gamev1.ListStoragePrototypesRequest) (*gamev1.ListStoragePrototypesResponse, error) {
	protos, err := h.queries.ListStoragePrototypes(ctx)
	if err != nil {
		return nil, dtos.StatusFromError(ctx, h.translator, err)
	}
	return &gamev1.ListStoragePrototypesResponse{Prototypes: dtos.StoragePrototypesToProto(protos)}, nil
}

func (h *StoragePrototypeHandler) GetStoragePrototype(ctx context.Context, req *gamev1.GetStoragePrototypeRequest) (*gamev1.GetStoragePrototypeResponse, error) {
	if req.GetId() <= 0 {
		return nil, status.Error(codes.InvalidArgument, "id is required and must be positive")
	}
	proto, err := h.queries.GetStoragePrototype(ctx, int(req.GetId()))
	if err != nil {
		return nil, dtos.StatusFromError(ctx, h.translator, err)
	}
	return &gamev1.GetStoragePrototypeResponse{Prototype: dtos.StoragePrototypeToProto(proto)}, nil
}

func (h *StoragePrototypeHandler) CreateStoragePrototype(ctx context.Context, req *gamev1.CreateStoragePrototypeRequest) (*gamev1.CreateStoragePrototypeResponse, error) {
	proto, err := dtos.StoragePrototypeFromProto(req.GetPrototype())
	if err != nil {
		return nil, err
	}
	created, err := h.commands.CreateStoragePrototype(ctx, proto)
	if err != nil {
		return nil, dtos.StatusFromError(ctx, h.translator, err)
	}
	return &gamev1.CreateStoragePrototypeResponse{Prototype: dtos.StoragePrototypeDomainToProto(created)}, nil
}

func (h *StoragePrototypeHandler) UpdateStoragePrototype(ctx context.Context, req *gamev1.UpdateStoragePrototypeRequest) (*gamev1.UpdateStoragePrototypeResponse, error) {
	proto, err := dtos.StoragePrototypeFromProto(req.GetPrototype())
	if err != nil {
		return nil, err
	}
	updated, err := h.commands.UpdateStoragePrototype(ctx, proto)
	if err != nil {
		return nil, dtos.StatusFromError(ctx, h.translator, err)
	}
	return &gamev1.UpdateStoragePrototypeResponse{Prototype: dtos.StoragePrototypeDomainToProto(updated)}, nil
}
