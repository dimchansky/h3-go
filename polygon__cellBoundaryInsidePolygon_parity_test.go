//go:build cgo && c2go

package h3

import "testing"

func Test_cellBoundaryInsidePolygon_ParityWithC(t *testing.T) {
	outer := GeoLoop{{Lat: 0, Lng: 0}, {Lat: 0, Lng: 2}, {Lat: 2, Lng: 2}, {Lat: 2, Lng: 0}}
	hole := GeoLoop{{Lat: 0.5, Lng: 0.5}, {Lat: 0.5, Lng: 1.5}, {Lat: 1.5, Lng: 1.5}, {Lat: 1.5, Lng: 0.5}}
	poly := GeoPolygon{GeoLoop: outer, Holes: []GeoLoop{hole}}
	bboxes := make([]bbox, 1+len(poly.Holes))
	bboxFromGeoLoop(poly.GeoLoop, &bboxes[0])
	for i := range poly.Holes {
		bboxFromGeoLoop(poly.Holes[i], &bboxes[i+1])
	}
	// Boundary inside outer but crossing hole; not contained
	boundary := CellBoundary{NumVerts: 2, Verts: []LatLng{{Lat: 0.25, Lng: 0.25}, {Lat: 0.75, Lng: 0.75}}}
	var boundaryBBox bbox
	bboxFromGeoLoop(boundary.Verts, &boundaryBBox)
	goVal := cellBoundaryInsidePolygon(poly, bboxes, &boundary, &boundaryBBox)
	cVal := cellBoundaryInsidePolygonC(poly, bboxes, boundary, boundaryBBox)
	if goVal != cVal {
		t.Fatalf("cellBoundaryInsidePolygon mismatch (inside/cross hole)")
	}
	// Boundary fully inside hole should be false
	boundary = CellBoundary{NumVerts: 2, Verts: []LatLng{{Lat: 1.0, Lng: 0.6}, {Lat: 1.4, Lng: 1.4}}}
	bboxFromGeoLoop(boundary.Verts, &boundaryBBox)
	goVal = cellBoundaryInsidePolygon(poly, bboxes, &boundary, &boundaryBBox)
	cVal = cellBoundaryInsidePolygonC(poly, bboxes, boundary, boundaryBBox)
	if goVal != cVal {
		t.Fatalf("cellBoundaryInsidePolygon mismatch (inside hole)")
	}
}
