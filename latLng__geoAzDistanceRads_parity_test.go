//go:build cgo

package h3

import (
	"math"
	"testing"
)

func Test_geoAzDistanceRads_ParityWithC(t *testing.T) {
	cases := []struct {
		p1 LatLng
		az float64
		d  float64
	}{
		{LatLng{0, 0}, 0, 0},
		{LatLng{0.1, -0.2}, 0, 0.5},       // due north
		{LatLng{0.1, -0.2}, math.Pi, 0.5}, // due south
		{LatLng{1.0, 2.0}, 1.2345, 0.789}, // general case
		{LatLng{-0.5, 0.5}, 2.5, 1.0},
	}
	for _, tc := range cases {
		var goP2 LatLng
		_geoAzDistanceRads(&tc.p1, tc.az, tc.d, &goP2)
		cP2 := _geoAzDistanceRadsC(tc.p1, tc.az, tc.d)
		if math.Abs(goP2.Lat.Rad()-cP2.Lat.Rad()) > 1e-15 || math.Abs(goP2.Lng.Rad()-cP2.Lng.Rad()) > 1e-15 {
			t.Fatalf("_geoAzDistanceRads mismatch: go=(%.17g,%.17g) c=(%.17g,%.17g)", goP2.Lat, goP2.Lng, cP2.Lat, cP2.Lng)
		}
	}
}
