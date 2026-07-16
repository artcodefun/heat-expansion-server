package grpcclient

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	gamev1 "github.com/artcodefun/heat-expansion-server/contracts/game/grpc/v1"
	"github.com/artcodefun/heat-expansion-server/internal/admin/application/ports"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
)

// GameClient implements ports.GamePrivateClient by calling the game module's
// private gRPC API. It dials lazily: the connection is established on the first
// RPC call, which avoids races during the shared errgroup startup sequence.
type GameClient struct {
	army        gamev1.ArmyPrototypeServiceClient
	build       gamev1.BuildPrototypeServiceClient
	storage     gamev1.StoragePrototypeServiceClient
	tech        gamev1.TechPrototypeServiceClient
	translation gamev1.TranslationServiceClient
}

func NewGameClient(addr, key string) (*GameClient, error) {
	conn, err := grpc.NewClient(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
		grpc.WithUnaryInterceptor(keyInterceptor(key)),
	)
	if err != nil {
		return nil, err
	}
	return &GameClient{
		army:        gamev1.NewArmyPrototypeServiceClient(conn),
		build:       gamev1.NewBuildPrototypeServiceClient(conn),
		storage:     gamev1.NewStoragePrototypeServiceClient(conn),
		tech:        gamev1.NewTechPrototypeServiceClient(conn),
		translation: gamev1.NewTranslationServiceClient(conn),
	}, nil
}

// keyInterceptor returns a client unary interceptor that injects the static
// module key into every outgoing RPC as x-internal-key metadata.
func keyInterceptor(key string) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		ctx = metadata.AppendToOutgoingContext(ctx, "x-internal-key", key)
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

// ── Army ─────────────────────────────────────────────────────────────────────

func (c *GameClient) ListArmyPrototypes(ctx context.Context) ([]*ports.ArmyPrototype, error) {
	resp, err := c.army.ListArmyPrototypes(ctx, &gamev1.ListArmyPrototypesRequest{})
	if err != nil {
		return nil, grpcErrToSentinel(err)
	}
	return armyProtosFromProto(resp.Prototypes), nil
}

func (c *GameClient) GetArmyPrototype(ctx context.Context, id int64) (*ports.ArmyPrototype, error) {
	resp, err := c.army.GetArmyPrototype(ctx, &gamev1.GetArmyPrototypeRequest{Id: id})
	if err != nil {
		return nil, grpcErrToSentinel(err)
	}
	return armyProtoFromProto(resp.Prototype), nil
}

func (c *GameClient) CreateArmyPrototype(ctx context.Context, p *ports.ArmyPrototype) (*ports.ArmyPrototype, error) {
	resp, err := c.army.CreateArmyPrototype(ctx, &gamev1.CreateArmyPrototypeRequest{Prototype: armyProtoToProto(p)})
	if err != nil {
		return nil, grpcErrToSentinel(err)
	}
	return armyProtoFromProto(resp.Prototype), nil
}

func (c *GameClient) UpdateArmyPrototype(ctx context.Context, p *ports.ArmyPrototype) (*ports.ArmyPrototype, error) {
	resp, err := c.army.UpdateArmyPrototype(ctx, &gamev1.UpdateArmyPrototypeRequest{Prototype: armyProtoToProto(p)})
	if err != nil {
		return nil, grpcErrToSentinel(err)
	}
	return armyProtoFromProto(resp.Prototype), nil
}

// ── Build ─────────────────────────────────────────────────────────────────────

func (c *GameClient) ListBuildPrototypes(ctx context.Context) ([]*ports.BuildPrototype, error) {
	resp, err := c.build.ListBuildPrototypes(ctx, &gamev1.ListBuildPrototypesRequest{})
	if err != nil {
		return nil, grpcErrToSentinel(err)
	}
	return buildProtosFromProto(resp.Prototypes), nil
}

