package dtos

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	gamev1 "github.com/artcodefun/heat-expansion-server/contracts/game/grpc/v1"
	"github.com/artcodefun/heat-expansion-server/internal/game/application/cqrs/readmodels"
	"github.com/artcodefun/heat-expansion-server/internal/game/domain"
)

// Validation in this package guards the proto->domain conversion boundary. These
// failures mean the client sent something the domain vocabulary doesn't accept —
// an integration bug, not something an admin can correct — so the messages are
// plain English gRPC status errors (InvalidArgument), not translated keys.

// hasNegative reports whether any of the given values is negative. Prototype
// numeric fields (stats, costs, durations) must never be below zero.
func hasNegative(vals ...int64) bool {
	for _, v := range vals {
		if v < 0 {
			return true
		}
	}
	return false
}

// validateFaction checks f against the known faction set. Factions are shared
// across prototype types and locations.
func validateFaction(f domain.Faction) error {
	switch f {
	case domain.FactionExoCoalition, domain.FactionMarauders, domain.FactionFerrousSwarm,
		domain.FactionTitanArachnids, domain.FactionVoidEcho, domain.FactionCustodianProtocol,
		domain.FactionScorchWalkers, domain.FactionObsidianSentinels, domain.FactionNeuralWormApex:
		return nil
	default:
		return status.Errorf(codes.InvalidArgument, "invalid faction: %q", string(f))
	}
}

// validateCreationSources checks every source against the known set, returning an
// error naming the first offending value. Shared by all prototype types, which
// carry the same creation-source vocabulary.
func validateCreationSources(sources []domain.CreationSource) error {
	for _, s := range sources {
		switch s {
		case domain.CreationSourcePlayerBase, domain.CreationSourceBlackMarket,
			domain.CreationSourceNPCLocation, domain.CreationSourceConsumableBox:
		default:
			return status.Errorf(codes.InvalidArgument, "invalid creation source: %q", string(s))
		}
	}
	return nil
}

func priceToProto(p readmodels.PriceModel) *gamev1.PriceModel {
	return &gamev1.PriceModel{
		Credits:    int64(p.Credits),
		Iron:       int64(p.Iron),
		Titanium:   int64(p.Titanium),
		Antimatter: int64(p.Antimatter),
	}
}

func priceFromProto(p *gamev1.PriceModel) domain.PriceModel {
	if p == nil {
		return domain.PriceModel{}
	}
	return domain.PriceModel{
		Credits:    int(p.Credits),
		Iron:       int(p.Iron),
		Titanium:   int(p.Titanium),
		Antimatter: int(p.Antimatter),
	}
}

func creationSourcesToStrings[T ~string](sources []T) []string {
	if len(sources) == 0 {
		return nil
	}
	out := make([]string, len(sources))
	for i, s := range sources {
		out[i] = string(s)
	}
	return out
}

func creationSourcesFromStrings(sources []string) []domain.CreationSource {
	if len(sources) == 0 {
		return nil
	}
	out := make([]domain.CreationSource, len(sources))
	for i, s := range sources {
		out[i] = domain.CreationSource(s)
	}
	return out
}
