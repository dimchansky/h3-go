package h3

// localIjkToCell produces an H3 index from IJK+ coordinates anchored by an origin.
// The coordinate space used by this function may have deleted regions or warping
// due to pentagonal distortion. Failure may occur if the coordinates are too far
// away from the origin or if the index is on the other side of a pentagon.
// Ported from H3 C: localij.c::localIjkToCell
func localIjkToCell(origin H3Index, ijk *CoordIJK, out *H3Index) H3Error {
	res := getResolution(origin)
	originBaseCell := getBaseCell(origin)
	if originBaseCell < 0 || originBaseCell >= NUM_BASE_CELLS {
		return E_CELL_INVALID
	}
	originOnPent := _isBaseCellPentagon(originBaseCell)

	// This logic is very similar to faceIjkToH3
	// initialize the index
	*out = H3Index(H3_INIT)
	*out = setMode(*out, H3_CELL_MODE)
	*out = setResolution(*out, res)

	// check for res 0/base cell
	if res == 0 {
		dir := _unitIjkToDigit(ijk)
		if dir == INVALID_DIGIT {
			// out of range input - not a unit vector or zero vector
			return E_FAILED
		}

		newBaseCell := _getBaseCellNeighbor(originBaseCell, dir)
		if newBaseCell == INVALID_BASE_CELL {
			// Moving in an invalid direction off a pentagon.
			return E_FAILED
		}
		*out = setBaseCell(*out, newBaseCell)
		return E_SUCCESS
	}

	// we need to find the correct base cell offset (if any) for this H3 index;
	// start with the passed in base cell and resolution res ijk coordinates
	// in that base cell's coordinate system
	ijkCopy := *ijk

	// build the H3Index from finest res up
	// adjust r for the fact that the res 0 base cell offsets the indexing
	// digits
	for r := res - 1; r >= 0; r-- {
		lastIJK := ijkCopy
		var lastCenter CoordIJK
		if isResolutionClassIII(r + 1) {
			// rotate ccw
			upAp7Error := _upAp7Checked(&ijkCopy)
			if upAp7Error != E_SUCCESS {
				return upAp7Error
			}
			lastCenter = ijkCopy
			_downAp7(&lastCenter)
		} else {
			// rotate cw
			upAp7rError := _upAp7rChecked(&ijkCopy)
			if upAp7rError != E_SUCCESS {
				return upAp7rError
			}
			lastCenter = ijkCopy
			_downAp7r(&lastCenter)
		}

		var diff CoordIJK
		_ijkSub(&lastIJK, &lastCenter, &diff)
		_ijkNormalize(&diff)

		*out = setIndexDigit(*out, r+1, int32(_unitIjkToDigit(&diff)))
	}

	// ijkCopy should now hold the IJK of the base cell in the
	// coordinate system of the current base cell

	if ijkCopy.I > 1 || ijkCopy.J > 1 || ijkCopy.K > 1 {
		// out of range input
		return E_FAILED
	}

	// lookup the correct base cell
	dir := _unitIjkToDigit(&ijkCopy)
	baseCell := _getBaseCellNeighbor(originBaseCell, dir)
	// If baseCell is invalid, it must be because the origin base cell is a
	// pentagon, and because pentagon base cells do not border each other,
	// baseCell must not be a pentagon.
	var indexOnPent bool
	if baseCell == INVALID_BASE_CELL {
		indexOnPent = false
	} else {
		indexOnPent = _isBaseCellPentagon(baseCell)
	}

	if dir != CENTER_DIGIT {
		// If the index is in a warped direction, we need to unwarp the base
		// cell direction. There may be further need to rotate the index digits.
		var pentagonRotations int32 = 0
		if originOnPent {
			originLeadingDigit := _h3LeadingNonZeroDigit(origin)
			if originLeadingDigit == int32(INVALID_DIGIT) {
				return E_CELL_INVALID
			}
			pentagonRotations = PENTAGON_ROTATIONS_REVERSE[originLeadingDigit][dir]
			for i := int32(0); i < pentagonRotations; i++ {
				dir = _rotate60ccw(dir)
			}
			// The pentagon rotations are being chosen so that dir is not the
			// deleted direction. If it still happens, it means we're moving
			// into a deleted subsequence, so there is no index here.
			if dir == K_AXES_DIGIT {
				return E_PENTAGON
			}
			baseCell = _getBaseCellNeighbor(originBaseCell, dir)

			// indexOnPent does not need to be checked again since no pentagon
			// base cells border each other.
			if baseCell == INVALID_BASE_CELL {
				return E_CELL_INVALID
			}
			indexOnPent = _isBaseCellPentagon(baseCell)
			if indexOnPent {
				return E_CELL_INVALID
			}
		}

		// Now we can determine the relation between the origin and target base
		// cell.
		baseCellRotations := int32(baseCellNeighbor60CCWRots[originBaseCell][dir])
		if baseCellRotations < 0 {
			return E_CELL_INVALID
		}

		// Adjust for pentagon warping within the base cell. The base cell
		// should be in the right location, so now we need to rotate the index
		// back. We might not need to check for errors since we would just be
		// double mapping.
		if indexOnPent {
			revDir := _getBaseCellDirection(baseCell, originBaseCell)
			if revDir == INVALID_DIGIT {
				return E_CELL_INVALID
			}

			// Adjust for the different coordinate space in the two base cells.
			// This is done first because we need to do the pentagon rotations
			// based on the leading digit in the pentagon's coordinate system.
			for i := int32(0); i < baseCellRotations; i++ {
				*out = _h3Rotate60ccw(*out)
			}

			indexLeadingDigit := _h3LeadingNonZeroDigit(*out)
			// This case should be unreachable because this function is building
			// *out, and should never generate an invalid digit, above.
			if indexLeadingDigit == int32(INVALID_DIGIT) {
				return E_CELL_INVALID
			}
			if _isBaseCellPolarPentagon(baseCell) {
				pentagonRotations = PENTAGON_ROTATIONS_REVERSE_POLAR[revDir][indexLeadingDigit]
			} else {
				pentagonRotations = PENTAGON_ROTATIONS_REVERSE_NONPOLAR[revDir][indexLeadingDigit]
			}
			// For this to occur, revDir would need to be 1. Since revDir is
			// from the index base cell (which is a pentagon) towards the
			// origin, this should never be the case.
			if pentagonRotations < 0 {
				return E_CELL_INVALID
			}

			for i := int32(0); i < pentagonRotations; i++ {
				*out = _h3RotatePent60ccw(*out)
			}
		} else {
			if pentagonRotations < 0 {
				return E_CELL_INVALID
			}
			for i := int32(0); i < pentagonRotations; i++ {
				*out = _h3Rotate60ccw(*out)
			}

			// Adjust for the different coordinate space in the two base cells.
			for i := int32(0); i < baseCellRotations; i++ {
				*out = _h3Rotate60ccw(*out)
			}
		}
	} else if originOnPent && indexOnPent {
		originLeadingDigit := _h3LeadingNonZeroDigit(origin)
		indexLeadingDigit := _h3LeadingNonZeroDigit(*out)

		if originLeadingDigit == int32(INVALID_DIGIT) || indexLeadingDigit == int32(INVALID_DIGIT) {
			return E_CELL_INVALID
		}
		withinPentagonRotations := PENTAGON_ROTATIONS_REVERSE[originLeadingDigit][indexLeadingDigit]
		if withinPentagonRotations < 0 {
			// This occurs when an invalid K axis digit is present
			return E_CELL_INVALID
		}

		for i := int32(0); i < withinPentagonRotations; i++ {
			*out = _h3Rotate60ccw(*out)
		}
	}

	if indexOnPent {
		// TODO: There are cases in cellToLocalIjk which are failed but not
		// accounted for here - instead just fail if the recovered index is
		// invalid.
		if _h3LeadingNonZeroDigit(*out) == int32(K_AXES_DIGIT) {
			return E_PENTAGON
		}
	}

	*out = setBaseCell(*out, baseCell)
	return E_SUCCESS
}
