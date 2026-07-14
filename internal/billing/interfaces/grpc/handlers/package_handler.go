package handlers

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	billingv1 "github.com/artcodefun/heat-expansion-server/contracts/billing/grpc/v1"
	"github.com/artcodefun/heat-expansion-server/internal/billing/application/cqrs"
	"github.com/artcodefun/heat-expansion-server/internal/billing/application/ports"
	"github.com/artcodefun/heat-expansion-server/internal/billing/interfaces/grpc/dtos"
	"github.com/google/uuid"
)

// PackageHandler implements billingv1.CrystalPackageServiceServer.
type PackageHandler struct {
	billingv1.UnimplementedCrystalPackageServiceServer
	cmd cqrs.CrystalPackageCommands
	qry cqrs.CrystalPackageQueries
	tr  ports.Translator
}

func NewPackageHandler(cmd cqrs.CrystalPackageCommands, qry cqrs.CrystalPackageQueries, tr ports.Translator) *PackageHandler {
	return &PackageHandler{cmd: cmd, qry: qry, tr: tr}
}

func (h *PackageHandler) ListCrystalPackages(ctx context.Context, _ *billingv1.ListCrystalPackagesRequest) (*billingv1.ListCrystalPackagesResponse, error) {
	pkgs, err := h.qry.ListAllCrystalPackages(ctx)
	if err != nil {
		return nil, dtos.StatusFromError(ctx, h.tr, err)
	}
	return &billingv1.ListCrystalPackagesResponse{Packages: dtos.PackagesToProto(pkgs)}, nil
}

func (h *PackageHandler) GetCrystalPackage(ctx context.Context, req *billingv1.GetCrystalPackageRequest) (*billingv1.GetCrystalPackageResponse, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid package id")
	}
	pkg, err := h.qry.GetCrystalPackage(ctx, id)
	if err != nil {
		return nil, dtos.StatusFromError(ctx, h.tr, err)
	}
	return &billingv1.GetCrystalPackageResponse{Package: dtos.PackageToProto(pkg)}, nil
}

func (h *PackageHandler) CreateCrystalPackage(ctx context.Context, req *billingv1.CreateCrystalPackageRequest) (*billingv1.CreateCrystalPackageResponse, error) {
	pkg, err := dtos.PackageFromCreateRequest(req)
	if err != nil {
		return nil, err
	}
	created, err := h.cmd.CreateCrystalPackage(ctx, pkg)
	if err != nil {
		return nil, dtos.StatusFromError(ctx, h.tr, err)
	}
	return &billingv1.CreateCrystalPackageResponse{Package: dtos.PackageDomainToProto(created)}, nil
}

func (h *PackageHandler) UpdateCrystalPackage(ctx context.Context, req *billingv1.UpdateCrystalPackageRequest) (*billingv1.UpdateCrystalPackageResponse, error) {
	pkg, err := dtos.PackageFromUpdateRequest(req)
	if err != nil {
		return nil, err
	}
	updated, err := h.cmd.UpdateCrystalPackage(ctx, pkg)
	if err != nil {
		return nil, dtos.StatusFromError(ctx, h.tr, err)
	}
	return &billingv1.UpdateCrystalPackageResponse{Package: dtos.PackageDomainToProto(updated)}, nil
}
