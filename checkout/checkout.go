package checkout

import (
	"errors"

	"github.com/Joshswooft/thinkmoney-test/currency"
	"github.com/Joshswooft/thinkmoney-test/quantity"
	"github.com/Joshswooft/thinkmoney-test/sku"
)

var (
	errUnknownItemScanned     = errors.New("an unknown item was scanned")
	errNoPricingRulesProvided = errors.New("no pricing rules was provided")
	errNoBasketProvided       = errors.New("no basket was provided")
)

type PricingRules interface {
	GetPrice(sku sku.SKU, quantity quantity.Quantity) currency.Pence
	PriceExists(sku sku.SKU) bool
}

type Basket interface {
	// Adds a new item or updates the existing item's quantity by its sku
	AddItem(sku sku.SKU, quantity quantity.Quantity) error

	// Gets an item by it's product SKU
	// if the item is not found then it returns a checkout.ErrItemNotFound error
	GetItem(sku sku.SKU) (qty quantity.Quantity, err error)
}

func NewCheckout(pricingRules PricingRules, basket Basket) (*checkout, error) {

	if pricingRules == nil {
		return nil, errNoPricingRulesProvided
	}

	if basket == nil {
		return nil, errNoBasketProvided
	}

	return &checkout{
		pricingRules: pricingRules,
		basket:       basket,
	}, nil
}

type checkout struct {
	pricingRules PricingRules
	basket       Basket
}

func (c *checkout) doScan(sku sku.SKU, quantity quantity.Quantity) error {
	// maybe its better to just add the item to the basket?
	if exists := c.pricingRules.PriceExists(sku); !exists {
		return errUnknownItemScanned
	}

	itemQuantity, err := c.basket.GetItem(sku)

	if err != nil && errors.Is(err, ErrItemNotFound) {
		return c.basket.AddItem(sku, quantity)
	}

	updatedQuantity := itemQuantity
	updatedQuantity.Add(quantity.Value())

	return c.basket.AddItem(sku, updatedQuantity)
}

// scans in a single item into the basket
func (c *checkout) Scan(sku sku.SKU, quantity quantity.Quantity) error {
	if quantity.Value() == 0 {
		return nil
	}
	return c.doScan(sku, quantity)
}

func (c *checkout) GetTotalPrice() currency.Pence {
	totalPrice := currency.Pence(10)

	return totalPrice

}
