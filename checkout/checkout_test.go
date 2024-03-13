package checkout

import (
	"errors"
	"reflect"

	"testing"

	"github.com/Joshswooft/thinkmoney-test/currency"
	"github.com/Joshswooft/thinkmoney-test/pricing"
	"github.com/Joshswooft/thinkmoney-test/quantity"
	"github.com/Joshswooft/thinkmoney-test/sku"
)

func skuGenerator(t *testing.T, r rune) sku.SKU {
	s, err := sku.New(r)
	if err != nil {
		t.Errorf("failed to make sku, input: %c, err: %v", r, err)
	}
	return s
}

func TestNewCheckout(t *testing.T) {

	scanner := &skuScanner{}
	basket := &basket{}
	pricingRules := &pricing.SimplePricing{}

	type args struct {
		pricingRules PricingRules
		basket       Basket
		scanner      Scanner
	}
	tests := []struct {
		name string
		args args
		err  error
	}{
		{
			name: "returns error when given no pricing rules",
			args: args{scanner: scanner, basket: basket},
			err:  errNoPricingRulesProvided,
		},
		{
			name: "returns error when given no scanner",
			args: args{pricingRules: pricingRules, basket: basket},
			err:  errNoScannerProvided,
		},
		{
			name: "returns error when given no basket",
			args: args{scanner: scanner, pricingRules: pricingRules},
			err:  errNoBasketProvided,
		},
		{
			name: "happy path",
			args: args{scanner: scanner, pricingRules: pricingRules, basket: basket},
			err:  nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewCheckout(tt.args.pricingRules, tt.args.basket, tt.args.scanner)
			if err != tt.err {
				t.Errorf("NewCheckout() error = %v, wantErr %v", err, tt.err)
				return
			}
		})
	}
}

func Test_checkout_Scan(t *testing.T) {

	skuA, err := sku.New('A')
	if err != nil {
		t.Errorf("failed to make sku, err: %v", err)
	}

	type fields struct {
		basket       Basket
		scanner      Scanner
		pricingRules PricingRules
	}
	type args struct {
		sku      sku.SKU
		quantity quantity.Quantity
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		qty     quantity.Quantity
	}{
		{
			name: "scans in a new item with quantity = 1",
			fields: fields{
				basket:       NewBasket(),
				scanner:      &skuScanner{},
				pricingRules: &pricing.SimplePricing{UnitPrices: map[sku.SKU]currency.Pence{skuA: 10}},
			},
			args:    args{sku: skuA, quantity: *quantity.New(1)},
			wantErr: false,
			qty:     *quantity.New(1),
		},
		{
			name:    "fails to scan an item where the pricing rules dont exist for that sku",
			fields:  fields{basket: NewBasket(), scanner: &skuScanner{}, pricingRules: &pricing.SimplePricing{}},
			args:    args{sku: skuA, quantity: *quantity.New(1)},
			wantErr: true,
		},
		{
			name: "updates an existing items quantity",
			fields: fields{
				basket:       &basket{items: map[sku.SKU]quantity.Quantity{skuA: *quantity.New(4)}},
				pricingRules: &pricing.SimplePricing{UnitPrices: map[sku.SKU]currency.Pence{skuA: 10}},
				scanner:      &skuScanner{},
			},
			args:    args{sku: skuA, quantity: *quantity.New(10)},
			wantErr: false,
			qty:     *quantity.New(14),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &checkout{
				basket:       tt.fields.basket,
				scanner:      tt.fields.scanner,
				pricingRules: tt.fields.pricingRules,
			}
			if err := c.Scan(tt.args.sku, tt.args.quantity); (err != nil) != tt.wantErr {
				t.Errorf("checkout.Scan() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				gotQty, _ := c.basket.GetItem(tt.args.sku)

				if gotQty != tt.qty {
					t.Errorf("basket quantity not expected, got: %d, expected: %d", gotQty, tt.qty)
				}

			}
		})
	}
}

