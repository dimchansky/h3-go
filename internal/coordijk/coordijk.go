// Package coordijk implements the H3 IJK axial coordinate system and
// operations on integer IJK coordinates used in face-local hex grids.
// IJK coordinates are a hexagonal coordinate system with 3 axes,
// each 120 degrees apart.
package coordijk

import (
	"math"
)

// Direction represents H3 digit values (0-7).
type Direction int

const (
	// CenterDigit represents the center of a hexagon (0).
	CenterDigit Direction = 0
	// KAxesDigit represents movement in the k-axis direction (1).
	KAxesDigit Direction = 1
	// JAxesDigit represents movement in the j-axis direction (2).
	JAxesDigit Direction = 2
	// JKAxesDigit represents movement in the j==k direction (3).
	JKAxesDigit Direction = 3
	// IAxesDigit represents movement in the i-axis direction (4).
	IAxesDigit Direction = 4
	// IKAxesDigit represents movement in the i==k direction (5).
	IKAxesDigit Direction = 5
	// IJAxesDigit represents movement in the i==j direction (6).
	IJAxesDigit Direction = 6
	// InvalidDigit represents an invalid direction (7).
	InvalidDigit Direction = 7
	// NumDigits is the number of valid digits (same as InvalidDigit).
	NumDigits = InvalidDigit
	// PentagonSkippedDigit is the digit skipped for pentagons (K_AXES_DIGIT).
	PentagonSkippedDigit = KAxesDigit
)

// CoordIJK represents a point in the IJK hexagonal coordinate system.
// Each axis is spaced 120 degrees apart.
type CoordIJK struct {
	I int // i component
	J int // j component
	K int // k component
}

// UnitVecs contains the IJK unit vectors corresponding to the 7 H3 digits.
var UnitVecs = [7]CoordIJK{
	{0, 0, 0}, // direction 0 (center)
	{0, 0, 1}, // direction 1 (k)
	{0, 1, 0}, // direction 2 (j)
	{0, 1, 1}, // direction 3 (j+k)
	{1, 0, 0}, // direction 4 (i)
	{1, 0, 1}, // direction 5 (i+k)
	{1, 1, 0}, // direction 6 (i+j)
}

// M_SQRT3_2 is sqrt(3) / 2.
const M_SQRT3_2 = 0.8660254037844386467637231707529361834714 //nolint:revive // keep H3 naming for constants

// M_RSIN60 is 1 / sin(60 degrees) = 2 / sqrt(3).
const M_RSIN60 = 1.1547005383792515290182975610039149269484 //nolint:revive // keep H3 naming for constants

// Add adds two IJK coordinates in-place and returns the modified receiver.
func (c *CoordIJK) Add(other CoordIJK) *CoordIJK {
	c.I += other.I
	c.J += other.J
	c.K += other.K
	return c
}

// Sub subtracts two IJK coordinates in-place and returns the modified receiver.
func (c *CoordIJK) Sub(other CoordIJK) *CoordIJK {
	c.I -= other.I
	c.J -= other.J
	c.K -= other.K
	return c
}

// Scale multiplies all components by a factor in-place and returns the modified receiver.
func (c *CoordIJK) Scale(factor int) *CoordIJK {
	c.I *= factor
	c.J *= factor
	c.K *= factor
	return c
}

// Normalize normalizes the IJK coordinates so that i+j+k=0.
// In the IJK coordinate system, valid coordinates always sum to 0.
// Returns the modified receiver for chaining.
func (c *CoordIJK) Normalize() *CoordIJK {
	// Exactly match H3 C _ijkNormalize algorithm
	// Step 1: Remove any negative values
	if c.I < 0 {
		c.J -= c.I
		c.K -= c.I
		c.I = 0
	}

	if c.J < 0 {
		c.I -= c.J
		c.K -= c.J
		c.J = 0
	}

	if c.K < 0 {
		c.I -= c.K
		c.J -= c.K
		c.K = 0
	}

	// Step 2: Remove the min value if needed
	minComponent := c.I
	if c.J < minComponent {
		minComponent = c.J
	}
	if c.K < minComponent {
		minComponent = c.K
	}
	if minComponent > 0 {
		c.I -= minComponent
		c.J -= minComponent
		c.K -= minComponent
	}
	return c
}

