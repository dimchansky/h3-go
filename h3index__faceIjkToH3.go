package h3

// _faceIjkToH3 converts a faceIJK address to an h3Index.
// Ported from H3 C: h3Index.c::_faceIjkToH3.
func _faceIjkToH3(fijk *faceIJK, res int32) h3Index {
	// initialize the index
	h := h3Index(h3Init)
	h = setMode(h, h3CellMode)
	// Set resolution
	h = setResolution(h, res)

	// check for res 0/base cell
	if res == 0 {
		if fijk.Coord.I > maxFaceCoord || fijk.Coord.J > maxFaceCoord ||
			fijk.Coord.K > maxFaceCoord {
			// out of range input
			return h3Null
		}

		h = setBaseCell(h, _faceIjkToBaseCell(fijk))
		return h
	}

	// we need to find the correct base cell faceIJK for this H3 index;
	// start with the passed in face and resolution res ijk coordinates
	// in that face's coordinate system
	fijkBC := *fijk

	// build the h3Index from finest res up
	// adjust r for the fact that the res 0 base cell offsets the indexing
	// digits
	ijk := &fijkBC.Coord
	for r := res - 1; r >= 0; r-- {
		lastIJK := *ijk
		var lastCenter coordIJK
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

		var diff coordIJK
		_ijkSub(&lastIJK, &lastCenter, &diff)
		_ijkNormalize(&diff)

		h = setIndexDigit(h, r+1, int32(_unitIjkToDigit(&diff)))
	}

	// fijkBC should now hold the IJK of the base cell in the
	// coordinate system of the current face

	if fijkBC.Coord.I > maxFaceCoord || fijkBC.Coord.J > maxFaceCoord ||
		fijkBC.Coord.K > maxFaceCoord {
		// out of range input
		return h3Null
	}

	// lookup the correct base cell
	baseCell := _faceIjkToBaseCell(&fijkBC)
	h = setBaseCell(h, baseCell)

	// rotate if necessary to get canonical base cell orientation
	// for this base cell
	numRots := _faceIjkToBaseCellCCWrot60(&fijkBC)
	if _isBaseCellPentagon(baseCell) {
		// force rotation out of missing k-axes sub-sequence
		if _h3LeadingNonZeroDigit(h) == int32(kAxesDigit) {
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
