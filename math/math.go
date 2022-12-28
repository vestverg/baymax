package math

import "github.com/vestverg/baymax/generics"

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
