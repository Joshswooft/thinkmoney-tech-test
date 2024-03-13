package checkout

import (
	"github.com/Joshswooft/thinkmoney-test/currency"
	"github.com/Joshswooft/thinkmoney-test/quantity"
	"github.com/Joshswooft/thinkmoney-test/sku"
)

// MockBasketStorage is a mock implementation of Basket for testing purposes.
type MockBasketStorage struct {
	Items           map[sku.SKU]quantity.Quantity
	AddItemErr      error
	GetTotalPriceFn func() int
}

// AddItem simulates adding an item to the mock basket storage.
func (m *MockBasketStorage) AddItem(sku sku.SKU, quantity quantity.Quantity) error {
	if m.AddItemErr != nil {
		return m.AddItemErr
	}
	m.Items[sku] = quantity
	return nil
}

func (m *MockBasketStorage) GetItem(sku sku.SKU) (qty quantity.Quantity, err error) {
	qty, exists := m.Items[sku]
	if !exists {
		return *quantity.New(0), ErrItemNotFound
	}
	return qty, nil
}

func (m *MockBasketStorage) Range(iterator func(id itemID, quantity quantity.Quantity)) {
	for id, qty := range m.Items {
		iterator(id, qty)
	}
}

// MockPricingRules is a mock implementation of PricingRules for testing purposes.
type MockPricingRules struct {
	Prices map[sku.SKU]currency.Pence
}

// GetPrice returns the price for a given SKU and quantity.
func (m *MockPricingRules) GetPrice(sku sku.SKU, quantity quantity.Quantity) currency.Pence {
	price, ok := m.Prices[sku]
	if !ok {
		return 0
	}
	qtyValue := quantity.Value()
	return price * currency.Pence(qtyValue)
}

// PriceExists checks if a price exists for a given SKU.
func (m *MockPricingRules) PriceExists(sku sku.SKU) bool {
	_, exists := m.Prices[sku]
	return exists
}
