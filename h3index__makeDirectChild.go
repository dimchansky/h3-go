package h3

// makeDirectChild returns the immediate child index for the given cell and child digit (0-6).
// Port of C makeDirectChild (bit operations only).
// Ported from H3 C: h3Index.c::makeDirectChild.
func makeDirectChild(h h3Index, cellNumber int32) h3Index {
	childRes := getResolution(h) + 1
	// Set resolution to childRes
	h = setResolution(h, childRes)
	// Set the digit at childRes
	h = h3SetIndexDigit(h, childRes, cellNumber)
	return h
}