func (c *GameClient) GetBuildPrototype(ctx context.Context, id int64) (*ports.BuildPrototype, error) {
	resp, err := c.build.GetBuildPrototype(ctx, &gamev1.GetBuildPrototypeRequest{Id: id})
	if err != nil {
		return nil, grpcErrToSentinel(err)
	}
	return buildProtoFromProto(resp.Prototype), nil
}

func (c *GameClient) CreateBuildPrototype(ctx context.Context, p *ports.BuildPrototype) (*ports.BuildPrototype, error) {
	resp, err := c.build.CreateBuildPrototype(ctx, &gamev1.CreateBuildPrototypeRequest{Prototype: buildProtoToProto(p)})
	if err != nil {
		return nil, grpcErrToSentinel(err)
	}
	return buildProtoFromProto(resp.Prototype), nil
}

func (c *GameClient) UpdateBuildPrototype(ctx context.Context, p *ports.BuildPrototype) (*ports.BuildPrototype, error) {
	resp, err := c.build.UpdateBuildPrototype(ctx, &gamev1.UpdateBuildPrototypeRequest{Prototype: buildProtoToProto(p)})
	if err != nil {
		return nil, grpcErrToSentinel(err)
	}
	return buildProtoFromProto(resp.Prototype), nil
}

// ── Storage ───────────────────────────────────────────────────────────────────

func (c *GameClient) ListStoragePrototypes(ctx context.Context) ([]*ports.StoragePrototype, error) {
	resp, err := c.storage.ListStoragePrototypes(ctx, &gamev1.ListStoragePrototypesRequest{})
	if err != nil {
		return nil, grpcErrToSentinel(err)
	}
	return storageProtosFromProto(resp.Prototypes), nil
}

func (c *GameClient) GetStoragePrototype(ctx context.Context, id int64) (*ports.StoragePrototype, error) {
	resp, err := c.storage.GetStoragePrototype(ctx, &gamev1.GetStoragePrototypeRequest{Id: id})
	if err != nil {
		return nil, grpcErrToSentinel(err)
	}
	return storageProtoFromProto(resp.Prototype), nil
}

func (c *GameClient) CreateStoragePrototype(ctx context.Context, p *ports.StoragePrototype) (*ports.StoragePrototype, error) {
	resp, err := c.storage.CreateStoragePrototype(ctx, &gamev1.CreateStoragePrototypeRequest{Prototype: storageProtoToProto(p)})
	if err != nil {
		return nil, grpcErrToSentinel(err)
	}
	return storageProtoFromProto(resp.Prototype), nil
}

func (c *GameClient) UpdateStoragePrototype(ctx context.Context, p *ports.StoragePrototype) (*ports.StoragePrototype, error) {
	resp, err := c.storage.UpdateStoragePrototype(ctx, &gamev1.UpdateStoragePrototypeRequest{Prototype: storageProtoToProto(p)})
	if err != nil {
		return nil, grpcErrToSentinel(err)
	}
	return storageProtoFromProto(resp.Prototype), nil
}

// ── Tech ──────────────────────────────────────────────────────────────────────

func (c *GameClient) ListTechPrototypes(ctx context.Context) ([]*ports.TechPrototype, error) {
	resp, err := c.tech.ListTechPrototypes(ctx, &gamev1.ListTechPrototypesRequest{})
	if err != nil {
		return nil, grpcErrToSentinel(err)
	}
	return techProtosFromProto(resp.Prototypes), nil
}

func (c *GameClient) GetTechPrototype(ctx context.Context, id int64) (*ports.TechPrototype, error) {
	resp, err := c.tech.GetTechPrototype(ctx, &gamev1.GetTechPrototypeRequest{Id: id})
	if err != nil {
		return nil, grpcErrToSentinel(err)
	}
	return techProtoFromProto(resp.Prototype), nil
}

