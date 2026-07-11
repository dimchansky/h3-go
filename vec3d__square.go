package h3

// _square returns the square of a number.
// Ported from H3 C: vec3d.c::_square.
//
//nolint:unused // ported from H3 C for parity completeness; exercised by cgo && c2go parity tests
func _square(x float64) float64 {
	return x * x
}
