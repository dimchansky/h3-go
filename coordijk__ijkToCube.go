package h3

// ijkToCube converts IJK coordinates to cube coordinates.
// Mirrors H3's coordijk.h::ijkToCube behavior.
// Ported from H3 C: coordijk.h::ijkToCube.
func ijkToCube(ijk *coordIJK) {
	ijk.I = -ijk.I + ijk.K
	ijk.J = ijk.J - ijk.K
	ijk.K = -ijk.I - ijk.J
}
