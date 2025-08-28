package h3

import "math"

// _upAp7r finds the normalized IJK coordinates of the indexing parent of a cell
// in a clockwise aperture 7 grid. Works in place.
// Mirrors H3's coordijk.c::_upAp7r behavior.
// Ported from H3 C: coordijk.c::_upAp7r
func _upAp7r(ijk *CoordIJK) {
	// convert to CoordIJ
	i := ijk.I - ijk.K
	j := ijk.J - ijk.K
	ijk.I = int32(math.Round((2*float64(i) + float64(j)) * M_ONESEVENTH))
	ijk.J = int32(math.Round((3*float64(j) - float64(i)) * M_ONESEVENTH))
	ijk.K = 0
	_ijkNormalize(ijk)
}
