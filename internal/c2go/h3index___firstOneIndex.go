package c2go

import "math/bits"

// _firstOneIndex returns the index of the highest set bit in an H3Index.
// This is equivalent to the C implementation using __builtin_clzll or portable fallback.
// Ported from H3 C: h3Index.c::_firstOneIndex
func _firstOneIndex(h H3Index) int {
	if h == 0 {
		return -1 // Handle edge case where no bits are set
	}
	return 63 - bits.LeadingZeros64(uint64(h))
}
