package ports

import "github.com/artcodefun/heat-expansion-server/internal/game/domain"

type GeneratedLocationContent struct {
	Name        domain.TranslationKey
	Description domain.TranslationKey
	ImageURL    string
}

// ContentGenerator defines the interface for generating all content for game assets at once.
type ContentGenerator interface {
	// Empty sector flavor
	GenerateEmptySectorContent(sector *domain.SectorModel) GeneratedLocationContent
	GenerateBaseContent(base *domain.UserBaseModel) GeneratedLocationContent
	GenerateResourceLocationContent(resource *domain.ResourceLocationModel) GeneratedLocationContent
	GenerateDangerousLocationContent(danger *domain.DangerousLocationModel) GeneratedLocationContent
}
