package h3

// localIjkToCell produces an H3 index from IJK+ coordinates anchored by an origin.
// The coordinate space used by this function may have deleted regions or warping
// due to pentagonal distortion. Failure may occur if the coordinates are too far
// away from the origin or if the index is on the other side of a pentagon.
// Ported from H3 C: localij.c::localIjkToCell.
func localIjkToCell(origin h3Index, ijk *coordIJK, out *h3Index) h3Error {
	res := getResolution(origin)
	originBaseCell := getBaseCell(origin)
	if originBaseCell < 0 || originBaseCell >= numBaseCells {
		return eCellInvalid
	}
	originOnPent := _isBaseCellPentagon(originBaseCell)

	// This logic is very similar to faceIjkToH3
	// initialize the index
	*out = h3Index(h3Init)
	*out = setMode(*out, h3CellMode)
	*out = setResolution(*out, res)

	// check for res 0/base cell
	if res == 0 {
		dir := _unitIjkToDigit(ijk)
		if dir == invalidDigit {
			// out of range input - not a unit vector or zero vector
			return eFailed
		}

		newBaseCell := _getBaseCellNeighbor(originBaseCell, dir)
		if newBaseCell == invalidBaseCell {
			// Moving in an invalid direction off a pentagon.
			return eFailed
		}
		*out = setBaseCell(*out, newBaseCell)
		return eSuccess
	}

	// we need to find the correct base cell offset (if any) for this H3 index;
	// start with the passed in base cell and resolution res ijk coordinates
	// in that base cell's coordinate system
	ijkCopy := *ijk

	// build the h3Index from finest res up
	// adjust r for the fact that the res 0 base cell offsets the indexing
	// digits
	for r := res - 1; r >= 0; r-- {
		lastIJK := ijkCopy
		var lastCenter coordIJK
		if isResolutionClassIII(r + 1) {
			// rotate ccw
			upAp7Error := _upAp7Checked(&ijkCopy)
			if upAp7Error != eSuccess {
				return upAp7Error
			}
			lastCenter = ijkCopy
			_downAp7(&lastCenter)
		} else {
			// rotate cw
			upAp7rError := _upAp7rChecked(&ijkCopy)
			if upAp7rError != eSuccess {
				return upAp7rError
			}
			lastCenter = ijkCopy
			_downAp7r(&lastCenter)
		}

		var diff coordIJK
		_ijkSub(&lastIJK, &lastCenter, &diff)
		_ijkNormalize(&diff)

		*out = h3SetIndexDigit(*out, r+1, int32(_unitIjkToDigit(&diff)))
	}

	// ijkCopy should now hold the IJK of the base cell in the
	// coordinate system of the current base cell

	if ijkCopy.I > 1 || ijkCopy.J > 1 || ijkCopy.K > 1 {
		// out of range input
		return eFailed
	}

	// lookup the correct base cell
	dir := _unitIjkToDigit(&ijkCopy)
	baseCell := _getBaseCellNeighbor(originBaseCell, dir)
	// If baseCell is invalid, it must be because the origin base cell is a
	// pentagon, and because pentagon base cells do not border each other,
	// baseCell must not be a pentagon.
	var indexOnPent bool
	if baseCell == invalidBaseCell {
		indexOnPent = false
	} else {
		indexOnPent = _isBaseCellPentagon(baseCell)
	}

	if dir != centerDigit {
		// If the index is in a warped direction, we need to unwarp the base
		// cell direction. There may be further need to rotate the index digits.
		var pentagonRotations int32 = 0
		if originOnPent {
			originLeadingDigit := _h3LeadingNonZeroDigit(origin)
			if originLeadingDigit == int32(invalidDigit) {
				return eCellInvalid
			}
			pentagonRotations = pentagonRotationsReverse[originLeadingDigit][dir]
			for i := int32(0); i < pentagonRotations; i++ {
				dir = _rotate60ccw(dir)
			}
			// The pentagon rotations are being chosen so that dir is not the
			// deleted direction. If it still happens, it means we're moving
			// into a deleted subsequence, so there is no index here.
			if dir == kAxesDigit {
				return ePentagon
			}
			baseCell = _getBaseCellNeighbor(originBaseCell, dir)

			// indexOnPent does not need to be checked again since no pentagon
			// base cells border each other.
			if baseCell == invalidBaseCell {
				return eCellInvalid
			}
			indexOnPent = _isBaseCellPentagon(baseCell)
			if indexOnPent {
				return eCellInvalid
			}
		}

		// Now we can determine the relation between the origin and target base
		// cell.
		baseCellRotations := int32(baseCellNeighbor60CCWRots[originBaseCell][dir])
		if baseCellRotations < 0 {
			return eCellInvalid
		}

		// Adjust for pentagon warping within the base cell. The base cell
		// should be in the right location, so now we need to rotate the index
		// back. We might not need to check for errors since we would just be
		// double mapping.
		if indexOnPent {
			revDir := _getBaseCellDirection(baseCell, originBaseCell)
			if revDir == invalidDigit {
				return eCellInvalid
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
			if indexLeadingDigit == int32(invalidDigit) {
				return eCellInvalid
			}
			if _isBaseCellPolarPentagon(baseCell) {
				pentagonRotations = pentagonRotationsReversePolar[revDir][indexLeadingDigit]
			} else {
				pentagonRotations = pentagonRotationsReverseNonpolar[revDir][indexLeadingDigit]
			}
			// For this to occur, revDir would need to be 1. Since revDir is
			// from the index base cell (which is a pentagon) towards the
			// origin, this should never be the case.
			if pentagonRotations < 0 {
				return eCellInvalid
			}

			for i := int32(0); i < pentagonRotations; i++ {
				*out = _h3RotatePent60ccw(*out)
			}
		} else {
			if pentagonRotations < 0 {
				return eCellInvalid
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

		if originLeadingDigit == int32(invalidDigit) || indexLeadingDigit == int32(invalidDigit) {
			return eCellInvalid
		}
		withinPentagonRotations := pentagonRotationsReverse[originLeadingDigit][indexLeadingDigit]
		if withinPentagonRotations < 0 {
			// This occurs when an invalid K axis digit is present
			return eCellInvalid
		}

		for i := int32(0); i < withinPentagonRotations; i++ {
			*out = _h3Rotate60ccw(*out)
		}
	}

	if indexOnPent {
		// TODO: There are cases in cellToLocalIjk which are failed but not
		// accounted for here - instead just fail if the recovered index is
		// invalid.
		if _h3LeadingNonZeroDigit(*out) == int32(kAxesDigit) {
			return ePentagon
		}
	}

	*out = setBaseCell(*out, baseCell)
	return eSuccess
}
