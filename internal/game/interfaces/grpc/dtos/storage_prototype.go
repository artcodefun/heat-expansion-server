package dtos

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	gamev1 "github.com/artcodefun/heat-expansion-server/contracts/game/grpc/v1"
	"github.com/artcodefun/heat-expansion-server/internal/game/application/cqrs/readmodels"
	"github.com/artcodefun/heat-expansion-server/internal/game/domain"
)

// StoragePrototypesToProto maps a slice of readmodel storage prototypes to their wire shapes.
func StoragePrototypesToProto(protos []*readmodels.StorageItemPrototype) []*gamev1.StoragePrototype {
	out := make([]*gamev1.StoragePrototype, len(protos))
	for i, p := range protos {
		out[i] = StoragePrototypeToProto(p)
	}
	return out
}

// StoragePrototypeToProto maps a readmodel storage prototype to its wire shape.
func StoragePrototypeToProto(p *readmodels.StorageItemPrototype) *gamev1.StoragePrototype {
	if p == nil {
		return nil
	}
	out := &gamev1.StoragePrototype{
		Id:               int64(p.ID),
		Name:             string(p.Name),
		Category:         string(p.Category),
		CreationSources:  creationSourcesToStrings(p.CreationSources),
		EstimatedWorth:   int32(p.EstimatedWorth),
		ShortDescription: string(p.ShortDescription),
		FullDescription:  string(p.FullDescription),
		ImageUrl:         p.ImageURL,
	}
	setCategoryDataStorageModel(out, p)
	return out
}

// StoragePrototypeDomainToProto maps a domain storage prototype to its wire shape (for command responses).
func StoragePrototypeDomainToProto(p *domain.StorageItemPrototype) *gamev1.StoragePrototype {
	if p == nil {
		return nil
	}
	out := &gamev1.StoragePrototype{
		Id:               int64(p.ID),
		Name:             string(p.Name),
		Category:         string(p.Category),
		CreationSources:  creationSourcesToStrings(p.CreationSources),
		EstimatedWorth:   int32(p.EstimatedWorth),
		ShortDescription: string(p.ShortDescription),
		FullDescription:  string(p.FullDescription),
		ImageUrl:         p.ImageURL,
	}
	setCategoryDataStorageDomain(out, p)
	return out
}

// StoragePrototypeFromProto maps a wire storage prototype to the domain aggregate (for
// commands), validating the input at the conversion boundary. Callers supply the
// id explicitly (prototypes are ordered and grouped into id ranges by type), so a
// missing/non-positive id is rejected. Validation failures are returned as plain
// English InvalidArgument status errors.
func StoragePrototypeFromProto(p *gamev1.StoragePrototype) (*domain.StorageItemPrototype, error) {
	if p == nil {
		return nil, status.Error(codes.InvalidArgument, "prototype is required")
	}
	if p.Id <= 0 {
		return nil, status.Error(codes.InvalidArgument, "prototype id is required and must be positive")
	}
	if p.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "prototype name is required")
	}
	category := domain.StorageCategory(p.Category)
	if err := validateStorageCategory(category); err != nil {
		return nil, err
	}
	sources := creationSourcesFromStrings(p.CreationSources)
	if err := validateCreationSources(sources); err != nil {
		return nil, err
	}
	if p.EstimatedWorth < 0 {
		return nil, status.Error(codes.InvalidArgument, "estimated_worth must not be negative")
	}

	buff, intel, damaged, artifact, consumable, err := categoryDataStorageFromProto(p)
	if err != nil {
		return nil, err
	}

	out := &domain.StorageItemPrototype{
		ID:               int(p.Id),
		Name:             p.Name,
		Category:         category,
		CreationSources:  sources,
		EstimatedWorth:   int(p.EstimatedWorth),
		ShortDescription: p.ShortDescription,
		FullDescription:  p.FullDescription,
		ImageURL:         p.ImageUrl,
		BuffData:         buff,
		IntelData:        intel,
		DamagedData:      damaged,
		ArtifactData:     artifact,
		ConsumableData:   consumable,
	}
	return out, nil
}

// validateStorageCategory checks c against the known storage category set.
func validateStorageCategory(c domain.StorageCategory) error {
	switch c {
	case domain.StorageCategoryBuff, domain.StorageCategoryIntel, domain.StorageCategoryDamaged,
		domain.StorageCategoryArtifact, domain.StorageCategoryConsumable:
		return nil
	default:
		return status.Errorf(codes.InvalidArgument, "invalid storage category: %q", string(c))
	}
}

