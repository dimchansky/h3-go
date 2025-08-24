package c2go

// areNeighborCells returns whether or not the provided H3Indexes are neighbors.
// This function determines if two H3 cells are adjacent to each other.
// Uses an optimized approach for cells that share the same parent, falling back
// to gridDisk check for cells that don't share a parent.
// Ported from H3 C: directedEdge.c::areNeighborCells
func areNeighborCells(origin, destination H3Index) (bool, H3Error) {
	// Make sure they're hexagon indexes
	if getMode(origin) != H3_CELL_MODE ||
		getMode(destination) != H3_CELL_MODE {
		return false, E_CELL_INVALID
	}

	// Hexagons cannot be neighbors with themselves
	if origin == destination {
		return false, E_SUCCESS
	}

	// Only hexagons in the same resolution can be neighbors
	if getResolution(origin) != getResolution(destination) {
		return false, E_RES_MISMATCH
	}

	// H3 Indexes that share the same parent are very likely to be neighbors
	// Child 0 is neighbor with all of its parent's 'offspring', the other
	// children are neighbors with 3 of the 7 children. So a simple comparison
	// of origin and destination parents and then a lookup table of the children
	// is a super-cheap way to possibly determine they are neighbors.
	parentRes := getResolution(origin) - 1
	if parentRes > 0 {
		originParent, err1 := cellToParent(origin, int32(parentRes))
		if err1 != E_SUCCESS {
			return false, err1
		}
		destinationParent, err2 := cellToParent(destination, int32(parentRes))
		if err2 != E_SUCCESS {
			return false, err2
		}
		if originParent == destinationParent {
			originResDigit := Direction(getIndexDigit(origin, parentRes+1))
			destinationResDigit := Direction(getIndexDigit(destination, parentRes+1))
			if originResDigit == CENTER_DIGIT ||
				destinationResDigit == CENTER_DIGIT {
				return true, E_SUCCESS
			}
			if originResDigit >= INVALID_DIGIT {
				// Prevent indexing off the end of the array below
				return false, E_CELL_INVALID
			}
			isPent := isPentagon(originParent)
			if (originResDigit == K_AXES_DIGIT ||
				destinationResDigit == K_AXES_DIGIT) &&
				isPent {
				// If these are invalid cells, fail rather than incorrectly
				// reporting neighbors. For pentagon cells that are actually
				// neighbors across the deleted subsequence, they will fail the
				// optimized check below, but they will be accepted by the
				// gridDisk check below that.
				return false, E_CELL_INVALID
			}
			// These sets are the relevant neighbors in the clockwise
			// and counter-clockwise
			neighborSetClockwise := []Direction{
				CENTER_DIGIT, JK_AXES_DIGIT, IJ_AXES_DIGIT, J_AXES_DIGIT,
				IK_AXES_DIGIT, K_AXES_DIGIT, I_AXES_DIGIT}
			neighborSetCounterclockwise := []Direction{
				CENTER_DIGIT, IK_AXES_DIGIT, JK_AXES_DIGIT, K_AXES_DIGIT,
				IJ_AXES_DIGIT, I_AXES_DIGIT, J_AXES_DIGIT}
			if neighborSetClockwise[originResDigit] == destinationResDigit ||
				neighborSetCounterclockwise[originResDigit] == destinationResDigit {
				return true, E_SUCCESS
			}
		}
	}

	// Otherwise, we have to determine the neighbor relationship the "hard" way.
	neighborRing := make([]H3Index, 7)
	err := gridDisk(origin, 1, neighborRing)
	if err != E_SUCCESS {
		// If gridDisk fails, assume they are not neighbors (C behavior)
		// rather than propagating the error
		return false, E_SUCCESS
	}
	for i := 0; i < 7; i++ {
		if neighborRing[i] == destination {
			return true, E_SUCCESS
		}
	}

	// Made it here, they definitely aren't neighbors
	return false, E_SUCCESS
}
