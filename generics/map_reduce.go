package generics

// Map perform mapFunc on each element of slice and return new slice
func Map[T any, R any](slice []T, mapFunc Function[T, R]) []R {
	ts := make([]R, len(slice))
	for i, t := range slice {
		ts[i] = mapFunc(t)
	}
	return ts
}

// Reduce  perform reduce function on slice
func Reduce[T any, R any](slice []T, identity R, reduceFunc ReduceFunction[T, R]) R {
	res := identity
	for _, v := range slice {
		res = reduceFunc(v, res)
	}
	return res
}
