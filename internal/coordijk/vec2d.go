package coordijk

import "math"

// Vec2d represents a 2D cartesian coordinate vector.
type Vec2d struct {
	X float64
	Y float64
}

// Magnitude returns the magnitude (length) of the vector.
func (v Vec2d) Magnitude() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

// Add adds two vectors.
func (v Vec2d) Add(other Vec2d) Vec2d {
	return Vec2d{
		X: v.X + other.X,
		Y: v.Y + other.Y,
	}
}

// Sub subtracts two vectors.
func (v Vec2d) Sub(other Vec2d) Vec2d {
	return Vec2d{
		X: v.X - other.X,
		Y: v.Y - other.Y,
	}
}

// Scale multiplies the vector by a scalar.
func (v Vec2d) Scale(s float64) Vec2d {
	return Vec2d{
		X: v.X * s,
		Y: v.Y * s,
	}
}

// Normalize returns a unit vector in the same direction.
func (v Vec2d) Normalize() Vec2d {
	mag := v.Magnitude()
	if mag == 0 {
		return v
	}
	return v.Scale(1.0 / mag)
}
