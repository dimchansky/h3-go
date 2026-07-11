package h3

// getMode returns the H3 mode (port of H3_GET_MODE)
// Ported from H3 C: h3Index.h::H3_GET_MODE.
func getMode(h h3Index) int32 {
	return int32((uint64(h) & h3ModeMask) >> h3ModeOffset)
}

// setMode sets the H3 mode (port of H3_SET_MODE)
// Ported from H3 C: h3Index.h::H3_SET_MODE.
func setMode(h h3Index, v int32) h3Index {
	x := uint64(h)
	x &^= h3ModeMask
	x |= (uint64(v) & 15) << h3ModeOffset
	return h3Index(x)
}

// getHighBit returns the highest bit (port of H3_GET_HIGH_BIT)
// Ported from H3 C: h3Index.h::H3_GET_HIGH_BIT.
//
//nolint:unused // ported from H3 C for parity completeness; exercised by cgo && c2go parity tests
func getHighBit(h h3Index) int32 {
	if (uint64(h)&h3HighBitMask)>>h3MaxOffset != 0 {
		return 1
	}
	return 0
}

// setHighBit sets the highest bit (port of H3_SET_HIGH_BIT)
// Ported from H3 C: h3Index.h::H3_SET_HIGH_BIT.
func setHighBit(h h3Index, v int32) h3Index {
	x := uint64(h)
	x &^= h3HighBitMask
	x |= (uint64(v) & 1) << h3MaxOffset
	return h3Index(x)
}
