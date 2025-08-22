package c2go

import "math"

// _v3dNormalize returns the unit-length direction of v (or zero vector if |v|==0).
// Mirrors H3's vec3d.c::_v3dNormalize behavior.
func _v3dNormalize(v *Vec3d) Vec3d {
	m := math.Sqrt(v.X*v.X + v.Y*v.Y + v.Z*v.Z)
	if m == 0 {
		return Vec3d{}
	}
	inv := 1 / m
	return Vec3d{X: v.X * inv, Y: v.Y * inv, Z: v.Z * inv}
}
