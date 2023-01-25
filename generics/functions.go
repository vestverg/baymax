package generics

// Function mapping function T type -> R type
type Function[T any, R any] func(val T) R

// Predicate check condition T -> bool
type Predicate[T any] func(val T) bool

// ReduceFunction (T type,R type ) -> R type reduce two values in one
type ReduceFunction[T any, R any] func(v1 T, v2 R) R

// Consumer consume value
type Consumer[T any] func(val T)

// Comparator compare function signature
type Comparator[T any] func(val1 T, val2 T) int
