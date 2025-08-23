package c2go

import "math"

// _upAp7Checked finds the normalized ijk coordinates of the indexing parent of a cell
// in a counter-clockwise aperture 7 grid with overflow checking. Works in place.
//
// Returns E_FAILED if signed integer overflow could occur.
// Assumes ijk is IJK+ coordinates (no negative numbers).
//
// Mirrors H3's coordijk.c::_upAp7Checked behavior.
// Ported from H3 C: coordijk.c::upAp7Checked
func _upAp7Checked(ijk *CoordIJK) H3Error {
	// Doesn't need to be checked because i, j, and k must all be non-negative
	i := int32(ijk.I - ijk.K)
	j := int32(ijk.J - ijk.K)

	// <0 is checked because the input must all be non-negative, but some
	// negative inputs are used in unit tests to exercise the below.
	if i >= INT32_MAX_3 || j >= INT32_MAX_3 || i < 0 || j < 0 {
		if addInt32sOverflows(i, i) {
			return E_FAILED
		}
		i2 := i + i
		if addInt32sOverflows(i2, i) {
			return E_FAILED
		}
		i3 := i2 + i
		if addInt32sOverflows(j, j) {
			return E_FAILED
		}
		j2 := j + j

		if subInt32sOverflows(i3, j) {
			return E_FAILED
		}
		if addInt32sOverflows(i, j2) {
			return E_FAILED
		}
	}

	ijk.I = int32(math.Round(float64(i*3-j) * M_ONESEVENTH))
	ijk.J = int32(math.Round(float64(i+j*2) * M_ONESEVENTH))
	ijk.K = 0

	// Expected not to be reachable, because max + min or max - min would need
	// to overflow.
	if _ijkNormalizeCouldOverflow(ijk) {
		return E_FAILED
	}
	_ijkNormalize(ijk)
	return E_SUCCESS
}
