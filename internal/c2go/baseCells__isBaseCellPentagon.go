package c2go

// _isBaseCellPentagon returns true if the base cell is a pentagon, else false.
// Ported from baseCells.c::_isBaseCellPentagon using the local lookup table.
// Ported from H3 C: baseCells.c::isBaseCellPentagon
func _isBaseCellPentagon(baseCell int) bool {
	if baseCell < 0 || baseCell >= len(isBaseCellPentagonArr) {
		return false
	}
	return isBaseCellPentagonArr[baseCell]
}