// validateBuffType checks t against the known buff type set.
func validateBuffType(t domain.BuffType) error {
	switch t {
	case domain.BuffTypeCreditsProduction, domain.BuffTypeIronProduction, domain.BuffTypeTitaniumProduction,
		domain.BuffTypeAttackIncrease, domain.BuffTypeDefenceIncrease, domain.BuffTypeStealthIncrease,
		domain.BuffTypeCapacityIncrease, domain.BuffTypeSpeedIncrease:
		return nil
	default:
		return status.Errorf(codes.InvalidArgument, "invalid buff type: %q", string(t))
	}
}

// validateHiddenLocationType checks t against the known hidden location type set.
func validateHiddenLocationType(t domain.HiddenLocationType) error {
	switch t {
	case domain.HiddenLocationTypeResourceful, domain.HiddenLocationTypeDangerous, domain.HiddenLocationTypeUserBase:
		return nil
	default:
		return status.Errorf(codes.InvalidArgument, "invalid hidden location type: %q", string(t))
	}
}

// validateArtifactEffectType checks t against the known artifact effect type set.
func validateArtifactEffectType(t domain.ArtifactEffectType) error {
	switch t {
	case domain.ArtifactEffectTypeCreditsProduction, domain.ArtifactEffectTypeIronProduction,
		domain.ArtifactEffectTypeTitaniumProduction, domain.ArtifactEffectTypeAttackIncrease,
		domain.ArtifactEffectTypeDefenceIncrease, domain.ArtifactEffectTypeStealthIncrease,
		domain.ArtifactEffectTypeCapacityIncrease, domain.ArtifactEffectTypeSpeedIncrease:
		return nil
	default:
		return status.Errorf(codes.InvalidArgument, "invalid artifact effect type: %q", string(t))
	}
}

// validateConsumableType checks t against the known consumable type set.
func validateConsumableType(t domain.ConsumableType) error {
	switch t {
	case domain.ConsumableTypeBox, domain.ConsumableTypeWarpCapsule:
		return nil
	default:
		return status.Errorf(codes.InvalidArgument, "invalid consumable type: %q", string(t))
	}
}

// validateConsumableBoxContents checks each box content entry against the known set.
func validateConsumableBoxContents(contents []domain.ConsumableBoxContents) error {
	for _, c := range contents {
		switch c {
		case domain.ConsumableContentsCredits, domain.ConsumableContentsIron, domain.ConsumableContentsTitanium,
			domain.ConsumableContentsAntimatter, domain.ConsumableContentsCrystals, domain.ConsumableContentsBuff,
			domain.ConsumableContentsIntel, domain.ConsumableContentsDamaged, domain.ConsumableContentsArtifact:
		default:
			return status.Errorf(codes.InvalidArgument, "invalid box contents value: %q", string(c))
		}
	}
	return nil
}

// ---- proto -> domain (validating) ----

// categoryDataStorageFromProto unpacks the oneof and validates the active branch.
func categoryDataStorageFromProto(p *gamev1.StoragePrototype) (
	buff *domain.BuffStorageData,
	intel *domain.IntelStorageData,
	damaged *domain.DamagedStorageData,
	artifact *domain.ArtifactStorageData,
	consumable *domain.ConsumableStorageData,
	err error,
) {
	switch v := p.GetCategoryData().(type) {
	case *gamev1.StoragePrototype_BuffData:
		buff, err = buffDataFromProto(v.BuffData)
	case *gamev1.StoragePrototype_IntelData:
		intel, err = intelDataFromProto(v.IntelData)
	case *gamev1.StoragePrototype_DamagedData:
		damaged, err = damagedDataFromProto(v.DamagedData)
	case *gamev1.StoragePrototype_ArtifactData:
		artifact, err = artifactDataFromProto(v.ArtifactData)
	case *gamev1.StoragePrototype_ConsumableData:
		consumable, err = consumableDataFromProto(v.ConsumableData)
	}
	return
}

func buffDataFromProto(d *gamev1.StorageBuffData) (*domain.BuffStorageData, error) {
	if d == nil {
		return nil, nil
	}
	t := domain.BuffType(d.Type)
	if err := validateBuffType(t); err != nil {
		return nil, err
	}
	if d.Value < 0 {
		return nil, status.Error(codes.InvalidArgument, "buff value must not be negative")
	}
	if d.DurationSeconds < 0 {
		return nil, status.Error(codes.InvalidArgument, "buff duration_seconds must not be negative")
	}
	return &domain.BuffStorageData{
		Type:            t,
		Value:           d.Value,
		DurationSeconds: d.DurationSeconds,
	}, nil
}

