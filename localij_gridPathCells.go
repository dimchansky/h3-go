package h3

// gridPathCells returns a line of H3 indexes between two H3 indexes (inclusive).
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
// Ported from H3 C: localij.c::gridPathCells.
func gridPathCells(out []H3Index, start H3Index, end H3Index) H3Error {
	var distance int64
	if err := gridDistance(start, end, &distance); err != E_SUCCESS {
		// Early exit if we can't calculate the line
		return err
	}

	required := distance + 1
	if int64(len(out)) < required {
		// Mirror C contract (caller allocates) but keep Go safe.
		return E_FAILED
	}

	// Get IJK coords for the start and end.
	var startIjk, endIjk CoordIJK
	if err := cellToLocalIjk(start, start, &startIjk); err != E_SUCCESS {
		// Unreachable in C path (was already validated by gridDistance)
		return err
	}
	if err := cellToLocalIjk(start, end, &endIjk); err != E_SUCCESS {
		// Unreachable in C path (was already validated by gridDistance)
		return err
	}

	// Convert to cube coordinates for linear interpolation.
	ijkToCube(&startIjk)
	ijkToCube(&endIjk)

	invDistance := 0.0
	if distance != 0 {
		invDistance = 1.0 / float64(distance)
	}

	iStep := float64(endIjk.I-startIjk.I) * invDistance
	jStep := float64(endIjk.J-startIjk.J) * invDistance
	kStep := float64(endIjk.K-startIjk.K) * invDistance

	currentIjk := CoordIJK{I: startIjk.I, J: startIjk.J, K: startIjk.K}
	for n := int64(0); n <= distance; n++ {
		cubeRound(
			float64(startIjk.I)+iStep*float64(n),
			float64(startIjk.J)+jStep*float64(n),
			float64(startIjk.K)+kStep*float64(n),
			&currentIjk,
		)
		// Convert cube -> ijk -> H3 index
		cubeToIjk(&currentIjk)
		if err := localIjkToCell(start, &currentIjk, &out[n]); err != E_SUCCESS {
			// Cells between start and end may cross pentagon distortion.
			return err
		}
	}

	return E_SUCCESS
}
