package c2go

// getResolution returns the H3 resolution (port of H3_GET_RESOLUTION)
// Ported from H3 C: h3Index.c::getResolution
func getResolution(h H3Index) int32 {
	return int32((uint64(h) & H3_RES_MASK) >> H3_RES_OFFSET)
}

// getBaseCell returns the base cell (port of H3_GET_BASE_CELL)
// Ported from H3 C: h3Index.c::H3_GET_BASE_CELL
func getBaseCell(h H3Index) int32 {
	return int32((uint64(h) & H3_BC_MASK) >> H3_BC_OFFSET)
}

// setResolution sets the H3 resolution (port of H3_SET_RESOLUTION)
// Ported from H3 C: h3Index.h::H3_SET_RESOLUTION
func setResolution(h H3Index, res int32) H3Index {
	x := uint64(h)
	x &^= H3_RES_MASK
	x |= (uint64(res) & 15) << H3_RES_OFFSET
	return H3Index(x)
}

// setBaseCell sets the base cell number (port of H3_SET_BASE_CELL)
// Ported from H3 C: h3Index.h::H3_SET_BASE_CELL
func setBaseCell(h H3Index, baseCell int32) H3Index {
	x := uint64(h)
	x &^= H3_BC_MASK
	x |= (uint64(baseCell) & 127) << H3_BC_OFFSET
	return H3Index(x)
}