func intelDataFromProto(d *gamev1.StorageIntelData) (*domain.IntelStorageData, error) {
	if d == nil {
		return nil, nil
	}
	t := domain.HiddenLocationType(d.Type)
	if err := validateHiddenLocationType(t); err != nil {
		return nil, err
	}
	if d.DecryptionSeconds < 0 {
		return nil, status.Error(codes.InvalidArgument, "intel decryption_seconds must not be negative")
	}
	return &domain.IntelStorageData{
		Type:              t,
		DecryptionSeconds: d.DecryptionSeconds,
	}, nil
}

func damagedDataFromProto(d *gamev1.StorageDamagedData) (*domain.DamagedStorageData, error) {
	if d == nil {
		return nil, nil
	}
	price := priceFromProto(d.RestorePrice)
	if hasNegative(int64(price.Credits), int64(price.Iron), int64(price.Titanium), int64(price.Antimatter)) {
		return nil, status.Error(codes.InvalidArgument, "damaged restore_price must not be negative")
	}
	if d.RestorationSeconds < 0 {
		return nil, status.Error(codes.InvalidArgument, "damaged restoration_seconds must not be negative")
	}
	if d.OriginalUnitId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "damaged original_unit_id must be positive")
	}
	return &domain.DamagedStorageData{
		RestorePrice:       price,
		RestorationSeconds: d.RestorationSeconds,
		OriginalUnitID:     int(d.OriginalUnitId),
	}, nil
}

func artifactDataFromProto(d *gamev1.StorageArtifactData) (*domain.ArtifactStorageData, error) {
	if d == nil {
		return nil, nil
	}
	t := domain.ArtifactEffectType(d.Type)
	if err := validateArtifactEffectType(t); err != nil {
		return nil, err
	}
	if d.Value < 0 {
		return nil, status.Error(codes.InvalidArgument, "artifact value must not be negative")
	}
	return &domain.ArtifactStorageData{
		Type:  t,
		Value: d.Value,
	}, nil
}

func consumableDataFromProto(d *gamev1.StorageConsumableData) (*domain.ConsumableStorageData, error) {
	if d == nil {
		return nil, nil
	}
	t := domain.ConsumableType(d.Type)
	if err := validateConsumableType(t); err != nil {
		return nil, err
	}
	contents := make([]domain.ConsumableBoxContents, len(d.BoxContents))
	for i, c := range d.BoxContents {
		contents[i] = domain.ConsumableBoxContents(c)
	}
	if err := validateConsumableBoxContents(contents); err != nil {
		return nil, err
	}
	if d.BoxSize < 0 {
		return nil, status.Error(codes.InvalidArgument, "consumable box_size must not be negative")
	}
	return &domain.ConsumableStorageData{
		Type:        t,
		BoxContents: contents,
		BoxSize:     int(d.BoxSize),
	}, nil
}

// ---- domain -> proto ----

// setCategoryDataStorageDomain sets out.CategoryData to the oneof wrapper matching
// the active domain data block.
func setCategoryDataStorageDomain(out *gamev1.StoragePrototype, p *domain.StorageItemPrototype) {
	switch {
	case p.BuffData != nil:
		out.CategoryData = &gamev1.StoragePrototype_BuffData{BuffData: buffDataDomainToProto(p.BuffData)}
	case p.IntelData != nil:
		out.CategoryData = &gamev1.StoragePrototype_IntelData{IntelData: intelDataDomainToProto(p.IntelData)}
	case p.DamagedData != nil:
		out.CategoryData = &gamev1.StoragePrototype_DamagedData{DamagedData: damagedDataDomainToProto(p.DamagedData)}
	case p.ArtifactData != nil:
		out.CategoryData = &gamev1.StoragePrototype_ArtifactData{ArtifactData: artifactDataDomainToProto(p.ArtifactData)}
	case p.ConsumableData != nil:
		out.CategoryData = &gamev1.StoragePrototype_ConsumableData{ConsumableData: consumableDataDomainToProto(p.ConsumableData)}
	}
}

func buffDataDomainToProto(d *domain.BuffStorageData) *gamev1.StorageBuffData {
	if d == nil {
		return nil
	}
	return &gamev1.StorageBuffData{
		Type:            string(d.Type),
		Value:           d.Value,
		DurationSeconds: d.DurationSeconds,
	}
}

func intelDataDomainToProto(d *domain.IntelStorageData) *gamev1.StorageIntelData {
	if d == nil {
		return nil
	}
	return &gamev1.StorageIntelData{
		Type:              string(d.Type),
		DecryptionSeconds: d.DecryptionSeconds,
	}
}

