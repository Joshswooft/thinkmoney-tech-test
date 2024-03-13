package checkout

import (
	"testing"

	"github.com/Joshswooft/thinkmoney-test/pricing"
)

func TestNewCheckout(t *testing.T) {

	pricingRules := &pricing.SpecialPricing{}

	type args struct {
		pricingRules PricingRules
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
			name: "happy path",
			args: args{pricingRules: pricingRules},
			err:  nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewCheckout(tt.args.pricingRules)
			if err != tt.err {
				t.Errorf("NewCheckout() error = %v, wantErr %v", err, tt.err)
				return
			}
		})
	}
}
