package h3

// _h3ToFaceIjk converts an H3 index to a FaceIJK address.
// This handles coordinate transformations and overage adjustments.
// Ported from H3 C: h3Index.c::_h3ToFaceIjk.
func _h3ToFaceIjk(h H3Index, fijk *FaceIJK) H3Error {
	baseCell := getBaseCell(h)
	if baseCell < 0 || baseCell >= NUM_BASE_CELLS {
		// Base cells less than zero can not be represented in an index
		// To prevent reading uninitialized memory, we zero the output.
		fijk.Face = 0
		fijk.Coord.I = 0
		fijk.Coord.J = 0
		fijk.Coord.K = 0
		return E_CELL_INVALID
	}

	// adjust for the pentagonal missing sequence; all of sub-sequence 5 needs
	// to be adjusted (and some of sub-sequence 4 below)
	if _isBaseCellPentagon(baseCell) && _h3LeadingNonZeroDigit(h) == 5 {
		h = _h3Rotate60cw(h)
	}

	// start with the "home" face and ijk+ coordinates for the base cell of c
	*fijk = baseCellData[baseCell].HomeFijk
	if _h3ToFaceIjkWithInitializedFijk(h, fijk) == 0 {
		return E_SUCCESS // no overage is possible; h lies on this face
	}

	// if we're here we have the potential for an "overage"; i.e., it is
	// possible that c lies on an adjacent face

	origIJK := fijk.Coord

	// if we're in Class III, drop into the next finer Class II grid
	res := getResolution(h)
	if isResolutionClassIII(res) {
		// Class III
		_downAp7r(&fijk.Coord)
		res++
	}

	// adjust for overage if needed
	// a pentagon base cell with a leading 4 digit requires special handling
	pentLeading4 := (_isBaseCellPentagon(baseCell) && _h3LeadingNonZeroDigit(h) == 4)
	if _adjustOverageClassII(fijk, res, pentLeading4, false) != NO_OVERAGE {
		// if the base cell is a pentagon we have the potential for secondary
		// overages
		if _isBaseCellPentagon(baseCell) {
			for _adjustOverageClassII(fijk, res, false, false) != NO_OVERAGE {
				continue
			}
		}

		if res != getResolution(h) {
			_upAp7r(&fijk.Coord)
		}
	} else if res != getResolution(h) {
		fijk.Coord = origIJK
	}
	return E_SUCCESS
}
