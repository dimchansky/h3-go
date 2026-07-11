package h3

// setH3Index initializes an H3 index (port of C setH3Index).
// Mirrors behavior: start from h3Init (all digits 7), set cell mode, resolution,
// base cell, and initialize digits 1..res to initDigit.
// Ported from H3 C: h3Index.c::setH3Index.
func setH3Index(hp *h3Index, res int32, baseCell int32, initDigit int32) {
	h := h3Index(h3Init)
	// Set mode to cell mode
	h = setMode(h, h3CellMode)
	// Set resolution
	h = setResolution(h, res)
	// Set base cell
	h = setBaseCell(h, baseCell)
	// Initialize digits 1..res
	for r := int32(1); r <= res; r++ {
		h = setIndexDigit(h, r, initDigit)
	}
	*hp = h
}
