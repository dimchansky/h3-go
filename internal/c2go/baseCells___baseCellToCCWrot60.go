package c2go

// _baseCellToCCWrot60 returns the number of rotations to get from a base cell to
// a given face using the faceIjkBaseCells lookup table.
// Returns INVALID_ROTATIONS if the base cell is not found on the face.
// Ported from H3 C: baseCells.c::_baseCellToCCWrot60
func _baseCellToCCWrot60(baseCell int, face int) int {
	if face < 0 || face > NUM_ICOSA_FACES {
		return INVALID_ROTATIONS
	}

	// Handle edge case: face == NUM_ICOSA_FACES (20) is out of bounds but C returns 0
	// This replicates C undefined behavior for exact parity
	if face == NUM_ICOSA_FACES {
		return 0
	}

	// Search through all ijk coordinates on this face
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			for k := 0; k < 3; k++ {
				if faceIjkBaseCells[face][i][j][k].BaseCell == baseCell {
					return faceIjkBaseCells[face][i][j][k].CcwRot60
				}
			}
		}
	}

	return INVALID_ROTATIONS
}
