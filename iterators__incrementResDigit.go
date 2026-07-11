package h3

// incrementResDigit increments the digit (0-7) at location `res` in an H3 index.
// This function is used internally by the H3 iterator system to advance through
// child cells by incrementing the appropriate resolution digit.
// Ported from H3 C: iterators.c::_incrementResDigit.
func incrementResDigit(h *h3Index, res int32) {
	val := h3Index(1)
	val <<= h3PerDigitOffset * (maxH3Res - res)
	*h += val
}
