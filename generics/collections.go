package generics

// Filter slice filter with predicate
func Filter[T any](slice []T, predicate Predicate[T]) []T {
	var res []T
	for _, el := range slice {
		if predicate(el) {
			res = append(res, el)
		}
	}
	return res
}

// ToMap convert slice to map with provided key function
func ToMap[K comparable, V any](slice []V, keyFunc Function[V, K]) map[K]V {
	m := make(map[K]V, len(slice))
	return Reduce(slice, m, func(value V, val map[K]V) map[K]V {
		val[keyFunc(value)] = value
		return val
	})
}

// Sum return sum of slice numbers
func Sum[T Number](slice []T) T {
	var res T
	for _, t := range slice {
		res += t
	}
	return res
}
