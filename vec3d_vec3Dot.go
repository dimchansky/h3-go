package h3

// vec3Dot returns the dot product of two 3D vectors.
// Ported from H3 C: vec3d.h::vec3Dot.
//
// The explicit float64 conversions force each product to be rounded
// before the additions: the Go spec permits the compiler to fuse a*b+c
// into an FMA unless an explicit conversion forces rounding, and gc
// does fuse on arm64 (even across statements). The parity harness pins
// the reference C to strict IEEE with -ffp-contract=off (see the
// test-c2go comment in the Makefile), so bit-exact parity requires the
// same strict evaluation on the Go side, which these conversions
// guarantee.
func vec3Dot(v1, v2 vec3d) float64 {
	return float64(v1.X*v2.X) + float64(v1.Y*v2.Y) + float64(v1.Z*v2.Z)
}
