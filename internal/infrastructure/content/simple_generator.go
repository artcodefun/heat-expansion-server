package content

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/artcodefun/heat-expansion-api/internal/core/domain"
	"github.com/artcodefun/heat-expansion-api/internal/core/ports"
)

const maxContentIndex = 1000

// SimpleGenerator provides deterministic placeholder content for locations.
type SimpleGenerator struct {
	contentDir string
	staticBase string
	rand       *rand.Rand
}

type placeholderData struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// NewSimpleGenerator creates a generator that knows how to build URLs backed by a local content directory.
func NewSimpleGenerator(contentDir, staticBaseURL string) *SimpleGenerator {
	g := &SimpleGenerator{
		contentDir: strings.TrimSpace(contentDir),
		staticBase: strings.TrimSpace(staticBaseURL),
		rand:       rand.New(rand.NewSource(time.Now().UnixNano())),
	}
	if g.contentDir != "" {
		_ = os.MkdirAll(g.contentDir, 0o755)
	}
	return g
}

func (g *SimpleGenerator) randomIndex() int {
	if g.rand == nil {
		g.rand = rand.New(rand.NewSource(time.Now().UnixNano()))
	}
	return g.rand.Intn(maxContentIndex + 1)
}

func (g *SimpleGenerator) assetURL(kind string, index int) string {
	rel := path.Join(kind, fmt.Sprintf("%d_image.png", index))
	g.ensureLocalPath(rel)
	return g.buildPublicURL(rel)
}

func (g *SimpleGenerator) ensureLocalPath(rel string) {
	if g.contentDir == "" {
		return
	}
	local := filepath.Join(g.contentDir, filepath.FromSlash(rel))
	if err := os.MkdirAll(filepath.Dir(local), 0o755); err != nil {
		return
	}
}

func (g *SimpleGenerator) metadata(kind string, index int) placeholderData {
	if g.contentDir == "" {
		return placeholderData{}
	}
	file := filepath.Join(g.contentDir, kind, fmt.Sprintf("%d_data.json", index))
	data, err := os.ReadFile(file)
	if err != nil {
		return placeholderData{}
	}
	var meta placeholderData
	if err := json.Unmarshal(data, &meta); err != nil {
		return placeholderData{}
	}
	return meta
}

func (g *SimpleGenerator) buildPublicURL(rel string) string {
	if g.staticBase == "" {
		return ""
	}
	rel = strings.TrimLeft(path.Clean("/"+rel), "/")
	return fmt.Sprintf("%s/%s", strings.TrimRight(g.staticBase, "/"), rel)
}

func (g *SimpleGenerator) generate(kind string, fallbackName, fallbackDescription string) ports.GeneratedLocationContent {
	idx := g.randomIndex()
	meta := g.metadata(kind, idx)
	name := fallbackName
	if meta.Name != "" {
		name = meta.Name
	}
	desc := fallbackDescription
	if meta.Description != "" {
		desc = meta.Description
	}
	return ports.GeneratedLocationContent{
		Name:        name,
		Description: desc,
		ImageURL:    g.assetURL(kind, idx),
	}
}

func (g *SimpleGenerator) GenerateEmptySectorContent(sector *domain.SectorModel) ports.GeneratedLocationContent {
	return g.generate("empty_sectors", "Uncharted Sector", "A quiet, empty stretch of space.")
}

func (g *SimpleGenerator) GenerateBaseContent(base *domain.UserBaseModel) ports.GeneratedLocationContent {
	fallbackName := "Unknown Base"
	if base != nil {
		fallbackName = fmt.Sprintf("Base #%d", base.UserID)
	}
	return g.generate("user_bases", fallbackName, "A sturdy outpost in the sector.")
}

func (g *SimpleGenerator) GenerateResourceLocationContent(resource *domain.ResourceLocationModel) ports.GeneratedLocationContent {
	return g.generate("resource_locs", "Resource Field", "Rich in valuable materials.")
}

func (g *SimpleGenerator) GenerateDangerousLocationContent(danger *domain.DangerousLocationModel) ports.GeneratedLocationContent {
	return g.generate("dangerous_locs", "Hazard Zone", "Hostile entities detected.")
}
