package h3

// vec3Dot returns the dot product of two 3D vectors.
// Ported from H3 C: vec3d.h::vec3Dot.
//
// The explicit float64 conversions force each product to be rounded
// before the additions: the Go spec otherwise permits the compiler to
// fuse a*b+c into an FMA (which gc does on arm64, even across
// statements), while the parity harness compiles the reference C with
// -ffp-contract=off. Bit-exact parity requires strict IEEE evaluation on
// both sides; a conversion is the spec-guaranteed way to force rounding.
func vec3Dot(v1, v2 vec3d) float64 {
	return float64(v1.X*v2.X) + float64(v1.Y*v2.Y) + float64(v1.Z*v2.Z)
}
