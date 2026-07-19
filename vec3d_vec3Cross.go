package h3

// vec3Cross returns the cross product of two 3D vectors.
// Ported from H3 C: vec3d.h::vec3Cross.
//
// The explicit float64 conversions force the products to round before
// the subtractions, defeating gc's arm64 FMA fusion; see vec3Dot.
func vec3Cross(v1, v2 vec3d) vec3d {
	out := vec3d{
		X: float64(v1.Y*v2.Z) - float64(v1.Z*v2.Y),
		Y: float64(v1.Z*v2.X) - float64(v1.X*v2.Z),
		Z: float64(v1.X*v2.Y) - float64(v1.Y*v2.X),
	}
	return out
}
