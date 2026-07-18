package h3

// ijkToIj transforms coordinates from the IJK coordinate system to the quadIJ coordinate system.
// Mirrors H3's coordijk.h::ijkToIj behavior.
// Ported from H3 C: coordijk.h::ijkToIj.
func ijkToIj(ijk *coordIJK, ij *CoordIJ) {
	ij.I = ijk.I - ijk.K
	ij.J = ijk.J - ijk.K
}
