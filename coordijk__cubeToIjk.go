package h3

// cubeToIjk converts cube coordinates to IJK coordinates.
// Mirrors H3's coordijk.c::cubeToIjk behavior.
// Ported from H3 C: coordijk.c::cubeToIjk.
func cubeToIjk(ijk *coordIJK) {
	ijk.I = -ijk.I
	ijk.K = 0
	_ijkNormalize(ijk)
}
