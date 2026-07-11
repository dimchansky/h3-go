package h3

// _h3ToFaceIjkWithInitializedFijk converts h3Index to faceIJK address with initialized face.
// The faceIJK address should be initialized with the desired face and normalized base cell coordinates.
// Returns 1 if the possibility of overage exists, otherwise 0.
// Ported from H3 C: h3Index.c::_h3ToFaceIjkWithInitializedFijk.
func _h3ToFaceIjkWithInitializedFijk(h h3Index, fijk *faceIJK) int32 {
	ijk := &fijk.Coord
	res := getResolution(h)

	// Center base cell hierarchy is entirely on this face
	possibleOverage := int32(1)
	if !_isBaseCellPentagon(getBaseCell(h)) &&
		(res == 0 ||
			(fijk.Coord.I == 0 && fijk.Coord.J == 0 && fijk.Coord.K == 0)) {
		possibleOverage = 0
	}

	for r := int32(1); r <= res; r++ {
		if isResolutionClassIII(r) {
			// Class III == rotate ccw
			_downAp7(ijk)
		} else {
			// Class II == rotate cw
			_downAp7r(ijk)
		}

		_neighbor(ijk, direction(getIndexDigit(h, r)))
	}

	return possibleOverage
}
