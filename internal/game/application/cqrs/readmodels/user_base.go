package readmodels

import (
	"github.com/google/uuid"
)

// UserBaseModel represents a military base in a sector.
type UserBaseModel struct {
	ID          int
	Coordinates Vector2i
	UserID      int
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

// BaseOwnedItem is embedded in all items that belong to a user base.
type BaseOwnedItem struct {
	ID         uuid.UUID
	UserBaseID int
}
