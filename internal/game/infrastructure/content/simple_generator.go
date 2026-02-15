package content

import (
	"fmt"
	"strings"

	"github.com/artcodefun/heat-expansion-server/internal/game/application/ports"
	"github.com/artcodefun/heat-expansion-server/internal/game/domain"
)

// SimpleGenerator provides content for locations using hardcoded metadata.
type SimpleGenerator struct {
	staticBase string
	imageCount int
}

// NewSimpleGenerator creates a generator that uses hardcoded metadata for content.
func NewSimpleGenerator(staticBaseURL string) *SimpleGenerator {
	return &SimpleGenerator{
		staticBase: strings.TrimRight(staticBaseURL, "/"),
		imageCount: 25, // Assuming ~25 images per type as per user request
	}
}

func (g *SimpleGenerator) buildURL(folder, filename string, index int) string {
	if filename == "" {
		return ""
	}
	// Use absolute value for index to handle negative coordinates consistently
	if index < 0 {
		index = -index
	}
	// filename is expected to be without extension and index, e.g. "empty_sector"
	fullFilename := fmt.Sprintf("%s_%d.png", filename, index%g.imageCount)
	return fmt.Sprintf("%s/images/locations/%s/%s", g.staticBase, folder, fullFilename)
}

func (g *SimpleGenerator) getCoordIndex(coords domain.Vector2i) int {
	// Deterministic index from coordinates: (X * prime1 + Y)
	// Using 16381 as a prime to create a spread across sectors
	return coords.X*16381 + coords.Y
}

func (g *SimpleGenerator) GenerateEmptySectorContent(sector *domain.SectorModel) ports.GeneratedLocationContent {
	idx := 0
	if sector != nil {
		idx = g.getCoordIndex(sector.Coordinates)
	}
	return ports.GeneratedLocationContent{
		Name:        "location.empty_sector.name",
		Description: "location.empty_sector.description",
		ImageURL:    g.buildURL("empty_sectors", "empty_sector", idx),
	}
}

func (g *SimpleGenerator) GenerateBaseContent(base *domain.UserBaseModel) ports.GeneratedLocationContent {
	idx := 0
	if base != nil {
		idx = g.getCoordIndex(base.Coordinates)
	}
	return ports.GeneratedLocationContent{
		Name:        "location.user_base.name",
		Description: "location.user_base.description",
		ImageURL:    g.buildURL("user_bases", "user_base", idx),
	}
}

func (g *SimpleGenerator) GenerateResourceLocationContent(resource *domain.ResourceLocationModel) ports.GeneratedLocationContent {
	name := domain.TranslationKey("location.resource.default.name")
	desc := domain.TranslationKey("location.resource.default.description")
	fileBase := ""
	idx := 0

	if resource != nil {
		idx = g.getCoordIndex(resource.Coordinates)
		switch resource.Type {
		case domain.ResourceTypeIron:
			name = "location.resource.iron.name"
			desc = "location.resource.iron.description"
			fileBase = "resource_iron"
		case domain.ResourceTypeTitanium:
			name = "location.resource.titanium.name"
			desc = "location.resource.titanium.description"
			fileBase = "resource_titanium"
		case domain.ResourceTypeAntimatter:
			name = "location.resource.antimatter.name"
			desc = "location.resource.antimatter.description"
			fileBase = "resource_antimatter"
		case domain.ResourceTypeCredits:
			name = "location.resource.credits.name"
			desc = "location.resource.credits.description"
			fileBase = "resource_credits"
		}
	}

	return ports.GeneratedLocationContent{
		Name:        name,
		Description: desc,
		ImageURL:    g.buildURL("resource", fileBase, idx),
	}
}

func (g *SimpleGenerator) GenerateDangerousLocationContent(danger *domain.DangerousLocationModel) ports.GeneratedLocationContent {
	name := domain.TranslationKey("location.dangerous.default.name")
	desc := domain.TranslationKey("location.dangerous.default.description")
	fileBase := ""
	idx := 0

	if danger != nil {
		idx = g.getCoordIndex(danger.Coordinates)
		switch danger.DefenderFaction {
		case domain.FactionCustodianProtocol:
			name = "location.dangerous.vault.name"
			desc = "location.dangerous.vault.description"
			fileBase = "dangerous_vault"
		case domain.FactionNeuralWormApex:
			name = "location.dangerous.caverns.name"
			desc = "location.dangerous.caverns.description"
			fileBase = "dangerous_caverns"
		case domain.FactionObsidianSentinels:
			name = "location.dangerous.spires.name"
			desc = "location.dangerous.spires.description"
			fileBase = "dangerous_spires"
		case domain.FactionScorchWalkers:
			name = "location.dangerous.monolith.name"
			desc = "location.dangerous.monolith.description"
			fileBase = "dangerous_monolith"
		}
	}

	return ports.GeneratedLocationContent{
		Name:        name,
		Description: desc,
		ImageURL:    g.buildURL("dangerous", fileBase, idx),
	}
}
