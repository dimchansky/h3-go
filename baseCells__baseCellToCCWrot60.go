package h3

// _baseCellToCCWrot60 returns the number of rotations to get from a base cell to
// a given face using the faceIjkBaseCells lookup table.
// Returns invalidRotations if the base cell is not found on the face.
// Ported from H3 C: baseCells.c::_baseCellToCCWrot60.
func _baseCellToCCWrot60(baseCell int32, face int32) int32 {
	if face < 0 || face > numIcosaFaces {
		return invalidRotations
	}

	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			for k := 0; k < 3; k++ {
				if faceIjkBaseCells[face][i][j][k].BaseCell == baseCell {
					return faceIjkBaseCells[face][i][j][k].CcwRot60
				}
			}
		}
	}

	return invalidRotations
}
