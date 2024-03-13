package checkout

import (
	"testing"

	"github.com/Joshswooft/thinkmoney-test/currency"
	"github.com/Joshswooft/thinkmoney-test/pricing"
	"github.com/Joshswooft/thinkmoney-test/quantity"
	"github.com/Joshswooft/thinkmoney-test/sku"
)

func TestNewCheckout(t *testing.T) {

	pricingRules := &pricing.SpecialPricing{}
	basket := NewBasket()

	type args struct {
		pricingRules PricingRules
		basket       Basket
	}
	tests := []struct {
		name string
		args args
		err  error
	}{
		{
			name: "returns error when given no pricing rules",
			err:  errNoPricingRulesProvided,
		},
		{
			name: "returns error when given no basket",
			args: args{pricingRules: pricingRules},
			err:  errNoBasketProvided,
		},
		{
			name: "happy path",
			args: args{pricingRules: pricingRules, basket: basket},
			err:  nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewCheckout(tt.args.pricingRules, tt.args.basket)
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
			name: "scanning in a product with quantity of 0 does nothing",
			fields: fields{
				basket:       NewBasket(),
				pricingRules: &pricing.SimplePricing{UnitPrices: map[sku.SKU]currency.Pence{skuA: 10}},
			},
			args:    args{sku: skuA, quantity: *quantity.New(0)},
			wantErr: false,
			qty:     *quantity.New(0),
		},
		{
			name: "scans in a new item with quantity = 1",
			fields: fields{
				basket:       NewBasket(),
				pricingRules: &pricing.SimplePricing{UnitPrices: map[sku.SKU]currency.Pence{skuA: 10}},
			},
			args:    args{sku: skuA, quantity: *quantity.New(1)},
			wantErr: false,
			qty:     *quantity.New(1),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &checkout{
				basket:       tt.fields.basket,
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
