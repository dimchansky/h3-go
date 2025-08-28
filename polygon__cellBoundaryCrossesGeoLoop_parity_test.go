//go:build cgo

package h3

import "testing"

func Test_cellBoundaryCrossesGeoLoop_ParityWithC(t *testing.T) {
	loop := GeoLoop{{Lat: 0, Lng: 0}, {Lat: 0, Lng: 2}, {Lat: 2, Lng: 2}, {Lat: 2, Lng: 0}}
	var loopBBox BBox
	bboxFromGeoLoop(loop, &loopBBox)
	// Boundary crossing the square
	boundary := CellBoundary{NumVerts: 2, Verts: []LatLng{{Lat: -1, Lng: 1}, {Lat: 3, Lng: 1}}}
	var boundaryBBox BBox
	bboxFromGeoLoop(boundary.Verts, &boundaryBBox)
	goVal := cellBoundaryCrossesGeoLoop(loop, &loopBBox, &boundary, &boundaryBBox)
	cVal := cellBoundaryCrossesGeoLoopC(loop, loopBBox, boundary, boundaryBBox)
	if goVal != cVal {
		t.Fatalf("cellBoundaryCrossesGeoLoop mismatch (crossing)")
	}
	// Boundary outside
	boundary = CellBoundary{NumVerts: 2, Verts: []LatLng{{Lat: -1, Lng: -1}, {Lat: -2, Lng: -2}}}
	bboxFromGeoLoop(boundary.Verts, &boundaryBBox)
	goVal = cellBoundaryCrossesGeoLoop(loop, &loopBBox, &boundary, &boundaryBBox)
	cVal = cellBoundaryCrossesGeoLoopC(loop, loopBBox, boundary, boundaryBBox)
	if goVal != cVal {
		t.Fatalf("cellBoundaryCrossesGeoLoop mismatch (outside)")
	}
}
