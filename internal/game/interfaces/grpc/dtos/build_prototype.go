package dtos

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	gamev1 "github.com/artcodefun/heat-expansion-server/contracts/game/grpc/v1"
	"github.com/artcodefun/heat-expansion-server/internal/game/application/cqrs/readmodels"
	"github.com/artcodefun/heat-expansion-server/internal/game/domain"
)

// BuildPrototypesToProto maps a slice of readmodel build prototypes to their wire shapes.
func BuildPrototypesToProto(protos []*readmodels.BuildItemPrototype) []*gamev1.BuildPrototype {
	out := make([]*gamev1.BuildPrototype, len(protos))
	for i, p := range protos {
		out[i] = BuildPrototypeToProto(p)
	}
	return out
}

// BuildPrototypeToProto maps a readmodel build prototype to its wire shape.
func BuildPrototypeToProto(p *readmodels.BuildItemPrototype) *gamev1.BuildPrototype {
	if p == nil {
		return nil
	}
	out := &gamev1.BuildPrototype{
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
	}
	if p.UnlockTechnologyID != nil {
		v := int64(*p.UnlockTechnologyID)
		out.UnlockTechnologyId = &v
	}
	setCategoryDataModel(out, p)
	return out
}

// BuildPrototypeDomainToProto maps a domain build prototype to its wire shape (for command responses).
func BuildPrototypeDomainToProto(p *domain.BuildItemPrototype) *gamev1.BuildPrototype {
	if p == nil {
		return nil
	}
	out := &gamev1.BuildPrototype{
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
	}
	if p.UnlockTechnologyID != nil {
		v := int64(*p.UnlockTechnologyID)
		out.UnlockTechnologyId = &v
	}
	setCategoryDataDomain(out, p)
	return out
}

// BuildPrototypeFromProto maps a wire build prototype to the domain aggregate (for
// commands), validating the input at the conversion boundary. Callers supply the
// id explicitly (prototypes are ordered and grouped into id ranges by type), so a
// missing/non-positive id is rejected. Validation failures are returned as plain
// English InvalidArgument status errors (see helpers.go for the rationale).
func BuildPrototypeFromProto(p *gamev1.BuildPrototype) (*domain.BuildItemPrototype, error) {
	if p == nil {
		return nil, status.Error(codes.InvalidArgument, "prototype is required")
	}
	if p.Id <= 0 {
		return nil, status.Error(codes.InvalidArgument, "prototype id is required and must be positive")
	}
	if p.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "prototype name is required")
	}
	category := domain.BuildCategory(p.Category)
	if err := validateBuildCategory(category); err != nil {
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
	if hasNegative(p.ProductionTime, int64(p.Space)) {
		return nil, status.Error(codes.InvalidArgument, "prototype numeric fields must not be negative")
	}
	price := priceFromProto(p.Price)
	if hasNegative(int64(price.Credits), int64(price.Iron), int64(price.Titanium), int64(price.Antimatter)) {
		return nil, status.Error(codes.InvalidArgument, "prototype price must not be negative")
	}

	control, resources, defense, military, intelligence, err := categoryDataFromProto(p)
	if err != nil {
		return nil, err
	}

	out := &domain.BuildItemPrototype{
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
		ControlData:      control,
		ResourcesData:    resources,
		DefenseData:      defense,
		MilitaryData:     military,
		IntelligenceData: intelligence,
	}
	if p.UnlockTechnologyId != nil {
		v := int(*p.UnlockTechnologyId)
		out.UnlockTechnologyID = &v
	}
	return out, nil
}

// validateBuildCategory checks c against the known build category set.
func validateBuildCategory(c domain.BuildCategory) error {
	switch c {
	case domain.BuildCategoryControl, domain.BuildCategoryResources, domain.BuildCategoryDefense,
		domain.BuildCategoryMilitary, domain.BuildCategoryIntelligence:
		return nil
	default:
		return status.Errorf(codes.InvalidArgument, "invalid build category: %q", string(c))
	}
}

