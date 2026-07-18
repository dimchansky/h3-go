//go:build cgo && c2go

package h3

// Shared float-comparison helpers for the parity suite. They previously
// lived in latLng__cellAreaRads2_parity_test.go, which is gated !h3v450
// (docs/sync/h3v450-exclusion-inventory.md), while ungated tests such as
// the cellAreaKm2/M2 parity tests still need them in both configurations.

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
