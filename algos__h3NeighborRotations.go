package h3

// h3NeighborRotations returns the hexagon index neighboring the origin, in the direction dir.
//
// Implementation note: The only reachable case where this returns 0 is if the
// origin is a pentagon and the translation is in the k direction. Thus,
// 0 can only be returned if origin is a pentagon.
//
// This function traverses the H3 index structure by modifying digits at each resolution
// level, handling base cell transitions, and managing pentagon special cases with
// missing K-subsequences.
//
// Ported from H3 C: algos.c::h3NeighborRotations
func h3NeighborRotations(origin H3Index, dir Direction, rotations *int32, out *H3Index) H3Error {
	current := origin

	if dir < CENTER_DIGIT || dir >= INVALID_DIGIT {
		return E_FAILED
	}

	// Ensure that rotations is modulo'd by 6 before any possible addition,
	// to protect against signed integer overflow.
	*rotations = *rotations % 6
	for i := int32(0); i < *rotations; i++ {
		dir = _rotate60ccw(dir)
	}

	newRotations := int32(0)
	oldBaseCell := int32(getBaseCell(current))
	if oldBaseCell < 0 || oldBaseCell >= NUM_BASE_CELLS {
		// Base cells less than zero can not be represented in an index
		return E_CELL_INVALID
	}
	oldLeadingDigit := _h3LeadingNonZeroDigit(current)

	// Adjust the indexing digits and, if needed, the base cell.
	r := getResolution(current) - 1
	for {
		if r == -1 {
			current = setBaseCell(current, baseCellNeighbors[oldBaseCell][dir])
			newRotations = int32(baseCellNeighbor60CCWRots[oldBaseCell][dir])

			if getBaseCell(current) == INVALID_BASE_CELL {
				// Adjust for the deleted k vertex at the base cell level.
				// This edge actually borders a different neighbor.
				current = setBaseCell(current, baseCellNeighbors[oldBaseCell][IK_AXES_DIGIT])
				newRotations = int32(baseCellNeighbor60CCWRots[oldBaseCell][IK_AXES_DIGIT])

				// perform the adjustment for the k-subsequence we're skipping
				// over.
				current = _h3Rotate60ccw(current)
				*rotations = *rotations + 1
			}

			break
		} else {
			oldDigit := Direction(getIndexDigit(current, r+1))
			var nextDir Direction
			if oldDigit == INVALID_DIGIT {
				// Only possible on invalid input
				return E_CELL_INVALID
			} else if isResolutionClassIII(r + 1) {
				current = setIndexDigit(current, r+1, int32(NEW_DIGIT_II[oldDigit][dir]))
				nextDir = NEW_ADJUSTMENT_II[oldDigit][dir]
			} else {
				current = setIndexDigit(current, r+1, int32(NEW_DIGIT_III[oldDigit][dir]))
				nextDir = NEW_ADJUSTMENT_III[oldDigit][dir]
			}

			if nextDir != CENTER_DIGIT {
				dir = nextDir
				r--
			} else {
				// No more adjustment to perform
				break
			}
		}
	}

	newBaseCell := int32(getBaseCell(current))
	if _isBaseCellPentagon(newBaseCell) {
		alreadyAdjustedKSubsequence := false

		// force rotation out of missing k-axes sub-sequence
		if Direction(_h3LeadingNonZeroDigit(current)) == K_AXES_DIGIT {
			if oldBaseCell != newBaseCell {
				// in this case, we traversed into the deleted
				// k subsequence of a pentagon base cell.
				// We need to rotate out of that case depending
				// on how we got here.
				// check for a cw/ccw offset face; default is ccw

				if _baseCellIsCwOffset(newBaseCell, baseCellData[oldBaseCell].HomeFijk.Face) {
					current = _h3Rotate60cw(current)
				} else {
					// See cwOffsetPent in testGridDisk.c for why this is
					// unreachable.
					current = _h3Rotate60ccw(current)
				}
				alreadyAdjustedKSubsequence = true
			} else {
				// In this case, we traversed into the deleted
				// k subsequence from within the same pentagon
				// base cell.
				if Direction(oldLeadingDigit) == CENTER_DIGIT {
					// Undefined: the k direction is deleted from here
					return E_PENTAGON
				} else if Direction(oldLeadingDigit) == JK_AXES_DIGIT {
					// Rotate out of the deleted k subsequence
					// We also need an additional change to the direction we're
					// moving in
					current = _h3Rotate60ccw(current)
					*rotations = *rotations + 1
				} else if Direction(oldLeadingDigit) == IK_AXES_DIGIT {
					// Rotate out of the deleted k subsequence
					// We also need an additional change to the direction we're
					// moving in
					current = _h3Rotate60cw(current)
					*rotations = *rotations + 5
				} else {
					// TODO: Should never occur, but is reachable by fuzzer
					return E_FAILED
				}
			}
		}

		for i := int32(0); i < newRotations; i++ {
			current = _h3RotatePent60ccw(current)
		}

		// Account for differing orientation of the base cells (this edge
		// might not follow properties of some other edges.)
		if oldBaseCell != newBaseCell {
			if _isBaseCellPolarPentagon(newBaseCell) {
				// 'polar' base cells behave differently because they have all
				// i neighbors.
				if oldBaseCell != 118 && oldBaseCell != 8 &&
					Direction(_h3LeadingNonZeroDigit(current)) != JK_AXES_DIGIT {
					*rotations = *rotations + 1
				}
			} else if Direction(_h3LeadingNonZeroDigit(current)) == IK_AXES_DIGIT &&
				!alreadyAdjustedKSubsequence {
				// account for distortion introduced to the 5 neighbor by the
				// deleted k subsequence.
				*rotations = *rotations + 1
			}
		}
	} else {
		for i := int32(0); i < newRotations; i++ {
			current = _h3Rotate60ccw(current)
		}
	}

	*rotations = (*rotations + newRotations) % 6
	*out = current

	return E_SUCCESS
}

