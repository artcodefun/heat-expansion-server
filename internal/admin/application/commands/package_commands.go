package commands

import (
	"context"

	"github.com/artcodefun/heat-expansion-server/internal/admin/application/cqrs"
	"github.com/artcodefun/heat-expansion-server/internal/admin/application/cqrs/readmodels"
	"github.com/artcodefun/heat-expansion-server/internal/admin/application/ports"
)

// PackageCommands delegates crystal package writes to the billing private gRPC API.
type PackageCommands struct {
	billing ports.BillingPrivateClient
}

func NewPackageCommands(billing ports.BillingPrivateClient) *PackageCommands {
	return &PackageCommands{billing: billing}
}

func (c *PackageCommands) CreateCrystalPackage(ctx context.Context, _ cqrs.Actor, p *readmodels.CrystalPackage) (*readmodels.CrystalPackage, error) {
	created, err := c.billing.CreateCrystalPackage(ctx, p)
	return created, clientErr(err)
}

func (c *PackageCommands) UpdateCrystalPackage(ctx context.Context, _ cqrs.Actor, p *readmodels.CrystalPackage) (*readmodels.CrystalPackage, error) {
	updated, err := c.billing.UpdateCrystalPackage(ctx, p)
	return updated, clientErr(err)
}
