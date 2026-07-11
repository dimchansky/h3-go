package h3

// _h3Rotate60cw rotates an h3Index 60 degrees clockwise.
// Rotates each digit at every resolution using the _rotate60cw transform.
// Ported from H3 C: h3Index.c::_h3Rotate60cw.
func _h3Rotate60cw(h h3Index) h3Index {
	res := getResolution(h)
	for r := int32(1); r <= res; r++ {
		oldDigit := h3GetIndexDigit(h, r)
		h = h3SetIndexDigit(h, r, int32(_rotate60cw(direction(oldDigit))))
	}
	return h
}