// Lookup tables for neighbor traversal algorithm

// NEW_DIGIT_II: New digit when traversing along class II grids.
// Current digit -> direction -> new digit.
var NEW_DIGIT_II = [7][7]Direction{
	{CENTER_DIGIT, K_AXES_DIGIT, J_AXES_DIGIT, JK_AXES_DIGIT, I_AXES_DIGIT,
		IK_AXES_DIGIT, IJ_AXES_DIGIT},
	{K_AXES_DIGIT, I_AXES_DIGIT, JK_AXES_DIGIT, IJ_AXES_DIGIT, IK_AXES_DIGIT,
		J_AXES_DIGIT, CENTER_DIGIT},
	{J_AXES_DIGIT, JK_AXES_DIGIT, K_AXES_DIGIT, I_AXES_DIGIT, IJ_AXES_DIGIT,
		CENTER_DIGIT, IK_AXES_DIGIT},
	{JK_AXES_DIGIT, IJ_AXES_DIGIT, I_AXES_DIGIT, IK_AXES_DIGIT, CENTER_DIGIT,
		K_AXES_DIGIT, J_AXES_DIGIT},
	{I_AXES_DIGIT, IK_AXES_DIGIT, IJ_AXES_DIGIT, CENTER_DIGIT, J_AXES_DIGIT,
		JK_AXES_DIGIT, K_AXES_DIGIT},
	{IK_AXES_DIGIT, J_AXES_DIGIT, CENTER_DIGIT, K_AXES_DIGIT, JK_AXES_DIGIT,
		IJ_AXES_DIGIT, I_AXES_DIGIT},
	{IJ_AXES_DIGIT, CENTER_DIGIT, IK_AXES_DIGIT, J_AXES_DIGIT, K_AXES_DIGIT,
		I_AXES_DIGIT, JK_AXES_DIGIT},
}

// NEW_ADJUSTMENT_II: New traversal direction when traversing along class II grids.
// Current digit -> direction -> new ap7 move (at coarser level).
var NEW_ADJUSTMENT_II = [7][7]Direction{
	{CENTER_DIGIT, CENTER_DIGIT, CENTER_DIGIT, CENTER_DIGIT, CENTER_DIGIT,
		CENTER_DIGIT, CENTER_DIGIT},
	{CENTER_DIGIT, K_AXES_DIGIT, CENTER_DIGIT, K_AXES_DIGIT, CENTER_DIGIT,
		IK_AXES_DIGIT, CENTER_DIGIT},
	{CENTER_DIGIT, CENTER_DIGIT, J_AXES_DIGIT, JK_AXES_DIGIT, CENTER_DIGIT,
		CENTER_DIGIT, J_AXES_DIGIT},
	{CENTER_DIGIT, K_AXES_DIGIT, JK_AXES_DIGIT, JK_AXES_DIGIT, CENTER_DIGIT,
		CENTER_DIGIT, CENTER_DIGIT},
	{CENTER_DIGIT, CENTER_DIGIT, CENTER_DIGIT, CENTER_DIGIT, I_AXES_DIGIT,
		I_AXES_DIGIT, IJ_AXES_DIGIT},
	{CENTER_DIGIT, IK_AXES_DIGIT, CENTER_DIGIT, CENTER_DIGIT, I_AXES_DIGIT,
		IK_AXES_DIGIT, CENTER_DIGIT},
	{CENTER_DIGIT, CENTER_DIGIT, J_AXES_DIGIT, CENTER_DIGIT, IJ_AXES_DIGIT,
		CENTER_DIGIT, IJ_AXES_DIGIT},
}

