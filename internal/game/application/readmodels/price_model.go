package readmodels

// PriceModel represents the price for items.
type PriceModel struct {
	Credits    int
	Iron       int
	Titanium   int
	Antimatter int
}

func (p PriceModel) MultiplyFloat(f float64) PriceModel {
	return PriceModel{
		Credits:    int(float64(p.Credits) * f),
		Iron:       int(float64(p.Iron) * f),
		Titanium:   int(float64(p.Titanium) * f),
		Antimatter: int(float64(p.Antimatter) * f),
	}
}
