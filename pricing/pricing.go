package pricing

type sku = rune
type pence = int
type quantity = uint

type PricingData struct {
	UnitPrice       pence
	SpecialPrice    pence
	SpecialQuantity quantity
}

func (p *PricingData) HasSpecialOffer() bool {
	if p == nil {
		return false
	}
	return p.SpecialPrice != 0 || p.SpecialQuantity != 0
}

// exposing pricing config for convenience
type SpecialPricing struct {
	Config map[sku]PricingData
}

func (p *SpecialPricing) calculatePrice(data PricingData, quantity quantity) pence {
	if p == nil {
		return 0
	}

	if data.HasSpecialOffer() {

		specialQuantity := data.SpecialQuantity
		if specialQuantity <= 0 {
			specialQuantity = 1
		}

		// the number of items that have qualified for the special price
		bundleQuantity := quantity / specialQuantity

		// left over items that will have unit price applied
		remainingQuantity := quantity % specialQuantity

		bundlePrice := pence(bundleQuantity) * data.SpecialPrice
		regularPrice := pence(remainingQuantity) * data.UnitPrice

		return bundlePrice + pence(regularPrice)
	}

	return data.UnitPrice * pence(quantity)

}

func (p *SpecialPricing) GetPrice(sku sku, quantity quantity) pence {

	if p == nil || p.Config == nil {
		return 0
	}

	pricingData, exists := p.Config[sku]
	if !exists {
		return 0
	}

	return p.calculatePrice(pricingData, quantity)

}

func (p *SpecialPricing) PriceExists(sku sku) bool {
	if p == nil {
		return false
	}

	_, exists := p.Config[sku]

	return exists
}
