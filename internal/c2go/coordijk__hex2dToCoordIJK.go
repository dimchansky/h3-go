package c2go

import "math"

// _hex2dToCoordIJK converts 2D hex coordinates to IJK coordinates.
// Mirrors H3's coordijk.c::_hex2dToCoordIJK behavior.
// Ported from H3 C: coordijk.c::_hex2dToCoordIJK
func _hex2dToCoordIJK(v *Vec2d, h *CoordIJK) {
	var a1, a2 float64
	var x1, x2 float64
	var m1, m2 int32
	var r1, r2 float64

	// quantize into the ij system and then normalize
	h.K = 0

	a1 = math.Abs(v.X)
	a2 = math.Abs(v.Y)

	// first do a reverse conversion
	x2 = a2 * M_RSIN60
	x1 = a1 + x2/2.0

	// check if we have the center of a hex
	m1 = int32(x1)
	m2 = int32(x2)

	// otherwise round correctly
	r1 = x1 - float64(m1)
	r2 = x2 - float64(m2)

	if r1 < 0.5 {
		if r1 < 1.0/3.0 {
			if r2 < (1.0+r1)/2.0 {
				h.I = int32(m1)
				h.J = int32(m2)
			} else {
				h.I = int32(m1)
				h.J = int32(m2 + 1)
			}
		} else {
			if r2 < (1.0 - r1) {
				h.J = int32(m2)
			} else {
				h.J = int32(m2 + 1)
			}

			if (1.0-r1) <= r2 && r2 < (2.0*r1) {
				h.I = int32(m1 + 1)
			} else {
				h.I = int32(m1)
			}
		}
	} else {
		if r1 < 2.0/3.0 {
			if r2 < (1.0 - r1) {
				h.J = int32(m2)
			} else {
				h.J = int32(m2 + 1)
			}

			if (2.0*r1-1.0) < r2 && r2 < (1.0-r1) {
				h.I = int32(m1)
			} else {
				h.I = int32(m1 + 1)
			}
		} else {
			if r2 < (r1 / 2.0) {
				h.I = int32(m1 + 1)
				h.J = int32(m2)
			} else {
				h.I = int32(m1 + 1)
				h.J = int32(m2 + 1)
			}
		}
	}

	// now fold across the axes if necessary

	if v.X < 0.0 {
		if (h.J % 2) == 0 { // even
			axisi := int64(h.J / 2)
			diff := int64(h.I) - axisi
			h.I = int32(int64(h.I) - 2*diff)
		} else {
			axisi := int64((h.J + 1) / 2)
			diff := int64(h.I) - axisi
			h.I = int32(int64(h.I) - (2*diff + 1))
		}
	}

	if v.Y < 0.0 {
		h.I = h.I - (2*h.J+1)/2
		h.J = -1 * h.J
	}

	_ijkNormalize(h)
}