// Neighbor moves the IJK coordinate in the specified direction and returns the modified receiver.
func (c *CoordIJK) Neighbor(dir Direction) *CoordIJK {
	// Mirror H3 C _neighbor: only act for digits (1..6). CENTER_DIGIT is no-op.
	if dir > CenterDigit && dir < NumDigits {
		c.Add(UnitVecs[dir]).Normalize()
	}
	return c
}

// Rotate60CCW rotates the IJK coordinates 60 degrees counter-clockwise.
// Returns the modified receiver for chaining.
func (c *CoordIJK) Rotate60CCW() *CoordIJK {
	// Match H3 C _ijkRotate60ccw by recombining unit vectors and normalizing.
	iVec := CoordIJK{1, 1, 0}
	jVec := CoordIJK{0, 1, 1}
	kVec := CoordIJK{1, 0, 1}

	iVec.Scale(c.I)
	jVec.Scale(c.J)
	kVec.Scale(c.K)

	*c = iVec
	return c.Add(jVec).Add(kVec).Normalize()
}

// Rotate60CW rotates the IJK coordinates 60 degrees clockwise.
// Returns the modified receiver for chaining.
func (c *CoordIJK) Rotate60CW() *CoordIJK {
	// Match H3 C _ijkRotate60cw by recombining unit vectors and normalizing.
	iVec := CoordIJK{1, 0, 1}
	jVec := CoordIJK{1, 1, 0}
	kVec := CoordIJK{0, 1, 1}

	iVec.Scale(c.I)
	jVec.Scale(c.J)
	kVec.Scale(c.K)

	*c = iVec
	return c.Add(jVec).Add(kVec).Normalize()
}

// DownAp7 transforms coordinates for a resolution decrease from Class III to Class II.
// This is the aperture 7 transformation for going down in resolution.
// Returns the modified receiver for chaining.
func (c *CoordIJK) DownAp7() *CoordIJK {
	// From H3 source: res r unit vectors in res r+1
	iVec := CoordIJK{3, 0, 1}
	jVec := CoordIJK{1, 3, 0}
	kVec := CoordIJK{0, 1, 3}

	iVec.Scale(c.I)
	jVec.Scale(c.J)
	kVec.Scale(c.K)

	*c = iVec
	return c.Add(jVec).Add(kVec).Normalize()
}

// DownAp7r transforms coordinates for a resolution decrease from Class II to Class III.
// This is the reverse aperture 7 transformation.
// Returns the modified receiver for chaining.
func (c *CoordIJK) DownAp7r() *CoordIJK {
	// From H3 source: res r unit vectors in res r+1
	iVec := CoordIJK{3, 1, 0}
	jVec := CoordIJK{0, 3, 1}
	kVec := CoordIJK{1, 0, 3}

	iVec.Scale(c.I)
	jVec.Scale(c.J)
	kVec.Scale(c.K)

	*c = iVec
	return c.Add(jVec).Add(kVec).Normalize()
}

// DownAp3 transforms coordinates for a resolution decrease in Class III.
// This is the aperture 3 transformation (used for pentagons).
// Returns the modified receiver for chaining.
func (c *CoordIJK) DownAp3() *CoordIJK {
	// From H3 source: res r unit vectors in res r+1
	iVec := CoordIJK{2, 0, 1}
	jVec := CoordIJK{1, 2, 0}
	kVec := CoordIJK{0, 1, 2}

	iVec.Scale(c.I)
	jVec.Scale(c.J)
	kVec.Scale(c.K)

	*c = iVec
	return c.Add(jVec).Add(kVec).Normalize()
}

// DownAp3r transforms coordinates for a resolution decrease in Class II.
// This is the reverse aperture 3 transformation.
// Returns the modified receiver for chaining.
func (c *CoordIJK) DownAp3r() *CoordIJK {
	// From H3 source: res r unit vectors in res r+1
	iVec := CoordIJK{2, 1, 0}
	jVec := CoordIJK{0, 2, 1}
	kVec := CoordIJK{1, 0, 2}

	iVec.Scale(c.I)
	jVec.Scale(c.J)
	kVec.Scale(c.K)

	*c = iVec
	return c.Add(jVec).Add(kVec).Normalize()
}

// UnitIJKToDigit determines the H3 digit corresponding to a unit IJK coordinate.
// Normalizes the input coordinate before checking, following H3 C _unitIjkToDigit.
func UnitIJKToDigit(ijk CoordIJK) Direction {
	// Normalize the coordinate first, matching H3 C behavior
	ijk.Normalize()

	// Find which unit vector this matches
	for dir := Direction(0); dir < NumDigits; dir++ {
		if ijk == UnitVecs[dir] {
			return dir
		}
	}
	return InvalidDigit
}