func (c *GameClient) CreateTechPrototype(ctx context.Context, p *ports.TechPrototype) (*ports.TechPrototype, error) {
	resp, err := c.tech.CreateTechPrototype(ctx, &gamev1.CreateTechPrototypeRequest{Prototype: techProtoToProto(p)})
	if err != nil {
		return nil, grpcErrToSentinel(err)
	}
	return techProtoFromProto(resp.Prototype), nil
}

func (c *GameClient) UpdateTechPrototype(ctx context.Context, p *ports.TechPrototype) (*ports.TechPrototype, error) {
	resp, err := c.tech.UpdateTechPrototype(ctx, &gamev1.UpdateTechPrototypeRequest{Prototype: techProtoToProto(p)})
	if err != nil {
		return nil, grpcErrToSentinel(err)
	}
	return techProtoFromProto(resp.Prototype), nil
}

// ── Translation ───────────────────────────────────────────────────────────────

func (c *GameClient) UpsertTranslation(ctx context.Context, locale, key, value string) (*ports.Translation, error) {
	resp, err := c.translation.UpsertTranslation(ctx, &gamev1.UpsertTranslationRequest{
		Entry: &gamev1.TranslationEntry{Key: key, Locale: locale, Value: value},
	})
	if err != nil {
		return nil, grpcErrToSentinel(err)
	}
	e := resp.GetEntry()
	return &ports.Translation{Key: e.Key, Locale: e.Locale, Value: e.Value}, nil
}

func (c *GameClient) ListTranslations(ctx context.Context) ([]*ports.Translation, error) {
	resp, err := c.translation.ListTranslations(ctx, &gamev1.ListTranslationsRequest{})
	if err != nil {
		return nil, grpcErrToSentinel(err)
	}
	out := make([]*ports.Translation, len(resp.Entries))
	for i, e := range resp.Entries {
		out[i] = &ports.Translation{Key: e.Key, Locale: e.Locale, Value: e.Value}
	}
	return out, nil
}

// ── Mapping helpers ───────────────────────────────────────────────────────────

func priceFromProto(p *gamev1.PriceModel) ports.PriceModel {
	if p == nil {
		return ports.PriceModel{}
	}
	return ports.PriceModel{Credits: p.Credits, Iron: p.Iron, Titanium: p.Titanium, Antimatter: p.Antimatter}
}

func priceToProto(p ports.PriceModel) *gamev1.PriceModel {
	return &gamev1.PriceModel{Credits: p.Credits, Iron: p.Iron, Titanium: p.Titanium, Antimatter: p.Antimatter}
}

// ── Army mappings ─────────────────────────────────────────────────────────────

func armyProtoFromProto(p *gamev1.ArmyPrototype) *ports.ArmyPrototype {
	if p == nil {
		return nil
	}
	m := &ports.ArmyPrototype{
		ID:               p.Id,
		Name:             p.Name,
		Category:         p.Category,
		CreationSources:  p.CreationSources,
		Faction:          p.Faction,
		ShortDescription: p.ShortDescription,
		FullDescription:  p.FullDescription,
		Price:            priceFromProto(p.Price),
		ProductionTime:   p.ProductionTime,
		Space:            p.Space,
		ImageURL:         p.ImageUrl,
		Attack:           p.Attack,
		Defence:          p.Defence,
		Capacity:         p.Capacity,
		Stealth:          p.Stealth,
		Speed:            p.Speed,
	}
	m.UnlockTechnologyID = p.UnlockTechnologyId
	return m
}

func armyProtosFromProto(ps []*gamev1.ArmyPrototype) []*ports.ArmyPrototype {
	out := make([]*ports.ArmyPrototype, len(ps))
	for i, p := range ps {
		out[i] = armyProtoFromProto(p)
	}
	return out
}

