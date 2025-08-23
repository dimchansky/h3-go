//go:build cgo

package c2go

import (
	"math"
	"testing"
)

func Test_greatCircleDistanceKm_ParityWithC(t *testing.T) {
	cases := []struct {
		a, b LatLng
	}{
		{LatLng{0, 0}, LatLng{0, 0}},
		{LatLng{0.1, -0.2}, LatLng{-0.3, 0.4}},
		{LatLng{1.0, 2.0}, LatLng{-1.0, -2.0}},
		{LatLng{0.5, 0.5}, LatLng{-0.5, -0.5}},
	}
	for _, tc := range cases {
		goVal := greatCircleDistanceKm(&tc.a, &tc.b)
		cVal := greatCircleDistanceKmC(tc.a, tc.b)
		if math.Abs(goVal-cVal) > 1e-11 {
			t.Fatalf("greatCircleDistanceKm mismatch: go=%.17g c=%.17g", goVal, cVal)
		}
	}
}
