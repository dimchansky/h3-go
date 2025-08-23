package c2go

// _h3RotatePent60cw rotates an H3Index 60 degrees clockwise about a pentagonal center.
// This function handles pentagon-specific rotation logic, skipping leading 1 digits (k-axis)
// and adjusting for deleted k-axes sequence when necessary.
// Ported from H3 C: h3Index.c::_h3RotatePent60cw
func _h3RotatePent60cw(h H3Index) H3Index {
	// rotate in place; skips any leading 1 digits (k-axis)
	foundFirstNonZeroDigit := false
	res := getResolution(h)

	for r := 1; r <= res; r++ {
		// rotate this digit
		oldDigit := getIndexDigit(h, r)
		h = setIndexDigit(h, r, int(_rotate60cw(Direction(oldDigit))))

		// look for the first non-zero digit so we can adjust for deleted k-axes sequence if necessary
		if !foundFirstNonZeroDigit && getIndexDigit(h, r) != 0 {
			foundFirstNonZeroDigit = true

			// adjust for deleted k-axes sequence
			if Direction(_h3LeadingNonZeroDigit(h)) == K_AXES_DIGIT {
				h = _h3Rotate60cw(h)
			}
		}
	}
	return h
}
