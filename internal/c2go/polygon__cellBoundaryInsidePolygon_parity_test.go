//go:build c2go

package c2go

import "testing"

func Test_cellBoundaryInsidePolygon_ParityWithC(t *testing.T) {
	outer := GeoLoop{{Lat: 0, Lng: 0}, {Lat: 0, Lng: 2}, {Lat: 2, Lng: 2}, {Lat: 2, Lng: 0}}
	hole := GeoLoop{{Lat: 0.5, Lng: 0.5}, {Lat: 0.5, Lng: 1.5}, {Lat: 1.5, Lng: 1.5}, {Lat: 1.5, Lng: 0.5}}
	poly := GeoPolygon{Geoloop: outer, Holes: []GeoLoop{hole}}
	bboxes := make([]BBox, 1+len(poly.Holes))
	bboxes[0] = bboxFromGeoLoop(poly.Geoloop)
	for i := range poly.Holes {
		bboxes[i+1] = bboxFromGeoLoop(poly.Holes[i])
	}
	// Boundary inside outer but crossing hole; not contained
	boundary := CellBoundary{NumVerts: 2, Verts: []LatLng{{Lat: 0.25, Lng: 0.25}, {Lat: 0.75, Lng: 0.75}}}
	boundaryBBox := bboxFromGeoLoop(boundary.Verts)
	goVal := cellBoundaryInsidePolygon(poly, bboxes, boundary, boundaryBBox)
	cVal := cellBoundaryInsidePolygonC(poly, bboxes, boundary, boundaryBBox)
	if goVal != cVal {
		t.Fatalf("cellBoundaryInsidePolygon mismatch (inside/cross hole)")
	}
	// Boundary fully inside hole should be false
	boundary = CellBoundary{NumVerts: 2, Verts: []LatLng{{Lat: 1.0, Lng: 0.6}, {Lat: 1.4, Lng: 1.4}}}
	boundaryBBox = bboxFromGeoLoop(boundary.Verts)
	goVal = cellBoundaryInsidePolygon(poly, bboxes, boundary, boundaryBBox)
	cVal = cellBoundaryInsidePolygonC(poly, bboxes, boundary, boundaryBBox)
	if goVal != cVal {
		t.Fatalf("cellBoundaryInsidePolygon mismatch (inside hole)")
	}
}
