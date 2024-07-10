package math

import (
	"errors"
	"math/rand"
	"time"

	"golang.org/x/exp/constraints"

	"github.com/vestverg/baymax/generics"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type Number interface {
	constraints.Float | constraints.Integer
}

// WeightedRandom  will pick up random value from slice with provided weight
func WeightedRandom[T any, W Number](weights []W, values []T) (T, error) {
	var res T
	if len(weights) == 0 || len(weights) != len(values) {
		return res, errors.New("weights and values must have equal length and cannot be empty")
	}

	normalized, err := normalize(weights)
	if err != nil {
		return res, err
	}
	target := rand.Float64()
	cumulative := 0.0

	for i, norm := range normalized {
		if norm < 0 {
			return res, errors.New("negative weight")
		}
		cumulative += norm
		if target <= cumulative {
			res = values[i]
			break
		}
	}
	return res, nil
}

func normalize[W Number](weights []W) ([]float64, error) {
	sum := generics.Sum(weights)
	if sum == 0 {
		return nil, errors.New("sum of weights cannot be zero")
	}

	return generics.Map(weights, func(val W) float64 {
		return float64(val) / float64(sum)
	}), nil
}
