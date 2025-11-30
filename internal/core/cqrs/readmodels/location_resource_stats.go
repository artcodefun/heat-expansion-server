package readmodels

// LocationResourceStats represents the available resources at a non-user-base location
// (e.g., resource nodes or dangerous locations). This is similar in spirit to
// UserBaseStats but intentionally scoped down to just the available resource pool
// at that region, which can be looted or otherwise consumed by operations.
type LocationResourceStats struct {
	Credits    int
	Iron       int
	Titanium   int
	Antimatter int

	// Optional bookkeeping to support time-based accumulation if needed later.
	// For now, this lets us compute deltas similarly to how bases do, without
	// introducing production/capacity semantics until required.
	CalculationTimestamp int64 // Unix timestamp of last resource calculation
}
