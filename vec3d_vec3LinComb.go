package h3

// vec3LinComb returns the linear combination a*v1 + b*v2.
// Ported from H3 C: vec3d.h::vec3LinComb.
//
// The explicit float64 conversions force the products to round before
// the additions, defeating gc's arm64 FMA fusion; see vec3Dot.
func vec3LinComb(a float64, v1 vec3d, b float64, v2 vec3d) vec3d {
	out := vec3d{
		X: float64(a*v1.X) + float64(b*v2.X),
		Y: float64(a*v1.Y) + float64(b*v2.Y),
		Z: float64(a*v1.Z) + float64(b*v2.Z),
	}
	return out
}