// validateControlSubtype checks s against the known control-building subtype set.
func validateControlSubtype(s domain.ControlSubtype) error {
	switch s {
	case domain.ControlSubtypeRepairCenter, domain.ControlSubtypeCryptographyLab,
		domain.ControlSubtypeArtifactLab, domain.ControlSubtypeTradingTerminal,
		domain.ControlSubtypeMailingTerminal:
		return nil
	default:
		return status.Errorf(codes.InvalidArgument, "invalid control subtype: %q", string(s))
	}
}

// validateIntelligenceSubtype checks s against the known intelligence-building subtype set.
func validateIntelligenceSubtype(s domain.IntelligenceSubtype) error {
	switch s {
	case domain.IntelligenceSubtypeScanner, domain.IntelligenceSubtypeRadar,
		domain.IntelligenceSubtypeCloaking, domain.IntelligenceSubtypeScanInterceptor:
		return nil
	default:
		return status.Errorf(codes.InvalidArgument, "invalid intelligence subtype: %q", string(s))
	}
}

// ---- proto -> domain (validating) ----

// categoryDataFromProto unpacks the oneof and validates the active branch.
// It returns pointers for all five data types; exactly one will be non-nil
// (matching the active case), or all nil if no case is set.
// p.GetCategoryData() returns the unexported interface, so we accept *gamev1.BuildPrototype
// and call GetCategoryData() internally to keep the unexported type out of the signature.
func categoryDataFromProto(p *gamev1.BuildPrototype) (
	control *domain.ControlBuildingData,
	resources *domain.ResourcesBuildingData,
	defense *domain.DefenseBuildingData,
	military *domain.MilitaryBuildingData,
	intelligence *domain.IntelligenceBuildingData,
	err error,
) {
	switch v := p.GetCategoryData().(type) {
	case *gamev1.BuildPrototype_ControlData:
		control, err = controlDataFromProto(v.ControlData)
	case *gamev1.BuildPrototype_ResourcesData:
		resources, err = resourcesDataFromProto(v.ResourcesData)
	case *gamev1.BuildPrototype_DefenseData:
		defense, err = defenseDataFromProto(v.DefenseData)
	case *gamev1.BuildPrototype_MilitaryData:
		military, err = militaryDataFromProto(v.MilitaryData)
	case *gamev1.BuildPrototype_IntelligenceData:
		intelligence, err = intelligenceDataFromProto(v.IntelligenceData)
	}
	return
}

func controlDataFromProto(d *gamev1.BuildControlData) (*domain.ControlBuildingData, error) {
	if d == nil {
		return nil, nil
	}
	subtype := domain.ControlSubtype(d.Subtype)
	if err := validateControlSubtype(subtype); err != nil {
		return nil, err
	}
	return &domain.ControlBuildingData{Subtype: subtype}, nil
}

func resourcesDataFromProto(d *gamev1.BuildResourcesData) (*domain.ResourcesBuildingData, error) {
	if d == nil {
		return nil, nil
	}
	if hasNegativeFloat(d.CreditsProduction, d.IronProduction, d.TitaniumProduction, d.AntimatterProduction) {
		return nil, status.Error(codes.InvalidArgument, "resource production rates must not be negative")
	}
	if hasNegative(int64(d.CreditsCapacity), int64(d.IronCapacity), int64(d.TitaniumCapacity), int64(d.AntimatterCapacity)) {
		return nil, status.Error(codes.InvalidArgument, "resource capacities must not be negative")
	}
	return &domain.ResourcesBuildingData{
		CreditsProduction:    d.CreditsProduction,
		IronProduction:       d.IronProduction,
		TitaniumProduction:   d.TitaniumProduction,
		AntimatterProduction: d.AntimatterProduction,
		CreditsCapacity:      int(d.CreditsCapacity),
		IronCapacity:         int(d.IronCapacity),
		TitaniumCapacity:     int(d.TitaniumCapacity),
		AntimatterCapacity:   int(d.AntimatterCapacity),
	}, nil
}