func Test_checkout_ScanItems(t *testing.T) {

	errScanItem := errors.New("error scanning item")

	skuGen := func(r rune) sku.SKU {
		return skuGenerator(t, r)
	}

	tests := []struct {
		name           string
		pricingRules   PricingRules
		scanner        Scanner
		basketStorage  Basket
		expectedErr    error
		expectedBasket map[sku.SKU]quantity.Quantity
	}{
		{
			name: "Scan all items successfully into an empty basket",
			pricingRules: &MockPricingRules{
				Prices: map[sku.SKU]currency.Pence{skuGen('C'): 10, skuGen('A'): 5, skuGen('B'): 2},
			},
			scanner: &MockScanner{
				Items: []sku.SKU{skuGen('C'), skuGen('A'), skuGen('B')},
			},
			basketStorage: &MockBasketStorage{
				Items: map[sku.SKU]quantity.Quantity{},
			},
			expectedErr:    nil,
			expectedBasket: map[sku.SKU]quantity.Quantity{skuGen('A'): *quantity.New(1), skuGen('B'): *quantity.New(1), skuGen('C'): *quantity.New(1)},
		},
		{
			name:         "Error while scanning items",
			pricingRules: &MockPricingRules{Prices: map[sku.SKU]currency.Pence{skuGen('C'): 30}},
			scanner: &MockScanner{
				Err:   errScanItem,
				Items: []sku.SKU{skuGen('C')},
			},
			basketStorage: &MockBasketStorage{
				Items: map[sku.SKU]quantity.Quantity{},
			},
			expectedErr:    errScanItem,
			expectedBasket: map[sku.SKU]quantity.Quantity{},
		},
		{
			name: "fails to scan in items without pricing rules",
			pricingRules: &MockPricingRules{
				Prices: map[sku.SKU]currency.Pence{skuGen('A'): 40},
			},
			scanner: &MockScanner{
				Items: []sku.SKU{skuGen('C'), skuGen('A')},
			},
			basketStorage: &MockBasketStorage{
				Items: map[sku.SKU]quantity.Quantity{},
			},
			expectedErr:    nil,
			expectedBasket: map[sku.SKU]quantity.Quantity{skuGen('A'): *quantity.New(1)},
		},
		{
			name: "scans in items that appear multiple times",
			pricingRules: &MockPricingRules{
				Prices: map[sku.SKU]currency.Pence{skuGen('A'): 40},
			},
			scanner: &MockScanner{
				Items: []sku.SKU{skuGen('A'), skuGen('A')},
			},
			basketStorage: &MockBasketStorage{
				Items: map[sku.SKU]quantity.Quantity{},
			},
			expectedErr:    nil,
			expectedBasket: map[sku.SKU]quantity.Quantity{skuGen('A'): *quantity.New(2)},
		},
		{
			name: "scans in item to existing basket",
			pricingRules: &MockPricingRules{
				Prices: map[sku.SKU]currency.Pence{skuGen('A'): 40},
			},
			scanner: &MockScanner{
				Items: []sku.SKU{skuGen('A'), skuGen('A')},
			},
			basketStorage: &MockBasketStorage{
				Items: map[sku.SKU]quantity.Quantity{skuGen('A'): *quantity.New(5), skuGen('B'): *quantity.New(10)},
			},
			expectedErr:    nil,
			expectedBasket: map[sku.SKU]quantity.Quantity{skuGen('A'): *quantity.New(7), skuGen('B'): *quantity.New(10)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ch, err := NewCheckout(tt.pricingRules, tt.basketStorage, tt.scanner)
			if err != nil {
				t.Errorf("failed to init checkout: %v", err)
			}

			err = ch.ScanItems()

			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("unexpected error: got %v, want %v", err, tt.expectedErr)
			}

			for sku, quantity := range tt.expectedBasket {

				qty, _ := tt.basketStorage.GetItem(sku)

				if qty != quantity {
					t.Errorf("unexpected quantity for item %s: got %d, want %d", sku, qty, quantity)
				}
			}
		})
	}
}

func Test_checkout_GetTotalPrice(t *testing.T) {

	skuGen := func(r rune) sku.SKU {
		return skuGenerator(t, r)
	}

	type fields struct {
		basket       Basket
		scanner      Scanner
		pricingRules PricingRules
	}
	tests := []struct {
		name   string
		fields fields
		want   currency.Pence
	}{
		{
			name: "gets the price of a single item",
			fields: fields{
				basket: &MockBasketStorage{
					Items: map[sku.SKU]quantity.Quantity{skuGen('A'): *quantity.New(1)},
				},
				scanner: &skuScanner{},
				pricingRules: &pricing.SimplePricing{
					UnitPrices: map[sku.SKU]currency.Pence{skuGen('A'): 10},
				},
			},
			want: 10,
		},
		{
			name: "adds total price for multiple items",
			fields: fields{
				basket: &MockBasketStorage{
					Items: map[sku.SKU]quantity.Quantity{skuGen('A'): *quantity.New(1), skuGen('B'): *quantity.New(2)},
				},
				scanner: &skuScanner{},
				pricingRules: &pricing.SimplePricing{
					UnitPrices: map[sku.SKU]currency.Pence{skuGen('A'): 10, skuGen('B'): 20},
				},
			},
			want: 50,
		},
		{
			name: "calculates special pricing",
			fields: fields{
				basket: &MockBasketStorage{
					Items: map[sku.SKU]quantity.Quantity{skuGen('A'): *quantity.New(3), skuGen('B'): *quantity.New(2)},
				},
				scanner: &skuScanner{},
				pricingRules: &pricing.SpecialPricing{
					Config: map[sku.SKU]pricing.PricingData{
						skuGen('A'): {
							UnitPrice:       10,
							SpecialPrice:    5,
							SpecialQuantity: *quantity.New(3),
						},
					},
				},
			},
			want: 5,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &checkout{
				basket:       tt.fields.basket,
				scanner:      tt.fields.scanner,
				pricingRules: tt.fields.pricingRules,
			}
			if got := c.GetTotalPrice(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("checkout.GetTotalPrice() = %v, want %v", got, tt.want)
			}
		})
	}
}
