//go:build cgo

package h3

import (
	"testing"
)

// Test degenerate polygons (point, line) to verify C and Go behavior consistency
// This test was created to investigate why testPolygonToCells.c tests were failing.
// Finding: Both C and Go implementations return E_FAILED for degenerate polygons
// at the origin, but return E_SUCCESS for polygons with actual coordinates.
func Test_maxPolygonToCellsSize_degeneratePolygons_parity(t *testing.T) {
	t.Parallel()

	// Point polygon with single vertex at origin
	// Both C and Go return E_FAILED for this case
	pointAtOrigin := GeoPolygon{
		GeoLoop: []LatLng{{Lat: 0, Lng: 0}},
		Holes:   nil,
	}

	// Line polygon with two vertices from origin
	// Both C and Go return E_FAILED for this case
	linePolygon := GeoPolygon{
		GeoLoop: []LatLng{{Lat: 0, Lng: 0}, {Lat: 1, Lng: 0}},
		Holes:   nil,
	}

	// Line polygon with actual coordinates (not at origin)
	// Both C and Go return E_SUCCESS for this case
	linePolygonActual := GeoPolygon{
		GeoLoop: []LatLng{
			{Lat: 0.6595072188743, Lng: -2.1371053983433},
			{Lat: 0.6591482046471, Lng: -2.1373141048153},
		},
		Holes: nil,
	}

	testCases := []struct {
		name    string
		polygon *GeoPolygon
		res     int32
		flags   uint32
	}{
		{
			name:    "point polygon at origin",
			polygon: &pointAtOrigin,
			res:     9,
			flags:   0,
		},
		{
			name:    "line polygon (0,0) to (1,0)",
			polygon: &linePolygon,
			res:     9,
			flags:   0,
		},
		{
			name:    "line polygon with actual coords",
			polygon: &linePolygonActual,
			res:     9,
			flags:   0,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// Test Go implementation
			var goNumHexagons int64
			goErr := maxPolygonToCellsSize(tc.polygon, tc.res, tc.flags, &goNumHexagons)

			// Test C implementation via oracle
			var cNumHexagons int64
			cErr := maxPolygonToCellsSizeC(tc.polygon, tc.res, tc.flags, &cNumHexagons)

			// Check if behaviors match
			if goErr != cErr {
				t.Errorf("BEHAVIOR DIFFERENCE: Go returned %v (numHexagons=%d), C returned %v (numHexagons=%d)",
					goErr, goNumHexagons, cErr, cNumHexagons)
			} else if goErr == E_SUCCESS && goNumHexagons != cNumHexagons {
				t.Errorf("Size mismatch: Go returned %d, C returned %d", goNumHexagons, cNumHexagons)
			} else {
				// Success - behaviors match
				if goErr == E_FAILED {
					t.Logf("Both correctly reject degenerate polygon with E_FAILED")
				} else {
					t.Logf("Both accept polygon: numHexagons=%d", goNumHexagons)
				}
			}
		})
	}
}