package bits

import (
	"golang.org/x/exp/constraints"
)

func SetBit[T constraints.Integer](val T, idx int) T {
	return val | 1<<idx
}

func IsSet[T constraints.Integer](val T, idx int) bool {
	return (val>>idx)%2 != 0
}

func NumberOfTrailingZeros64[T constraints.Integer](val T) int {
	x := int32(val)
	if x == 0 {
		return 32 + NumberOfTrailingZeros32(int32(val>>32))
	}
	return NumberOfTrailingZeros32(x)
}

func NumberOfTrailingZeros32[T ~int32 | ~uint32](val T) int {
	i := ^val & (val - 1)
	if i <= 0 {
		return int(i & 32)
	}
	n := 1
	if i > 1<<16 {
		n += 16
		i >>= 16
	}
	if i > 1<<8 {
		n += 8
		i >>= 8
	}
	if i > 1<<4 {
		n += 4
		i >>= 4
	}
	if i > 1<<2 {
		n += 2
		i >>= 2
	}
	return n + int(i>>1)
}
