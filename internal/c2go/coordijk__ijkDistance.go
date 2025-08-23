package c2go

// abs32 returns the absolute value of x for int32.
func abs32(x int32) int32 {
	if x < 0 {
		return -x
	}
	return x
}

// max32 returns the maximum of a and b for int32.
func max32(a, b int32) int32 {
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
	absDiff := CoordIJK{I: abs32(diff.I), J: abs32(diff.J), K: abs32(diff.K)}
	return int(max32(absDiff.I, max32(absDiff.J, absDiff.K)))
}
