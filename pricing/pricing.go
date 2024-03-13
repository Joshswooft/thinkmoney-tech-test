package pricing

import (
	"github.com/Joshswooft/thinkmoney-test/currency"
	"github.com/Joshswooft/thinkmoney-test/quantity"
	"github.com/Joshswooft/thinkmoney-test/sku"
)

type SimplePricing struct {
	UnitPrices map[sku.SKU]currency.Pence
}

func (p *SimplePricing) GetPrice(sku sku.SKU, quantity quantity.Quantity) currency.Pence {
	if p == nil {
		return 0
	}

	if price, exists := p.UnitPrices[sku]; exists {
		return price * currency.Pence(quantity.Value())
	}

	return 0

}

func (p *SimplePricing) PriceExists(sku sku.SKU) bool {
	if p == nil {
		return false
	}

	_, exists := p.UnitPrices[sku]

	return exists
}

// exposing this just for convenience - in reality we would use a factory or config to set the data
type PricingData struct {
	UnitPrice       currency.Pence
	SpecialPrice    currency.Pence
	SpecialQuantity quantity.Quantity
}

func (p *PricingData) HasSpecialOffer() bool {
	if p == nil {
		return false
	}
	return p.SpecialPrice != 0 || p.SpecialQuantity.Value() != 0
}

// exposing pricing config for convenience
type SpecialPricing struct {
	Config map[sku.SKU]PricingData
}

func (p *SpecialPricing) calculatePrice(data PricingData, quantity quantity.Quantity) currency.Pence {
	if p == nil {
		return 0
	}

	qty := quantity.Value()

	if data.HasSpecialOffer() {

		specialQuantity := data.SpecialQuantity.Value()
		if specialQuantity <= 0 {
			specialQuantity = 1
		}

		// the number of items that have qualified for the special price
		bundleQuantity := qty / specialQuantity

		// left over items that will have unit price applied
		remainingQuantity := qty % specialQuantity

		bundlePrice := currency.Pence(bundleQuantity) * data.SpecialPrice
		regularPrice := currency.Pence(remainingQuantity) * data.UnitPrice

		return bundlePrice + currency.Pence(regularPrice)
	}

	return data.UnitPrice * currency.Pence(quantity.Value())

}

func (p *SpecialPricing) GetPrice(sku sku.SKU, quantity quantity.Quantity) currency.Pence {

	if p == nil || p.Config == nil {
		return 0
	}

	pricingData, exists := p.Config[sku]
	if !exists {
		return 0
	}

	return p.calculatePrice(pricingData, quantity)

}

func (p *SpecialPricing) PriceExists(sku sku.SKU) bool {
	if p == nil {
		return false
	}

	_, exists := p.Config[sku]

	return exists
}