// NEW_DIGIT_III: New traversal direction when traversing along class III grids.
// Current digit -> direction -> new ap7 move (at coarser level).
var NEW_DIGIT_III = [7][7]Direction{
	{CENTER_DIGIT, K_AXES_DIGIT, J_AXES_DIGIT, JK_AXES_DIGIT, I_AXES_DIGIT,
		IK_AXES_DIGIT, IJ_AXES_DIGIT},
	{K_AXES_DIGIT, J_AXES_DIGIT, JK_AXES_DIGIT, I_AXES_DIGIT, IK_AXES_DIGIT,
		IJ_AXES_DIGIT, CENTER_DIGIT},
	{J_AXES_DIGIT, JK_AXES_DIGIT, I_AXES_DIGIT, IK_AXES_DIGIT, IJ_AXES_DIGIT,
		CENTER_DIGIT, K_AXES_DIGIT},
	{JK_AXES_DIGIT, I_AXES_DIGIT, IK_AXES_DIGIT, IJ_AXES_DIGIT, CENTER_DIGIT,
		K_AXES_DIGIT, J_AXES_DIGIT},
	{I_AXES_DIGIT, IK_AXES_DIGIT, IJ_AXES_DIGIT, CENTER_DIGIT, K_AXES_DIGIT,
		J_AXES_DIGIT, JK_AXES_DIGIT},
	{IK_AXES_DIGIT, IJ_AXES_DIGIT, CENTER_DIGIT, K_AXES_DIGIT, J_AXES_DIGIT,
		JK_AXES_DIGIT, I_AXES_DIGIT},
	{IJ_AXES_DIGIT, CENTER_DIGIT, K_AXES_DIGIT, J_AXES_DIGIT, JK_AXES_DIGIT,
		I_AXES_DIGIT, IK_AXES_DIGIT},
}

// NEW_ADJUSTMENT_III: New traversal direction when traversing along class III grids.
// Current digit -> direction -> new ap7 move (at coarser level).
var NEW_ADJUSTMENT_III = [7][7]Direction{
	{CENTER_DIGIT, CENTER_DIGIT, CENTER_DIGIT, CENTER_DIGIT, CENTER_DIGIT,
		CENTER_DIGIT, CENTER_DIGIT},
	{CENTER_DIGIT, K_AXES_DIGIT, CENTER_DIGIT, JK_AXES_DIGIT, CENTER_DIGIT,
		K_AXES_DIGIT, CENTER_DIGIT},
	{CENTER_DIGIT, CENTER_DIGIT, J_AXES_DIGIT, J_AXES_DIGIT, CENTER_DIGIT,
		CENTER_DIGIT, IJ_AXES_DIGIT},
	{CENTER_DIGIT, JK_AXES_DIGIT, J_AXES_DIGIT, JK_AXES_DIGIT, CENTER_DIGIT,
		CENTER_DIGIT, CENTER_DIGIT},
	{CENTER_DIGIT, CENTER_DIGIT, CENTER_DIGIT, CENTER_DIGIT, I_AXES_DIGIT,
		IK_AXES_DIGIT, I_AXES_DIGIT},
	{CENTER_DIGIT, K_AXES_DIGIT, CENTER_DIGIT, CENTER_DIGIT, IK_AXES_DIGIT,
		IK_AXES_DIGIT, CENTER_DIGIT},
	{CENTER_DIGIT, CENTER_DIGIT, IJ_AXES_DIGIT, CENTER_DIGIT, I_AXES_DIGIT,
		CENTER_DIGIT, IJ_AXES_DIGIT},
}
