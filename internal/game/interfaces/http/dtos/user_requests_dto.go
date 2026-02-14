package dtos

// UserCrystalBalanceResponse represents the current crystal balance for the authenticated user.
type UserCrystalBalanceResponse struct {
	Crystals int `json:"crystals"`
}
