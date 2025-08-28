package h3

// _isBaseCellPolarPentagon returns true if the base cell is a polar pentagon.
// Ported from H3 C: baseCells.c::_isBaseCellPolarPentagon
func _isBaseCellPolarPentagon(baseCell int32) bool {
	return baseCell == 4 || baseCell == 117
}
