package c2go

// cubeToIjk converts cube coordinates to IJK coordinates.
// Mirrors H3's coordijk.c::cubeToIjk behavior.
func cubeToIjk(ijk *CoordIJK) {
	ijk.I = -ijk.I
	ijk.K = 0
	_ijkNormalize(ijk)
}
