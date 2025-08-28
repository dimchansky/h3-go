//go:build cgo

package h3

import (
	"math"
	"testing"
)

func Test_greatCircleDistanceRads_ParityWithC(t *testing.T) {
	cases := []struct {
		a, b LatLng
	}{
		{LatLng{0, 0}, LatLng{0, 0}},
		{LatLng{0.1, -0.2}, LatLng{-0.3, 0.4}},
		{LatLng{1.0, 2.0}, LatLng{-1.0, -2.0}},
		{LatLng{0.5, 0.5}, LatLng{-0.5, -0.5}},
	}
	for _, tc := range cases {
		goVal := greatCircleDistanceRads(&tc.a, &tc.b)
		cVal := greatCircleDistanceRadsC(tc.a, tc.b)
		if math.Abs(goVal-cVal) > 1e-15 {
			t.Fatalf("greatCircleDistanceRads mismatch: go=%.17g c=%.17g", goVal, cVal)
		}
	}
}
