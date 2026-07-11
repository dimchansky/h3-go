package h3

// _baseCellToFaceIjk converts a base cell number to the faceIJK address of that base cell.
// This is a simple lookup in the baseCellData table: h = baseCellData[baseCell].homeFijk
// Ported from H3 C: baseCells.c::_baseCellToFaceIjk.
func _baseCellToFaceIjk(baseCell int32, h *faceIJK) {
	*h = baseCellData[baseCell].HomeFijk
}
