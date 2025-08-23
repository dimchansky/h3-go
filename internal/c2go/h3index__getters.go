package c2go

// getResolution returns the H3 resolution (port of H3_GET_RESOLUTION)
// Ported from H3 C: h3Index.c::getResolution
func getResolution(h H3Index) int {
	return int((uint64(h) & H3_RES_MASK) >> H3_RES_OFFSET)
}

// getBaseCellNumber returns the base cell (port of H3_GET_BASE_CELL)
// Ported from H3 C: h3Index.c::getBaseCellNumber
func getBaseCellNumber(h H3Index) int {
	return int((uint64(h) & H3_BC_MASK) >> H3_BC_OFFSET)
}
