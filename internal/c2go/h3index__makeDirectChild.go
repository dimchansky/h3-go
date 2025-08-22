package c2go

// makeDirectChild returns the immediate child index for the given cell and child digit (0-6).
// Port of C makeDirectChild (bit operations only).
func makeDirectChild(h H3Index, cellNumber int) H3Index {
	childRes := getResolution(h) + 1
	// Set resolution to childRes
	const resOffset = 52
	const resMask = uint64(15) << resOffset
	x := uint64(h)
	x &^= resMask
	x |= (uint64(childRes) & 15) << resOffset
	out := H3Index(x)
	// Set the digit at childRes
	out = setIndexDigit(out, childRes, cellNumber)
	return out
}
