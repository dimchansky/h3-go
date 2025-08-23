package c2go

// _h3Rotate60cw rotates an H3Index 60 degrees clockwise.
// Rotates each digit at every resolution using the _rotate60cw transform.
// Ported from H3 C: h3Index.c::_h3Rotate60cw
func _h3Rotate60cw(h H3Index) H3Index {
	res := getResolution(h)
	for r := 1; r <= res; r++ {
		oldDigit := getIndexDigit(h, r)
		h = setIndexDigit(h, r, int(_rotate60cw(Direction(oldDigit))))
	}
	return h
}
