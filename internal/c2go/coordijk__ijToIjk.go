package c2go

// ijToIjk transforms coordinates from the IJ coordinate system to the IJK+ coordinate system.
// Returns E_SUCCESS on success, E_FAILED if signed integer overflow would have occurred.
// Mirrors H3's coordijk.c::ijToIjk behavior.
// Ported from H3 C: coordijk.c::ijToIjk
func ijToIjk(ij *CoordIJ, ijk *CoordIJK) H3Error {
	ijk.I = ij.I
	ijk.J = ij.J
	ijk.K = 0

	if _ijkNormalizeCouldOverflow(ijk) {
		return E_FAILED
	}

	_ijkNormalize(ijk)
	return E_SUCCESS
}
