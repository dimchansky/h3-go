package h3

// _h3LeadingNonZeroDigit returns the highest resolution non-zero digit.
// Ported from H3 C: h3Index.c::_h3LeadingNonZeroDigit.
func _h3LeadingNonZeroDigit(h h3Index) int32 {
	res := getResolution(h)
	for r := int32(1); r <= res; r++ {
		d := getIndexDigit(h, r)
		if d != 0 {
			return d
		}
	}
	// All zeros => centerDigit (0)
	return 0
}
