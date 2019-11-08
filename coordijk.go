package h3

import (
	"math"
)

//lint:file-ignore U1000 Ignore all unused code

// CoordIJK holds IJK hexagon coordinates.
// Each axis is spaced 120 degrees apart.
type CoordIJK struct {
	i int // i component
	j int // j component
	k int // k component
}

var (
	// UNIT_VECS are CoordIJK unit vectors corresponding to the 7 H3 digits.
	UNIT_VECS = []CoordIJK{
		{0, 0, 0}, // direction 0
		{0, 0, 1}, // direction 1
		{0, 1, 0}, // direction 2
		{0, 1, 1}, // direction 3
		{1, 0, 0}, // direction 4
		{1, 0, 1}, // direction 5
		{1, 1, 0}, // direction 6
	}
)

// Direction is H3 digit representing ijk+ axes direction.
// Values will be within the lowest 3 bits of an integer.
type Direction int

const (
	CENTER_DIGIT  Direction = 0                           // H3 digit in center
	K_AXES_DIGIT  Direction = 1                           // H3 digit in k-axes direction
	J_AXES_DIGIT  Direction = 2                           // H3 digit in j-axes direction
	JK_AXES_DIGIT           = J_AXES_DIGIT | K_AXES_DIGIT // 3 - H3 digit in j == k direction
	I_AXES_DIGIT  Direction = 4                           // H3 digit in i-axes direction
	IK_AXES_DIGIT           = I_AXES_DIGIT | K_AXES_DIGIT // 5 - H3 digit in i == k direction
	IJ_AXES_DIGIT           = I_AXES_DIGIT | J_AXES_DIGIT // 6 - H3 digit in i == j direction
	INVALID_DIGIT Direction = 7                           // H3 digit in the invalid direction
	NUM_DIGITS              = INVALID_DIGIT               // Valid digits will be less than this value. Same value as INVALID_DIGIT.
)

// _setIJK sets an IJK coordinate to the specified component values.
// `ijk`: The IJK coordinate to set.
// `i`: The desired i component value.
// `j`: The desired j component value.
// `k`: The desired k component value.
func _setIJK(ijk *CoordIJK, i int, j int, k int) {
	ijk.i = i
	ijk.j = j
	ijk.k = k
}

// _hex2dToCoordIJK determines the containing hex in ijk+ coordinates for a 2D cartesian
// coordinate vector (from DGGRID).
// `v`: The 2D cartesian coordinate vector.
// `h`: The ijk+ coordinates of the containing hex.
func _hex2dToCoordIJK(v *Vec2d, h *CoordIJK) {
	// quantize into the ij system and then normalize
	h.k = 0

	a1 := math.Abs(v.x)
	a2 := math.Abs(v.y)

	// first do a reverse conversion
	x2 := a2 / M_SIN60
	x1 := a1 + x2/2.0

	// check if we have the center of a hex
	m1 := int(x1)
	m2 := int(x2)

	// otherwise round correctly
	r1 := x1 - float64(m1)
	r2 := x2 - float64(m2)

	if r1 < 0.5 {
		if r1 < 1.0/3.0 {
			if r2 < (1.0+r1)/2.0 {
				h.i = m1
				h.j = m2
			} else {
				h.i = m1
				h.j = m2 + 1
			}
		} else {
			if r2 < (1.0 - r1) {
				h.j = m2
			} else {
				h.j = m2 + 1
			}

			if (1.0-r1) <= r2 && r2 < (2.0*r1) {
				h.i = m1 + 1
			} else {
				h.i = m1
			}
		}
	} else {
		if r1 < 2.0/3.0 {
			if r2 < (1.0 - r1) {
				h.j = m2
			} else {
				h.j = m2 + 1
			}

			if (2.0*r1-1.0) < r2 && r2 < (1.0-r1) {
				h.i = m1
			} else {
				h.i = m1 + 1
			}
		} else {
			if r2 < (r1 / 2.0) {
				h.i = m1 + 1
				h.j = m2
			} else {
				h.i = m1 + 1
				h.j = m2 + 1
			}
		}
	}

	// now fold across the axes if necessary

	if v.x < 0.0 {
		if (h.j % 2) == 0 { // even
			var axisi = int64(h.j / 2)
			var diff = int64(h.i) - axisi
			h.i = int(int64(h.i) - 2*diff)
		} else {
			var axisi = int64((h.j + 1) / 2)
			var diff = int64(h.i) - axisi
			h.i = int(int64(h.i) - (2*diff + 1))
		}
	}

	if v.y < 0.0 {
		h.i = h.i - (2*h.j+1)/2
		h.j = -1 * h.j
	}

	_ijkNormalize(h)
}

