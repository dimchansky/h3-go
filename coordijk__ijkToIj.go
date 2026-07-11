package h3

// ijkToIj transforms coordinates from the IJK coordinate system to the IJ coordinate system.
// Mirrors H3's coordijk.c::ijkToIj behavior.
// Ported from H3 C: coordijk.c::ijkToIj.
func ijkToIj(ijk *CoordIJK, ij *CoordIJ) {
	ij.I = ijk.I - ijk.K
	ij.J = ijk.J - ijk.K
}