func armyProtoToProto(m *ports.ArmyPrototype) *gamev1.ArmyPrototype {
	p := &gamev1.ArmyPrototype{
		Id:               m.ID,
		Name:             m.Name,
		Category:         m.Category,
		CreationSources:  m.CreationSources,
		Faction:          m.Faction,
		ShortDescription: m.ShortDescription,
		FullDescription:  m.FullDescription,
		Price:            priceToProto(m.Price),
		ProductionTime:   m.ProductionTime,
		Space:            m.Space,
		ImageUrl:         m.ImageURL,
		Attack:           m.Attack,
		Defence:          m.Defence,
		Capacity:         m.Capacity,
		Stealth:          m.Stealth,
		Speed:            m.Speed,
	}
	p.UnlockTechnologyId = m.UnlockTechnologyID
	return p
}

// ── Build mappings ────────────────────────────────────────────────────────────

func buildProtoFromProto(p *gamev1.BuildPrototype) *ports.BuildPrototype {
	if p == nil {
		return nil
	}
	m := &ports.BuildPrototype{
		ID:               p.Id,
		Name:             p.Name,
		Category:         p.Category,
		CreationSources:  p.CreationSources,
		Faction:          p.Faction,
		ShortDescription: p.ShortDescription,
		FullDescription:  p.FullDescription,
		Price:            priceFromProto(p.Price),
		ProductionTime:   p.ProductionTime,
		Space:            p.Space,
		ImageURL:         p.ImageUrl,
	}
	m.UnlockTechnologyID = p.UnlockTechnologyId
	switch v := p.GetCategoryData().(type) {
	case *gamev1.BuildPrototype_ControlData:
		if v.ControlData != nil {
			m.ControlData = &ports.BuildControlData{Subtype: v.ControlData.Subtype}
		}
	case *gamev1.BuildPrototype_ResourcesData:
		if v.ResourcesData != nil {
			d := v.ResourcesData
			m.ResourcesData = &ports.BuildResourcesData{
				CreditsProduction:    d.CreditsProduction,
				IronProduction:       d.IronProduction,
				TitaniumProduction:   d.TitaniumProduction,
				AntimatterProduction: d.AntimatterProduction,
				CreditsCapacity:      d.CreditsCapacity,
				IronCapacity:         d.IronCapacity,
				TitaniumCapacity:     d.TitaniumCapacity,
				AntimatterCapacity:   d.AntimatterCapacity,
			}
		}
	case *gamev1.BuildPrototype_DefenseData:
		if v.DefenseData != nil {
			m.DefenseData = &ports.BuildDefenseData{DefenceBonus: v.DefenseData.DefenceBonus}
		}
	case *gamev1.BuildPrototype_MilitaryData:
		if v.MilitaryData != nil {
			m.MilitaryData = &ports.BuildMilitaryData{UnlockArmyCategory: v.MilitaryData.UnlockArmyCategory}
		}
	case *gamev1.BuildPrototype_IntelligenceData:
		if v.IntelligenceData != nil {
			d := v.IntelligenceData
			m.IntelligenceData = &ports.BuildIntelligenceData{
				Subtype:         d.Subtype,
				StealthStrength: d.StealthStrength,
				ScanRange:       d.ScanRange,
				ScanCooldown:    d.ScanCooldown,
			}
		}
	}
	return m
}

func buildProtosFromProto(ps []*gamev1.BuildPrototype) []*ports.BuildPrototype {
	out := make([]*ports.BuildPrototype, len(ps))
	for i, p := range ps {
		out[i] = buildProtoFromProto(p)
	}
	return out
}

