package c2go

// _baseCellIsCwOffset returns true if the given face is a clockwise offset face for the given pentagon base cell.
// For non-pentagon base cells, this always returns false.
// Ported from H3 C: baseCells.c::_baseCellIsCwOffset
func _baseCellIsCwOffset(baseCell int, testFace int) bool {
	return baseCellData[baseCell].CwOffsetPent[0] == testFace ||
		baseCellData[baseCell].CwOffsetPent[1] == testFace
}
