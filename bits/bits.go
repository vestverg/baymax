package bits

import (
	"golang.org/x/exp/constraints"
)

func SetBit[T constraints.Integer](val T, idx int) T {
	if idx < 0 || idx > 63 {
		panic("index out of range")
	}
	return val | 1<<idx
}

func IsSet[T constraints.Integer](val T, idx int) bool {
	if idx < 0 || idx > 63 {
		panic("index out of range")
	}

	return (val>>idx)%2 != 0
}

func NumberOfTrailingZeros64[T constraints.Integer](val T) int {
	x := int32(val)
	if x == 0 {
		return 32 + NumberOfTrailingZeros32(int32(val>>32))
	}
	return NumberOfTrailingZeros32(x)
}

// NumberOfTrailingZeros32 finds the number of trailing zeros in a 32-bit value.
func NumberOfTrailingZeros32[T ~int32 | ~uint32](val T) int {
	i := computeI(val)
	if i <= 0 {
		return int(i & 32)
	}
	n := computeN(i)
	return n
}

// computeI computes the bitwise operation needed for the calculation.
func computeI[T ~int32 | ~uint32](val T) (i T) {
	i = ^val & (val - 1)
	return i
}

// computeN further computes the number of trailing zeros.
func computeN[T ~int32 | ~uint32](i T) (n int) {
	n = 1
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
	n += int(i >> 1)
	return n
}
