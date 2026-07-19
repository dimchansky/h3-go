package h3

// cubeToIjk converts cube coordinates to IJK coordinates.
// Mirrors H3's coordijk.h::cubeToIjk behavior.
// Ported from H3 C: coordijk.h::cubeToIjk.
func cubeToIjk(ijk *coordIJK) {
	ijk.I = -ijk.I
	ijk.K = 0
	_ijkNormalize(ijk)
}
