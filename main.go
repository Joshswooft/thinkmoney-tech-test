package main

import (
	"fmt"
	"log"

	"strings"

	"github.com/Joshswooft/thinkmoney-test/checkout"
	"github.com/Joshswooft/thinkmoney-test/pricing"
	"github.com/Joshswooft/thinkmoney-test/quantity"
	"github.com/Joshswooft/thinkmoney-test/sku"
)

func main() {
	input := "a69B$42*0(Cdb"
	fmt.Println("input: ", input)
	scanner, err := checkout.NewSkuScanner(strings.NewReader(input))

	if err != nil {
		log.Fatal(err)
	}

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
	ch, err := checkout.NewCheckout(&pricingRules, basket, scanner)

	if err != nil {
		log.Fatal(err)
	}

	if err := ch.ScanItems(); err != nil {
		log.Fatal(err)
	}

	total := ch.GetTotalPrice()

	fmt.Printf("total is: %d pence \n", total)

}
