package math

import (
	"golang.org/x/exp/constraints"

	"github.com/vestverg/baymax/generics"
)

// Max ...
func Max[T generics.Number](first T, rest ...T) T {
	var ans = first
	for _, value := range rest {
		if value > ans {
			ans = value
		}
	}
	return ans
}

// Min ...
func Min[T generics.Number](first T, rest ...T) T {
	var ans = first
	for _, value := range rest {
		if value < ans {
			ans = value
		}
	}
	return ans
}

func Compare[T constraints.Ordered](val1, val2 T) int {
	if val1 < val2 {
		return -1
	} else if val1 == val2 {
		return 0
	}
	return 1
}
