//go:build cgo && c2go

package h3

import (
	"math"
	"testing"
)

func Test_validatePolygonFlags_ParityWithC(t *testing.T) {
	cases := []uint32{0, 1, 2, 3, 4, 15, 16, 31, 255}
	for _, f := range cases {
		goVal := validatePolygonFlags(f)
		cVal := validatePolygonFlagsC(f)
		if (goVal == 0) != (cVal == 0) {
			t.Fatalf("validatePolygonFlags mismatch: flags=%d go=%d c=%d", f, goVal, cVal)
		}
	}
}

func Test_lineCrossesLine_ParityWithC(t *testing.T) {
	a1 := LatLng{Lat: 0, Lng: 0}
	a2 := LatLng{Lat: 1, Lng: 1}
	b1 := LatLng{Lat: 0, Lng: 1}
	b2 := LatLng{Lat: 1, Lng: 0}
	if lineCrossesLine(&a1, &a2, &b1, &b2) != lineCrossesLineC(a1, a2, b1, b2) {
		t.Fatalf("lineCrossesLine mismatch (cross)")
	}
	b1 = LatLng{Lat: 0, Lng: 2}
	b2 = LatLng{Lat: 1, Lng: 2}
	if lineCrossesLine(&a1, &a2, &b1, &b2) != lineCrossesLineC(a1, a2, b1, b2) {
		t.Fatalf("lineCrossesLine mismatch (no cross)")
	}
}

// NaN/Inf coordinates make the intersection parameter NaN; C's final
// `return (test >= 0 && test <= 1)` is then false, and the Go port must
// agree (issue #3 investigation finding).
func Test_lineCrossesLine_nonFinite_ParityWithC(t *testing.T) {
	nan := Angle(math.NaN())
	posInf := Angle(math.Inf(1))
	negInf := Angle(math.Inf(-1))
	cases := []struct {
		name           string
		a1, a2, b1, b2 LatLng
	}{
		{"NaN lat on a1", LatLng{Lat: nan, Lng: 0}, LatLng{Lat: 0.1, Lng: 0.1}, LatLng{Lat: 0, Lng: 0.05}, LatLng{Lat: 0.2, Lng: 0.05}},
		{"+Inf lat on a1", LatLng{Lat: posInf, Lng: 0}, LatLng{Lat: 0.1, Lng: 0.1}, LatLng{Lat: 0, Lng: 0.05}, LatLng{Lat: 0.2, Lng: 0.05}},
		{"NaN lng on a1", LatLng{Lat: 0.2, Lng: nan}, LatLng{Lat: 0.1, Lng: 0.1}, LatLng{Lat: 0, Lng: 0.05}, LatLng{Lat: 0.2, Lng: 0.05}},
		{"-Inf lng on b1", LatLng{Lat: 0, Lng: 0}, LatLng{Lat: 0.2, Lng: 0.2}, LatLng{Lat: 0.1, Lng: negInf}, LatLng{Lat: 0.1, Lng: 0.1}},
		{"huge-magnitude lats", LatLng{Lat: Angle(3.6864520893115203e+267), Lng: 0}, LatLng{Lat: Angle(-4.782802678075986e+287), Lng: 0.1}, LatLng{Lat: 0, Lng: 0.05}, LatLng{Lat: 0.2, Lng: 0.05}},
	}
	for _, tc := range cases {
		goVal := lineCrossesLine(&tc.a1, &tc.a2, &tc.b1, &tc.b2)
		cVal := lineCrossesLineC(tc.a1, tc.a2, tc.b1, tc.b2)
		if goVal != cVal {
			t.Errorf("%s: lineCrossesLine mismatch: go=%v c=%v", tc.name, goVal, cVal)
		}
	}
}
