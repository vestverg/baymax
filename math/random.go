package math

import (
	"math/rand"
	"time"

	"github.com/vestverg/baymax/generics"
)

// WeightedRandom  will pick up random value from slice with provided weight
func WeightedRandom[T any, W generics.Number](weights []W, values []T) T {
	if len(weights) == 0 || len(weights) != len(values) {
		panic("")
	}
	rand.Seed(time.Now().UnixNano())
	normalized := normalize(weights)
	target := rand.Float64()
	cumulative := 0.0
	var res T
	for i := 0; i < len(weights); i++ {
		if normalized[i] < 0 {
			panic("negative weight")
		}
		cumulative += normalized[i]
		if target <= cumulative {
			res = values[i]
			break
		}
	}
	return res
}

func normalize[W generics.Number](weights []W) []float64 {
	sum := generics.Sum(weights)
	return generics.Map(weights, func(val W) float64 {
		return float64(val) / float64(sum)
	})
}
