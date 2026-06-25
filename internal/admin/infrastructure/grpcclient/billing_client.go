package grpcclient

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	billingv1 "github.com/artcodefun/heat-expansion-server/contracts/billing/grpc/v1"
	"github.com/artcodefun/heat-expansion-server/internal/admin/application/cqrs/readmodels"
	"github.com/google/uuid"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
)

// BillingClient implements ports.BillingPrivateClient by calling the billing
// module's private gRPC API. It dials lazily: the connection is established on
// the first RPC call, which avoids races during the shared errgroup startup.
type BillingClient struct {
	pkg billingv1.CrystalPackageServiceClient
}

func NewBillingClient(addr, key string) (*BillingClient, error) {
	conn, err := grpc.NewClient(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
		grpc.WithUnaryInterceptor(keyInterceptor(key)),
	)
	if err != nil {
		return nil, err
	}
	return &BillingClient{pkg: billingv1.NewCrystalPackageServiceClient(conn)}, nil
}

func (c *BillingClient) ListCrystalPackages(ctx context.Context) ([]*readmodels.CrystalPackage, error) {
	resp, err := c.pkg.ListCrystalPackages(ctx, &billingv1.ListCrystalPackagesRequest{})
	if err != nil {
		return nil, grpcErrToSentinel(err)
	}
	out := make([]*readmodels.CrystalPackage, len(resp.Packages))
	for i, p := range resp.Packages {
		out[i] = packageFromProto(p)
	}
	return out, nil
}

func (c *BillingClient) GetCrystalPackage(ctx context.Context, id uuid.UUID) (*readmodels.CrystalPackage, error) {
	resp, err := c.pkg.GetCrystalPackage(ctx, &billingv1.GetCrystalPackageRequest{Id: id.String()})
	if err != nil {
		return nil, grpcErrToSentinel(err)
	}
	return packageFromProto(resp.Package), nil
}

func (c *BillingClient) CreateCrystalPackage(ctx context.Context, p *readmodels.CrystalPackage) (*readmodels.CrystalPackage, error) {
	resp, err := c.pkg.CreateCrystalPackage(ctx, &billingv1.CreateCrystalPackageRequest{
		Name:            p.Name,
		Crystals:        p.Crystals,
		PriceMinorUnits: p.PriceMinorUnits,
		Currency:        p.Currency,
		ImageUrl:        p.ImageURL,
		IsActive:        p.IsActive,
	})
	if err != nil {
		return nil, grpcErrToSentinel(err)
	}
	return packageFromProto(resp.Package), nil
}

func (c *BillingClient) UpdateCrystalPackage(ctx context.Context, p *readmodels.CrystalPackage) (*readmodels.CrystalPackage, error) {
	resp, err := c.pkg.UpdateCrystalPackage(ctx, &billingv1.UpdateCrystalPackageRequest{
		Id:              p.ID.String(),
		Name:            p.Name,
		Crystals:        p.Crystals,
		PriceMinorUnits: p.PriceMinorUnits,
		Currency:        p.Currency,
		ImageUrl:        p.ImageURL,
		IsActive:        p.IsActive,
	})
	if err != nil {
		return nil, grpcErrToSentinel(err)
	}
	return packageFromProto(resp.Package), nil
}

// ── Mapping helpers ───────────────────────────────────────────────────────────

func packageFromProto(p *billingv1.CrystalPackage) *readmodels.CrystalPackage {
	if p == nil {
		return nil
	}
	id, _ := uuid.Parse(p.Id)
	return &readmodels.CrystalPackage{
		ID:              id,
		Name:            p.Name,
		Crystals:        p.Crystals,
		PriceMinorUnits: p.PriceMinorUnits,
		Currency:        p.Currency,
		ImageURL:        p.ImageUrl,
		IsActive:        p.IsActive,
	}
}
