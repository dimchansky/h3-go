//go:build cgo && c2go

package h3

import "testing"

func Test_geoAlmostEqualThreshold_ParityWithC(t *testing.T) {
	cases := []struct {
		a   LatLng
		b   LatLng
		thr float64
	}{
		{LatLng{0, 0}, LatLng{0, 0}, 1e-9},
		{LatLng{0.1, -0.2}, LatLng{0.1 + 1e-12, -0.2 - 1e-12}, 1e-9},
		{LatLng{1, 2}, LatLng{1.0, 2.0 + 1e-8}, 1e-9},
		{LatLng{-1.5, 3.2}, LatLng{-1.5 + 1e-10, 3.2 + 1e-10}, 1e-9},
	}
	for _, tc := range cases {
		goVal := geoAlmostEqualThreshold(&tc.a, &tc.b, tc.thr)
		cVal := geoAlmostEqualThresholdC(tc.a, tc.b, tc.thr)
		if goVal != cVal {
			t.Fatalf("geoAlmostEqualThreshold parity mismatch: go=%v c=%v", goVal, cVal)
		}
	}
}
