package domain

import (
	"github.com/google/uuid"
)

// User represents a player in the game.
type User struct {
	EventProducer
	ID       uuid.UUID
	Name     string
	Crystals int // Global in-game currency for the user
}

// Default values for new users
const (
	DefaultCrystalsBalance = 50
)

// CrystalCreditReason identifies why crystals were added to a user's balance.
type CrystalCreditReason string

const (
	// CrystalCreditReasonSignupGrant is the starting balance granted on account creation.
	CrystalCreditReasonSignupGrant CrystalCreditReason = "SIGNUP_GRANT"
	// CrystalCreditReasonPackPurchase is a crystal pack bought with real money via billing.
	CrystalCreditReasonPackPurchase CrystalCreditReason = "PACK_PURCHASE"
)

// CrystalSpendReason identifies why crystals were deducted from a user's balance.
type CrystalSpendReason string

const (
	CrystalSpendReasonSpeedupBuilding     CrystalSpendReason = "SPEEDUP_BUILDING"
	CrystalSpendReasonSpeedupArmy         CrystalSpendReason = "SPEEDUP_ARMY"
	CrystalSpendReasonSpeedupTech         CrystalSpendReason = "SPEEDUP_TECH"
	CrystalSpendReasonSpeedupTrade        CrystalSpendReason = "SPEEDUP_TRADE"
	CrystalSpendReasonSpeedupMilitary     CrystalSpendReason = "SPEEDUP_MILITARY"
	CrystalSpendReasonBlackMarketResource CrystalSpendReason = "BLACK_MARKET_RESOURCE"
	CrystalSpendReasonBlackMarketBuilding CrystalSpendReason = "BLACK_MARKET_BUILDING"
	CrystalSpendReasonBlackMarketArmy     CrystalSpendReason = "BLACK_MARKET_ARMY"
	CrystalSpendReasonBlackMarketStorage  CrystalSpendReason = "BLACK_MARKET_STORAGE"
)

// NewUser creates a new user with default settings and a creation event.
func NewUser(id uuid.UUID, name string) *User {
	u := &User{
		ID:       id,
		Name:     name,
		Crystals: DefaultCrystalsBalance,
	}
	u.AddEvent(NewUserAccountCreatedEvent(u.ID))
	u.AddEvent(NewCrystalsCreditedEvent(u.ID, DefaultCrystalsBalance, CrystalCreditReasonSignupGrant, "", u.Crystals))
	return u
}

// SpendCrystals deducts the given amount from the user's crystal balance.
// It returns an error if the amount is non-positive or if the user does not
// have enough crystals available.
func (u *User) SpendCrystals(amount int, reason CrystalSpendReason, reference string) error {
	if amount <= 0 {
		return NewError("error.domain.user.invalid_crystal_spend_amount", H{"amount": amount})
	}
	if u.Crystals < amount {
		return NewError("error.domain.user.not_enough_crystals", nil)
	}
	u.Crystals -= amount
	u.AddEvent(NewCrystalsSpentEvent(u.ID, amount, reason, reference, u.Crystals))
	return nil
}

// AddCrystals credits the given amount to the user's crystal balance.
// It returns an error if amount is non-positive.
func (u *User) AddCrystals(amount int, reason CrystalCreditReason, reference string) error {
	if amount <= 0 {
		return NewError("error.domain.user.invalid_crystal_add_amount", H{"amount": amount})
	}
	u.Crystals += amount
	u.AddEvent(NewCrystalsCreditedEvent(u.ID, amount, reason, reference, u.Crystals))
	return nil
}
