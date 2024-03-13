package checkout

import (
	"errors"
	"fmt"
	"io"
	"log"

	"github.com/Joshswooft/thinkmoney-test/currency"
	"github.com/Joshswooft/thinkmoney-test/quantity"
	"github.com/Joshswooft/thinkmoney-test/sku"
)

var (
	errUnknownItemScanned     = errors.New("an unknown item was scanned")
	errNoPricingRulesProvided = errors.New("no pricing rules was provided")
	errNoScannerProvided      = errors.New("no scanner was provided")
	errNoBasketProvided       = errors.New("no basket was provided")
)

type Scanner interface {
	Scan() (sku.SKU, error)
}

// this interface could be a lot better but does the job for now
type Basket interface {
	// Adds a new item or updates the existing item's quantity by its sku
	AddItem(sku sku.SKU, quantity quantity.Quantity) error

	// Gets an item by it's product SKU
	// if the item is not found then it returns a checkout.ErrItemNotFound error
	GetItem(sku sku.SKU) (qty quantity.Quantity, err error)
	// runs the iterator func over every item in the basket
	Range(iterator func(id itemID, quantity quantity.Quantity))
}

type PricingRules interface {
	GetPrice(sku sku.SKU, quantity quantity.Quantity) currency.Pence
	PriceExists(sku sku.SKU) bool
}

type checkout struct {
	basket       Basket
	scanner      Scanner
	pricingRules PricingRules
}

func NewCheckout(pricingRules PricingRules, basket Basket, scanner Scanner) (*checkout, error) {
	if pricingRules == nil {
		return nil, errNoPricingRulesProvided
	}

	if scanner == nil {
		return nil, errNoScannerProvided
	}

	if basket == nil {
		return nil, errNoBasketProvided
	}

	return &checkout{
		pricingRules: pricingRules,
		basket:       basket,
		scanner:      scanner,
	}, nil
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

// reads in everything from the scanner and adds to the basket
// doesnt stop reading until it hits an io.EOF error
func (c *checkout) ScanItems() error {
	for {
		skuInstance, err := c.scanner.Scan()
		if err == io.EOF {
			break
		}

		if errors.Is(err, sku.ErrNoSpecialCharacters) {
			continue
		}

		if err != nil {
			// unknown error
			log.Println(fmt.Errorf("failed to read items from scanner, err=%w", err))
			// not sure if its better to continue or return an error here?
			// continue
			return err
		}

		if scanErr := c.doScan(skuInstance, *quantity.New(1)); scanErr != nil {
			log.Println(fmt.Errorf("failed to scan items into basket, err=%v, sku=%s", scanErr, skuInstance))
			// if the error is an unknown item then continue else setup retry logic...
			continue
		}
	}
	return nil
}

func (c *checkout) GetTotalPrice() currency.Pence {
	totalPrice := currency.Pence(0)

	adder := func(sku itemID, qty quantity.Quantity) {
		price := c.pricingRules.GetPrice(sku, qty)
		totalPrice += price
	}

	c.basket.Range(adder)

	return totalPrice

}
