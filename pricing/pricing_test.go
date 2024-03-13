package pricing

import (
	"reflect"
	"testing"
)

func TestSpecialPricing_GetPrice(t *testing.T) {

	skuA := 'A'

	type fields struct {
		pricing map[sku]interface{}
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
