package readmodels

import (
	"github.com/google/uuid"
)

// UserBaseModel represents a military base in a sector.
type UserBaseModel struct {
	ID          int
	Coordinates Vector2i
	UserID      uuid.UUID
	LocationDetails
}

// UserBaseStats represents current properties of a base.
type UserBaseStats struct {
	Credits               float64
	CreditsCapacity       int
	CreditsProduction     float64
	Iron                  float64
	IronCapacity          int
	IronProduction        float64
	Titanium              float64
	TitaniumCapacity      int
	TitaniumProduction    float64
	Antimatter            float64
	AntimatterCapacity    int
	AntimatterProduction  float64
	Defence               int
	Attack                int
	Space                 int
	MaxSpace              int
	MaxOperations         int
	MaxActiveBuffs        int
	MaxActiveArtifacts    int
	MaxBuildingProduction int
	MaxActiveRestorations int // Numeric bonus (prev. DamagedRestorationBonus)
	MaxActiveDecryptions  int
	CalculationTimestamp  int64 // Unix timestamp of last resource calculation
}

// CurrentResources projects resource amounts forward from CalculationTimestamp to now,
// applying production rates and clamping to capacities.
func (s *UserBaseStats) CurrentResources(now int64) PriceModel {
	delta := now - s.CalculationTimestamp
	credits := s.Credits
	iron := s.Iron
	titanium := s.Titanium
	antimatter := s.Antimatter
	if delta > 0 {
		credits = min(credits+s.CreditsProduction*float64(delta), float64(s.CreditsCapacity))
		iron = min(iron+s.IronProduction*float64(delta), float64(s.IronCapacity))
		titanium = min(titanium+s.TitaniumProduction*float64(delta), float64(s.TitaniumCapacity))
		antimatter = min(antimatter+s.AntimatterProduction*float64(delta), float64(s.AntimatterCapacity))
	}
	return PriceModel{
		Credits:    int(credits),
		Iron:       int(iron),
		Titanium:   int(titanium),
		Antimatter: int(antimatter),
	}
}

// BaseOwnedItem is embedded in all items that belong to a user base.
type BaseOwnedItem struct {
	ID         uuid.UUID
	UserBaseID int
}
