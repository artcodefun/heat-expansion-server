package dtos

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	gamev1 "github.com/artcodefun/heat-expansion-server/contracts/game/grpc/v1"
	"github.com/artcodefun/heat-expansion-server/internal/game/application/cqrs/readmodels"
	"github.com/artcodefun/heat-expansion-server/internal/game/domain"
)

// TechPrototypesToProto maps a slice of readmodel tech prototypes to their wire shapes.
func TechPrototypesToProto(protos []*readmodels.TechItemPrototype) []*gamev1.TechPrototype {
	out := make([]*gamev1.TechPrototype, len(protos))
	for i, p := range protos {
		out[i] = TechPrototypeToProto(p)
	}
	return out
}

// TechPrototypeToProto maps a readmodel tech prototype to its wire shape.
func TechPrototypeToProto(p *readmodels.TechItemPrototype) *gamev1.TechPrototype {
	if p == nil {
		return nil
	}
	out := &gamev1.TechPrototype{
		Id:               int64(p.ID),
		Name:             string(p.Name),
		Category:         string(p.Category),
		ShortDescription: string(p.ShortDescription),
		FullDescription:  string(p.FullDescription),
		Price:            priceToProto(p.Price),
		ResearchTime:     p.ResearchTime,
		ImageUrl:         p.ImageURL,
		Improvement:      techImprovementModelToProto(p.Improvement),
	}
	if p.UnlockTechnologyID != nil {
		out.UnlockTechnologyId = int64(*p.UnlockTechnologyID)
	}
	return out
}

// TechPrototypeDomainToProto maps a domain tech prototype to its wire shape (for command responses).
func TechPrototypeDomainToProto(p *domain.TechItemPrototype) *gamev1.TechPrototype {
	if p == nil {
		return nil
	}
	out := &gamev1.TechPrototype{
		Id:               int64(p.ID),
		Name:             string(p.Name),
		Category:         string(p.Category),
		ShortDescription: string(p.ShortDescription),
		FullDescription:  string(p.FullDescription),
		Price: &gamev1.PriceModel{
			Credits:    int64(p.Price.Credits),
			Iron:       int64(p.Price.Iron),
			Titanium:   int64(p.Price.Titanium),
			Antimatter: int64(p.Price.Antimatter),
		},
		ResearchTime: p.ResearchTime,
		ImageUrl:     p.ImageURL,
		Improvement:  techImprovementDomainToProto(p.Improvement),
	}
	if p.UnlockTechnologyID != nil {
		out.UnlockTechnologyId = int64(*p.UnlockTechnologyID)
	}
	return out
}

// TechPrototypeFromProto maps a wire tech prototype to the domain aggregate (for
// commands), validating the input at the conversion boundary. Validation failures
// are returned as plain English InvalidArgument status errors.
func TechPrototypeFromProto(p *gamev1.TechPrototype) (*domain.TechItemPrototype, error) {
	if p == nil {
		return nil, status.Error(codes.InvalidArgument, "prototype is required")
	}
	if p.Id <= 0 {
		return nil, status.Error(codes.InvalidArgument, "prototype id is required and must be positive")
	}
	if p.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "prototype name is required")
	}
	category := domain.TechCategory(p.Category)
	if err := validateTechCategory(category); err != nil {
		return nil, err
	}
	if p.ResearchTime < 0 {
		return nil, status.Error(codes.InvalidArgument, "research_time must not be negative")
	}
	if p.Price != nil && hasNegative(p.Price.Credits, p.Price.Iron, p.Price.Titanium, p.Price.Antimatter) {
		return nil, status.Error(codes.InvalidArgument, "price must not contain negative values")
	}

	var unlockID *int
	if p.UnlockTechnologyId > 0 {
		v := int(p.UnlockTechnologyId)
		unlockID = &v
	}

	improvement, err := techImprovementFromProto(p.Improvement)
	if err != nil {
		return nil, err
	}

	out := &domain.TechItemPrototype{
		ID:                 int(p.Id),
		Name:               p.Name,
		Category:           category,
		UnlockTechnologyID: unlockID,
		ShortDescription:   p.ShortDescription,
		FullDescription:    p.FullDescription,
		Price:              priceFromProto(p.Price),
		ResearchTime:       p.ResearchTime,
		ImageURL:           p.ImageUrl,
		Improvement:        improvement,
	}
	return out, nil
}

// validateTechCategory checks c against the known tech category set.
func validateTechCategory(c domain.TechCategory) error {
	switch c {
	case domain.TechCategoryArmy, domain.TechCategoryBuild, domain.TechCategoryBase, domain.TechCategoryPolitics:
		return nil
	default:
		return status.Errorf(codes.InvalidArgument, "invalid tech category: %q", string(c))
	}
}

// validateImprovementType checks t against the known improvement type set.
func validateImprovementType(t domain.ImprovementType) error {
	switch t {
	case domain.ImprovementTypeSpaceCapacity, domain.ImprovementTypeOperationsCount,
		domain.ImprovementTypeActiveBuffsCount, domain.ImprovementTypeActiveArtifactsCount,
		domain.ImprovementTypeActiveRestorationsCount, domain.ImprovementTypeBuildingProductionCount,
		domain.ImprovementTypeActiveDecryptionsCount:
		return nil
	default:
		return status.Errorf(codes.InvalidArgument, "invalid improvement type: %q", string(t))
	}
}

// ---- proto -> domain ----

func techImprovementFromProto(m *gamev1.TechImprovementModel) (*domain.TechImprovement, error) {
	if m == nil {
		return nil, nil
	}
	t := domain.ImprovementType(m.Type)
	if err := validateImprovementType(t); err != nil {
		return nil, err
	}
	if m.Value < 0 {
		return nil, status.Error(codes.InvalidArgument, "improvement value must not be negative")
	}
	var maxLevel *int
	if m.MaxLevel > 0 {
		v := int(m.MaxLevel)
		maxLevel = &v
	}
	return &domain.TechImprovement{
		Type:     t,
		Value:    int(m.Value),
		MaxLevel: maxLevel,
	}, nil
}

// ---- domain -> proto ----

func techImprovementDomainToProto(imp *domain.TechImprovement) *gamev1.TechImprovementModel {
	if imp == nil {
		return nil
	}
	m := &gamev1.TechImprovementModel{
		Type:  string(imp.Type),
		Value: int32(imp.Value),
	}
	if imp.MaxLevel != nil {
		m.MaxLevel = int32(*imp.MaxLevel)
	}
	return m
}

// ---- readmodel -> proto ----

func techImprovementModelToProto(imp *readmodels.TechImprovement) *gamev1.TechImprovementModel {
	if imp == nil {
		return nil
	}
	m := &gamev1.TechImprovementModel{
		Type:  string(imp.Type),
		Value: int32(imp.Value),
	}
	if imp.MaxLevel != nil {
		m.MaxLevel = int32(*imp.MaxLevel)
	}
	return m
}
