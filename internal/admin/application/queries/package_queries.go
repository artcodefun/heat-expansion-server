package queries

import (
	"context"

	"github.com/artcodefun/heat-expansion-server/internal/admin/application/cqrs"
	"github.com/artcodefun/heat-expansion-server/internal/admin/application/ports"
	"github.com/google/uuid"
)

// PackageQueries delegates crystal package reads to the billing private gRPC API.
type PackageQueries struct {
	billing ports.BillingPrivateClient
}

func NewPackageQueries(billing ports.BillingPrivateClient) *PackageQueries {
	return &PackageQueries{billing: billing}
}

func (q *PackageQueries) ListCrystalPackages(ctx context.Context, _ cqrs.Actor) ([]*ports.CrystalPackage, error) {
	list, err := q.billing.ListCrystalPackages(ctx)
	return list, clientErr(err)
}

func (q *PackageQueries) GetCrystalPackage(ctx context.Context, _ cqrs.Actor, id uuid.UUID) (*ports.CrystalPackage, error) {
	p, err := q.billing.GetCrystalPackage(ctx, id)
	return p, clientErr(err)
}