// _ijkToHex2d finds the center point in 2D cartesian coordinates of a hex.
// `h`: The ijk coordinates of the hex.
// `v`: The 2D cartesian coordinates of the hex center point.
func _ijkToHex2d(h *CoordIJK, v *Vec2d) {
	i := h.i - h.k
	j := h.j - h.k

	v.x = float64(i) - 0.5*float64(j)
	v.y = float64(j) * M_SQRT3_2
}

// _ijkMatches returns whether or not two ijk coordinates contain exactly the same
// component values.
// `c1`: The first set of ijk coordinates.
// `c2`: The second set of ijk coordinates.
// Returns true if the two addresses match, false if they do not.
func _ijkMatches(c1 *CoordIJK, c2 *CoordIJK) bool {
	return c1.i == c2.i && c1.j == c2.j && c1.k == c2.k
}

// _ijkAdd adds two ijk coordinates.
// `h1`: The first set of ijk coordinates.
// `h2`: The second set of ijk coordinates.
// `sum`: The sum of the two sets of ijk coordinates.
func _ijkAdd(h1 *CoordIJK, h2 *CoordIJK, sum *CoordIJK) {
	sum.i = h1.i + h2.i
	sum.j = h1.j + h2.j
	sum.k = h1.k + h2.k
}

// _ijkSub subtracts two ijk coordinates.
// `h1`: The first set of ijk coordinates.
// `h2`: The second set of ijk coordinates.
// `diff`: The difference of the two sets of ijk coordinates (`h1` - `h2`).
func _ijkSub(h1 *CoordIJK, h2 *CoordIJK, diff *CoordIJK) {
	diff.i = h1.i - h2.i
	diff.j = h1.j - h2.j
	diff.k = h1.k - h2.k
}

// _ijkScale uniformly scales ijk coordinates by a scalar. Works in place.
// `c`: The ijk coordinates to scale.
// `factor`: The scaling factor.
func _ijkScale(c *CoordIJK, factor int) {
	c.i *= factor
	c.j *= factor
	c.k *= factor
}

// _ijkNormalize normalizes ijk coordinates by setting the components to the smallest possible
// values. Works in place.
// `c`: The ijk coordinates to normalize.
func _ijkNormalize(c *CoordIJK) {
	// remove any negative values
	if c.i < 0 {
		c.j -= c.i
		c.k -= c.i
		c.i = 0
	}

	if c.j < 0 {
		c.i -= c.j
		c.k -= c.j
		c.j = 0
	}

	if c.k < 0 {
		c.i -= c.k
		c.j -= c.k
		c.k = 0
	}

	// remove the min value if needed
	min := c.i
	if c.j < min {
		min = c.j
	}
	if c.k < min {
		min = c.k
	}
	if min > 0 {
		c.i -= min
		c.j -= min
		c.k -= min
	}
}

// _unitIjkToDigit determines the H3 digit corresponding to a unit vector in ijk coordinates.
// `ijk`: The ijk coordinates; must be a unit vector.
// Returns the H3 digit (0-6) corresponding to the ijk unit vector, or
// INVALID_DIGIT on failure.
func _unitIjkToDigit(ijk *CoordIJK) Direction {
	c := *ijk
	_ijkNormalize(&c)

	digit := INVALID_DIGIT
	for i := CENTER_DIGIT; i < NUM_DIGITS; i++ {
		if _ijkMatches(&c, &UNIT_VECS[i]) {
			digit = i
			break
		}
	}

	return digit
}

// _upAp7 finds the normalized ijk coordinates of the indexing parent of a cell in a
// counter-clockwise aperture 7 grid. Works in place.
// `ijk`: The ijk coordinates.
func _upAp7(ijk *CoordIJK) {
	// convert to CoordIJ
	i := ijk.i - ijk.k
	j := ijk.j - ijk.k

	ijk.i = int(math.Round(float64(3*i-j) / 7.0))
	ijk.j = int(math.Round(float64(i+2*j) / 7.0))
	ijk.k = 0
	_ijkNormalize(ijk)
}

