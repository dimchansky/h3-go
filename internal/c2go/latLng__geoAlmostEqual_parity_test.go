//go:build cgo

package c2go

import "testing"

func Test_geoAlmostEqual_ParityWithC(t *testing.T) {
	cases := []struct {
		a LatLng
		b LatLng
	}{
		{LatLng{0, 0}, LatLng{0, 0}},
		{LatLng{0.1, -0.2}, LatLng{0.1 + 1e-12, -0.2 - 1e-12}},
		{LatLng{1, 2}, LatLng{1.0, 2.0 + 1e-8}},
		{LatLng{-1.5, 3.2}, LatLng{-1.5 + 1e-10, 3.2 + 1e-10}},
	}
	for _, tc := range cases {
		goVal := geoAlmostEqual(&tc.a, &tc.b)
		cVal := geoAlmostEqualC(tc.a, tc.b)
		if goVal != cVal {
			t.Fatalf("geoAlmostEqual parity mismatch: go=%v c=%v", goVal, cVal)
		}
	}
}
