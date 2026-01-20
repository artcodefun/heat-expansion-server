package domain

func newBaseWithDefaults(id int) *UserBaseModel {
	b := &UserBaseModel{ID: id}
	b.recalculateStats()
	// Override with some initial resources for tests (Lorentz ratios)
	b.Stats.Credits = 1000
	b.Stats.Iron = 250
	b.Stats.Titanium = 50
	b.Stats.Antimatter = 3
	b.Stats.CalculationTimestamp = NowUnix()
	return b
}
