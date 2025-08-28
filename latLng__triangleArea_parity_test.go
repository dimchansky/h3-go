//go:build cgo

package h3

import (
	"math"
	"testing"
)

func Test_triangleArea_ParityWithC(t *testing.T) {
	tri := [3]LatLng{
		{Lat: 0.1, Lng: -0.2},
		{Lat: -0.3, Lng: 0.4},
		{Lat: 0.5, Lng: -0.6},
	}
	goVal := triangleArea(&tri[0], &tri[1], &tri[2])
	cVal := triangleAreaC(tri[0], tri[1], tri[2])
	if math.Abs(goVal-cVal) > 1e-15 {
		t.Fatalf("triangleArea mismatch: go=%.17g c=%.17g", goVal, cVal)
	}
}
