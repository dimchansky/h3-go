package h3

import "math"

// _upAp7 finds the normalized IJK coordinates of the indexing parent of a cell
// in a counter-clockwise aperture 7 grid. Works in place.
// Mirrors H3's coordijk.c::_upAp7 behavior.
// Ported from H3 C: coordijk.c::_upAp7.
func _upAp7(ijk *coordIJK) {
	// convert to CoordIJ
	i := ijk.I - ijk.K
	j := ijk.J - ijk.K
	ijk.I = int32(math.Round((3*float64(i) - float64(j)) * mOneseventh))
	ijk.J = int32(math.Round((float64(i) + 2*float64(j)) * mOneseventh))
	ijk.K = 0
	_ijkNormalize(ijk)
}