func damagedDataDomainToProto(d *domain.DamagedStorageData) *gamev1.StorageDamagedData {
	if d == nil {
		return nil
	}
	return &gamev1.StorageDamagedData{
		RestorePrice: &gamev1.PriceModel{
			Credits:    int64(d.RestorePrice.Credits),
			Iron:       int64(d.RestorePrice.Iron),
			Titanium:   int64(d.RestorePrice.Titanium),
			Antimatter: int64(d.RestorePrice.Antimatter),
		},
		RestorationSeconds: d.RestorationSeconds,
		OriginalUnitId:     int64(d.OriginalUnitID),
	}
}

func artifactDataDomainToProto(d *domain.ArtifactStorageData) *gamev1.StorageArtifactData {
	if d == nil {
		return nil
	}
	return &gamev1.StorageArtifactData{
		Type:  string(d.Type),
		Value: d.Value,
	}
}

func consumableDataDomainToProto(d *domain.ConsumableStorageData) *gamev1.StorageConsumableData {
	if d == nil {
		return nil
	}
	contents := make([]string, len(d.BoxContents))
	for i, c := range d.BoxContents {
		contents[i] = string(c)
	}
	return &gamev1.StorageConsumableData{
		Type:        string(d.Type),
		BoxContents: contents,
		BoxSize:     int32(d.BoxSize),
	}
}

// ---- readmodel -> proto ----

// setCategoryDataStorageModel sets out.CategoryData to the oneof wrapper matching
// the active readmodel data block.
func setCategoryDataStorageModel(out *gamev1.StoragePrototype, p *readmodels.StorageItemPrototype) {
	switch {
	case p.BuffData != nil:
		out.CategoryData = &gamev1.StoragePrototype_BuffData{BuffData: buffDataModelToProto(p.BuffData)}
	case p.IntelData != nil:
		out.CategoryData = &gamev1.StoragePrototype_IntelData{IntelData: intelDataModelToProto(p.IntelData)}
	case p.DamagedData != nil:
		out.CategoryData = &gamev1.StoragePrototype_DamagedData{DamagedData: damagedDataModelToProto(p.DamagedData)}
	case p.ArtifactData != nil:
		out.CategoryData = &gamev1.StoragePrototype_ArtifactData{ArtifactData: artifactDataModelToProto(p.ArtifactData)}
	case p.ConsumableData != nil:
		out.CategoryData = &gamev1.StoragePrototype_ConsumableData{ConsumableData: consumableDataModelToProto(p.ConsumableData)}
	}
}

func buffDataModelToProto(d *readmodels.BuffStorageData) *gamev1.StorageBuffData {
	if d == nil {
		return nil
	}
	return &gamev1.StorageBuffData{
		Type:            string(d.Type),
		Value:           d.Value,
		DurationSeconds: d.DurationSeconds,
	}
}

func intelDataModelToProto(d *readmodels.IntelStorageData) *gamev1.StorageIntelData {
	if d == nil {
		return nil
	}
	return &gamev1.StorageIntelData{
		Type:              string(d.Type),
		DecryptionSeconds: d.DecryptionSeconds,
	}
}

func damagedDataModelToProto(d *readmodels.DamagedStorageData) *gamev1.StorageDamagedData {
	if d == nil {
		return nil
	}
	return &gamev1.StorageDamagedData{
		RestorePrice: &gamev1.PriceModel{
			Credits:    int64(d.RestorePrice.Credits),
			Iron:       int64(d.RestorePrice.Iron),
			Titanium:   int64(d.RestorePrice.Titanium),
			Antimatter: int64(d.RestorePrice.Antimatter),
		},
		RestorationSeconds: d.RestorationSeconds,
		OriginalUnitId:     int64(d.OriginalUnitID),
	}
}

func artifactDataModelToProto(d *readmodels.ArtifactStorageData) *gamev1.StorageArtifactData {
	if d == nil {
		return nil
	}
	return &gamev1.StorageArtifactData{
		Type:  string(d.Type),
		Value: d.Value,
	}
}

func consumableDataModelToProto(d *readmodels.ConsumableStorageData) *gamev1.StorageConsumableData {
	if d == nil {
		return nil
	}
	contents := make([]string, len(d.BoxContents))
	for i, c := range d.BoxContents {
		contents[i] = string(c)
	}
	return &gamev1.StorageConsumableData{
		Type:        string(d.Type),
		BoxContents: contents,
		BoxSize:     int32(d.BoxSize),
	}
}
