package h3

// getRes0Cells populates the provided slice with all base cells (resolution 0 cells).
// This is the dst-buffer pattern version that returns all 122 base cells.
// Ported from H3 C: baseCells.c::getRes0Cells.
func getRes0Cells(out []H3Index) H3Error {
	if len(out) != NUM_BASE_CELLS {
		return E_FAILED // Need exactly NUM_BASE_CELLS slots
	}

	for bc := int32(0); bc < NUM_BASE_CELLS; bc++ {
		baseCell := H3Index(H3_INIT)
		baseCell = setMode(baseCell, H3_CELL_MODE)
		baseCell = setBaseCell(baseCell, bc)
		out[bc] = baseCell
	}
	return E_SUCCESS
}
