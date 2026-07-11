package h3

// _h3RotatePent60ccw rotates an h3Index 60 degrees counterclockwise about a pentagonal center.
// This function handles pentagon-specific rotation logic, skipping leading 1 digits (k-axis)
// and adjusting for deleted k-axes sequence when necessary.
// Ported from H3 C: h3Index.c::_h3RotatePent60ccw.
func _h3RotatePent60ccw(h h3Index) h3Index {
	// rotate in place; skips any leading 1 digits (k-axis)
	foundFirstNonZeroDigit := false
	res := getResolution(h)

	for r := int32(1); r <= res; r++ {
		// rotate this digit
		oldDigit := h3GetIndexDigit(h, r)
		h = h3SetIndexDigit(h, r, int32(_rotate60ccw(direction(oldDigit))))

		// look for the first non-zero digit so we can adjust for deleted k-axes sequence if necessary
		if !foundFirstNonZeroDigit && h3GetIndexDigit(h, r) != 0 {
			foundFirstNonZeroDigit = true

			// adjust for deleted k-axes sequence
			if direction(_h3LeadingNonZeroDigit(h)) == kAxesDigit {
				h = _h3Rotate60ccw(h)
			}
		}
	}
	return h
}