// _upAp7r finds the normalized ijk coordinates of the indexing parent of a cell in a
// clockwise aperture 7 grid. Works in place.
// `ijk`: The ijk coordinates.
func _upAp7r(ijk *CoordIJK) {
	// convert to CoordIJ
	i := ijk.i - ijk.k
	j := ijk.j - ijk.k

	ijk.i = int(math.Round(float64(2*i+j) / 7.0))
	ijk.j = int(math.Round(float64(3*j-i) / 7.0))
	ijk.k = 0
	_ijkNormalize(ijk)
}

// _downAp7 finds the normalized ijk coordinates of the hex centered on the indicated
// hex at the next finer aperture 7 counter-clockwise resolution. Works in
// place.
// `ijk`: The ijk coordinates.
func _downAp7(ijk *CoordIJK) {
	// res r unit vectors in res r+1
	iVec := CoordIJK{3, 0, 1}
	jVec := CoordIJK{1, 3, 0}
	kVec := CoordIJK{0, 1, 3}

	_ijkScale(&iVec, ijk.i)
	_ijkScale(&jVec, ijk.j)
	_ijkScale(&kVec, ijk.k)

	_ijkAdd(&iVec, &jVec, ijk)
	_ijkAdd(ijk, &kVec, ijk)

	_ijkNormalize(ijk)
}

// _downAp7r finds the normalized ijk coordinates of the hex centered on the indicated
// hex at the next finer aperture 7 clockwise resolution. Works in place.
// `ijk`: The ijk coordinates.
func _downAp7r(ijk *CoordIJK) {
	// res r unit vectors in res r+1
	iVec := CoordIJK{3, 1, 0}
	jVec := CoordIJK{0, 3, 1}
	kVec := CoordIJK{1, 0, 3}

	_ijkScale(&iVec, ijk.i)
	_ijkScale(&jVec, ijk.j)
	_ijkScale(&kVec, ijk.k)

	_ijkAdd(&iVec, &jVec, ijk)
	_ijkAdd(ijk, &kVec, ijk)

	_ijkNormalize(ijk)
}

// _neighbor finds the normalized ijk coordinates of the hex in the specified digit
// direction from the specified ijk coordinates. Works in place.
// `ijk`: The ijk coordinates.
// `digit`: The digit direction from the original ijk coordinates.
func _neighbor(ijk *CoordIJK, digit Direction) {
	if digit > CENTER_DIGIT && digit < NUM_DIGITS {
		_ijkAdd(ijk, &UNIT_VECS[digit], ijk)
		_ijkNormalize(ijk)
	}
}

// _ijkRotate60ccw rotates ijk coordinates 60 degrees counter-clockwise. Works in place.
// `ijk`: The ijk coordinates.
func _ijkRotate60ccw(ijk *CoordIJK) {
	// unit vector rotations
	iVec := CoordIJK{1, 1, 0}
	jVec := CoordIJK{0, 1, 1}
	kVec := CoordIJK{1, 0, 1}

	_ijkScale(&iVec, ijk.i)
	_ijkScale(&jVec, ijk.j)
	_ijkScale(&kVec, ijk.k)

	_ijkAdd(&iVec, &jVec, ijk)
	_ijkAdd(ijk, &kVec, ijk)

	_ijkNormalize(ijk)
}

// _ijkRotate60cw rotates ijk coordinates 60 degrees clockwise. Works in place.
// `ijk`: The ijk coordinates.
func _ijkRotate60cw(ijk *CoordIJK) {
	// unit vector rotations
	iVec := CoordIJK{1, 0, 1}
	jVec := CoordIJK{1, 1, 0}
	kVec := CoordIJK{0, 1, 1}

	_ijkScale(&iVec, ijk.i)
	_ijkScale(&jVec, ijk.j)
	_ijkScale(&kVec, ijk.k)

	_ijkAdd(&iVec, &jVec, ijk)
	_ijkAdd(ijk, &kVec, ijk)

	_ijkNormalize(ijk)
}

// _rotate60ccw rotates indexing digit 60 degrees counter-clockwise. Returns result.
// `digit`: Indexing digit (between 1 and 6 inclusive)
func _rotate60ccw(digit Direction) Direction {
	switch digit {
	case K_AXES_DIGIT:
		return IK_AXES_DIGIT
	case IK_AXES_DIGIT:
		return I_AXES_DIGIT
	case I_AXES_DIGIT:
		return IJ_AXES_DIGIT
	case IJ_AXES_DIGIT:
		return J_AXES_DIGIT
	case J_AXES_DIGIT:
		return JK_AXES_DIGIT
	case JK_AXES_DIGIT:
		return K_AXES_DIGIT
	default:
		return digit
	}
}

