package h3

// gridPathCells returns a line of H3 indexes between two H3 indexes (inclusive).
//
// This function relies on gridDistance(start, end) to determine the
// expected path length, and returns the same error if gridDistance fails.
//
// Path construction is performed by straight-line interpolation in the
// origin-anchored local IJK coordinate space: first anchored at start;
// if that fails, retried anchored at end with the sequence reversed into
// out. If both attempts fail, the first attempt's error is returned.
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
func gridPathCells(out []h3Index, start h3Index, end h3Index) h3Error {
	var distance int64
	if err := gridDistance(start, end, &distance); err != eSuccess {
		// Early exit if we can't calculate the line
		return err
	}

	required := distance + 1
	if int64(len(out)) < required {
		// Mirror C contract (caller allocates) but keep Go safe.
		return eFailed
	}

	if distance == 0 {
		out[0] = start
		return eSuccess
	}

	// Straight-line interpolation in local IJK space anchored at `start`.
	interpolateErr := gridPathCellsInterpolate(start, end, distance, out, 0, 1)
	if interpolateErr == eSuccess {
		return eSuccess
	}

	// Retry interpolation anchored at `end` and reverse the output.
	// This can resolve cases where the local IJK chart is discontinuous
	// relative to one origin but not the other.
	reverseErr := gridPathCellsInterpolate(end, start, distance, out, distance, -1)
	if reverseErr == eSuccess {
		return eSuccess
	}

	return interpolateErr
}
