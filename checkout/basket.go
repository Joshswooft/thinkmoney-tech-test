package checkout

import (
	"errors"
	"sync"

	"github.com/Joshswooft/thinkmoney-test/quantity"
	"github.com/Joshswooft/thinkmoney-test/sku"
)

// exposed this so other basket implementations can return the correct errors
var ErrItemNotFound = errors.New("item not found in the basket")

// was getting a compile error when using sku.SKU and quantity.Quantity as the direct map type
type itemID = sku.SKU
type qty = quantity.Quantity

// simple in-memory basket that is go-routine safe
func NewBasket() *basket {
	return &basket{items: make(map[itemID]quantity.Quantity)}
}

type basket struct {
	mu    sync.RWMutex
	items map[itemID]quantity.Quantity
}

// adds an item to the basket, operation is go-routine safe
// if the item already exists then the item is updated
func (b *basket) AddItem(sku sku.SKU, quantity quantity.Quantity) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.items == nil {
		b.items = make(map[itemID]qty)
	}

	b.items[sku] = quantity

	return nil
}

// GetItem from a basket by its sku, operation is go-routine safe
func (b *basket) GetItem(sku sku.SKU) (qty quantity.Quantity, err error) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	if b.items == nil {
		return qty, ErrItemNotFound
	}

	qty, found := b.items[sku]

	if !found {
		return *quantity.New(0), ErrItemNotFound
	}

	return qty, nil

}