func buildProtoToProto(m *ports.BuildPrototype) *gamev1.BuildPrototype {
	p := &gamev1.BuildPrototype{
		Id:               m.ID,
		Name:             m.Name,
		Category:         m.Category,
		CreationSources:  m.CreationSources,
		Faction:          m.Faction,
		ShortDescription: m.ShortDescription,
		FullDescription:  m.FullDescription,
		Price:            priceToProto(m.Price),
		ProductionTime:   m.ProductionTime,
		Space:            m.Space,
		ImageUrl:         m.ImageURL,
	}
	p.UnlockTechnologyId = m.UnlockTechnologyID
	switch {
	case m.ControlData != nil:
		p.CategoryData = &gamev1.BuildPrototype_ControlData{ControlData: &gamev1.BuildControlData{Subtype: m.ControlData.Subtype}}
	case m.ResourcesData != nil:
		d := m.ResourcesData
		p.CategoryData = &gamev1.BuildPrototype_ResourcesData{ResourcesData: &gamev1.BuildResourcesData{
			CreditsProduction:    d.CreditsProduction,
			IronProduction:       d.IronProduction,
			TitaniumProduction:   d.TitaniumProduction,
			AntimatterProduction: d.AntimatterProduction,
			CreditsCapacity:      d.CreditsCapacity,
			IronCapacity:         d.IronCapacity,
			TitaniumCapacity:     d.TitaniumCapacity,
			AntimatterCapacity:   d.AntimatterCapacity,
		}}
	case m.DefenseData != nil:
		p.CategoryData = &gamev1.BuildPrototype_DefenseData{DefenseData: &gamev1.BuildDefenseData{DefenceBonus: m.DefenseData.DefenceBonus}}
	case m.MilitaryData != nil:
		p.CategoryData = &gamev1.BuildPrototype_MilitaryData{MilitaryData: &gamev1.BuildMilitaryData{UnlockArmyCategory: m.MilitaryData.UnlockArmyCategory}}
	case m.IntelligenceData != nil:
		d := m.IntelligenceData
		p.CategoryData = &gamev1.BuildPrototype_IntelligenceData{IntelligenceData: &gamev1.BuildIntelligenceData{
			Subtype:         d.Subtype,
			StealthStrength: d.StealthStrength,
			ScanRange:       d.ScanRange,
			ScanCooldown:    d.ScanCooldown,
		}}
	}
	return p
}

// ── Storage mappings ──────────────────────────────────────────────────────────

func storageProtoFromProto(p *gamev1.StoragePrototype) *ports.StoragePrototype {
	if p == nil {
		return nil
	}
	m := &ports.StoragePrototype{
		ID:               p.Id,
		Name:             p.Name,
		Category:         p.Category,
		CreationSources:  p.CreationSources,
		EstimatedWorth:   p.EstimatedWorth,
		ShortDescription: p.ShortDescription,
		FullDescription:  p.FullDescription,
		ImageURL:         p.ImageUrl,
	}
	switch v := p.GetCategoryData().(type) {
	case *gamev1.StoragePrototype_BuffData:
		if v.BuffData != nil {
			m.BuffData = &ports.StorageBuffData{Type: v.BuffData.Type, Value: v.BuffData.Value, DurationSeconds: v.BuffData.DurationSeconds}
		}
	case *gamev1.StoragePrototype_IntelData:
		if v.IntelData != nil {
			m.IntelData = &ports.StorageIntelData{Type: v.IntelData.Type, DecryptionSeconds: v.IntelData.DecryptionSeconds}
		}
	case *gamev1.StoragePrototype_DamagedData:
		if v.DamagedData != nil {
			d := v.DamagedData
			m.DamagedData = &ports.StorageDamagedData{
				RestorePrice:       priceFromProto(d.RestorePrice),
				RestorationSeconds: d.RestorationSeconds,
				OriginalUnitID:     d.OriginalUnitId,
			}
		}
	case *gamev1.StoragePrototype_ArtifactData:
		if v.ArtifactData != nil {
			m.ArtifactData = &ports.StorageArtifactData{Type: v.ArtifactData.Type, Value: v.ArtifactData.Value}
		}
	case *gamev1.StoragePrototype_ConsumableData:
		if v.ConsumableData != nil {
			d := v.ConsumableData
			m.ConsumableData = &ports.StorageConsumableData{Type: d.Type, BoxContents: d.BoxContents, BoxSize: d.BoxSize}
		}
	}
	return m
}

func storageProtosFromProto(ps []*gamev1.StoragePrototype) []*ports.StoragePrototype {
	out := make([]*ports.StoragePrototype, len(ps))
	for i, p := range ps {
		out[i] = storageProtoFromProto(p)
	}
	return out
}