func defenseDataFromProto(d *gamev1.BuildDefenseData) (*domain.DefenseBuildingData, error) {
	if d == nil {
		return nil, nil
	}
	if hasNegative(int64(d.DefenceBonus)) {
		return nil, status.Error(codes.InvalidArgument, "defence bonus must not be negative")
	}
	return &domain.DefenseBuildingData{DefenceBonus: int(d.DefenceBonus)}, nil
}

func militaryDataFromProto(d *gamev1.BuildMilitaryData) (*domain.MilitaryBuildingData, error) {
	if d == nil {
		return nil, nil
	}
	category := domain.ArmyCategory(d.UnlockArmyCategory)
	if err := validateArmyCategory(category); err != nil {
		return nil, err
	}
	return &domain.MilitaryBuildingData{UnlockArmyCategory: category}, nil
}

func intelligenceDataFromProto(d *gamev1.BuildIntelligenceData) (*domain.IntelligenceBuildingData, error) {
	if d == nil {
		return nil, nil
	}
	subtype := domain.IntelligenceSubtype(d.Subtype)
	if err := validateIntelligenceSubtype(subtype); err != nil {
		return nil, err
	}
	if hasNegative(int64(d.StealthStrength), int64(d.ScanRange), d.ScanCooldown) {
		return nil, status.Error(codes.InvalidArgument, "intelligence numeric fields must not be negative")
	}
	return &domain.IntelligenceBuildingData{
		Subtype:         subtype,
		StealthStrength: int(d.StealthStrength),
		ScanRange:       int(d.ScanRange),
		ScanCooldown:    d.ScanCooldown,
	}, nil
}

// ---- domain -> proto ----

// setCategoryDataDomain sets out.CategoryData to the oneof wrapper matching
// the active domain data block. The oneof interface is unexported so we mutate
// the proto struct directly rather than returning the interface type.
func setCategoryDataDomain(out *gamev1.BuildPrototype, p *domain.BuildItemPrototype) {
	switch {
	case p.ControlData != nil:
		out.CategoryData = &gamev1.BuildPrototype_ControlData{ControlData: controlDataDomainToProto(p.ControlData)}
	case p.ResourcesData != nil:
		out.CategoryData = &gamev1.BuildPrototype_ResourcesData{ResourcesData: resourcesDataDomainToProto(p.ResourcesData)}
	case p.DefenseData != nil:
		out.CategoryData = &gamev1.BuildPrototype_DefenseData{DefenseData: defenseDataDomainToProto(p.DefenseData)}
	case p.MilitaryData != nil:
		out.CategoryData = &gamev1.BuildPrototype_MilitaryData{MilitaryData: militaryDataDomainToProto(p.MilitaryData)}
	case p.IntelligenceData != nil:
		out.CategoryData = &gamev1.BuildPrototype_IntelligenceData{IntelligenceData: intelligenceDataDomainToProto(p.IntelligenceData)}
	}
}

func controlDataDomainToProto(d *domain.ControlBuildingData) *gamev1.BuildControlData {
	if d == nil {
		return nil
	}
	return &gamev1.BuildControlData{Subtype: string(d.Subtype)}
}

func resourcesDataDomainToProto(d *domain.ResourcesBuildingData) *gamev1.BuildResourcesData {
	if d == nil {
		return nil
	}
	return &gamev1.BuildResourcesData{
		CreditsProduction:    d.CreditsProduction,
		IronProduction:       d.IronProduction,
		TitaniumProduction:   d.TitaniumProduction,
		AntimatterProduction: d.AntimatterProduction,
		CreditsCapacity:      int32(d.CreditsCapacity),
		IronCapacity:         int32(d.IronCapacity),
		TitaniumCapacity:     int32(d.TitaniumCapacity),
		AntimatterCapacity:   int32(d.AntimatterCapacity),
	}
}

func defenseDataDomainToProto(d *domain.DefenseBuildingData) *gamev1.BuildDefenseData {
	if d == nil {
		return nil
	}
	return &gamev1.BuildDefenseData{DefenceBonus: int32(d.DefenceBonus)}
}

