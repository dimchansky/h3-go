package c2go

// _faceIjkToH3 converts a FaceIJK address to an H3Index.
// Ported from H3 C: h3Index.c::_faceIjkToH3
func _faceIjkToH3(fijk *FaceIJK, res int32) H3Index {
	// initialize the index
	h := H3Index(H3_INIT)
	h = setMode(h, H3_CELL_MODE)
	// Set resolution (inline H3_SET_RESOLUTION)
	x := uint64(h)
	x &^= H3_RES_MASK
	x |= (uint64(res) & 15) << H3_RES_OFFSET
	h = H3Index(x)

	// check for res 0/base cell
	if res == 0 {
		if fijk.Coord.I > MAX_FACE_COORD || fijk.Coord.J > MAX_FACE_COORD ||
			fijk.Coord.K > MAX_FACE_COORD {
			// out of range input
			return H3_NULL
		}

		h = H3Index(setBaseCell(uint64(h), _faceIjkToBaseCell(fijk)))
		return h
	}

	// we need to find the correct base cell FaceIJK for this H3 index;
	// start with the passed in face and resolution res ijk coordinates
	// in that face's coordinate system
	fijkBC := *fijk

	// build the H3Index from finest res up
	// adjust r for the fact that the res 0 base cell offsets the indexing
	// digits
	ijk := &fijkBC.Coord
	for r := res - 1; r >= 0; r-- {
		lastIJK := *ijk
		var lastCenter CoordIJK
		if isResolutionClassIII(r + 1) {
			// rotate ccw
			_upAp7(ijk)
			lastCenter = *ijk
			_downAp7(&lastCenter)
		} else {
			// rotate cw
			_upAp7r(ijk)
			lastCenter = *ijk
			_downAp7r(&lastCenter)
		}

		var diff CoordIJK
		_ijkSub(&lastIJK, &lastCenter, &diff)
		_ijkNormalize(&diff)

		h = setIndexDigit(h, r+1, int32(_unitIjkToDigit(&diff)))
	}

	// fijkBC should now hold the IJK of the base cell in the
	// coordinate system of the current face

	if fijkBC.Coord.I > MAX_FACE_COORD || fijkBC.Coord.J > MAX_FACE_COORD ||
		fijkBC.Coord.K > MAX_FACE_COORD {
		// out of range input
		return H3_NULL
	}

	// lookup the correct base cell
	baseCell := _faceIjkToBaseCell(&fijkBC)
	h = H3Index(setBaseCell(uint64(h), baseCell))

	// rotate if necessary to get canonical base cell orientation
	// for this base cell
	numRots := _faceIjkToBaseCellCCWrot60(&fijkBC)
	if _isBaseCellPentagon(baseCell) {
		// force rotation out of missing k-axes sub-sequence
		if _h3LeadingNonZeroDigit(h) == int32(K_AXES_DIGIT) {
			// check for a cw/ccw offset face; default is ccw
			if _baseCellIsCwOffset(baseCell, fijkBC.Face) {
				h = _h3Rotate60cw(h)
			} else {
				h = _h3Rotate60ccw(h)
			}
		}
		for i := int32(0); i < numRots; i++ {
			h = _h3RotatePent60ccw(h)
		}
	} else {
		for i := int32(0); i < numRots; i++ {
			h = _h3Rotate60ccw(h)
		}
	}

	return h
}