func storageProtoToProto(m *ports.StoragePrototype) *gamev1.StoragePrototype {
	p := &gamev1.StoragePrototype{
		Id:               m.ID,
		Name:             m.Name,
		Category:         m.Category,
		CreationSources:  m.CreationSources,
		EstimatedWorth:   m.EstimatedWorth,
		ShortDescription: m.ShortDescription,
		FullDescription:  m.FullDescription,
		ImageUrl:         m.ImageURL,
	}
	switch {
	case m.BuffData != nil:
		p.CategoryData = &gamev1.StoragePrototype_BuffData{BuffData: &gamev1.StorageBuffData{Type: m.BuffData.Type, Value: m.BuffData.Value, DurationSeconds: m.BuffData.DurationSeconds}}
	case m.IntelData != nil:
		p.CategoryData = &gamev1.StoragePrototype_IntelData{IntelData: &gamev1.StorageIntelData{Type: m.IntelData.Type, DecryptionSeconds: m.IntelData.DecryptionSeconds}}
	case m.DamagedData != nil:
		d := m.DamagedData
		p.CategoryData = &gamev1.StoragePrototype_DamagedData{DamagedData: &gamev1.StorageDamagedData{
			RestorePrice:       priceToProto(d.RestorePrice),
			RestorationSeconds: d.RestorationSeconds,
			OriginalUnitId:     d.OriginalUnitID,
		}}
	case m.ArtifactData != nil:
		p.CategoryData = &gamev1.StoragePrototype_ArtifactData{ArtifactData: &gamev1.StorageArtifactData{Type: m.ArtifactData.Type, Value: m.ArtifactData.Value}}
	case m.ConsumableData != nil:
		d := m.ConsumableData
		p.CategoryData = &gamev1.StoragePrototype_ConsumableData{ConsumableData: &gamev1.StorageConsumableData{Type: d.Type, BoxContents: d.BoxContents, BoxSize: d.BoxSize}}
	}
	return p
}

// ── Tech mappings ─────────────────────────────────────────────────────────────

func techProtoFromProto(p *gamev1.TechPrototype) *ports.TechPrototype {
	if p == nil {
		return nil
	}
	m := &ports.TechPrototype{
		ID:                 p.Id,
		Name:               p.Name,
		Category:           p.Category,
		UnlockTechnologyID: p.UnlockTechnologyId,
		ShortDescription:   p.ShortDescription,
		FullDescription:    p.FullDescription,
		Price:              priceFromProto(p.Price),
		ResearchTime:       p.ResearchTime,
		ImageURL:           p.ImageUrl,
	}
	if imp := p.Improvement; imp != nil {
		m.Improvement = &ports.TechImprovement{Type: imp.Type, Value: imp.Value, MaxLevel: imp.MaxLevel}
	}
	return m
}

func techProtosFromProto(ps []*gamev1.TechPrototype) []*ports.TechPrototype {
	out := make([]*ports.TechPrototype, len(ps))
	for i, p := range ps {
		out[i] = techProtoFromProto(p)
	}
	return out
}

func techProtoToProto(m *ports.TechPrototype) *gamev1.TechPrototype {
	p := &gamev1.TechPrototype{
		Id:                 m.ID,
		Name:               m.Name,
		Category:           m.Category,
		UnlockTechnologyId: m.UnlockTechnologyID,
		ShortDescription:   m.ShortDescription,
		FullDescription:    m.FullDescription,
		Price:              priceToProto(m.Price),
		ResearchTime:       m.ResearchTime,
		ImageUrl:           m.ImageURL,
	}
	if m.Improvement != nil {
		p.Improvement = &gamev1.TechImprovementModel{Type: m.Improvement.Type, Value: m.Improvement.Value, MaxLevel: m.Improvement.MaxLevel}
	}
	return p
}