// _rotate60cw rotates indexing digit 60 degrees clockwise. Returns result.
// `digit`: Indexing digit (between 1 and 6 inclusive)
func _rotate60cw(digit Direction) Direction {
	switch digit {
	case K_AXES_DIGIT:
		return JK_AXES_DIGIT
	case JK_AXES_DIGIT:
		return J_AXES_DIGIT
	case J_AXES_DIGIT:
		return IJ_AXES_DIGIT
	case IJ_AXES_DIGIT:
		return I_AXES_DIGIT
	case I_AXES_DIGIT:
		return IK_AXES_DIGIT
	case IK_AXES_DIGIT:
		return K_AXES_DIGIT
	default:
		return digit
	}
}

// _downAp3 finds the normalized ijk coordinates of the hex centered on the indicated
// hex at the next finer aperture 3 counter-clockwise resolution. Works in
// place.
// `ijk`: The ijk coordinates.
func _downAp3(ijk *CoordIJK) {
	// res r unit vectors in res r+1
	iVec := CoordIJK{2, 0, 1}
	jVec := CoordIJK{1, 2, 0}
	kVec := CoordIJK{0, 1, 2}

	_ijkScale(&iVec, ijk.i)
	_ijkScale(&jVec, ijk.j)
	_ijkScale(&kVec, ijk.k)

	_ijkAdd(&iVec, &jVec, ijk)
	_ijkAdd(ijk, &kVec, ijk)

	_ijkNormalize(ijk)
}

// _downAp3r finds the normalized ijk coordinates of the hex centered on the indicated
// hex at the next finer aperture 3 clockwise resolution. Works in place.
// `ijk`: The ijk coordinates.
func _downAp3r(ijk *CoordIJK) {
	// res r unit vectors in res r+1
	iVec := CoordIJK{2, 1, 0}
	jVec := CoordIJK{0, 2, 1}
	kVec := CoordIJK{1, 0, 2}

	_ijkScale(&iVec, ijk.i)
	_ijkScale(&jVec, ijk.j)
	_ijkScale(&kVec, ijk.k)

	_ijkAdd(&iVec, &jVec, ijk)
	_ijkAdd(ijk, &kVec, ijk)

	_ijkNormalize(ijk)
}

// ijkDistance finds the distance between the two coordinates. Returns result.
// `c1`: The first set of ijk coordinates.
// `c2`: The second set of ijk coordinates.
func ijkDistance(c1 *CoordIJK, c2 *CoordIJK) int {
	var diff CoordIJK
	_ijkSub(c1, c2, &diff)
	_ijkNormalize(&diff)
	absDiff := CoordIJK{absInt(diff.i), absInt(diff.j), absInt(diff.k)}
	return maxInt(absDiff.i, maxInt(absDiff.j, absDiff.k))
}

// ijkToIj transforms coordinates from the IJK+ coordinate system to the IJ coordinate
// system.
// `ijk`: The input IJK+ coordinates
// `ij`: The output IJ coordinates
func ijkToIj(ijk *CoordIJK, ij *CoordIJ) {
	ij.i = ijk.i - ijk.k
	ij.j = ijk.j - ijk.k
}

// ijToIjk transforms coordinates from the IJ coordinate system to the IJK+ coordinate
// system.
// `ij`: The input IJ coordinates
// `ijk`: The output IJK+ coordinates
func ijToIjk(ij *CoordIJ, ijk *CoordIJK) {
	ijk.i = ij.i
	ijk.j = ij.j
	ijk.k = 0

	_ijkNormalize(ijk)
}

// ijkToCube converts IJK coordinates to cube coordinates, in place
// `ijk`: Coordinate to convert
func ijkToCube(ijk *CoordIJK) {
	ijk.i = -ijk.i + ijk.k
	ijk.j = ijk.j - ijk.k
	ijk.k = -ijk.i - ijk.j
}

// cubeToIjk convert cube coordinates to IJK coordinates, in place
// `ijk`: Coordinate to convert
func cubeToIjk(ijk *CoordIJK) {
	ijk.i = -ijk.i
	ijk.k = 0
	_ijkNormalize(ijk)
}
