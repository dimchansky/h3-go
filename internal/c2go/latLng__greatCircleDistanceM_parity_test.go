//go:build c2go

package c2go

import (
    "math"
    "testing"
)

func Test_greatCircleDistanceM_ParityWithC(t *testing.T) {
    cases := []struct{
        a, b LatLng
    }{
        {LatLng{0,0}, LatLng{0,0}},
        {LatLng{0.1, -0.2}, LatLng{-0.3, 0.4}},
        {LatLng{1.0, 2.0}, LatLng{-1.0, -2.0}},
        {LatLng{0.5, 0.5}, LatLng{-0.5, -0.5}},
    }
    for _, tc := range cases {
        goVal := greatCircleDistanceM(tc.a, tc.b)
        cVal := greatCircleDistanceMC(tc.a, tc.b)
        if math.Abs(goVal-cVal) > 1e-8 {
            t.Fatalf("greatCircleDistanceM mismatch: go=%.17g c=%.17g", goVal, cVal)
        }
    }
}
