package h3

// _isBaseCellPentagon returns true if the base cell is a pentagon, else false.
// Ported from baseCells.c::_isBaseCellPentagon using the local lookup table.
// Ported from H3 C: baseCells.c::_isBaseCellPentagon.
func _isBaseCellPentagon(baseCell int32) bool {
	if baseCell < 0 || int(baseCell) >= len(isBaseCellPentagonArr) {
		return false
	}
	return isBaseCellPentagonArr[baseCell]
}
