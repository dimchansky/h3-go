package h3

// _h3Rotate60ccw rotates an H3Index 60 degrees counterclockwise.
// Rotates each digit at every resolution using the _rotate60ccw transform.
// Ported from H3 C: h3Index.c::_h3Rotate60ccw
func _h3Rotate60ccw(h H3Index) H3Index {
	res := getResolution(h)
	for r := int32(1); r <= res; r++ {
		oldDigit := getIndexDigit(h, r)
		h = setIndexDigit(h, r, int32(_rotate60ccw(Direction(oldDigit))))
	}
	return h
}
