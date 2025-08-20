// Package coordijk implements the H3 IJK coordinate system.
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

// M_SQRT3_2 is sqrt(3) / 2
const M_SQRT3_2 = 0.8660254037844386467637231707529361834714

// M_RSIN60 is 1 / sin(60 degrees) = 2 / sqrt(3)
const M_RSIN60 = 1.1547005383792515290182975610039149269484

// Add adds two IJK coordinates.
func (c CoordIJK) Add(other CoordIJK) CoordIJK {
	return CoordIJK{
		I: c.I + other.I,
		J: c.J + other.J,
		K: c.K + other.K,
	}
}

// Sub subtracts two IJK coordinates.
func (c CoordIJK) Sub(other CoordIJK) CoordIJK {
	return CoordIJK{
		I: c.I - other.I,
		J: c.J - other.J,
		K: c.K - other.K,
	}
}

// Scale multiplies all components by a factor.
func (c CoordIJK) Scale(factor int) CoordIJK {
	return CoordIJK{
		I: c.I * factor,
		J: c.J * factor,
		K: c.K * factor,
	}
}

// Normalize normalizes the IJK coordinates so that i+j+k=0.
// In the IJK coordinate system, valid coordinates always sum to 0.
func (c *CoordIJK) Normalize() {
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
	min := c.I
	if c.J < min {
		min = c.J
	}
	if c.K < min {
		min = c.K
	}
	if min > 0 {
		c.I -= min
		c.J -= min
		c.K -= min
	}
}

// Neighbor moves the IJK coordinate in the specified direction.
func (c *CoordIJK) Neighbor(dir Direction) {
	if dir >= 0 && dir < NumDigits {
		*c = c.Add(UnitVecs[dir])
		c.Normalize()
	}
}

// Rotate60CCW rotates the IJK coordinates 60 degrees counter-clockwise.
func (c *CoordIJK) Rotate60CCW() {
	// Rotation: i' = -k, j' = -i, k' = -j
	i := -c.K
	j := -c.I
	k := -c.J
	c.I = i
	c.J = j
	c.K = k
}

// Rotate60CW rotates the IJK coordinates 60 degrees clockwise.
func (c *CoordIJK) Rotate60CW() {
	// Rotation: i' = -j, j' = -k, k' = -i
	i := -c.J
	j := -c.K
	k := -c.I
	c.I = i
	c.J = j
	c.K = k
}

// DownAp7 transforms coordinates for a resolution decrease from Class III to Class II.
// This is the aperture 7 transformation for going down in resolution.
func (c *CoordIJK) DownAp7() {
	// From H3 source: res r unit vectors in res r+1
	iVec := CoordIJK{3, 0, 1}
	jVec := CoordIJK{1, 3, 0}
	kVec := CoordIJK{0, 1, 3}
	
	*c = iVec.Scale(c.I).Add(jVec.Scale(c.J)).Add(kVec.Scale(c.K))
	c.Normalize()
}

// DownAp7r transforms coordinates for a resolution decrease from Class II to Class III.
// This is the reverse aperture 7 transformation.
func (c *CoordIJK) DownAp7r() {
	// From H3 source: res r unit vectors in res r+1
	iVec := CoordIJK{3, 1, 0}
	jVec := CoordIJK{0, 3, 1}
	kVec := CoordIJK{1, 0, 3}
	
	*c = iVec.Scale(c.I).Add(jVec.Scale(c.J)).Add(kVec.Scale(c.K))
	c.Normalize()
}

// DownAp3 transforms coordinates for a resolution decrease in Class III.
// This is the aperture 3 transformation (used for pentagons).
func (c *CoordIJK) DownAp3() {
	// From H3 source: res r unit vectors in res r+1
	iVec := CoordIJK{2, 0, 1}
	jVec := CoordIJK{1, 2, 0}
	kVec := CoordIJK{0, 1, 2}
	
	*c = iVec.Scale(c.I).Add(jVec.Scale(c.J)).Add(kVec.Scale(c.K))
	c.Normalize()
}

// DownAp3r transforms coordinates for a resolution decrease in Class II.
// This is the reverse aperture 3 transformation.
func (c *CoordIJK) DownAp3r() {
	// From H3 source: res r unit vectors in res r+1
	iVec := CoordIJK{2, 1, 0}
	jVec := CoordIJK{0, 2, 1}
	kVec := CoordIJK{1, 0, 2}
	
	*c = iVec.Scale(c.I).Add(jVec.Scale(c.J)).Add(kVec.Scale(c.K))
	c.Normalize()
}

// UnitIJKToDigit determines the H3 digit corresponding to a unit IJK coordinate.
// Normalizes the input coordinate before checking, following H3 C _unitIjkToDigit.
func UnitIJKToDigit(ijk CoordIJK) Direction {
	// Normalize the coordinate first, matching H3 C behavior
	normalized := ijk
	normalized.Normalize()
	
	// Find which unit vector this matches
	for dir := Direction(0); dir < NumDigits; dir++ {
		if normalized == UnitVecs[dir] {
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
	diff := a.Sub(b)
	diff.Normalize() // Important: normalize the difference as H3 C does
	return (abs(diff.I) + abs(diff.J) + abs(diff.K)) / 2
}

// UpAp7 finds the normalized IJK coordinates of the indexing parent
// in a counter-clockwise aperture 7 grid. Works in place.
func (c *CoordIJK) UpAp7() {
	// Convert to CoordIJ
	i := c.I - c.K
	j := c.J - c.K
	
	// Apply aperture 7 transformation
	c.I = int(math.Round(float64(3*i-j) / 7.0))
	c.J = int(math.Round(float64(i+2*j) / 7.0))
	c.K = 0
	c.Normalize()
}

// UpAp7r finds the normalized IJK coordinates of the indexing parent
// in a clockwise aperture 7 grid. Works in place.
func (c *CoordIJK) UpAp7r() {
	// Convert to CoordIJ  
	i := c.I - c.K
	j := c.J - c.K
	
	// Apply reverse aperture 7 transformation
	c.I = int(math.Round(float64(2*i+j) / 7.0))
	c.J = int(math.Round(float64(3*j-i) / 7.0))
	c.K = 0
	c.Normalize()
}

// ToHex2d converts IJK coordinates to 2D hex coordinates.
// The Vec2d represents a point in the hex2d coordinate system.
func (c CoordIJK) ToHex2d() Vec2d {
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
			h.I = h.I - 2*diff
		} else {
			axisi := (h.J + 1) / 2
			diff := h.I - axisi
			h.I = h.I - (2*diff + 1)
		}
	}
	
	if v.Y < 0.0 {
		h.I = h.I - (2*h.J+1)/2
		h.J = -h.J
	}
	
	h.Normalize()
	return h
}