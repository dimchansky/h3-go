package c2go

// setH3Index initializes an H3 index (port of C setH3Index).
// Mirrors behavior: start from H3_INIT (all digits 7), set cell mode, resolution,
// base cell, and initialize digits 1..res to initDigit.
func setH3Index(hp *H3Index, res int, baseCell int, initDigit int) {
    const (
        h3Init     = uint64(35184372088831) // H3_INIT = 2^45 - 1
        modeOffset = 59
        modeMask   = uint64(15) << modeOffset
        resOffset  = 52
        resMask    = uint64(15) << resOffset
        bcOffset   = 45
        bcMask     = uint64(127) << bcOffset
    )

	var h uint64 = h3Init
	// Set mode to cell mode
	h &^= modeMask
    h |= uint64(H3_CELL_MODE) << modeOffset
	// Set resolution
	h &^= resMask
	h |= (uint64(res) & 15) << resOffset
	// Set base cell
	h &^= bcMask
	h |= (uint64(baseCell) & 127) << bcOffset
	// Initialize digits 1..res
	out := H3Index(h)
	for r := 1; r <= res; r++ {
		out = setIndexDigit(out, r, initDigit)
	}
	*hp = out
}
