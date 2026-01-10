package domain

func newBaseWithDefaults(id int) *UserBaseModel {
	b := &UserBaseModel{ID: id}
	b.recalculateStats()
	// Override with some initial resources for tests
	b.Stats.Credits = 10_000
	b.Stats.Iron = 10_000
	b.Stats.Titanium = 10_000
	b.Stats.Antimatter = 10_000
	b.Stats.CalculationTimestamp = NowUnix()
	return b
}
