package pricing

type sku = rune
type pence = int
type quantity = uint

type PricingData struct {
	UnitPrice pence
}

// exposing pricing config for convenience
type SpecialPricing struct {
	Config map[sku]PricingData
}

func (p *SpecialPricing) GetPrice(sku sku, quantity quantity) pence {

	if p == nil || p.Config == nil {
		return 0
	}

	pricingData, exists := p.Config[sku]
	if !exists {
		return 0
	}

	return pricingData.UnitPrice

}
