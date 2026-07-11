package h3

// res0CellCount returns the number of resolution 0 base cells.
// Ported from H3 C: baseCells.c::res0CellCount.
func res0CellCount() int32 {
	return NUM_BASE_CELLS
}
