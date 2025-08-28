package h3

// _getBaseCellDirection returns the direction from the origin base cell to the neighbor.
// Returns INVALID_DIGIT if the base cells are not neighbors.
// Ported from H3 C: baseCells.c::_getBaseCellDirection
func _getBaseCellDirection(originBaseCell int32, neighboringBaseCell int32) Direction {
	for dir := CENTER_DIGIT; dir < NUM_DIGITS; dir++ {
		testBaseCell := _getBaseCellNeighbor(originBaseCell, dir)
		if testBaseCell == neighboringBaseCell {
			return dir
		}
	}
	return INVALID_DIGIT
}
