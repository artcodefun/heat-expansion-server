package domain

// PriceModel represents the price for items.
type PriceModel struct {
	Credits    int
	Iron       int
	Titanium   int
	Antimatter int
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
