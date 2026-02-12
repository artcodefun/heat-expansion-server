package readmodels

// Faction represents the origin of an army unit or location defenders.
type Faction string

const (
	FactionExoCoalition      Faction = "EXO_COALITION"   // Playable (Human)
	FactionMarauders         Faction = "MARAUDERS"       // NPC: Credits
	FactionFerrousSwarm      Faction = "FERROUS_SWARM"   // NPC: Iron
	FactionTitanArachnids    Faction = "TITAN_ARACHNIDS" // NPC: Titanium
	FactionVoidEcho          Faction = "VOID_ECHO"       // NPC: Antimatter
	FactionCustodianProtocol Faction = "CUSTODIAN"       // NPC: Dangerous (Artifacts)
	FactionScorchWalkers     Faction = "SCORCH_WALKERS"  // NPC: Dangerous (Buffs)
	FactionObsidianSentinels Faction = "OBSIDIAN"        // NPC: Dangerous (Trophies)
	FactionNeuralWormApex    Faction = "NEURAL_WORM"     // NPC: Dangerous (Intel)
)
