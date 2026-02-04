package content

import (
	"fmt"
	"strings"

	"github.com/artcodefun/heat-expansion-api/internal/core/domain"
	"github.com/artcodefun/heat-expansion-api/internal/core/ports"
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
		Name:        "Empty Sector",
		Description: "A desolate volcanic wasteland with no detectable resources or inhabitants.",
		ImageURL:    g.buildURL("empty_sectors", "empty_sector", idx),
	}
}

func (g *SimpleGenerator) GenerateBaseContent(base *domain.UserBaseModel) ports.GeneratedLocationContent {
	idx := 0
	name := "Human Expedition Base"
	if base != nil {
		idx = g.getCoordIndex(base.Coordinates)
		name = fmt.Sprintf("Base #%d", base.UserID)
	}
	return ports.GeneratedLocationContent{
		Name:        name,
		Description: "A fortified expeditionary outpost established by the Exo-Coalition.",
		ImageURL:    g.buildURL("user_bases", "user_base", idx),
	}
}

func (g *SimpleGenerator) GenerateResourceLocationContent(resource *domain.ResourceLocationModel) ports.GeneratedLocationContent {
	name := "Resource Field"
	desc := "Rich in valuable materials."
	fileBase := ""
	idx := 0

	if resource != nil {
		idx = g.getCoordIndex(resource.Coordinates)
		switch resource.Type {
		case domain.ResourceTypeIron:
			name = "Iron Rich Location"
			desc = "Large-scale metallic formations indicate high concentrations of raw iron ore."
			fileBase = "resource_iron"
		case domain.ResourceTypeTitanium:
			name = "Titanium Rich Location"
			desc = "Volcanic vents surrounding jagged orange crystal clusters, a sign of rich titanium deposits."
			fileBase = "resource_titanium"
		case domain.ResourceTypeAntimatter:
			name = "Antimatter Rich Location"
			desc = "An unstable rift leaking pure antimatter, causing localized gravitational distortions."
			fileBase = "resource_antimatter"
		case domain.ResourceTypeCredits:
			name = "Credit Rich Location"
			desc = "A ramshackle outpost built over the ruins of an old merchant hub, likely containing hoarded credits."
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
	name := "Hazard Zone"
	desc := "Hostile entities detected."
	fileBase := ""
	idx := 0

	if danger != nil {
		idx = g.getCoordIndex(danger.Coordinates)
		switch danger.DefenderFaction {
		case domain.FactionCustodianProtocol:
			name = "Dangerous Vault"
			desc = "A monolithic precursor structure guarding the entrance to a long-forgotten vault."
			fileBase = "dangerous_vault"
		case domain.FactionNeuralWormApex:
			name = "Dangerous Caverns"
			desc = "Subterranean depths filled with bioluminescent neural filaments and ancient data archives."
			fileBase = "dangerous_caverns"
		case domain.FactionObsidianSentinels:
			name = "Dangerous Spires"
			desc = "A forest of floating obsidian spires that hum with a rhythmic, defensive energy."
			fileBase = "dangerous_spires"
		case domain.FactionScorchWalkers:
			name = "Dangerous Monolith"
			desc = "A massive, featureless black monolith that seems to absorb all light and nearby scans."
			fileBase = "dangerous_monolith"
		}
	}

	return ports.GeneratedLocationContent{
		Name:        name,
		Description: desc,
		ImageURL:    g.buildURL("dangerous", fileBase, idx),
	}
}
