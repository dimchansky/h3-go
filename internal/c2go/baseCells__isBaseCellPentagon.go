package c2go

// _isBaseCellPentagon returns 1 if the base cell is a pentagon, else 0.
// Ported from baseCells.c::_isBaseCellPentagon using the local lookup table.
// Ported from H3 C: baseCells.c::isBaseCellPentagon
func _isBaseCellPentagon(baseCell int) int {
	if baseCell < 0 || baseCell >= len(isBaseCellPentagonArr) {
		return 0
	}
	if isBaseCellPentagonArr[baseCell] {
		return 1
	}
	return 0
}
