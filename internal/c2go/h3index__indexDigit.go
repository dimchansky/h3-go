package c2go

// getIndexDigit returns the direction digit at res (port of H3_GET_INDEX_DIGIT).
// Ported from H3 C: h3Index.h::H3_GET_INDEX_DIGIT
func getIndexDigit(h H3Index, res int32) int32 {
	shift := (MAX_H3_RES - res) * H3_PER_DIGIT_OFFSET
	return int32((uint64(h) >> shift) & H3_DIGIT_MASK)
}

// setIndexDigit sets the direction digit at res (port of H3_SET_INDEX_DIGIT).
// Ported from H3 C: h3Index.h::H3_SET_INDEX_DIGIT
func setIndexDigit(h H3Index, res int32, digit int32) H3Index {
	shift := (MAX_H3_RES - res) * H3_PER_DIGIT_OFFSET
	mask := H3_DIGIT_MASK << shift
	x := uint64(h)
	x &^= mask
	x |= (uint64(digit) & H3_DIGIT_MASK) << shift
	return H3Index(x)
}
