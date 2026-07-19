package h3

// gridPathCellsInterpolate attempts to generate a shortest-length path by
// interpolating through an origin-anchored local IJK coordinate space.
//
// It can fail if interpolation lands on intermediate IJK coordinates that
// cannot be mapped back to valid H3 cells: the origin-anchored local
// IJ(K) space is not globally continuous, and some intermediate
// coordinates have no inverse mapping back to a cell in the chosen chart
// (for example, due to discontinuities or warping near pentagons).
//
// The output is written to out[outOffset + outStep*n], allowing callers
// to fill the path in either direction without an intermediate buffer.
// Ported from H3 C: localij.c::gridPathCellsInterpolate.
func gridPathCellsInterpolate(start, end h3Index, distance int64, out []h3Index, outOffset, outStep int64) h3Error {
	// Get IJK coords for the start and end. We've already confirmed
	// that these can be calculated with the distance check above.
	startIjk := coordIJK{}
	endIjk := coordIJK{}

	// Convert H3 addresses to IJK coords
	startError := cellToLocalIjk(start, start, &startIjk)
	// NEVER(startError) in C - unreachable because this was called as part
	// of gridDistance
	if startError != eSuccess {
		return startError
	}
	endError := cellToLocalIjk(start, end, &endIjk)
	// NEVER(endError) in C - unreachable because this was called as part
	// of gridDistance
	if endError != eSuccess {
		return endError
	}

	// Convert IJK to cube coordinates suitable for linear interpolation
	ijkToCube(&startIjk)
	ijkToCube(&endIjk)

	invDistance := 1.0 / float64(distance)
	iStep := float64(endIjk.I-startIjk.I) * invDistance
	jStep := float64(endIjk.J-startIjk.J) * invDistance
	kStep := float64(endIjk.K-startIjk.K) * invDistance

	currentIjk := coordIJK{I: startIjk.I, J: startIjk.J, K: startIjk.K}
	for n := int64(0); n <= distance; n++ {
		cubeRound(
			float64(startIjk.I)+iStep*float64(n),
			float64(startIjk.J)+jStep*float64(n),
			float64(startIjk.K)+kStep*float64(n),
			&currentIjk,
		)
		// Convert cube -> ijk -> h3 index
		cubeToIjk(&currentIjk)
		idx := outOffset + outStep*n
		currentError := localIjkToCell(start, &currentIjk, &out[idx])
		if currentError != eSuccess {
			// The cells between `start` and `end` may fall in pentagon
			// distortion.
			return currentError
		}
	}

	return eSuccess
}
