//go:build c2go

package c2go

import "testing"

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
