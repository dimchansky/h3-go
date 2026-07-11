package h3

// getResolution returns the H3 resolution (port of H3_GET_RESOLUTION)
// Ported from H3 C: h3Index.c::getResolution.
func getResolution(h h3Index) int32 {
	return int32((uint64(h) & h3ResMask) >> h3ResOffset)
}

// getBaseCell returns the base cell (port of H3_GET_BASE_CELL)
// Ported from H3 C: h3Index.c::H3_GET_BASE_CELL.
func getBaseCell(h h3Index) int32 {
	return int32((uint64(h) & h3BcMask) >> h3BcOffset)
}

// setResolution sets the H3 resolution (port of H3_SET_RESOLUTION)
// Ported from H3 C: h3Index.h::H3_SET_RESOLUTION.
func setResolution(h h3Index, res int32) h3Index {
	x := uint64(h)
	x &^= h3ResMask
	x |= (uint64(res) & 15) << h3ResOffset
	return h3Index(x)
}

// setBaseCell sets the base cell number (port of H3_SET_BASE_CELL)
// Ported from H3 C: h3Index.h::H3_SET_BASE_CELL.
func setBaseCell(h h3Index, baseCell int32) h3Index {
	x := uint64(h)
	x &^= h3BcMask
	x |= (uint64(baseCell) & 127) << h3BcOffset
	return h3Index(x)
}
