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
