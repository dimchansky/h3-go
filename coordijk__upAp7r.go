package h3

import "math"

// _upAp7r finds the normalized IJK coordinates of the indexing parent of a cell
// in a clockwise aperture 7 grid. Works in place.
// Mirrors H3's coordijk.h::_upAp7r behavior.
// Ported from H3 C: coordijk.h::_upAp7r.
func _upAp7r(ijk *coordIJK) {
	// convert to CoordIJ
	i := ijk.I - ijk.K
	j := ijk.J - ijk.K
	ijk.I = int32(math.Round((2*float64(i) + float64(j)) * mOneseventh))
	ijk.J = int32(math.Round((3*float64(j) - float64(i)) * mOneseventh))
	ijk.K = 0
	_ijkNormalize(ijk)
}
