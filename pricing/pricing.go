package pricing

type sku = rune
type pence = int
type quantity = uint

// exposing pricing config for convenience
type SpecialPricing struct {
	Config map[sku]interface{}
}

func (p *SpecialPricing) GetPrice(sku sku, quantity quantity) pence {

	if p == nil || p.Config == nil {
		return 0
	}

	return 0

}
