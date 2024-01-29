package main

import (
	"log"
	"math"

	"github.com/pkg/errors"
)

const (
	BasicWeight      = 1
	BasicPrice       = 18.0
	InsurePercentage = 0.01
)

func calc(weight float64) (int, error) {
	w := int(math.Round(weight))

	if w <= 1 {
		log.Println(BasicPrice)
		return BasicPrice, nil
	}
	if w > 100 {
		return 0, errors.Errorf("too large weight: %d\n", w)
	}

	sum := 0.0
	for i := 1; i <= w; i++ {
		additionalWeight := w - 1
		shippingPrice := BasicPrice + (5 * additionalWeight)
		insurePrice := (float64(shippingPrice) + sum) * InsurePercentage
		sum = float64(shippingPrice) + +insurePrice
	}

	res := int(math.Round(sum))

	return res, nil
}
