package pricing

import (
	"reflect"
	"testing"

	"github.com/Joshswooft/thinkmoney-test/currency"
	"github.com/Joshswooft/thinkmoney-test/quantity"
	"github.com/Joshswooft/thinkmoney-test/sku"
)

func TestSpecialPricing_GetPrice(t *testing.T) {
	type fields struct {
		pricing map[sku.SKU]PricingData
	}
	type args struct {
		sku      sku.SKU
		quantity quantity.Quantity
	}

	skuA, err := sku.New('A')

	if err != nil {
		t.Fatalf("creating sku A failed: %v", err)
	}

	skuB, err := sku.New('B')

	if err != nil {
		t.Fatalf("creating sku B failed: %v", err)
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   currency.Pence
	}{
		{
			name:   "should return 0 pence when no pricing rules exist",
			fields: fields{pricing: nil},
			args:   args{sku: skuA, quantity: *quantity.New(5)},
			want:   currency.Pence(0),
		},
		{
			name: "returns 0 pence for a sku that doesnt exist in the pricing rules",
			fields: fields{
				pricing: map[sku.SKU]PricingData{skuA: {UnitPrice: 50}},
			},
			args: args{
				sku:      skuB,
				quantity: *quantity.New(1),
			},
			want: currency.Pence(0),
		},
		{
			name: "returns the price of the product - pineapple = 50p",
			fields: fields{
				pricing: map[sku.SKU]PricingData{skuA: {UnitPrice: 50}},
			},
			args: args{sku: skuA, quantity: *quantity.New(1)},
			want: currency.Pence(50),
		},
		{
			name: "buying 3 pineapples for 130p special offer met",
			fields: fields{
				pricing: map[sku.SKU]PricingData{skuA: {UnitPrice: 50, SpecialPrice: 130, SpecialQuantity: *quantity.New(3)}},
			},
			args: args{sku: skuA, quantity: *quantity.New(3)},
			want: currency.Pence(130),
		},
		{
			name: "buying more pineapples than the special quantity still applies the special price + the unit price",
			fields: fields{
				pricing: map[sku.SKU]PricingData{skuA: {UnitPrice: 50, SpecialPrice: 130, SpecialQuantity: *quantity.New(3)}},
			},
			args: args{sku: skuA, quantity: *quantity.New(4)},
			want: 130 + 50,
		},
		{
			name: "special price for this day replaces the unit price",
			fields: fields{
				pricing: map[sku.SKU]PricingData{skuA: {UnitPrice: 50, SpecialPrice: 35, SpecialQuantity: *quantity.New(0)}},
			},
			args: args{sku: skuA, quantity: *quantity.New(2)},
			want: 35 + 35,
		},
		{
			name: "Free item",
			fields: fields{
				pricing: map[sku.SKU]PricingData{skuA: {UnitPrice: 50, SpecialQuantity: *quantity.New(1)}},
			},
			args: args{sku: skuA, quantity: *quantity.New(2)},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &SpecialPricing{
				Config: tt.fields.pricing,
			}
			if got := p.GetPrice(tt.args.sku, tt.args.quantity); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SpecialPricing.GetPrice() = %v, want %v", got, tt.want)
			}
		})
	}
}
