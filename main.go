package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/Joshswooft/thinkmoney-test/checkout"
	"github.com/Joshswooft/thinkmoney-test/pricing"
	"github.com/Joshswooft/thinkmoney-test/quantity"
	"github.com/Joshswooft/thinkmoney-test/sku"
)

func main() {
	input := "a69B$42*0(Cdb"

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

	ch, err := checkout.NewCheckout(&pricingRules, basket)

	if err != nil {
		log.Fatal(err)
	}

	buf := make([]byte, 64)

	for {
		bytesRead, err := scanner.Read(buf)

		if err == io.EOF {
			break
		}

		if errors.Is(err, sku.ErrNoSpecialCharacters) {
			continue
		}

		for i := 0; i < bytesRead; i++ {
			r := rune(input[i])
			skuInstance, err := sku.New(r)

			if err != nil {
				log.Println(err)
			}

			ch.Scan(skuInstance, *quantity.New(1))
		}

	}

	total := ch.GetTotalPrice()

	fmt.Printf("total: %d pence", total)
}
