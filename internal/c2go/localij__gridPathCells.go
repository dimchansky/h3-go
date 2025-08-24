package c2go

// _gridPathCells returns a line of H3 indexes between two H3 indexes (inclusive).
//
// This function may fail to find the line between two indexes, for example if they are
// very far apart. It may also fail when finding distances for indexes on opposite sides
// of a pentagon.
//
// Notes:
//   - The specific output of this function should not be considered stable across library
//     versions. The only guarantees the library provides are that the line length will be
//     `gridDistance(start, end) + 1` and that every index in the line will be a neighbor
//     of the preceding index.
//   - Lines are drawn in grid space, and may not correspond exactly to either Cartesian
//     lines or great arcs.
//
// The dst slice will be used as the output buffer if it has sufficient capacity,
// otherwise a new slice will be allocated.
//
// Ported from H3 C: localij.c::gridPathCells
func gridPathCells(dst []H3Index, start H3Index, end H3Index) ([]H3Index, H3Error) {
	var distance int64
	err := gridDistance(start, end, &distance)
	// Early exit if we can't calculate the line
	if err != E_SUCCESS {
		return nil, err
	}

	// Calculate required size
	size := distance + 1

	// Ensure dst has sufficient capacity
	if int64(cap(dst)) < size {
		dst = make([]H3Index, size)
	} else {
		dst = dst[:size]
	}

	// Get IJK coords for the start and end. We've already confirmed
	// that these can be calculated with the distance check above.
	var startIjk CoordIJK
	var endIjk CoordIJK

	// Convert H3 addresses to IJK coords
	startError := cellToLocalIjk(start, start, &startIjk)
	if startError != E_SUCCESS {
		// Unreachable because this was called as part of gridDistance
		return nil, startError
	}
	endError := cellToLocalIjk(start, end, &endIjk)
	if endError != E_SUCCESS {
		// Unreachable because this was called as part of gridDistance
		return nil, endError
	}

	// Convert IJK to cube coordinates suitable for linear interpolation
	ijkToCube(&startIjk)
	ijkToCube(&endIjk)

	var invDistance float64
	if distance != 0 {
		invDistance = 1.0 / float64(distance)
	} else {
		invDistance = 0
	}

	iStep := float64(endIjk.I-startIjk.I) * invDistance
	jStep := float64(endIjk.J-startIjk.J) * invDistance
	kStep := float64(endIjk.K-startIjk.K) * invDistance

	var currentIjk CoordIJK
	for n := int64(0); n <= distance; n++ {
		cubeRound(float64(startIjk.I)+iStep*float64(n),
			float64(startIjk.J)+jStep*float64(n),
			float64(startIjk.K)+kStep*float64(n), &currentIjk)
		// Convert cube -> ijk -> h3 index
		cubeToIjk(&currentIjk)
		currentError := localIjkToCell(start, &currentIjk, &dst[n])
		if currentError != E_SUCCESS {
			// The cells between `start` and `end` may fall in pentagon distortion.
			return nil, currentError
		}
	}

	return dst, E_SUCCESS
}
