package dtos

import (
	billingv1 "github.com/artcodefun/heat-expansion-server/contracts/billing/grpc/v1"
	"github.com/artcodefun/heat-expansion-server/internal/billing/application/cqrs/readmodels"
	"github.com/artcodefun/heat-expansion-server/internal/billing/domain"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// PackageToProto converts a read model to the proto wire type.
func PackageToProto(p *readmodels.CrystalPackage) *billingv1.CrystalPackage {
	return &billingv1.CrystalPackage{
		Id:              p.ID.String(),
		Name:            p.Name,
		Crystals:        int32(p.Crystals),
		PriceMinorUnits: p.PriceMinorUnits,
		Currency:        p.Currency,
		ImageUrl:        p.ImageURL,
		IsActive:        p.IsActive,
	}
}

// PackagesToProto converts a slice of read models to proto wire types.
func PackagesToProto(ps []*readmodels.CrystalPackage) []*billingv1.CrystalPackage {
	out := make([]*billingv1.CrystalPackage, len(ps))
	for i, p := range ps {
		out[i] = PackageToProto(p)
	}
	return out
}

// PackageDomainToProto converts a domain type (returned by Create/Update) to the proto wire type.
func PackageDomainToProto(p *domain.CrystalPackage) *billingv1.CrystalPackage {
	return &billingv1.CrystalPackage{
		Id:              p.ID.String(),
		Name:            p.Name,
		Crystals:        int32(p.Crystals),
		PriceMinorUnits: p.PriceMinorUnits,
		Currency:        p.Currency,
		ImageUrl:        p.ImageURL,
		IsActive:        p.IsActive,
	}
}

// PackageFromCreateRequest constructs a domain CrystalPackage from a Create RPC request.
func PackageFromCreateRequest(req *billingv1.CreateCrystalPackageRequest) (*domain.CrystalPackage, error) {
	if err := validatePackageFields(req.Name, req.Currency); err != nil {
		return nil, err
	}
	return &domain.CrystalPackage{
		ID:              uuid.New(),
		Name:            req.Name,
		Crystals:        int(req.Crystals),
		PriceMinorUnits: req.PriceMinorUnits,
		Currency:        req.Currency,
		ImageURL:        req.ImageUrl,
		IsActive:        req.IsActive,
	}, nil
}

// PackageFromUpdateRequest constructs a domain CrystalPackage from an Update RPC request.
func PackageFromUpdateRequest(req *billingv1.UpdateCrystalPackageRequest) (*domain.CrystalPackage, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid package id")
	}
	if err := validatePackageFields(req.Name, req.Currency); err != nil {
		return nil, err
	}
	return &domain.CrystalPackage{
		ID:              id,
		Name:            req.Name,
		Crystals:        int(req.Crystals),
		PriceMinorUnits: req.PriceMinorUnits,
		Currency:        req.Currency,
		ImageURL:        req.ImageUrl,
		IsActive:        req.IsActive,
	}, nil
}

func validatePackageFields(name, currency string) error {
	if name == "" {
		return status.Error(codes.InvalidArgument, "name is required")
	}
	if currency == "" {
		return status.Error(codes.InvalidArgument, "currency is required")
	}
	return nil
}
