package h3

// getRes0Cells populates the provided slice with all base cells (resolution 0 cells).
// This is the dst-buffer pattern version that returns all 122 base cells.
// Ported from H3 C: baseCells.c::getRes0Cells.
func getRes0Cells(out []h3Index) h3Error {
	if len(out) != numBaseCells {
		return eFailed // Need exactly numBaseCells slots
	}

	for bc := int32(0); bc < numBaseCells; bc++ {
		baseCell := h3Index(h3Init)
		baseCell = setMode(baseCell, h3CellMode)
		baseCell = setBaseCell(baseCell, bc)
		out[bc] = baseCell
	}
	return eSuccess
}
