package services

import (
	"github.com/artcodefun/heat-expansion-api/internal/game/application/ports"
	"github.com/artcodefun/heat-expansion-api/internal/game/domain"
)

// SectorProvisioningService centralizes lazy creation of sectors with generated empty flavor content.
// Use from application layer when a sector must exist.
type SectorProvisioningService struct {
	ContentGenerator ports.ContentGenerator
}

func NewSectorProvisioningService(content ports.ContentGenerator) *SectorProvisioningService {
	return &SectorProvisioningService{ContentGenerator: content}
}

// EnsureSectorExists locks an existing sector by coordinates or creates one with generated empty details.
// sRepo must be transaction-scoped (obtained via SectorRepo.Tx(tx)); FOR UPDATE requires a transaction.
func (s *SectorProvisioningService) EnsureSectorExists(sRepo ports.SectorRepository, x, y int) (*domain.SectorModel, error) {
	// Try to lock existing sector first to serialize concurrent creators.
	sector, err := sRepo.FindByCoordinatesForUpdate(x, y)
	if err == ports.ErrNotFound {
		sector = &domain.SectorModel{Coordinates: domain.Vector2i{X: x, Y: y}}
		gen := s.ContentGenerator.GenerateEmptySectorContent(sector)
		sector.Details = domain.LocationDetails{Name: gen.Name, Description: gen.Description, ImageURL: gen.ImageURL}
		if err = sRepo.Create(sector); err != nil {
			// In case another tx inserted concurrently, re-read with lock and return.
			if locked, rerr := sRepo.FindByCoordinatesForUpdate(x, y); rerr == nil {
				return locked, nil
			}
			return nil, repoErr(err)
		}
		return sector, nil
	}
	if err != nil {
		return nil, repoErr(err)
	}
	return sector, nil
}

// CreateResourceLocationIfEmpty ensures sector exists and creates a resource location if the sector is empty.
// Expects repositories to be transaction-scoped when called inside a transaction.
func (s *SectorProvisioningService) CreateResourceLocationIfEmpty(sRepo ports.SectorRepository, rRepo ports.ResourceLocationRepository, loc *domain.ResourceLocationModel) error {
	if loc == nil {
		return nil
	}
	if _, err := s.EnsureSectorExists(sRepo, loc.Coordinates.X, loc.Coordinates.Y); err != nil {
		return err
	}
	if _, err := sRepo.FindByCoordinatesForUpdate(loc.Coordinates.X, loc.Coordinates.Y); err != nil {
		return repoErr(err)
	}
	lt, err := sRepo.GetLocationTypeByCoordinates(loc.Coordinates.X, loc.Coordinates.Y)
	if err != nil {
		return repoErr(err)
	}
	if lt != domain.LocationTypeEmpty {
		return nil
	}
	gen := s.ContentGenerator.GenerateResourceLocationContent(loc)
	loc.LocationDetails = domain.LocationDetails{Name: gen.Name, Description: gen.Description, ImageURL: gen.ImageURL}
	if err := rRepo.Create(loc); err != nil {
		return repoErr(err)
	}
	return nil
}

// CreateDangerousLocationIfEmpty ensures sector exists and creates a dangerous location if the sector is empty.
// Expects repositories to be transaction-scoped when called inside a transaction.
func (s *SectorProvisioningService) CreateDangerousLocationIfEmpty(sRepo ports.SectorRepository, dRepo ports.DangerousLocationRepository, loc *domain.DangerousLocationModel) error {
	if loc == nil {
		return nil
	}
	if _, err := s.EnsureSectorExists(sRepo, loc.Coordinates.X, loc.Coordinates.Y); err != nil {
		return err
	}
	if _, err := sRepo.FindByCoordinatesForUpdate(loc.Coordinates.X, loc.Coordinates.Y); err != nil {
		return repoErr(err)
	}
	lt, err := sRepo.GetLocationTypeByCoordinates(loc.Coordinates.X, loc.Coordinates.Y)
	if err != nil {
		return repoErr(err)
	}
	if lt != domain.LocationTypeEmpty {
		return nil
	}
	gen := s.ContentGenerator.GenerateDangerousLocationContent(loc)
	loc.LocationDetails = domain.LocationDetails{Name: gen.Name, Description: gen.Description, ImageURL: gen.ImageURL}
	if err := dRepo.Create(loc); err != nil {
		return repoErr(err)
	}
	return nil
}

// CreateUserBaseIfEmpty ensures the sector exists, locks it, re-checks emptiness, and creates the provided base.
// Returns true if the base was created, false if the sector was not empty.
// Expects repositories to be transaction-scoped when called inside a transaction.
func (s *SectorProvisioningService) CreateUserBaseIfEmpty(sRepo ports.SectorRepository, bRepo ports.UserBaseRepository, base *domain.UserBaseModel) (bool, error) {
	if base == nil {
		return false, nil
	}
	x, y := base.Coordinates.X, base.Coordinates.Y
	if _, err := s.EnsureSectorExists(sRepo, x, y); err != nil {
		return false, err
	}
	if _, err := sRepo.FindByCoordinatesForUpdate(x, y); err != nil {
		return false, repoErr(err)
	}
	lt, err := sRepo.GetLocationTypeByCoordinates(x, y)
	if err != nil {
		return false, repoErr(err)
	}
	if lt != domain.LocationTypeEmpty {
		return false, nil
	}
	content := s.ContentGenerator.GenerateBaseContent(base)
	base.Name = content.Name
	base.Description = content.Description
	base.ImageURL = content.ImageURL
	if err := bRepo.Create(base); err != nil {
		return false, repoErr(err)
	}
	return true, nil
}
