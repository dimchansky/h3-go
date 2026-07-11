package h3

// ijToIjk transforms coordinates from the quadIJ coordinate system to the IJK+ coordinate system.
// Returns eSuccess on success, eFailed if signed integer overflow would have occurred.
// Mirrors H3's coordijk.c::ijToIjk behavior.
// Ported from H3 C: coordijk.c::ijToIjk.
func ijToIjk(ij *CoordIJ, ijk *coordIJK) h3Error {
	ijk.I = ij.I
	ijk.J = ij.J
	ijk.K = 0

	if _ijkNormalizeCouldOverflow(ijk) {
		return eFailed
	}

	_ijkNormalize(ijk)
	return eSuccess
}
