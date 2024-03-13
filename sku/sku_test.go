package sku

import (
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	type args struct {
		value rune
	}
	tests := []struct {
		name    string
		args    args
		want    SKU
		wantErr error
	}{
		{
			name:    "should return error when given empty sku",
			wantErr: ErrNoSpecialCharacters,
		},
		{
			name:    "should return error when given symbol",
			args:    args{value: '$'},
			wantErr: ErrNoSpecialCharacters,
		},
		{
			name:    "should return error when given number",
			args:    args{value: '7'},
			wantErr: ErrNoSpecialCharacters,
		},
		{
			name:    "converts given sku to uppercase",
			args:    args{value: 'd'},
			want:    SKU{value: 'D'},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.value)
			if err != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}
