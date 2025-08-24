package c2go

// ijkToCube converts IJK coordinates to cube coordinates.
// Mirrors H3's coordijk.c::ijkToCube behavior.
// Ported from H3 C: coordijk.c::ijkToCube
func ijkToCube(ijk *CoordIJK) {
	ijk.I = -ijk.I + ijk.K
	ijk.J = ijk.J - ijk.K
	ijk.K = -ijk.I - ijk.J
}
