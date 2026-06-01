package domain

// PriceModel represents the price for items.
type PriceModel struct {
	Credits    int
	Iron       int
	Titanium   int
	Antimatter int
}

func (p PriceModel) Multiply(n int) PriceModel {
	return PriceModel{
		Credits:    p.Credits * n,
		Iron:       p.Iron * n,
		Titanium:   p.Titanium * n,
		Antimatter: p.Antimatter * n,
	}
}

func (p PriceModel) MultiplyFloat(f float64) PriceModel {
	return PriceModel{
		Credits:    int(float64(p.Credits) * f),
		Iron:       int(float64(p.Iron) * f),
		Titanium:   int(float64(p.Titanium) * f),
		Antimatter: int(float64(p.Antimatter) * f),
	}
}

func (p PriceModel) Divide(n int) PriceModel {
	if n <= 0 {
		return p
	}
	return PriceModel{
		Credits:    p.Credits / n,
		Iron:       p.Iron / n,
		Titanium:   p.Titanium / n,
		Antimatter: p.Antimatter / n,
	}
}

func (p PriceModel) CreditsWorth() float64 {
	return float64(p.Credits)*WorthCredit +
		float64(p.Iron)*WorthIron +
		float64(p.Titanium)*WorthTitanium +
		float64(p.Antimatter)*WorthAntimatter
}
