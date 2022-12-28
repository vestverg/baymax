package generics

import "golang.org/x/exp/constraints"

// Number ...
type Number interface {
	constraints.Float | constraints.Integer
}
