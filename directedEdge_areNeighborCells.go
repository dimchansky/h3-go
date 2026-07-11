package h3

// areNeighborCells returns whether or not the provided H3Indexes are neighbors.
// This function determines if two H3 cells are adjacent to each other.
// Uses an optimized approach for cells that share the same parent, falling back
// to gridDisk check for cells that don't share a parent.
// Ported from H3 C: directedEdge.c::areNeighborCells.
func areNeighborCells(origin, destination h3Index) (bool, h3Error) {
	// Make sure they're hexagon indexes
	if getMode(origin) != h3CellMode ||
		getMode(destination) != h3CellMode {
		return false, eCellInvalid
	}

	// Hexagons cannot be neighbors with themselves
	if origin == destination {
		return false, eSuccess
	}

	// Only hexagons in the same resolution can be neighbors
	if getResolution(origin) != getResolution(destination) {
		return false, eResMismatch
	}

	// H3 Indexes that share the same parent are very likely to be neighbors
	// Child 0 is neighbor with all of its parent's 'offspring', the other
	// children are neighbors with 3 of the 7 children. So a simple comparison
	// of origin and destination parents and then a lookup table of the children
	// is a super-cheap way to possibly determine they are neighbors.
	parentRes := getResolution(origin) - 1
	if parentRes > 0 {
		originParent, err1 := cellToParent(origin, int32(parentRes))
		if err1 != eSuccess {
			return false, err1
		}
		destinationParent, err2 := cellToParent(destination, int32(parentRes))
		if err2 != eSuccess {
			return false, err2
		}
		if originParent == destinationParent {
			originResDigit := direction(h3GetIndexDigit(origin, parentRes+1))
			destinationResDigit := direction(h3GetIndexDigit(destination, parentRes+1))
			if originResDigit == centerDigit ||
				destinationResDigit == centerDigit {
				return true, eSuccess
			}
			if originResDigit >= invalidDigit {
				// Prevent indexing off the end of the array below
				return false, eCellInvalid
			}
			isPent := isPentagon(originParent)
			if (originResDigit == kAxesDigit ||
				destinationResDigit == kAxesDigit) &&
				isPent {
				// If these are invalid cells, fail rather than incorrectly
				// reporting neighbors. For pentagon cells that are actually
				// neighbors across the deleted subsequence, they will fail the
				// optimized check below, but they will be accepted by the
				// gridDisk check below that.
				return false, eCellInvalid
			}
			// These sets are the relevant neighbors in the clockwise
			// and counter-clockwise
			neighborSetClockwise := []direction{
				centerDigit, jkAxesDigit, ijAxesDigit, jAxesDigit,
				ikAxesDigit, kAxesDigit, iAxesDigit}
			neighborSetCounterclockwise := []direction{
				centerDigit, ikAxesDigit, jkAxesDigit, kAxesDigit,
				ijAxesDigit, iAxesDigit, jAxesDigit}
			if neighborSetClockwise[originResDigit] == destinationResDigit ||
				neighborSetCounterclockwise[originResDigit] == destinationResDigit {
				return true, eSuccess
			}
		}
	}

	// Otherwise, we have to determine the neighbor relationship the "hard" way.
	neighborRing := make([]h3Index, 7)
	err := gridDisk(origin, 1, neighborRing)
	if err != eSuccess {
		// If gridDisk fails, assume they are not neighbors (C behavior)
		// rather than propagating the error
		return false, eSuccess
	}
	for i := 0; i < 7; i++ {
		if neighborRing[i] == destination {
			return true, eSuccess
		}
	}

	// Made it here, they definitely aren't neighbors
	return false, eSuccess
}
