//go:build cgo && c2go

package h3

import "testing"

func Test_cellBoundaryCrossesGeoLoop_ParityWithC(t *testing.T) {
	loop := GeoLoop{{Lat: 0, Lng: 0}, {Lat: 0, Lng: 2}, {Lat: 2, Lng: 2}, {Lat: 2, Lng: 0}}
	var loopBBox bbox
	bboxFromGeoLoop(loop, &loopBBox)
	// Boundary crossing the square
	boundary := CellBoundary{numVerts: 2, verts: [MaxCellBoundaryVerts]LatLng{{Lat: -1, Lng: 1}, {Lat: 3, Lng: 1}}}
	var boundaryBBox bbox
	bboxFromGeoLoop(boundary.verts[:boundary.numVerts], &boundaryBBox)
	goVal := cellBoundaryCrossesGeoLoop(loop, &loopBBox, &boundary, &boundaryBBox)
	cVal := cellBoundaryCrossesGeoLoopC(loop, loopBBox, boundary, boundaryBBox)
	if goVal != cVal {
		t.Fatalf("cellBoundaryCrossesGeoLoop mismatch (crossing)")
	}
	// Boundary outside
	boundary = CellBoundary{numVerts: 2, verts: [MaxCellBoundaryVerts]LatLng{{Lat: -1, Lng: -1}, {Lat: -2, Lng: -2}}}
	bboxFromGeoLoop(boundary.verts[:boundary.numVerts], &boundaryBBox)
	goVal = cellBoundaryCrossesGeoLoop(loop, &loopBBox, &boundary, &boundaryBBox)
	cVal = cellBoundaryCrossesGeoLoopC(loop, loopBBox, boundary, boundaryBBox)
	if goVal != cVal {
		t.Fatalf("cellBoundaryCrossesGeoLoop mismatch (outside)")
	}
}
