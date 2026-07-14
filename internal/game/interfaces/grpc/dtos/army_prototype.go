package dtos

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	gamev1 "github.com/artcodefun/heat-expansion-server/contracts/game/grpc/v1"
	"github.com/artcodefun/heat-expansion-server/internal/game/application/cqrs/readmodels"
	"github.com/artcodefun/heat-expansion-server/internal/game/domain"
)

// ArmyPrototypesToProto maps a slice of readmodel army prototypes to their wire shapes.
func ArmyPrototypesToProto(protos []*readmodels.ArmyItemPrototype) []*gamev1.ArmyPrototype {
	out := make([]*gamev1.ArmyPrototype, len(protos))
	for i, p := range protos {
		out[i] = ArmyPrototypeToProto(p)
	}
	return out
}

// ArmyPrototypeToProto maps a readmodel army prototype to its wire shape.
func ArmyPrototypeToProto(p *readmodels.ArmyItemPrototype) *gamev1.ArmyPrototype {
	if p == nil {
		return nil
	}
	out := &gamev1.ArmyPrototype{
		Id:               int64(p.ID),
		Name:             string(p.Name),
		Category:         string(p.Category),
		CreationSources:  creationSourcesToStrings(p.CreationSources),
		Faction:          string(p.Faction),
		ShortDescription: string(p.ShortDescription),
		FullDescription:  string(p.FullDescription),
		Price:            priceToProto(p.Price),
		ProductionTime:   p.ProductionTime,
		Space:            int32(p.Space),
		ImageUrl:         p.ImageURL,
		Attack:           int32(p.Attack),
		Defence:          int32(p.Defence),
		Capacity:         int32(p.Capacity),
		Stealth:          int32(p.Stealth),
		Speed:            int32(p.Speed),
	}
	if p.UnlockTechnologyID != nil {
		v := int64(*p.UnlockTechnologyID)
		out.UnlockTechnologyId = &v
	}
	return out
}

// ArmyPrototypeDomainToProto maps a domain army prototype to its wire shape (for command responses).
func ArmyPrototypeDomainToProto(p *domain.ArmyItemPrototype) *gamev1.ArmyPrototype {
	if p == nil {
		return nil
	}
	out := &gamev1.ArmyPrototype{
		Id:               int64(p.ID),
		Name:             string(p.Name),
		Category:         string(p.Category),
		CreationSources:  creationSourcesToStrings(p.CreationSources),
		Faction:          string(p.Faction),
		ShortDescription: string(p.ShortDescription),
		FullDescription:  string(p.FullDescription),
		Price: &gamev1.PriceModel{
			Credits:    int64(p.Price.Credits),
			Iron:       int64(p.Price.Iron),
			Titanium:   int64(p.Price.Titanium),
			Antimatter: int64(p.Price.Antimatter),
		},
		ProductionTime: p.ProductionTime,
		Space:          int32(p.Space),
		ImageUrl:       p.ImageURL,
		Attack:         int32(p.Attack),
		Defence:        int32(p.Defence),
		Capacity:       int32(p.Capacity),
		Stealth:        int32(p.Stealth),
		Speed:          int32(p.Speed),
	}
	if p.UnlockTechnologyID != nil {
		v := int64(*p.UnlockTechnologyID)
		out.UnlockTechnologyId = &v
	}
	return out
}

// ArmyPrototypeFromProto maps a wire army prototype to the domain aggregate (for
// commands), validating the input at the conversion boundary. Callers supply the
// id explicitly (prototypes are ordered and grouped into id ranges by type), so a
// missing/non-positive id is rejected. Validation failures are returned as plain
// English InvalidArgument status errors (see helpers.go for the rationale).
func ArmyPrototypeFromProto(p *gamev1.ArmyPrototype) (*domain.ArmyItemPrototype, error) {
	if p == nil {
		return nil, status.Error(codes.InvalidArgument, "prototype is required")
	}
	if p.Id <= 0 {
		return nil, status.Error(codes.InvalidArgument, "prototype id is required and must be positive")
	}
	if p.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "prototype name is required")
	}
	category := domain.ArmyCategory(p.Category)
	if err := validateArmyCategory(category); err != nil {
		return nil, err
	}
	faction := domain.Faction(p.Faction)
	if err := validateFaction(faction); err != nil {
		return nil, err
	}
	sources := creationSourcesFromStrings(p.CreationSources)
	if err := validateCreationSources(sources); err != nil {
		return nil, err
	}
	if hasNegative(p.ProductionTime, int64(p.Space), int64(p.Attack), int64(p.Defence), int64(p.Capacity), int64(p.Stealth), int64(p.Speed)) {
		return nil, status.Error(codes.InvalidArgument, "prototype numeric fields must not be negative")
	}
	price := priceFromProto(p.Price)
	if hasNegative(int64(price.Credits), int64(price.Iron), int64(price.Titanium), int64(price.Antimatter)) {
		return nil, status.Error(codes.InvalidArgument, "prototype price must not be negative")
	}

	out := &domain.ArmyItemPrototype{
		ID:               int(p.Id),
		Name:             p.Name,
		Category:         category,
		CreationSources:  sources,
		Faction:          faction,
		ShortDescription: p.ShortDescription,
		FullDescription:  p.FullDescription,
		Price:            price,
		ProductionTime:   p.ProductionTime,
		Space:            int(p.Space),
		ImageURL:         p.ImageUrl,
		Attack:           int(p.Attack),
		Defence:          int(p.Defence),
		Capacity:         int(p.Capacity),
		Stealth:          int(p.Stealth),
		Speed:            int(p.Speed),
	}
	if p.UnlockTechnologyId != nil {
		v := int(*p.UnlockTechnologyId)
		out.UnlockTechnologyID = &v
	}
	return out, nil
}

// validateArmyCategory checks c against the known army category set.
func validateArmyCategory(c domain.ArmyCategory) error {
	switch c {
	case domain.ArmyCategoryInfantry, domain.ArmyCategoryArmored, domain.ArmyCategoryArtillery,
		domain.ArmyCategoryAviation, domain.ArmyCategorySpy, domain.ArmyCategorySpecial:
		return nil
	default:
		return status.Errorf(codes.InvalidArgument, "invalid army category: %q", string(c))
	}
}
