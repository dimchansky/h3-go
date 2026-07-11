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
// Ported from H3 C: algos.c::h3NeighborRotations.
func h3NeighborRotations(origin h3Index, dir direction, rotations *int32, out *h3Index) h3Error {
	current := origin

	if dir < centerDigit || dir >= invalidDigit {
		return eFailed
	}

	// Ensure that rotations is modulo'd by 6 before any possible addition,
	// to protect against signed integer overflow.
	*rotations = *rotations % 6
	for i := int32(0); i < *rotations; i++ {
		dir = _rotate60ccw(dir)
	}

	newRotations := int32(0)
	oldBaseCell := int32(getBaseCell(current))
	if oldBaseCell < 0 || oldBaseCell >= numBaseCells {
		// Base cells less than zero can not be represented in an index
		return eCellInvalid
	}
	oldLeadingDigit := _h3LeadingNonZeroDigit(current)

	// Adjust the indexing digits and, if needed, the base cell.
	r := getResolution(current) - 1
	for {
		if r == -1 {
			current = setBaseCell(current, baseCellNeighbors[oldBaseCell][dir])
			newRotations = int32(baseCellNeighbor60CCWRots[oldBaseCell][dir])

			if getBaseCell(current) == invalidBaseCell {
				// Adjust for the deleted k vertex at the base cell level.
				// This edge actually borders a different neighbor.
				current = setBaseCell(current, baseCellNeighbors[oldBaseCell][ikAxesDigit])
				newRotations = int32(baseCellNeighbor60CCWRots[oldBaseCell][ikAxesDigit])

				// perform the adjustment for the k-subsequence we're skipping
				// over.
				current = _h3Rotate60ccw(current)
				*rotations = *rotations + 1
			}

			break
		} else {
			oldDigit := direction(getIndexDigit(current, r+1))
			var nextDir direction
			if oldDigit == invalidDigit {
				// Only possible on invalid input
				return eCellInvalid
			} else if isResolutionClassIII(r + 1) {
				current = setIndexDigit(current, r+1, int32(newDigitII[oldDigit][dir]))
				nextDir = newAdjustmentII[oldDigit][dir]
			} else {
				current = setIndexDigit(current, r+1, int32(newDigitIII[oldDigit][dir]))
				nextDir = newAdjustmentIII[oldDigit][dir]
			}

			if nextDir != centerDigit {
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
		if direction(_h3LeadingNonZeroDigit(current)) == kAxesDigit {
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
				if direction(oldLeadingDigit) == centerDigit {
					// Undefined: the k direction is deleted from here
					return ePentagon
				} else if direction(oldLeadingDigit) == jkAxesDigit {
					// Rotate out of the deleted k subsequence
					// We also need an additional change to the direction we're
					// moving in
					current = _h3Rotate60ccw(current)
					*rotations = *rotations + 1
				} else if direction(oldLeadingDigit) == ikAxesDigit {
					// Rotate out of the deleted k subsequence
					// We also need an additional change to the direction we're
					// moving in
					current = _h3Rotate60cw(current)
					*rotations = *rotations + 5
				} else {
					// TODO: Should never occur, but is reachable by fuzzer
					return eFailed
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
					direction(_h3LeadingNonZeroDigit(current)) != jkAxesDigit {
					*rotations = *rotations + 1
				}
			} else if direction(_h3LeadingNonZeroDigit(current)) == ikAxesDigit &&
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

	return eSuccess
}

// Lookup tables for neighbor traversal algorithm

// newDigitII: New digit when traversing along class II grids.
// Current digit -> direction -> new digit.
var newDigitII = [7][7]direction{
	{centerDigit, kAxesDigit, jAxesDigit, jkAxesDigit, iAxesDigit,
		ikAxesDigit, ijAxesDigit},
	{kAxesDigit, iAxesDigit, jkAxesDigit, ijAxesDigit, ikAxesDigit,
		jAxesDigit, centerDigit},
	{jAxesDigit, jkAxesDigit, kAxesDigit, iAxesDigit, ijAxesDigit,
		centerDigit, ikAxesDigit},
	{jkAxesDigit, ijAxesDigit, iAxesDigit, ikAxesDigit, centerDigit,
		kAxesDigit, jAxesDigit},
	{iAxesDigit, ikAxesDigit, ijAxesDigit, centerDigit, jAxesDigit,
		jkAxesDigit, kAxesDigit},
	{ikAxesDigit, jAxesDigit, centerDigit, kAxesDigit, jkAxesDigit,
		ijAxesDigit, iAxesDigit},
	{ijAxesDigit, centerDigit, ikAxesDigit, jAxesDigit, kAxesDigit,
		iAxesDigit, jkAxesDigit},
}

// newAdjustmentII: New traversal direction when traversing along class II grids.
// Current digit -> direction -> new ap7 move (at coarser level).
var newAdjustmentII = [7][7]direction{
	{centerDigit, centerDigit, centerDigit, centerDigit, centerDigit,
		centerDigit, centerDigit},
	{centerDigit, kAxesDigit, centerDigit, kAxesDigit, centerDigit,
		ikAxesDigit, centerDigit},
	{centerDigit, centerDigit, jAxesDigit, jkAxesDigit, centerDigit,
		centerDigit, jAxesDigit},
	{centerDigit, kAxesDigit, jkAxesDigit, jkAxesDigit, centerDigit,
		centerDigit, centerDigit},
	{centerDigit, centerDigit, centerDigit, centerDigit, iAxesDigit,
		iAxesDigit, ijAxesDigit},
	{centerDigit, ikAxesDigit, centerDigit, centerDigit, iAxesDigit,
		ikAxesDigit, centerDigit},
	{centerDigit, centerDigit, jAxesDigit, centerDigit, ijAxesDigit,
		centerDigit, ijAxesDigit},
}

// newDigitIII: New traversal direction when traversing along class III grids.
// Current digit -> direction -> new ap7 move (at coarser level).
var newDigitIII = [7][7]direction{
	{centerDigit, kAxesDigit, jAxesDigit, jkAxesDigit, iAxesDigit,
		ikAxesDigit, ijAxesDigit},
	{kAxesDigit, jAxesDigit, jkAxesDigit, iAxesDigit, ikAxesDigit,
		ijAxesDigit, centerDigit},
	{jAxesDigit, jkAxesDigit, iAxesDigit, ikAxesDigit, ijAxesDigit,
		centerDigit, kAxesDigit},
	{jkAxesDigit, iAxesDigit, ikAxesDigit, ijAxesDigit, centerDigit,
		kAxesDigit, jAxesDigit},
	{iAxesDigit, ikAxesDigit, ijAxesDigit, centerDigit, kAxesDigit,
		jAxesDigit, jkAxesDigit},
	{ikAxesDigit, ijAxesDigit, centerDigit, kAxesDigit, jAxesDigit,
		jkAxesDigit, iAxesDigit},
	{ijAxesDigit, centerDigit, kAxesDigit, jAxesDigit, jkAxesDigit,
		iAxesDigit, ikAxesDigit},
}

// newAdjustmentIII: New traversal direction when traversing along class III grids.
// Current digit -> direction -> new ap7 move (at coarser level).
var newAdjustmentIII = [7][7]direction{
	{centerDigit, centerDigit, centerDigit, centerDigit, centerDigit,
		centerDigit, centerDigit},
	{centerDigit, kAxesDigit, centerDigit, jkAxesDigit, centerDigit,
		kAxesDigit, centerDigit},
	{centerDigit, centerDigit, jAxesDigit, jAxesDigit, centerDigit,
		centerDigit, ijAxesDigit},
	{centerDigit, jkAxesDigit, jAxesDigit, jkAxesDigit, centerDigit,
		centerDigit, centerDigit},
	{centerDigit, centerDigit, centerDigit, centerDigit, iAxesDigit,
		ikAxesDigit, iAxesDigit},
	{centerDigit, kAxesDigit, centerDigit, centerDigit, ikAxesDigit,
		ikAxesDigit, centerDigit},
	{centerDigit, centerDigit, ijAxesDigit, centerDigit, iAxesDigit,
		centerDigit, ijAxesDigit},
}
