package c2go

// setH3Index initializes an H3 index (port of C setH3Index).
// Mirrors behavior: start from H3_INIT (all digits 7), set cell mode, resolution,
// base cell, and initialize digits 1..res to initDigit.
func setH3Index(hp *H3Index, res int, baseCell int, initDigit int) {
	var h uint64 = H3_INIT
	// Set mode to cell mode
	h &^= H3_MODE_MASK
	h |= uint64(H3_CELL_MODE) << H3_MODE_OFFSET
	// Set resolution
	h &^= H3_RES_MASK
	h |= (uint64(res) & 15) << H3_RES_OFFSET
	// Set base cell
	h &^= H3_BC_MASK
	h |= (uint64(baseCell) & 127) << H3_BC_OFFSET
	// Initialize digits 1..res
	out := H3Index(h)
	for r := 1; r <= res; r++ {
		out = setIndexDigit(out, r, initDigit)
	}
	*hp = out
}
