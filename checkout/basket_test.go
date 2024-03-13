package checkout

import (
	"reflect"
	"testing"

	"github.com/Joshswooft/thinkmoney-test/quantity"
	"github.com/Joshswooft/thinkmoney-test/sku"

	"errors"
)

func Test_basket_AddItem(t *testing.T) {
	type fields struct {
		items map[itemID]quantity.Quantity
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
		name      string
		fields    fields
		args      args
		wantItems map[itemID]quantity.Quantity
	}{
		{
			name:      "adds an item to a nil basket",
			fields:    fields{items: nil},
			args:      args{sku: skuA, quantity: *quantity.New(0)},
			wantItems: map[itemID]quantity.Quantity{skuA: *quantity.New(0)},
		},
		{
			name:      "adds an item to an empty basket",
			fields:    fields{items: make(map[itemID]quantity.Quantity)},
			args:      args{sku: skuA, quantity: *quantity.New(0)},
			wantItems: map[itemID]quantity.Quantity{skuA: *quantity.New(0)},
		},
		{
			name:      "adds an item to an existing basket",
			fields:    fields{items: map[itemID]quantity.Quantity{skuB: *quantity.New(1)}},
			args:      args{sku: skuA, quantity: *quantity.New(2)},
			wantItems: map[itemID]quantity.Quantity{skuB: *quantity.New(1), skuA: *quantity.New(2)},
		},
		{
			name:      "updates an existing item in the basket with the given quantity",
			fields:    fields{items: map[itemID]quantity.Quantity{skuB: *quantity.New(1)}},
			args:      args{sku: skuB, quantity: *quantity.New(2)},
			wantItems: map[itemID]quantity.Quantity{skuB: *quantity.New(2)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &basket{
				items: tt.fields.items,
			}
			b.AddItem(tt.args.sku, tt.args.quantity)
		})
	}
}

func Test_basket_GetItem(t *testing.T) {
	type fields struct {
		items map[itemID]quantity.Quantity
	}
	type args struct {
		sku sku.SKU
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
		name    string
		fields  fields
		args    args
		wantQty quantity.Quantity
		wantErr error
	}{
		{
			name:    "returns a not found error when getting from a nil basket",
			fields:  fields{items: nil},
			args:    args{sku: skuA},
			wantQty: *quantity.New(0),
			wantErr: ErrItemNotFound,
		},
		{
			name:    "returns a not found error when getting from an empty basket",
			fields:  fields{items: make(map[itemID]quantity.Quantity)},
			args:    args{sku: skuA},
			wantQty: *quantity.New(0),
			wantErr: ErrItemNotFound,
		},
		{
			name:    "returns a not found error when item doesnt exist in the basket",
			fields:  fields{items: map[itemID]quantity.Quantity{skuB: *quantity.New(42)}},
			args:    args{sku: skuA},
			wantQty: *quantity.New(0),
			wantErr: ErrItemNotFound,
		},
		{
			name:    "returns the found items quantity from the basket",
			fields:  fields{items: map[itemID]quantity.Quantity{skuB: *quantity.New(42)}},
			args:    args{sku: skuB},
			wantQty: *quantity.New(42),
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &basket{
				items: tt.fields.items,
			}
			gotQty, err := b.GetItem(tt.args.sku)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("errors dont match = %v, want %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(gotQty, tt.wantQty) {
				t.Errorf("basket.GetItem() = %v, want %v", gotQty, tt.wantQty)
			}
		})
	}
}
