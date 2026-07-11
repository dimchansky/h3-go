package h3

import "math"

// _upAp7rChecked finds the normalized ijk coordinates of the indexing parent of a cell
// in a clockwise aperture 7 grid with overflow checking. Works in place.
//
// Returns eFailed if signed integer overflow could occur.
// Assumes ijk is IJK+ coordinates (no negative numbers).
//
// Mirrors H3's coordijk.c::_upAp7rChecked behavior.
// Ported from H3 C: coordijk.c::_upAp7rChecked.
func _upAp7rChecked(ijk *coordIJK) h3Error {
	// Doesn't need to be checked because i, j, and k must all be non-negative
	i := ijk.I - ijk.K
	j := ijk.J - ijk.K

	// <0 is checked because the input must all be non-negative, but some
	// negative inputs are used in unit tests to exercise the below.
	if i >= int32Max3 || j >= int32Max3 || i < 0 || j < 0 {
		if addInt32sOverflows(i, i) {
			return eFailed
		}
		i2 := i + i
		if addInt32sOverflows(j, j) {
			return eFailed
		}
		j2 := j + j
		if addInt32sOverflows(j2, j) {
			return eFailed
		}
		j3 := j2 + j

		if addInt32sOverflows(i2, j) {
			return eFailed
		}
		if subInt32sOverflows(j3, i) {
			return eFailed
		}
	}

	ijk.I = int32(math.Round(float64(i*2+j) * mOneseventh))
	ijk.J = int32(math.Round(float64(j*3-i) * mOneseventh))
	ijk.K = 0

	// Expected not to be reachable, because max + min or max - min would need
	// to overflow.
	if _ijkNormalizeCouldOverflow(ijk) {
		return eFailed
	}
	_ijkNormalize(ijk)
	return eSuccess
}
