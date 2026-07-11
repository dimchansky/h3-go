// Tests ported from testCellToBoundaryEdgeCases.c
package h3

import (
	"testing"
)

func Test_doublePrecisionVertex(t *testing.T) {
	t.Parallel()

	// The carefully constructed case here:
	// - A res 1 pentagon cell with distortion vertexes that change
	//   when we use a double instead of a float in _v2dIntersect
	// - One of the previous (float-based) distortion vertexes
	// This is the only case yet found where a point indexed to the
	// cell is shown to be incorrectly outside of the geo boundary
	// when we use the float version. Presumably more could be found.
	cell := h3Index(0x81083ffffffffff)
	point := LatLng{
		Lat: Deg(61.890838431),
		Lng: Deg(8.644221328),
	}

	var boundary CellBoundary
	err := cellToBoundary(cell, &boundary)
	if err != eSuccess {
		t.Fatalf("cellToBoundary failed: %v", err)
	}

	// Convert CellBoundary to GeoLoop (slice of LatLng)
	geoloop := make([]LatLng, boundary.NumVerts)
	copy(geoloop, boundary.Verts[:boundary.NumVerts])

	var bbox bbox
	bboxFromGeoLoop(geoloop, &bbox)

	var cell2 h3Index
	err = latLngToCell(&point, 1, &cell2)
	if err != eSuccess {
		t.Fatalf("latLngToCell failed: %v", err)
	}

	// Check whether the point is physically inside the geo boundary
	if cell2 == cell {
		if !pointInsideGeoLoop(geoloop, &bbox, &point) {
			t.Error("Boundary contains input point")
		}
	} else {
		if pointInsideGeoLoop(geoloop, &bbox, &point) {
			t.Error("Boundary does not contain input point")
		}
	}
}