// abs returns the absolute value of an integer.
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// Distance computes the grid distance between two IJK coordinates.
// This matches H3 C ijkDistance function behavior.
func Distance(a, b CoordIJK) int {
	// Match H3 C ijkDistance: normalize(a-b), then return max(|i|,|j|,|k|).
	diff := a.Sub(b).Normalize()
	ai, aj, ak := abs(diff.I), abs(diff.J), abs(diff.K)
	if ai < aj {
		ai = aj
	}
	if ai < ak {
		ai = ak
	}
	return ai
}

// UpAp7 finds the normalized IJK coordinates of the indexing parent
// in a counter-clockwise aperture 7 grid. Works in place.
// Returns the modified receiver for chaining.
func (c *CoordIJK) UpAp7() *CoordIJK {
	// Convert to CoordIJ
	i := c.I - c.K
	j := c.J - c.K

	// Apply aperture 7 transformation
	c.I = int(math.Round(float64(3*i-j) / 7.0))
	c.J = int(math.Round(float64(i+2*j) / 7.0))
	c.K = 0
	return c.Normalize()
}

// UpAp7r finds the normalized IJK coordinates of the indexing parent
// in a clockwise aperture 7 grid. Works in place.
// Returns the modified receiver for chaining.
func (c *CoordIJK) UpAp7r() *CoordIJK {
	// Convert to CoordIJ
	i := c.I - c.K
	j := c.J - c.K

	// Apply reverse aperture 7 transformation
	c.I = int(math.Round(float64(2*i+j) / 7.0))
	c.J = int(math.Round(float64(3*j-i) / 7.0))
	c.K = 0
	return c.Normalize()
}

// ToHex2d converts IJK coordinates to 2D hex coordinates.
// The Vec2d represents a point in the hex2d coordinate system.
func (c *CoordIJK) ToHex2d() Vec2d {
	i := c.I - c.K
	j := c.J - c.K

	return Vec2d{
		X: float64(i) - 0.5*float64(j),
		Y: float64(j) * M_SQRT3_2,
	}
}

// Hex2dToCoordIJK converts 2D hex coordinates to IJK coordinates.
func Hex2dToCoordIJK(v Vec2d) CoordIJK {
	a1 := math.Abs(v.X)
	a2 := math.Abs(v.Y)

	// Reverse conversion
	x2 := a2 * M_RSIN60
	x1 := a1 + x2/2.0

	// Quantize
	m1 := int(x1)
	m2 := int(x2)

	// Round correctly
	r1 := x1 - float64(m1)
	r2 := x2 - float64(m2)

	var h CoordIJK
	h.K = 0

	if r1 < 0.5 {
		if r1 < 1.0/3.0 {
			if r2 < (1.0+r1)/2.0 {
				h.I = m1
				h.J = m2
			} else {
				h.I = m1
				h.J = m2 + 1
			}
		} else {
			if r2 < (1.0 - r1) {
				h.J = m2
			} else {
				h.J = m2 + 1
			}

			if (1.0-r1) <= r2 && r2 < (2.0*r1) {
				h.I = m1 + 1
			} else {
				h.I = m1
			}
		}
	} else {
		if r1 < 2.0/3.0 {
			if r2 < (1.0 - r1) {
				h.J = m2
			} else {
				h.J = m2 + 1
			}

			if (2.0*r1-1.0) < r2 && r2 < (1.0-r1) {
				h.I = m1
			} else {
				h.I = m1 + 1
			}
		} else {
			if r2 < r1/2.0 {
				h.I = m1 + 1
				h.J = m2
			} else {
				h.I = m1 + 1
				h.J = m2 + 1
			}
		}
	}

	// Fold across axes if necessary
	if v.X < 0.0 {
		if h.J%2 == 0 {
			axisi := h.J / 2
			diff := h.I - axisi
			h.I -= 2 * diff
		} else {
			axisi := (h.J + 1) / 2
			diff := h.I - axisi
			h.I -= (2*diff + 1)
		}
	}

	if v.Y < 0.0 {
		h.I -= (2*h.J + 1) / 2
		h.J = -h.J
	}

	h.Normalize()
	return h
}
