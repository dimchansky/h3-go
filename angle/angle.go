// Package angle provides a safe angle type with internal storage in radians.
package angle

import (
	"fmt"
	"math"
)

// Angle represents an angle stored internally in radians.
// It provides type safety and convenience methods for angle operations.
type Angle float64

// Common angle constants and conversion factors.
const (
	// RadPerDeg is the conversion factor from degrees to radians.
	RadPerDeg = math.Pi / 180
	
	// DegPerRad is the conversion factor from radians to degrees.
	DegPerRad = 180 / math.Pi

	// Pi represents π radians as an Angle.
	Pi = Angle(math.Pi)
	
	// TwoPi represents 2π radians as an Angle.
	TwoPi = Angle(2 * math.Pi)
)

// Rad creates an Angle from a value in radians.
func Rad(v float64) Angle { return Angle(v) }

// Deg creates an Angle from a value in degrees.
func Deg(v float64) Angle { return Angle(v * RadPerDeg) }

// Rad returns the angle value in radians.
func (a Angle) Rad() float64 { return float64(a) }

// Deg returns the angle value in degrees.
func (a Angle) Deg() float64 { return float64(a) * DegPerRad }

// Add returns the sum of two angles.
func (a Angle) Add(b Angle) Angle { return a + b }

// Sub returns the difference between two angles.
func (a Angle) Sub(b Angle) Angle { return a - b }

// Mul returns the angle multiplied by a scalar.
func (a Angle) Mul(s float64) Angle { return Angle(float64(a) * s) }

// Div returns the angle divided by a scalar.
func (a Angle) Div(s float64) Angle { return Angle(float64(a) / s) }

// Neg returns the negation of the angle.
func (a Angle) Neg() Angle { return -a }

// WrapTwoPi normalizes the angle to the range [0, 2π).
func (a Angle) WrapTwoPi() Angle {
	r := math.Mod(float64(a), 2*math.Pi)
	if r < 0 {
		r += 2 * math.Pi
	}
	return Angle(r)
}

// WrapPi normalizes the angle to the range (-π, π].
func (a Angle) WrapPi() Angle {
	r := math.Mod(float64(a), 2*math.Pi)
	if r <= -math.Pi {
		r += 2 * math.Pi
	} else if r > math.Pi {
		r -= 2 * math.Pi
	}
	return Angle(r)
}

// EqualApprox compares two angles for approximate equality within the given epsilon.
// If eps is zero or negative, a default strict epsilon of 1e-12 is used.
func (a Angle) EqualApprox(b Angle, eps float64) bool {
	if eps <= 0 {
		eps = 1e-12
	}
	return math.Abs(float64(a-b)) <= eps
}

// Sin returns the sine of the angle.
func (a Angle) Sin() float64 { return math.Sin(float64(a)) }

// Cos returns the cosine of the angle.
func (a Angle) Cos() float64 { return math.Cos(float64(a)) }

// Tan returns the tangent of the angle.
func (a Angle) Tan() float64 { return math.Tan(float64(a)) }

// SinCos returns the sine and cosine of the angle.
func (a Angle) SinCos() (float64, float64) { return math.Sincos(float64(a)) }

// String returns a human-readable representation of the angle in degrees with the ° symbol.
func (a Angle) String() string { return fmt.Sprintf("%.6f°", a.Deg()) }