func militaryDataDomainToProto(d *domain.MilitaryBuildingData) *gamev1.BuildMilitaryData {
	if d == nil {
		return nil
	}
	return &gamev1.BuildMilitaryData{UnlockArmyCategory: string(d.UnlockArmyCategory)}
}

func intelligenceDataDomainToProto(d *domain.IntelligenceBuildingData) *gamev1.BuildIntelligenceData {
	if d == nil {
		return nil
	}
	return &gamev1.BuildIntelligenceData{
		Subtype:         string(d.Subtype),
		StealthStrength: int32(d.StealthStrength),
		ScanRange:       int32(d.ScanRange),
		ScanCooldown:    d.ScanCooldown,
	}
}

// ---- readmodel -> proto ----

// setCategoryDataModel sets out.CategoryData to the oneof wrapper matching
// the active readmodel data block. The oneof interface is unexported so we
// mutate the proto struct directly rather than returning the interface type.
func setCategoryDataModel(out *gamev1.BuildPrototype, p *readmodels.BuildItemPrototype) {
	switch {
	case p.ControlData != nil:
		out.CategoryData = &gamev1.BuildPrototype_ControlData{ControlData: controlDataModelToProto(p.ControlData)}
	case p.ResourcesData != nil:
		out.CategoryData = &gamev1.BuildPrototype_ResourcesData{ResourcesData: resourcesDataModelToProto(p.ResourcesData)}
	case p.DefenseData != nil:
		out.CategoryData = &gamev1.BuildPrototype_DefenseData{DefenseData: defenseDataModelToProto(p.DefenseData)}
	case p.MilitaryData != nil:
		out.CategoryData = &gamev1.BuildPrototype_MilitaryData{MilitaryData: militaryDataModelToProto(p.MilitaryData)}
	case p.IntelligenceData != nil:
		out.CategoryData = &gamev1.BuildPrototype_IntelligenceData{IntelligenceData: intelligenceDataModelToProto(p.IntelligenceData)}
	}
}

func controlDataModelToProto(d *readmodels.ControlBuildingData) *gamev1.BuildControlData {
	if d == nil {
		return nil
	}
	return &gamev1.BuildControlData{Subtype: string(d.Subtype)}
}

func resourcesDataModelToProto(d *readmodels.ResourcesBuildingData) *gamev1.BuildResourcesData {
	if d == nil {
		return nil
	}
	return &gamev1.BuildResourcesData{
		CreditsProduction:    d.CreditsProduction,
		IronProduction:       d.IronProduction,
		TitaniumProduction:   d.TitaniumProduction,
		AntimatterProduction: d.AntimatterProduction,
		CreditsCapacity:      int32(d.CreditsCapacity),
		IronCapacity:         int32(d.IronCapacity),
		TitaniumCapacity:     int32(d.TitaniumCapacity),
		AntimatterCapacity:   int32(d.AntimatterCapacity),
	}
}

func defenseDataModelToProto(d *readmodels.DefenseBuildingData) *gamev1.BuildDefenseData {
	if d == nil {
		return nil
	}
	return &gamev1.BuildDefenseData{DefenceBonus: int32(d.DefenceBonus)}
}

func militaryDataModelToProto(d *readmodels.MilitaryBuildingData) *gamev1.BuildMilitaryData {
	if d == nil {
		return nil
	}
	return &gamev1.BuildMilitaryData{UnlockArmyCategory: string(d.UnlockArmyCategory)}
}

func intelligenceDataModelToProto(d *readmodels.IntelligenceBuildingData) *gamev1.BuildIntelligenceData {
	if d == nil {
		return nil
	}
	return &gamev1.BuildIntelligenceData{
		Subtype:         string(d.Subtype),
		StealthStrength: int32(d.StealthStrength),
		ScanRange:       int32(d.ScanRange),
		ScanCooldown:    d.ScanCooldown,
	}
}
