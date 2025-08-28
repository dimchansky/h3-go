package h3

// cellToLocalIjk produces IJK+ coordinates for an index anchored by an origin.
//
// The coordinate space used by this function may have deleted regions or warping
// due to pentagonal distortion. Coordinates are only comparable if they come from
// the same origin index.
//
// Failure may occur if the index is too far away from the origin or if the index
// is on the other side of a pentagon.
//
// Ported from H3 C: localij.c::cellToLocalIjk
func cellToLocalIjk(origin H3Index, h3 H3Index, out *CoordIJK) H3Error {
	res := getResolution(origin)

	if res != getResolution(h3) {
		return E_RES_MISMATCH
	}

	originBaseCell := getBaseCell(origin)
	baseCell := getBaseCell(h3)

	// NEVER macro checks - base cells less than zero can not be represented in an index
	if originBaseCell < 0 || originBaseCell >= NUM_BASE_CELLS {
		return E_CELL_INVALID
	}
	if baseCell < 0 || baseCell >= NUM_BASE_CELLS {
		return E_CELL_INVALID
	}

	// Direction from origin base cell to index base cell
	dir := CENTER_DIGIT
	revDir := CENTER_DIGIT
	if originBaseCell != baseCell {
		dir = _getBaseCellDirection(originBaseCell, baseCell)
		if dir == INVALID_DIGIT {
			// Base cells are not neighbors, can't unfold
			return E_FAILED
		}
		revDir = _getBaseCellDirection(baseCell, originBaseCell)
		// assert(revDir != INVALID_DIGIT) - this should never happen if dir is valid
	}

	originOnPent := _isBaseCellPentagon(originBaseCell)
	indexOnPent := _isBaseCellPentagon(baseCell)

	var indexFijk FaceIJK
	if dir != CENTER_DIGIT {
		// Rotate index into the orientation of the origin base cell.
		// cw because we are undoing the rotation into that base cell.
		baseCellRotations := baseCellNeighbor60CCWRots[originBaseCell][dir]
		if indexOnPent {
			for i := 0; i < baseCellRotations; i++ {
				h3 = _h3RotatePent60cw(h3)

				revDir = _rotate60cw(revDir)
				if revDir == K_AXES_DIGIT {
					revDir = _rotate60cw(revDir)
				}
			}
		} else {
			for i := 0; i < baseCellRotations; i++ {
				h3 = _h3Rotate60cw(h3)

				revDir = _rotate60cw(revDir)
			}
		}
	}

	// Face is unused. This produces coordinates in base cell coordinate space.
	_h3ToFaceIjkWithInitializedFijk(h3, &indexFijk)

	if dir != CENTER_DIGIT {
		// assert(baseCell != originBaseCell)
		// assert(!(originOnPent && indexOnPent))

		var pentagonRotations int32 = 0
		var directionRotations int32 = 0

		if originOnPent {
			originLeadingDigit := _h3LeadingNonZeroDigit(origin)

			if originLeadingDigit == int32(INVALID_DIGIT) {
				return E_CELL_INVALID
			}
			if FAILED_DIRECTIONS[originLeadingDigit][dir] {
				// TODO: We may be unfolding the pentagon incorrectly in this
				// case; return an error code until this is guaranteed to be
				// correct.
				return E_FAILED
			}

			directionRotations = PENTAGON_ROTATIONS[originLeadingDigit][dir]
			pentagonRotations = directionRotations
		} else if indexOnPent {
			indexLeadingDigit := _h3LeadingNonZeroDigit(h3)

			if indexLeadingDigit == int32(INVALID_DIGIT) {
				return E_CELL_INVALID
			}
			if FAILED_DIRECTIONS[indexLeadingDigit][revDir] {
				// TODO: We may be unfolding the pentagon incorrectly in this
				// case; return an error code until this is guaranteed to be
				// correct.
				return E_FAILED
			}

			pentagonRotations = PENTAGON_ROTATIONS[revDir][indexLeadingDigit]
		}

		if pentagonRotations < 0 || directionRotations < 0 {
			// This occurs when an invalid K axis digit is present
			return E_CELL_INVALID
		}

		for i := int32(0); i < pentagonRotations; i++ {
			_ijkRotate60cw(&indexFijk.Coord)
		}

		var offset CoordIJK
		_neighbor(&offset, dir)
		// Scale offset based on resolution
		for r := res - 1; r >= 0; r-- {
			if isResolutionClassIII(r + 1) {
				// rotate ccw
				_downAp7(&offset)
			} else {
				// rotate cw
				_downAp7r(&offset)
			}
		}

		for i := int32(0); i < directionRotations; i++ {
			_ijkRotate60cw(&offset)
		}

		// Perform necessary translation
		_ijkAdd(&indexFijk.Coord, &offset, &indexFijk.Coord)
		_ijkNormalize(&indexFijk.Coord)
	} else if originOnPent && indexOnPent {
		// If the origin and index are on pentagon, and we checked that the base
		// cells are the same or neighboring, then they must be the same base
		// cell.
		// assert(baseCell == originBaseCell)

		originLeadingDigit := _h3LeadingNonZeroDigit(origin)
		indexLeadingDigit := _h3LeadingNonZeroDigit(h3)

		if originLeadingDigit == int32(INVALID_DIGIT) ||
			indexLeadingDigit == int32(INVALID_DIGIT) {
			return E_CELL_INVALID
		}
		if FAILED_DIRECTIONS[originLeadingDigit][indexLeadingDigit] {
			// TODO: We may be unfolding the pentagon incorrectly in this case;
			// return an error code until this is guaranteed to be correct.
			return E_FAILED
		}

		withinPentagonRotations :=
			PENTAGON_ROTATIONS[originLeadingDigit][indexLeadingDigit]

		for i := int32(0); i < withinPentagonRotations; i++ {
			_ijkRotate60cw(&indexFijk.Coord)
		}
	}

	*out = indexFijk.Coord
	return E_SUCCESS
}
