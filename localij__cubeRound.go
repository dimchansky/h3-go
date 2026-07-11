package h3

import "math"

// cubeRound rounds cube coordinates to the nearest integer cube coordinates.
// Maintains the constraint that i + j + k = 0 for valid cube coordinates.
// Ported from H3 C: localij.c::cubeRound.
func cubeRound(i, j, k float64, ijk *coordIJK) {
	ri := int32(math.Round(i))
	rj := int32(math.Round(j))
	rk := int32(math.Round(k))

	iDiff := math.Abs(float64(ri) - i)
	jDiff := math.Abs(float64(rj) - j)
	kDiff := math.Abs(float64(rk) - k)

	// Round, maintaining valid cube coords
	if iDiff > jDiff && iDiff > kDiff {
		ri = -rj - rk
	} else if jDiff > kDiff {
		rj = -ri - rk
	} else {
		rk = -ri - rj
	}

	ijk.I = ri
	ijk.J = rj
	ijk.K = rk
}
