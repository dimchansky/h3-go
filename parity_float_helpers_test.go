//go:build cgo && c2go

package h3

// Shared float-comparison helpers for the parity suite, used across
// several parity test files (e.g. the cellAreaKm2/M2 and vec3 suites).

// floatAlmostEqual is a helper function for floating point comparison.
func floatAlmostEqual(a, b, tolerance float64) bool {
	return abs(a-b) <= tolerance
}

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}
