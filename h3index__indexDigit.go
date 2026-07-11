package h3

// h3GetIndexDigit returns the direction digit at res (port of H3_GET_INDEX_DIGIT).
// Ported from H3 C: h3Index.h::H3_GET_INDEX_DIGIT.
func h3GetIndexDigit(h h3Index, res int32) int32 {
	shift := (maxH3Res - res) * h3PerDigitOffset
	return int32((uint64(h) >> shift) & h3DigitMask)
}

// h3SetIndexDigit sets the direction digit at res (port of H3_SET_INDEX_DIGIT).
// Ported from H3 C: h3Index.h::H3_SET_INDEX_DIGIT.
func h3SetIndexDigit(h h3Index, res int32, digit int32) h3Index {
	shift := (maxH3Res - res) * h3PerDigitOffset
	mask := h3DigitMask << shift
	x := uint64(h)
	x &^= mask
	x |= (uint64(digit) & h3DigitMask) << shift
	return h3Index(x)
}
