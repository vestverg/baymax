package bits

import (
	"fmt"
	"strconv"
	"testing"
)

func TestNumberOfTrailingZeros64(t *testing.T) {

	type tc struct {
		val      int
		expected int
	}
	tests := []tc{
		{1 << 0, 0},
		{1 << 1, 1},
		{1 << 2, 2},
		{1 << 3, 3},
		{1 << 31, 31},
		{1 << 32, 32},
		{1 << 33, 33},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s", strconv.FormatInt(int64(tt.val), 2)), func(t *testing.T) {
			if got := NumberOfTrailingZeros64(tt.val); got != tt.expected {
				t.Errorf("NumberOfTrailingZeros64() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestSetBit(t *testing.T) {

	type tc struct {
		idx   int
		isSet bool
	}
	tests := []tc{
		{0, true},
		{1, true},
		{2, true},
		{3, true},
		{31, true},
		{32, true},
		{33, true},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%d", tt.idx), func(t *testing.T) {
			if got := IsSet(SetBit(1, tt.idx), tt.idx); got != tt.isSet {
				t.Errorf("TestSetBit() = %v, want %v", got, tt.isSet)
			}
		})
	}
}
func TestIsSet(t *testing.T) {

	type tc struct {
		val      int
		idx      int
		expected bool
	}
	tests := []tc{
		{2 << 1, 2, true},
		{2 << 2, 3, true},
		{2 << 3, 4, true},
		{2 << 32, 33, true},
		{2 << 33, 34, true},
		{2 << 33, 1, false},
		{2 << 33, 32, false},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			if got := IsSet(tt.val, tt.idx); got != tt.expected {
				t.Errorf("TestIsSet() = %v, want %v", got, tt.expected)
			}
		})
	}
}
