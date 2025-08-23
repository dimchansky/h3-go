package c2go

// abs returns the absolute value of x.
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// max returns the maximum of a and b.
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// ijkDistance computes the distance between two IJK coordinates.
// Mirrors H3's coordijk.c::ijkDistance behavior.
// Ported from H3 C: coordijk.c::ijkDistance
func ijkDistance(c1, c2 *CoordIJK) int {
	var diff CoordIJK
	_ijkSub(c1, c2, &diff)
	_ijkNormalize(&diff)
	absDiff := CoordIJK{I: abs(diff.I), J: abs(diff.J), K: abs(diff.K)}
	return max(absDiff.I, max(absDiff.J, absDiff.K))
}
