package quantity

import (
	"testing"
)

func TestNew(t *testing.T) {
	type args struct {
		qty int
	}
	tests := []struct {
		name      string
		args      args
		wantValue int
	}{
		{
			name:      "a negative quantity gets turned into 0",
			args:      args{qty: -10},
			wantValue: 0,
		},
		{
			name:      "converts an int into a valid quantity",
			args:      args{qty: 99},
			wantValue: 99,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.qty); got.value != tt.wantValue {
				t.Errorf("unexpected quantity value, got: %s, want: %d", got, tt.wantValue)
			}
		})
	}
}

func TestQuantity_Add(t *testing.T) {
	type args struct {
		amount int
	}
	tests := []struct {
		name string
		qty  Quantity
		args args
		want int
	}{
		{
			name: "updates the internal value of the quantity by the given amount",
			qty:  *New(5),
			args: args{amount: 1},
			want: 6,
		},
		{
			name: "can add negative amount",
			qty:  *New(9),
			args: args{amount: -2},
			want: 7,
		},
		{
			name: "adding negative amount cant take quantity below 0",
			qty:  *New(1),
			args: args{amount: -50},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.qty.Add(tt.args.amount); got != tt.want {
				t.Errorf("Quantity.Add() = %v, want %v", got, tt.want)
			}

			if tt.qty.value != tt.want {
				t.Errorf("internal quantity value not correct, got: %d, want: %d", tt.qty.value, tt.want)
			}
		})
	}
}
