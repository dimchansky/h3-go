//go:build c2go

package c2go

import "testing"

func Test_cellBoundaryCrossesPolygon_ParityWithC(t *testing.T) {
    outer := GeoLoop{{Lat: 0, Lng: 0}, {Lat: 0, Lng: 2}, {Lat: 2, Lng: 2}, {Lat: 2, Lng: 0}}
    hole := GeoLoop{{Lat: 0.5, Lng: 0.5}, {Lat: 0.5, Lng: 1.5}, {Lat: 1.5, Lng: 1.5}, {Lat: 1.5, Lng: 0.5}}
    poly := GeoPolygon{Geoloop: outer, Holes: []GeoLoop{hole}}
    bboxes := make([]BBox, 1+len(poly.Holes))
    bboxes[0] = bboxFromGeoLoop(poly.Geoloop)
    for i := range poly.Holes {
        bboxes[i+1] = bboxFromGeoLoop(poly.Holes[i])
    }
    // Crossing outer loop
    boundary := CellBoundary{NumVerts: 2, Verts: []LatLng{{Lat: -1, Lng: 1}, {Lat: 3, Lng: 1}}}
    boundaryBBox := bboxFromGeoLoop(boundary.Verts)
    goVal := cellBoundaryCrossesPolygon(poly, bboxes, boundary, boundaryBBox)
    cVal := cellBoundaryCrossesPolygonC(poly, bboxes, boundary, boundaryBBox)
    if goVal != cVal {
        t.Fatalf("cellBoundaryCrossesPolygon mismatch (outer)")
    }
    // Crossing hole only (inside outer but crossing hole)
    boundary = CellBoundary{NumVerts: 2, Verts: []LatLng{{Lat: 0.75, Lng: 0.4}, {Lat: 0.75, Lng: 1.6}}}
    boundaryBBox = bboxFromGeoLoop(boundary.Verts)
    goVal = cellBoundaryCrossesPolygon(poly, bboxes, boundary, boundaryBBox)
    cVal = cellBoundaryCrossesPolygonC(poly, bboxes, boundary, boundaryBBox)
    if goVal != cVal {
        t.Fatalf("cellBoundaryCrossesPolygon mismatch (hole)")
    }
}

