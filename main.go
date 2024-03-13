package main

import (
	"fmt"
	"log"

	"github.com/Joshswooft/thinkmoney-test/checkout"
	"github.com/Joshswooft/thinkmoney-test/pricing"
	"github.com/Joshswooft/thinkmoney-test/quantity"
	"github.com/Joshswooft/thinkmoney-test/sku"
)

func main() {
	input := "a69B$42*0(Cdb"

	skuA, _ := sku.New('A')
	skuB, _ := sku.New('B')
	skuC, _ := sku.New('C')

	pricingRules := pricing.SpecialPricing{
		Config: map[sku.SKU]pricing.PricingData{
			skuA: {UnitPrice: 10},
			skuB: {UnitPrice: 20, SpecialPrice: 10, SpecialQuantity: *quantity.New(2)},
			skuC: {UnitPrice: 50, SpecialPrice: 30, SpecialQuantity: *quantity.New(5)},
		},
	}

	basket := checkout.NewBasket()

	ch, err := checkout.NewCheckout(&pricingRules, basket)

	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < len(input); i++ {
		r := rune(input[i])
		skuInstance, err := sku.New(r)

		if err != nil {
			log.Println(err)
		}

		ch.Scan(skuInstance, *quantity.New(1))
	}

	total := ch.GetTotalPrice()

	fmt.Printf("total: %d pence", total)
}
