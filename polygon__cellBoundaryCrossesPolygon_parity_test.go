//go:build cgo && c2go

package h3

import "testing"

func Test_cellBoundaryCrossesPolygon_ParityWithC(t *testing.T) {
	outer := GeoLoop{{Lat: 0, Lng: 0}, {Lat: 0, Lng: 2}, {Lat: 2, Lng: 2}, {Lat: 2, Lng: 0}}
	hole := GeoLoop{{Lat: 0.5, Lng: 0.5}, {Lat: 0.5, Lng: 1.5}, {Lat: 1.5, Lng: 1.5}, {Lat: 1.5, Lng: 0.5}}
	poly := GeoPolygon{GeoLoop: outer, Holes: []GeoLoop{hole}}
	bboxes := make([]bbox, 1+len(poly.Holes))
	bboxFromGeoLoop(poly.GeoLoop, &bboxes[0])
	for i := range poly.Holes {
		bboxFromGeoLoop(poly.Holes[i], &bboxes[i+1])
	}
	// Crossing outer loop
	boundary := CellBoundary{numVerts: 2, verts: [MaxCellBoundaryVerts]LatLng{{Lat: -1, Lng: 1}, {Lat: 3, Lng: 1}}}
	var boundaryBBox bbox
	bboxFromGeoLoop(boundary.verts[:boundary.numVerts], &boundaryBBox)
	goVal := cellBoundaryCrossesPolygon(poly, bboxes, &boundary, &boundaryBBox)
	cVal := cellBoundaryCrossesPolygonC(poly, bboxes, boundary, boundaryBBox)
	if goVal != cVal {
		t.Fatalf("cellBoundaryCrossesPolygon mismatch (outer)")
	}
	// Crossing hole only (inside outer but crossing hole)
	boundary = CellBoundary{numVerts: 2, verts: [MaxCellBoundaryVerts]LatLng{{Lat: 0.75, Lng: 0.4}, {Lat: 0.75, Lng: 1.6}}}
	bboxFromGeoLoop(boundary.verts[:boundary.numVerts], &boundaryBBox)
	goVal = cellBoundaryCrossesPolygon(poly, bboxes, &boundary, &boundaryBBox)
	cVal = cellBoundaryCrossesPolygonC(poly, bboxes, boundary, boundaryBBox)
	if goVal != cVal {
		t.Fatalf("cellBoundaryCrossesPolygon mismatch (hole)")
	}
}
