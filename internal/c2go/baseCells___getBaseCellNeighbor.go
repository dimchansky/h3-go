package c2go

// _getBaseCellNeighbor returns the neighboring base cell in the given direction.
// Returns the base cell number of the neighbor, or INVALID_BASE_CELL if there
// is no neighbor in that direction (for pentagon base cells in certain directions).
// Ported from H3 C: baseCells.c::_getBaseCellNeighbor
func _getBaseCellNeighbor(baseCell int32, dir Direction) int32 {
	return baseCellNeighbors[baseCell][dir]
}
