package pricing

import (
	"reflect"
	"testing"
)

func TestSpecialPricing_GetPrice(t *testing.T) {

	skuA := 'A'
	skuB := 'B'

	type fields struct {
		pricing map[sku]PricingData
	}
	type args struct {
		sku      sku
		quantity quantity
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   pence
	}{
		{
			name:   "should return 0 pence when no pricing rules exist",
			fields: fields{pricing: nil},
			args:   args{sku: skuA, quantity: 5},
			want:   0,
		},
		{
			name: "returns 0 pence for a sku that doesnt exist in the pricing rules",
			fields: fields{
				pricing: map[sku]PricingData{skuA: {UnitPrice: 50}},
			},
			args: args{
				sku:      skuB,
				quantity: 1,
			},
			want: 0,
		},
		{
			name: "returns the price of the product - pineapple = 50p",
			fields: fields{
				pricing: map[sku]PricingData{skuA: {UnitPrice: 50}},
			},
			args: args{sku: skuA, quantity: 1},
			want: 50,
		},
		{
			name: "buying 3 pineapples for 130p special offer met",
			fields: fields{
				pricing: map[sku]PricingData{skuA: {UnitPrice: 50, SpecialPrice: 130, SpecialQuantity: 3}},
			},
			args: args{sku: skuA, quantity: 3},
			want: 130,
		},
		{
			name: "buying more pineapples than the special quantity still applies the special price + the unit price",
			fields: fields{
				pricing: map[sku]PricingData{skuA: {UnitPrice: 50, SpecialPrice: 130, SpecialQuantity: 3}},
			},
			args: args{sku: skuA, quantity: 4},
			want: 130 + 50,
		},
		{
			name: "special price for this day replaces the unit price",
			fields: fields{
				pricing: map[sku]PricingData{skuA: {UnitPrice: 50, SpecialPrice: 35, SpecialQuantity: 0}},
			},
			args: args{sku: skuA, quantity: 2},
			want: 35 + 35,
		},
		{
			name: "Free item",
			fields: fields{
				pricing: map[sku]PricingData{skuA: {UnitPrice: 50, SpecialQuantity: 1}},
			},
			args: args{sku: skuA, quantity: 2},
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
