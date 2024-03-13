package checkout

import (
	"errors"

	"github.com/Joshswooft/thinkmoney-test/currency"
	"github.com/Joshswooft/thinkmoney-test/quantity"
	"github.com/Joshswooft/thinkmoney-test/sku"
)

var (
	errNoPricingRulesProvided = errors.New("no pricing rules was provided")
)

type PricingRules interface {
	GetPrice(sku sku.SKU, quantity quantity.Quantity) currency.Pence
	PriceExists(sku sku.SKU) bool
}

func NewCheckout(pricingRules PricingRules) (*checkout, error) {

	if pricingRules == nil {
		return nil, errNoPricingRulesProvided
	}

	return &checkout{
		pricingRules: pricingRules,
	}, nil
}

type checkout struct {
	pricingRules PricingRules
}
