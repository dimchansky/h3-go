// Package angles provides utilities for angle conversions and normalization.
package angles

import "math"

// Constants for angle conversions and comparisons.
const (
	// EpsRad is a small epsilon in radians used for angular comparisons.
	EpsRad = 1e-12

	// EpsDeg is a small epsilon in degrees used for geographic comparisons.
	EpsDeg = 1e-9

	// DegToRad is the conversion factor from degrees to radians.
	DegToRad = math.Pi / 180.0

	// RadToDeg is the conversion factor from radians to degrees.
	RadToDeg = 180.0 / math.Pi
)

// DegreesToRadians converts degrees to radians.
func DegreesToRadians(deg float64) float64 {
	return deg * DegToRad
}

// RadiansToDegrees converts radians to degrees.
func RadiansToDegrees(rad float64) float64 {
	return rad * RadToDeg
}

// NormalizeRadians normalizes an angle in radians to the range [0, 2π).
func NormalizeRadians(rad float64) float64 {
	// First reduce to [-2π, 2π]
	rad = math.Mod(rad, 2*math.Pi)
	
	// Then shift to [0, 2π)
	if rad < 0 {
		rad += 2 * math.Pi
	}
	
	return rad
}

// NormalizeDegrees normalizes an angle in degrees to the range [0, 360).
func NormalizeDegrees(deg float64) float64 {
	// First reduce to [-360, 360]
	deg = math.Mod(deg, 360.0)
	
	// Then shift to [0, 360)
	if deg < 0 {
		deg += 360.0
	}
	
	return deg
}

// NormalizeLongitude normalizes longitude to the range (-180, 180].
func NormalizeLongitude(lng float64) float64 {
	// Normalize to [0, 360)
	lng = NormalizeDegrees(lng)
	
	// Shift to (-180, 180]
	if lng > 180.0 {
		lng -= 360.0
	}
	
	return lng
}

// ClampLatitude clamps latitude to the range [-90, 90].
func ClampLatitude(lat float64) float64 {
	if lat < -90.0 {
		return -90.0
	}
	if lat > 90.0 {
		return 90.0
	}
	return lat
}

// ApproxEqualRad checks if two angles in radians are approximately equal.
func ApproxEqualRad(a, b float64) bool {
	return math.Abs(a-b) < EpsRad
}

// ApproxEqualDeg checks if two angles in degrees are approximately equal.
func ApproxEqualDeg(a, b float64) bool {
	return math.Abs(a-b) < EpsDeg
}

// WrapAngleRad wraps an angle in radians to the range [-π, π].
func WrapAngleRad(rad float64) float64 {
	rad = NormalizeRadians(rad)
	if rad > math.Pi {
		rad -= 2 * math.Pi
	}
	return rad
}

// WrapAngleDeg wraps an angle in degrees to the range [-180, 180].
func WrapAngleDeg(deg float64) float64 {
	deg = NormalizeDegrees(deg)
	if deg > 180.0 {
		deg -= 360.0
	}
	return deg
}