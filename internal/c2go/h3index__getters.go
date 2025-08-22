package c2go

// getResolution returns the H3 resolution (port of H3_GET_RESOLUTION)
func getResolution(h H3Index) int {
    const H3_RES_OFFSET = 52
    const H3_RES_MASK = uint64(15) << H3_RES_OFFSET
    return int((uint64(h) & H3_RES_MASK) >> H3_RES_OFFSET)
}

// getBaseCellNumber returns the base cell (port of H3_GET_BASE_CELL)
func getBaseCellNumber(h H3Index) int {
    const H3_BC_OFFSET = 45
    const H3_BC_MASK = uint64(127) << H3_BC_OFFSET
    return int((uint64(h) & H3_BC_MASK) >> H3_BC_OFFSET)
}

