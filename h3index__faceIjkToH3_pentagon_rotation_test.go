//go:build cgo && c2go

package h3

import (
	"testing"
)

// Test_h3index_faceIjkToH3_PentagonRotation_ParityWithC tests a specific case
// that exposed a bug in pentagon rotation logic. The issue was using
// _isBaseCellPolarPentagon instead of _baseCellIsCwOffset for determining
// rotation direction, and incorrect loop structure for pentagon rotations.
func Test_h3index_faceIjkToH3_PentagonRotation_ParityWithC(t *testing.T) {
	// This specific case (0.0, 0.0) at resolution 15 revealed the pentagon
	// rotation bug where Go returned 0x8f745c4dbb49490 but C returned 0x8f754e64992d6d8
	testCases := []struct {
		lat, lng float64
		res      int32
		desc     string
	}{
		{0.0, 0.0, 15, "origin at high resolution"},
		{0.0, 0.0, 10, "origin at medium resolution"},
		{0.0, 0.0, 5, "origin at low resolution"},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			g := &LatLng{Lat: Rad(tc.lat), Lng: Rad(tc.lng)}

			var goOut H3Index
			var cOut H3Index

			goErr := latLngToCell(g, tc.res, &goOut)
			cErr := latLngToCellC(g, tc.res, &cOut)

			if uint32(goErr) != cErr {
				t.Fatalf("latLngToCell error mismatch for %s: go=%d c=%d",
					tc.desc, uint32(goErr), cErr)
			}

			if goOut != cOut {
				t.Fatalf("latLngToCell result mismatch for %s: go=%x c=%x",
					tc.desc, uint64(goOut), uint64(cOut))
			}
		})
	}
}
