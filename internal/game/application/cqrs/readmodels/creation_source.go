package readmodels

// CreationSource identifies where a prototype is allowed to originate from.
type CreationSource string

const (
	CreationSourcePlayerBase    CreationSource = "PLAYER_BASE"
	CreationSourceBlackMarket   CreationSource = "BLACK_MARKET"
	CreationSourceNPCLocation   CreationSource = "NPC_LOCATION"
	CreationSourceConsumableBox CreationSource = "CONSUMABLE_BOX"
)
